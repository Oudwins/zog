package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ PrimitiveZogSchema[bool] = &BoolSchema{}

type BoolSchema struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *bool
	required       *p.Test
	catch          *bool
	coercer        conf.CoercerFunc
}

// ! INTERNALS

// Returns the type of the schema
func (v *BoolSchema) getType() zconst.ZogType {
	return zconst.TypeBool
}

// Sets the coercer for the schema
func (v *BoolSchema) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// Internal function to process the data
func (v *BoolSchema) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
	primitiveProcessor(val, dest, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// ! USER FACING FUNCTIONS

// Returns a new Bool Schema
func Bool(opts ...SchemaOption) *BoolSchema {
	b := &BoolSchema{
		coercer: conf.Coercers.Bool, // default coercer
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Parse data into destination pointer
func (v *BoolSchema) Parse(data any, dest *bool, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.List
}

// Validate data against schema
func (v *BoolSchema) Validate(val *bool, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.validate(val, path, ctx)
	return errs.List
}

// Internal function to validate data
func (v *BoolSchema) validate(val any, path p.PathBuilder, ctx ParseCtx) {
	primitiveValidator(val, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch)
}

// GLOBAL METHODS

// Adds pretransform function to schema
func (v *BoolSchema) PreTransform(transform p.PreTransform) *BoolSchema {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *BoolSchema) PostTransform(transform p.PostTransform) *BoolSchema {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS
// marks field as required
func (v *BoolSchema) Required(options ...TestOption) *BoolSchema {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *BoolSchema) Optional() *BoolSchema {
	v.required = nil
	return v
}

// sets the default value
func (v *BoolSchema) Default(val bool) *BoolSchema {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *BoolSchema) Catch(val bool) *BoolSchema {
	v.catch = &val
	return v
}

// UNIQUE METHODS

func (v *BoolSchema) True() *BoolSchema {
	v.tests = append(v.tests, p.EQ[bool](true))
	return v
}

func (v *BoolSchema) False() *BoolSchema {
	v.tests = append(v.tests, p.EQ[bool](false))
	return v
}
