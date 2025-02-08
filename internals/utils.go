package internals

import (
	"fmt"
)

func SafeString(x any) string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", x)
}

func SafeError(x error) string {
	if x == nil {
		return "<nil>"
	}
	return x.Error()
}

func AddTest(testArr []Test, t Test, isNot bool) []Test {
	if !isNot {
		return append(testArr, t)
	}

	oldFn := t.ValidateFunc
	t.ValidateFunc = func(val any, ctx ParseCtx) bool {
		return !oldFn(val, ctx)
	}

	return append(testArr, t)
}
