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

	groupSize := 3
	var pSum int64
	for i := 0; i < len(lines); i += groupSize {
		groupLines := lines[i : i+groupSize]

		contents := make(map[rune]uint8)
		for g, ln := range groupLines {
			items := []rune(ln)

			groupBit := uint8(1 << g)
			for _, c := range items {
				contents[c] |= groupBit
			}
		}

		for k, v := range contents {
			if v == 7 {
				pSum += priority(k)
				break
			}
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
