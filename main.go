package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jdelkins/shorten_path/spath"
	"github.com/pborman/getopt/v2"
)

// TODO: also consider the length of lead-in and lead-out strings
func pathLength(pe spath.PathElements, use []bool) int {
	res := 0
	for i := 0; i < len(pe); i++ {
		pe_len := len(pe[i].OrigElement)
		if use[i] {
			pe_len = len(pe[i].ShortElement)
		}
		if i > 0 {
			// add the /
			pe_len++
		}
		res += pe_len
	}
	return res
}

func main() {
	helpFlag := getopt.BoolLong("help", '?', "display help")
	shortLeadIn := getopt.StringLong("short-lead-in", 'i', "", "character sequence to begin abbreviated elements")
	shortLeadOut := getopt.StringLong("short-lead-out", 'o', "", "character sequence to end abbreviated elements")
	headLeadIn := getopt.StringLong("head-lead-in", 'h', "", "character sequence to begin the first element")
	headLeadOut := getopt.StringLong("head-lead-out", 'H', "", "character sequence to end the first element")
	tailLeadIn := getopt.StringLong("tail-lead-in", 't', "", "character sequence to begin the last element")
	tailLeadOut := getopt.StringLong("tail-lead-out", 'T', "", "character sequence to end the last element")
	length := getopt.IntLong("length", 'l', 1, "length of path above which shortening will be attempted")
	minSavings := getopt.IntLong("minimum-element-savings", 'm', 1, "don't abbreviate a path element "+
		"unless doing so will result in at least this many charcaters saved; "+
		"use to compensate for printable lead-in or lead-out strings, if any")
	getopt.Parse()

	if *helpFlag {
		getopt.PrintUsage(os.Stdout)
		return
	}

	pth := strings.Join(getopt.Args(), " ")
	if pth == "" {
		return
	}

	short := spath.Shorten(spath.Homealize(spath.Components(path.Clean(pth))))
	use := make([]bool, len(short))
	for i := 1; i < len(short)-1; i++ {
		if pathLength(short, use) > *length && short[i].Shortened && len(short[i].OrigElement)-len(short[i].ShortElement) >= *minSavings {
			use[i] = true
		}
	}

	formatter := func(i int, pe *spath.PathElement) string {
		if i == 0 {
			return *headLeadIn + pe.OrigElement + *headLeadOut
		}
		if i == len(short)-1 {
			return *tailLeadIn + pe.OrigElement + *tailLeadOut
		}
		if use[i] {
			return *shortLeadIn + pe.ShortElement + *shortLeadOut
		}
		return pe.OrigElement
	}

	fmt.Println(path.Join(short.Map(formatter)...))
}
