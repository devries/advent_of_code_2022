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
	monkeys["humn"] = Job{"x", "", ""}

	equality := monkeys["root"]
	listA := rpnize(equality.A, monkeys)
	listB := rpnize(equality.B, monkeys)

	var listVar []string
	var listPlain []string

	// Find which list contains the variable
	found := false
	for _, v := range listA {
		if v == "x" {
			listVar = listA
			listPlain = listB
			found = true
			break
		}
	}
	if !found {
		listVar = listB
		listPlain = listA
	}

	// Calculate the right side of the equality
	s := doCalculation(listPlain)

	// Split the other side into pre and post variable
	pre, post := splitAroundVariable(listVar)

	// Calculate what comes before the variable and grab the stack
	pcalc := doCalculation(pre)

	// Condense the post variable RPN list by doing any possible portions
	// of the calculation
	condy := condense(post)

	complete := undoCalculation(condy, pcalc, s[0])

	return complete[0]
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

func rpnize(root string, monkeys map[string]Job) []string {
	rpnList := []string{root}

	// Expand out
	keepGoing := true
	for keepGoing {
		// Do this until there are no more monkey names in stack
		keepGoing = false
		newList := []string{}
		// Move from bottom to top
		for _, m := range rpnList {
			if strings.ContainsAny(m, "0123456789+-/*") || m == "x" {
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

	return rpnList
}

func doCalculation(rpnList []string) IntStack {
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
	return nStack
}

func undoCalculation(rpnList []string, nStack IntStack, equality int64) IntStack {
	// Reverse stack
	// Reverse the stack
	for i := 0; i < len(nStack)/2; i++ {
		nStack[i], nStack[len(nStack)-1-i] = nStack[len(nStack)-1-i], nStack[i]
	}

	// Append equality to stack
	nStack = append(nStack, equality)

	// Reverse the condensed list mostly, but keep numbers paired with their operations
	// Condensing already swapped the numbers so their operations preceed them
	for i := 0; i < len(rpnList)/2; i++ {
		rpnList[i], rpnList[len(rpnList)-1-i] = rpnList[len(rpnList)-1-i], rpnList[i]
	}

	// undo calculations
	for i := 0; i < len(rpnList); i++ {
		v := rpnList[i]
		switch v {
		case "+":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a - b)
		case "-":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a + b)
		case "*":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a / b)
		case "/":
			a := nStack.Pop()
			b := nStack.Pop()
			nStack.Push(a * b)
		default:
			n, err := strconv.ParseInt(v, 10, 64)
			utils.Check(err, "unable to convert %s to integer", v)
			i++
			v = rpnList[i]
			l := nStack.Pop()
			switch v {
			case "+":
				nStack.Push(l - n)
			case "*":
				nStack.Push(l / n)
			case "-":
				nStack.Push(n - l)
			case "/":
				nStack.Push(n / l)
			}
		}
	}

	return nStack
}

func doRevCalculation(rpnList []string, nStack IntStack) IntStack {
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
			t := nStack.Pop()
			nStack.Push(n)
			nStack.Push(t)
		}
	}
	return nStack
}

func splitAroundVariable(rpnList []string) ([]string, []string) {
	for i, p := range rpnList {
		if p == "x" {
			return rpnList[:i], rpnList[i+1:]
		}
	}

	empty := []string{}
	return rpnList, empty
}

func condense(rpnList []string) []string {
	nStack := IntStack{}
	remain := []string{}

	for _, v := range rpnList {
		switch v {
		case "+", "-", "/", "*":
			switch len(nStack) {
			case 0:
				remain = append(remain, v)
			case 1:
				remain = append(remain, v)
				remain = append(remain, fmt.Sprintf("%d", nStack.Pop()))
			default:
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
		default:
			n, err := strconv.ParseInt(v, 10, 64)
			utils.Check(err, "unable to convert %s to integer", v)
			nStack.Push(n)
		}
	}

	return remain
}
