package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day22.txt"

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

	ysize := len(lines) - 2
	directions := lines[len(lines)-1]

	xsize := 0
	grid := make(map[utils.Point]rune)
	for j := 0; j < ysize; j++ {
		for i, v := range []rune(lines[ysize-1-j]) {
			if v != ' ' {
				grid[utils.Point{X: i, Y: j}] = v
				if i > xsize {
					xsize = i
				}
			}
		}
	}
	xsize++

	startx := 0
	for i := 0; i < xsize; i++ {
		if grid[utils.Point{X: i, Y: ysize - 1}] != 0 {
			startx = i
			break
		}
	}

	corners := findCorners(grid, xsize, ysize)
	edgemap := make(map[Warp]Warp)
	for p, d := range corners {
		// Run along edges to find where map points meet up.
		r := getRunners(p, d)

		for {
			r0edge := r[0].pos.Add(r[0].point)
			r1edge := r[1].pos.Add(r[1].point)
			edgemap[Warp{r0edge, r[0].point.Scale(-1)}] = Warp{r1edge, r[1].point}
			edgemap[Warp{r1edge, r[1].point.Scale(-1)}] = Warp{r0edge, r[0].point}

			r[0].pos = r[0].pos.Add(r[0].dir)
			r[1].pos = r[1].pos.Add(r[1].dir)

			next0edge := r[0].pos.Add(r[0].point)
			next1edge := r[1].pos.Add(r[1].point)
			if grid[next0edge] == 0 && grid[next1edge] == 0 {
				break
			}

			if grid[next0edge] == 0 {
				r[0].pos = next0edge
				r[0].dir, r[0].point = r[0].point, r[0].dir.Scale(-1)
			}

			if grid[next1edge] == 0 {
				r[1].pos = next1edge
				r[1].dir, r[1].point = r[1].point, r[1].dir.Scale(-1)
			}
		}
	}

	pos := utils.Point{X: startx, Y: ysize - 1}
	dir := utils.Point{X: 1, Y: 0}

	if utils.Verbose {
		printGridandCorners(grid, xsize, ysize, corners)
	}
	movements := parseDirections(directions)

	for _, m := range movements {
		pos, dir = step(pos, dir, grid, m, edgemap)
	}

	row := ysize - pos.Y
	column := pos.X + 1
	dirState := 0
	switch dir {
	case utils.East:
		dirState = 0
	case utils.South:
		dirState = 1
	case utils.West:
		dirState = 2
	case utils.North:
		dirState = 3
	}

	if utils.Verbose {
		fmt.Println("Row:", ysize-pos.Y)
		fmt.Println("Column:", pos.X+1)
		fmt.Println("Direction:", dirState)
	}
	return 1000*row + 4*column + dirState
}

func printGrid(grid map[utils.Point]rune, xsize int, ysize int) {
	for j := ysize - 1; j >= 0; j-- {
		for i := 0; i < xsize; i++ {
			v := grid[utils.Point{X: i, Y: j}]
			switch v {
			case 0:
				fmt.Printf(" ")
			default:
				fmt.Printf("%c", v)
			}
		}
		fmt.Printf("\n")
	}
}

func printGridandCorners(grid map[utils.Point]rune, xsize int, ysize int, corners map[utils.Point][]utils.Point) {
	for j := ysize - 1; j >= 0; j-- {
		for i := 0; i < xsize; i++ {
			p := utils.Point{X: i, Y: j}
			v := grid[p]
			switch v {
			case 0:
				if _, ok := corners[p]; ok {
					fmt.Printf("X")
				} else {
					fmt.Printf(" ")
				}
			default:
				fmt.Printf("%c", v)
			}
		}
		fmt.Printf("\n")
	}
}

type Movement struct {
	Steps int
	Turn  rune
}

func parseDirections(s string) []Movement {
	characters := []rune(s)

	ret := []Movement{}
	start := 0
	for i, v := range characters {
		switch v {
		case 'R', 'L':
			sv := string(characters[start:i])
			n, err := strconv.Atoi(sv)
			utils.Check(err, "unable to convert %s to integer", sv)
			m := Movement{n, v}
			ret = append(ret, m)
			start = i + 1
		}
	}

	// Final motion is just steps
	sv := string(characters[start:])
	n, err := strconv.Atoi(sv)
	utils.Check(err, "unable to convert %s to integer", sv)
	m := Movement{n, 0}
	ret = append(ret, m)

	return ret
}

func step(pos utils.Point, dir utils.Point, grid map[utils.Point]rune, m Movement, edgemap map[Warp]Warp) (utils.Point, utils.Point) {
	for i := 0; i < m.Steps; i++ {
		np := pos.Add(dir)
		switch grid[np] {
		case '.':
			pos = np
		case 0:
			w := edgemap[Warp{pos, dir}]
			if grid[w.pos] == '.' {
				pos = w.pos
				dir = w.dir
			}
		}
	}

	switch m.Turn {
	case 'R':
		// turn right
		dir = dir.Right()
	case 'L':
		// turn left
		dir = dir.Left()
	}

	return pos, dir
}

func findCorners(grid map[utils.Point]rune, xsize int, ysize int) map[utils.Point][]utils.Point {
	corners := make(map[utils.Point][]utils.Point)

	for j := ysize - 1; j >= 0; j-- {
		for i := 0; i < xsize; i++ {
			pos := utils.Point{X: i, Y: j}
			v := grid[pos]
			if v == 0 {
				drs := []utils.Point{}
				for _, d := range utils.Directions {
					p := pos.Add(d)
					if grid[p] != 0 {
						drs = append(drs, d)
					}
				}
				if len(drs) == 2 {
					// Corner found
					corners[pos] = drs
				}
			}
		}
	}

	return corners
}

type Runner struct {
	pos   utils.Point
	dir   utils.Point
	point utils.Point
}

func getRunners(corner utils.Point, directions []utils.Point) []Runner {
	ret := make([]Runner, 2)

	for i := 0; i < 2; i++ {
		o := (i + 1) % 2

		ret[i] = Runner{corner, directions[o].Scale(-1), directions[i]}
	}

	return ret
}

type Warp struct {
	pos utils.Point
	dir utils.Point
}
