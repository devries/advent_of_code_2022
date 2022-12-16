package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/devries/combs"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day16.txt"
const limit = 26

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

	// Find the max gas for every combination of valves

	// First fing the significant valves (that have flow > 0)
	sigvalves := []string{}
	for k, v := range valves {
		if v.Flow > 0 {
			sigvalves = append(sigvalves, k)
		}
	}

	// We define a state which includes time step, position, open valves, and valve set
	// assigned to the user. We will be splitting up the work so that one person
	// has a combination of valves and the other has the remaining
	stateholder := make(map[State]IsDone)
	// Maximum maximum gas flow
	maxMax := 0

	// The first user will have i valves to open
	for i := 1; i < len(sigvalves); i++ {
		// Find all combinations of i valves
		c := combs.Combinations(i, sigvalves)
		for combo := range c {
			currentValves := make(map[string]Valve)
			currentValves["AA"] = valves["AA"]
			valveSetA := Bitset(0)
			// Give player A currentValves and the used valveSet
			for _, name := range combo {
				currentValves[name] = valves[name]
				valveSetA.Add(valves[name].Bit)
			}
			maxA := findMaxGas(0, "AA", Bitset(0), valveSetA, currentValves, ds, stateholder)

			// The second user will take all the valves that remain
			valveSetB := Bitset(0)
			remainingValves := make(map[string]Valve)
			remainingValves["AA"] = valves["AA"] // Everyone starts at AA
			for _, name := range sigvalves {
				v := valves[name]
				if !valveSetA.Contains(v.Bit) && name != "AA" {
					remainingValves[name] = v
					valveSetB.Add(v.Bit)
				}
			}

			maxB := findMaxGas(0, "AA", Bitset(0), valveSetB, remainingValves, ds, stateholder)

			// The total gas is what both users were able to accomplish
			sum := maxA + maxB
			if sum > maxMax {
				maxMax = sum
			}
		}
	}
	return maxMax
}

type State struct {
	time     int
	position string
	open     Bitset // Open valve set
	total    Bitset // User's valve set
}

// Store if this state was done and what the value is
type IsDone struct {
	done bool
	val  int
}

func findMaxGas(time int, position string, open Bitset, total Bitset, valves map[string]Valve, distances map[string]map[string]int, precalc map[State]IsDone) int {
	// Check to see if we've done this state
	s := State{time, position, open, total}
	pc := precalc[s]
	if pc.done {
		return pc.val
	}

	// Otherwise calculate it
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
			gas = gasReleased(open, valves, deltat) + findMaxGas(time+deltat, next, nopen, total, valves, distances, precalc)
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

	precalc[s] = IsDone{true, max}

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
