package zog

import (
	"log"
	"reflect"
	"time"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type Processor interface {
	process(val any, dest any, path p.PathBuilder, ctx ParseCtx)
}

// ! Passing Types through

type ParseCtx = p.ParseCtx
type Test = p.Test

type ZogError = p.ZogError
type ZogErrMap = p.ZogErrMap
type ZogErrList = p.ZogErrList

// ! TESTS
func TestFunc(errCode zconst.ZogErrCode, validateFunc p.TestFunc) p.Test {
	t := p.Test{
		ErrCode:      errCode,
		ValidateFunc: validateFunc,
	}
	return t
}

// ! ERRORS
type errHelpers struct {
}

// Beware this API may change
var Errors = errHelpers{}

// Create error from (originValue any, destinationValue any, test *p.Test)
func (e *errHelpers) FromTest(o any, destType zconst.ZogType, t *p.Test, p ParseCtx) p.ZogError {
	er := e.New(t.ErrCode, o, destType, t.Params, "", nil)
	if t.ErrFmt != nil {
		t.ErrFmt(er, p)
	}
	return er
}

// Create error from
func (e *errHelpers) FromErr(o any, destType zconst.ZogType, err error) p.ZogError {
	return e.New(zconst.ErrCodeCustom, o, destType, nil, "", err)
}

func (e *errHelpers) WrapUnknown(o any, destType zconst.ZogType, err error) p.ZogError {
	zerr, ok := err.(p.ZogError)
	if !ok {
		return e.FromErr(o, destType, err)
	}
	return zerr
}

func (e *errHelpers) New(code zconst.ZogErrCode, o any, destType zconst.ZogType, params map[string]any, msg string, err error) p.ZogError {
	return &p.ZogErr{
		C:       code,
		ParamsM: params,
		Val:     o,
		Typ:     destType,
		Msg:     msg,
		Err:     err,
	}
}

func (e *errHelpers) SanitizeMap(m p.ZogErrMap) map[string][]string {
	errs := make(map[string][]string, len(m))
	for k, v := range m {
		errs[k] = e.SanitizeList(v)
	}
	return errs
}

func (e *errHelpers) SanitizeList(l p.ZogErrList) []string {
	errs := make([]string, len(l))
	for i, err := range l {
		errs[i] = err.Message()
	}
	return errs
}

// ! Data Providers

// Creates a new map data provider
func NewMapDataProvider[T any](m map[string]T) p.DataProvider {
	return p.NewMapDataProvider(m)
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
