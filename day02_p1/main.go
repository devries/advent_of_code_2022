package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day02.txt"

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

	points := 0
	for _, ln := range lines {
		hands := strings.Fields(ln)

		switch hands[1] {
		case "X":
			points += 1
		case "Y":
			points += 2
		case "Z":
			points += 3
		}

		switch hands[0] {
		case "A":
			switch hands[1] {
			case "X":
				points += 3 // rock v rock
			case "Y":
				points += 6 // rock v paper
			}
		case "B":
			switch hands[1] {
			case "Y":
				points += 3 // paper v paper
			case "Z":
				points += 6 // paper v scissors
			}
		case "C":
			switch hands[1] {
			case "X":
				points += 6 // scissors v rock
			case "Z":
				points += 3 // scissors v scissors
			}
		}
	}

	return points
}
