package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {

	http.HandleFunc("/hello", func(rw http.ResponseWriter, req *http.Request) {
		io.WriteString(rw, "Hello World")
	})

	log.Fatal(http.Serve(autocert.NewListener("example.com"), nil))

}
