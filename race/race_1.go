package main

import (
	"fmt"
	"sync"
)

func main() {
	a := make([]int, 0)
	var m sync.Mutex // если убрать мьютекс, получим гонку запустив go run -race race_1.go
	wg := &sync.WaitGroup{}

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			worker(wg)

			m.Lock()
			a = append(a, i)
			m.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println(a)
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
}
