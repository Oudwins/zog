package internals

import (
	"fmt"
	"strings"

	"github.com/Oudwins/zog/zconst"
)

const (
	notPrefix     = "not_"
	defaultString = "<nil>"
)

func SafeString(x any) string {
	if x == nil {
		return defaultString
	}
	return fmt.Sprintf("%v", x)
}

func SafeError(x error) string {
	if x == nil {
		return defaultString
	}
	return x.Error()
}

func NotIssueCode(e zconst.ZogIssueCode) string {
	if strings.HasPrefix(e, notPrefix) {
		return zconst.ZogIssueCode(strings.TrimPrefix(e, notPrefix))
	}
	return zconst.ZogIssueCode(notPrefix + e)
}

func PtrOf[T any](v T) *T {
	return &v
}
