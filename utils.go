package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
)

// ! Passing Types through

// This is the context that is passed through an entire execution of `schema.Parse()` or `schema.Validate()`.
// You can use it to pass a key/value for a specific execution. More about context in the [docs](https://zog.dev/context)
type Ctx = p.Ctx

// This is a type for the ZogIssue type. It is the type of all the errors returned from zog.
type ZogIssue = p.ZogIssue

// ZogIssueList is the type returned by all schema Parse/Validate operations.
// It is a slice of pointers to ZogIssue.
type ZogIssueList = p.ZogIssueList

type CoercerFunc = conf.CoercerFunc

// ! TESTS

// Test is the test object. It is the struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.
type Test[T any] p.Test[T]

type issueHelpers struct {
}

var Issues = issueHelpers{}

func (i *issueHelpers) Flatten(issues ZogIssueList) map[string][]string {
	return p.Flatten(issues)
}

func (i *issueHelpers) FlattenAndCollect(issues ZogIssueList) map[string][]string {
	flattened := i.Flatten(issues)
	i.Collect(issues)
	return flattened
}

func (i *issueHelpers) GroupByFlattenedPath(issues ZogIssueList) map[string]ZogIssueList {
	return p.GroupByFlattenedPath(issues)
}

// Collect returns issues to the pool for reuse.
// This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) Collect(issues ZogIssueList) {
	for _, iss := range issues {
		i.CollectOne(iss)
	}
}

// CollectOne returns a single issue to the pool for reuse.
func (i *issueHelpers) CollectOne(issue *ZogIssue) {
	p.FreeIssue(issue)
}

// =========== DEPRECATED METHODS ===========

// Deprecated: Use Flatten instead. SanitizeList is kept for backward compatibility.
func (i *issueHelpers) SanitizeList(l ZogIssueList) []string {
	return i.Sanitize(l)
}

// Deprecated: Use FlattenAndCollect instead.
func (i issueHelpers) SanitizeListAndCollect(l ZogIssueList) []string {
	return i.SanitizeAndCollect(l)
}

// Deprecated: Use Collect instead.
func (i *issueHelpers) CollectList(issues ZogIssueList) {
	i.Collect(issues)
}

// Deprecated: Use flatten instead or write this function yourself
// Sanitize returns a slice of issue messages from a ZogIssueList.
// This is the primary sanitization method now that all schemas return ZogIssueList.
func (i *issueHelpers) Sanitize(l ZogIssueList) []string {
	errs := make([]string, len(l))
	for idx, err := range l {
		errs[idx] = err.Message
	}
	return errs
}

// Deprecated: Use FlattenAndCollect instead
// SanitizeAndCollect sanitizes the issues and returns them to the pool for reuse.
func (i issueHelpers) SanitizeAndCollect(l ZogIssueList) []string {
	errs := i.Sanitize(l)
	i.Collect(l)
	return errs
}
