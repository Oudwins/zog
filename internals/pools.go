package internals

import (
	"strings"
	"sync"
)

var ExecCtxPool = sync.Pool{
	New: func() any {
		return &ExecCtx{}
	},
}

var SchemaCtxPool = sync.Pool{
	New: func() any {
		return &SchemaCtx{}
	},
}

var InternalIssueListPool = sync.Pool{
	New: func() any {
		return &ErrsList{}
	},
}

var InternalIssueMapPool = sync.Pool{
	New: func() any {
		return &ErrsMap{}
	},
}

var ZogIssuePool = sync.Pool{
	New: func() any {
		return &ZogIssue{}
	},
}

var PathBuilderPool = sync.Pool{
	New: func() any {
		pb := make(PathBuilder, 0, 5)
		return &pb
	},
}

var StringBuilderPool = sync.Pool{
	New: func() any {
		sb := strings.Builder{}
		return &sb
	},
}

func NewStringBuilder() *strings.Builder {
	sb := StringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

func FreeStringBuilder(sb *strings.Builder) {
	StringBuilderPool.Put(sb)
}

func ClearPools() {
	ExecCtxPool = sync.Pool{
		New: func() any {
			return &ExecCtx{}
		},
	}
	SchemaCtxPool = sync.Pool{
		New: func() any {
			return &SchemaCtx{}
		},
	}
	InternalIssueListPool = sync.Pool{
		New: func() any {
			return &ErrsList{}
		},
	}
	InternalIssueMapPool = sync.Pool{
		New: func() any {
			return &ErrsMap{}
		},
	}
	ZogIssuePool = sync.Pool{
		New: func() any {
			return &ZogIssue{}
		},
	}
	PathBuilderPool = sync.Pool{
		New: func() any {
			pb := make(PathBuilder, 0, 5)
			return &pb
		},
	}
	StringBuilderPool = sync.Pool{
		New: func() any {
			sb := strings.Builder{}
			return &sb
		},
	}
}

func Clear() {
	ClearPools()
}
