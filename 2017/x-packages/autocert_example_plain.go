package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", func(rw http.ResponseWriter, req *http.Request) {
		io.WriteString(rw, "Hello World")
	})

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("couldn't listen: %v", err)
	}

	log.Fatal(http.Serve(ln, nil))

}
