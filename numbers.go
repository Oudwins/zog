package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
	"golang.org/x/exp/constraints"
)

type Numeric = constraints.Ordered

var _ PrimitiveZogSchema[int] = &NumberSchema[int]{}

type NumberSchema[T Numeric] struct {
	preTransforms  []PreTransform
	tests          []Test
	postTransforms []PostTransform
	defaultVal     *T
	required       *Test
	catch          *T
	coercer        conf.CoercerFunc
}

// ! INTERNALS

// Returns the type of the schema
func (v *NumberSchema[T]) getType() zconst.ZogType {
	return zconst.TypeNumber
}

// Sets the coercer for the schema
func (v *NumberSchema[T]) setCoercer(c CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Deprecated: Use Float64 instead
// creates a new float64 schema
func Float(opts ...SchemaOption) *NumberSchema[float64] {
	return Float64(opts...)
}

func Float64(opts ...SchemaOption) *NumberSchema[float64] {
	s := &NumberSchema[float64]{
		coercer: conf.Coercers.Float64,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func Float32(opts ...SchemaOption) *NumberSchema[float32] {
	s := &NumberSchema[float32]{
		coercer: func(data any) (any, error) {
			x, err := conf.Coercers.Float64(data)
			if err != nil {
				return nil, err
			}
			if n, ok := x.(float64); ok {
				return float32(n), nil
			}
			return x, nil
		},
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

func Int64(opts ...SchemaOption) *NumberSchema[int64] {
	s := &NumberSchema[int64]{
		coercer: func(data any) (any, error) {
			x, err := conf.Coercers.Int(data)
			if err != nil {
				return nil, err
			}
			if n, ok := x.(int); ok {
				return int64(n), nil
			}
			return x, nil
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func Int32(opts ...SchemaOption) *NumberSchema[int32] {
	s := &NumberSchema[int32]{
		coercer: func(data any) (any, error) {
			x, err := conf.Coercers.Int(data)
			if err != nil {
				return nil, err
			}
			if n, ok := x.(int); ok {
				return int32(n), nil
			}
			return x, nil
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// parses the value and stores it in the destination
func (v *NumberSchema[T]) Parse(data any, dest *T, options ...ExecOption) ZogIssueList {
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

// Internal function to process the data
func (v *NumberSchema[T]) process(ctx *p.SchemaCtx) {
	primitiveProcessor(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Validates a number pointer
func (v *NumberSchema[T]) Validate(data *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}

	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, data, path, v.getType())
	defer sctx.Free()
	v.validate(sctx)
	return errs.List
}

func (v *NumberSchema[T]) validate(ctx *p.SchemaCtx) {
	primitiveValidator(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch)
}

// GLOBAL METHODS

func (v *NumberSchema[T]) PreTransform(transform PreTransform) *NumberSchema[T] {
	if v.preTransforms == nil {
		v.preTransforms = []PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *NumberSchema[T]) PostTransform(transform PostTransform) *NumberSchema[T] {
	if v.postTransforms == nil {
		v.postTransforms = []PostTransform{}
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
func (v *NumberSchema[T]) Test(t Test, opts ...TestOption) *NumberSchema[T] {
	for _, opt := range opts {
		opt(&t)
	}
	t.ValidateFunc = customTestBackwardsCompatWrapper(t.ValidateFunc)
	v.tests = append(v.tests, t)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *NumberSchema[T]) TestFunc(testFunc p.TestFunc, options ...TestOption) *NumberSchema[T] {
	test := TestFunc("", testFunc)
	v.Test(test, options...)
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
