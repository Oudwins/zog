package internals

import (
	"reflect"

	zconst "github.com/Oudwins/zog/zconst"
	"golang.org/x/exp/constraints"
)

// A TestFunc that takes the data as input and returns a boolean indicating if it is valid or not
type BoolTFunc = func(val any, ctx Ctx) bool

// TFunc is the function that tests hold that execute on the data for validation. They use the z.Ctx to add issues if needed
type TFunc = func(val any, ctx Ctx)

// TestOption is the option for a test
type TestOption = func(test *Test)

func TestFuncFromBool(fn BoolTFunc, test *Test) {
	test.Func = func(val any, ctx Ctx) {
		if !fn(val, ctx) {
			c := ctx.(*SchemaCtx)
			ctx.AddIssue(c.IssueFromTest(c.Test, val))
		}
	}
}

func NewTestFunc(IssueCode zconst.ZogIssueCode, fn BoolTFunc, options ...TestOption) *Test {
	t := &Test{
		IssueCode: IssueCode,
	}
	for _, opt := range options {
		opt(t)
	}
	TestFuncFromBool(fn, t)
	return t
}

// Test is a struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.
type Test struct {
	// The issue code that will be used to create the issue if the test fails
	IssueCode zconst.ZogIssueCode
	// The path to the issue
	IssuePath string
	// The params for the test which are passed to the issue and can be used for customizing issue messages. For example `z.String().Min(3)` will have a param with the key `min` and the value `3`.
	Params map[string]any
	// The formatter for the issue. You can use this if you want to customize the message formatter for issues generated by this specific test. Not super useful but just here for completeness.
	IssueFmtFunc IssueFmtFunc
	// The function that will be executed when evaluating the test.
	Func TFunc
}

// returns a required test to be used for processor.Required() method
func Required() Test {
	t := Test{
		IssueCode: zconst.IssueCodeRequired,
		// this is not an accident. required is only a test because it makes it easier to handle error messages. But the function to check if the value is a zero value is out of the scope of this test.
		Func: nil,
	}
	return t
}

type LengthCapable[K any] interface {
	~[]any | ~[]K | ~string | map[any]any | ~chan any
}

func LenMin[T LengthCapable[any]](n int) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		x := val.(*T)
		return len(*x) >= n
	}

	t := Test{
		IssueCode: zconst.IssueCodeMin,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMin] = n
	return t, fn
}

func LenMax[T LengthCapable[any]](n int) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		x := val.(*T)
		return len(*x) <= n
	}

	t := Test{
		IssueCode: zconst.IssueCodeMax,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMax] = n
	return t, fn
}

func Len[T LengthCapable[any]](n int) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		x := val.(*T)
		return len(*x) == n
	}

	t := Test{
		IssueCode: zconst.IssueCodeLen,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeLen] = n
	return t, fn
}

func In[T any](values []T) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		for _, value := range values {
			v := val.(*T)
			if reflect.DeepEqual(*v, value) {
				return true
			}
		}
		return false
	}

	t := Test{
		IssueCode: zconst.IssueCodeOneOf,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeOneOf] = values
	return t, fn
}

func EQ[T comparable](n T) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		v, ok := val.(*T)
		if !ok {
			return false
		}
		return *v == n
	}

	t := Test{
		IssueCode: zconst.IssueCodeEQ,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeEQ] = n
	return t, fn
}

func LTE[T constraints.Ordered](n T) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		v, ok := val.(*T)
		if !ok {
			return false
		}
		return *v <= n
	}

	t := Test{
		IssueCode: zconst.IssueCodeLTE,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeLTE] = n
	return t, fn
}

func GTE[T constraints.Ordered](n T) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		v, ok := val.(*T)
		if !ok {
			return false
		}
		return *v >= n
	}

	t := Test{
		IssueCode: zconst.IssueCodeGTE,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeGTE] = n
	return t, fn
}

func LT[T constraints.Ordered](n T) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		v, ok := val.(*T)
		if !ok {
			return false
		}
		return *v < n
	}

	t := Test{
		IssueCode: zconst.IssueCodeLT,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeLT] = n
	return t, fn
}

func GT[T constraints.Ordered](n T) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		v, ok := val.(*T)
		if !ok {
			return false
		}
		return *v > n
	}

	t := Test{
		IssueCode: zconst.IssueCodeGT,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeGT] = n
	return t, fn
}
