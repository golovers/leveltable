# leveltable

leveltable is a wrapper of leveldb which supports table, batches and some other utilities...

The original code is from go-ethereum... I just copy it here with a little modification for easier to use for my purpose...please use the original one if you like.

# Usage:

Code:
```go

db, err := leveltable.NewLDBDatabase("leveltable.db", 1000, 1000)
if err != nil {
	panic(err)
}
table1 := leveltable.NewTable(db, "table1")
for i := 'A'; i <= 'E'; i++ {
	table1.Put([]byte(string(i)), []byte("table 1 data "+string(i)))
}

table2 := leveltable.NewTable(db, "table2")
for i := 'F'; i < 'J'; i++ {
	table2.Put([]byte(string(i)), []byte("table 2 data "+string(i)))
}
table2.Put([]byte("my-prefix 123"), []byte("table2 with prefix key data 123"))

fmt.Println("--------------table1-------------------------")
for it := table1.NewIterator(); it.Next(); {
	fmt.Printf("%s\n", it.Value())
}
fmt.Println("--------------table2-------------------------")
for it := table2.NewIterator(); it.Next(); {
	fmt.Printf("%s\n", it.Value())
}
fmt.Println("--------table2 with prefix key---------------")
for it := table2.NewPrefixIterator("my-prefix"); it.Next(); {
	fmt.Printf("%s\n", it.Value())
}
```
Output:
```
--------------table1-------------------------
table 1 data A
table 1 data B
table 1 data C
table 1 data D
table 1 data E
--------------table2-------------------------
table 2 data F
table 2 data G
table 2 data H
table 2 data I
table2 with prefix key data 123
--------table2 with prefix key---------------
table2 with prefix key data 123
```
