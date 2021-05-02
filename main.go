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
func partialLength(pe spath.PathElements, use []bool) int {
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
		if partialLength(short, use) > *length {
			use[i] = true
		}
	}

	formatter := func(i int, pe *spath.PathElement) string {
		if pe.Shortened && use[i] {
			return *leadIn + pe.ShortElement + *leadOut
		}
		return pe.OrigElement
	}

	fmt.Println(path.Join(short.Map(formatter)...))
}
