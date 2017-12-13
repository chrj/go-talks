package main

import "testing"

func TestSum(t *testing.T) {

	// BEGIN test OMIT

	cases := []struct {
		Numbers  []int
		Expected int
	}{
		{[]int{2, 15, 25}, 42},
		{[]int{-5, 5, 20, 28, 32}, 80},
		{[]int{-2, -8, -10, -18, -22}, -60},
	}

	for _, c := range cases {
		if s := Sum(c.Numbers...); s != c.Expected {
			t.Errorf("got: '%d', expected: '%d' (input: %v)", s, c.Expected, c.Numbers)
		}
	}

	// END test OMIT

}
