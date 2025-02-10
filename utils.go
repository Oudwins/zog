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

// This is a type for the ZogError interface. It is the interface that all errors returned from zog implement.
type ZogError = p.ZogError

type ZogIssue = p.ZogError

// This is a type for the ZogErrList type. It is a list of ZogErrors returned from parsing primitive schemas. The type is []ZogError
type ZogErrList = p.ZogErrList

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
type ZogErrMap = p.ZogErrMap

// ! TESTS

// Test is the test object. It is the struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.
type Test = p.Test

// TestFunc is a helper function to define a custom test. It takes the error code which will be used for the error message and a validate function. Usage:
//
//	schema.Test(z.TestFunc(zconst.ErrCodeCustom, func(val any, ctx ParseCtx) bool {
//		return val == "hello"
//	}))
func TestFunc(errCode zconst.ZogErrCode, validateFunc p.TestFunc) p.Test {
	t := p.Test{
		ErrCode:      errCode,
		ValidateFunc: validateFunc,
	}
	return t
}

// ! ERRORS
// Deprecated: This will be removed in the future.
type errHelpers struct {
}

// Deprecated: This will be removed in the future.
// Helper struct for dealing with zog errors. Beware this API may change
var Errors = errHelpers{}

// Deprecated: This will be removed in the future.
// Create error from (originValue any, destinationValue any, test *p.Test)
func (e *errHelpers) FromTest(o any, destType zconst.ZogType, t *p.Test, p ParseCtx) p.ZogError {
	er := e.New(t.ErrCode, o, destType, t.Params, "", nil)
	if t.ErrFmt != nil {
		t.ErrFmt(er, p)
	}
	return er
}

// Deprecated: This will be removed in the future.
// Create error from
func (e *errHelpers) FromErr(o any, destType zconst.ZogType, err error) p.ZogError {
	return e.New(zconst.ErrCodeCustom, o, destType, nil, "", err)
}

// Deprecated: This will be removed in the future.
func (e *errHelpers) WrapUnknown(o any, destType zconst.ZogType, err error) p.ZogError {
	zerr, ok := err.(p.ZogError)
	if !ok {
		return e.FromErr(o, destType, err)
	}
	return zerr
}

// Deprecated: This will be removed in the future.
func (e *errHelpers) New(code zconst.ZogErrCode, o any, destType zconst.ZogType, params map[string]any, msg string, err error) p.ZogError {
	return &p.ZogErr{
		C:       code,
		ParamsM: params,
		Val:     o,
		Typ:     destType,
		Msg:     msg,
		Err:     err,
	}
}

func (e *errHelpers) SanitizeMap(m p.ZogErrMap) map[string][]string {
	errs := make(map[string][]string, len(m))
	for k, v := range m {
		errs[k] = e.SanitizeList(v)
	}
	return errs
}

func (e *errHelpers) SanitizeList(l p.ZogErrList) []string {
	errs := make([]string, len(l))
	for i, err := range l {
		errs[i] = err.Message()
	}
	return errs
}

// ! Data Providers

// Deprecated: This will be removed in the future.
// You should just pass your map[string]T to the schema.Parse() function directly without using this:
// old: schema.Parse(z.NewMapDataProvider(m), &dest)
// new: schema.Parse(m, &dest)
func NewMapDataProvider[T any](m map[string]T) p.DataProvider {
	return p.NewMapDataProvider(m)
}
