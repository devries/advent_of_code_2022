package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/devries/advent_of_code_2022/utils"
	"github.com/spf13/pflag"
)

const inputfile = "../inputs/day07.txt"
const maxsize = 100000

func main() {
	pflag.Parse()
	f, err := os.Open(inputfile)
	utils.Check(err, "error opening input")
	defer f.Close()

	r := solve(f)
	fmt.Println(r)
}

func solve(r io.Reader) int64 {
	root := parse(r)

	dirList := []*directory{root}
	var totalSize int64

	for len(dirList) > 0 {
		d := dirList[0]
		dirList = dirList[1:]

		s := d.GetSize()
		if s <= maxsize {
			totalSize += s
		}
		dirList = append(dirList, d.Children...)
	}

	return totalSize
}

func parse(r io.Reader) *directory {
	lines := utils.ReadLines(r)

	root := directory{Name: ""}
	current := &root

lineParser:
	for _, ln := range lines {
		components := strings.Fields(ln)
		switch components[0] {
		case "$":
			// command
			switch components[1] {
			case "cd":
				if components[2] == "/" {
					current = &root
				} else if components[2] == ".." {
					current = current.Parent
				} else {
					// check if directory exists first
					name := components[2]
					for _, c := range current.Children {
						if c.Name == name {
							// directory exists, just change to it
							current = c
							continue lineParser
						}
					}
					newdir := directory{Name: name, Parent: current}
					current.Children = append(current.Children, &newdir)
					current = &newdir
				}
			}
		default:
			switch components[0] {
			case "dir":
				// directory in listing
				for _, c := range current.Children {
					if c.Name == components[1] {
						// directory exists already
						continue lineParser
					}
				}
				newdir := directory{Name: components[1], Parent: current}
				current.Children = append(current.Children, &newdir)
			default:
				size, err := strconv.ParseInt(components[0], 10, 64)
				utils.Check(err, "unable to parse %s to integer", components[0])
				for _, f := range current.Files {
					if f.Name == components[1] {
						// file exists already
						continue lineParser
					}
				}
				current.Files = append(current.Files, file{components[1], size})
			}
		}
	}

	return &root
}

type directory struct {
	Name     string
	Parent   *directory
	Files    []file
	Children []*directory
	size     int64
	once     sync.Once
}

type file struct {
	Name string
	Size int64
}

func (d *directory) GetSize() int64 {
	d.once.Do(func() {
		for _, f := range d.Files {
			d.size += f.Size
		}
		for _, subd := range d.Children {
			d.size += subd.GetSize()
		}
		if utils.Verbose {
			fmt.Printf("%s: %d\n", d.Name, d.size)
		}
	})
	return d.size
}
