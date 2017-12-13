package main

import "fmt"

func main() {

	// BEGIN demo OMIT

	s := "Hello World"
	fmt.Printf("Hello %s\n", s[6:])
	fmt.Printf("%s World\n", s[:5])

	l := []int{5, 10, 15, 20}
	fmt.Printf("Slice: %v\n", l[1:3])

	// END demo OMIT

}
