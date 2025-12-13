package internals

import (
	"strconv"
	"strings"

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
	result := map[string]any{
		"errors":     []string{},
		"properties": map[string]any{},
	}

	for _, issue := range issues {
		// Root level errors (empty or nil path)
		if len(issue.Path) == 0 {
			errors := result["errors"].([]string)
			result["errors"] = append(errors, issue.Message)
			continue
		}

		// Navigate to the target location in the tree
		properties := result["properties"].(map[string]any)
		current := any(properties)

		// Process all path segments except the last one to build the structure
		for i := 0; i < len(issue.Path); i++ {
			segment := issue.Path[i]
			idx, isArrayIndex := parseArrayIndex(segment)

			if isArrayIndex {
				// This is an array index
				// Ensure current is a map with "items" key
				currentMap, ok := current.(map[string]any)
				if !ok {
					// This shouldn't happen in normal flow, but handle it
					break
				}

				items, exists := currentMap["items"]
				if !exists {
					items = []any{}
					currentMap["items"] = items
				}

				itemsSlice := items.([]any)
				// Ensure the slice is large enough
				for len(itemsSlice) <= idx {
					itemsSlice = append(itemsSlice, nil)
				}

				// If this is the last segment, create error structure
				if i == len(issue.Path)-1 {
					item := itemsSlice[idx]
					if item == nil {
						item = map[string]any{"errors": []string{}}
						itemsSlice[idx] = item
					}
					itemMap := item.(map[string]any)
					errors := itemMap["errors"].([]string)
					itemMap["errors"] = append(errors, issue.Message)
				} else {
					// Not the last segment, continue navigating
					item := itemsSlice[idx]
					if item == nil {
						item = map[string]any{"errors": []string{}}
						itemsSlice[idx] = item
					}
					current = item
				}

				currentMap["items"] = itemsSlice
			} else {
				// This is a property name
				currentMap, ok := current.(map[string]any)
				if !ok {
					// This shouldn't happen in normal flow, but handle it
					break
				}

				// If this is the last segment, add error to this property
				if i == len(issue.Path)-1 {
					prop, exists := currentMap[segment]
					if !exists {
						prop = map[string]any{"errors": []string{}}
						currentMap[segment] = prop
					}
					propMap := prop.(map[string]any)
					errors := propMap["errors"].([]string)
					propMap["errors"] = append(errors, issue.Message)
				} else {
					// Not the last segment, continue navigating
					prop, exists := currentMap[segment]
					if !exists {
						prop = map[string]any{"errors": []string{}}
						currentMap[segment] = prop
					}
					current = prop
				}
			}
		}
	}

	return result
}

// parseArrayIndex attempts to parse a path segment as an array index.
// Returns (index, true) if it's an array index, (0, false) otherwise.
// Handles both numeric strings (like "1") and bracket notation (like "[1]").
func parseArrayIndex(segment string) (int, bool) {
	// Check if it's bracket notation like "[1]"
	if strings.HasPrefix(segment, "[") && strings.HasSuffix(segment, "]") {
		idxStr := segment[1 : len(segment)-1]
		idx, err := strconv.Atoi(idxStr)
		if err == nil {
			return idx, true
		}
		return 0, false
	}

	// Check if it's a plain numeric string
	idx, err := strconv.Atoi(segment)
	if err == nil {
		return idx, true
	}

	return 0, false
}
