package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day12.txt"

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) int {
	// To find the 'a' with the shortest path we just find the shortest path from the end
	// to the first a we encounter

	lines := utils.ReadLines(r)

	maze := parse(lines)
	seen := make(map[utils.Point]bool)

	start := SearchPoint{maze.End, 0}
	search := []SearchPoint{start}

	for len(search) > 0 {
		sp := search[0]
		search = search[1:]

		pos := sp.Pos
		steps := sp.Steps
		if seen[pos] {
			continue
		}

		seen[pos] = true
		if maze.Altitudes[pos] == 'a' {
			return steps
		}

		for _, d := range utils.Directions {
			p := pos.Add(d)

			if maze.Altitudes[p] == 0 {
				continue
			}

			if maze.Altitudes[pos]-1 <= maze.Altitudes[p] {
				search = append(search, SearchPoint{p, steps + 1})
			}
		}
	}

	return 0
}

type Maze struct {
	Altitudes map[utils.Point]rune
	Start     utils.Point
	End       utils.Point
}

type SearchPoint struct {
	Pos   utils.Point
	Steps int
}

func parse(lines []string) Maze {
	ret := Maze{make(map[utils.Point]rune), utils.Point{}, utils.Point{}}

	for y := 0; y < len(lines); y++ {
		ln := lines[len(lines)-y-1]
		for x, c := range []rune(ln) {
			switch c {
			case 'S':
				ret.Start = utils.Point{X: x, Y: y}
				ret.Altitudes[ret.Start] = 'a'
			case 'E':
				ret.End = utils.Point{X: x, Y: y}
				ret.Altitudes[ret.End] = 'z'
			default:
				ret.Altitudes[utils.Point{X: x, Y: y}] = c
			}
		}
	}

	return ret
}
