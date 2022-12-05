package main

import (
	"strings"
	"testing"
)

var testInput = `    [D]    
[N] [C]    
[Z] [M] [P]
 1   2   3 

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer string
	}{
		{testInput, "CMZ"},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %s, got %s", test.answer, result)
		}
	}
}
