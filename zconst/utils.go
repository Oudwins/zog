package zconst

import "strings"

const notPrefix = "not_"

func NotIssueCode(e ZogIssueCode) string {
	if strings.HasPrefix(e, notPrefix) {
		return ZogIssueCode(strings.TrimPrefix(e, notPrefix))
	}
	return ZogIssueCode(notPrefix + e)
}
