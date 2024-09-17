package internals

import (
	"reflect"

	zconst "github.com/Oudwins/zog/zconst"
	"golang.org/x/exp/constraints"
)

func Required() Test {
	t := Test{
		ErrCode: zconst.ErrCodeRequired,
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			return !IsZeroValue(val)
		},
	}
	return t
}

type LengthCapable[K any] interface {
	~[]any | ~[]K | string | map[any]any | ~chan any
}

func LenMin[T LengthCapable[any]](n int) Test {
	t := Test{
		ErrCode: zconst.ErrCodeMin,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			x := val.(T)
			return len(x) >= n
		},
	}
	t.Params[zconst.ErrCodeMin] = n
	return t
}

func LenMax[T LengthCapable[any]](n int) Test {
	t := Test{
		ErrCode: zconst.ErrCodeMax,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(T)
			if !ok {
				return false
			}
			return len(val) <= n
		},
	}
	t.Params[zconst.ErrCodeMax] = n
	return t
}

func Len[T LengthCapable[any]](n int) Test {
	t := Test{
		ErrCode: zconst.ErrCodeLen,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(T)
			if !ok {
				return false
			}
			return len(val) == n
		},
	}
	t.Params[zconst.ErrCodeLen] = n
	return t
}

func In[T any](values []T) Test {
	t := Test{
		ErrCode: zconst.ErrCodeOneOf,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			for _, value := range values {
				v := val.(T)
				if reflect.DeepEqual(v, value) {
					return true
				}
			}
			return false
		},
	}
	t.Params[zconst.ErrCodeOneOf] = values
	return t
}

func EQ[T comparable](n T) Test {
	t := Test{
		ErrCode: zconst.ErrCodeEQ,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v == n
		},
	}
	t.Params[zconst.ErrCodeEQ] = n
	return t
}

func LTE[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: zconst.ErrCodeLTE,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v <= n
		},
	}
	t.Params[zconst.ErrCodeLTE] = n
	return t
}

func GTE[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: zconst.ErrCodeGTE,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v >= n
		},
	}
	t.Params[zconst.ErrCodeGTE] = n
	return t
}

func LT[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: zconst.ErrCodeLT,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v < n
		},
	}
	t.Params[zconst.ErrCodeLT] = n
	return t
}

func GT[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: zconst.ErrCodeGT,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v > n
		},
	}
	t.Params[zconst.ErrCodeGT] = n
	return t
}
