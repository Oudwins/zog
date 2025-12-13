package internals

import (
	"github.com/Oudwins/zog/zconst"
)

func Flatten(issues ZogIssueList) map[string][]string {
	flattened := make(map[string][]string)
	for _, issue := range issues {
		path := FlattenPath(issue.Path)
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
		path := FlattenPath(issue.Path)
		flattened[path] = append(flattened[path], issue)
	}
	return flattened
}

func Treeify(issues ZogIssueList) map[string]any {
	// response := make(map[string][]string)
	// for _, issue := range issues {

	// }
	return nil
}
