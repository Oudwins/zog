package zog

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"

	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type ExMetaRegistry map[any]map[string]any

func registryAdd(r ExMetaRegistry, key any, path string, value any) {
	if _, ok := r[key]; !ok {
		r[key] = map[string]any{}
	}
	r[key][path] = value
}

func getGenericTypeName[T any]() string {
	var zero T
	t := reflect.TypeOf(zero)

	// If T is a pointer or interface, handle nil case
	if t == nil {
		t = reflect.TypeOf((*T)(nil)).Elem()
	}

	return t.Name()
}

func (r ExMetaRegistry) Add(key any, path string, value any) {
	if _, ok := r[key]; !ok {
		r[key] = map[string]any{}
	}

}

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
var EX_META_REGISTRY = map[any]map[string]any{}

// TODO make zog schemas for all of these to validate them!
type JsonProcessor struct {
	Type string // "transform", "validator", "required"

	// Validator
	ID        string // issue code or transform ID
	IssuePath *string
	Message   *string
	Params    map[string]any
}

type JsonTest struct {
	Type string // "test"
	// Validator
	IssueCode *string
	IssuePath *string
	Params    map[string]any
}

type JsonTransformer struct {
	Type        string // "transformer"
	TransformId *string
}

type JsonZogSchema struct {
	Type         string // "string"
	Processors   []any  // JsonTest or JsonTransformer
	Format       *string
	Child        any // *JsonZogSchema | map[string]JsonZogSchema
	Required     *JsonTest
	DefaultValue any
	CatchValue   any
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

func (s *StringSchema[T]) toJson() *JsonZogSchema {
	rvP := reflect.ValueOf(s.processors)
	j := JsonZogSchema{
		Type:         zconst.TypeString,
		Required:     toJsonTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToJson(rvP),
	}
	return &j
}

func (s *NumberSchema[T]) toJson() *JsonZogSchema {
	rvP := reflect.ValueOf(s.processors)
	j := JsonZogSchema{
		Type:         zconst.TypeNumber,
		Required:     toJsonTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToJson(rvP),
	}
	return &j
}

func (s *BoolSchema[T]) toJson() *JsonZogSchema {
	rvP := reflect.ValueOf(s.processors)
	j := JsonZogSchema{
		Type:         zconst.TypeBool,
		Required:     toJsonTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToJson(rvP),
	}
	return &j
}

func (s *TimeSchema) toJson() *JsonZogSchema {
	// rvP := reflect.ValueOf(s.processors)
	j := JsonZogSchema{
		Type: zconst.TypeTime,
		// Required:     toJsonTest(s.required),
		// DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		// CatchValue:   deepCopyPrimitivePtr(s.catch),
		// Processors:   processorsToJson(rvP),
	}
	exmeta, ok := EX_META_REGISTRY[s]
	if ok {
		x := exmeta["format"].(string)
		j.Format = &x
	}
	return &j
}

func (s *PointerSchema) toJson() *JsonZogSchema {
	j := JsonZogSchema{
		Type:     zconst.TypePtr,
		Required: toJsonTest(s.required),
		Child:    s.schema.toJson(),
	}
	return &j
}

func (s *SliceSchema) toJson() *JsonZogSchema {
	rvP := reflect.ValueOf(s.processors)
	j := JsonZogSchema{
		Type:         zconst.TypeSlice,
		Required:     toJsonTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		Processors:   processorsToJson(rvP),
		Child:        s.schema.toJson(),
	}
	return &j
}

func (s *StructSchema) toJson() *JsonZogSchema {
	rvP := reflect.ValueOf(s.processors)
	j := JsonZogSchema{
		Type:       zconst.TypeSlice,
		Required:   toJsonTest(s.required),
		Processors: processorsToJson(rvP),
		Child:      toJsonShape(s.schema),
	}
	return &j
}

func (s *Custom[T]) toJson() *JsonZogSchema {
	j := JsonZogSchema{
		Type: "custom",
		// TODO not sure this is the right place for this info
		Required: toJsonTest(&s.test),
	}
	return &j
}

func toJsonShape(s Shape) (m map[string]JsonZogSchema) {
	// iterate and return
	// TODO forgot how to fucking do this
	return m
}

func (s *PreprocessSchema[F, T]) toJson() *JsonZogSchema {
	j := JsonZogSchema{
		Type:  "preprocess",
		Child: s.schema.toJson(),
	}
	return &j
}

func (s *BoxedSchema[B, T]) toJson() *JsonZogSchema {
	j := JsonZogSchema{
		Type:  "boxed",
		Child: s.schema.toJson(),
	}
	return &j
}

func processRVtoJson(rv reflect.Value) any {

	if !rv.CanInterface() {
		// TODO add assert here
		fmt.Println("THIS SHOULD NEVER HAPPEN")
		return nil
	}

	rvi := rv.Interface()

	var out any

	if test, ok := rvi.(internals.TestInterface); ok {
		out = toJsonTest(test)
	} else if trans, ok := rvi.(internals.TransformerInterface); ok {
		out = toJsonTransformer(trans)
	} else {
		// TODO add assert here
		fmt.Println("THIS SHOULD NEVER HAPPEN")
	}
	return out
}

func processorsToJson(l reflect.Value) []any {
	if l.IsNil() {
		return nil
	}
	ln := l.Len()
	out := []any{}
	for i := 0; i < ln; i++ {
		p := l.Index(i)
		fmt.Println(l.CanInterface())
		result := processRVtoJson(p)
		if result == nil {
			continue
		}
		out = append(out, result)
	}
	return out
}

func toJsonTest(test internals.TestInterface) *JsonTest {
	if test == nil {
		return nil
	}

	j := JsonTest{}
	j.Type = zconst.ProcessorTest
	c := test.GetIssueCode()
	j.IssueCode = &c
	path := test.GetIssuePath()
	j.IssuePath = &path
	params := test.GetParams()
	newParams := map[string]any{}
	maps.Copy(newParams, params)
	j.Params = newParams
	return &j
}

func toJsonTransformer(transformer internals.TransformerInterface) *JsonTransformer {
	// TODO issue here is that I can't get the code for the transformer and we currently do not have IDs so no way to actually know what this will be
	if transformer == nil {
		return nil
	}
	j := JsonTransformer{}
	j.Type = zconst.ProcessorTransform
	return &j
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
