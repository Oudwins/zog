package internals

import (
	"fmt"
	"strings"

	"github.com/Oudwins/zog/zconst"
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

	// Saving old functions required here to prevent recursive call during assignment.
	oldFn := t.ValidateFunc
	t.ValidateFunc = func(val any, ctx Ctx) bool {
		return !oldFn(val, ctx)
	}
	t.IssueCode = NotIssueCode(t.IssueCode)

	return append(testArr, t)
}

func NotIssueCode(e zconst.ZogIssueCode) string {
	if strings.HasPrefix(e, "not_") {
		return zconst.ZogIssueCode(strings.TrimPrefix(e, "not_"))
	}
	return zconst.ZogIssueCode(fmt.Sprintf("not_%s", e))
}
