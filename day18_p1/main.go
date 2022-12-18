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

const inputfile = "../inputs/day18.txt"

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
	grid := make(map[Point]bool)

	for _, ln := range lines {
		nums := strings.Split(ln, ",")

		x, err := strconv.Atoi(nums[0])
		utils.Check(err, "Error converting %s", ln)

		y, err := strconv.Atoi(nums[1])
		utils.Check(err, "Error converting %s", ln)

		z, err := strconv.Atoi(nums[2])
		utils.Check(err, "Error converting %s", ln)

		grid[Point{x, y, z}] = true
	}

	totalSides := len(grid) * 6

	for k := range grid {
		// Subtract two for any adjacency in + directions
		if grid[k.Add(Xhat)] {
			totalSides -= 2
		}
		if grid[k.Add(Yhat)] {
			totalSides -= 2
		}
		if grid[k.Add(Zhat)] {
			totalSides -= 2
		}
	}

	return totalSides
}

type Point struct {
	X int
	Y int
	Z int
}

var Xhat = Point{1, 0, 0}
var Yhat = Point{0, 1, 0}
var Zhat = Point{0, 0, 1}

func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}
