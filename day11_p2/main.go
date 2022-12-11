package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day11.txt"

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) int64 {
	lines := utils.ReadLines(r)

	monkeys := parse(lines)

	// find LCM
	lcm := int64(1)
	for _, m := range monkeys {
		lcm = utils.Lcm(m.Test, lcm)
	}

	// rounds
	for i := 0; i < 10000; i++ {
		for j := 0; j < len(monkeys); j++ {
			turn(monkeys, j, lcm)
		}
	}

	sort.Sort(ByCount(monkeys))
	var ret int64 = 1
	for j := len(monkeys) - 2; j < len(monkeys); j++ {
		ret *= monkeys[j].Counter
	}

	return ret
}

type Monkey struct {
	Items     []int64
	Operation func(int64, int64) int64
	Test      int64
	True      int
	False     int
	Counter   int64
}

func parse(lines []string) []Monkey {
	ret := []Monkey{}

	monkeyRe := regexp.MustCompile(`^Monkey\s(\d+)`)
	startingRe := regexp.MustCompile(`Starting\sitems:\s`)
	operationRe := regexp.MustCompile(`Operation:\snew\s=\s(\w+)\s([*+])\s(\w+)`)
	testRe := regexp.MustCompile(`Test:\sdivisible\sby\s(\d+)`)
	ifRe := regexp.MustCompile(`If\s(true|false):\sthrow\sto\smonkey\s(\d+)`)
	csnRe := regexp.MustCompile(`(\d+)`)

	var m Monkey
	first := true
	for _, ln := range lines {
		switch {
		case monkeyRe.MatchString(ln):
			if !first {
				ret = append(ret, m)
			} else {
				first = false
			}
			m = Monkey{}

		case startingRe.MatchString(ln):
			for _, v := range csnRe.FindAllString(ln, -1) {
				n, err := strconv.ParseInt(v, 10, 64)
				utils.Check(err, "unable to parse %s to integer", v)

				m.Items = append(m.Items, n)
			}

		case operationRe.MatchString(ln):
			matches := operationRe.FindStringSubmatch(ln)

			switch matches[2] {
			case "+":
				if matches[3] == "old" {
					m.Operation = func(x int64, lcm int64) int64 {
						return (x + x) % lcm
					}
				} else {
					n, err := strconv.ParseInt(matches[3], 10, 64)
					utils.Check(err, "unable to parse %s to integer", matches[3])
					m.Operation = func(x int64, lcm int64) int64 {
						return (x + n) % lcm
					}
				}
			case "*":
				if matches[3] == "old" {
					m.Operation = func(x int64, lcm int64) int64 {
						return (x * x) % lcm
					}
				} else {
					n, err := strconv.ParseInt(matches[3], 10, 64)
					utils.Check(err, "unable to parse %s to integer", matches[3])
					m.Operation = func(x int64, lcm int64) int64 {
						return (x * n) % lcm
					}
				}
			}

		case testRe.MatchString(ln):
			matches := testRe.FindStringSubmatch(ln)

			n, err := strconv.ParseInt(matches[1], 10, 64)
			utils.Check(err, "unable to parse %s to integer", matches[1])

			m.Test = n

		case ifRe.MatchString(ln):
			matches := ifRe.FindStringSubmatch(ln)

			n, err := strconv.Atoi(matches[2])
			utils.Check(err, "unable to parse %s to integer", matches[2])

			switch matches[1] {
			case "true":
				m.True = n
			case "false":
				m.False = n
			}
		}
	}

	ret = append(ret, m)

	return ret
}

func turn(monkeys []Monkey, n int, lcm int64) {
	m := monkeys[n]

	for {
		if len(m.Items) == 0 {
			break
		}

		i := m.Items[0]
		m.Items = m.Items[1:]
		m.Counter++

		// Perform operation
		i = m.Operation(i, lcm)
		// i = i / 3

		// Test
		if i%m.Test == 0 {
			monkeys[m.True].Items = append(monkeys[m.True].Items, i)
		} else {
			monkeys[m.False].Items = append(monkeys[m.False].Items, i)
		}
	}

	monkeys[n] = m
}

type ByCount []Monkey

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Counter < a[j].Counter }
