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

func solve(r io.Reader) int {
	lines := utils.ReadLines(r)

	groupSize := 3
	var pSum int
	for i := 0; i < len(lines); i += groupSize {
		groupLines := lines[i : i+groupSize]

		var itemSet uint64 = 0xfffffffffffff

		for _, ln := range groupLines {
			bm := bitmask([]rune(ln))
			itemSet &= bm
		}

		pSum += maskToPriority(itemSet)
	}

	return pSum
}

func priority(c rune) int {
	if c >= 'a' && c <= 'z' {
		return int(c - 'a' + 1)
	} else {
		return int(c - 'A' + 27)
	}
}

func bitmask(letters []rune) uint64 {
	var ret uint64

	for _, c := range letters {
		ret |= 1 << (priority(c) - 1)
	}

	return ret
}

func maskToPriority(v uint64) int {
	var ret int

	for l := v; l > 0; l >>= 1 {
		ret++
	}

	return ret
}
