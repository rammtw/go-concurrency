package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var ErrTimedOut = errors.New("timed out")

func main() {
	in := make(chan int)
	out := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go func() {
		defer close(in)
		for i := range 100 {
			select {
			case in <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	now := time.Now()
	processParallel(ctx, in, out, 5)

	for val := range out {
		fmt.Printf("read %v \n", val)
	}
	fmt.Println(time.Since(now))

	fmt.Println("goroutines alive:", runtime.NumGoroutine())
}

func processData(ctx context.Context, v int) (int, error) {
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(rand.Intn(22)) * time.Second)
		close(ch)
	}()

	select {
	case <-ch:
		return v * 2, nil
	case <-ctx.Done():
		return 0, ErrTimedOut
	}
}

func processParallel(ctx context.Context, in <-chan int, out chan<- int, n int) {
	wg := sync.WaitGroup{}
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case v, ok := <-in:
					if !ok {
						return
					}

					val, err := processData(ctx, v)
					if errors.Is(err, ErrTimedOut) {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- val:
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}
