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

const inputfile = "../inputs/day14.txt"

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

	// Read in topology
	filled := make(map[utils.Point]bool)
	maxy := 0

	for _, ln := range lines {
		pts, b := parseLine(ln)
		if b > maxy {
			maxy = b
		}
		for _, p := range pts {
			filled[p] = true
		}
	}

	falldir := []utils.Point{{X: 0, Y: 1}, {X: -1, Y: 1}, {X: 1, Y: 1}}

	units := 0
outerLoop:
	for {
		// Release Sand
		units++
		s := utils.Point{X: 500, Y: 0}

		for s.Y < maxy {
			// Let it drop
			moved := false
			for _, dir := range falldir {
				d := s.Add(dir)
				if !filled[d] {
					s = d
					moved = true
					break
				}
			}

			if !moved {
				filled[s] = true
				continue outerLoop
			}
		}

		units--
		break
	}

	return units
}

// Parse a line of input and get the points and the bottom
func parseLine(ln string) ([]utils.Point, int) {
	corners := strings.Split(ln, " -> ")

	cornerPoints := []utils.Point{}

	for _, v := range corners {
		coords := strings.Split(v, ",")

		x, err := strconv.Atoi(coords[0])
		utils.Check(err, "Unable to convert %s to integer", coords[0])
		y, err := strconv.Atoi(coords[1])
		utils.Check(err, "Unable to convert %s to integer", coords[1])

		cornerPoints = append(cornerPoints, utils.Point{X: x, Y: y})
	}

	pts := []utils.Point{}
	bottom := 0

	for i := 0; i < len(cornerPoints)-1; i++ {
		s := cornerPoints[i]
		e := cornerPoints[i+1]

		if s.Y > bottom {
			bottom = s.Y
		}
		if e.Y > bottom {
			bottom = e.Y
		}

		d := e.Add(s.Scale(-1))

		inc := utils.Point{X: d.X / d.Manhattan(), Y: d.Y / d.Manhattan()}

		for p := s; p != e; p = p.Add(inc) {
			pts = append(pts, p)
		}
	}

	pts = append(pts, cornerPoints[len(cornerPoints)-1])

	return pts, bottom
}
