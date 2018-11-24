package main

import (
	"fmt"

	"github.com/golovers/leveltable"
)

func main() {
	db, err := leveltable.New("leveltable.db", 1000, 1000)
	if err != nil {
		panic(err)
	}

	table1 := db.Table("table1")
	for i := 'A'; i <= 'E'; i++ {
		table1.Put([]byte(string(i)), []byte("table 1 data "+string(i)))
	}
	table2 := db.Table("table2")
	for i := 'F'; i < 'J'; i++ {
		table2.Put([]byte(string(i)), []byte("table 2 data "+string(i)))
	}
	table2.Put([]byte("my-prefix 123"), []byte("table2 with prefix key data 123"))

	for it := table1.NewIterator(); it.Next(); {
		fmt.Printf("%s\n", it.Value())
	}

	for it := table2.NewIterator(); it.Next(); {
		fmt.Printf("%s\n", it.Value())
	}

	for it := table2.NewPrefixIterator("my-prefix"); it.Next(); {
		fmt.Printf("%s\n", it.Value())
	}

	table1.Table("table-3-over-table-1").Put([]byte("test"), []byte("test"))

	v, _ := table1.Table("table-3-over-table-1").Get([]byte("test"))
	fmt.Println("table-3-over-table-1: ", string(v))

}
