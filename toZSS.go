package zog

import (
	"maps"
	"reflect"

	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/pkgs/internals"
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
)

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
type ExMetaRegistry map[any]map[string]any

const (
	EX_META_KEY_FORMAT  = "format"
	EX_META_KEY_MESSAGE = "message"
	EX_META_KEY_ID      = "id"
)

func getZSSGoType[T any]() zss.ZSSGoType {
	var zero T
	t := reflect.TypeOf(zero)

	// If T is a pointer or interface, handle nil case
	if t == nil {
		t = reflect.TypeOf((*T)(nil)).Elem()
	}

	return zss.ZSSGoType{
		PkgPath: t.PkgPath(),
		Name:    t.Name(),
		Display: t.String(),
	}
}

type ZSSSerializable interface {
	toZSS() *zss.ZSSSchema
}

func EXPERIMENTAL_TO_ZSS(s ZSSSerializable) zss.ZSSDocument {
	j := s.toZSS()
	return zss.ZSSDocument{
		Version: zss.ZSS_VERSION_LATEST,
		Root:    j,
	}
}

func (s *StringSchema[T]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeString,
		Required:     toZSSRequired(s.required, zconst.TypeString),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeString),
	}

	if EXHAUSTIVE_METADATA {
		j.GoTypes = []zss.ZSSGoType{getZSSGoType[T]()}
	}

	return &j
}

func (s *NumberSchema[T]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeNumber,
		Required:     toZSSRequired(s.required, zconst.TypeNumber),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeNumber),
	}

	if EXHAUSTIVE_METADATA {
		j.GoTypes = []zss.ZSSGoType{getZSSGoType[T]()}
	}
	return &j
}

func (s *BoolSchema[T]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeBool,
		Required:     toZSSRequired(s.required, zconst.TypeBool),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeBool),
	}

	if EXHAUSTIVE_METADATA {
		j.GoTypes = []zss.ZSSGoType{getZSSGoType[T]()}
	}
	return &j
}

func (s *TimeSchema) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeTime,
		Required:     toZSSRequired(s.required, zconst.TypeTime),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeTime),
	}
	if x, ok := RegistryGet(exMetaRegistry, s, EX_META_KEY_FORMAT); ok {
		str := x.(string)
		j.Format = &str
	}
	return &j
}

func (s *PointerSchema) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind:     zconst.TypePtr,
		Required: toZSSRequired(s.required, s.schema.getType()),
		Childs:   []zss.ZSSSchemaChild{{Kind: zss.ZSSSchemaChildKindSchema, Schema: s.schema.toZSS()}},
	}
	return &j
}

func (s *SliceSchema) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeSlice,
		Required:     toZSSRequired(s.required, zconst.TypeSlice),
		DefaultValue: defaultValueFromAnyFunc(s.defaultFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeSlice),
		Childs:       []zss.ZSSSchemaChild{{Kind: zss.ZSSSchemaChildKindSchema, Schema: s.schema.toZSS()}},
	}
	return &j
}

func (s *MapSchema[K, V]) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	childMap := map[string]zss.ZSSSchema{
		"key":   *s.keySchema.toZSS(),
		"value": *s.valueSchema.toZSS(),
	}
	defaultValue := shallowCopyMapFromFunc(s.defaultFunc)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeMap,
		Required:     toZSSRequired(s.required, zconst.TypeMap),
		DefaultValue: defaultValue,
		Processors:   processorsToZSS(rvP, zconst.TypeMap),
		Childs:       []zss.ZSSSchemaChild{{Kind: zss.ZSSSchemaChildKindShape, Shape: childMap}},
	}
	return &j
}

// Helper function to shallow copy a map for ZSS
func shallowCopyMap[K comparable, V any](m map[K]V) any {
	if m == nil {
		return nil
	}
	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

func shallowCopyMapFromFunc[K comparable, V any](defaultFunc func() map[K]V) any {
	if defaultFunc == nil {
		return nil
	}
	return shallowCopyMap(defaultFunc())
}

func (s *StructSchema) toZSS() *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:       zconst.TypeStruct,
		Required:   toZSSRequired(s.required, zconst.TypeStruct),
		Processors: processorsToZSS(rvP, zconst.TypeStruct),
		Childs:     []zss.ZSSSchemaChild{{Kind: zss.ZSSSchemaChildKindShape, Shape: toZSSShape(s.schema)}},
	}
	return &j
}

func (s *Custom[T]) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind: zconst.TypeCustom,
		// TODO not sure this is the right place for this info
		// TODO THIS create this
		Processors: toZSSProcessorList(&s.test, zconst.TypeCustom),
	}
	if EXHAUSTIVE_METADATA {
		j.GoTypes = []zss.ZSSGoType{getZSSGoType[T]()}
	}
	return &j
}

func toZSSShape(m Shape) map[string]zss.ZSSSchema {
	out := map[string]zss.ZSSSchema{}
	for k, v := range m {
		out[k] = *v.toZSS()
	}
	return out
}

func (s *PreprocessSchema[F, T]) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind:   zconst.TypePreprocess,
		Childs: []zss.ZSSSchemaChild{{Kind: zss.ZSSSchemaChildKindSchema, Schema: s.schema.toZSS()}},
	}

	if EXHAUSTIVE_METADATA {
		j.GoTypes = []zss.ZSSGoType{getZSSGoType[F](), getZSSGoType[T]()}
	}
	return &j
}

func (s *BoxedSchema[B, T]) toZSS() *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind:   zconst.TypeBoxed,
		Childs: []zss.ZSSSchemaChild{{Kind: zss.ZSSSchemaChildKindSchema, Schema: s.schema.toZSS()}},
	}

	if EXHAUSTIVE_METADATA {
		j.GoTypes = []zss.ZSSGoType{getZSSGoType[B](), getZSSGoType[T]()}
	}
	return &j
}

func processRVtoZSS(rv reflect.Value, dtype zconst.ZogType) *zss.ZSSProcessor {

	if !rv.CanInterface() {
		// TODO better error messages + maybe not panic
		panic("[Zog] - This should never happen. processRVtoZSS: rv.CanInterface() is false")
	}

	rvi := rv.Interface()

	out := zss.ZSSProcessor{}

	if test, ok := rvi.(internals.TestInterface); ok {
		out.Test = toZSSTest(test, dtype)
		out.Kind = zconst.ZogProcessorTest
	} else if trans, ok := rvi.(internals.TransformerInterface); ok {
		out.Transformer = toZSSTransformer(trans)
		out.Kind = zconst.ZogProcessorTransform
	} else {
		// TODO better error messages + maybe not panic
		panic("[Zog] - This should never happen. processRVtoZSS: rvi is not a TestInterface or TransformerInterface")
	}
	return &out
}

func toZSSRequired(test any, dtype zconst.ZogType) *zss.ZSSTest {
	if test == nil {
		return nil
	}

	// Check if the underlying value is actually nil using reflection
	// This handles the case where a nil pointer is passed as an interface
	rv := reflect.ValueOf(test)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return nil
	}

	j := toZSSTest(test.(internals.TestInterface), dtype)
	(*j).ID = zconst.ZogProcessorRequired
	return j
}

func processorsToZSS(l reflect.Value, dtype zconst.ZogType) []zss.ZSSProcessor {
	if l.IsNil() {
		return nil
	}
	ln := l.Len()
	out := []zss.ZSSProcessor{}
	for i := 0; i < ln; i++ {
		p := l.Index(i)
		result := processRVtoZSS(p, dtype)
		if result == nil {
			continue
		}
		out = append(out, *result)
	}
	return out
}

func toZSSProcessorList(test any, dtype zconst.ZogType) []zss.ZSSProcessor {
	if test == nil {
		return nil
	}

	// Check if the underlying value is actually nil using reflection
	// This handles the case where a nil pointer is passed as an interface
	rv := reflect.ValueOf(test)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return nil
	}
	p := processRVtoZSS(rv, dtype)
	if p == nil {
		return nil
	}
	return []zss.ZSSProcessor{*p}

}

// fakeCtx is a minimal implementation of Ctx interface for extracting default messages
type fakeCtx struct{}

func (f *fakeCtx) Get(key string) any                                       { return nil }
func (f *fakeCtx) AddIssue(e *internals.ZogIssue)                           {}
func (f *fakeCtx) Issue() *internals.ZogIssue                               { return internals.NewZogIssue() }
func (f *fakeCtx) NewError(p *internals.PathBuilder, e *internals.ZogIssue) {}
func (f *fakeCtx) HasErrored() bool                                         { return false }

var fakeCtxInstance = &fakeCtx{}

func toZSSTest(test internals.TestInterface, dtype zconst.ZogType) *zss.ZSSTest {
	if test == nil {
		return nil
	}

	j := zss.ZSSTest{}
	c := test.GetIssueCode()
	j.ID = c
	path := test.GetIssuePath()
	if len(path) > 0 {
		j.IssuePath = path
	}
	params := test.GetParams()
	newParams := map[string]any{}
	maps.Copy(newParams, params)
	j.Params = newParams

	// Check for custom message in registry first
	if message, ok := RegistryGet(exMetaRegistry, test, EX_META_KEY_MESSAGE); ok {
		j.Message = message.(string)
	}

	// If no message is set, extract message using the test's formatter or default
	if j.Message == "" {
		fakeIssue := internals.NewZogIssue().
			SetCode(c).
			SetDType(dtype).
			SetParams(params)
		// Use the test's custom formatter if available, otherwise use default
		if customFmter := test.GetIssueFmtFunc(); customFmter != nil {
			customFmter(fakeIssue, fakeCtxInstance)
		} else {
			conf.DefaultIssueFormatter(fakeIssue, fakeCtxInstance)
		}
		j.Message = fakeIssue.Message
		internals.FreeIssue(fakeIssue)
	}

	return &j
}

func toZSSTransformer(transformer internals.TransformerInterface) *zss.ZSSTransformer {
	// TODO issue here is that I can't get the code for the transformer and we currently do not have IDs so no way to actually know what this will be
	if transformer == nil {
		return nil
	}
	j := zss.ZSSTransformer{
		ID: zconst.ZogTransformIDCustom,
	}

	// extra
	if id, ok := RegistryGet(exMetaRegistry, transformer, EX_META_KEY_ID); ok {
		j.ID = id.(zconst.ZogTransformID)
	}
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
	return ptr.Interface()
}

func defaultValueFromFunc[T any](defaultFunc func() T) any {
	if defaultFunc == nil {
		return nil
	}
	val := defaultFunc()
	return deepCopyPrimitivePtr(&val)
}

func defaultValueFromAnyFunc(defaultFunc func() any) any {
	if defaultFunc == nil {
		return nil
	}
	val := defaultFunc()
	return deepCopyPrimitivePtr(&val)
}
