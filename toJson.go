package zog

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"

	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type JsonProcessor struct {
	Type string // "transform", "validator", "required"

	// Validator
	IssueCode *string
	IssuePath *string
	Params    map[string]any

	// transform
	TransformId *string
}

type JsonZogSchema struct {
	Type         string // "string"
	Processors   []JsonProcessor
	Child        *JsonZogSchema
	Required     *JsonProcessor
	DefaultValue any
	CatchValue   any
}

func (s *StringSchema[T]) toJson() *JsonZogSchema {
	rv := reflect.ValueOf(s).Elem()
	j := JsonZogSchema{
		Type: zconst.TypeString,
		// Required:     processorToJson(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToJson(rv),
	}
	return &j
}

type JsonifyableSchema interface {
	toJson() *JsonZogSchema
}

func ToJson(s JsonifyableSchema) ([]byte, error) {
	j := s.toJson()
	jsonSchema, err := json.Marshal(j)
	fmt.Println(string(jsonSchema))
	return jsonSchema, err
}

func processorToJson(p any) *JsonProcessor {
	if p == nil {
		return nil
	}

	rv := reflect.ValueOf(p)
	return processRVtoJson(rv)
}

func processRVtoJson(rv reflect.Value) *JsonProcessor {
	// for rv.Kind() == reflect.Pointer {
	// 	if rv.IsNil() {
	// 		return nil
	// 	}
	// 	rv = rv.Elem()
	// }

	j := &JsonProcessor{}

	// THE ISSUE I'M HAVING IS THAT ITS ACTUALLY A ZProcessor. THATS THE TYPE. So cannot tell between transformer and Test :(
	/*
		Above is solved, I can do one of:
		1. Make interface for Transfomer (actually not even needed because I already have one for the Test and we only have transformers and tests)
		2. rv.FieldByName("IssueCode").IsValid() vs rv.FieldByName("Transform").IsValid()

	*/
	// Actually the issue is that the processor value was obtained by accessing a private field. Which causes value.Interface() to panic
	// Maybe I just need to make ZProcessors property available in the schema? :(

	if isTest(rv) {
		fmt.Println("Its a test")
		j.Type = zconst.ProcessorTest
		c := rv.FieldByName("IssueCode").Interface().(string)
		j.IssueCode = &c

		path := rv.FieldByName("IssuePath").Interface().(string)
		j.IssuePath = &path
		params := rv.FieldByName("Params").Interface().(map[string]any)
		newParams := map[string]any{}
		maps.Copy(newParams, params)
		j.Params = newParams
	} else if isTransformer(rv) {
		fmt.Println("Its a transform")
		j.Type = zconst.ProcessorTransform
	} else {
		fmt.Println("This should neverh happen")
		return nil
	}
	return j
}

func processorsToJson(v reflect.Value) []JsonProcessor {
	l := v.FieldByName("processors")
	if l.IsNil() {
		return nil
	}
	ln := l.Len()
	out := []JsonProcessor{}
	for i := 0; i < ln; i++ {
		p := l.Index(i)
		fmt.Println(l.CanInterface())
		result := processRVtoJson(p)
		if result == nil {
			continue
		}
		out = append(out, *result)
	}
	return out
}

func deepCopyPrimitivePtr(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil
	}
	e := rv.Elem()

	ptr := reflect.New(e.Type())

	ptr.Elem().Set(e)
	return ptr
}

func isTransformer(rv reflect.Value) bool {
	_, ok := rv.Interface().(internals.TransfomerProcessor)
	return ok
}
func isTest(rv reflect.Value) bool {
	// return !isTransformer(rv)
	fmt.Println(rv.CanInterface())
	_, ok := rv.Interface().(internals.TestInterface)
	return ok
}
