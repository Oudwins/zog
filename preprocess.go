package zog

import (
	"fmt"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type PreprocessSchema[F any, T any] struct {
	schema ZogSchema
	fn     func(data F, ctx Ctx) (out T, err error)
}

// out should never be a pointer type
func Preprocess[F any, T any](fn func(data F, ctx Ctx) (out T, err error), schema ZogSchema) *PreprocessSchema[F, T] {
	return &PreprocessSchema[F, T]{fn: fn, schema: schema}
}

func (s *PreprocessSchema[F, T]) process(ctx *p.SchemaCtx) {
	v, ok := ctx.Data.(F)
	if !ok {
		ctx.AddIssue(ctx.IssueFromCoerce(fmt.Errorf("preprocess expected %T but got %T", v, ctx.Data)))
		return
	}
	out, err := s.fn(v, ctx)
	if err != nil {
		ctx.AddIssue(ctx.IssueFromUnknownError(err))
		return
	}
	ctx.Data = p.UnwrapPtr(out)
	s.schema.process(ctx)
}

func (s *PreprocessSchema[F, T]) validate(ctx *p.SchemaCtx) {
	v, ok := ctx.ValPtr.(F)
	if !ok {
		p.Panicf(p.PanicTypeCast, ctx.String(), new(F), ctx.ValPtr)
	}
	out, err := s.fn(v, ctx)
	if err != nil {
		ctx.AddIssue(ctx.Issue().SetMessage(err.Error()))
		return
	}
	switch v := ctx.ValPtr.(type) {
	case *T:
		*v = out
	case **T:
		*v = &out
	default:
		panic(fmt.Sprintf("Preprocessed should be passed in schema.Validate() a value pointer that is compatible with its returned type T. Either *T or **T. Got %T", v))
	}
	s.schema.validate(ctx)
}

func (s *PreprocessSchema[F, T]) getType() zconst.ZogType {
	return s.schema.getType()
}

func (s *PreprocessSchema[F, T]) setCoercer(coercer CoercerFunc) {
	s.schema.setCoercer(coercer)
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
