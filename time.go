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
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *time.Time
	required       *p.Test
	catch          *time.Time
	coercer        conf.CoercerFunc
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

// Returns a new Time Schema
var Time TimeFunc = func(opts ...SchemaOption) *TimeSchema {
	t := &TimeSchema{
		coercer: conf.Coercers.Time,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// WARNING ONLY SUPPOORTS Schema.Parse!
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

// WARNING ONLY SUPPOORTS Schema.Parse!
// Sets the string format for the  time schema
// Usage is:
// z.Time(z.Time.Format(time.RFC3339))
func (t TimeFunc) Format(format string) SchemaOption {
	return t.FormatFunc(func(data string) (time.Time, error) {
		return time.Parse(format, data)
	})
}

// Parses the data into the destination time.Time. Returns a list of errors
func (v *TimeSchema) Parse(data any, dest *time.Time, options ...ExecOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewExecCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(ctx.NewSchemaCtx(data, dest, path, v.getType()))

	return errs.List
}

// internal processes the data
func (v *TimeSchema) process(ctx *p.SchemaCtx) {
	primitiveProcessor(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Validates an existing time.Time
func (v *TimeSchema) Validate(data *time.Time, options ...ExecOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewExecCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	v.validate(ctx.NewValidateSchemaCtx(data, p.PathBuilder(""), v.getType()))
	return errs.List
}

// Internal function to validate the data
func (v *TimeSchema) validate(ctx *p.SchemaCtx) {
	primitiveValidator(ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch)
}

// Adds pretransform function to schema
func (v *TimeSchema) PreTransform(transform p.PreTransform) *TimeSchema {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *TimeSchema) PostTransform(transform p.PostTransform) *TimeSchema {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *TimeSchema) Required(options ...TestOption) *TimeSchema {
	r := p.Required()
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

// custom test function call it -> schema.Test("error_code", func(val any, ctx ParseCtx) bool {return true})
func (v *TimeSchema) Test(t p.Test, opts ...TestOption) *TimeSchema {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *TimeSchema) TestFunc(testFunc p.TestFunc, options ...TestOption) *TimeSchema {
	test := TestFunc("", testFunc)
	v.Test(test, options...)
	return v
}

// UNIQUE METHODS

// Checks that the value is after the given time
func (v *TimeSchema) After(t time.Time, opts ...TestOption) *TimeSchema {
	r := p.Test{
		IssueCode: zconst.IssueCodeAfter,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(time.Time)
			if !ok {
				return false
			}
			return val.After(t)
		},
	}
	r.Params[zconst.IssueCodeAfter] = t
	for _, opt := range opts {
		opt(&r)
	}
	v.tests = append(v.tests, r)
	return v
}

// Checks that the value is before the given time
func (v *TimeSchema) Before(t time.Time, opts ...TestOption) *TimeSchema {
	r :=
		p.Test{
			IssueCode: zconst.IssueCodeBefore,
			Params:    make(map[string]any, 1),
			ValidateFunc: func(v any, ctx ParseCtx) bool {
				val, ok := v.(time.Time)
				if !ok {
					return false
				}
				return val.Before(t)
			},
		}
	r.Params[zconst.IssueCodeBefore] = t
	for _, opt := range opts {
		opt(&r)
	}
	v.tests = append(v.tests, r)

	return v
}

// Checks that the value is equal to the given time
func (v *TimeSchema) EQ(t time.Time, opts ...TestOption) *TimeSchema {
	r := p.Test{
		IssueCode: zconst.IssueCodeEQ,
		Params:    make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(time.Time)
			if !ok {
				return false
			}
			return val.Equal(t)
		},
	}
	r.Params[zconst.IssueCodeEQ] = t
	for _, opt := range opts {
		opt(&r)
	}
	v.tests = append(v.tests, r)

	return v
}
