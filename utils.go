package zog

import (
	"reflect"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// ! Passing Types through

// Deprecated: Use z.Ctx instead.
// ParseCtx will be removed in the future since it is used for both validation and parsing and is a confusing name.
type ParseCtx = p.Ctx

// This is the context that is passed through an entire execution of `schema.Parse()` or `schema.Validate()`.
// You can use it to pass a key/value for a specific execution. More about context in the [docs](https://zog.dev/context)
type Ctx = p.Ctx

// This is a type for the ZogIssue type. It is the type of all the errors returned from zog.
type ZogIssue = p.ZogIssue

// Deprecated: This will be removed in the future. Use z.ZogIssueList instead
// This is a type for the ZogErrList type. It is a list of ZogIssues returned from parsing primitive schemas. The type is []ZogError
type ZogErrList = p.ZogIssueList

// This is a type for the ZogErrList type. It is a list of ZogIssues returned from parsing primitive schemas. The type is []ZogIssue
type ZogIssueList = p.ZogIssueList

// Deprecated: This will be removed in the future. Use z.ZogIssueMap instead
// This is a type for the ZogErrMap type. It is a map[string][]ZogError returned from parsing complex schemas. The type is map[string][]ZogError
// All errors are returned in a flat map, not matter how deep the schema is. For example:
/*
schema := z.Struct(z.Schema{
  "address": z.Struct(z.Schema{
    "street": z.String().Min(3).Max(10),
    "city": z.String().Min(3).Max(10),
  }),
  "fields": z.Slice(z.String().Min(3).Max(10)),
})
errors = map[string][]ZogError{
  "address.street": []ZogError{....}, // error for the street field in the address struct
  "fields[0]": []ZogError{...}, // error for the first field in the slice
}

*/
type ZogErrMap = p.ZogIssueMap

// This is a type for the ZogIssueMap type. It is a map[string][]ZogIssue returned from parsing complex schemas. The type is map[string][]ZogIssue
// All errors are returned in a flat map, not matter how deep the schema is. For example:
/*
schema := z.Struct(z.Schema{
  "address": z.Struct(z.Schema{
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

// ! TESTS

// Test is the test object. It is the struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.
type Test = p.Test

// TestFunc is a helper function to define a custom test. It takes the error code which will be used for the error message and a validate function. Usage:
//
//	schema.Test(z.TestFunc(zconst.IssueCodeCustom, func(val any, ctx ParseCtx) bool {
//		return val == "hello"
//	}))
func TestFunc(IssueCode zconst.ZogIssueCode, validateFunc p.TestFunc) p.Test {
	t := p.Test{
		IssueCode:    IssueCode,
		ValidateFunc: validateFunc,
	}
	return t
}

type issueHelpers struct {
}

var Issues = issueHelpers{}

func (i *issueHelpers) SanitizeMap(m ZogIssueMap) map[string][]string {
	errs := make(map[string][]string, len(m))
	for k, v := range m {
		errs[k] = i.SanitizeList(v)
	}
	return errs
}

func (i issueHelpers) SanitizeMapAndCollect(m ZogIssueMap) map[string][]string {
	errs := i.SanitizeMap(m)
	i.CollectMap(m)
	return errs
}

func (i *issueHelpers) SanitizeList(l ZogIssueList) []string {
	errs := make([]string, len(l))
	for i, err := range l {
		errs[i] = err.Message
	}
	return errs
}

func (i *issueHelpers) SanitizeListAndCollect(l ZogIssueList) []string {
	errs := i.SanitizeList(l)
	i.CollectList(l)
	return errs
}

func (i *issueHelpers) CollectMap(issues ZogIssueMap) {
	for _, list := range issues {
		i.CollectList(list)
	}
}

func (i *issueHelpers) CollectList(issues ZogIssueList) {
	for _, err := range issues {
		err.Free()
	}
}

// ! Data Providers

// Deprecated: This will be removed in the future.
// You should just pass your map[string]T to the schema.Parse() function directly without using this:
// old: schema.Parse(z.NewMapDataProvider(m), &dest)
// new: schema.Parse(m, &dest)
func NewMapDataProvider[T any](m map[string]T) p.DataProvider {
	return p.NewMapDataProvider(m)
}

// Backwards Compatibility

func customTestBackwardsCompatWrapper(testFunc p.TestFunc) func(val any, ctx Ctx) bool {
	return func(val any, ctx Ctx) bool {
		refVal := reflect.ValueOf(val).Elem()
		return testFunc(refVal.Interface(), ctx)
	}
}
