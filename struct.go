package zog

import (
	"fmt"
	"reflect"
	"unicode"

	p "github.com/Oudwins/zog/primitives"
)

type Schema map[string]Processor

func Struct(schema Schema) *structProcessor {
	return &structProcessor{
		schema: schema,
	}
}

type structProcessor struct {
	preTransforms  []p.PreTransform
	schema         Schema
	postTransforms []p.PostTransform
	defaultVal     *string
	required       *p.Test
	catch          *string
}

// only supports val = map[string]any & dest = &struct
func (v *structProcessor) Parse(val any, dest any) p.ZogSchemaErrors {
	// create context
	// handle options
	// empty path
	// TODO create context -> but for single field
	var ctx = p.NewParseCtx()
	errs := p.NewErrsMap()
	path := p.Pather("")

	v.process(val, dest, errs, path, ctx)

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
	// 2. cast data as map[string]any
	m := val.(map[string]any)

	// for each field in the struct we process it
	for fieldName, fieldProcessor := range v.schema {
		publicFieldName := string(unicode.ToUpper(rune(fieldName[0]))) + fieldName[1:]
		// TODO HERE I NEED TO CHECK IF FIELD IS NIL & CREATE IT IF IT DOESN'T EXIST
		destPtr, err := getStructFieldPointerByName(dest, publicFieldName)
		if err != nil {
			panic(err)
		}
		fieldVal := m[fieldName]
		p := path.Push(fieldName)
		fieldProcessor.process(fieldVal, destPtr, errs, p, ctx)
	}
	// TODO custom tests

	// 4. postTransforms
	if v.postTransforms != nil {
		for _, fn := range v.postTransforms {
			err := fn(dest, ctx)
			if err != nil {
				errs.Add(path, Errors.WrapUnknown(err))
				return
			}
		}
	}
}

func getStructFieldPointerByName(v any, name string) (any, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected dest to be a pointer, got %s", val.Kind())
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected dest to point to a struct, got %s", val.Kind())
	}

	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("field %s not found", name)
	}
	// TODO HERE WE PROBABLY NEED TO CREATE structs, maps, slices, etc if they don't exist and that is the kind of field
	// if its an optional value
	fmt.Println("Field vals")
	fmt.Println(fieldVal.Kind(), fieldVal.Type())

	// HERE WE NEED TO CHECK IF THE VALUE IS A POINTER. If it is, we should check the underlying value for nil
	// if its nil we need to create it based on the type

	return fieldVal.Addr().Interface(), nil
}
