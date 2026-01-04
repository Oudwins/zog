package zog

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"

	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
	"github.com/Oudwins/zog/zss"
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

type ZSSSerializable interface {
	toZSS() *zss.ZSSSchema
}

func EXPERIMENTAL_TO_ZSS(s ZSSSerializable) ([]byte, error) {
	j := s.toZSS()
	jsonSchema, err := json.Marshal(j)
	fmt.Println(string(jsonSchema))
	return jsonSchema, err
}

func (s *StringSchema[T]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		GoType:       zconst.TypeString,
		Required:     toZSSTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToZSS(rvP),
	}
	return &j
}

func (s *NumberSchema[T]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		GoType:       zconst.TypeNumber,
		Required:     toZSSTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToZSS(rvP),
	}
	return &j
}

func (s *BoolSchema[T]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		GoType:       zconst.TypeBool,
		Required:     toZSSTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		CatchValue:   deepCopyPrimitivePtr(s.catch),
		Processors:   processorsToZSS(rvP),
	}
	return &j
}

func (s *TimeSchema) toZSS() *zss.ZSSSchema {
	// rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		GoType: zconst.TypeTime,
		// Required:     toZSSTest(s.required),
		// DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		// CatchValue:   deepCopyPrimitivePtr(s.catch),
		// Processors:   processorsToZSS(rvP),
	}
	exmeta, ok := EX_META_REGISTRY[s]
	if ok {
		x := exmeta["format"].(string)
		j.Format = &x
	}
	return &j
}

func (s *PointerSchema) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		GoType:   zconst.TypePtr,
		Required: toZSSTest(s.required),
		Child:    s.schema.toZSS(),
	}
	return &j
}

func (s *SliceSchema) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		GoType:       zconst.TypeSlice,
		Required:     toZSSTest(s.required),
		DefaultValue: deepCopyPrimitivePtr(s.defaultVal),
		Processors:   processorsToZSS(rvP),
		Child:        s.schema.toZSS(),
	}
	return &j
}

func (s *StructSchema) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		GoType:     zconst.TypeSlice,
		Required:   toZSSTest(s.required),
		Processors: processorsToZSS(rvP),
		Child:      toZSSShape(s.schema),
	}
	return &j
}

func (s *Custom[T]) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		GoType: "custom",
		// TODO not sure this is the right place for this info
		Required: toZSSTest(&s.test),
	}
	return &j
}

func toZSSShape(s Shape) (m map[string]zss.ZSSSchema) {
	// iterate and return
	// TODO forgot how to fucking do this
	return m
}

func (s *PreprocessSchema[F, T]) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		GoType: "preprocess",
		Child:  s.schema.toZSS(),
	}
	return &j
}

func (s *BoxedSchema[B, T]) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		GoType: "boxed",
		Child:  s.schema.toZSS(),
	}
	return &j
}

func processRVtoZSS(rv reflect.Value) *zss.ZSSProcessor {

	if !rv.CanInterface() {
		// TODO add assert here
		fmt.Println("THIS SHOULD NEVER HAPPEN")
		return nil
	}

	rvi := rv.Interface()

	out := zss.ZSSProcessor{}

	if test, ok := rvi.(internals.TestInterface); ok {
		out.Test = toZSSTest(test)
		out.Kind = zconst.ZogProcessorTest
	} else if trans, ok := rvi.(internals.TransformerInterface); ok {
		out.Transformer = toZSSTransformer(trans)
		out.Kind = zconst.ZogProcessorTransform
	} else {
		// TODO add assert here
		fmt.Println("THIS SHOULD NEVER HAPPEN")
	}
	return &out
}

func processorsToZSS(l reflect.Value) []zss.ZSSProcessor {
	if l.IsNil() {
		return nil
	}
	ln := l.Len()
	out := []zss.ZSSProcessor{}
	for i := 0; i < ln; i++ {
		p := l.Index(i)
		fmt.Println(l.CanInterface())
		result := processRVtoZSS(p)
		if result == nil {
			continue
		}
		out = append(out, *result)
	}
	return out
}

func toZSSRequired(test internals.TestInterface) *zss.ZSSTest {
	if test == nil {
		return nil
	}

	j := toZSSTest(test)
	(*j).ID = zconst.ZogProcessorRequired
	return j
}

func toZSSTest(test internals.TestInterface) *zss.ZSSTest {
	if test == nil {
		return nil
	}

	j := zss.ZSSTest{}
	c := test.GetIssueCode()
	j.ID = c
	path := test.GetIssuePath()
	j.IssuePath = &path
	params := test.GetParams()
	newParams := map[string]any{}
	maps.Copy(newParams, params)
	j.Params = newParams
	return &j
}

func toZSSTransformer(transformer internals.TransformerInterface) *zss.ZSSTransformer {
	// TODO issue here is that I can't get the code for the transformer and we currently do not have IDs so no way to actually know what this will be
	if transformer == nil {
		return nil
	}
	j := zss.ZSSTransformer{}
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
