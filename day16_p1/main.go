package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day16.txt"
const limit = 30

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

	valves := parseInput(lines)

	ds := findDistances(valves)

	// Recursive search up to time limit
	ret := findMaxGas(0, "AA", Bitset(0), valves, ds)
	return ret
}

func findMaxGas(time int, position string, open Bitset, valves map[string]Valve, distances map[string]map[string]int) int {
	max := 0

	for next, vnext := range valves {
		if vnext.Flow == 0 || open.Contains(vnext.Bit) || next == position {
			continue
		}

		deltat := distances[position][next] + 1
		var gas int
		if time+deltat >= limit {
			gas = gasReleased(open, valves, limit-time)
		} else {
			nopen := open
			nopen.Add(vnext.Bit)
			gas = gasReleased(open, valves, deltat) + findMaxGas(time+deltat, next, nopen, valves, distances)
		}

		if gas > max {
			max = gas
		}
	}

	// Also check if you have nowhere else to go
	gas := gasReleased(open, valves, limit-time)
	if gas > max {
		max = gas
	}

	return max
}

type Valve struct {
	Flow    int
	Tunnels []string
	Bit     Bitset
}

func parseInput(lines []string) map[string]Valve {
	ret := make(map[string]Valve)

	bs := Bitset(1)
	for _, ln := range lines {
		components := strings.Fields(ln)

		name := components[1]
		var flow int
		fmt.Sscanf(components[4], "rate=%d;", &flow)

		destinations := make([]string, len(components[9:]))
		for i := 0; i < len(components[9:]); i++ {
			destinations[i] = strings.Trim(components[i+9], ",")
		}

		ret[name] = Valve{flow, destinations, bs}
		bs <<= 1
	}

	return ret
}

type DistanceState struct {
	ValveName string
	Steps     int
}

func findDistances(valves map[string]Valve) map[string]map[string]int {
	ret := make(map[string]map[string]int)

	for name := range valves {
		ret[name] = make(map[string]int)
	}

	for start := range valves {
		// Do BFS to find connection to all valves
		queue := []DistanceState{{start, 0}}
		seen := make(map[string]bool)

		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]

			if seen[p.ValveName] {
				continue
			}
			seen[p.ValveName] = true

			if p.ValveName != start {
				ret[start][p.ValveName] = p.Steps
			}

			for _, n := range valves[p.ValveName].Tunnels {
				queue = append(queue, DistanceState{n, p.Steps + 1})
			}
		}
	}

	return ret
}

type Bitset uint64

func (b *Bitset) Add(v Bitset) {
	*b |= v
}

func (b *Bitset) Remove(v Bitset) {
	*b &= ^v
}

func (b Bitset) Contains(v Bitset) bool {
	return (b & v) > 0
}

func gasReleased(open Bitset, valves map[string]Valve, increment int) int {
	sum := 0

	for _, v := range valves {
		if open.Contains(v.Bit) {
			sum += increment * v.Flow
		}
	}

	return sum
}
