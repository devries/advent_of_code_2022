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

	// For state we will save the time elapsed, position of the current player, open valves,
	// and which player. We will run player A and then player B.
	stateholder := make(map[State]IsDone)
	// Maximum maximum gas flow
	ret := findMaxGas(0, "AA", Bitset(0), "A", valves, ds, stateholder)

	return ret
}

type State struct {
	time     int
	position string
	open     Bitset // Open valve set
	player   string
}

// Store if this state was done and what the value is
type IsDone struct {
	done bool
	val  int
}

func findMaxGas(time int, position string, open Bitset, player string, valves map[string]Valve, distances map[string]map[string]int, precalc map[State]IsDone) int {
	// Check to see if we've done this state
	s := State{time, position, open, player}
	pc := precalc[s]
	if pc.done {
		return pc.val
	}

	// Otherwise calculate it. If player=="A" do those turns and then calculate player "B".
	max := 0

	currentGas := valves[position].Flow * (limit - time) // This valve will release gas for the rest of the time
	for next, vnext := range valves {
		if vnext.Flow == 0 || open.Contains(vnext.Bit) || next == position {
			continue
		}

		deltat := distances[position][next] + 1
		var gas int
		if time+deltat >= limit {
			if player == "A" {
				// Let second player go as well
				gas = currentGas + findMaxGas(0, "AA", open, "B", valves, distances, precalc)
			}
		} else {
			nopen := open
			nopen.Add(vnext.Bit)
			gas = currentGas + findMaxGas(time+deltat, next, nopen, player, valves, distances, precalc)
		}

		if gas > max {
			max = gas
		}
	}

	// Also check if you have nowhere else to go
	gas := currentGas
	if player == "A" {
		// Let second player go
		gas += findMaxGas(0, "AA", open, "B", valves, distances, precalc)
	}
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
