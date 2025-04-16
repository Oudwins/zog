package zconst

import "strings"

const NotIssuePrefix = "not_"

func NotIssueCode(e ZogIssueCode) string {
	if strings.HasPrefix(e, NotIssuePrefix) {
		return ZogIssueCode(strings.TrimPrefix(e, NotIssuePrefix))
	}
	return ZogIssueCode(NotIssuePrefix + e)
}
