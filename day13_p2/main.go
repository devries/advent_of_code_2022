package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

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

	packets := PacketList{}
	for _, ln := range lines {
		if ln != "" {
			p, _ := parsePacket([]rune(ln), 0)
			packets = append(packets, p)
		}
	}
	pda, _ := parsePacket([]rune("[[2]]"), 0)
	pdb, _ := parsePacket([]rune("[[6]]"), 0)
	packets = append(packets, pda, pdb)

	sort.Sort(packets)
	solution := 1
	for i, p := range packets {
		if p.String() == "[[2]]" || p.String() == "[[6]]" {
			solution *= i + 1
		}
		if utils.Verbose {
			fmt.Printf("%d: %s\n", i+1, p)
		}
	}

	return solution
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Item interface {
	Compare(Item) int
	String() string
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

	default:
		panic("this line should not be reachable")
	}
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

	default:
		panic("this line should not be reachable")
	}
}

func (a PacketList) Len() int      { return len(a) }
func (a PacketList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PacketList) Less(i, j int) bool {
	c := a[i].Compare(a[j])
	if c < 0 {
		return true
	}
	return false
}

func (a PacketList) String() string {
	els := []string{}

	for _, s := range a {
		els = append(els, s.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(els, ","))
}

func (a PacketNumber) String() string {
	return fmt.Sprintf("%d", a)
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
