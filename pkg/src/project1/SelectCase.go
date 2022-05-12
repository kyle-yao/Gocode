package main

import "fmt"

func fibonacci(ch1, ch2 chan int) {
	x, y := 1, 1
	for {
		select {
		case ch1 <- x:
			sum := x + y
			x = y
			y = sum
		case <-ch2:
			fmt.Println("fibonacci is ended.")
			return
			//default:
			//	fmt.Println(1)
			//	return
		}
	}

}

func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		defer fmt.Println("Goroutine is ended")
		for i := 0; i < 6; i++ {
			fmt.Println(<-ch1)
		}
		ch2 <- 0
	}()

	fibonacci(ch1, ch2)
}
