package main

import (
	"strings"
	"testing"
)

var testInput = `30373
25512
65332
33549
35390`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer int
	}{
		{testInput, 21},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %d, got %d", test.answer, result)
		}
	}
}
