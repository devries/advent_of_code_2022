package main

import (
	"strings"
	"testing"
)

var testInput = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer int
	}{
		{testInput, 13},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %d, got %d", test.answer, result)
		}
	}
}
