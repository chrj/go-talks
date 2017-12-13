package main

import "fmt"

func main() {

	// BEGIN demo OMIT

	m := make(map[string]string)
	m["a"] = "Hello"

	fmt.Printf("a: %s, b: %s\n", m["a"], m["b"])

	if value, found := m["b"]; !found {
		fmt.Println("'b' not found in map")
	} else {
		fmt.Printf("b: %s", value)
	}

	if _, found := m["b"]; !found {
		fmt.Println("'b' really not found in map")

	}
	// END demo OMIT

}
