package main

import (
	"fmt"
	"io"
	"math"
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
	min := Point{math.MaxInt, math.MaxInt, math.MaxInt}
	max := Point{math.MinInt, math.MinInt, math.MinInt}
	for _, ln := range lines {
		nums := strings.Split(ln, ",")

		x, err := strconv.Atoi(nums[0])
		utils.Check(err, "Error converting %s", ln)
		if x < min.X {
			min.X = x
		}
		if x > max.X {
			max.X = x
		}

		y, err := strconv.Atoi(nums[1])
		utils.Check(err, "Error converting %s", ln)
		if y < min.Y {
			min.Y = y
		}
		if y > max.Y {
			max.Y = y
		}

		z, err := strconv.Atoi(nums[2])
		utils.Check(err, "Error converting %s", ln)
		if z < min.Z {
			min.Z = z
		}
		if z > max.Z {
			max.Z = z
		}

		grid[Point{x, y, z}] = true
	}

	water := make(map[Point]bool)
	// Expand like water
	start := Point{min.X - 1, min.Y - 1, min.Z - 1}
	queue := []Point{start}

	for len(queue) > 0 {
		w := queue[0]
		queue = queue[1:]

		if water[w] {
			// Already seen
			continue
		}

		water[w] = true
		for _, d := range directions {
			n := w.Add(d)

			if n.X > max.X+1 || n.X < min.X-1 {
				continue
			}
			if n.Y > max.Y+1 || n.Y < min.Y-1 {
				continue
			}
			if n.Z > max.Z+1 || n.Z < min.Z-1 {
				continue
			}

			if !grid[n] && !water[n] {
				queue = append(queue, n)
			}
		}
	}

	totalSides := 0
	for k := range grid {
		for _, d := range directions {
			if water[k.Add(d)] {
				totalSides++
			}
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

var directions = []Point{Xhat, Yhat, Zhat, {-1, 0, 0}, {0, -1, 0}, {0, 0, -1}}

func (a Point) Add(b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}
