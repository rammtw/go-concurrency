package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 10)

	// Отправляем 10 запросов в канал
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)
	fmt.Println("Запуск rate limiter (3 запроса в секунду):")
	rateLimiter(requests, 3)
}

func rateLimiter(requests <-chan int, limit int) {
	ticker := time.NewTicker(3 * time.Second)

	defer ticker.Stop()

	for req := range requests {
		<-ticker.C
		fmt.Printf("[%s] Обработан запрос #%d\n",
			time.Now().Format("15:04:05.000"), req)
	}
}
