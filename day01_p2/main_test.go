package main

import (
	"strings"
	"testing"
)

var testInput = `1000
2000
3000

4000

5000
6000

7000
8000
9000

10000`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer int
	}{
		{testInput, 45000},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %d, got %d", test.answer, result)
		}
	}
}
