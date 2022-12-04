package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day04.txt"

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

	re := regexp.MustCompile(`(\d+)-(\d+),(\d+)-(\d+)`)

	count := 0
	for _, ln := range lines {
		ssm := re.FindStringSubmatch(ln)

		mina, err := strconv.Atoi(ssm[1])
		utils.Check(err, "unable to convert string to int")
		maxa, err := strconv.Atoi(ssm[2])
		utils.Check(err, "unable to convert string to int")
		minb, err := strconv.Atoi(ssm[3])
		utils.Check(err, "unable to convert string to int")
		maxb, err := strconv.Atoi(ssm[4])
		utils.Check(err, "unable to convert string to int")

		workA := make(map[int]bool)
		for i := mina; i <= maxa; i++ {
			workA[i] = true
		}

		for i := minb; i <= maxb; i++ {
			if workA[i] {
				count++
				break
			}
		}
	}

	return count
}
