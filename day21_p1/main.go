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

const inputfile = "../inputs/day21.txt"

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

	monkeys := make(map[string]Job)

	for _, ln := range lines {
		name, j := parseLine(ln)
		monkeys[name] = j
	}

	rpnList := []string{"root"}

	// Expand out
	keepGoing := true
	for keepGoing {
		// Do this until there are no more monkey names in stack
		keepGoing = false
		newList := []string{}
		// Move from bottom to top
		for _, m := range rpnList {
			if strings.ContainsAny(m, "0123456789+-/*") {
				newList = append(newList, m)
				continue
			}

			// Split the operation op and append this operation to the existing op stack
			keepGoing = true
			j := monkeys[m]
			switch j.Operation {
			case "":
				newList = append(newList, j.A)
			default:
				newList = append(newList, j.B, j.A, j.Operation)
			}
		}
		rpnList = newList
	}

	// Do calculation
	nStack := IntStack{}

	for _, v := range rpnList {
		switch v {
		case "+":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a + b)
		case "-":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a - b)
		case "/":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a / b)
		case "*":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a * b)
		default:
			n, err := strconv.ParseInt(v, 10, 64)
			utils.Check(err, "unable to convert %s to integer", v)
			nStack.Push(n)
		}
	}

	return nStack[0]
}

type Job struct {
	A         string
	B         string
	Operation string
}

func parseLine(ln string) (string, Job) {
	components := strings.Split(ln, ": ")

	jobParts := strings.Fields(components[1])

	switch len(jobParts) {
	case 1:
		job := Job{jobParts[0], "", ""}
		return components[0], job

	default:
		job := Job{jobParts[0], jobParts[2], jobParts[1]}
		return components[0], job
	}
}

type IntStack []int64

func (s *IntStack) Pop() int64 {
	idx := len(*s) - 1
	r := (*s)[idx]
	*s = (*s)[:idx]

	return r
}

func (s *IntStack) Push(n int64) {
	*s = append(*s, n)
}
