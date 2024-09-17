package internals

import "strings"

type PathBuilder string

func (p PathBuilder) Push(path string) PathBuilder {
	if p == "" {
		return PathBuilder(path)
	}
	if path[0] == '[' {
		return p + PathBuilder(path)
	}
	return p + PathBuilder("."+path)
}
func (p PathBuilder) Pop() PathBuilder {
	return p[:strings.LastIndex(string(p), ".")]
}

func (p PathBuilder) String() string {
	return string(p)
}
