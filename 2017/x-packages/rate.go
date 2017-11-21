package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/time/rate"
)

func main() {

	l := rate.NewLimiter(2, 1)
	ctx := context.Background()

	for i := 0; i < 10; i++ {

		go func(i int) {

			for {
				l.Wait(ctx)
				log.Printf("working in %v", i)
			}

		}(i)

	}

	time.Sleep(10 * time.Second)

}
