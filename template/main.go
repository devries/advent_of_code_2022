package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
)

const inputfile = "../inputs/day{{ printf `%02d` . }}.txt"

func main() {
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) int {
	lines := utils.ReadLines(r)

	_ = lines

	return 0
}
