package leveltable

import (
	"github.com/sirupsen/logrus"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var OpenFileLimit = 64

type ldbDatabase struct {
	db *leveldb.DB // LevelDB instance
}

// New returns a LevelDB wrapped object.
func New(file string, cache int, handles int) (Database, error) {
	// Ensure we have some minimal caching and file guarantees
	if cache < 16 {
		cache = 16
	}
	if handles < 16 {
		handles = 16
	}
	logrus.Info("Allocated cache and file handles", "cache", cache, "handles", handles)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &ldbDatabase{
		db: db,
	}, nil
}

// Put puts the given key / value to the queue
func (db *ldbDatabase) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// Has return true if the key is present
func (db *ldbDatabase) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get returns the given key if it's present.
func (db *ldbDatabase) Get(key []byte) ([]byte, error) {
	// Retrieve the key and increment the miss counter if not found
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Delete deletes the key from the queue and database
func (db *ldbDatabase) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *ldbDatabase) NewIterator() iterator.Iterator {
	return db.db.NewIterator(nil, nil)
}

func (db *ldbDatabase) NewPrefixIterator(prefix string) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
}

func (db *ldbDatabase) Close() {
	err := db.db.Close()
	if err == nil {
		logrus.Info("Database closed")
	} else {
		logrus.Error("Failed to close database", "err", err)
	}
}

func (db *ldbDatabase) LDB() *leveldb.DB {
	return db.db
}

func (db *ldbDatabase) NewBatch() Batch {
	return &ldbBatch{db: db.db, b: new(leveldb.Batch)}
}

func (db *ldbDatabase) Table(name string) Database {
	return newTable(db, name)
}

func (db *ldbDatabase) NewTableBatch(name string) Batch {
	return newTableBatch(db, name)
}

type ldbBatch struct {
	db   *leveldb.DB
	b    *leveldb.Batch
	size int
}

func (b *ldbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(value)
	return nil
}

func (b *ldbBatch) Write() error {
	return b.db.Write(b.b, nil)
}

func (b *ldbBatch) ValueSize() int {
	return b.size
}

func (b *ldbBatch) Reset() {
	b.b.Reset()
	b.size = 0
}

type table struct {
	db     Database
	prefix string
}

// NewTable returns a Database object that prefixes all keys with a given
// string.
func newTable(db Database, prefix string) Database {
	return &table{
		db:     db,
		prefix: prefix,
	}
}

func (dt *table) Put(key []byte, value []byte) error {
	return dt.db.Put(append([]byte(dt.prefix), key...), value)
}

func (dt *table) Has(key []byte) (bool, error) {
	return dt.db.Has(append([]byte(dt.prefix), key...))
}

func (dt *table) Get(key []byte) ([]byte, error) {
	return dt.db.Get(append([]byte(dt.prefix), key...))
}

func (dt *table) Delete(key []byte) error {
	return dt.db.Delete(append([]byte(dt.prefix), key...))
}

func (dt *table) NewIterator() iterator.Iterator {
	return dt.db.NewPrefixIterator(dt.prefix)
}

func (dt *table) NewPrefixIterator(prefix string) iterator.Iterator {
	return dt.db.NewPrefixIterator(dt.prefix + prefix)
}

func (dt *table) Table(name string) Database {
	return dt.db.Table(dt.prefix + name)
}

func (dt *table) NewTableBatch(name string) Batch {
	return dt.db.NewTableBatch(dt.prefix + name)
}

func (dt *table) Close() {
	// Do nothing; don't close the underlying DB.
}

type tableBatch struct {
	batch  Batch
	prefix string
}

// NewTableBatch returns a Batch object which prefixes all keys with a given string.
func newTableBatch(db Database, prefix string) Batch {
	return &tableBatch{db.NewBatch(), prefix}
}

func (dt *table) NewBatch() Batch {
	return &tableBatch{dt.db.NewBatch(), dt.prefix}
}

func (tb *tableBatch) Put(key, value []byte) error {
	return tb.batch.Put(append([]byte(tb.prefix), key...), value)
}

func (tb *tableBatch) Write() error {
	return tb.batch.Write()
}

func (tb *tableBatch) ValueSize() int {
	return tb.batch.ValueSize()
}

func (tb *tableBatch) Reset() {
	tb.batch.Reset()
}
