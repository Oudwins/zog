package zog

import (
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/pkgs/internals"
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
)

var _ ZogSchema = &AnySchema{}

type AnySchema struct {
	processors []p.ZProcessor[*any]
	defaultVal *any
	required   *p.Test[*any]
	catch      *any
}

// ! INTERNALS

// Returns the type of the schema
func (v *AnySchema) getType() zconst.ZogType {
	return zconst.TypeAny
}

// Sets the coercer for the schema (no-op for Any schema)
func (v *AnySchema) setCoercer(c CoercerFunc) {
	// no-op - Any schema doesn't need coercion
}

// ! USER FACING FUNCTIONS

// Returns a new EXPERIMENTAL_ANY Schema
// Do not use unless you know what you are doing & are okay with possible breaking changes & bugs.
func EXPERIMENTAL_ANY(opts ...SchemaOption) *AnySchema {
	a := &AnySchema{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Parse data into destination pointer
func (v *AnySchema) Parse(data any, dest *any, options ...ExecOption) ZogIssueList {
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
func (v *AnySchema) process(ctx *p.SchemaCtx) {
	ctx.CanCatch = v.catch != nil

	destPtr, ok := ctx.ValPtr.(*any)
	if !ok {
		p.Panicf(p.PanicTypeCast, ctx.String(), ctx.DType, ctx.ValPtr)
	}

	// Handle default/required for nil values
	isZeroVal := p.IsParseZeroValue(ctx.Data, ctx)
	if isZeroVal {
		if v.defaultVal != nil {
			*destPtr = *v.defaultVal
		} else if v.required == nil {
			// This handles optional case
			return
		} else {
			// is required & zero value
			if ctx.CanCatch {
				*destPtr = *v.catch
				return
			} else {
				ctx.AddIssue(ctx.IssueFromTest(v.required, *destPtr))
				return
			}
		}
	} else {
		// For Any schema, we just assign the value as-is
		*destPtr = ctx.Data
	}

	// Process all processors (tests and transforms)
	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(destPtr, ctx)
		if ctx.Exit {
			if ctx.CanCatch {
				*destPtr = *v.catch
				return
			}
			return
		}
	}
}

// Validate data against schema
func (v *AnySchema) Validate(val *any, options ...ExecOption) ZogIssueList {
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
func (v *AnySchema) validate(ctx *p.SchemaCtx) {
	ctx.CanCatch = v.catch != nil

	valPtr, ok := ctx.ValPtr.(*any)
	if !ok {
		p.Panicf(p.PanicTypeCast, ctx.String(), ctx.DType, ctx.ValPtr)
	}

	// Handle default/required for zero values
	isZeroVal := p.IsZeroValue(*valPtr)
	if isZeroVal {
		if v.defaultVal != nil {
			*valPtr = *v.defaultVal
		} else if v.required == nil {
			// This handles optional case
			return
		} else {
			// is required & zero value
			if ctx.CanCatch {
				*valPtr = *v.catch
				return
			} else {
				ctx.AddIssue(ctx.IssueFromTest(v.required, *valPtr))
				return
			}
		}
	}

	// Process all processors (tests and transforms)
	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(valPtr, ctx)
		if ctx.Exit {
			if ctx.CanCatch {
				*valPtr = *v.catch
				return
			}
			return
		}
	}
}

// GLOBAL METHODS

func (v *AnySchema) Test(t p.Test[*any]) *AnySchema {
	v.processors = append(v.processors, &t)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *AnySchema) TestFunc(testFunc p.BoolTFunc[*any], options ...p.TestOption) *AnySchema {
	test := p.NewTestFunc("", testFunc, options...)
	v.Test(*test)
	return v
}

// Adds a transform function to the schema. Runs in the order it is called
func (v *AnySchema) Transform(transform p.Transform[*any]) *AnySchema {
	v.processors = append(v.processors, &p.TransformProcessor[*any]{Transform: transform})
	return v
}

// ! MODIFIERS

// marks field as required
func (v *AnySchema) Required(options ...TestOption) *AnySchema {
	r := p.Required[*any]()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *AnySchema) Optional() *AnySchema {
	v.required = nil
	return v
}

// sets the default value
func (v *AnySchema) Default(val any) *AnySchema {
	v.defaultVal = &val
	return v
}

// sets the catch value (i.e the value to use if the validation fails)
func (v *AnySchema) Catch(val any) *AnySchema {
	v.catch = &val
	return v
}

// toZSS converts the schema to ZSS format
func (v *AnySchema) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(v.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeAny,
		Required:     toZSSRequired(v.required, zconst.TypeAny),
		DefaultValue: deepCopyPrimitivePtr(v.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(v.catch),
		Processors:   processorsToZSS(rvP, zconst.TypeAny),
	}
	return &j
}
