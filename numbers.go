package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type Numeric = p.Numeric

var _ PrimitiveZogSchema[int] = &NumberSchema[int]{}

type NumberSchema[T Numeric] struct {
	processors []p.ZProcessor[*T]
	defaultVal *T
	required   *p.Test[*T]
	catch      *T
	coercer    CoercerFunc
	isNot      bool
}

type NotNumberSchema[T Numeric] interface {
	OneOf(enum []T, options ...TestOption) *NumberSchema[T]
	EQ(n T, options ...TestOption) *NumberSchema[T]
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
// creates a new float64 schema. No plans to remove it but recommended to use Float64 instead.
func Float(opts ...SchemaOption) *NumberSchema[float64] {
	return Float64(opts...)
}

func FloatLike[T Numeric](opts ...SchemaOption) *NumberSchema[T] {
	s := &NumberSchema[T]{
		coercer: func(data any) (any, error) {
			x, err := conf.Coercers.Float64(data)
			if err != nil {
				return nil, err
			}
			return T(x.(float64)), nil
		},
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
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

func IntLike[T Numeric](opts ...SchemaOption) *NumberSchema[T] {
	s := &NumberSchema[T]{
		coercer: func(data any) (any, error) {
			x, err := conf.Coercers.Int(data)
			if err != nil {
				return nil, err
			}
			return T(x.(int)), nil
		},
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

// creates a uint schema
func Uint(opts ...SchemaOption) *NumberSchema[uint] {
	s := &NumberSchema[uint]{
		coercer: conf.Coercers.Uint,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func UintLike[T Numeric](opts ...SchemaOption) *NumberSchema[T] {
	s := &NumberSchema[T]{
		coercer: func(data any) (any, error) {
			x, err := conf.Coercers.Uint(data)
			if err != nil {
				return nil, err
			}
			return T(x.(uint)), nil
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
	primitiveParsing(ctx, v.processors, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
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
	primitiveValidation(ctx, v.processors, v.defaultVal, v.required, v.catch)
}

// GLOBAL METHODS

// Adds a transform function to the schema. Runs in the order it is called
func (v *NumberSchema[T]) Transform(transform p.Transform[*T]) *NumberSchema[T] {
	v.processors = append(v.processors, &p.TransformProcessor[*T]{Transform: transform})
	return v
}

// ! MODIFIERS

// marks field as required
func (v *NumberSchema[T]) Required(options ...TestOption) *NumberSchema[T] {
	r := p.Required[*T]()
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
func (v *NumberSchema[T]) Test(t Test[*T]) *NumberSchema[T] {
	x := p.Test[*T](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *NumberSchema[T]) TestFunc(testFunc BoolTFunc[*T], options ...TestOption) *NumberSchema[T] {
	test := p.NewTestFunc("", p.BoolTFunc[*T](testFunc), options...)
	v.Test(Test[*T](*test))
	return v
}

// UNIQUE METHODS

// Check that the value is one of the enum values
func (v *NumberSchema[T]) OneOf(enum []T, options ...TestOption) *NumberSchema[T] {
	t, fn := p.In(enum)

	return v.addTest(&t, fn, options...)
}

// checks for equality
func (v *NumberSchema[T]) EQ(n T, options ...TestOption) *NumberSchema[T] {
	t, fn := p.EQ(n)

	return v.addTest(&t, fn, options...)
}

// checks for lesser or equal
func (v *NumberSchema[T]) LTE(n T, options ...TestOption) *NumberSchema[T] {
	t, fn := p.LTE(n)

	return v.addTest(&t, fn, options...)
}

// checks for greater or equal
func (v *NumberSchema[T]) GTE(n T, options ...TestOption) *NumberSchema[T] {
	t, fn := p.GTE(n)

	return v.addTest(&t, fn, options...)
}

// checks for lesser
func (v *NumberSchema[T]) LT(n T, options ...TestOption) *NumberSchema[T] {
	t, fn := p.LT(n)

	return v.addTest(&t, fn, options...)
}

// checks for greater
func (v *NumberSchema[T]) GT(n T, options ...TestOption) *NumberSchema[T] {
	t, fn := p.GT(n)

	return v.addTest(&t, fn, options...)
}

func (v *NumberSchema[T]) Not() NotNumberSchema[T] {
	v.isNot = true
	return v
}

func (v *NumberSchema[T]) addTest(t *p.Test[*T], fn p.BoolTFunc[*T], options ...TestOption) *NumberSchema[T] {
	if v.isNot {
		p.TestNotFuncFromBool(fn, t)
		t.IssueCode = zconst.NotIssueCode(t.IssueCode)
		v.isNot = false
	} else {
		p.TestFuncFromBool(fn, t)
	}

	for _, opt := range options {
		opt(t)
	}

	v.processors = append(v.processors, t)
	return v
}
