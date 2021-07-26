// -*- coding: utf-8 -*-
package main

import (
	"fmt"
	"log"
	"os"
)

// init 方法每次都会在 main 方法之前运行
func init() {
	fmt.Println("execute init method.")
}

func main() {
	log.SetOutput(os.Stdout)
	fmt.Println("Hello, Tabuyos")
	fmt.Println("execute main method.")
}
