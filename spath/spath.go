package spath

import (
	"fmt"
	"os"
	"path"
)

type PathElement struct {
	OrigElement  string
	ShortElement string
	Shortened    bool
}
type PathElements []PathElement

func (pe PathElements) Map(trans func(int, *PathElement) string) []string {
	res := make([]string, len(pe))
	for i, v := range pe {
		res[i] = trans(i, &v)
	}
	return res
}

func Homealize(p []string) []string {
	hd, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get user home directory")
		return p
	}
	hc := Components(hd)
	if len(p) < len(hc) {
		return p
	}
	for i, v := range hc {
		if p[i] != v {
			return p
		}
	}
	return append([]string{"~"}, p[len(hc):]...)
}

func Dehomealize(p []string) []string {
	if len(p) == 0 || p[0] != "~" {
		return p
	}
	hd, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't determine user home directory")
		return p
	}
	hc := Components(hd)
	return append(hc, p[1:]...)
}

func reverse(l []string) {
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
}

func Components(p string) []string {
	res := []string{}
	if p == "" || p == "." || p == "/" {
		return []string{p}
	}
	dir, file := p, ""
	for dir != "." && dir != "" && dir != "/" {
		dir, file = path.Split(dir)
		dir = path.Clean(dir)
		if dir != "/" {
			res = append(res, file)
		}
	}
	if dir == "/" && file != "" {
		res = append(res, path.Join(dir, file))
	}
	reverse(res)
	return res
}

func unique(dir []string, pfx string) bool {
	if pfx == "." || pfx == ".." {
		return false
	}
	files, err := os.ReadDir(path.Join(dir...))
	if err != nil {
		return pfx != "."
	}
	cnt := 0
	for _, f := range files {
		n := f.Name()
		if len(n) >= len(pfx) && n[:len(pfx)] == pfx {
			cnt++
			if cnt > 1 {
				return false
			}
		}
	}
	return cnt < 2
}

func Shorten(p []string) PathElements {
	res := make(PathElements, len(p))
	res[0] = PathElement{p[0], p[0], false}
	for i, v := range p {
		if i == 0 || i == len(p)-1 {
			continue
		}
		found_abbrev := false
		for j := 1; j < len(v); j++ {
			pfx := v[:j]
			if unique(Dehomealize(p[:i]), pfx) {
				found_abbrev = true
				res[i] = PathElement{v, pfx, true}
				break
			}
		}
		if !found_abbrev {
			res[i] = PathElement{v, v, false}
		}
	}
	res[len(p)-1] = PathElement{p[len(p)-1], p[len(p)-1], false}
	return res
}
