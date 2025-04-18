package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type PreprocessSchema[F any, T any] struct {
	schema ZogSchema
	fn     func(data F, ctx Ctx) (out T, err error)
}

func Preprocess[F any, T any](fn func(data F, ctx Ctx) (out T, err error), schema ZogSchema) *PreprocessSchema[F, T] {
	return &PreprocessSchema[F, T]{fn: fn, schema: schema}
}

func (s *PreprocessSchema[F, T]) process(ctx *p.SchemaCtx) {
	out, err := s.fn(ctx.Data.(F), ctx)
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(err.Error()))
		return
	}
	ctx.Data = out
	// ptr := ctx.ValPtr.(*T)
	// *ptr = out
	// ptrVal := reflect.ValueOf(ctx.ValPtr).Elem()
	// ptrVal.Set(reflect.ValueOf(out))
	// ctx.ValPtr = out
	s.schema.process(ctx)
}

func (s *PreprocessSchema[F, T]) validate(ctx *p.SchemaCtx) {
	out, err := s.fn(ctx.ValPtr.(F), ctx)
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(err.Error()))
		return
	}
	ptr := ctx.ValPtr.(*T)
	*ptr = out
	// if out != ctx.ValPtr {
	// 	reflect.ValueOf(ctx.ValPtr).Elem().Set(reflect.ValueOf(out))
	// }
	s.schema.validate(ctx)
}

func (s *PreprocessSchema[F, T]) getType() zconst.ZogType {
	return s.schema.getType()
}

func (s *PreprocessSchema[F, T]) Parse(data F, destPtr *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, destPtr, path, s.getType())
	defer sctx.Free()
	s.process(sctx)
	return errs.List
}

func (s *PreprocessSchema[F, T]) Validate(data *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, data, path, s.getType())
	defer sctx.Free()
	s.validate(sctx)
	return errs.List
}
