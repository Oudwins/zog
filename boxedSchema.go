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

func (s *BoxedSchema[B, T]) Parse(data any, dest any, options ...ExecOption) ZogIssueMap {
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
		p.Panicf("BoxedSchema[%T, %T]: Expected valPtr type to correspond with type defined in schema. But it does not. Expected type: %T, got: %T", new(B), new(T), new(*B), ctx.ValPtr)
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
	boxPtr, ok := ctx.ValPtr.(*B)
	if !ok {
		p.Panicf("BoxedSchema[%T, %T]: Expected valPtr type to correspond with type defined in schema. But it does not. Expected type: %T, got: %T", new(B), new(T), new(*B), ctx.ValPtr)
	}

	// 1. Handle ctx.Data - could be B, *B, or raw data
	var innerData any
	switch d := ctx.Data.(type) {
	case B:
		// Unbox B to get T
		unboxed, err := s.unbox(d, ctx)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		innerData = unboxed
	case *B:
		// Dereference and unbox
		unboxed, err := s.unbox(*d, ctx)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		innerData = unboxed
	default:
		// Raw data - pass directly to inner schema
		innerData = d
	}

	// 2. Create temporary T for inner schema and pass data
	var inner T
	ctx.Data = innerData
	ctx.ValPtr = &inner

	// 3. Process through inner schema (keeps pointer to inner)
	s.schema.process(ctx)

	// 4. Re-box and set to original destination
	if s.box != nil {
		newBox, err := s.box(inner, ctx)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		*boxPtr = newBox
	}
	// TODO maybe some kind of flag that you executed process with boxFunc is nil
}

func (s *BoxedSchema[B, T]) getType() zconst.ZogType {
	return s.schema.getType()
}

func (s *BoxedSchema[B, T]) setCoercer(c CoercerFunc) {
	s.schema.setCoercer(c)
}
