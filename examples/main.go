package main

import (
	"fmt"

	"github.com/golovers/leveltable"
)

func main() {
	db, err := leveltable.NewLDBDatabase("mydb.db", 1000, 1000)
	if err != nil {
		panic(err)
	}
	mytable := leveltable.NewTable(db, "mytable")

	mykey := []byte("my key")
	mytable.Put(mykey, []byte("my data"))

	data, err := mytable.Get(mykey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("data: %s\n", data)
}
