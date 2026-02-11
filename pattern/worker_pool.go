package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	numWorkers := 3
	numJobs := 20

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	wg := &sync.WaitGroup{}

	for i := range numWorkers {
		wg.Add(1)
		go worker(i, jobs, results, wg)
	}

	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println("result:", r)
	}
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		fmt.Printf("worker %d start to do job %d\n", id, j)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		fmt.Printf("worker %d finished job %d\n", id, j)

		results <- j * 2
	}
}
