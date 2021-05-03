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
		if (i > 0) {
			// add the /
			pe_len++
		}
		res += pe_len
	}
	return res
}

func main() {
	helpFlag := getopt.BoolLong("help", 'h', "display help")
	leadIn := getopt.StringLong("lead-in", 'i', "", "character sequence to begin abbreviated elements")
	leadOut := getopt.StringLong("lead-out", 'o', "", "character sequence to end abbreviated elements")
	length := getopt.IntLong("length", 'l', 1, "length of path above which shortening will be attempted")
	minSavings := getopt.IntLong("minimum-element-savings", 'm', 1, "don't abbreviate a path element " +
								 "unless doing so will result in at least this many charcaters saved; " +
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
		if pathLength(short, use) > *length && short[i].Shortened && len(short[i].OrigElement) - len(short[i].ShortElement) >= *minSavings {
			use[i] = true
		}
	}

	formatter := func(i int, pe *spath.PathElement) string {
		if use[i] {
			return *leadIn + pe.ShortElement + *leadOut
		}
		return pe.OrigElement
	}

	fmt.Println(path.Join(short.Map(formatter)...))
}
