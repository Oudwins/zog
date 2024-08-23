package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/primitives"
)

type Processor interface {
	process(val any, dest any, errs p.ZogErrors, path p.PathBuilder, ctx p.ParseCtx)
}

// ! Parse Context

type ParseCtx = p.ParseCtx

// ! ERRORS
type errHelpers struct {
}

// Creates a new error with a message. Assuming you only care about the message
func (e *errHelpers) New(msg string) p.ZogError {
	return p.ZogError{
		Message: msg,
	}
}

// Wraps an error with a message
func (e *errHelpers) Wrap(err error, msg string) p.ZogError {
	return p.ZogError{
		Message: msg,
		Err:     err,
	}
}

// Wraps an error but first checks if it is a zog error. Sets the message to equal the error message if it is not a zog error
func (e *errHelpers) WrapUnknown(err error) p.ZogError {
	zerr, ok := err.(p.ZogError)
	if !ok {
		return p.ZogError{
			Message: err.Error(),
			Err:     err,
		}
	}
	return zerr
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
		errs[i] = err.Message
	}
	return errs
}

var Errors = errHelpers{}

type ZogError = p.ZogError
type ZogErrMap = p.ZogErrMap
type ZogErrList = p.ZogErrList

// ! Data Providers

// Creates a new map data provider
func NewMapDataProvider[T any](m map[string]T) p.DataProvider {
	return p.NewMapDataProvider(m)
}

// ! PRIMITIVE PROCESSING

func primitiveProcessor[T p.ZogPrimitive](val any, dest any, errs p.ZogErrors, path p.PathBuilder, ctx p.ParseCtx, preTransforms []p.PreTransform, tests []p.Test, postTransforms []p.PostTransform, defaultVal *T, required *p.Test, catch *T, coercer conf.CoercerFunc) {
	canCatch := catch != nil
	hasCatched := false

	destPtr := dest.(*T)
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
			errs.Add(path, Errors.WrapUnknown(err))
			return
		}
		val = nVal
	}
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if errs.IsEmpty() {
			for _, fn := range postTransforms {
				err := fn(destPtr, ctx)
				if err != nil {
					errs.Add(path, Errors.WrapUnknown(err))
					return
				}
			}
		}
	}()

	if !hasCatched {
		// 2. cast data to string & handle default/required
		isZeroVal := p.IsZeroValue(val)

		if isZeroVal {
			if defaultVal != nil {
				*destPtr = *defaultVal
			} else if required == nil {
				// This handles optional case
				return
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
					errs.Add(path, Errors.Wrap(err, "failed to validate field"))
					return
				}
			}
		}
	}

	if !hasCatched {
		// required
		if required != nil && !required.ValidateFunc(*destPtr, ctx) {
			if catch != nil {
				*destPtr = *catch
				hasCatched = true
			} else {
				errs.Add(path, Errors.New(required.ErrorFunc(*destPtr, ctx)))
				return
			}
		}
	}

	if !hasCatched {
		// 3. tests
		for _, test := range tests {
			if !test.ValidateFunc(*destPtr, ctx) {
				// catching the first error if catch is set
				if catch != nil {
					*destPtr = *catch
					hasCatched = true //nolint
					break
				}
				//
				errs.Add(path, Errors.New(test.ErrorFunc(*destPtr, ctx)))
			}
		}
	}

	// 4. postTransforms -> Done above on defer
}
