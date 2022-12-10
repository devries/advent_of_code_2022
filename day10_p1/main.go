package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day10.txt"

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

	values, _ := findDeviations(lines)

	var initial int64 = 1

	var sum int64
	for i := 20; i <= 220; i += 40 {
		sum += int64(i) * (values[i-1] + initial)
	}

	return sum
}

func findDeviations(lines []string) ([]int64, int64) {
	values := []int64{}

	var deviation int64
	for _, ln := range lines {
		components := strings.Fields(ln)

		switch components[0] {
		case "noop":
			values = append(values, deviation)

		case "addx":
			delta, err := strconv.ParseInt(components[1], 10, 64)
			utils.Check(err, "unable to convert %s to integer", components[1])

			values = append(values, deviation, deviation)
			deviation += delta
		}
	}

	return values, deviation
}
