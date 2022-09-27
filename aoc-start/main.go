package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/spf13/pflag"
)

func myUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	pflag.PrintDefaults()
}

func main() {
	now := time.Now()
	day := pflag.IntP("day", "d", now.Day(), "Advent of Code Day")
	help_flag := pflag.BoolP("help", "h", false, "show help")

	pflag.Usage = myUsage
	pflag.Parse()

	if *help_flag {
		myUsage()
		os.Exit(0)
	}

	dirname := fmt.Sprintf("day%02d_p1", *day)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		panic(err)
	}

	tpl, err := template.ParseFiles("template/main.go")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Join(dirname, "main.go"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tpl.Execute(f, *day); err != nil {
		panic(err)
	}

	ftest, err := os.Open("template/main_test.go")
	if err != nil {
		panic(err)
	}

	ftestout, err := os.Create(path.Join(dirname, "main_test.go"))
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(ftestout, ftest); err != nil {
		panic(err)
	}
}
