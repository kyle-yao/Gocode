package main

import (
	"fmt"
	"time"
)

func gotask() {
	i := 0
	for {
		i++
		fmt.Printf("new Goroutine: i=%d", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go gotask()
	i := 0
	for {
		i++
		fmt.Printf("man goroutine:i=%d\n", i)
		time.Sleep(1 * time.Second)
	}
}
