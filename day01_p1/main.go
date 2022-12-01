package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day01.txt"

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

	max := 0
	total := 0
	for _, r := range lines {
		if len(r) != 0 {
			v, err := strconv.Atoi(r)
			utils.Check(err, "error converting string to number")
			total += v
		} else {
			if total > max {
				max = total
			}
			total = 0
		}
	}
	if total > max {
		max = total
	}

	return max
}
