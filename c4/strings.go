package main

import "fmt"

func main2() {
	s0 := `This is a raw string \n`
	fmt.Println(s0)
	str := "Beginning of the string " +
		"second part of the string"
	fmt.Println(str)
}
