package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		for i := range 100 {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for v := range ch {
			fmt.Println(v, "worker1")
		}
	}()

	for v := range ch {
		fmt.Println(v, "worker2")
	}
}
