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

	e := Employee{
		Name:       "Christian",
		Department: "Development",
		Salary:     1000,
	}

	if encoded, err := json.Marshal(e); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", encoded)
	}

}

// END demo OMIT
