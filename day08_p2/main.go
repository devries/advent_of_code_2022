package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day08.txt"

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

	f := parse(lines)

	maxScenic := 0
	for k, v := range f.Trees {
		// Count trees in each cardinal direction
		score := 1
		for _, d := range utils.Directions {
			n := 0
			for p := k.Add(d); p.X >= 0 && p.X < f.Width && p.Y >= 0 && p.Y < f.Height; p = p.Add(d) {
				n++
				if f.Trees[p] >= v {
					break
				}
			}
			score *= n
		}
		if score > maxScenic {
			maxScenic = score
		}
	}

	return maxScenic
}

func parse(lines []string) Forest {
	ret := Forest{make(map[utils.Point]int), len(lines), len([]rune(lines[0]))}

	for j, ln := range lines {
		y := len(lines) - j - 1
		characters := []rune(ln)
		for i, c := range characters {
			ret.Trees[utils.Point{X: i, Y: y}] = int(c - '0')
		}
	}

	return ret
}

type Forest struct {
	Trees  map[utils.Point]int
	Height int
	Width  int
}
