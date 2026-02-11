package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool, 1)

	go func() {
		time.Sleep(3 * time.Second)
		done <- true
	}()

	<-done
	fmt.Println("program exited")
}
