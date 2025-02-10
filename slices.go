package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// ! INTERNALS
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

// Returns the type of the schema
func (v *SliceSchema) getType() zconst.ZogType {
	return zconst.TypeSlice
}

// Sets the coercer for the schema
func (v *SliceSchema) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Creates a slice schema. That is a Zog representation of a slice.
// It takes a ZogSchema which will be used to validate against all the items in the slice.
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

// Validates a pointer pointer
func (v *SliceSchema) Validate(data any, options ...ExecOption) p.ZogErrMap {
	errs := p.NewErrsMap()

	ctx := p.NewExecCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}

	v.validate(ctx.NewValidateSchemaCtx(data, p.PathBuilder(""), v.getType()))
	return errs.M
}

// Internal function to validate the data
func (v *SliceSchema) validate(ctx *p.SchemaCtx) {
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range v.postTransforms {
				err := fn(ctx.Val, ctx)
				if err != nil {
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	refVal := reflect.ValueOf(ctx.Val).Elem() // we use this to set the value to the ptr. But we still reference the ptr everywhere. This is correct even if it seems confusing.
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(refVal.Interface(), ctx)
			// bail if error in preTransform
			if err != nil {
				ctx.AddIssue(ctx.IssueFromUnknownError(err))
				return
			}
			refVal.Set(reflect.ValueOf(nVal))
		}
	}

	// 2. cast data to string & handle default/required
	isZeroVal := p.IsZeroValue(ctx.Val)

	if isZeroVal || refVal.Len() == 0 {
		if v.defaultVal != nil {
			refVal.Set(reflect.ValueOf(v.defaultVal))
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
			return
		}
	}

	// 3.1 tests for slice items
	for idx := 0; idx < refVal.Len(); idx++ {
		item := refVal.Index(idx).Addr().Interface()
		path := ctx.Path.Push(fmt.Sprintf("[%d]", idx))
		v.schema.validate(ctx.NewValidateSchemaCtx(item, path, v.schema.getType()))
	}

	// 3. tests for slice
	for _, test := range v.tests {
		if !test.ValidateFunc(ctx.Val, ctx) {
			// catching the first error if catch is set
			// if v.catch != nil {
			// 	dest = v.catch
			// 	break
			// }
			//
			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.Val))
		}
	}
	// 4. postTransforms -> defered see above
}

// Only supports parsing from data=slice[any] to a dest =&slice[] (this can be typed. Doesn't have to be any)
func (v *SliceSchema) Parse(data any, dest any, options ...ExecOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewExecCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")
	v.process(ctx.NewSchemaCtx(data, dest, path, v.getType()))

	return errs.M
}

// Internal function to process the data
func (v *SliceSchema) process(ctx *p.SchemaCtx) {
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(ctx.Val, ctx)
			// bail if error in preTransform
			if err != nil {
				ctx.AddIssue(ctx.IssueFromUnknownError(err))
				return
			}
			ctx.Val = nVal
		}
	}

	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range v.postTransforms {
				err := fn(ctx.DestPtr, ctx)
				if err != nil {
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	// 2. cast data to string & handle default/required
	isZeroVal := p.IsParseZeroValue(ctx.Val, ctx)
	destVal := reflect.ValueOf(ctx.DestPtr).Elem()
	var refVal reflect.Value

	if isZeroVal {
		if v.defaultVal != nil {
			refVal = reflect.ValueOf(v.defaultVal)
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
			return
		}
	} else {
		// make sure val is a slice if not try to make it one
		v, err := v.coercer(ctx.Val)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromCoerce(err))
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
			path := ctx.Path.Push(fmt.Sprintf("[%d]", idx))

			v.schema.process(ctx.NewSchemaCtx(item, ptr, path, v.schema.getType()))
		}
	}

	// 3. tests for slice
	for _, test := range v.tests {
		if !test.ValidateFunc(ctx.DestPtr, ctx) {
			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.DestPtr))
		}
	}
	// 4. postTransforms -> defered see above
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
