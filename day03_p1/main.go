package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day03.txt"

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) int64 {
	lines := utils.ReadLines(r)

	var pSum int64
	for _, ln := range lines {
		items := []rune(ln)
		cSize := len(items) / 2

		compA := items[:cSize]
		compB := items[cSize:]

		contents := make(map[rune]bool)
		inBoth := make(map[rune]bool)

		for _, c := range compA {
			contents[c] = true
		}

		for _, c := range compB {
			if contents[c] {
				inBoth[c] = true
			}
		}

		for k := range inBoth {
			pSum += priority(k)
		}
	}

	return pSum
}

func priority(c rune) int64 {
	if c >= 'a' && c <= 'z' {
		return int64(c - 'a' + 1)
	} else {
		return int64(c - 'A' + 27)
	}
}
