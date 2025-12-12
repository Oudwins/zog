package internals

import (
	"github.com/Oudwins/zog/zconst"
)

func Flatten(issues ZogIssueList) map[string][]string {
	flattened := make(map[string][]string)
	for _, issue := range issues {
		path := issue.Path
		if path == "" {
			path = zconst.ISSUE_KEY_ROOT
		}
		flattened[path] = append(flattened[path], issue.Message)
	}
	return flattened
}

func GroupByFlattenedPath(issues ZogIssueList) map[string]ZogIssueList {
	flattened := make(map[string]ZogIssueList)
	for _, issue := range issues {
		flattened[issue.Path] = append(flattened[issue.Path], issue)
	}
	return flattened
}
