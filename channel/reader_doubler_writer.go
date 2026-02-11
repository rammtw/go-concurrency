package main

import (
	"fmt"
	"time"
)

func main() {
	a := writer()
	b := doubler(a)
	reader(b)
}

func reader(a <-chan int) {
	for v := range a {
		fmt.Println(v)
		time.Sleep(500 * time.Millisecond)
	}
}

func doubler(a <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range a {
			out <- v * 2
		}
		close(out)
	}()
	return out
}

func writer() <-chan int {
	out := make(chan int)

	go func() {
		for i := range 10 {
			out <- i + 1
		}
		close(out)
	}()
	return out
}
