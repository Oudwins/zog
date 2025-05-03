package zog

import (
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// The ZogSchema is the interface all schemas must implement
// This is most useful for internal use. If you are looking to pass schemas around, use the ComplexZogSchema or PrimitiveZogSchema interfaces if possible.
type ZogSchema interface {
	process(ctx *p.SchemaCtx)
	validate(ctx *p.SchemaCtx)
	getType() zconst.ZogType
	setCoercer(c CoercerFunc)
}

// This is a common interface for all complex schemas (i.e structs, slices, pointers...)
// You can use this to pass any complex schema around
type ComplexZogSchema interface {
	ZogSchema
	Parse(val any, dest any, options ...ExecOption) ZogIssueMap
}

// This is a common interface for all primitive schemas (i.e strings, numbers, booleans, time.Time...)
// You can use this to pass any primitive schema around
type PrimitiveZogSchema[T p.ZogPrimitive] interface {
	ZogSchema
	Parse(val any, dest *T, options ...ExecOption) ZogIssueList
}

// Shape Parts Export

// Function signature for transforms. Takes the value pointer and the context and returns an optional error.
type Transform[T any] p.Transform[T]

// Function signature for issue formatters. Takes the issue and the context and returns the formatted issue.
type IssueFmtFunc = p.IssueFmtFunc

// Function signature for tests. Takes the value and the context and returns a boolean.
// This used to be a function you could pass to the schema.Test method -> `s.Test(z.TFunc(fn))`. But that has been deprecated. Use `schema.TFunc(fn)` instead.
type TFunc[T any] p.TFunc[T]

// Function signature for bool tests. Takes the value and the context and returns a boolean. This is the function passed to the TestFunc method.
type BoolTFunc[T any] p.BoolTFunc[T]

// Creates a reusable testFunc you can add to schemas by doing schema.Test(z.TestFunc()). Has the same API as schema.TestFunc() so it is recommended you use that one for non reusable tests.
func TestFunc[T any](IssueCode zconst.ZogIssueCode, fn BoolTFunc[T], options ...p.TestOption) Test[T] {
	return Test[T](*p.NewTestFunc(IssueCode, p.BoolTFunc[T](fn), options...))
}

// ! PRIMITIVE PROCESSING -> Not userspace code

func primitiveParsing[T p.ZogPrimitive](ctx *p.SchemaCtx, processors []p.ZProcessor[*T], defaultVal *T, required *p.Test[*T], catch *T, coercer CoercerFunc, isZeroFunc p.IsZeroValueFunc) {
	ctx.CanCatch = catch != nil

	destPtr, ok := ctx.ValPtr.(*T)
	if !ok {
		p.Panicf(p.PanicTypeCast, ctx.String(), ctx.DType, ctx.ValPtr)
	}

	// 2. cast data to string & handle default/required
	isZeroVal := isZeroFunc(ctx.Data, ctx)
	if isZeroVal {
		if defaultVal != nil {
			*destPtr = *defaultVal
		} else if required == nil {
			// This handles optional case
			return
		} else {
			// is required & zero value
			// required
			if ctx.CanCatch {
				*destPtr = *catch
				return
			} else {
				ctx.AddIssue(ctx.IssueFromTest(required, *destPtr))
				return
			}
		}
	} else {
		v, err := coercer(ctx.Data)
		if err != nil {
			if ctx.CanCatch {
				*destPtr = *catch
				return
			}
			ctx.AddIssue(ctx.IssueFromCoerce(err))
			return
		}
		x, ok := v.(T)
		if !ok {
			p.Panicf(p.PanicTypeCastCoercer, ctx.String(), ctx.DType, v)
		}
		*destPtr = x
	}

	for _, processor := range processors {
		ctx.Processor = processor
		processor.ZProcess(destPtr, ctx)
		if ctx.Exit {
			if ctx.CanCatch {
				*destPtr = *catch
				return
			}
			return
		}
	}
}

func primitiveValidation[T p.ZogPrimitive](ctx *p.SchemaCtx, processors []p.ZProcessor[*T], defaultVal *T, required *p.Test[*T], catch *T) {
	ctx.CanCatch = catch != nil

	valPtr, ok := ctx.ValPtr.(*T)
	if !ok {
		p.Panicf(p.PanicTypeCast, ctx.String(), ctx.DType, ctx.ValPtr)
	}

	// 2. cast data to string & handle default/required
	// Warning. This uses generic IsZeroValue because for Validate we treat zero values as invalid for required fields. This is different from Parse.
	isZeroVal := p.IsZeroValue(*valPtr)

	if isZeroVal {
		if defaultVal != nil {
			*valPtr = *defaultVal
		} else if required == nil {
			// This handles optional case
			return
		} else {
			// is required & zero value
			// required
			if ctx.CanCatch {
				*valPtr = *catch
				return
			} else {
				ctx.AddIssue(ctx.IssueFromTest(required, *valPtr))
				return
			}
		}
	}

	for _, processor := range processors {
		ctx.Processor = processor
		processor.ZProcess(valPtr, ctx)
		if ctx.Exit {
			if ctx.CanCatch {
				*valPtr = *catch
				return
			}
			return
		}
	}
}
