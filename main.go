package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jdelkins/shorten_path/spath"
	"github.com/pborman/getopt/v2"
)

func main() {
	helpFlag := getopt.BoolLong("help", 'h', "display help")
	leadIn := getopt.StringLong("lead-in", 'i', "", "character sequence to begin abbreviated elements")
	leadOut := getopt.StringLong("lead-out", 'o', "", "character sequence to end abbreviated elements")
	getopt.Parse()

	if *helpFlag {
		getopt.PrintUsage(os.Stdout)
		return
	}

	pth := strings.Join(getopt.Args(), " ")
	short := spath.Shorten(spath.Homealize(spath.Components(path.Clean(pth))))
	formatter := func(pe *spath.PathElement) string {
		if pe.Shortened {
			return *leadIn + pe.ShortElement + *leadOut
		}
		return pe.OrigElement
	}
	fmt.Println(path.Join(short.Map(formatter)...))
}
