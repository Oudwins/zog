package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/zconst"
)

// ! INTERNALS
var _ ComplexZogSchema = &MapSchema[string, any]{}

type MapSchema[K p.ZogPrimitive, V any] struct {
	processors  []p.ZProcessor[any]
	keySchema   PrimitiveZogSchema[K]
	valueSchema ZogSchema
	required    *p.Test[any]
	defaultFunc func() map[K]V
}

// Returns the type of the schema
func (v *MapSchema[K, V]) getType() zconst.ZogType {
	return zconst.TypeMap
}

// Sets the coercer for the schema (no-op for maps, kept for interface compliance)
func (v *MapSchema[K, V]) setCoercer(c conf.CoercerFunc) {
	// Maps don't need coercers as we use generics
}

// ! USER FACING FUNCTIONS

// Creates a map schema. That is a Zog representation of a map.
// It takes a PrimitiveZogSchema for keys and a ZogSchema for values.
func EXPERIMENTAL_MAP[K p.ZogPrimitive, V any](keySchema PrimitiveZogSchema[K], valueSchema ZogSchema, opts ...SchemaOption) *MapSchema[K, V] {
	s := &MapSchema[K, V]{
		keySchema:   keySchema,
		valueSchema: valueSchema,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Validates a map
func (v *MapSchema[K, V]) Validate(data *map[K]V, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
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
	return errs.List
}

// Internal function to validate the data
func (v *MapSchema[K, V]) validate(ctx *p.SchemaCtx) {
	refVal := reflect.ValueOf(ctx.ValPtr).Elem()
	isZeroVal := p.IsZeroValue(ctx.ValPtr)

	if isZeroVal || refVal.Len() == 0 {
		if v.defaultFunc != nil {
			refVal.Set(reflect.ValueOf(v.defaultFunc()))
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.ValPtr))
			return
		}
	}

	// Create contexts once for reuse
	keySubCtx := ctx.NewValidateSchemaCtx(ctx.ValPtr, ctx.Path, v.keySchema.getType())
	defer keySubCtx.Free()
	subCtx := ctx.NewValidateSchemaCtx(ctx.ValPtr, ctx.Path, v.valueSchema.getType())
	defer subCtx.Free()

	// Validate map entries
	for _, key := range refVal.MapKeys() {
		keyVal := key.Interface().(K)
		k := fmt.Sprintf(`["%v"]`, keyVal)

		// Validate key
		keyPtr := &keyVal
		keySubCtx.ValPtr = keyPtr
		keySubCtx.Path.Push(&k)
		keySubCtx.Exit = false
		v.keySchema.validate(keySubCtx)
		keySubCtx.Path.Pop()

		// Validate value - use map's value type, not runtime type of the value
		// This ensures proper handling of interface types like `any`
		valueType := refVal.Type().Elem()
		valuePtr := reflect.New(valueType).Interface()
		reflect.ValueOf(valuePtr).Elem().Set(refVal.MapIndex(key))
		subCtx.ValPtr = valuePtr
		subCtx.Path.Push(&k)
		subCtx.Exit = false
		v.valueSchema.validate(subCtx)
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

// Parse the data into the destination map
func (v *MapSchema[K, V]) Parse(data any, dest any, options ...ExecOption) ZogIssueList {
	errs := p.NewErrsList()
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

	return errs.List
}

// Internal function to process the data
func (v *MapSchema[K, V]) process(ctx *p.SchemaCtx) {
	isZeroVal := p.IsParseZeroValue(ctx.Data, ctx)
	var inputMap reflect.Value

	if isZeroVal {
		if v.defaultFunc != nil {
			inputMap = reflect.ValueOf(v.defaultFunc())
		} else if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Data))
			return
		}
	} else {
		// Verify input is a map
		inputMap = reflect.ValueOf(ctx.Data)
		if inputMap.Kind() != reflect.Map {
			ctx.AddIssue(ctx.Issue().SetCode(zconst.IssueCodeCoerce).SetMessage(fmt.Sprintf("expected map, got %v", inputMap.Kind())))
			return
		}
	}

	// Create destination map
	destVal := reflect.ValueOf(ctx.ValPtr).Elem()
	destType := destVal.Type()
	if destType.Kind() != reflect.Map {
		p.Panicf(p.PanicTypeCast, ctx.String(), ctx.DType, ctx.ValPtr)
	}
	destMap := reflect.MakeMap(destType)

	// Create contexts once for reuse
	keySubCtx := ctx.NewSchemaCtx(ctx.Data, ctx.ValPtr, ctx.Path, v.keySchema.getType())
	defer keySubCtx.Free()
	subCtx := ctx.NewSchemaCtx(ctx.Data, ctx.ValPtr, ctx.Path, v.valueSchema.getType())
	defer subCtx.Free()

	// Process map entries
	for _, key := range inputMap.MapKeys() {
		keyData := key.Interface()
		valueData := inputMap.MapIndex(key).Interface()
		k := fmt.Sprintf(`["%v"]`, keyData)

		// Parse key - create a zero value and get pointer to it
		keyPtr := reflect.New(destType.Key()).Interface()
		keySubCtx.Data = keyData
		keySubCtx.ValPtr = keyPtr
		keySubCtx.Path.Push(&k)
		keySubCtx.Exit = false
		v.keySchema.process(keySubCtx)
		keySubCtx.Path.Pop()
		if keySubCtx.Exit {
			continue
		}
		parsedKey := reflect.ValueOf(keyPtr).Elem().Interface().(K)

		// Parse value - create a zero value and get pointer to it
		valuePtr := reflect.New(destType.Elem()).Interface()
		subCtx.Data = valueData
		subCtx.ValPtr = valuePtr
		subCtx.Path.Push(&k)
		subCtx.Exit = false
		v.valueSchema.process(subCtx)
		subCtx.Path.Pop()

		// Only add to map if no errors occurred
		// Use reflect.Value directly to avoid type assertion issues with nil interfaces
		if !subCtx.Exit {
			parsedValueReflect := reflect.ValueOf(valuePtr).Elem()
			destMap.SetMapIndex(reflect.ValueOf(parsedKey), parsedValueReflect)
		}
	}

	destVal.Set(destMap)

	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(ctx.ValPtr, ctx)
		if ctx.Exit {
			return
		}
	}
}

// Adds transform function to schema.
func (v *MapSchema[K, V]) Transform(transform Transform[any]) *MapSchema[K, V] {
	v.processors = append(v.processors, &p.TransformProcessor[any]{
		Transform: p.Transform[any](transform),
	})
	return v
}

// !MODIFIERS

// marks field as required
func (v *MapSchema[K, V]) Required(options ...TestOption) *MapSchema[K, V] {
	r := p.Required[any]()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *MapSchema[K, V]) Optional() *MapSchema[K, V] {
	v.required = nil
	return v
}

// sets the default value
func (v *MapSchema[K, V]) Default(val map[K]V) *MapSchema[K, V] {
	return v.DefaultFunc(func() map[K]V {
		return val
	})
}

// sets the default value using a function
func (v *MapSchema[K, V]) DefaultFunc(defaultFunc func() map[K]V) *MapSchema[K, V] {
	v.defaultFunc = defaultFunc
	return v
}

// !TESTS

// custom test function call it -> schema.Test(t z.Test)
func (v *MapSchema[K, V]) Test(t Test[any]) *MapSchema[K, V] {
	x := p.Test[any](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *MapSchema[K, V]) TestFunc(testFunc BoolTFunc[any], opts ...TestOption) *MapSchema[K, V] {
	t := p.NewTestFunc("", p.BoolTFunc[any](testFunc), opts...)
	v.Test(Test[any](*t))
	return v
}

// Minimum number of entries
func (v *MapSchema[K, V]) Min(n int, options ...TestOption) *MapSchema[K, V] {
	t, fn := mapMin(n)
	return v.addTest(&t, fn, options...)
}

// Maximum number of entries
func (v *MapSchema[K, V]) Max(n int, options ...TestOption) *MapSchema[K, V] {
	t, fn := mapMax(n)
	return v.addTest(&t, fn, options...)
}

// Exact number of entries
func (v *MapSchema[K, V]) Len(n int, options ...TestOption) *MapSchema[K, V] {
	t, fn := mapLength(n)
	return v.addTest(&t, fn, options...)
}

func mapMin(n int) (p.Test[any], p.BoolTFunc[any]) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Map {
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

func mapMax(n int) (p.Test[any], p.BoolTFunc[any]) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Map {
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

func mapLength(n int) (p.Test[any], p.BoolTFunc[any]) {
	fn := func(val any, ctx Ctx) bool {
		rv := reflect.ValueOf(val).Elem()
		if rv.Kind() != reflect.Map {
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

func (v *MapSchema[K, V]) addTest(t *p.Test[any], fn p.BoolTFunc[any], options ...TestOption) *MapSchema[K, V] {
	p.TestFuncFromBool(fn, t)

	for _, opt := range options {
		opt(t)
	}

	v.processors = append(v.processors, t)
	return v
}
