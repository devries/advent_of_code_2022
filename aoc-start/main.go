package main

import (
	"fmt"
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
	templ := pflag.StringP("templates", "t", "template", "Directory where templates are stored")
	helpFlag := pflag.BoolP("help", "h", false, "show help")

	pflag.Usage = myUsage
	pflag.Parse()

	if *helpFlag {
		myUsage()
		os.Exit(0)
	}

	dirname := fmt.Sprintf("day%02d_p1", *day)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating destination directory %s: %s\n", dirname, err)
		os.Exit(1)
	}

	tpls := template.Must(template.ParseGlob(path.Join(*templ, "*")))

	for _, tpl := range tpls.Templates() {
		f, err := os.Create(path.Join(dirname, tpl.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file %s: %s\n", tpl.Name(), err)
			os.Exit(2)
		}
		defer f.Close()

		if err := tpl.Execute(f, *day); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing template for file %s: %s\n", tpl.Name(), err)
			os.Exit(3)
		}
	}
}
