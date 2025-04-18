package zog

import (
	"fmt"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

func defaultCustomCoercer[T any](data any) (any, error) {
	v, ok := data.(T)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", new(T), data)
	}
	return v, nil
}

type Custom[T any] struct {
	test Test
}

func CustomFunc[T any](fn func(ptr *T, ctx Ctx) bool, opts ...TestOption) *Custom[T] {
	test := &Test{}
	p.TestFuncFromBool(func(val any, ctx Ctx) bool {
		return fn(val.(*T), ctx)
	}, test)
	for _, opt := range opts {
		opt(test)
	}
	return &Custom[T]{test: *test}
}

func (c *Custom[T]) Parse(data any, destPtr *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, destPtr, path, c.getType())
	defer sctx.Free()
	c.process(sctx)
	return errs.List
}

func (c *Custom[T]) process(ctx *p.SchemaCtx) {
	ctx.Test = &c.test

	// set the value
	d, ok := ctx.Data.(T)
	if !ok {
		ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("expected %T, got %T", new(T), ctx.Data)))
		return
	}
	ptr := ctx.ValPtr.(*T)
	*ptr = d

	// run the test
	c.test.Func(ptr, ctx)
}

func (c *Custom[T]) Validate(dataPtr *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(dataPtr, dataPtr, path, c.getType())
	defer sctx.Free()
	c.validate(sctx)
	return errs.List
}

func (c *Custom[T]) validate(ctx *p.SchemaCtx) {
	ctx.Test = &c.test
	c.test.Func(ctx.ValPtr, ctx)
}

func (c *Custom[T]) setCoercer(coercer CoercerFunc) {
	// no op
}

func (c *Custom[T]) getType() zconst.ZogType {
	return "custom"
}
