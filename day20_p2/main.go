package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day20.txt"
const key int64 = 811589153

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

	length := len(lines)
	original := make([]*Element, length)
	var p *Element
	var first *Element
	var last *Element
	var zero *Element
	for i, ln := range lines {
		v, err := strconv.Atoi(ln)
		utils.Check(err, "unable to convert %s to integer", ln)

		e := Element{value: key * int64(v)}
		original[i] = &e
		if v == 0 {
			zero = &e
		}

		if p != nil {
			e.previous = p
			p.next = &e
		} else {
			first = &e
		}

		p = &e
		last = &e
	}
	last.next = first
	first.previous = last

	// printList(zero, length)

	for repeats := 0; repeats < 10; repeats++ {
		for _, p := range original {
			delta := p.value % (int64(length) - 1)

			// remove element
			p.next.previous = p.previous
			p.previous.next = p.next

			ip := p.previous // insertion point

			// find destination
			if delta > 0 {
				for i := int64(0); i < delta; i++ {
					ip = ip.next
				}
			}
			if delta < 0 {
				for i := int64(0); i < -delta; i++ {
					ip = ip.previous
				}
			}

			p.previous = ip
			p.next = ip.next
			ip.next = p
			p.next.previous = p

			// printList(zero, length)
		}
	}
	sum := int64(0)
	sum += findValue(zero, 1000, length)
	sum += findValue(zero, 2000, length)
	sum += findValue(zero, 3000, length)

	return sum
}

type Element struct {
	value    int64
	next     *Element
	previous *Element
}

func printList(e *Element, l int) {
	for i := 0; i < l; i++ {
		fmt.Println(e.value)
		e = e.next
	}
}

func findValue(start *Element, delta int, length int) int64 {
	d := delta % length

	p := start
	for i := 0; i < d; i++ {
		p = p.next
	}

	return p.value
}
