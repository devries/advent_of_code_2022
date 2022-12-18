package main

import (
	"fmt"
	"hash/maphash"
	"io"
	"os"
	"strings"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day17.txt"
const totalSteps int64 = 1000000000000

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
	minoStep := 0
	stepCount := 0
	totalHeight := 0
	heightUpToCycle := int64(0)
	stepsUpToCycle := int64(0)
	heightOfCycle := int64(0)
	stepsInCycle := int64(0)

	var h maphash.Hash
	states := make(map[State]Track)

	for {
		// plotGame(grid)
		h.Write(grid)
		state := State{motionStep, minoStep, h.Sum64()}
		h.Reset()
		l, ok := states[state]
		if ok {
			heightUpToCycle = l.Height
			stepsUpToCycle = l.Steps
			heightOfCycle = int64(totalHeight) - l.Height
			stepsInCycle = int64(stepCount) - l.Steps
			break
		}
		states[state] = Track{int64(stepCount), int64(totalHeight)}

		var ht int
		grid, motionStep, ht = step(grid, minos[minoStep], motionStep, motions)
		totalHeight += ht
		minoStep++
		minoStep %= 5
		stepCount++
	}
	repeats := totalSteps / stepsInCycle
	remains := totalSteps % stepsInCycle

	if remains < stepsUpToCycle {
		repeats--
	}

	extraTurns := totalSteps - repeats*stepsInCycle - stepsUpToCycle
	extraHeight := int64(0)

	for i := int64(0); i < extraTurns; i++ {
		var ht int
		grid, motionStep, ht = step(grid, minos[minoStep], motionStep, motions)
		minoStep++
		minoStep %= 5
		extraHeight += int64(ht)
	}

	return repeats*heightOfCycle + heightUpToCycle + extraHeight
}

func step(grid []uint8, origMino Mino, motionStep int, motions []rune) ([]uint8, int, int) {
	originalHeight := len(grid)

	offset := len(grid) + 3
	grid = append(grid, 0, 0, 0)

	mino := make(Mino, len(origMino))
	for j := 0; j < len(mino); j++ {
		grid = append(grid, 0)
		mino[j] = origMino[j]
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

	var chopped int
	grid, chopped = trim(grid)
	finalHeight := len(grid) - originalHeight + chopped

	return grid, motionStep % len(motions), finalHeight
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

type State struct {
	Motion   int
	Mino     int
	GridHash uint64
}

type Track struct {
	Steps  int64
	Height int64
}

func equalGrids(a []uint8, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Trim unused space from grid returning new grid
// and number of lines trimmed from bottom
func trim(g []uint8) ([]uint8, int) {
	accum := uint8(0)
	max := len(g)
	min := 0
	for i := len(g) - 1; i >= 0; i-- {
		if g[i] == 0 {
			max = i
		}
		accum |= g[i]

		if accum == 0b1111111 {
			min = i
			break
		}
	}

	return g[min:max], min
}
