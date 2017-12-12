package main

import (
	"fmt"
	"time"
)

// START demo OMIT

func main() {

	c := make(chan string)
	go pinger(c)

	for msg := range c {
		fmt.Println(msg)
	}

}

func pinger(c chan string) {
	for {
		time.Sleep(time.Second)
		c <- fmt.Sprintf("Ping: %s", time.Now().Format(time.RFC3339))
	}
}

// END demo OMIT
