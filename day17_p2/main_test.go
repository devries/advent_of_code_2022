package main

import (
	"strings"
	"testing"
)

var testInput = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer int64
	}{
		{testInput, 1514285714288},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %d, got %d", test.answer, result)
		}
	}
}
