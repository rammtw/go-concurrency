package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

func riskyOperation(id int) error {
	if id%4 == 0 {
		return fmt.Errorf("failed on %d", id)
	}
	fmt.Println("processed", id)
	return nil
}

func main() {
	var g errgroup.Group
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}

	for _, id := range items {
		g.Go(func() error {
			if err := riskyOperation(id); err != nil {
				// Оборачиваем ошибку контекстом
				return fmt.Errorf("item %d: %w", id, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("got error:", err)
	} else {
		fmt.Println("all ok")
	}
}
