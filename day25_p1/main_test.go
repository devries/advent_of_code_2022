package main

import (
	"strings"
	"testing"
)

var testInput = `1=-0-2
12111
2=0=
21
2=01
111
20012
112
1=-1=
1-12
12
1=
122`

func TestSolution(t *testing.T) {
	tests := []struct {
		input  string
		answer string
	}{
		{testInput, "2=-1=0"},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)

		result := solve(r)

		if result != test.answer {
			t.Errorf("Expected %s, got %s", test.answer, result)
		}
	}
}
