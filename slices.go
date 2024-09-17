package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
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
func (v *sliceProcessor) Parse(data any, dest any, options ...ParsingOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")
	v.process(data, dest, path, ctx)

	return errs.M
}

func (v *sliceProcessor) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
	destType := zconst.TypeSlice
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(val, ctx)
			// bail if error in preTransform
			if err != nil {
				ctx.NewError(path, Errors.WrapUnknown(val, destType, err))
				return
			}
			val = nVal
		}
	}

	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range v.postTransforms {
				err := fn(dest, ctx)
				if err != nil {
					ctx.NewError(path, Errors.WrapUnknown(val, destType, err))
					return
				}
			}
		}
	}()

	// 2. cast data to string & handle default/required
	isZeroVal := p.IsZeroValue(val)
	destVal := reflect.ValueOf(dest).Elem()
	var refVal reflect.Value

	if isZeroVal {
		if v.defaultVal != nil {
			refVal = reflect.ValueOf(v.defaultVal)
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.NewError(path, Errors.FromTest(val, destType, v.required, ctx))
			return
		}
	} else {
		// make sure val is a slice if not try to make it one
		v, err := conf.Coercers.Slice(val)
		if err != nil {
			ctx.NewError(path, Errors.New(zconst.ErrCodeCoerce, val, destType, nil, "", err))
			return
		}
		refVal = reflect.ValueOf(v)
	}

	destVal.Set(reflect.MakeSlice(destVal.Type(), refVal.Len(), refVal.Len()))

	// 3.1 tests for slice items
	if v.schema != nil {
		for idx := 0; idx < refVal.Len(); idx++ {
			item := refVal.Index(idx).Interface()
			ptr := destVal.Index(idx).Addr().Interface()
			path := path.Push(fmt.Sprintf("[%d]", idx))
			v.schema.process(item, ptr, path, ctx)
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
			ctx.NewError(path, Errors.FromTest(val, destType, &test, ctx))
		}
	}
	// 4. postTransforms -> defered see above
}

// Adds pretransform function to schema
func (v *sliceProcessor) PreTransform(transform p.PreTransform) *sliceProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *sliceProcessor) PostTransform(transform p.PostTransform) *sliceProcessor {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// !MODIFIERS

// marks field as required
func (v *sliceProcessor) Required(options ...TestOption) *sliceProcessor {
	r := p.Required()
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

// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *sliceProcessor) Test(t p.Test, opts ...TestOption) *sliceProcessor {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Minimum number of items
func (v *sliceProcessor) Min(n int, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		sliceMin(n),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}

	return v
}

// Maximum number of items
func (v *sliceProcessor) Max(n int, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		sliceMax(n),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

// Exact number of items
func (v *sliceProcessor) Len(n int, options ...TestOption) *sliceProcessor {
	v.tests = append(v.tests,
		sliceLength(n),
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
			ErrCode: zconst.ErrCodeContains,
			Params:  make(map[string]any, 1),
			ValidateFunc: func(val any, ctx ParseCtx) bool {
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
	v.tests[len(v.tests)-1].Params[zconst.ErrCodeContains] = value
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

func sliceMin(n int) p.Test {
	t := p.Test{
		ErrCode: zconst.ErrCodeMin,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			rv := reflect.ValueOf(val).Elem()
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() >= n
		},
	}
	t.Params[zconst.ErrCodeMin] = n
	return t
}
func sliceMax(n int) p.Test {
	t := p.Test{
		ErrCode: zconst.ErrCodeMax,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			rv := reflect.ValueOf(val).Elem()
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() <= n
		},
	}
	t.Params[zconst.ErrCodeMax] = n
	return t
}
func sliceLength(n int) p.Test {
	t := p.Test{
		ErrCode: zconst.ErrCodeLen,
		Params:  make(map[string]any, 1),
		ValidateFunc: func(val any, ctx ParseCtx) bool {
			rv := reflect.ValueOf(val).Elem()
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() == n
		},
	}
	t.Params[zconst.ErrCodeLen] = n
	return t
}
