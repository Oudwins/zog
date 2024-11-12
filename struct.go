package zog

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

var _ ZogSchema = &structProcessor{}

type structProcessor struct {
	preTransforms  []p.PreTransform
	schema         Schema
	postTransforms []p.PostTransform
	tests          []p.Test
	// defaultVal     any
	required *p.Test
	// catch          any
}

// Returns the type of the schema
func (v *structProcessor) getType() zconst.ZogType {
	return zconst.TypeStruct
}

// Sets the coercer for the schema
func (v *structProcessor) setCoercer(c conf.CoercerFunc) {
	// noop
}

func (v *structProcessor) process(data any, dest any, path p.PathBuilder, ctx ParseCtx) {
	destType := zconst.TypeStruct
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(data, ctx)
			// bail if error in preTransform
			if err != nil {
				ctx.NewError(path, Errors.WrapUnknown(data, destType, err))
				return
			}
			data = nVal
		}
	}

	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if !ctx.HasErrored() {
			for _, fn := range v.postTransforms {
				err := fn(dest, ctx)
				if err != nil {
					ctx.NewError(path, Errors.WrapUnknown(data, destType, err))
					return
				}
			}
		}
	}()

	// 2. cast data as DataProvider
	_, isEmpty := data.(*p.EmptyDataProvider)
	if isEmpty {
		if v.required != nil {
			ctx.NewError(path, Errors.FromTest(data, destType, v.required, ctx))
			return
		}
		return
	}
	dataProv, err := p.TryNewAnyDataProvider(data)

	// 2.5 check for required & errors
	if err != nil {
		code := err.Code()
		// This means its optional and we got an error coercing the value to a DataProvider, so we can ignore it
		if v.required == nil && code == zconst.ErrCodeCoerce {
			return
		}
		// This means that its required but we got an error coercing the value or a factory errored with required
		if v.required != nil && (code == zconst.ErrCodeCoerce || code == zconst.ErrCodeRequired) {
			ctx.NewError(path, Errors.FromTest(data, destType, v.required, ctx))
			return
		}
		// Some other error happened. Coercion error
		ctx.NewError(path, err.SDType(destType).SValue(data))
		return
	}

	// 3. Process / validate struct fields
	structVal := reflect.ValueOf(dest).Elem()
	//

	for key, processor := range v.schema {
		fieldKey := key
		key = strings.ToUpper(string(key[0])) + key[1:]

		fieldMeta, ok := structVal.Type().FieldByName(key)
		if !ok {
			panic(fmt.Sprintf("Struct is missing expected schema key: %s", key))
		}
		destPtr := structVal.FieldByName(key).Addr().Interface()

		fieldTag, ok := fieldMeta.Tag.Lookup(zconst.ZogTag)
		if ok {
			fieldKey = fieldTag
		}

		switch schema := processor.(type) {
		case *structProcessor:
			schema.process(dataProv.GetNestedProvider(fieldKey), destPtr, path.Push(fieldKey), ctx)

		default:
			schema.process(dataProv.Get(fieldKey), destPtr, path.Push(fieldKey), ctx)
		}

	}

	// 3. Tests for struct
	for _, test := range v.tests {
		if !test.ValidateFunc(dest, ctx) {
			ctx.NewError(path, Errors.FromTest(data, destType, &test, ctx))
		}
	}

}

// ! USER FACING FUNCTIONS

// A map of field names to zog schemas
type Schema map[string]ZogSchema

// Returns a new structProcessor which can be used to parse input data into a struct
func Struct(schema Schema) *structProcessor {
	return &structProcessor{
		schema: schema,
	}
}

// Parses val into destPtr and validates each field based on the schema. Only supports val = map[string]any & dest = &struct
func (v *structProcessor) Parse(data any, destPtr any, options ...ParsingOption) p.ZogErrMap {
	errs := p.NewErrsMap()
	ctx := p.NewParseCtx(errs, conf.ErrorFormatter)
	for _, opt := range options {
		opt(ctx)
	}
	path := p.PathBuilder("")

	v.process(data, destPtr, path, ctx)

	return errs.M
}

// Add a pretransform step to the schema
func (v *structProcessor) PreTransform(transform p.PreTransform) *structProcessor {
	if v.preTransforms == nil {
		v.preTransforms = []p.PreTransform{}
	}
	v.preTransforms = append(v.preTransforms, transform)
	return v
}

// Adds posttransform function to schema
func (v *structProcessor) PostTransform(transform p.PostTransform) *structProcessor {
	if v.postTransforms == nil {
		v.postTransforms = []p.PostTransform{}
	}
	v.postTransforms = append(v.postTransforms, transform)
	return v
}

// ! MODIFIERS

// marks field as required
func (v *structProcessor) Required(options ...TestOption) *structProcessor {
	r := p.Required()
	for _, opt := range options {
		opt(&r)
	}
	v.required = &r
	return v
}

// marks field as optional
func (v *structProcessor) Optional() *structProcessor {
	v.required = nil
	return v
}

// // sets the default value
// func (v *structProcessor) Default(val any) *structProcessor {
// 	v.defaultVal = val
// 	return v
// }

// // sets the catch value (i.e the value to use if the validation fails)
// func (v *structProcessor) Catch(val any) *structProcessor {
// 	v.catch = val
// 	return v
// }

// ! VALIDATORS
// custom test function call it -> schema.Test(t z.Test, opts ...TestOption)
func (v *structProcessor) Test(t p.Test, opts ...TestOption) *structProcessor {
	for _, opt := range opts {
		opt(&t)
	}
	v.tests = append(v.tests, t)
	return v
}
