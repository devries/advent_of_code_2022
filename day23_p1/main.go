package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day23.txt"

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

	g := parseGrid(lines)

	d := []utils.Point{utils.North, utils.South, utils.West, utils.East}

	for i := 0; i < 10; i++ {
		g = step(g, d)
		d[0], d[1], d[2], d[3] = d[1], d[2], d[3], d[0]
	}
	if utils.Verbose {
		printGrid(g)
	}
	res := countGroundSpaces(g)
	return res
}

func parseGrid(lines []string) map[utils.Point]utils.Point {
	g := make(map[utils.Point]utils.Point)

	for j, ln := range lines {
		y := len(lines) - 1 - j
		for i, c := range []rune(ln) {
			if c == '#' {
				g[utils.Point{X: i, Y: y}] = utils.Point{X: i, Y: y}
			}
		}
	}

	return g
}

func findMinMax(g map[utils.Point]utils.Point) (utils.Point, utils.Point) {
	min := utils.Point{X: math.MaxInt, Y: math.MaxInt}
	max := utils.Point{X: math.MinInt, Y: math.MinInt}

	for p := range g {
		if p.X < min.X {
			min.X = p.X
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}

	return min, max
}

func step(g map[utils.Point]utils.Point, directions []utils.Point) map[utils.Point]utils.Point {
	ng := make(map[utils.Point]utils.Point) // keep track of previous position after move
	banned := make(map[utils.Point]bool)    // Points where no elf should move because of elf collision

	// find surrounding points
	surrounds := []utils.Point{{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 0}, {X: 1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1}}
	adjacencies := map[utils.Point][]utils.Point{
		utils.North: {{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1}},
		utils.South: {{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1}},
		utils.East:  {{X: 1, Y: -1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
		utils.West:  {{X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}},
	}

	for p := range g {
		// check if elf should stay
		neighbors := 0
		for _, n := range surrounds {
			c := p.Add(n)
			if _, ok := g[c]; ok {
				// elf has a neighbor
				neighbors++
			}
		}
		if neighbors == 0 {
			// no neighbors, elf stays put
			ng[p] = p
			continue // move to next elf
		}

		// now check each direction and see if the coast is clear
		placed := false
	searchLoop:
		for _, d := range directions {
			checkpoints := adjacencies[d]
			np := p.Add(d)
			for _, cp := range checkpoints {
				if _, ok := g[p.Add(cp)]; ok {
					// This direction has a nearby elf
					continue searchLoop
				}
			}
			// if you get here the direction is clear
			if banned[np] {
				// This position was determined blocked for elf motion
				placed = true
				ng[p] = p // do not move
				break
			}

			e, ok := ng[np]
			if ok {
				// This spot is a place where another elf wants to go
				banned[np] = true // ban for further checks
				delete(ng, np)    // Remove previous elf from that position
				ng[e] = e         // move other elf back
				ng[p] = p         // keep current elf in current position
				placed = true
				break
			}

			// If you make it this far, you can probably move in your next turn
			ng[np] = p
			placed = true
			break
		}
		if !placed {
			// Elf found no place to move
			ng[p] = p
		}
	}

	return ng
}

func printGrid(g map[utils.Point]utils.Point) {
	min, max := findMinMax(g)

	for j := max.Y; j >= min.Y; j-- {
		for i := min.X; i <= max.X; i++ {
			_, ok := g[utils.Point{X: i, Y: j}]
			switch ok {
			case true:
				fmt.Printf("#")
			default:
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

func countGroundSpaces(g map[utils.Point]utils.Point) int {
	min, max := findMinMax(g)
	count := 0

	for j := min.Y; j <= max.Y; j++ {
		for i := min.X; i <= max.X; i++ {
			if _, ok := g[utils.Point{X: i, Y: j}]; !ok {
				count++
			}
		}
	}

	return count
}
