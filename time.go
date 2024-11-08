package zog

import (
	"time"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// ! INTERNALS
var _ ZogSchema = &timeProcessor{}

type timeProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *time.Time
	required       *p.Test
	catch          *time.Time
	coercer        conf.CoercerFunc
}

// internal processes the data
func (v *timeProcessor) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
	primitiveProcessor(val, dest, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, v.coercer, p.IsParseZeroValue)
}

// Returns the type of the schema
func (v *timeProcessor) getType() zconst.ZogType {
	return zconst.TypeTime
}

// Sets the coercer for the schema
func (v *timeProcessor) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

type TimeFunc func(opts ...SchemaOption) *timeProcessor

// ! USER FACING FUNCTIONS

// Returns a new Time Schema
var Time TimeFunc = func(opts ...SchemaOption) *timeProcessor {
	t := &timeProcessor{
		coercer: conf.Coercers.Time,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Sets the format function for the time schema
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

// Sets the string format for the  time schema
// Usage is:
// z.Time(z.Time.Format(time.RFC3339))
func (t TimeFunc) Format(format string) SchemaOption {
	return t.FormatFunc(func(data string) (time.Time, error) {
		return time.Parse(format, data)
	})
}

// Parses the data into the destination time.Time. Returns a list of errors
func (v *timeProcessor) Parse(data any, dest *time.Time, options ...ParsingOption) p.ZogErrList {
	errs := p.NewErrsList()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)

	for _, opt := range options {
		opt(ctx)
	}

	path := p.PathBuilder("")

	v.process(data, dest, path, ctx)

	return errs.List
}

// Adds pretransform function to schema
func (v *timeProcessor) PreTransform(transform p.PreTransform) *timeProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *timeProcessor) PostTransform(transform p.PostTransform) *timeProcessor {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *timeProcessor) Required(options ...TestOption) *timeProcessor {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *timeProcessor) Optional() *timeProcessor {
	v.required = nil
	return v
}

// sets the default value
func (v *timeProcessor) Default(val time.Time) *timeProcessor {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *timeProcessor) Catch(val time.Time) *timeProcessor {
	v.catch = &val
	return v
}

// GLOBAL METHODS

// custom test function call it -> schema.Test("error_code", func(val any, ctx ParseCtx) bool {return true})
func (v *timeProcessor) Test(t p.Test, opts ...TestOption) *timeProcessor {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// UNIQUE METHODS

// Checks that the value is after the given time
func (v *timeProcessor) After(t time.Time, opts ...TestOption) *timeProcessor {
	r := p.Test{
		ErrCode: zconst.ErrCodeAfter,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(time.Time)
			if !ok {
				return false
			}
			return val.After(t)
		},
	}
	r.Params[zconst.ErrCodeAfter] = t
	for _, opt := range opts {
		opt(&r)
	}
	for _, opt := range v.tests {
		r.ErrFmt = opt.ErrFmt
	}
	v.tests = append(v.tests, r)
	return v
}

// Checks that the value is before the given time
func (v *timeProcessor) Before(t time.Time, opts ...TestOption) *timeProcessor {
	r :=
		p.Test{
			ErrCode: zconst.ErrCodeBefore,
			Params:  make(map[string]any, 1),
			ValidateFunc: func(v any, ctx ParseCtx) bool {
				val, ok := v.(time.Time)
				if !ok {
					return false
				}
				return val.Before(t)
			},
		}
	r.Params[zconst.ErrCodeBefore] = t
	for _, opt := range opts {
		opt(&r)
	}
	v.tests = append(v.tests, r)

	return v
}

// Checks that the value is equal to the given time
func (v *timeProcessor) EQ(t time.Time, opts ...TestOption) *timeProcessor {
	r := p.Test{
		ErrCode: zconst.ErrCodeEQ,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(v any, ctx ParseCtx) bool {
			val, ok := v.(time.Time)
			if !ok {
				return false
			}
			return val.Equal(t)
		},
	}
	r.Params[zconst.ErrCodeEQ] = t
	for _, opt := range opts {
		opt(&r)
	}
	v.tests = append(v.tests, r)

	return v
}
