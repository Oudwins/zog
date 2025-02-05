package zog

import (
	"log"
	"reflect"
	"time"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// The ZogSchema is the interface all schemas must implement
// This is most useful for internal use. If you are looking to pass schemas around, use the ComplexZogSchema or PrimitiveZogSchema interfaces if possible.
type ZogSchema interface {
	process(val any, dest any, path p.PathBuilder, ctx ParseCtx)
	validate(val any, path p.PathBuilder, ctx ParseCtx)
	setCoercer(c conf.CoercerFunc)
	getType() zconst.ZogType
}

// This is a common interface for all complex schemas (i.e structs, slices, pointers...)
// You can use this to pass any complex schema around
type ComplexZogSchema interface {
	ZogSchema
	Parse(val any, dest any, options ...ParsingOption) ZogErrMap
}

// This is a common interface for all primitive schemas (i.e strings, numbers, booleans, time.Time...)
// You can use this to pass any primitive schema around
type PrimitiveZogSchema[T p.ZogPrimitive] interface {
	ZogSchema
	Parse(val any, dest *T, options ...ParsingOption) ZogErrList
}

// ! PRIMITIVE PROCESSING

func getDestType(dest any) zconst.ZogType {
	switch reflect.TypeOf(dest).Kind() {
	case reflect.Slice:
		return zconst.TypeSlice
	case reflect.Struct:
		if reflect.TypeOf(dest) == reflect.TypeOf(time.Time{}) {
			return zconst.TypeTime
		}
		return zconst.TypeStruct
	case reflect.Float64:
	case reflect.Int:
		return zconst.TypeNumber
	case reflect.Bool:
		return zconst.TypeBool
	case reflect.String:
		return zconst.TypeString
	default:
		log.Fatal("Unsupported destination type")
	}
	// should never get here
	return zconst.TypeString
}

func primitiveProcessor[T p.ZogPrimitive](val any, dest any, path p.PathBuilder, ctx ParseCtx, preTransforms []p.PreTransform, tests []p.Test, postTransforms []p.PostTransform, defaultVal *T, required *p.Test, catch *T, coercer conf.CoercerFunc, isZeroFunc p.IsZeroValueFunc) {
	canCatch := catch != nil
	hasCatched := false

	destPtr := dest.(*T)
	destType := getDestType(*destPtr)
	// 1. preTransforms
	for _, fn := range preTransforms {
		nVal, err := fn(val, ctx)
		// bail if error in preTransform
		if err != nil {
			if canCatch {
				*destPtr = *catch
				hasCatched = true
				break
			}
			ctx.NewError(path, Errors.WrapUnknown(val, destType, err))
			return
		}
		val = nVal
	}
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range postTransforms {
				err := fn(destPtr, ctx)
				if err != nil {
					ctx.NewError(path, Errors.WrapUnknown(val, destType, err))
					return
				}
			}
		}
	}()

	if hasCatched {
		return
	}

	// 2. cast data to string & handle default/required
	isZeroVal := isZeroFunc(val, ctx)

	if isZeroVal {
		if defaultVal != nil {
			*destPtr = *defaultVal
		} else if required == nil {
			// This handles optional case
			return
		} else {
			// is required & zero value
			// required
			if catch != nil {
				*destPtr = *catch
				hasCatched = true
			} else {
				ctx.NewError(path, Errors.FromTest(val, destType, required, ctx))
				return
			}
		}
	} else {
		newVal, err := coercer(val)
		if err == nil {
			*destPtr = newVal.(T)
		} else {
			if canCatch {
				*destPtr = *catch
				hasCatched = true
			} else {
				ctx.NewError(path, Errors.New(zconst.ErrCodeCoerce, val, destType, nil, "", err))
				return
			}
		}
	}

	if hasCatched {
		return
	}
	// 3. tests
	for _, test := range tests {
		if !test.ValidateFunc(*destPtr, ctx) {
			// catching the first error if catch is set
			if catch != nil {
				*destPtr = *catch
				hasCatched = true //nolint
				break
			}
			ctx.NewError(path, Errors.FromTest(val, destType, &test, ctx))
		}
	}

	// 4. postTransforms -> Done above on defer
}

func primitiveValidator[T p.ZogPrimitive](val any, path p.PathBuilder, ctx ParseCtx, preTransforms []p.PreTransform, tests []p.Test, postTransforms []p.PostTransform, defaultVal *T, required *p.Test, catch *T) {

	canCatch := catch != nil

	valPtr := val.(*T)
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range postTransforms {
				err := fn(val, ctx)
				if err != nil {
					ctx.NewError(path, Errors.WrapUnknown(val, zconst.TypeBool, err))
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
			if canCatch {
				*valPtr = *catch
				return
			}
			ctx.NewError(path, Errors.WrapUnknown(val, zconst.TypeBool, err))
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
			if catch != nil {
				*valPtr = *catch
				return
			} else {
				ctx.NewError(path, Errors.FromTest(val, zconst.TypeBool, required, ctx))
				return
			}
		}
	}
	// 3. tests
	for _, test := range tests {
		if !test.ValidateFunc(*valPtr, ctx) {
			// catching the first error if catch is set
			if canCatch {
				*valPtr = *catch
				return
			}
			ctx.NewError(path, Errors.FromTest(val, zconst.TypeBool, &test, ctx))
		}
	}

}
