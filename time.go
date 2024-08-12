package zog

import (
	"fmt"
	"time"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/primitives"
)

type timeProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	postTransforms []p.PostTransform
	defaultVal     *time.Time
	required       *p.Test
	catch          *time.Time
}

func Time() *timeProcessor {
	return &timeProcessor{}
}

func (v *timeProcessor) Parse(val any, dest *time.Time) p.ZogErrorList {
	var ctx = p.NewParseCtx()
	errs := p.NewErrsList()
	path := p.PathBuilder("")

	v.process(val, dest, errs, path, ctx)

	return errs.List
}

func (v *timeProcessor) process(val any, dest any, errs p.ZogErrors, path p.PathBuilder, ctx *p.ParseCtx) {
	primitiveProcessor(val, dest, errs, path, ctx, v.preTransforms, v.tests, v.postTransforms, v.defaultVal, v.required, v.catch, conf.Coercers["time"])
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
	r := p.Required(p.DErrorFunc("is a required field"))
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

// custom test function call it -> schema.Test("test_name", z.Message(""), func(val any, ctx *p.ParseCtx) bool {return true})
func (v *timeProcessor) Test(ruleName string, errorMsg TestOption, validateFunc p.TestFunc) *timeProcessor {
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

func (v *timeProcessor) After(t time.Time) *timeProcessor {
	r := p.Test{
		Name:      "timeAfter",
		ErrorFunc: p.DErrorFunc(fmt.Sprintf("is not after %v", t)),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(time.Time)
			if !ok {
				return false
			}
			return val.After(t)
		},
	}
	for _, opt := range v.tests {
		r.ErrorFunc = opt.ErrorFunc
	}
	v.tests = append(v.tests, r)
	return v
}

func (v *timeProcessor) Before(t time.Time) *timeProcessor {
	r :=
		p.Test{
			Name:      "timeBefore",
			ErrorFunc: p.DErrorFunc(fmt.Sprintf("is not before %v", t)),
			ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
				val, ok := v.(time.Time)
				if !ok {
					return false
				}
				return val.Before(t)
			},
		}
	for _, opt := range v.tests {
		r.ErrorFunc = opt.ErrorFunc
	}
	v.tests = append(v.tests, r)

	return v
}

func (v *timeProcessor) Is(t time.Time) *timeProcessor {
	r := p.Test{
		Name:      "timeIs",
		ErrorFunc: p.DErrorFunc(fmt.Sprintf("is not %v", t)),
		ValidateFunc: func(v any, ctx *p.ParseCtx) bool {
			val, ok := v.(time.Time)
			if !ok {
				return false
			}
			return val.Equal(t)
		},
	}

	for _, opt := range v.tests {
		r.ErrorFunc = opt.ErrorFunc
	}
	v.tests = append(v.tests, r)

	return v
}
