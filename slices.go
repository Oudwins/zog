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
	tests          []Test
	schema         ZogSchema
	postTransforms []PostTransform
	required       *Test
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

// Validates a slice
func (v *SliceSchema) Validate(data any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()

	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, data, path, v.getType())
	defer sctx.Free()
	v.validate(sctx)
	return errs.M
}

// Internal function to validate the data
func (v *SliceSchema) validate(ctx *p.SchemaCtx) {
	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range v.postTransforms {
				err := fn(ctx.ValPtr, ctx)
				if err != nil {
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	refVal := reflect.ValueOf(ctx.ValPtr).Elem() // we use this to set the value to the ptr. But we still reference the ptr everywhere. This is correct even if it seems confusing.
	// 2. cast data to string & handle default/required
	isZeroVal := p.IsZeroValue(ctx.ValPtr)

	if isZeroVal || refVal.Len() == 0 {
		if v.defaultVal != nil {
			refVal.Set(reflect.ValueOf(v.defaultVal))
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.ValPtr))
			return
		}
	}

	// 3.1 tests for slice items
	subCtx := ctx.NewValidateSchemaCtx(ctx.ValPtr, ctx.Path, v.schema.getType())
	defer subCtx.Free()
	for idx := 0; idx < refVal.Len(); idx++ {
		item := refVal.Index(idx).Addr().Interface()
		k := fmt.Sprintf("[%d]", idx)
		subCtx.ValPtr = item
		subCtx.Path.Push(&k)
		subCtx.Exit = false
		v.schema.validate(subCtx)
		subCtx.Path.Pop()
	}

	// 3. tests for slice
	for _, test := range v.tests {
		ctx.Test = &test
		test.Func(ctx.ValPtr, ctx)
		if ctx.Exit {
			// catching the first error if catch is set
			// if v.catch != nil {
			// 	dest = v.catch
			// 	break
			// }
			//
			return
		}
	}
	// 4. postTransforms -> defered see above
}

// Only supports parsing from data=slice[any] to a dest =&slice[] (this can be typed. Doesn't have to be any)
func (v *SliceSchema) Parse(data any, dest any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, dest, path, v.getType())
	defer sctx.Free()
	v.process(sctx)

	return errs.M
}

// Internal function to process the data
func (v *SliceSchema) process(ctx *p.SchemaCtx) {

	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range v.postTransforms {
				err := fn(ctx.ValPtr, ctx)
				if err != nil {
					ctx.AddIssue(ctx.IssueFromUnknownError(err))
					return
				}
			}
		}
	}()

	// 2. cast data to string & handle default/required
	isZeroVal := p.IsParseZeroValue(ctx.Data, ctx)
	var refVal reflect.Value

	if isZeroVal {
		if v.defaultVal != nil {
			refVal = reflect.ValueOf(v.defaultVal)
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Data))
			return
		}
	} else {
		// make sure val is a slice if not try to make it one
		v, err := v.coercer(ctx.Data)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromCoerce(err))
			return
		}
		refVal = reflect.ValueOf(v)
	}

	destVal := reflect.ValueOf(ctx.ValPtr).Elem()
	destVal.Set(reflect.MakeSlice(destVal.Type(), refVal.Len(), refVal.Len()))

	// 3.1 tests for slice items
	subCtx := ctx.NewSchemaCtx(ctx.Data, ctx.ValPtr, ctx.Path, v.schema.getType())
	defer subCtx.Free()
	for idx := 0; idx < refVal.Len(); idx++ {
		item := refVal.Index(idx).Interface()
		ptr := destVal.Index(idx).Addr().Interface()
		k := fmt.Sprintf("[%d]", idx)
		subCtx.Data = item
		subCtx.ValPtr = ptr
		subCtx.Path.Push(&k)
		v.schema.process(subCtx)
		subCtx.Path.Pop()
	}

	// 3. tests for slice
	for _, test := range v.tests {
		ctx.Test = &test
		test.Func(ctx.ValPtr, ctx)
		if ctx.Exit {
			// catch here
			return
		}
	}
	// 4. postTransforms -> defered see above
}

// Adds posttransform function to schema
func (v *SliceSchema) PostTransform(transform PostTransform) *SliceSchema {
	if v.postTransforms == nil {
		v.postTransforms = []PostTransform{}
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

// custom test function call it -> schema.Test(t z.Test)
func (v *SliceSchema) Test(t Test) *SliceSchema {
	v.tests = append(v.tests, t)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *SliceSchema) TestFunc(testFunc BoolTFunc, opts ...TestOption) *SliceSchema {
	t := p.NewTestFunc("", testFunc, opts...)
	v.Test(*t)
	return v
}

// Minimum number of items
func (v *SliceSchema) Min(n int, options ...TestOption) *SliceSchema {
	t, fn := sliceMin(n)
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Maximum number of items
func (v *SliceSchema) Max(n int, options ...TestOption) *SliceSchema {
	t, fn := sliceMax(n)
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Exact number of items
func (v *SliceSchema) Len(n int, options ...TestOption) *SliceSchema {
	t, fn := sliceLength(n)
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Slice contains a specific value
func (v *SliceSchema) Contains(value any, options ...TestOption) *SliceSchema {
	fn := func(val any, ctx Ctx) bool {
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
	}
	t := Test{
		IssueCode: zconst.IssueCodeContains,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeContains] = value
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

func sliceMin(n int) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Slice {
			return false
		}
		return rv.Len() >= n
	}

	t := Test{
		IssueCode: zconst.IssueCodeMin,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMin] = n
	return t, fn
}

func sliceMax(n int) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Slice {
			return false
		}
		return rv.Len() <= n
	}

	t := Test{
		IssueCode: zconst.IssueCodeMax,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMax] = n
	return t, fn
}
func sliceLength(n int) (Test, BoolTFunc) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Slice {
			return false
		}
		return rv.Len() == n
	}
	t := Test{
		IssueCode: zconst.IssueCodeLen,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeLen] = n
	return t, fn
}
