package primitives

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

type LengthCapable[K any] interface {
	~[]any | ~[]K | string | map[any]any | ~chan any
}

func IsType[T any](msg string) Rule {
	return Rule{
		Name:         "isType",
		ErrorMessage: msg,
		ValidateFunc: func(set Rule) bool {
			_, ok := set.FieldValue.(T)
			return ok
		},
	}
}

func LenMin[T LengthCapable[any]](n int, msg string) Rule {
	return Rule{
		Name:      "min",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			val, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return len(val) >= n
		},
		ErrorMessage: msg,
	}
}

func LenMax[T LengthCapable[any]](n int, msg string) Rule {
	return Rule{
		Name:      "max",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			val, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return len(val) <= n
		},
		ErrorMessage: msg,
	}
}

func Length[T LengthCapable[any]](n int, msg string) Rule {
	return Rule{
		Name:      "length",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			val, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return len(val) == n
		},
		ErrorMessage: msg,
	}
}

func In[T any](values []T, msg string) Rule {
	return Rule{
		Name:      "in",
		RuleValue: values,
		ValidateFunc: func(set Rule) bool {
			for _, value := range values {
				v := set.FieldValue.(T)
				if reflect.DeepEqual(v, value) {
					return true
				}
			}
			return false
		},
		ErrorMessage: msg,
	}
}

func EQ[T comparable](n T, msg string) Rule {
	return Rule{
		Name:      "eq",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			v, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return v == n
		},
		ErrorMessage: msg,
	}
}

func LTE[T constraints.Ordered](n T, msg string) Rule {
	return Rule{
		Name:      "lte",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			v, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return v <= n
		},
		ErrorMessage: msg,
	}
}

func GTE[T constraints.Ordered](n T, msg string) Rule {
	return Rule{
		Name:      "gte",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			v, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return v >= n
		},
		ErrorMessage: msg,
	}
}

func LT[T constraints.Ordered](n T, msg string) Rule {
	return Rule{
		Name:      "lt",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			v, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return v < n
		},
		ErrorMessage: msg,
	}
}

func GT[T constraints.Ordered](n T, msg string) Rule {
	return Rule{
		Name:      "gt",
		RuleValue: n,
		ValidateFunc: func(set Rule) bool {
			v, ok := set.FieldValue.(T)
			if !ok {
				return false
			}
			return v > n
		},
		ErrorMessage: msg,
	}
}
