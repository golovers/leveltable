# leveltable

leveltable is a wrapper of leveldb which supports table, batches and some other utilities...

The original code is from go-ethereum... I just copy it here with a little modification for easier to use for my purpose...please use the original one if you like.

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
