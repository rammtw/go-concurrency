package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func main() {
	done := make(chan struct{})
	var once sync.Once
	var wg sync.WaitGroup

	for i, v := range NewData(10) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-done:
				return
			default:
			}

			if !v {
				fmt.Printf("%d: exit\n", i)
				once.Do(func() { close(done) })
				return
			}
			fmt.Printf("%d: ok\n", i)
		}()
	}

	wg.Wait()

	select {
	case <-done:
		fmt.Println("error: got false")
	default:
		fmt.Println("all ok")
	}
}

func NewData(n int) []bool {
	var data []bool

	if n > 15 {
		n = 15
	}

	for range n {
		data = append(data, true)
	}

	i := rand.IntN(len(data))
	data[i] = false

	return data
}
