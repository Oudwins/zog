package zog

import (
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

// Deprecated: This will be removed in the future. Use z.ZogIssue instead
// This is a type for the ZogError interface. It is the interface that all errors returned from zog implement.
type ZogError = p.ZogError

// This is a type for the ZogIssue= interface. It is the interface that all errors returned from zog implement.
type ZogIssue = p.ZogError

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

func (i *issueHelpers) SanitizeList(l ZogIssueList) []string {
	errs := make([]string, len(l))
	for i, err := range l {
		errs[i] = err.Message()
	}
	return errs
}

// ! ERRORS -> Deprecated
// Deprecated: This will be removed in the future.
type errHelpers struct {
}

// Deprecated: This will be removed in the future.
// Use z.Issues instead
// Helper struct for dealing with zog errors. Beware this API may change
var Errors = errHelpers{}

// Deprecated: This will be removed in the future.
// Create error from (originValue any, destinationValue any, test *p.Test)
func (e *errHelpers) FromTest(o any, destType zconst.ZogType, t *p.Test, p ParseCtx) ZogIssue {
	er := e.New(t.IssueCode, o, destType, t.Params, "", nil)
	if t.IssueFmtFunc != nil {
		t.IssueFmtFunc(er, p)
	}
	return er
}

// Deprecated: This will be removed in the future.
// Create error from
func (e *errHelpers) FromErr(o any, destType zconst.ZogType, err error) ZogIssue {
	return e.New(zconst.IssueCodeCustom, o, destType, nil, "", err)
}

// Deprecated: This will be removed in the future.
func (e *errHelpers) WrapUnknown(o any, destType zconst.ZogType, err error) ZogIssue {
	zerr, ok := err.(ZogIssue)
	if !ok {
		return e.FromErr(o, destType, err)
	}
	return zerr
}

// Deprecated: This will be removed in the future.
func (e *errHelpers) New(code zconst.ZogIssueCode, o any, destType zconst.ZogType, params map[string]any, msg string, err error) p.ZogError {
	return &p.ZogErr{
		C:       code,
		ParamsM: params,
		Val:     o,
		Typ:     destType,
		Msg:     msg,
		Err:     err,
	}
}

// Deprecated: This will be removed in the future. Use z.Issues.SanitizeMap instead
func (e *errHelpers) SanitizeMap(m p.ZogIssueMap) map[string][]string {
	return Issues.SanitizeMap(m)
}

// Deprecated: This will be removed in the future. Use z.Issues.SanitizeList instead
func (e *errHelpers) SanitizeList(l p.ZogIssueList) []string {
	return Issues.SanitizeList(l)
}

// ! Data Providers

// Deprecated: This will be removed in the future.
// You should just pass your map[string]T to the schema.Parse() function directly without using this:
// old: schema.Parse(z.NewMapDataProvider(m), &dest)
// new: schema.Parse(m, &dest)
func NewMapDataProvider[T any](m map[string]T) p.DataProvider {
	return p.NewMapDataProvider(m)
}
