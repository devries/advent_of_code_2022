package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day15.txt"

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f, 4000000)
	fmt.Println(r)
}

func solve(r io.Reader, maxCoord int64) int64 {
	lines := utils.ReadLines(r)

	possibleRange := Range{0, maxCoord}

	sensors := []Sensor{}
	beacons := []Pos{}
	for _, ln := range lines {
		s, b := parseLine(ln)
		d := b.Sub(s)
		sensors = append(sensors, Sensor{s, d.Metropolitan()})

		newbeacon := true
		for _, p := range beacons {
			if b == p {
				newbeacon = false
				break
			}
		}
		if newbeacon {
			beacons = append(beacons, b)
		}

	}

	var solution int64
	for y := int64(0); y <= maxCoord; y++ {
		ranges := []Range{}

		for _, s := range sensors {
			// distance from sensor to row
			d := abs(y - s.Position.Y)
			if d > s.Range {
				// sensor does not see row
				continue
			}

			subdistance := s.Range - d
			srange := Range{s.Position.X - subdistance, s.Position.X + subdistance}

			newranges := []Range{}
			for _, r := range ranges {
				nr, ok := overlap(r, srange)
				if ok {
					srange = nr
					continue
				}

				newranges = append(newranges, r)
			}
			ranges = newranges
			ranges = append(ranges, srange)
		}

		overlapCount := 0
		for _, r := range ranges {
			_, ok := overlap(possibleRange, r)
			if ok {
				overlapCount++
			}
		}

		if overlapCount == 2 {
			x := min(ranges[0].Max, ranges[1].Max) + 1
			solution = 4000000*x + y
			break
		}
	}

	return solution
}

type Pos struct {
	X int64
	Y int64
}

type Range struct {
	Min int64
	Max int64
}

func (r Range) Count() int64 {
	return r.Max - r.Min + 1
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func (a Pos) Metropolitan() int64 {
	return abs(a.X) + abs(a.Y)
}

func (a Pos) Sub(b Pos) Pos {
	return Pos{a.X - b.X, a.Y - b.Y}
}

func parseLine(ln string) (Pos, Pos) {
	s := Pos{}
	b := Pos{}

	fmt.Sscanf(ln, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &s.X, &s.Y, &b.X, &b.Y)

	return s, b
}

func max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Returns overlap and true if they overlap, false if they dont
func overlap(a Range, b Range) (Range, bool) {
	ret := Range{}

	if a.Min > b.Max+1 || a.Max < b.Min-1 {
		return ret, false
	}

	ret.Min = min(a.Min, b.Min)
	ret.Max = max(a.Max, b.Max)

	return ret, true
}

type Sensor struct {
	Position Pos
	Range    int64
}
