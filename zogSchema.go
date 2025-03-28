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
	setCoercer(c CoercerFunc)
	getType() zconst.ZogType
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

// Schema Parts Export

// Function signature for preTransforms. Takes the value and the context and returns the new value and an error.
type PreTransform = p.PreTransform

// Function signature for postTransforms. Takes the value pointer and the context and returns an error.
type PostTransform = p.PostTransform

// Function signature for issue formatters. Takes the issue and the context and returns the formatted issue.
type IssueFmtFunc = p.IssueFmtFunc

// Function signature for tests. Takes the value and the context and returns a boolean.
// This used to be a function you could pass to the schema.Test method -> `s.Test(z.TestFunc(fn))`. But that has been deprecated. Use `schema.TestFunc(fn)` instead.
type TestFunc = p.TestFunc

// Function signature for bool tests. Takes the value and the context and returns a boolean. This is the function passed to the TestFunc method.
type BoolTestFunc = p.BoolTestFunc

// ! PRIMITIVE PROCESSING -> Not userspace code

func primitiveProcessor[T p.ZogPrimitive](ctx *p.SchemaCtx, preTransforms []PreTransform, tests []Test, postTransforms []PostTransform, defaultVal *T, required *Test, catch *T, coercer CoercerFunc, isZeroFunc p.IsZeroValueFunc) {
	ctx.CanCatch = catch != nil

	destPtr := ctx.DestPtr.(*T)
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range postTransforms {
				err := fn(destPtr, ctx)
				if err != nil {
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	// 1. preTransforms
	for _, fn := range preTransforms {
		nVal, err := fn(ctx.Val, ctx)
		// bail if error in preTransform
		if err != nil || ctx.Exit {
			if ctx.CanCatch {
				*destPtr = *catch
				return
			}
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		ctx.Val = nVal
	}

	// 2. cast data to string & handle default/required
	isZeroVal := isZeroFunc(ctx.Val, ctx)

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
				ctx.AddIssue(ctx.IssueFromTest(required, ctx.Val))
				return
			}
		}
	} else {
		newVal, err := coercer(ctx.Val)
		if err == nil {
			*destPtr = newVal.(T)
		} else {
			if ctx.CanCatch {
				*destPtr = *catch
				return
			} else {
				ctx.AddIssue(ctx.IssueFromCoerce(err))
				return
			}
		}
	}

	// 3. tests
	for _, test := range tests {
		test.Func(destPtr, ctx)
		if ctx.Exit {
			if ctx.CanCatch {
				*destPtr = *catch
				return
			}
		}
	}

	// 4. postTransforms -> Done above on defer
}

func primitiveValidator[T p.ZogPrimitive](ctx *p.SchemaCtx, preTransforms []PreTransform, tests []Test, postTransforms []PostTransform, defaultVal *T, required *Test, catch *T) {
	ctx.CanCatch = catch != nil

	valPtr := ctx.Val.(*T)
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range postTransforms {
				err := fn(valPtr, ctx)
				if err != nil {
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	// 1. preTransforms
	for _, fn := range preTransforms {
		nVal, err := fn(*valPtr, ctx)
		// bail if error in preTransform
		if err != nil || ctx.Exit {
			if ctx.CanCatch {
				*valPtr = *catch
				return
			}
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		*valPtr = nVal.(T)
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
				ctx.AddIssue(ctx.IssueFromTest(required, ctx.Val))
				return
			}
		}
	}

	// 3. tests
	for _, test := range tests {
		test.Func(valPtr, ctx)
		if ctx.Exit {
			if ctx.CanCatch {
				*valPtr = *catch
				return
			}
		}
	}

}
