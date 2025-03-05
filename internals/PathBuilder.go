package internals

// type PathBuilder string

// func (p PathBuilder) Push(path string) PathBuilder {
// 	if p == "" {
// 		return PathBuilder(path)
// 	}
// 	if path[0] == '[' {
// 		return p + PathBuilder(path)
// 	}
// 	return p + PathBuilder("."+path)
// }

// func (p PathBuilder) String() string {
// 	return string(p)
// }

func NewPathBuilder() *PathBuilder {
	pb := PathBuilderPool.Get().(*PathBuilder)
	*pb = (*pb)[:1]
	return pb
}

type PathBuilder []string

func (p *PathBuilder) Push(path *string) *PathBuilder {
	*p = append(*p, *path)
	return p
}

func (p *PathBuilder) Pop() {
	if len(*p) == 0 {
		return
	}
	*p = (*p)[:len(*p)-1]
}

func (p *PathBuilder) String() string {
	sb := NewStringBuilder()
	defer FreeStringBuilder(sb)
	for i, v := range *p {
		if i > 0 && (*p)[i-1] != "" && v[0] != '[' {
			sb.WriteString(".")
		}
		sb.WriteString(v)
	}
	return sb.String()
}

func (p *PathBuilder) Free() {
	PathBuilderPool.Put(p)
}
