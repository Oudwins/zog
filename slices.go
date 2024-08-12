package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/primitives"
)

type sliceProcessor struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	schema         Processor
	postTransforms []p.PostTransform
	required       *p.Test
	defaultVal     any
	// catch          any
}

func Slice(schema Processor) *sliceProcessor {
	return &sliceProcessor{
		schema: schema,
		tests:  []p.Test{},
	}
}

// only supports val = slice[any] & dest = &slice[]
func (v *sliceProcessor) Parse(val any, dest any) p.ZogSchemaErrors {
	var ctx = p.NewParseCtx()
	errs := p.NewErrsMap()
	path := p.PathBuilder("")
	v.process(val, dest, errs, path, ctx)

	return errs.M
}

func (v *sliceProcessor) process(val any, dest any, errs p.ZogErrors, path p.PathBuilder, ctx *p.ParseCtx) {
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(val, ctx)
			// bail if error in preTransform
			if err != nil {
				errs.Add(path, Errors.WrapUnknown(err))
				return
			}
			val = nVal
		}
	}

	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if errs.IsEmpty() {
			for _, fn := range v.postTransforms {
				err := fn(dest, ctx)
				if err != nil {
					errs.Add(path, Errors.WrapUnknown(err))
					return
				}
			}
		}
	}()

	// 2. cast data to string & handle default/required
	isZeroVal := p.IsZeroValue(val)
	destVal := reflect.ValueOf(dest).Elem()
	// WHAT IF THIS IS NOT A SLICE???? TODO ! FUCK in default we set an invalid type
	var refVal reflect.Value

	if isZeroVal {
		if v.defaultVal != nil {
			refVal = reflect.ValueOf(v.defaultVal)
		} else if v.required == nil {
			return
		}
	} else {
		// make sure val is a slice if not try to make it one
		v, err := conf.Coercers["slice"](val)
		if err != nil {
			errs.Add(path, Errors.Wrap(err, "failed to validate field"))
		}
		refVal = reflect.ValueOf(v)
	}

	destVal.Set(reflect.MakeSlice(destVal.Type(), refVal.Len(), refVal.Len()))

	// required
	if v.required != nil && !v.required.ValidateFunc(dest, ctx) {
		errs.Add(path, Errors.New(v.required.ErrorFunc(dest, ctx)))
	}

	// 3.1 tests for slice items
	if v.schema != nil {
		for idx := 0; idx < refVal.Len(); idx++ {
			item := refVal.Index(idx).Interface()
			ptr := destVal.Index(idx).Addr().Interface()
			p := path.Push(fmt.Sprint(idx))
			v.schema.process(item, ptr, errs, p, ctx)
		}
	}

	// 3. tests for slice
	for _, test := range v.tests {
		if !test.ValidateFunc(dest, ctx) {
			// catching the first error if catch is set
			// if v.catch != nil {
			// 	dest = v.catch
			// 	break
			// }
			//
			errs.Add(path, Errors.New(test.ErrorFunc(dest, ctx)))
		}
	}
	// 4. postTransforms -> defered see above
}

// !MODIFIERS

// marks field as required
func (v *sliceProcessor) Required(options ...TestOption) *sliceProcessor {
	r := p.Required(p.DErrorFunc("is a required field"))
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *sliceProcessor) Optional() *sliceProcessor {
	v.required = nil
	return v
}

// sets the default value
func (v *sliceProcessor) Default(val any) *sliceProcessor {
	v.defaultVal = val
	return v
}

// NOT IMPLEMENTED YET
// sets the catch value (i.e the value to use if the validation fails)
// func (v *sliceProcessor) Catch(val string) *sliceProcessor {
// 	v.catch = &val
// 	return v
// }

// !TESTS

// custom test function call it -> schema.Test("test_name", z.Message(""), func(val any, ctx *p.ParseCtx) bool {return true})
func (v *sliceProcessor) Test(ruleName string, errorMsg TestOption, validateFunc p.TestFunc) *sliceProcessor {
	v.tests = append(v.tests, p.Test{
		Name:         ruleName,
		ErrorFunc:    nil,
		ValidateFunc: validateFunc,
	})
	errorMsg(&v.tests[len(v.tests)-1])
	return v
}

// Minimum number of items
func (v *sliceProcessor) Min(n int, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		sliceMin(n, fmt.Sprintf("should be at least %d items long", n)),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}

	return v
}

// Maximum number of items
func (v *sliceProcessor) Max(n int, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		sliceMax(n, fmt.Sprintf("should be at maximum %d items long", n)),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

// Exact number of items
func (v *sliceProcessor) Len(n int, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		sliceLength(n, fmt.Sprintf("should be exactly %d items long", n)),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

// Slice contains a specific value
func (v *sliceProcessor) Contains(value any, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		p.Test{
			Name:      "contains",
			ErrorFunc: p.DErrorFunc(fmt.Sprintf("should contain %v", value)),
			ValidateFunc: func(val any, ctx *p.ParseCtx) bool {
				rv := reflect.ValueOf(val).Elem()
				if rv.Kind() != reflect.Slice {
					return false
				}
				for idx := 0; idx < rv.Len(); idx++ {
					v := rv.Index(idx).Interface()

					if reflect.DeepEqual(v, value) {
						return true
					}
				}

				return false
			},
		},
	)

	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

func sliceMin(n int, errMsg string) p.Test {
	return p.Test{
		Name:      "sliceMin",
		ErrorFunc: p.DErrorFunc(errMsg),
		ValidateFunc: func(val any, ctx *p.ParseCtx) bool {
			rv := reflect.ValueOf(val).Elem()
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() >= n
		},
	}
}
func sliceMax(n int, errMsg string) p.Test {
	return p.Test{
		Name:      "sliceMax",
		ErrorFunc: p.DErrorFunc(errMsg),
		ValidateFunc: func(val any, ctx *p.ParseCtx) bool {
			rv := reflect.ValueOf(val).Elem()
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() <= n
		},
	}
}
func sliceLength(n int, errMsg string) p.Test {
	return p.Test{
		Name:      "sliceLength",
		ErrorFunc: p.DErrorFunc(errMsg),
		ValidateFunc: func(val any, ctx *p.ParseCtx) bool {
			rv := reflect.ValueOf(val).Elem()
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() == n
		},
	}
}
