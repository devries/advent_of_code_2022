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
			// Lose round
			switch hands[0] {
			case "A":
				// rock v scissors
				points += 3
			case "B":
				// paper v rock
				points += 1
			case "C":
				// scissors v paper
				points += 2
			}
		case "Y":
			// Draw round pick same as other player
			points += 3
			switch hands[0] {
			case "A":
				points += 1
			case "B":
				points += 2
			case "C":
				points += 3
			}
		case "Z":
			// Win round
			points += 6
			switch hands[0] {
			case "A":
				// rock v paper
				points += 2
			case "B":
				// paper v scissors
				points += 3
			case "C":
				// scissors v rock
				points += 1
			}
		}
	}

	return points
}
