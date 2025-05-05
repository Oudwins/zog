package zog

import (
	"fmt"
	"reflect"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ComplexZogSchema = &StructSchema{}

type StructSchema struct {
	schema     Shape
	processors []p.ZProcessor[any]
	// defaultVal     any
	required *p.Test[any]
	// catch          any
}

// Returns the type of the schema
func (v *StructSchema) getType() zconst.ZogType {
	return zconst.TypeStruct
}

// Sets the coercer for the schema
func (v *StructSchema) setCoercer(c conf.CoercerFunc) {
	// noop
}

// ! USER FACING FUNCTIONS

// A map of field names to zog schemas
type Shape map[string]ZogSchema

// Deprecated: use z.Struct(z.Shape{}) instead
// A map of field names to zog schemas
type Schema = Shape

// Returns a new StructSchema which can be used to parse input data into a struct
func Struct(schema Shape) *StructSchema {
	return &StructSchema{
		schema: schema,
	}
}

// Parses val into destPtr and validates each field based on the schema. Only supports val = map[string]any & dest = &struct
func (v *StructSchema) Parse(data any, destPtr any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(data, destPtr, path, v.getType())
	defer sctx.Free()
	v.process(sctx)

	return errs.M
}

func (v *StructSchema) process(ctx *p.SchemaCtx) {

	var dataProv p.DataProvider
	// 2. cast data as DataProvider
	if factory, ok := ctx.Data.(p.DpFactory); ok {
		newDp, err := factory()
		// This is a little bit hacky. But we want to exit here because the error came from zhttp. Meaning we had an error trying to parse the request.
		// I'm not sure if this is the best behaviour? Do we want to exit here or do we want to continue processing (ofc we add the error always)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromUnknownError(err))
			return
		}
		dataProv = newDp
	} else {
		newDp, err := p.TryNewAnyDataProvider(ctx.Data)
		if err != nil {
			ctx.AddIssue(ctx.IssueFromCoerce(err))
			return
		}
		dataProv = newDp
	}

	// 3. Process / validate struct fields
	structVal := reflect.ValueOf(ctx.ValPtr).Elem()
	subCtx := ctx.NewSchemaCtx(ctx.Data, ctx.ValPtr, ctx.Path, v.getType())
	defer subCtx.Free()
	for key, processor := range v.schema {
		originalKey := key
		if key[0] >= 'a' && key[0] <= 'z' {
			var b [32]byte // Use a size that fits your max key length
			copy(b[:], key)
			b[0] -= 32
			key = string(b[:len(key)])
		}

		fieldMeta, ok := structVal.Type().FieldByName(key)
		if !ok {
			p.Panicf(p.PanicMissingStructField, ctx.String(), key)
		}
		destPtr := structVal.FieldByName(key).Addr().Interface()

		subValue, fieldKey := dataProv.GetByField(fieldMeta, originalKey)
		subCtx.Data = subValue
		subCtx.ValPtr = destPtr
		subCtx.Path.Push(&fieldKey)
		subCtx.DType = processor.getType()
		subCtx.Exit = false
		processor.process(subCtx)
		subCtx.Path.Pop()
	}

	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(ctx.ValPtr, ctx)
		if ctx.Exit {
			// Catch here
			return
		}
	}

}

// Validate a struct pointer given the struct schema. Usage:
// userSchema.Validate(&User, ...options)
func (v *StructSchema) Validate(dataPtr any, options ...ExecOption) ZogIssueMap {
	errs := p.NewErrsMap()
	defer errs.Free()
	ctx := p.NewExecCtx(errs, conf.IssueFormatter)
	defer ctx.Free()
	for _, opt := range options {
		opt(ctx)
	}
	path := p.NewPathBuilder()
	defer path.Free()
	sctx := ctx.NewSchemaCtx(dataPtr, dataPtr, path, v.getType())
	defer sctx.Free()
	v.validate(sctx)

	return errs.M
}

// Internal function to validate the data
func (v *StructSchema) validate(ctx *p.SchemaCtx) {
	refVal := reflect.ValueOf(ctx.ValPtr).Elem()

	// 2. cast data to string & handle default/required

	// 3.1 tests for struct fields
	subCtx := ctx.NewSchemaCtx(ctx.Data, ctx.ValPtr, ctx.Path, v.getType())
	defer subCtx.Free()
	for key, schema := range v.schema {
		fieldKey := key
		if key[0] >= 'a' && key[0] <= 'z' {
			var b [32]byte // Use a size that fits your max key length
			copy(b[:], key)
			b[0] -= 32
			key = string(b[:len(key)])
		}

		fieldMeta, ok := refVal.Type().FieldByName(key)
		if !ok {
			panic(fmt.Sprintf("Struct is missing expected schema key: %s", key))
		}
		destPtr := refVal.FieldByName(key).Addr().Interface()

		fieldTag, ok := fieldMeta.Tag.Lookup(zconst.ZogTag)
		if ok {
			fieldKey = fieldTag
		}
		subCtx.Data = destPtr
		subCtx.ValPtr = destPtr
		subCtx.Path.Push(&fieldKey)
		subCtx.DType = schema.getType()
		schema.validate(subCtx)
		subCtx.Path.Pop()
	}

	for _, processor := range v.processors {
		ctx.Processor = processor
		processor.ZProcess(ctx.ValPtr, ctx)
		if ctx.Exit {
			// Catch here
			return
		}
	}
}

// Adds posttransform function to schema
func (v *StructSchema) Transform(transform p.Transform[any]) *StructSchema {
	v.processors = append(v.processors, &p.TransformProcessor[any]{Transform: transform})
	return v
}

// ! MODIFIERS

// Deprecated: structs are not required or optional. They pass through to the fields. If you want to say that an entire struct may not exist you should use z.Ptr(z.Struct(...))
// This now is a noop. But I believe most people expect it to work how it does now.
// marks field as required
func (v *StructSchema) Required(options ...TestOption) *StructSchema {
	return v
}

// Deprecated: structs are not required or optional. They pass through to the fields. If you want to say that an entire struct may not exist you should use z.Ptr(z.Struct(...))
// marks field as optional
func (v *StructSchema) Optional() *StructSchema {
	return v
}

// // sets the default value
// func (v *StructSchema) Default(val any) *StructSchema {
// 	v.defaultVal = val
// 	return v
// }

// // sets the catch value (i.e the value to use if the validation fails)
// func (v *StructSchema) Catch(val any) *StructSchema {
// 	v.catch = val
// 	return v
// }

// ! VALIDATORS
// custom test function call it -> schema.Test(t z.Test)
func (v *StructSchema) Test(t Test[any]) *StructSchema {
	x := p.Test[any](t)
	v.processors = append(v.processors, &x)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *StructSchema) TestFunc(testFunc BoolTFunc[any], options ...TestOption) *StructSchema {
	test := p.NewTestFunc("", p.BoolTFunc[any](testFunc), options...)
	v.Test(Test[any](*test))
	return v
}
