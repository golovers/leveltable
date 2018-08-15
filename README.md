# leveltable

leveltable is a wrapper of leveldb which supports table, batches and some other utilities...

# Usage:

Code:
```go
db, err := leveltable.NewLevelTableDB("leveltable.db", 1000, 10)
if err != nil {
	panic(err)
}

db.Table("table1").Put([]byte("key"), []byte("val"))
db.Table("table2").Put([]byte("key"), []byte("val")

```
