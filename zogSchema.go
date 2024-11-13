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
type ZogSchema interface {
	process(val any, dest any, path p.PathBuilder, ctx ParseCtx)
	setCoercer(c conf.CoercerFunc)
	getType() zconst.ZogType
}

// ! Passing Types through

// ParseCtx is the context passed through the parser
type ParseCtx = p.ParseCtx

// ZogError is the ZogError interface
type ZogError = p.ZogError

// ZogErrList is a []ZogError returned from parsing primitive schemas
type ZogErrList = p.ZogErrList

// ZogErrMap is a map[string][]ZogError returned from parsing complex schemas
type ZogErrMap = p.ZogErrMap

// ! TESTS

// Test is the test object
type Test = p.Test

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
