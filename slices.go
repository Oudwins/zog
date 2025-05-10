package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// ! INTERNALS
var _ ComplexZogSchema = &SliceSchema[int]{}

type SliceSchema[T comparable] struct {
	processors []p.ZProcessor[[]T]
	schema     ZogSchema
	required   *p.Test[[]T]
	defaultVal []T
	// catch          any
	coercer conf.CoercerFunc
}

// Returns the type of the schema
func (v *SliceSchema[T]) getType() zconst.ZogType {
	return zconst.TypeSlice
}

// Sets the coercer for the schema
func (v *SliceSchema[T]) setCoercer(c conf.CoercerFunc) {
	v.coercer = c
}

// ! USER FACING FUNCTIONS

// Creates a slice schema. That is a Zog representation of a slice.
// It takes a ZogSchema which will be used to validate against all the items in the slice.
func Slice[T comparable](schema ZogSchema, opts ...SchemaOption) *SliceSchema[T] {
	s := &SliceSchema[T]{
		schema:  schema,
		coercer: conf.Coercers.Slice, // default coercer
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Validates a slice
func (v *SliceSchema[T]) Validate(data any, options ...ExecOption) ZogIssueMap {
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
func (v *SliceSchema[T]) validate(ctx *p.SchemaCtx) {

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
func (v *SliceSchema[T]) Parse(data any, dest any, options ...ExecOption) ZogIssueMap {
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
func (v *SliceSchema[T]) process(ctx *p.SchemaCtx) {

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
func (v *SliceSchema[T]) Transform(transform Transform[[]T]) *SliceSchema[T] {
	v.processors = append(v.processors, &p.TransformProcessor[[]T]{
		Transform: p.Transform[[]T](transform),
	})
	return v
}

// !MODIFIERS

// marks field as required
func (v *SliceSchema[T]) Required(options ...TestOption) *SliceSchema[T] {
	r := p.Required[[]T]()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *SliceSchema[T]) Optional() *SliceSchema[T] {
	v.required = nil
	return v
}

// sets the default value
func (v *SliceSchema[T]) Default(val []T) *SliceSchema[T] {
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
func (v *SliceSchema[T]) Test(t Test[[]T]) *SliceSchema[T] {
	x := p.Test[[]T](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *SliceSchema[T]) TestFunc(testFunc BoolTFunc[[]T], opts ...TestOption) *SliceSchema[T] {
	t := p.NewTestFunc("", p.BoolTFunc[[]T](testFunc), opts...)
	v.Test(Test[[]T](*t))
	return v
}

// Minimum number of items
func (v *SliceSchema[T]) Min(n int, options ...TestOption) *SliceSchema[T] {
	t, fn := sliceMin[T](n)
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.processors = append(v.processors, &t)
	return v
}

// Maximum number of items
func (v *SliceSchema[T]) Max(n int, options ...TestOption) *SliceSchema[T] {
	t, fn := sliceMax[T](n)
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.processors = append(v.processors, &t)
	return v
}

// Exact number of items
func (v *SliceSchema[T]) Len(n int, options ...TestOption) *SliceSchema[T] {
	t, fn := sliceLength[T](n)
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.processors = append(v.processors, &t)
	return v
}

// Slice contains a specific value
func (v *SliceSchema[T]) Contains(value T, options ...TestOption) *SliceSchema[T] {
	fn := func(val []T, ctx Ctx) bool {
		for _, v := range val {
			if v == value {
				return true
			}
		}

		return false
	}
	t := p.Test[[]T]{
		IssueCode: zconst.IssueCodeContains,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeContains] = value
	p.TestFuncFromBool(fn, &t)
	for _, opt := range options {
		opt(&t)
	}
	v.processors = append(v.processors, &t)
	return v
}

func sliceMin[T comparable](n int) (p.Test[[]T], p.BoolTFunc[[]T]) {
	fn := func(val []T, ctx Ctx) bool {
		return len(val) >= n
	}

	t := p.Test[[]T]{
		IssueCode: zconst.IssueCodeMin,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMin] = n
	return t, fn
}

func sliceMax[T comparable](n int) (p.Test[[]T], p.BoolTFunc[[]T]) {
	fn := func(val []T, ctx Ctx) bool {
		return len(val) <= n
	}

	t := p.Test[[]T]{
		IssueCode: zconst.IssueCodeMax,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeMax] = n
	return t, fn
}

func sliceLength[T comparable](n int) (p.Test[[]T], p.BoolTFunc[[]T]) {
	fn := func(val []T, ctx Ctx) bool {
		return len(val) == n
	}
	t := p.Test[[]T]{
		IssueCode: zconst.IssueCodeLen,
		Params:    make(map[string]any, 1),
	}
	t.Params[zconst.IssueCodeLen] = n
	return t, fn
}
