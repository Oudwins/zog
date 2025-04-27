package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ PrimitiveZogSchema[bool] = &BoolSchema[bool]{}

type BoolSchema[T ~bool] struct {
	processors []p.ZProcessor[*T]
	defaultVal *T
	required   *p.Test[*T]
	catch      *T
	coercer    CoercerFunc
}

// ! INTERNALS

// Returns the type of the schema
func (v *BoolSchema[T]) getType() zconst.ZogType {
	return zconst.TypeBool
}

// Sets the coercer for the schema
func (v *BoolSchema[T]) setCoercer(c CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Returns a new Bool Shape
func Bool(opts ...SchemaOption) *BoolSchema[bool] {
	b := &BoolSchema[bool]{
		coercer: conf.Coercers.Bool, // default coercer
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Parse data into destination pointer
func (v *BoolSchema[T]) Parse(data any, dest *T, options ...ExecOption) ZogIssueList {
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
func (v *BoolSchema[T]) process(ctx *p.SchemaCtx) {
	primitiveParsing(ctx, v.processors, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Validate data against schema
func (v *BoolSchema[T]) Validate(val *T, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}

	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(val, val, path, v.getType())
	defer sctx.Free()
	v.validate(sctx)
	return errs.List
}

// Internal function to validate data
func (v *BoolSchema[T]) validate(ctx *p.SchemaCtx) {
	primitiveValidation(ctx, v.processors, v.defaultVal, v.required, v.catch)
}

// GLOBAL METHODS

func (v *BoolSchema[T]) Test(t p.Test[*T]) *BoolSchema[T] {
	v.processors = append(v.processors, &t)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *BoolSchema[T]) TestFunc(testFunc p.BoolTFunc[*T], options ...p.TestOption) *BoolSchema[T] {
	test := p.NewTestFunc("", testFunc, options...)
	v.Test(*test)
	return v
}

// Adds a transform function to the schema. Runs in the order it is called (i.e schema.True().Transform(...) will run the transform after the True test)
func (v *BoolSchema[T]) Transform(transform p.Transform[*T]) *BoolSchema[T] {
	v.processors = append(v.processors, &p.TransformProcessor[*T]{Transform: transform})
	return v
}

// ! MODIFIERS
// marks field as required
func (v *BoolSchema[T]) Required(options ...TestOption) *BoolSchema[T] {
	r := p.Required[*T]()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *BoolSchema[T]) Optional() *BoolSchema[T] {
	v.required = nil
	return v
}

// sets the default value
func (v *BoolSchema[T]) Default(val T) *BoolSchema[T] {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *BoolSchema[T]) Catch(val T) *BoolSchema[T] {
	v.catch = &val
	return v
}

// UNIQUE METHODS

func (v *BoolSchema[T]) True() *BoolSchema[T] {
	t, fn := p.EQ[T](T(true))
	p.TestFuncFromBool(fn, &t)
	v.processors = append(v.processors, &t)
	return v
}

func (v *BoolSchema[T]) False() *BoolSchema[T] {
	t, fn := p.EQ[T](T(false))
	p.TestFuncFromBool(fn, &t)
	v.processors = append(v.processors, &t)
	return v
}

func (v *BoolSchema[T]) EQ(val T) *BoolSchema[T] {
	t, fn := p.EQ[T](val)
	p.TestFuncFromBool(fn, &t)
	v.processors = append(v.processors, &t)
	return v
}
