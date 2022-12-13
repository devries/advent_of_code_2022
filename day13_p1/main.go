package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day13.txt"

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

	idxsum := 0
	for i := 0; i < len(lines); i += 3 {
		pa, _ := parsePacket([]rune(lines[i]), 0)
		pb, _ := parsePacket([]rune(lines[i+1]), 0)

		cmp := pa.Compare(pb)
		if cmp < 0 {
			idxsum += i/3 + 1
		}
	}

	return idxsum
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Item interface {
	Compare(Item) int
}

type PacketList []Item
type PacketNumber int

func (a PacketList) Compare(b Item) int {
	switch v := b.(type) {
	case PacketList:
		// Comparing list to list
		for i := 0; i < min(len(a), len(v)); i++ {
			c := a[i].Compare(v[i])
			if c != 0 {
				return c
			}
		}
		return len(a) - len(v)

	case PacketNumber:
		// Comparing list to number
		bl := PacketList{b}
		return a.Compare(bl)
	}

	return 0
}

func (a PacketNumber) Compare(b Item) int {
	switch v := b.(type) {
	case PacketList:
		// Comparing number to list
		al := PacketList{a}
		return al.Compare(b)

	case PacketNumber:
		// Comparing number to number
		return int(a - v)
	}

	return 0
}

func parsePacket(chars []rune, pos int) (Item, int) {
	if chars[pos] == ',' {
		pos++
	}
	c := pos

	switch chars[c] {
	case '[':
		// This is a list
		ret := PacketList{}
		var el Item
		c++
		for chars[c] != ']' {
			el, c = parsePacket(chars, c)
			ret = append(ret, el)
		}
		return ret, c + 1

	default:
		// This is a number
		for chars[c] != ',' && chars[c] != ']' {
			c++
		}
		ns := string(chars[pos:c])
		n, err := strconv.Atoi(ns)
		utils.Check(err, "Unable to convert %s to integer", ns)
		return PacketNumber(n), c
	}
}
