package zog

import (
	p "github.com/Oudwins/zog/primitives"
)

type boolProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *bool
	required       *p.Test
	catch          *bool
}

func Bool() *boolProcessor {
	return &boolProcessor{
		tests: []p.Test{},
	}
}

func (v *boolProcessor) Parse(data any, dest *bool) p.ZogErrorList {
	var ctx = p.NewParseCtx()
	errs := p.NewErrsList()
	path := p.Pather("")

	v.process(data, dest, errs, path, ctx)

	return errs.List
}

func (v *boolProcessor) process(val any, dest *bool, errs p.ZogErrors, path p.Pather, ctx *p.ParseCtx) {
	primitiveProcess(val, dest, errs, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, p.Coercers["bool"])
}

// GLOBAL METHODS

// Adds pretransform function to schema
func (v *boolProcessor) PreTransform(transform p.PreTransform) *boolProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *boolProcessor) PostTransform(transform p.PostTransform) *boolProcessor {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *boolProcessor) Required() *boolProcessor {
	r := p.Required(p.DErrorFunc("is a required field"))
	v.required = &r
	return v
}

// marks field as optional
func (v *boolProcessor) Optional() *boolProcessor {
	v.required = nil
	return v
}

// sets the default value
func (v *boolProcessor) Default(val bool) *boolProcessor {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *boolProcessor) Catch(val bool) *boolProcessor {
	v.catch = &val
	return v
}

// UNIQUE METHODS

func (v *boolProcessor) True() *boolProcessor {
	v.tests = append(v.tests, p.EQ[bool](true, "should be true"))
	return v
}

func (v *boolProcessor) False() *boolProcessor {
	v.tests = append(v.tests, p.EQ[bool](false, "should be false"))
	return v
}
