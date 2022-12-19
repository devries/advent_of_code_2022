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

const inputfile = "../inputs/day19.txt"

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

	sum := 0
	for i, ln := range lines {
		bp := parseBlueprint(ln)
		s := State{Resources{1, 0, 0, 0}, Resources{0, 0, 0, 0}}

		queue := []TimeState{{0, s}}
		max := buildGeode(queue, 24, bp)
		sum += (i + 1) * max
	}

	return sum
}

func buildGeode(queue []TimeState, timeLimit int, bp Blueprint) int {
	maxR := MaxResources(bp)
	max := 0

	for len(queue) > 0 {
		ts := queue[0]
		start := ts.S
		t := ts.Time
		queue = queue[1:]

		// Next robot to build selection
		if start.Robots.Obsidian > 0 {
			// Build Geode Robot
			next := start
			deltat := 0
			for {
				deltat++
				if t+deltat == timeLimit {
					next = AddResources(next)
					if next.Minerals.Geode > max {
						max = next.Minerals.Geode
					}
					break
				}

				if next, ok := buildGeodeRobot(next, bp); ok {
					queue = append(queue, TimeState{t + deltat, next})
					break
				}
				next = AddResources(next)
			}
		}

		if start.Robots.Clay > 0 && start.Robots.Obsidian < maxR.Obsidian {
			// Build an Obsidian Robot
			next := start
			deltat := 0
			for {
				deltat++
				if t+deltat == timeLimit {
					next = AddResources(next)
					if next.Minerals.Geode > max {
						max = next.Minerals.Geode
					}
					break
				}

				if next, ok := buildObsidianRobot(next, bp); ok {
					queue = append(queue, TimeState{t + deltat, next})
					break
				}
				next = AddResources(next)
			}
		}

		if start.Robots.Clay < maxR.Clay {
			// Build a clay robot
			next := start
			deltat := 0
			for {
				deltat++
				if t+deltat == timeLimit {
					next = AddResources(next)
					if next.Minerals.Geode > max {
						max = next.Minerals.Geode
					}
					break
				}

				if next, ok := buildClayRobot(next, bp); ok {
					queue = append(queue, TimeState{t + deltat, next})
					break
				}
				next = AddResources(next)
			}
		}

		if start.Robots.Ore < maxR.Ore {
			next := start
			deltat := 0
			for {
				deltat++
				if t+deltat == timeLimit {
					next = AddResources(next)
					if next.Minerals.Geode > max {
						max = next.Minerals.Geode
					}
					break
				}

				if next, ok := buildOreRobot(next, bp); ok {
					queue = append(queue, TimeState{t + deltat, next})
					break
				}
				next = AddResources(next)
			}
		}
	}

	return max
}
func MaxResources(bp Blueprint) Resources {
	// Find max of each resource required
	max := Resources{}

	for _, r := range []Resources{bp.Ore, bp.Clay, bp.Obsidian, bp.Geode} {
		if r.Ore > max.Ore {
			max.Ore = r.Ore
		}
		if r.Clay > max.Clay {
			max.Clay = r.Clay
		}
		if r.Obsidian > max.Obsidian {
			max.Obsidian = r.Obsidian
		}
	}

	return max
}

func AddResources(s State) State {
	s.Minerals.Ore += s.Robots.Ore
	s.Minerals.Clay += s.Robots.Clay
	s.Minerals.Obsidian += s.Robots.Obsidian
	s.Minerals.Geode += s.Robots.Geode

	return s
}

func buildOreRobot(s State, bp Blueprint) (State, bool) {
	if s.Minerals.Ore >= bp.Ore.Ore {
		s.Minerals.Ore -= bp.Ore.Ore
		s = AddResources(s)
		s.Robots.Ore++
		return s, true
	}

	return s, false
}

func buildClayRobot(s State, bp Blueprint) (State, bool) {
	if s.Minerals.Ore >= bp.Clay.Ore {
		s.Minerals.Ore -= bp.Clay.Ore
		s = AddResources(s)
		s.Robots.Clay++
		return s, true
	}

	return s, false
}

func buildObsidianRobot(s State, bp Blueprint) (State, bool) {
	if s.Minerals.Ore >= bp.Obsidian.Ore && s.Minerals.Clay >= bp.Obsidian.Clay {
		s.Minerals.Ore -= bp.Obsidian.Ore
		s.Minerals.Clay -= bp.Obsidian.Clay
		s = AddResources(s)
		s.Robots.Obsidian++
		return s, true
	}

	return s, false
}

func buildGeodeRobot(s State, bp Blueprint) (State, bool) {
	if s.Minerals.Ore >= bp.Geode.Ore && s.Minerals.Obsidian >= bp.Geode.Obsidian {
		s.Minerals.Ore -= bp.Geode.Ore
		s.Minerals.Obsidian -= bp.Geode.Obsidian
		s = AddResources(s)
		s.Robots.Geode++
		return s, true
	}

	return s, false
}

type Resources struct {
	Ore      int
	Clay     int
	Obsidian int
	Geode    int
}

type Blueprint struct {
	Ore      Resources
	Clay     Resources
	Obsidian Resources
	Geode    Resources
}

func parseBlueprint(ln string) Blueprint {
	components := strings.Fields(ln)

	// 6 = # ore for ore robot
	// 12 = # ore for clay robot
	// 18 = # ore for obsidian robot
	// 21 = # clay for obsidian robot
	// 27 = # clay for geode robot
	// 30 = # obsidian for geode robot

	var ret Blueprint
	var err error
	ret.Ore.Ore, err = strconv.Atoi(components[6])
	utils.Check(err, "Unable to get ore robot components from %s", ln)

	ret.Clay.Ore, err = strconv.Atoi(components[12])
	utils.Check(err, "Unable to get clay robot components from %s", ln)

	ret.Obsidian.Ore, err = strconv.Atoi(components[18])
	utils.Check(err, "Unable to get obsidian robot components from %s", ln)
	ret.Obsidian.Clay, err = strconv.Atoi(components[21])
	utils.Check(err, "Unable to get obsidian robot components from %s", ln)

	ret.Geode.Ore, err = strconv.Atoi(components[27])
	utils.Check(err, "Unable to get geode robot components from %s", ln)
	ret.Geode.Obsidian, err = strconv.Atoi(components[30])
	utils.Check(err, "Unable to get geode robot components from %s", ln)

	return ret
}

type State struct {
	Robots   Resources
	Minerals Resources
}

type TimeState struct {
	Time int
	S    State
}
