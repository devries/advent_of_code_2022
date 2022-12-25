package main

import (
	"fmt"
	"io"
	"os"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day25.txt"

var conversion = map[rune]int64{'2': 2,
	'1': 1,
	'0': 0,
	'-': -1,
	'=': -2,
}

var deconversion = map[int64]rune{0: '0',
	1: '1',
	2: '2',
	3: '=',
	4: '-',
}

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) string {
	lines := utils.ReadLines(r)

	var sum int64
	for _, ln := range lines {
		n := convertNumber(ln)
		sum += n
	}

	/*
		for i := int64(0); i < 20; i++ {
			fmt.Println(i, backConvert(i))
		}

		j := int64(314159265)
		fmt.Println(j, backConvert(j))
	*/

	ret := backConvert(sum)
	return ret
}

func convertNumber(s string) int64 {
	r := []rune(s)

	ret := int64(0)
	power := int64(1)
	for i := len(r) - 1; i >= 0; i-- {
		n := conversion[r[i]] * power
		ret += n
		power *= 5
	}

	return ret
}

func backConvert(n int64) string {
	vals := []int64{}
	rem := n

	for {
		mod := rem % 5
		rem = rem / 5

		vals = append(vals, mod)
		if rem == 0 {
			break
		}
	}

	number := []rune{}
	carry := int64(0)
	for i := 0; i < len(vals); i++ {
		v := vals[i]
		v = v + carry

		if v >= 5 {
			carry = 1
			v = v - 5
		} else {
			carry = 0
		}
		digit := deconversion[v]
		if digit == '-' || digit == '=' {
			carry++
		}

		number = append(number, digit)
	}
	if carry != 0 {
		number = append(number, deconversion[carry])
	}

	// reverse number
	for i := 0; i < len(number)/2; i++ {
		number[i], number[len(number)-1-i] = number[len(number)-1-i], number[i]
	}

	return string(number)
}
