package zog

import (
	"fmt"
	"reflect"
	"unicode"

	p "github.com/Oudwins/zog/primitives"
)

type StructParser interface {
	Parse(val any, destPtr any) p.ZogSchemaErrors
}

// A map of field names to zog schemas
type Schema map[string]Processor

// Returns a new structProcessor which can be used to parse input data into a struct
func Struct(schema Schema) *structProcessor {
	return &structProcessor{
		schema: schema,
	}
}

type structProcessor struct {
	preTransforms  []p.PreTransform
	schema         Schema
	postTransforms []p.PostTransform
	tests          []p.Test
	// defaultVal     any
	required *p.Test
	// catch          any
}

// Parses val into destPtr and validates each field based on the schema. Only supports val = map[string]any & dest = &struct
func (v *structProcessor) Parse(val any, destPtr any) p.ZogSchemaErrors {
	var ctx = p.NewParseCtx()
	errs := p.NewErrsMap()
	path := p.Pather("")

	v.process(val, destPtr, errs, path, ctx)

	if errs.IsEmpty() {
		return nil
	}
	return errs.M
}

func (v *structProcessor) process(val any, dest any, errs p.ZogErrors, path p.Pather, ctx *p.ParseCtx) {
	// 1. preTransforms
	if v.preTransforms != nil {
		for _, fn := range v.preTransforms {
			nVal, err := fn(val, ctx)
			// bail if error in preTransform
			if err != nil {
				errs.Add(path, Errors.WrapUnknown(err))
				return
			}
			val = nVal
		}
	}

	// 4. postTransforms
	defer func() {
		// only run posttransforms on success
		if errs.IsEmpty() {
			for _, fn := range v.postTransforms {
				err := fn(dest, ctx)
				if err != nil {
					errs.Add(path, Errors.WrapUnknown(err))
					return
				}
			}
		}
	}()

	isZeroVal := p.IsZeroValue(val)

	if isZeroVal && v.required == nil {
		return
	}

	// 2. cast data as map[string]any
	m, ok := val.(map[string]any)
	if !ok {
		errs.Add(path, Errors.Wrap(fmt.Errorf("expected map[string]any at path %s", path), "failed to validate field"))
		return
	}

	// required
	if v.required != nil && isZeroVal {
		errs.Add(path, Errors.New(v.required.ErrorFunc(dest, ctx)))
		return
	}

	// 3. Process / validate struct fields
	structVal := reflect.ValueOf(dest).Elem()
	for i := 0; i < structVal.NumField(); i++ {
		fieldMeta := structVal.Type().Field(i)

		// skip private fields
		if !fieldMeta.IsExported() {
			continue
		}

		fieldKey := string(unicode.ToLower(rune(fieldMeta.Name[0]))) + fieldMeta.Name[1:]
		processor, ok := v.schema[fieldKey]
		if !ok {
			continue
		}
		fieldTag, ok := fieldMeta.Tag.Lookup(p.ZogTag)
		if ok {
			fieldKey = fieldTag
		}
		input := m[fieldKey]
		destPtr := structVal.Field(i).Addr().Interface()
		processor.process(input, destPtr, errs, path.Push(fieldKey), ctx)
	}

	// 3. Tests for struct
	for _, test := range v.tests {
		if !test.ValidateFunc(dest, ctx) {
			errs.Add(path, Errors.New(test.ErrorFunc(dest, ctx)))
		}
	}

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
	r := p.Required(p.DErrorFunc("is a required field"))
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
// custom test function call it -> schema.Test("test_name", z.Message(""), func(val any, ctx *p.ParseCtx) bool {return true})
func (v *structProcessor) Test(ruleName string, errorMsg TestOption, validateFunc p.TestFunc) *structProcessor {
	t := p.Test{
		Name:         ruleName,
		ErrorFunc:    nil,
		ValidateFunc: validateFunc,
	}
	errorMsg(&t)
	v.tests = append(v.tests, t)

	return v
}
