package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// BEGIN demo OMIT

	valuec := make(chan int)
	errorc := make(chan error)

	go worker(valuec, errorc)

	for {
		select {
		case v := <-valuec:
			fmt.Printf("got value: %d\n", v)
		case err := <-errorc:
			fmt.Printf("got error: %v\n", err)
		}
	}

	// END demo OMIT

}

func worker(valuec chan int, errorc chan error) {

	for {
		select {
		case valuec <- rand.Intn(10):
		case errorc <- errors.New("error"):
		}
		time.Sleep(time.Second)
	}

}
