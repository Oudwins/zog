package internals

func Flatten(issues ZogIssueList) map[string][]string {
	flattened := make(map[string][]string)
	for _, issue := range issues {
		flattened[issue.Path] = append(flattened[issue.Path], issue.Message)
	}
	return flattened
}
