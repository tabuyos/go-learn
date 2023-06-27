package main

import (
	"fmt"
	"time"
)

func main2() {
	start := time.Now()
	time.Sleep(3 * time.Second)
	end := time.Now()
	delta := end.Sub(start)
	fmt.Printf("longCalculation took this amount of time: %s\n", delta)
}
