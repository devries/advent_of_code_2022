package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day17.txt"

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

	motions := []rune(lines[0])
	grid := []uint8{0b1111111}
	motionStep := 0

	for i := 0; i < 2022; i++ {
		mino := make(Mino, len(minos[i%5]))
		// plotGame(grid)

		offset := len(grid) + 3
		grid = append(grid, 0, 0, 0)
		for j := 0; j < len(mino); j++ {
			grid = append(grid, 0)
			mino[j] = minos[i%5][j]
		}

		for {
			m := motions[motionStep%len(motions)]
			next := make(Mino, len(mino))

			switch m {
			case '>':
				// move right
				for j := 0; j < len(mino); j++ {
					if mino[j]&0b1000000 > 0 {
						next = mino
						break
					}
					next[j] = mino[j] << 1
				}

			case '<':
				// move left
				for j := 0; j < len(mino); j++ {
					if mino[j]&0b1 > 0 {
						next = mino
						break
					}
					next[j] = mino[j] >> 1
				}
			}
			motionStep++

			if check_possible(next, grid, offset) {
				mino = next
			}

			// Now drop one unit
			offset--

			if !check_possible(mino, grid, offset) {
				for j := 0; j < len(mino); j++ {
					grid[offset+j+1] |= mino[j]
				}
				break
			}
		}

		for i := len(grid) - 1; i >= 0; i-- {
			if grid[i] != 0 {
				grid = grid[:i+1]
				break
			}
		}
	}

	return int64(len(grid) - 1)
}

type Mino []uint8

// Encoded backwards and bottom to top
var minos = []Mino{{0b111100},
	{0b01000,
		0b11100,
		0b01000},
	{0b11100,
		0b10000,
		0b10000},
	{0b100,
		0b100,
		0b100,
		0b100},
	{0b1100,
		0b1100}}

func check_possible(next []uint8, grid []uint8, offset int) bool {
	for i, row := range next {
		if row&grid[offset+i] > 0 { // Intersection
			return false
		}
	}

	return true
}

func plotGame(grid []uint8) {
	for i := len(grid) - 1; i >= 0; i-- {
		s := fmt.Sprintf("%07b", grid[i])
		fmt.Printf("%s\n", reverse(s))
	}
	fmt.Printf("\n")
}

func reverse(s string) string {
	r := []rune(s)

	rev := make([]rune, len(r))

	for i := 0; i < len(r); i++ {
		rev[len(r)-1-i] = r[i]
	}

	return strings.ReplaceAll(strings.ReplaceAll(string(rev), "1", "#"), "0", ".")
}
