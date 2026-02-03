package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result := make(chan string, 1)
	read(ctx, result)

	select {
	case <-ctx.Done():
		fmt.Println("Ctx done")
	case res := <-result:
		fmt.Println("Результат:", res)
	default:
		fmt.Println("default")
	}
}

func read(ctx context.Context, result chan<- string) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	select {
	case <-time.After(4500 * time.Millisecond):
		result <- "готово"
	case <-ctx.Done():
		fmt.Println("Операция отменена:", ctx.Err())
		return
	}
}
