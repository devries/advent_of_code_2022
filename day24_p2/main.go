package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day24.txt"

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	if utils.Verbose {
		initialize()
		r, _ := resize()
		interruptHandling(r, 0)
		clear()
		hideCursor()
	}

	r := solve(f)
	if utils.Verbose {
		showCursor()
		r, _ := resize()
		move(r, 0)
		cleanup()
	}

	fmt.Println(r)
}

func solve(r io.Reader) int {
	lines := utils.ReadLines(r)

	in := parseInput(lines)

	open := findMapStates(in)

	start := State{in.Start, 0, nil}
	p1 := partialSolve(open, start, in.End)
	p2 := partialSolve(open, p1, in.Start)
	p3 := partialSolve(open, p2, in.End)

	if utils.Verbose {
		path := []utils.Point{}
		s := &p3
		for s != nil {
			p := s.pos
			path = append(path, p)
			s = s.prev
		}

		// reverse path
		for i := 0; i < len(path)/2; i++ {
			path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
		}

		for i, p := range path {
			clear()
			move(0, 0)
			printMap(open[i%len(open)], in.XSize, in.YSize)
			move(in.YSize-1-p.Y, p.X)
			fmt.Printf("█")
			time.Sleep(100 * time.Millisecond)
		}
	}

	return p3.step
}

// Solve from one point to another starting at a time
// Return the final state of the system
func partialSolve(open []map[utils.Point]bool, start State, end utils.Point) State {
	configurations := len(open)
	options := []utils.Point{utils.South, utils.East, utils.West, utils.North, {X: 0, Y: 0}} // moves one can make

	seen := make(map[State]bool)
	queue := []State{start}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]

		next := open[(s.step+1)%configurations]

		for _, o := range options {
			pn := s.pos.Add(o)
			if pn == end {
				return State{pn, s.step + 1, &s}
			}

			// Seen positions will be those states that have the same step modulo
			// the total number of blizzard configurations
			seeState := State{pn, (s.step + 1) % configurations, nil}

			if next[pn] && !seen[seeState] {
				seen[seeState] = true
				queue = append(queue, State{pn, s.step + 1, &s})
			}
		}
	}
	return State{utils.Point{X: 0, Y: 0}, 0, nil}
}

type State struct {
	pos  utils.Point
	step int
	prev *State
}

func findMapStates(in Input) []map[utils.Point]bool {
	vx := int64(in.XSize - 2)
	vy := int64(in.YSize - 2)
	cycles := int(utils.Lcm(vx, vy)) // cycles until blizzards repeat

	b := getBlizzardList(in.Blizzards)

	ret := []map[utils.Point]bool{}
	for i := 0; i < cycles; i++ {
		m := make(map[utils.Point]bool)
		m[in.Start] = true
		m[in.End] = true
		bmap := getBlizzardPositions(b)

		for j := 1; j < in.YSize-1; j++ {
			for i := 1; i < in.XSize-1; i++ {
				p := utils.Point{X: i, Y: j}

				if !bmap[p] {
					m[p] = true
				}
			}
		}
		ret = append(ret, m)
		b = blizzardStep(b, in.XSize, in.YSize)
	}

	return ret
}

type Blizzard struct {
	pos utils.Point
	dir utils.Point
}

func getBlizzardList(b map[utils.Point]utils.Point) []Blizzard {
	ret := []Blizzard{}

	for pos, dir := range b {
		ret = append(ret, Blizzard{pos, dir})
	}

	return ret
}

func getBlizzardPositions(blist []Blizzard) map[utils.Point]bool {
	ret := make(map[utils.Point]bool)

	for _, b := range blist {
		ret[b.pos] = true
	}

	return ret
}

func blizzardStep(blist []Blizzard, xsize int, ysize int) []Blizzard {
	ret := make([]Blizzard, len(blist))

	for i, b := range blist {
		n := b.pos.Add(b.dir)
		switch n.Y {
		case 0:
			n.Y = ysize - 2
		case ysize - 1:
			n.Y = 1
		}
		switch n.X {
		case 0:
			n.X = xsize - 2
		case xsize - 1:
			n.X = 1
		}

		ret[i] = Blizzard{n, b.dir}
	}

	return ret
}

func printOpen(o map[utils.Point]bool, xsize int, ysize int) {
	for j := ysize - 1; j >= 0; j-- {
		for i := 0; i < xsize; i++ {
			p := utils.Point{X: i, Y: j}

			switch o[p] {
			case true:
				fmt.Printf(".")
			case false:
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func printMap(o map[utils.Point]bool, xsize int, ysize int) {
	for x := 0; x < xsize; x++ {
		switch x {
		case 1:
			fmt.Printf(" ")
		default:
			fmt.Printf("#")
		}
	}
	fmt.Printf("\n")

	for j := ysize - 2; j > 0; j-- {
		for i := 0; i < xsize; i++ {
			p := utils.Point{X: i, Y: j}

			switch p.X {
			case 0, xsize - 1:
				fmt.Printf("#")
			default:
				switch o[p] {
				case true:
					fmt.Printf(" ")
				case false:
					fmt.Printf("@")
				}
			}
		}
		fmt.Printf("\n")
	}

	for x := 0; x < xsize; x++ {
		switch x {
		case xsize - 2:
			fmt.Printf(" ")
		default:
			fmt.Printf("#")
		}
	}
}

type Input struct {
	XSize     int                         // size of valley along x
	YSize     int                         // size of valley along y
	Walls     map[utils.Point]bool        // Position of the walls
	Blizzards map[utils.Point]utils.Point // position and direction of all blizzards
	Start     utils.Point                 // starting point
	End       utils.Point                 // ending point
}

func parseInput(lines []string) Input {
	ret := Input{len([]rune(lines[0])), len(lines), make(map[utils.Point]bool), make(map[utils.Point]utils.Point), utils.Point{}, utils.Point{}}

	for j, ln := range lines {
		y := len(lines) - 1 - j
		for i, c := range []rune(ln) {
			p := utils.Point{X: i, Y: y}

			switch c {
			case '#':
				ret.Walls[p] = true
			case '>':
				ret.Blizzards[p] = utils.East
			case '^':
				ret.Blizzards[p] = utils.North
			case '<':
				ret.Blizzards[p] = utils.West
			case 'v':
				ret.Blizzards[p] = utils.South
			case '.':
				if j == 0 {
					ret.Start = p
				} else if y == 0 {
					ret.End = p
				}
			}
		}
	}

	return ret
}
