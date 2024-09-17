package zog

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type StructParser interface {
	Parse(val p.DataProvider, destPtr any) p.ZogErrMap
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

func (v *structProcessor) Merge(other *structProcessor) *structProcessor {
	new := &structProcessor{
		preTransforms:  make([]p.PreTransform, len(v.preTransforms)+len(other.preTransforms)),
		postTransforms: make([]p.PostTransform, len(v.postTransforms)+len(other.postTransforms)),
		tests:          make([]p.Test, len(v.tests)+len(other.tests)),
	}
	if v.preTransforms != nil {
		new.preTransforms = append(new.preTransforms, v.preTransforms...)
	}
	if other.preTransforms != nil {
		new.preTransforms = append(new.preTransforms, other.preTransforms...)
	}

	if v.postTransforms != nil {
		new.postTransforms = append(new.postTransforms, v.postTransforms...)
	}
	if other.postTransforms != nil {
		new.postTransforms = append(new.postTransforms, other.postTransforms...)
	}

	if v.tests != nil {
		new.tests = append(new.tests, v.tests...)
	}
	if other.tests != nil {
		new.tests = append(new.tests, other.tests...)
	}
	new.required = v.required
	new.schema = Schema{}
	maps.Copy(new.schema, v.schema)
	maps.Copy(new.schema, other.schema)
	return new
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

	_, isEmptyDP := data.(*p.EmptyDataProvider)

	if isEmptyDP || p.IsZeroValue(data) {
		if v.required == nil {
			return
		} else {
			ctx.NewError(path, Errors.FromTest(data, destType, v.required, ctx))
			return
		}
	}

	// 2. cast data as DataProvider
	dataProv, ok := data.(p.DataProvider)
	if !ok {
		if dataProv, ok = p.TryNewAnyDataProvider(data); !ok {
			ctx.NewError(path, Errors.New(zconst.ErrCodeCoerce, data, destType, nil, "", errors.New("could not convert data to a data provider")))
			return
		}
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
