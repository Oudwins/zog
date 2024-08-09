package zog

import (
	"fmt"

	p "github.com/Oudwins/zog/primitives"
)

type Numeric interface {
	~int | ~float64
}

type numberProcessor[T Numeric] struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *T
	required       *p.Test
	catch          *T
}

// creates a new float64 processor
func Float() *numberProcessor[float64] {
	return &numberProcessor[float64]{}
}

// creates a new int processor
func Int() *numberProcessor[int] {
	return &numberProcessor[int]{}
}

// parses the value and stores it in the destination
func (v *numberProcessor[T]) Parse(val any, dest *T) p.ZogErrorList {
	// TODO create context -> but for single field
	var ctx = p.NewParseCtx()
	errs := p.NewErrsList()
	path := p.PathBuilder("")
	// TODO handle options

	v.process(val, dest, errs, path, ctx)

	return errs.List
}

func (v *numberProcessor[T]) process(val any, dest any, errs p.ZogErrors, path p.PathBuilder, ctx *p.ParseCtx) {

	var coercer p.CoercerFunc
	switch any(dest).(type) {
	case *float64:
		coercer = p.Coercers["float64"]
	case *int:
		coercer = p.Coercers["int"]
	}

	primitiveProcessor(val, dest, errs, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, coercer)
}

// GLOBAL METHODS

func (v *numberProcessor[T]) PreTransform(transform p.PreTransform) *numberProcessor[T] {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *numberProcessor[T]) PostTransform(transform p.PostTransform) *numberProcessor[T] {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *numberProcessor[T]) Required(options ...TestOption) *numberProcessor[T] {
	r := p.Required(p.DErrorFunc("is a required field"))
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *numberProcessor[T]) Optional() *numberProcessor[T] {
	v.required = nil
	return v
}

// sets the default value
func (v *numberProcessor[T]) Default(val T) *numberProcessor[T] {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *numberProcessor[T]) Catch(val T) *numberProcessor[T] {
	v.catch = &val
	return v
}

func (v *numberProcessor[T]) Test(ruleName string, errorMsg TestOption, validateFunc p.TestFunc) *numberProcessor[T] {
	t := p.Test{
		Name:         ruleName,
		ErrorFunc:    nil,
		ValidateFunc: validateFunc,
	}
	errorMsg(&t)
	v.tests = append(v.tests, t)

	return v
}

// UNIQUE METHODS

func (v *numberProcessor[T]) OneOf(enum []T, options ...TestOption) *numberProcessor[T] {
	t := p.In(enum, fmt.Sprintf("should be one of %v", enum))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for equality
func (v *numberProcessor[T]) EQ(n T, options ...TestOption) *numberProcessor[T] {
	t := p.EQ(n, fmt.Sprintf("should be equal to %v", n))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for lesser or equal
func (v *numberProcessor[T]) LTE(n T, options ...TestOption) *numberProcessor[T] {
	t := p.LTE(n, fmt.Sprintf("should be lesser or equal than %v", n))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for greater or equal
func (v *numberProcessor[T]) GTE(n T, options ...TestOption) *numberProcessor[T] {
	t := p.GTE(n, fmt.Sprintf("should be greater or equal to %v", n))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for lesser
func (v *numberProcessor[T]) LT(n T, options ...TestOption) *numberProcessor[T] {
	t := p.LT(n, fmt.Sprintf("should be less than %v", n))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for greater
func (v *numberProcessor[T]) GT(n T, options ...TestOption) *numberProcessor[T] {
	t := p.GT(n, fmt.Sprintf("should be greater than %v", n))
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}
