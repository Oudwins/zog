package zog

import (
	"time"

	p "github.com/Oudwins/zog/primitives"
)

type Processor interface {
	process(val any, dest any, errs p.ZogErrors, path p.Pather, ctx *p.ParseCtx)
}

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

var Errors = errHelpers{}

// ! Tagger

// type tagger[T any] struct {
// 	Def     *T
// 	Req       *p.Test
// 	Catched          *T
// }

// func (v tagger[T]) Required() numberProcessor[T] {
// 	r := p.Required(p.DErrorFunc("is a required field"))
// 	v.Req = &r
// 	return v
// }

// // marks field as optional
// func (v numberProcessor[T]) Optional() numberProcessor[T] {
// 	v.required = nil
// 	return v
// }

// // sets the default value
// func (v numberProcessor[T]) Default(val T) numberProcessor[T] {
// 	v.defaultVal = &val
// 	return v
// }

// // sets the catch value (i.e the value to use if the validation fails)
// func (v numberProcessor[T]) Catch(val T) numberProcessor[T] {
// 	v.catch = &val
// 	return v
// }

// ! PRIMITIVE PROCESSING
type ZodPrimitive interface {
	~string | ~int | ~float64 | ~bool | time.Time
}

func primitiveProcess[T ZodPrimitive](val any, dest any, errs p.ZogErrors, path p.Pather, ctx *p.ParseCtx, preTransforms []p.PreTransform, tests []p.Test, postTransforms []p.PostTransform, defaultVal *T, required *p.Test, catch *T, coercer p.CoercerFunc) {

	destPtr := dest.(*T)
	// 1. preTransforms
	for _, fn := range preTransforms {
		nVal, err := fn(val, ctx)
		// bail if error in preTransform
		if err != nil {
			if catch != nil {
				*destPtr = *catch
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
		if err := coercer(val, destPtr); err != nil {
			if catch != nil {
				*destPtr = *catch
			} else {
				errs.Add(path, Errors.Wrap(err, "failed to validate string"))
				return
			}
		}
	}

	// required
	if required != nil && !required.ValidateFunc(*destPtr, ctx) {
		if catch != nil {
			*destPtr = *catch
		} else {
			errs.Add(path, Errors.New(required.ErrorFunc(*destPtr, ctx)))
			return
		}
	}

	// 3. tests
	for _, test := range tests {
		if !test.ValidateFunc(*destPtr, ctx) {
			// catching the first error if catch is set
			if catch != nil {
				*destPtr = *catch
				break
			}
			//
			errs.Add(path, Errors.New(test.ErrorFunc(*destPtr, ctx)))
		}
	}

	// 4. postTransforms -> Done above on defer
}
