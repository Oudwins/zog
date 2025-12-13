package zog

import (
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/stretchr/testify/assert"
)

// ============================================
// Tests for GroupByFlattenedPath
// ============================================

func TestGroupByFlattenedPath_EmptyList(t *testing.T) {
	issues := ZogIssueList{}
	result := Issues.GroupByFlattenedPath(issues)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestGroupByFlattenedPath_SingleIssueSimplePath(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"name"}, Message: "is required"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 1)
	assert.Len(t, result["name"], 1)
	assert.Equal(t, "is required", result["name"][0].Message)
}

func TestGroupByFlattenedPath_MultipleIssuesSamePath(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"email"}, Message: "is required"},
		{Path: []string{"email"}, Message: "must be valid"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 1)
	assert.Len(t, result["email"], 2)
	assert.Equal(t, "is required", result["email"][0].Message)
	assert.Equal(t, "must be valid", result["email"][1].Message)
}

func TestGroupByFlattenedPath_MultipleIssuesDifferentPaths(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"name"}, Message: "is required"},
		{Path: []string{"email"}, Message: "must be valid"},
		{Path: []string{"age"}, Message: "must be positive"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 3)
	assert.Len(t, result["name"], 1)
	assert.Len(t, result["email"], 1)
	assert.Len(t, result["age"], 1)
}

func TestGroupByFlattenedPath_NestedPath(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"user", "address", "city"}, Message: "is required"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 1)
	assert.Contains(t, result, "user.address.city")
	assert.Equal(t, "is required", result["user.address.city"][0].Message)
}

func TestGroupByFlattenedPath_EmptyPath(t *testing.T) {
	issues := ZogIssueList{
		{Path: nil, Message: "root error"},
		{Path: []string{}, Message: "another root error"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	// Both nil and empty slice paths should result in empty string key
	assert.Len(t, result, 1)
	assert.Contains(t, result, "")
	assert.Len(t, result[""], 2)
}

func TestGroupByFlattenedPath_MixedPaths(t *testing.T) {
	issues := ZogIssueList{
		{Path: nil, Message: "root error"},
		{Path: []string{"name"}, Message: "name error"},
		{Path: []string{"user", "email"}, Message: "email error"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 3)
	assert.Contains(t, result, "")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "user.email")
}

func TestGroupByFlattenedPath_ArrayIndexPath(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"users", "[0]", "name"}, Message: "is required"},
		{Path: []string{"users", "[1]", "name"}, Message: "too short"},
	}
	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 2)
	assert.Contains(t, result, "users[0].name")
	assert.Contains(t, result, "users[1].name")
}

// ============================================
// Tests for Treeify
// ============================================

func TestTreeify_EmptyList(t *testing.T) {
	issues := ZogIssueList{}
	result := Issues.Treeify(issues)

	assert.NotNil(t, result)
	assert.Contains(t, result, "errors")
	assert.Contains(t, result, "properties")
	assert.Empty(t, result["errors"])
	assert.Empty(t, result["properties"])
}

func TestTreeify_RootLevelError(t *testing.T) {
	issues := ZogIssueList{
		{Path: nil, Message: "validation failed"},
	}
	result := Issues.Treeify(issues)

	errors := result["errors"].([]string)
	assert.Len(t, errors, 1)
	assert.Equal(t, "validation failed", errors[0])
}

func TestTreeify_MultipleRootLevelErrors(t *testing.T) {
	issues := ZogIssueList{
		{Path: nil, Message: "error one"},
		{Path: []string{}, Message: "error two"},
	}
	result := Issues.Treeify(issues)

	errors := result["errors"].([]string)
	assert.Len(t, errors, 2)
	assert.Contains(t, errors, "error one")
	assert.Contains(t, errors, "error two")
}

func TestTreeify_SinglePropertyError(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"name"}, Message: "is required"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	assert.Contains(t, properties, "name")

	nameProp := properties["name"].(map[string]any)
	errors := nameProp["errors"].([]string)
	assert.Len(t, errors, 1)
	assert.Equal(t, "is required", errors[0])
}

func TestTreeify_NestedPropertyError(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"user", "name"}, Message: "too short"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	assert.Contains(t, properties, "user")

	userProp := properties["user"].(map[string]any)
	assert.Contains(t, userProp, "name")

	nameProp := userProp["name"].(map[string]any)
	errors := nameProp["errors"].([]string)
	assert.Len(t, errors, 1)
	assert.Equal(t, "too short", errors[0])
}

func TestTreeify_ArrayIndexError(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"users", "[0]"}, Message: "invalid user"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	assert.Contains(t, properties, "users")

	usersProp := properties["users"].(map[string]any)
	assert.Contains(t, usersProp, "items")

	items := usersProp["items"].([]any)
	assert.Len(t, items, 1)

	item0 := items[0].(map[string]any)
	errors := item0["errors"].([]string)
	assert.Len(t, errors, 1)
	assert.Equal(t, "invalid user", errors[0])
}

func TestTreeify_ArrayIndexWithNestedProperty(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"users", "[0]", "email"}, Message: "invalid email"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	usersProp := properties["users"].(map[string]any)
	items := usersProp["items"].([]any)
	item0 := items[0].(map[string]any)

	assert.Contains(t, item0, "email")
	emailProp := item0["email"].(map[string]any)
	errors := emailProp["errors"].([]string)
	assert.Len(t, errors, 1)
	assert.Equal(t, "invalid email", errors[0])
}

func TestTreeify_MultipleArrayIndices(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"users", "[0]", "name"}, Message: "error on first"},
		{Path: []string{"users", "[2]", "name"}, Message: "error on third"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	usersProp := properties["users"].(map[string]any)
	items := usersProp["items"].([]any)

	// Should have 3 items (indices 0, 1, 2), with index 1 being nil
	assert.Len(t, items, 3)
	assert.NotNil(t, items[0])
	assert.Nil(t, items[1])
	assert.NotNil(t, items[2])

	item0 := items[0].(map[string]any)
	name0 := item0["name"].(map[string]any)
	assert.Equal(t, "error on first", name0["errors"].([]string)[0])

	item2 := items[2].(map[string]any)
	name2 := item2["name"].(map[string]any)
	assert.Equal(t, "error on third", name2["errors"].([]string)[0])
}

func TestTreeify_MultipleErrorsSamePath(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"email"}, Message: "is required"},
		{Path: []string{"email"}, Message: "must be valid"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	emailProp := properties["email"].(map[string]any)
	errors := emailProp["errors"].([]string)

	assert.Len(t, errors, 2)
	assert.Contains(t, errors, "is required")
	assert.Contains(t, errors, "must be valid")
}

func TestTreeify_MixedRootAndNestedErrors(t *testing.T) {
	issues := ZogIssueList{
		{Path: nil, Message: "root error"},
		{Path: []string{"name"}, Message: "name error"},
		{Path: []string{"user", "email"}, Message: "email error"},
	}
	result := Issues.Treeify(issues)

	// Check root errors
	rootErrors := result["errors"].([]string)
	assert.Len(t, rootErrors, 1)
	assert.Equal(t, "root error", rootErrors[0])

	// Check nested errors
	properties := result["properties"].(map[string]any)
	assert.Contains(t, properties, "name")
	assert.Contains(t, properties, "user")

	nameProp := properties["name"].(map[string]any)
	assert.Equal(t, "name error", nameProp["errors"].([]string)[0])

	userProp := properties["user"].(map[string]any)
	emailProp := userProp["email"].(map[string]any)
	assert.Equal(t, "email error", emailProp["errors"].([]string)[0])
}

func TestTreeify_NumericStringPath(t *testing.T) {
	// Test that numeric strings are treated as array indices
	issues := ZogIssueList{
		{Path: []string{"items", "0", "name"}, Message: "error"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	itemsProp := properties["items"].(map[string]any)
	items := itemsProp["items"].([]any)

	assert.Len(t, items, 1)
	item0 := items[0].(map[string]any)
	nameProp := item0["name"].(map[string]any)
	assert.Equal(t, "error", nameProp["errors"].([]string)[0])
}

func TestTreeify_DeeplyNestedStructure(t *testing.T) {
	issues := ZogIssueList{
		{Path: []string{"a", "b", "c", "d"}, Message: "deep error"},
	}
	result := Issues.Treeify(issues)

	properties := result["properties"].(map[string]any)
	a := properties["a"].(map[string]any)
	b := a["b"].(map[string]any)
	c := b["c"].(map[string]any)
	d := c["d"].(map[string]any)

	errors := d["errors"].([]string)
	assert.Len(t, errors, 1)
	assert.Equal(t, "deep error", errors[0])
}

func TestTreeify_ComplexMixedStructure(t *testing.T) {
	issues := ZogIssueList{
		{Path: nil, Message: "form invalid"},
		{Path: []string{"users", "[0]", "name"}, Message: "name required"},
		{Path: []string{"users", "[0]", "email"}, Message: "email invalid"},
		{Path: []string{"users", "[1]", "name"}, Message: "name too short"},
		{Path: []string{"settings", "notifications"}, Message: "must be boolean"},
	}
	result := Issues.Treeify(issues)

	// Root error
	rootErrors := result["errors"].([]string)
	assert.Contains(t, rootErrors, "form invalid")

	properties := result["properties"].(map[string]any)

	// Check users array
	usersProp := properties["users"].(map[string]any)
	items := usersProp["items"].([]any)
	assert.Len(t, items, 2)

	// First user
	user0 := items[0].(map[string]any)
	assert.Equal(t, "name required", user0["name"].(map[string]any)["errors"].([]string)[0])
	assert.Equal(t, "email invalid", user0["email"].(map[string]any)["errors"].([]string)[0])

	// Second user
	user1 := items[1].(map[string]any)
	assert.Equal(t, "name too short", user1["name"].(map[string]any)["errors"].([]string)[0])

	// Check settings
	settingsProp := properties["settings"].(map[string]any)
	notifProp := settingsProp["notifications"].(map[string]any)
	assert.Equal(t, "must be boolean", notifProp["errors"].([]string)[0])
}

// ============================================
// Tests via public API (Issues helper)
// ============================================

func TestIssuesHelper_GroupByFlattenedPath(t *testing.T) {
	// Integration test using the public API
	issues := p.ZogIssueList{
		{Path: []string{"user", "name"}, Message: "error1"},
		{Path: []string{"user", "name"}, Message: "error2"},
		{Path: []string{"user", "email"}, Message: "error3"},
	}

	result := Issues.GroupByFlattenedPath(issues)

	assert.Len(t, result, 2)
	assert.Len(t, result["user.name"], 2)
	assert.Len(t, result["user.email"], 1)
}

func TestIssuesHelper_Treeify(t *testing.T) {
	// Integration test using the public API
	issues := p.ZogIssueList{
		{Path: nil, Message: "root"},
		{Path: []string{"field"}, Message: "field error"},
	}

	result := Issues.Treeify(issues)

	assert.NotNil(t, result)
	rootErrors := result["errors"].([]string)
	assert.Contains(t, rootErrors, "root")

	properties := result["properties"].(map[string]any)
	fieldProp := properties["field"].(map[string]any)
	assert.Contains(t, fieldProp["errors"].([]string), "field error")
}
