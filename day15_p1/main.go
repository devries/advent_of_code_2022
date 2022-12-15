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

	r := solve(f, 2000000)
	fmt.Println(r)
}

func solve(r io.Reader, row int64) int64 {
	lines := utils.ReadLines(r)

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

	ranges := []Range{}

	for _, s := range sensors {
		// distance from sensor to row
		d := abs(row - s.Position.Y)
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

	var count int64
	for _, r := range ranges {
		count += r.Count()
	}

	for _, b := range beacons {
		if b.Y == row {
			count--
		}
	}

	return count
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
