package internals

import (
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
		return &ZogErr{}
	},
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
			return &ZogErr{}
		},
	}
}

func Clear() {
	ClearPools()
}
