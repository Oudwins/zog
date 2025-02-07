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
	if isNot {
		t.ValidateFunc = func(val any, ctx ParseCtx) bool {
			return !t.ValidateFunc(val, ctx)
		}
	}

	return append(testArr, t)
}
