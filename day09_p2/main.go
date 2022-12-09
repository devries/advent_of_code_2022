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

const inputfile = "../inputs/day09.txt"

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
	knots := make([]utils.Point, 10)
	visited := make(map[utils.Point]bool)
	visited[utils.Point{X: 0, Y: 0}] = true

	for _, ln := range lines {
		components := strings.Fields(ln)

		steps, err := strconv.Atoi(components[1])
		utils.Check(err, "unable to convert %s to integer", components[1])

		var direction utils.Point
		for i := 0; i < steps; i++ {
			switch components[0] {
			case "R":
				direction = utils.East
			case "U":
				direction = utils.North
			case "L":
				direction = utils.West
			case "D":
				direction = utils.South
			}

			knots[0] = knots[0].Add(direction)

			for i := 0; i < 9; i++ {
				knots[i+1] = movePair(knots[i], knots[i+1])
			}

			visited[knots[9]] = true
		}
	}

	return len(visited)
}

func movePair(head utils.Point, tail utils.Point) utils.Point {
	// distance between head and tail
	d := utils.Point{X: head.X - tail.X, Y: head.Y - tail.Y}

	// Move in direction of separation if separated
	moveX := false
	moveY := false
	switch {
	case d.X > 1:
		tail = tail.Add(utils.East)
		moveX = true
	case d.X < -1:
		tail = tail.Add(utils.West)
		moveX = true
	case d.Y > 1:
		tail = tail.Add(utils.North)
		moveY = true
	case d.Y < -1:
		tail = tail.Add(utils.South)
		moveY = true
	}

	if moveX {
		// if moved in X, check if it needs to move in Y too
		switch {
		case d.Y > 0:
			tail = tail.Add(utils.North)
		case d.Y < 0:
			tail = tail.Add(utils.South)
		}
	}

	if moveY {
		// if moved in Y, check if it needs to move in X too
		switch {
		case d.X > 0:
			tail = tail.Add(utils.East)
		case d.X < 0:
			tail = tail.Add(utils.West)
		}
	}

	return tail
}
