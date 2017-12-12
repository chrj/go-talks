package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(rw, "Hello World!\n")
	})

	log.Fatal(http.ListenAndServe(":8123", nil))

}
