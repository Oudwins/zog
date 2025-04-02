package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ComplexZogSchema = &PointerSchema{}

type PointerSchema struct {
	// preTransforms  []PreTransform
	tests    []Test
	schema   ZogSchema
	required *Test
	// postTransforms []PostTransform
	// defaultVal     *any
	// catch          *any
}

func (v *PointerSchema) getType() zconst.ZogType {
	// return zconst.TypePtr
	return v.schema.getType()
}

func (v *PointerSchema) setCoercer(c conf.CoercerFunc) {
	v.schema.setCoercer(c)
}

// Ptr creates a pointer ZogSchema
func Ptr(schema ZogSchema) *PointerSchema {
	return &PointerSchema{
		tests:  []Test{},
		schema: schema,
	}
}

// Parse the data into the destination pointer
func (v *PointerSchema) Parse(data any, dest any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, dest, path, v.getType())
	defer sctx.Free()
	v.process(sctx)

	return errs.M
}

func (v *PointerSchema) process(ctx *p.SchemaCtx) {

	// TODO this is a mess. But couldn't figure out a simple way to support top level optional structs without doing this.
	// Companion code to this codde is in struct.go > process
	subCtx := ctx.NewSchemaCtx(ctx.Val, ctx.DestPtr, ctx.Path, v.schema.getType())
	defer subCtx.Free()
	if fn, ok := ctx.Val.(p.DpFactory); ok {
		val, err := fn()
		if err != nil {
			ctx.AddIssue(subCtx.IssueFromUnknownError(err))
			return
		}
		ctx.Val = val
	}
	// End of messy code

	isZero := p.IsParseZeroValue(ctx.Val, ctx)
	if isZero {
		if v.required != nil {
			// We set the destination type to the schema type because pointer doesn't have any issue messages. They pass through to the schema type
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val).SetDType(v.schema.getType()))
		}
		return
	}
	rv := reflect.ValueOf(ctx.DestPtr)
	destPtr := rv.Elem()
	if destPtr.IsNil() {
		// this sets the primitive also
		newVal := reflect.New(destPtr.Type().Elem())
		// this generates a new nil pointer
		//newVal := reflect.Zero(destPtr.Type())
		destPtr.Set(newVal)
	}
	di := destPtr.Interface()
	subCtx.DestPtr = di
	v.schema.process(subCtx)
}

// Validates a pointer pointer
func (v *PointerSchema) Validate(data any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	v.validate(ctx.NewValidateSchemaCtx(data, path, v.getType()))
	return errs.M
}

func (v *PointerSchema) validate(ctx *p.SchemaCtx) {
	rv := reflect.ValueOf(ctx.Val)
	destPtr := rv.Elem()
	if !destPtr.IsValid() || destPtr.IsNil() {
		if v.required != nil {
			// We set the destination type to the schema type because pointer doesn't have any issue messages. They pass through to the schema type
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val).SetDType(v.schema.getType()))
		}
		return
	}
	di := destPtr.Interface()
	ctx.Val = di
	v.schema.validate(ctx.NewValidateSchemaCtx(di, ctx.Path, v.schema.getType()))
}

// Validate Existing Pointer

func (v *PointerSchema) NotNil(options ...TestOption) *PointerSchema {
	r := Test{
		IssueCode: zconst.IssueCodeNotNil,
	}
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}
