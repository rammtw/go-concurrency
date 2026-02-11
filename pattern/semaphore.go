package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			fmt.Printf("Task %d started\n", i)
			time.Sleep(2 * time.Second)
			fmt.Printf("Task %d done\n", i)
		}()
	}

	wg.Wait()
}
