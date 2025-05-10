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
	processors []p.ZProcessor[any]
	schema     ZogSchema
	required   *p.Test[any]
	defaultVal any
	// catch          any
	coercer conf.CoercerFunc
	isNot   bool
}

type NotSliceSchema interface {
	Len(n int, options ...TestOption) *SliceSchema
	Contains(value any, options ...TestOption) *SliceSchema
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

	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(ctx.ValPtr, ctx)
		if ctx.Exit {
			// can catch here
			return
		}
	}
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

	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(ctx.ValPtr, ctx)
		if ctx.Exit {
			return
		}
	}

}

// Adds transform function to schema.
func (v *SliceSchema) Transform(transform Transform[any]) *SliceSchema {
	v.processors = append(v.processors, &p.TransformProcessor[any]{
		Transform: p.Transform[any](transform),
	})
	return v
}

// !MODIFIERS

// marks field as required
func (v *SliceSchema) Required(options ...TestOption) *SliceSchema {
	r := p.Required[any]()
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
func (v *SliceSchema) Test(t Test[any]) *SliceSchema {
	x := p.Test[any](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *SliceSchema) TestFunc(testFunc BoolTFunc[any], opts ...TestOption) *SliceSchema {
	t := p.NewTestFunc("", p.BoolTFunc[any](testFunc), opts...)
	v.Test(Test[any](*t))
	return v
}

// Minimum number of items
func (v *SliceSchema) Min(n int, options ...TestOption) *SliceSchema {
	t, fn := sliceMin(n)

	return v.addTest(&t, fn, options...)
}

// Maximum number of items
func (v *SliceSchema) Max(n int, options ...TestOption) *SliceSchema {
	t, fn := sliceMax(n)

	return v.addTest(&t, fn, options...)
}

// Exact number of items
func (v *SliceSchema) Len(n int, options ...TestOption) *SliceSchema {
	t, fn := sliceLength(n)

	return v.addTest(&t, fn, options...)
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
	t := &p.Test[any]{
		IssueCode: zconst.IssueCodeContains,
		Params: map[string]any{
			zconst.IssueCodeContains: value,
		},
	}

	return v.addTest(t, fn, options...)
}

func sliceMin(n int) (p.Test[any], p.BoolTFunc[any]) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Slice {
			return false
		}
		return rv.Len() >= n
	}

	t := p.Test[any]{
		IssueCode: zconst.IssueCodeMin,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMin] = n
	return t, fn
}

func sliceMax(n int) (p.Test[any], p.BoolTFunc[any]) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Slice {
			return false
		}
		return rv.Len() <= n
	}

	t := p.Test[any]{
		IssueCode: zconst.IssueCodeMax,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMax] = n
	return t, fn
}
func sliceLength(n int) (p.Test[any], p.BoolTFunc[any]) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Slice {
			return false
		}
		return rv.Len() == n
	}
	t := p.Test[any]{
		IssueCode: zconst.IssueCodeLen,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeLen] = n
	return t, fn
}

func (v *SliceSchema) Not() NotSliceSchema {
	v.isNot = true
	return v
}

func (v *SliceSchema) addTest(t *p.Test[any], fn p.BoolTFunc[any], options ...TestOption) *SliceSchema {
	if v.isNot {
		p.TestNotFuncFromBool(fn, t)
		t.IssueCode = zconst.NotIssueCode(t.IssueCode)
		v.isNot = false
	} else {
		p.TestFuncFromBool(fn, t)
	}

	for _, opt := range options {
		opt(t)
	}

	v.processors = append(v.processors, t)
	return v
}
