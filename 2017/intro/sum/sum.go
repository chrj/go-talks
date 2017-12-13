package main

import "fmt"

func main() {
	fmt.Printf("%d", Sum(2, 15, 25))
}

// BEGIN sum OMIT

func Sum(numbers ...int) int {

	result := 0

	for _, number := range numbers {
		result += number
	}

	return result

}

// END sum OMIT
