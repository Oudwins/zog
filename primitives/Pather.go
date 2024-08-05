package primitives

import "strings"

type Pather string

func (p Pather) Push(path string) Pather {
	if p == "" {
		return Pather(path)
	}
	return p + Pather("."+path)
}
func (p Pather) Pop() Pather {
	return p[:strings.LastIndex(string(p), ".")]
}

func (p Pather) String() string {
	return string(p)
}
