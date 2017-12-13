package main

import "fmt"

func main() {

	// BEGIN demo OMIT

	var a *int

	b := 42
	a = &b

	b = 60

	fmt.Printf("a: %d, b: %d\n", *a, b)

	// END demo OMIT

}
