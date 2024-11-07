package internals

import zconst "github.com/Oudwins/zog/zconst"

type ZogParseCtx struct {
	Fmter  ErrFmtFunc
	Errors ZogErrors
	m      map[string]any
}

func (c *ZogParseCtx) NewError(p PathBuilder, e ZogError) {
	c.Fmter(e, c)
	c.Errors.Add(p, e)
}

func (c *ZogParseCtx) HasErrored() bool {
	return !c.Errors.IsEmpty()
}

func (c *ZogParseCtx) SetErrFormatter(fmter ErrFmtFunc) {
	c.Fmter = fmter
}

func (c *ZogParseCtx) Set(key string, val any) {
	if c.m == nil {
		c.m = make(map[string]any)
	}
	c.m[key] = val
}

func (c *ZogParseCtx) Get(key string) any {
	return c.m[key]
}

type ParseCtx interface {
	// Get a value from the context
	Get(key string) any
	// Please don't depend on this interface it may change
	NewError(p PathBuilder, e ZogError)
	// Please don't depend on this interface it may change
	HasErrored() bool
}

func NewParseCtx(errs ZogErrors, fmter ErrFmtFunc) *ZogParseCtx {
	return &ZogParseCtx{
		Fmter:  fmter,
		Errors: errs,
	}
}

// TestFunc is a function that takes the data as input and returns a boolean indicating if it is valid or not
type TestFunc = func(val any, ctx ParseCtx) bool

// Test is a struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.
type Test struct {
	ErrCode      zconst.ZogErrCode
	Params       map[string]any
	ErrFmt       ErrFmtFunc
	ValidateFunc TestFunc
}
