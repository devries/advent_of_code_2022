package main

import (
	"strings"
	"testing"
)

var testInput = `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer int
	}{
		{testInput, 93},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %d, got %d", test.answer, result)
		}
	}
}
