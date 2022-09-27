package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestGcd(t *testing.T) {
	var tests = []struct {
		a      int64
		b      int64
		answer int64
	}{
		{1, 1, 1},
		{10, 5, 5},
		{12, 16, 4},
	}

	for _, test := range tests {
		result := Gcd(test.a, test.b)
		if result != test.answer {
			t.Errorf("For values %d and %d calculated %d, expected %d", test.a, test.b, result, test.answer)
		}
	}
}

func TestLcm(t *testing.T) {
	var tests = []struct {
		a      int64
		b      int64
		answer int64
	}{
		{1, 1, 1},
		{10, 5, 10},
		{12, 16, 48},
	}

	for _, test := range tests {
		result := Lcm(test.a, test.b)
		if result != test.answer {
			t.Errorf("For values %d and %d calculated %d, expected %d", test.a, test.b, result, test.answer)
		}
	}
}

func TestCountBits(t *testing.T) {
	var tests = []struct {
		n    uint32
		bits int
	}{
		{0b0, 0},
		{0b10, 1},
		{0b1011010110, 6},
	}

	for _, test := range tests {
		result := CountBits(test.n)
		if result != test.bits {
			t.Errorf("For bitfield %b and calculated %d bits, expected %d bits", test.n, result, test.bits)
		}
	}
}

func TestPoint(t *testing.T) {
	p := Point{0, 0}

	p2 := p.Add(North)
	p3 := p.Add(East)

	p4 := North.Add(South)

	if p2 != North {
		t.Errorf("0,0 + North should be North")
	}

	if p3 != East {
		t.Errorf("0,0 + East should be East")
	}

	if p4 != p {
		t.Errorf("North + South should be 0,0")
	}
}

func TestCheck(t *testing.T) {
	Check(nil, "test no error") // This should not panic

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The Check function did not panic")
		}
	}()

	err := fmt.Errorf("generic error")

	Check(err, "test error")
}

func TestReadLines(t *testing.T) {
	test := struct {
		Input  string
		Output []string
	}{"one\ntwo\nthree\n", []string{"one", "two", "three"}}

	r := strings.NewReader(test.Input)
	result := ReadLines(r)

	for i, v := range result {
		if v != test.Output[i] {
			t.Errorf("Expected %v, got %v", test.Output, result)
		}
	}
}
