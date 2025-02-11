package internals

import zconst "github.com/Oudwins/zog/zconst"

// Zog Context interface. This is the interface that is passed to schema tests, pre and post transforms
type Ctx interface {
	// Get a value from the context
	Get(key string) any
	// Deprecated: Use Ctx.AddIssue() instead
	// Please don't depend on this interface it may change
	NewError(p PathBuilder, e ZogIssue)
	// Adds an issue to the schema execution.
	AddIssue(e ZogIssue)
	// Please don't depend on this interface it may change
	HasErrored() bool
}

func NewExecCtx(errs ZogIssues, fmter IssueFmtFunc) *ExecCtx {
	return &ExecCtx{
		Fmter:  fmter,
		Errors: errs,
	}
}

type ExecCtx struct {
	Fmter  IssueFmtFunc
	Errors ZogIssues
	m      map[string]any
}

func (c *ExecCtx) HasErrored() bool {
	return !c.Errors.IsEmpty()
}

func (c *ExecCtx) SetIssueFormatter(fmter IssueFmtFunc) {
	c.Fmter = fmter
}

func (c *ExecCtx) Set(key string, val any) {
	if c.m == nil {
		c.m = make(map[string]any)
	}
	c.m[key] = val
}

func (c *ExecCtx) Get(key string) any {
	return c.m[key]
}

// Adds a ZogIssue to the execution context.
func (c *ExecCtx) AddIssue(e ZogIssue) {
	if e.Message() == "" {
		c.Fmter(e, c)
	}
	c.Errors.Add(e.Path(), e)
}

// Deprecated: Use Ctx.AddIssue() instead
// This is old interface. It will be removed soon
func (c *ExecCtx) NewError(path PathBuilder, e ZogIssue) {
	c.Errors.Add(path.String(), e)
}

// Internal. Used to format errors
func (c *ExecCtx) FmtErr(e ZogIssue) {
	if e.Message() != "" {
		return
	}
	c.Fmter(e, c)
}

func (c *ExecCtx) NewSchemaCtx(val any, destPtr any, path PathBuilder, dtype zconst.ZogType) *SchemaCtx {
	return &SchemaCtx{
		ExecCtx: c,
		Val:     val,
		DestPtr: destPtr,
		Path:    path,
		DType:   dtype,
	}
}

func (c *ExecCtx) NewValidateSchemaCtx(valPtr any, path PathBuilder, dtype zconst.ZogType) *SchemaCtx {
	return &SchemaCtx{
		ExecCtx: c,
		Val:     valPtr,
		Path:    path,
		DType:   dtype,
	}
}

type SchemaCtx struct {
	*ExecCtx
	Val     any
	DestPtr any
	Path    PathBuilder
	DType   zconst.ZogType
}
type TestCtx struct {
	*SchemaCtx
	Test *Test
}

func (c *SchemaCtx) Issue() ZogIssue {
	// TODO handle catch here
	return &ZogErr{
		EPath: c.Path.String(),
		Typ:   c.DType,
		Val:   c.Val,
	}
}

// Please don't depend on this method it may change
func (c *SchemaCtx) IssueFromTest(test *Test, val any) ZogIssue {
	e := &ZogErr{
		EPath:   c.Path.String(),
		Typ:     c.DType,
		Val:     val,
		C:       test.IssueCode,
		ParamsM: test.Params,
	}
	if test.IssueFmtFunc != nil {
		test.IssueFmtFunc(e, c)
	}
	return e
}

// Please don't depend on this method it may change
func (c *SchemaCtx) IssueFromCoerce(err error) ZogIssue {
	return &ZogErr{
		C:     zconst.IssueCodeCoerce,
		EPath: c.Path.String(),
		Typ:   c.DType,
		Val:   c.Val,
		Err:   err,
	}
}

// Please don't depend on this method it may change
// Wraps an error in a ZogIssue if it is not already a ZogIssue
func (c *SchemaCtx) IssueFromUnknownError(err error) ZogIssue {
	zerr, ok := err.(ZogIssue)
	if !ok {
		return c.Issue().SetError(err)
	}
	return zerr
}

func (c *TestCtx) Issue() ZogIssue {
	// TODO handle catch here
	return &ZogErr{
		EPath:   c.Path.String(),
		Typ:     c.DType,
		Val:     c.Val,
		C:       c.Test.IssueCode,
		ParamsM: c.Test.Params,
	}
}

func (c *TestCtx) FmtErr(e ZogIssue) {
	if e.Message() != "" {
		return
	}

	if c.Test.IssueFmtFunc != nil {
		c.Test.IssueFmtFunc(e, c)
		return
	}

	c.SchemaCtx.FmtErr(e)
}

func (c *TestCtx) AddIssue(e ZogIssue) {
	c.FmtErr(e)
	c.Errors.Add(c.Path.String(), e)
}
