package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		randomTimeWork()
		close(ch)
	}()

	select {
	case <-ch:
		fmt.Println("work done")
	case <-time.After(time.Second * 3):
		fmt.Println("time out")
	}
}

func randomTimeWork() {
	dur := time.Duration(rand.Intn(5)) * time.Second
	fmt.Println("random time", dur)
	time.Sleep(dur)
}
