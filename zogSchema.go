package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// The ZogSchema is the interface all schemas must implement
// This is most useful for internal use. If you are looking to pass schemas around, use the ComplexZogSchema or PrimitiveZogSchema interfaces if possible.
type ZogSchema interface {
	process(ctx *p.SchemaCtx)
	validate(ctx *p.SchemaCtx)
	setCoercer(c conf.CoercerFunc)
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

// ! PRIMITIVE PROCESSING

func primitiveProcessor[T p.ZogPrimitive](ctx *p.SchemaCtx, preTransforms []p.PreTransform, tests []p.Test, postTransforms []p.PostTransform, defaultVal *T, required *p.Test, catch *T, coercer conf.CoercerFunc, isZeroFunc p.IsZeroValueFunc) {
	ctx.CanCatch = catch != nil
	ctx.HasCaught = false

	destPtr := ctx.DestPtr.(*T)
	// 1. preTransforms
	for _, fn := range preTransforms {
		nVal, err := fn(ctx.Val, ctx)
		// bail if error in preTransform
		if err != nil || ctx.HasCaught {
			if ctx.CanCatch {
				ctx.HasCaught = true
				break
			}
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		ctx.Val = nVal
	}
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range postTransforms {
				err := fn(destPtr, ctx)
				if err != nil {
					// TODO unsure if we should also catch here
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	if ctx.HasCaught {
		*destPtr = *catch
		return
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
				ctx.HasCaught = true
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
				ctx.HasCaught = true
			} else {
				ctx.AddIssue(ctx.IssueFromCoerce(err))
				return
			}
		}
	}

	if ctx.HasCaught {
		*destPtr = *catch
		return
	}
	// 3. tests
	for _, test := range tests {
		if !test.ValidateFunc(destPtr, ctx) {
			// catching the first error if catch is set
			if ctx.CanCatch {
				*destPtr = *catch
				ctx.HasCaught = true
				break
			}
			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.Val))
		}
	}

	// 4. postTransforms -> Done above on defer
}

func primitiveValidator[T p.ZogPrimitive](ctx *p.SchemaCtx, preTransforms []p.PreTransform, tests []p.Test, postTransforms []p.PostTransform, defaultVal *T, required *p.Test, catch *T) {
	ctx.CanCatch = catch != nil
	ctx.HasCaught = false

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
		if err != nil {
			if ctx.CanCatch {
				ctx.HasCaught = true
				break
			}
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		*valPtr = nVal.(T)
	}

	if ctx.HasCaught {
		*valPtr = *catch
		return
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
				ctx.HasCaught = true
			} else {
				ctx.AddIssue(ctx.IssueFromTest(required, ctx.Val))
				return
			}
		}
	}

	if ctx.HasCaught {
		*valPtr = *catch
		return
	}

	// 3. tests
	for _, test := range tests {
		if !test.ValidateFunc(valPtr, ctx) {
			// catching the first error if catch is set
			if ctx.CanCatch {
				ctx.HasCaught = true
				*valPtr = *catch
				break
			}
			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.Val))
		}
	}

}
