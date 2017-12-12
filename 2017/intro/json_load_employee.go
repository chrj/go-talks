package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// BEGIN demo OMIT

type Employee struct {
	Name       string  `json:"name"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
}

func main() {

	data := []byte(`{"name":"Christian","department":"Development","salary":1000}`)

	e := Employee{}

	if err := json.Unmarshal(data, &e); err != nil {
		log.Fatalf("error while decoding data: %v", err)
	} else {
		fmt.Printf("Decoded: %s (salary: %.2f)\n", e.Name, e.Salary)
	}

}

// END demo OMIT
