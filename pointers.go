package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type pointerProcessor struct {
	// preTransforms  []p.PreTransform
	tests    []p.Test
	schema   ZogSchema
	required *p.Test
	// postTransforms []p.PostTransform
	// defaultVal     *any
	// catch          *any
}

func (v *pointerProcessor) getType() zconst.ZogType {
	return zconst.TypePtr
}

func (v *pointerProcessor) setCoercer(c conf.CoercerFunc) {
	v.schema.setCoercer(c)
}

// Ptr creates a pointer ZogSchema
func Ptr(schema ZogSchema) *pointerProcessor {
	return &pointerProcessor{
		tests:  []p.Test{},
		schema: schema,
	}
}

// Parse the data into the destination pointer
func (v *pointerProcessor) Parse(data any, dest any, options ...ParsingOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.M
}

func (v *pointerProcessor) process(data any, dest any, path p.PathBuilder, ctx ParseCtx) {
	isZero := p.IsParseZeroValue(data, ctx)
	if isZero {
		if v.required != nil {
			// ctx.AddError(v.required)
			ctx.NewError(path, Errors.FromTest(data, v.schema.getType(), v.required, ctx))
		}
		return
	}
	rv := reflect.ValueOf(dest)
	destPtr := rv.Elem()
	if destPtr.IsNil() {
		// this sets the primitive also
		newVal := reflect.New(destPtr.Type().Elem())
		// this generates a new nil pointer
		//newVal := reflect.Zero(destPtr.Type())
		destPtr.Set(newVal)
	}
	di := destPtr.Interface()
	v.schema.process(data, di, path, ctx)
}

func (v *pointerProcessor) NotNil(options ...TestOption) *pointerProcessor {
	r := p.Test{
		ErrCode: zconst.ErrCodeNotNil,
	}
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}
