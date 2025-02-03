package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type Numeric interface {
	~int | ~float64
}

var _ PrimitiveZogSchema[int] = &NumberSchema[int]{}

type NumberSchema[T Numeric] struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *T
	required       *p.Test
	catch          *T
	coercer        conf.CoercerFunc
}

// ! INTERNALS

// Returns the type of the schema
func (v *NumberSchema[T]) getType() zconst.ZogType {
	return zconst.TypeNumber
}

// Sets the coercer for the schema
func (v *NumberSchema[T]) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// Internal function to process the data
func (v *NumberSchema[T]) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
	primitiveProcessor(val, dest, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// ! USER FACING FUNCTIONS

// creates a new float64 schema
func Float(opts ...SchemaOption) *NumberSchema[float64] {
	s := &NumberSchema[float64]{
		coercer: conf.Coercers.Float64,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// creates a new int schema
func Int(opts ...SchemaOption) *NumberSchema[int] {
	s := &NumberSchema[int]{
		coercer: conf.Coercers.Int,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// parses the value and stores it in the destination
func (v *NumberSchema[T]) Parse(data any, dest *T, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}

	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.List
}

// Validates a number pointer
func (v *NumberSchema[T]) Validate(data *T) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)

	v.validate(data, p.PathBuilder(""), ctx)
	return errs.List
}

func (v *NumberSchema[T]) validate(val any, path p.PathBuilder, ctx ParseCtx) {
	primitiveValidator(val, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch)
}

// GLOBAL METHODS

func (v *NumberSchema[T]) PreTransform(transform p.PreTransform) *NumberSchema[T] {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *NumberSchema[T]) PostTransform(transform p.PostTransform) *NumberSchema[T] {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *NumberSchema[T]) Required(options ...TestOption) *NumberSchema[T] {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *NumberSchema[T]) Optional() *NumberSchema[T] {
	v.required = nil
	return v
}

// sets the default value
func (v *NumberSchema[T]) Default(val T) *NumberSchema[T] {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *NumberSchema[T]) Catch(val T) *NumberSchema[T] {
	v.catch = &val
	return v
}

// custom test function call it -> schema.Test(test, options)
func (v *NumberSchema[T]) Test(t p.Test, opts ...TestOption) *NumberSchema[T] {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// UNIQUE METHODS

// Check that the value is one of the enum values
func (v *NumberSchema[T]) OneOf(enum []T, options ...TestOption) *NumberSchema[T] {
	t := p.In(enum)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for equality
func (v *NumberSchema[T]) EQ(n T, options ...TestOption) *NumberSchema[T] {
	t := p.EQ(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for lesser or equal
func (v *NumberSchema[T]) LTE(n T, options ...TestOption) *NumberSchema[T] {
	t := p.LTE(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for greater or equal
func (v *NumberSchema[T]) GTE(n T, options ...TestOption) *NumberSchema[T] {
	t := p.GTE(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for lesser
func (v *NumberSchema[T]) LT(n T, options ...TestOption) *NumberSchema[T] {
	t := p.LT(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// checks for greater
func (v *NumberSchema[T]) GT(n T, options ...TestOption) *NumberSchema[T] {
	t := p.GT(n)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}
