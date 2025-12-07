package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type CreateBoxFunc[T any, B any] func(data T, ctx Ctx) (B, error)
type UnboxFunc[B any, T any] func(data B, ctx Ctx) (T, error)

var _ ComplexZogSchema = &BoxedSchema[any, any]{}

type BoxedSchema[B any, T any] struct {
	schema ZogSchema
	unbox  UnboxFunc[B, T]
	box    CreateBoxFunc[T, B]
}

func Boxed[B any, T any](schema ZogSchema, unboxFunc UnboxFunc[B, T], boxFunc CreateBoxFunc[T, B]) *BoxedSchema[B, T] {
	return &BoxedSchema[B, T]{schema: schema, unbox: unboxFunc, box: boxFunc}
}

func (s *BoxedSchema[B, T]) Parse(data any, dest B, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, dest, path, s.getType())
	defer sctx.Free()
	s.process(sctx)
	return errs.M
}

func (s *BoxedSchema[B, T]) Validate(dest *B, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(*dest, dest, path, s.getType())
	defer sctx.Free()
	s.validate(sctx)
	return errs.M
}

func (s *BoxedSchema[B, T]) validate(ctx *p.SchemaCtx) {
	boxPtr, ok := ctx.ValPtr.(*B)
	if !ok {
		p.Panicf("BoxedSchema[%T, %T]: Expected valPtr type to correspond with type defined in schema. But it does not. Expected type: %T, got: %T", new(T), new(B), new(*B), ctx.ValPtr)
	}
	unboxed, err := s.unbox(*boxPtr, ctx)
	if err != nil {
		ctx.AddIssue(ctx.IssueFromUnknownError(err))
		return
	}
	ctx.Data = &unboxed
	ctx.ValPtr = &unboxed
	s.schema.validate(ctx)

	// Re-box and propagate back
	if s.box != nil {
		newBox, err := s.box(unboxed, ctx)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		*boxPtr = newBox
	}
}

func (s *BoxedSchema[B, T]) process(ctx *p.SchemaCtx) {
	s.schema.process(ctx)
}

func (s *BoxedSchema[B, T]) getType() zconst.ZogType {
	return s.schema.getType()
}

func (s *BoxedSchema[B, T]) setCoercer(c CoercerFunc) {
	s.schema.setCoercer(c)
}
