package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day05.txt"

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) string {
	lines := utils.ReadLines(r)

	instRe := regexp.MustCompile(`^move\s(\d+)\sfrom\s(\d+)\sto\s(\d+)`)
	crateRe := regexp.MustCompile(`\[([A-Z])\]`)
	stackRe := regexp.MustCompile(`\d+`)

	var stacks [][]rune
	var stackLines []string

	for _, ln := range lines {
		switch {
		case instRe.MatchString(ln):
			// Instruction
			matches := instRe.FindStringSubmatch(ln)
			num, err := strconv.Atoi(matches[1])
			utils.Check(err, "Error converting string to int")
			start, err := strconv.Atoi(matches[2])
			utils.Check(err, "Error converting string to int")
			end, err := strconv.Atoi(matches[3])
			utils.Check(err, "Error converting string to int")

			s := stacks[start-1]
			e := stacks[end-1]

			e = append(e, s[len(s)-num:]...)
			s = s[:len(s)-num]

			stacks[start-1] = s
			stacks[end-1] = e

		case crateRe.MatchString(ln):
			// Crate line
			stackLines = append(stackLines, ln)

		case stackRe.MatchString(ln):
			// Stack numbers
			indecies := stackRe.FindAllStringSubmatchIndex(ln, -1)
			stacks = make([][]rune, len(indecies))

			for i := len(stackLines) - 1; i >= 0; i-- {
				sl := stackLines[i]
				for j, p := range indecies {
					crate := []rune(sl[p[0]:p[1]])
					if crate[0] != 32 {
						// no spaces
						stacks[j] = append(stacks[j], crate[0])
					}
				}
			}
			if utils.Verbose {
				printStacks(stacks)
			}
		}
	}

	solution := make([]rune, len(stacks))
	for i, s := range stacks {
		solution[i] = s[len(s)-1]
	}

	return string(solution)
}

func printStacks(stacks [][]rune) {
	for i, s := range stacks {
		fmt.Printf("%d: %s\n", i+1, string(s))
	}
}
