package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ComplexZogSchema = &PointerSchema{}

type PointerSchema struct {
	schema   ZogSchema
	required *p.Test[any]
	nullable bool
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
		schema: schema,
	}
}

// Parse the data into the destination pointer
func (v *PointerSchema) Parse(data any, dest any, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
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

	return errs.List
}

func (v *PointerSchema) process(ctx *p.SchemaCtx) {

	// TODO this is a mess. But couldn't figure out a simple way to support top level optional structs without doing this.
	// Companion code to this codde is in struct.go > process
	subCtx := ctx.NewSchemaCtx(ctx.Data, ctx.ValPtr, ctx.Path, v.schema.getType())
	defer subCtx.Free()
	if fn, ok := ctx.Data.(p.DpFactory); ok {
		val, err := fn()
		if err != nil {
			ctx.AddIssue(subCtx.IssueFromUnknownError(err))
			return
		}
		ctx.Data = val
	}

	// Nullable clears the destination on explicit null. NotNil cannot be set here
	// because Nullable and NotNil clear each other.
	if v.nullable && p.IsExplicitNull(ctx.Data) {
		rv := reflect.ValueOf(ctx.ValPtr)
		destPtr := rv.Elem()
		destPtr.Set(reflect.Zero(destPtr.Type()))
		return
	}

	_, isEmptyStruct := ctx.Data.(*p.EmptyDataProvider)
	// End of messy code

	isZero := p.IsParseZeroValue(ctx.Data, ctx) || isEmptyStruct
	if isZero {
		if v.required != nil {
			// We set the destination type to the schema type because pointer doesn't have any issue messages. They pass through to the schema type
			issueVal := ctx.Data
			if p.IsExplicitNull(issueVal) {
				// Observable user intent is nil; keep the internal sentinel out of
				// issues.
				issueVal = nil
			}
			ctx.AddIssue(ctx.IssueFromTest(v.required, issueVal).SetDType(v.schema.getType()))
		}
		return
	}
	rv := reflect.ValueOf(ctx.ValPtr)
	destPtr := rv.Elem()
	if destPtr.IsNil() {
		// this sets the primitive also
		newVal := reflect.New(destPtr.Type().Elem())
		// this generates a new nil pointer
		//newVal := reflect.Zero(destPtr.Type())
		destPtr.Set(newVal)
	}
	di := destPtr.Interface()
	subCtx.ValPtr = di
	v.schema.process(subCtx)
}

// Validates a pointer pointer
func (v *PointerSchema) Validate(data any, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	v.validate(ctx.NewValidateSchemaCtx(data, path, v.getType()))
	return errs.List
}

func (v *PointerSchema) validate(ctx *p.SchemaCtx) {
	rv := reflect.ValueOf(ctx.ValPtr)
	destPtr := rv.Elem()
	if !destPtr.IsValid() || destPtr.IsNil() {
		if v.required != nil {
			// We set the destination type to the schema type because pointer doesn't have any issue messages. They pass through to the schema type
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Data).SetDType(v.schema.getType()))
		}
		return
	}
	di := destPtr.Interface()
	ctx.ValPtr = di
	v.schema.validate(ctx.NewValidateSchemaCtx(di, ctx.Path, v.schema.getType()))
}

// Validate Existing Pointer

func (v *PointerSchema) NotNil(options ...TestOption) *PointerSchema {
	r := p.Test[any]{
		IssueCode: zconst.IssueCodeNotNil,
	}
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	v.nullable = false
	return v
}

// Nullable makes the pointer schema treat explicit null input (e.g. JSON
// null) as a request to clear the destination pointer. An absent key still
// preserves the destination; only an explicit null clears it. Last-writer
// wins against NotNil: calling Nullable clears a prior NotNil requirement,
// and calling NotNil afterwards clears Nullable.
func (v *PointerSchema) Nullable() *PointerSchema {
	v.nullable = true
	v.required = nil
	return v
}
