package primitives

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

// Default error func that takes the value as param. Expect msg = "%s is a required field"
func DErrorFuncWithVal(msg string) ErrorFunc {
	return func(val any, ctx *ParseCtx) string {
		return fmt.Sprintf(msg, val)
	}
}

// Default error func, doesn't take any params
func DErrorFunc(msg string) ErrorFunc {
	return func(val any, ctx *ParseCtx) string {
		return msg
	}
}

func Required(fn ErrorFunc) Test {
	t := Test{
		Name:      "required",
		ErrorFunc: fn,
		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			return !IsZeroValue(val)
		},
	}
	return t
}

type LengthCapable[K any] interface {
	~[]any | ~[]K | string | map[any]any | ~chan any
}

func LenMin[T LengthCapable[any]](n int, errFn ErrorFunc) Test {
	return Test{
		Name:      "min",
		ErrorFunc: errFn,
		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			x := val.(T)
			return len(x) >= n
		},
	}
}

func LenMax[T LengthCapable[any]](n int, errFn ErrorFunc) Test {
	return Test{
		Name: "max",
		ValidateFunc: func(v any, ctx *ParseCtx) bool {
			val, ok := v.(T)
			if !ok {
				return false
			}
			return len(val) <= n
		},
		ErrorFunc: errFn,
	}
}

func Len[T LengthCapable[any]](n int, errFn ErrorFunc) Test {
	return Test{
		Name: "length",
		ValidateFunc: func(v any, ctx *ParseCtx) bool {
			val, ok := v.(T)
			if !ok {
				return false
			}
			return len(val) == n
		},
		ErrorFunc: errFn,
	}
}

func In[T any](values []T, msg string) Test {
	return Test{
		Name: "oneof",
		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			for _, value := range values {
				v := val.(T)
				if reflect.DeepEqual(v, value) {
					return true
				}
			}
			return false
		},
		ErrorFunc: DErrorFunc(msg),
	}
}

func EQ[T comparable](n T, msg string) Test {
	return Test{
		Name: "eq",
		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v == n
		},
		ErrorFunc: DErrorFunc(msg),
	}
}

func LTE[T constraints.Ordered](n T, msg string) Test {
	return Test{
		Name: "lte",

		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v <= n
		},
		ErrorFunc: DErrorFunc(msg),
	}
}

func GTE[T constraints.Ordered](n T, msg string) Test {
	return Test{
		Name: "gte",

		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v >= n
		},
		ErrorFunc: DErrorFunc(msg),
	}
}

func LT[T constraints.Ordered](n T, msg string) Test {
	return Test{
		Name: "lt",

		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v < n
		},
		ErrorFunc: DErrorFunc(msg),
	}
}

func GT[T constraints.Ordered](n T, msg string) Test {
	return Test{
		Name: "gt",

		ValidateFunc: func(val any, ctx *ParseCtx) bool {
			v, ok := val.(T)
			if !ok {
				return false
			}
			return v > n
		},
		ErrorFunc: DErrorFunc(msg),
	}
}
