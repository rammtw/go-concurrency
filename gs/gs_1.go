package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // имитация долгого запроса
		w.Write([]byte("hello"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Канал для сигналов завершения (Ctrl+C, SIGTERM от Kubernetes и т.п.)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Println("server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Ждём сигнал
	<-stop
	log.Println("shutdown signal received")

	// Даём, например, 5 секунд на корректное завершение запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server exited gracefully")
}
