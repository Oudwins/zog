package zog

import (
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

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
type errHelpers struct {
}

// Helper struct for dealing with zog errors. Beware this API may change
var Errors = errHelpers{}

// Create error from (originValue any, destinationValue any, test *p.Test)
func (e *errHelpers) FromTest(o any, destType zconst.ZogType, t *p.Test, p ParseCtx) p.ZogError {
	er := e.New(t.ErrCode, o, destType, t.Params, "", nil)
	if t.ErrFmt != nil {
		t.ErrFmt(er, p)
	}
	return er
}

// Create error from
func (e *errHelpers) FromErr(o any, destType zconst.ZogType, err error) p.ZogError {
	return e.New(zconst.ErrCodeCustom, o, destType, nil, "", err)
}

func (e *errHelpers) WrapUnknown(o any, destType zconst.ZogType, err error) p.ZogError {
	zerr, ok := err.(p.ZogError)
	if !ok {
		return e.FromErr(o, destType, err)
	}
	return zerr
}

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
