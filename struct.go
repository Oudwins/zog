package zog

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ComplexZogSchema = &StructSchema{}

type StructSchema struct {
	preTransforms  []p.PreTransform
	schema         Schema
	postTransforms []p.PostTransform
	tests          []p.Test
	// defaultVal     any
	required *p.Test
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
type Schema map[string]ZogSchema

// Returns a new StructSchema which can be used to parse input data into a struct
func Struct(schema Schema) *StructSchema {
	return &StructSchema{
		schema: schema,
	}
}

// Parses val into destPtr and validates each field based on the schema. Only supports val = map[string]any & dest = &struct
func (v *StructSchema) Parse(data any, destPtr any, options ...ExecOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewExecCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}

	v.process(ctx.NewSchemaCtx(data, destPtr, p.PathBuilder(""), v.getType()))

	return errs.M
}

func (v *StructSchema) process(ctx *p.SchemaCtx) {
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(ctx.Val, ctx)
			// bail if error in preTransform
			if err != nil {
				ctx.AddIssue(ctx.Issue().SetError(err))
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
					ctx.AddIssue(ctx.Issue().SetError(err))
					return
				}
			}
		}
	}()

	// 2. cast data as DataProvider
	_, isEmpty := ctx.Val.(*p.EmptyDataProvider)
	if isEmpty {
		if v.required != nil {
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
			return
		}
		return
	}
	dataProv, err := p.TryNewAnyDataProvider(ctx.Val)

	// 2.5 check for required & errors
	if err != nil {
		code := err.Code()
		// This means its optional and we got an error coercing the value to a DataProvider, so we can ignore it
		if v.required == nil && code == zconst.ErrCodeCoerce {
			return
		}

		// Some other error happened. Coercion error. We don't know the cause
		if v.required == nil {
			ctx.AddIssue(ctx.Issue())
		}

		// This means that its required but we got an error coercing the value or a factory errored with required
		if code == zconst.ErrCodeCoerce || code == zconst.ErrCodeRequired {
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
		}
	}

	// 3. Process / validate struct fields
	structVal := reflect.ValueOf(ctx.DestPtr).Elem()
	//

	for key, processor := range v.schema {
		fieldKey := key
		key = strings.ToUpper(string(key[0])) + key[1:]

		fieldMeta, ok := structVal.Type().FieldByName(key)
		if !ok {
			panic(fmt.Sprintf("Struct is missing expected schema key: %s\n see the zog FAQ for more info", key))
		}
		destPtr := structVal.FieldByName(key).Addr().Interface()

		fieldTag, ok := fieldMeta.Tag.Lookup(zconst.ZogTag)
		if ok {
			fieldKey = fieldTag
		}

		switch schema := processor.(type) {
		case *StructSchema:
			schema.process(ctx.NewSchemaCtx(dataProv.GetNestedProvider(fieldKey), destPtr, ctx.Path.Push(fieldKey), schema.getType()))
		default:
			schema.process(ctx.NewSchemaCtx(dataProv.Get(fieldKey), destPtr, ctx.Path.Push(fieldKey), schema.getType()))
		}

	}

	// 3. Tests for struct
	for _, test := range v.tests {
		if !test.ValidateFunc(ctx.DestPtr, ctx) {
			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.DestPtr))
		}
	}

}

// Validate a struct pointer given the struct schema. Usage:
// userSchema.Validate(&User, ...options)
func (v *StructSchema) Validate(dataPtr any, options ...ExecOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewExecCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}

	v.validate(ctx.NewValidateSchemaCtx(dataPtr, p.PathBuilder(""), v.getType()))

	return errs.M
}

// Internal function to validate the data
func (v *StructSchema) validate(ctx *p.SchemaCtx) {
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
	refVal := reflect.ValueOf(ctx.Val).Elem()
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
	x := refVal.Interface()
	isZeroVal := p.IsZeroValue(x)

	if isZeroVal {
		if v.required == nil {
			return
		} else {
			// REQUIRED & ZERO VALUE
			ctx.AddIssue(ctx.IssueFromTest(v.required, ctx.Val))
			return
		}
	}

	// 3.1 tests for struct fields
	for key, schema := range v.schema {
		fieldKey := key
		key = strings.ToUpper(string(key[0])) + key[1:]

		fieldMeta, ok := refVal.Type().FieldByName(key)
		if !ok {
			panic(fmt.Sprintf("Struct is missing expected schema key: %s", key))
		}
		destPtr := refVal.FieldByName(key).Addr().Interface()

		fieldTag, ok := fieldMeta.Tag.Lookup(zconst.ZogTag)
		if ok {
			fieldKey = fieldTag
		}
		schema.validate(ctx.NewValidateSchemaCtx(destPtr, ctx.Path.Push(fieldKey), schema.getType()))
	}

	// 3. tests for slice
	for _, test := range v.tests {
		if !test.ValidateFunc(ctx.Val, ctx) {
			ctx.AddIssue(ctx.IssueFromTest(&test, ctx.Val))
		}
	}
	// 4. postTransforms -> defered see above
}

// Add a pretransform step to the schema
func (v *StructSchema) PreTransform(transform p.PreTransform) *StructSchema {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *StructSchema) PostTransform(transform p.PostTransform) *StructSchema {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *StructSchema) Required(options ...TestOption) *StructSchema {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *StructSchema) Optional() *StructSchema {
	v.required = nil
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
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *StructSchema) Test(t p.Test, opts ...TestOption) *StructSchema {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}

// Create a custom test function for the schema. This is similar to Zod's `.refine()` method.
func (v *StructSchema) TestFunc(testFunc p.TestFunc, options ...TestOption) *StructSchema {
	test := TestFunc("", testFunc)
	v.Test(test, options...)
	return v
}
