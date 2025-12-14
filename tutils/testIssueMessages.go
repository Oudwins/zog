package tutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Oudwins/zog/i18n/en"
	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// VerifyDefaultIssueMessages verifies that all issues have valid default messages
func VerifyDefaultIssueMessages(t *testing.T, errs internals.ZogIssueList) {
	for _, err := range errs {
		c := err.Code
		m, ok := en.Map[err.Dtype][c]
		if !ok {
			m, ok = en.Map[err.Dtype][zconst.IssueCodeFallback]
			if !ok {
				panic(fmt.Sprintf("no fallback message for type %s", err.Dtype))
			}
		}
		prefix := strings.Split(m, "{{")[0]
		postfix := prefix
		if strings.Contains(m, "}}") {
			postfix = strings.Split(m, "}}")[1]
		}
		assert.True(t, strings.HasPrefix(err.Message, prefix))
		assert.True(t, strings.HasSuffix(err.Message, postfix))
	}
}

// Deprecated: Use VerifyDefaultIssueMessages instead.
// All schemas now return ZogIssueList.
func VerifyDefaultIssueMessagesMap(t *testing.T, errs internals.ZogIssueMap) {
	for _, errList := range errs {
		VerifyDefaultIssueMessages(t, errList)
	}
}

// FindByPath returns all issues with the given path
func FindByPath(errs internals.ZogIssueList, path string) internals.ZogIssueList {
	var result internals.ZogIssueList
	for _, e := range errs {
		if internals.FlattenPath(e.Path) == path {
			result = append(result, e)
		}
	}
	return result
}

// HasPath checks if any issue has the given path
func HasPath(errs internals.ZogIssueList, path string) bool {
	for _, e := range errs {
		if internals.FlattenPath(e.Path) == path {
			return true
		}
	}
	return false
}
