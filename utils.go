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

// This is a type for the ZogErrList type. It is a list of ZogIssues returned from parsing primitive schemas. The type is []ZogIssue
type ZogIssueList = p.ZogIssueList

// This is a type for the ZogIssueMap type. It is a map[string][]ZogIssue returned from parsing complex schemas. The type is map[string][]ZogIssue
// All errors are returned in a flat map, not matter how deep the schema is. For example:
/*
schema := z.Struct(z.Shape{
  "address": z.Struct(z.Shape{
    "street": z.String().Min(3).Max(10),
    "city": z.String().Min(3).Max(10),
  }),
  "fields": z.Slice(z.String().Min(3).Max(10)),
})
errors = map[string][]ZogIssue{
  "address.street": []ZogIssue{....}, // error for the street field in the address struct
  "fields[0]": []ZogIssue{...}, // error for the first field in the slice
}
*/
type ZogIssueMap = p.ZogIssueMap

type CoercerFunc = conf.CoercerFunc

// ! TESTS

// Test is the test object. It is the struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.
type Test[T any] p.Test[T]

type issueHelpers struct {
}

var Issues = issueHelpers{}

// SanitizeMap returns a map of issue messages for each key in the map. It keeps only the issue messages and strips out any other issue data.
func (i *issueHelpers) SanitizeMap(m ZogIssueMap) map[string][]string {
	errs := make(map[string][]string, len(m))
	for k, v := range m {
		errs[k] = i.SanitizeList(v)
	}
	return errs
}

// Suger function that does both sanitize and collect.
func (i issueHelpers) SanitizeMapAndCollect(m ZogIssueMap) map[string][]string {
	errs := i.SanitizeMap(m)
	i.CollectMap(m)
	return errs
}

// SanitizeList returns a slice of issue messages for each issue in the list. It keeps only the issue messages and strips out any other issue data.
func (i *issueHelpers) SanitizeList(l ZogIssueList) []string {
	errs := make([]string, len(l))
	for i, err := range l {
		errs[i] = err.Message
	}
	return errs
}

// Suger function that does both sanitize and collect.
func (i issueHelpers) SanitizeListAndCollect(l ZogIssueList) []string {
	errs := i.SanitizeList(l)
	i.CollectList(l)
	return errs
}

// Collects a ZogIssueMap to be reused by Zog. This will "free" the issues in the map. This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) CollectMap(issues ZogIssueMap) {
	for _, list := range issues {
		i.CollectList(list)
	}
}

// Collects a ZogIssueList to be reused by Zog. This will "free" the issues in the list. This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) CollectList(issues ZogIssueList) {
	for _, iss := range issues {
		i.Collect(iss)
	}
}

// Collects a ZogIssue to be reused by Zog. This will "free" the issue. This can help make Zog more performant by reusing issue structs.
func (i *issueHelpers) Collect(issue *ZogIssue) {
	p.FreeIssue(issue)
}
