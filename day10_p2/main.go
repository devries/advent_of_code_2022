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
	if utils.Verbose {
		fmt.Println(r)
	}

	letters := utils.OCRLetters(r)
	fmt.Println(letters)
}

func solve(r io.Reader) string {
	lines := utils.ReadLines(r)

	values, _ := findDeviations(lines)

	var initial int64 = 1
	rows := make([][]rune, 6)

	for row := 0; row < 6; row++ {
		rows[row] = make([]rune, 40)
		for i := 0; i < 40; i++ {
			cycle := 40*row + i + 1
			offset := values[cycle-1] + initial

			if int64(i) >= offset-1 && int64(i) <= offset+1 {
				rows[row][i] = '#'
			} else {
				rows[row][i] = '.'
			}
		}
	}

	submessages := make([]string, 6)

	for i := 0; i < 6; i++ {
		submessages[i] = string(rows[i])
	}

	message := strings.Join(submessages, "\n")
	return message
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
