package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/devries/combs"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day06.txt"

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

	letters := []rune(lines[0])

	var retval int
outer:
	for i := 3; i < len(letters); i++ {
		marker := letters[i-3 : i+1]
		for pair := range combs.Combinations(2, marker) {
			if pair[0] == pair[1] {
				continue outer
			}
		}
		retval = i + 1
		break
	}

	return retval
}
