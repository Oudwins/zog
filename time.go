package zog

import (
	"time"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// ! INTERNALS
var _ PrimitiveZogSchema[time.Time] = &TimeSchema{}

type TimeSchema struct {
	processors []p.ZProcessor[*time.Time]
	defaultVal *time.Time
	required   *p.Test[*time.Time]
	catch      *time.Time
	coercer    conf.CoercerFunc
}

// Returns the type of the schema
func (v *TimeSchema) getType() zconst.ZogType {
	return zconst.TypeTime
}

// Sets the coercer for the schema
func (v *TimeSchema) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

type TimeFunc func(opts ...SchemaOption) *TimeSchema

// ! USER FACING FUNCTIONS

// Returns a new Time Shape
var Time TimeFunc = func(opts ...SchemaOption) *TimeSchema {
	t := &TimeSchema{
		coercer: conf.Coercers.Time,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// WARNING ONLY SUPPOORTS Shape.Parse!
// Sets the format function for the time schema.
// Usage is:
//
//	z.Time(z.Time.FormatFunc(func(data string) (time.Time, error) {
//		return time.Parse(time.RFC3339, data)
//	}))
func (t TimeFunc) FormatFunc(format func(data string) (time.Time, error)) SchemaOption {
	return func(s ZogSchema) {
		s.setCoercer(conf.TimeCoercerFactory(format))
	}
}

// WARNING ONLY SUPPOORTS Shape.Parse!
// Sets the string format for the  time schema
// Usage is:
// z.Time(z.Time.Format(time.RFC3339))
func (t TimeFunc) Format(format string) SchemaOption {
	return t.FormatFunc(func(data string) (time.Time, error) {
		return time.Parse(format, data)
	})
}

// Parses the data into the destination time.Time. Returns a list of errors
func (v *TimeSchema) Parse(data any, dest *time.Time, options ...ExecOption) ZogIssueList {
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

// internal processes the data
func (v *TimeSchema) process(ctx *p.SchemaCtx) {
	primitiveParsing(ctx, v.processors, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Validates an existing time.Time
func (v *TimeSchema) Validate(data *time.Time, options ...ExecOption) ZogIssueList {
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

// Internal function to validate the data
func (v *TimeSchema) validate(ctx *p.SchemaCtx) {
	primitiveValidation(ctx, v.processors, v.defaultVal, v.required, v.catch)
}

// Adds posttransform function to schema
func (v *TimeSchema) Transform(transform Transform[*time.Time]) *TimeSchema {
	v.processors = append(v.processors, &p.TransformProcessor[*time.Time]{Transform: p.Transform[*time.Time](transform)})
	return v
}

// ! MODIFIERS

// marks field as required
func (v *TimeSchema) Required(options ...TestOption) *TimeSchema {
	r := p.Required[*time.Time]()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *TimeSchema) Optional() *TimeSchema {
	v.required = nil
	return v
}

// sets the default value
func (v *TimeSchema) Default(val time.Time) *TimeSchema {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *TimeSchema) Catch(val time.Time) *TimeSchema {
	v.catch = &val
	return v
}

// GLOBAL METHODS

// custom test function call it -> schema.Test(z.Test{Func: func (val *time.Time, ctx z.Ctx) {
// my test
// }})
func (v *TimeSchema) Test(t Test[*time.Time]) *TimeSchema {
	x := p.Test[*time.Time](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *TimeSchema) TestFunc(testFunc BoolTFunc[*time.Time], options ...TestOption) *TimeSchema {
	test := p.NewTestFunc("", p.BoolTFunc[*time.Time](testFunc), options...)
	v.Test(Test[*time.Time](*test))
	return v
}

// UNIQUE METHODS

// Checks that the value is after the given time
func (v *TimeSchema) After(t time.Time, opts ...TestOption) *TimeSchema {
	fn := func(v *time.Time, ctx Ctx) bool {
		return (*v).After(t)
	}

	r := p.Test[*time.Time]{
		IssueCode: zconst.IssueCodeAfter,
		Params:    make(map[string]any, 1),
	}
	r.Params[zconst.IssueCodeAfter] = t
	p.TestFuncFromBool(fn, &r)
	for _, opt := range opts {
		opt(&r)
	}
	v.processors = append(v.processors, &r)
	return v
}

// Checks that the value is before the given time
func (v *TimeSchema) Before(t time.Time, opts ...TestOption) *TimeSchema {
	fn := func(v *time.Time, ctx Ctx) bool {
		return (*v).Before(t)
	}

	r := p.Test[*time.Time]{
		IssueCode: zconst.IssueCodeBefore,
		Params:    make(map[string]any, 1),
	}
	r.Params[zconst.IssueCodeBefore] = t
	p.TestFuncFromBool(fn, &r)
	for _, opt := range opts {
		opt(&r)
	}
	v.processors = append(v.processors, &r)
	return v
}

// Checks that the value is equal to the given time
func (v *TimeSchema) EQ(t time.Time, opts ...TestOption) *TimeSchema {
	fn := func(v *time.Time, ctx Ctx) bool {
		return (*v).Equal(t)
	}

	r := p.Test[*time.Time]{
		IssueCode: zconst.IssueCodeEQ,
		Params:    make(map[string]any, 1),
	}
	r.Params[zconst.IssueCodeEQ] = t
	p.TestFuncFromBool(fn, &r)
	for _, opt := range opts {
		opt(&r)
	}
	v.processors = append(v.processors, &r)

	return v
}
