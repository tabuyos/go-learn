package main

import "fmt"

type identifier struct {
	field1 string
	field2 string
}

func main() {
	fmt.Println("ok")
	var id identifier
	id.field1 = "one"
	id.field2 = "two"
}
