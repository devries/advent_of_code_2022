package main

import (
	"fmt"
	"io"
	"os"
	"sort"
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

	sums := []int{}
	total := 0
	for _, r := range lines {
		if len(r) != 0 {
			v, err := strconv.Atoi(r)
			utils.Check(err, "error converting string to number")
			total += v
		} else {
			sums = append(sums, total)
			total = 0
		}
	}
	sums = append(sums, total)

	sort.Ints(sums)

	total = 0
	for i := len(sums) - 3; i < len(sums); i++ {
		if utils.Verbose {
			fmt.Printf("%d: %d\n", i, sums[i])
		}
		total += sums[i]
	}

	return total
}
