package primitives

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

func Required() Test {
	t := Test{
		ErrCode: ErrCodeRequired,
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
		ErrCode: ErrCodeMin,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			x := val.(T)
			return len(x) >= n
		},
	}
	t.Params[ErrCodeMin] = n
	return t
}

func LenMax[T LengthCapable[any]](n int) Test {
	t := Test{
		ErrCode: ErrCodeMax,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(T)
			if !ok {
				return false
			}
			return len(val) <= n
		},
	}
	t.Params[ErrCodeMax] = n
	return t
}

func Len[T LengthCapable[any]](n int) Test {
	t := Test{
		ErrCode: ErrCodeLen,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(T)
			if !ok {
				return false
			}
			return len(val) == n
		},
	}
	t.Params[ErrCodeLen] = n
	return t
}

func In[T any](values []T) Test {
	t := Test{
		ErrCode: ErrCodeOneOf,
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
	t.Params[ErrCodeOneOf] = values
	return t
}

func EQ[T comparable](n T) Test {
	t := Test{
		ErrCode: ErrCodeEQ,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v == n
		},
	}
	t.Params[ErrCodeEQ] = n
	return t
}

func LTE[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: ErrCodeLTE,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v <= n
		},
	}
	t.Params[ErrCodeLTE] = n
	return t
}

func GTE[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: ErrCodeGTE,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v >= n
		},
	}
	t.Params[ErrCodeGTE] = n
	return t
}

func LT[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: ErrCodeLT,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v < n
		},
	}
	t.Params[ErrCodeLT] = n
	return t
}

func GT[T constraints.Ordered](n T) Test {
	t := Test{
		ErrCode: ErrCodeGT,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v > n
		},
	}
	t.Params[ErrCodeGT] = n
	return t
}
