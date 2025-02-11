package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ComplexZogSchema = &PointerSchema{}

type PointerSchema struct {
	// preTransforms  []p.PreTransform
	tests    []p.Test
	schema   ZogSchema
	required *p.Test
	// postTransforms []p.PostTransform
	// defaultVal     *any
	// catch          *any
}

func (v *PointerSchema) getType() zconst.ZogType {
	return zconst.TypePtr
}

func (v *PointerSchema) setCoercer(c conf.CoercerFunc) {
	v.schema.setCoercer(c)
}

// Ptr creates a pointer ZogSchema
func Ptr(schema ZogSchema) *PointerSchema {
	return &PointerSchema{
		tests:  []p.Test{},
		schema: schema,
	}
}

// Parse the data into the destination pointer
func (v *PointerSchema) Parse(data any, dest any, options ...ExecOption) p.ZogIssueMap {
	errs := p.NewErrsMap()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(ctx.NewSchemaCtx(data, dest, path, v.getType()))

	return errs.M
}

func (v *PointerSchema) process(ctx *p.SchemaCtx) {
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
	v.schema.process(ctx.NewSchemaCtx(ctx.Val, di, ctx.Path, v.schema.getType()))
}

// Validates a pointer pointer
func (v *PointerSchema) Validate(data any, options ...ExecOption) p.ZogIssueMap {
	errs := p.NewErrsMap()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	v.validate(ctx.NewValidateSchemaCtx(data, p.PathBuilder(""), v.getType()))
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
	r := p.Test{
		IssueCode: zconst.IssueCodeNotNil,
	}
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}
