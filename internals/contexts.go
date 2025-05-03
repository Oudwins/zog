package internals

import (
	"fmt"

	zconst "github.com/Oudwins/zog/zconst"
)

// Zog Context interface. This is the interface that is passed to schema tests, pre and post transforms
type Ctx interface {
	/**
	METHOD YOU ARE FREE TO USE
	*/
	// Get a value from the context
	Get(key string) any
	// Adds an issue to the schema execution.
	AddIssue(e *ZogIssue)

	// Returns a new issue with the current schema context's data prefilled
	/*
		Usage:

		func MyCustomTestFunc(val any, ctx z.Ctx) {
			if reason1 {
			   ctx.AddIssue(ctx.Issue().SetMessage("Reason 1"))
			} else if reason2 {
			   ctx.AddIssue(ctx.Issue().SetMessage("Reason 2"))
			} else {
			   ctx.AddIssue(ctx.Issue().SetMessage("Reason 3"))
			}
		}

	*/
	Issue() *ZogIssue

	/**
	METHOD YOU SHOULD NOT USE
	*/
	// Deprecated: Use Ctx.AddIssue() instead
	// Please don't depend on this interface it may change
	NewError(p *PathBuilder, e *ZogIssue)
	// Please don't depend on this interface it may change
	HasErrored() bool
}

func NewExecCtx(errs ZogIssues, fmter IssueFmtFunc) *ExecCtx {
	c := ExecCtxPool.Get().(*ExecCtx)
	c.Fmter = fmter
	c.Errors = errs
	return c
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
func (c *ExecCtx) AddIssue(e *ZogIssue) {
	if e.Message == "" {
		c.Fmter(e, c)
	}
	c.Errors.Add(e.Path, e)
}

func (c *ExecCtx) Issue() *ZogIssue {
	return NewZogIssue()
}

// Deprecated: Use Ctx.AddIssue() instead
// This is old interface. It will be removed soon
func (c *ExecCtx) NewError(path *PathBuilder, e *ZogIssue) {
	c.Errors.Add(path.String(), e)
}

// Internal. Used to format errors
func (c *ExecCtx) FmtErr(e *ZogIssue) {
	if e.Message != "" {
		return
	}
	c.Fmter(e, c)
}

func (c *ExecCtx) NewSchemaCtx(val any, destPtr any, path *PathBuilder, dtype zconst.ZogType) *SchemaCtx {
	c2 := SchemaCtxPool.Get().(*SchemaCtx)
	c2.ExecCtx = c
	c2.Data = val
	c2.ValPtr = destPtr
	c2.Path = path
	c2.DType = dtype
	c2.CanCatch = false
	c2.HasCaught = false
	c2.Exit = false
	return c2
}

func (c *ExecCtx) NewValidateSchemaCtx(valPtr any, path *PathBuilder, dtype zconst.ZogType) *SchemaCtx {
	c2 := SchemaCtxPool.Get().(*SchemaCtx)
	c2.ExecCtx = c
	c2.Data = nil
	c2.ValPtr = valPtr
	c2.Path = path
	c2.DType = dtype
	c2.CanCatch = false
	c2.HasCaught = false
	c2.Exit = false
	return c2
}

func (c *ExecCtx) Free() {
	ExecCtxPool.Put(c)
}

type SchemaCtx struct {
	*ExecCtx
	Data      any // input data in case of parse
	ValPtr    any // pointer to real output value in both validate & parse
	Path      *PathBuilder
	DType     zconst.ZogType
	CanCatch  bool
	Exit      bool
	HasCaught bool
	Processor any
}

func (c *SchemaCtx) AddIssue(e *ZogIssue) {
	if c.CanCatch {
		c.Exit = true
		FreeIssue(e)
		return
	}
	c.ExecCtx.AddIssue(e)
}

func (c *SchemaCtx) Issue() *ZogIssue {
	// e := ZogIssuePool.Get().(*ZogIssue)
	// e.Code = ""
	// e.Path = c.Path.String()
	// e.Err = nil
	// e.Message = ""
	// e.Params = nil
	// e.Dtype = c.DType
	// e.Value = c.Data
	return NewZogIssue().SetPath(c.Path.String()).SetDType(c.DType).SetValue(c.Data)
}

// Please don't depend on this method it may change
func (c *SchemaCtx) IssueFromTest(test TestInterface, val any) *ZogIssue {
	e := ZogIssuePool.Get().(*ZogIssue)
	e.Code = test.GetIssueCode()
	e.Path = c.Path.String()
	e.Err = nil
	e.Message = ""
	e.Dtype = c.DType
	e.Value = val
	e.Params = test.GetParams()
	if test.GetIssueFmtFunc() != nil {
		test.GetIssueFmtFunc()(e, c)
	}
	if test.GetIssuePath() != "" {
		e.Path = test.GetIssuePath()
	}
	return e
}

// Please don't depend on this method it may change
func (c *SchemaCtx) IssueFromCoerce(err error) *ZogIssue {
	e := ZogIssuePool.Get().(*ZogIssue)
	e.Code = zconst.IssueCodeCoerce
	e.Path = c.Path.String()
	e.Message = ""
	e.Dtype = c.DType
	e.Value = c.Data
	e.Err = err
	return e
}

// Please don't depend on this method it may change
// Wraps an error in a ZogIssue if it is not already a ZogIssue
func (c *SchemaCtx) IssueFromUnknownError(err error) *ZogIssue {
	zerr, ok := err.(*ZogIssue)
	if !ok {
		return c.Issue().SetError(err)
	}
	return zerr
}

// Frees the context to be reused
func (c *SchemaCtx) Free() {
	SchemaCtxPool.Put(c)
}

func (c *SchemaCtx) String() string {
	return fmt.Sprintf("z.Ctx{Data: %v, ValPtr: %v, Path: %v, DType: %v, CanCatch: %v, Exit: %v, HasCaught: %v }", SafeString(c.Data), SafeString(c.ValPtr), c.Path, c.DType, c.CanCatch, c.Exit, c.HasCaught)
}

// func (c *TestCtx) Issue() *ZogIssue {
// 	// TODO handle catch here
// 	zerr := ZogIssuePool.Get().(*ZogIssue)
// 	zerr.Code = c.Test.IssueCode
// 	zerr.Path = c.Path.String()
// 	zerr.Err = nil
// 	zerr.Message = ""
// 	zerr.Params = c.Test.Params
// 	zerr.Dtype = c.DType
// 	zerr.Value = c.Data
// 	return zerr
// }

// func (c *TestCtx) FmtErr(e *ZogIssue) {
// 	if e.Message != "" {
// 		return
// 	}

// 	if c.Test.IssueFmtFunc != nil {
// 		c.Test.IssueFmtFunc(e, c)
// 		return
// 	}

// 	c.SchemaCtx.FmtErr(e)
// }

// func (c *TestCtx) AddIssue(e *ZogIssue) {
// 	c.FmtErr(e)
// 	c.Errors.Add(c.Path.String(), e)
// }
