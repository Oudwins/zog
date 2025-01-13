package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ComplexZogSchema = &SliceSchema{}

type SliceSchema struct {
	preTransforms  []p.PreTransform
	tests          []p.Test
	schema         ZogSchema
	postTransforms []p.PostTransform
	required       *p.Test
	defaultVal     any
	// catch          any
	coercer conf.CoercerFunc
}

// ! INTERNALS

// Returns the type of the schema
func (v *SliceSchema) getType() zconst.ZogType {
	return zconst.TypeSlice
}

// Sets the coercer for the schema
func (v *SliceSchema) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// Internal function to process the data
func (v *SliceSchema) process(val any, dest any, path p.PathBuilder, ctx ParseCtx) {
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
	isZeroVal := p.IsParseZeroValue(val, ctx)
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
		v, err := v.coercer(val)
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

func Slice(schema ZogSchema, opts ...SchemaOption) *SliceSchema {
	s := &SliceSchema{
		schema:  schema,
		coercer: conf.Coercers.Slice, // default coercer
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// only supports val = slice[any] & dest = &slice[]
func (v *SliceSchema) Parse(data any, dest any, options ...ParsingOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")
	v.process(data, dest, path, ctx)

	return errs.M
}

// Adds pretransform function to schema
func (v *SliceSchema) PreTransform(transform p.PreTransform) *SliceSchema {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *SliceSchema) PostTransform(transform p.PostTransform) *SliceSchema {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// !MODIFIERS

// marks field as required
func (v *SliceSchema) Required(options ...TestOption) *SliceSchema {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *SliceSchema) Optional() *SliceSchema {
	v.required = nil
	return v
}

// sets the default value
func (v *SliceSchema) Default(val any) *SliceSchema {
	v.defaultVal = val
	return v
}

// NOT IMPLEMENTED YET
// sets the catch value (i.e the value to use if the validation fails)
// func (v *SliceSchema) Catch(val string) *SliceSchema {
// 	v.catch = &val
// 	return v
// }

// !TESTS

// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *SliceSchema) Test(t p.Test, opts ...TestOption) *SliceSchema {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Minimum number of items
func (v *SliceSchema) Min(n int, options ...TestOption) *SliceSchema {
	v.tests = append(v.tests,
		sliceMin(n),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}

	return v
}

// Maximum number of items
func (v *SliceSchema) Max(n int, options ...TestOption) *SliceSchema {
	v.tests = append(v.tests,
		sliceMax(n),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

// Exact number of items
func (v *SliceSchema) Len(n int, options ...TestOption) *SliceSchema {
	v.tests = append(v.tests,
		sliceLength(n),
	)
	for _, opt := range options {
		opt(&v.tests[len(v.tests)-1])
	}
	return v
}

// Slice contains a specific value
func (v *SliceSchema) Contains(value any, options ...TestOption) *SliceSchema {
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
