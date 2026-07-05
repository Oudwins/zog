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
	t := typeOf[T]()

	return zss.ZSSGoType{
		PkgPath: t.PkgPath(),
		Name:    t.Name(),
		Display: t.String(),
	}
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

type ZSSSerializable interface {
	toZSS(*ZSSSerializeCtx) *zss.ZSSSchema
}

type ZSSSerializeCtx struct {
	defs        map[string]*zss.ZSSSchema
	defNames    map[ZSSSerializable]int
	nextDef     *int
	currentType reflect.Type
}

func newZSSSerializeCtx() *ZSSSerializeCtx {
	return &ZSSSerializeCtx{
		defs:     map[string]*zss.ZSSSchema{},
		defNames: map[ZSSSerializable]int{},
		nextDef:  new(int),
	}
}

func (ctx *ZSSSerializeCtx) withType(t reflect.Type) *ZSSSerializeCtx {
	copy := *ctx
	copy.currentType = t
	return &copy
}

func (ctx *ZSSSerializeCtx) refFor(schema ZSSSerializable, typ reflect.Type) *zss.ZSSSchema {
	if name, ok := ctx.defNames[schema]; ok {
		ref := zss.ZSSRefFromKey(name)
		return &zss.ZSSSchema{Ref: &ref}
	}

	(*ctx.nextDef)++
	name := *ctx.nextDef
	ctx.defNames[schema] = name
	ctx.defs[zss.ZSSDefKeyFromKey(name)] = schema.toZSS(ctx.withType(typ))
	ref := zss.ZSSRefFromKey(name)
	return &zss.ZSSSchema{Ref: &ref}
}

func EXPERIMENTAL_TO_ZSS[T any](s ZSSSerializable) zss.ZSSDocument {
	ctx := newZSSSerializeCtx()
	ctx.currentType = reflect.TypeOf((*T)(nil)).Elem()
	j := s.toZSS(ctx)
	return zss.ZSSDocument{
		URI:  zss.ZSS_VERSION_LATEST,
		Root: j,
		Defs: ctx.defs,
	}
}

func addCurrentGoType(schema *zss.ZSSSchema, ctx *ZSSSerializeCtx) {
	if ctx.currentType == nil {
		return
	}
	if isEmptyInterface(ctx.currentType) {
		return
	}
	schema.GoTypes = []zss.ZSSGoType{goTypeFromReflect(ctx.currentType)}
}

func goTypeFromReflect(t reflect.Type) zss.ZSSGoType {
	return zss.ZSSGoType{PkgPath: t.PkgPath(), Name: t.Name(), Display: t.String()}
}

func indirectType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// returns field, found
func fieldTypeForShapeKey(t reflect.Type, key string) (reflect.StructField, bool) {
	t = indirectType(t)
	if t == nil || t.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}
	if key == "" {
		return reflect.StructField{}, false
	}
	fieldName := key
	if key[0] >= 'a' && key[0] <= 'z' {
		b := []byte(key)
		b[0] -= 32
		fieldName = string(b)
	}
	return t.FieldByName(fieldName)
}

func elementType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array:
		return t.Elem()
	default:
		t = indirectType(t)
		if t != nil && (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) {
			return t.Elem()
		}
		return nil
	}
}

func mapKeyType(t reflect.Type) reflect.Type {
	t = indirectType(t)
	if t == nil || t.Kind() != reflect.Map {
		return nil
	}
	return t.Key()
}

func mapValueType(t reflect.Type) reflect.Type {
	t = indirectType(t)
	if t == nil || t.Kind() != reflect.Map {
		return nil
	}
	return t.Elem()
}

func isEmptyInterface(t reflect.Type) bool {
	return t != nil && t.Kind() == reflect.Interface && t.NumMethod() == 0
}

func derivedType[T any](ctx *ZSSSerializeCtx) reflect.Type {
	if isEmptyInterface(ctx.currentType) {
		return nil
	}
	return typeOf[T]()
}

func (s *StringSchema[T]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeString,
		Required:     toZSSRequired(s.required, zconst.TypeString),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeString),
	}

	addCurrentGoType(&j, ctx)

	return &j
}

func (s *NumberSchema[T]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeNumber,
		Required:     toZSSRequired(s.required, zconst.TypeNumber),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeNumber),
	}

	addCurrentGoType(&j, ctx)
	return &j
}

func (s *BoolSchema[T]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeBool,
		Required:     toZSSRequired(s.required, zconst.TypeBool),
		DefaultValue: defaultValueFromFunc(s.defaultFunc),
		CatchValue:   defaultValueFromFunc(s.catchFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeBool),
	}

	addCurrentGoType(&j, ctx)
	return &j
}

func (s *TimeSchema) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
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
	addCurrentGoType(&j, ctx)
	return &j
}

func (s *PointerSchema) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind:     zconst.TypePtr,
		Required: toZSSRequired(s.required, s.schema.getType()),
		Element:  s.schema.toZSS(ctx.withType(elementType(ctx.currentType))),
	}
	addCurrentGoType(&j, ctx)
	return &j
}

func (s *SliceSchema) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeSlice,
		Required:     toZSSRequired(s.required, zconst.TypeSlice),
		DefaultValue: defaultValueFromAnyFunc(s.defaultFunc),
		Processors:   processorsToZSS(rvP, zconst.TypeSlice),
		Element:      s.schema.toZSS(ctx.withType(elementType(ctx.currentType))),
	}
	addCurrentGoType(&j, ctx)
	return &j
}

func (s *MapSchema[K, V]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	defaultValue := shallowCopyMapFromFunc(s.defaultFunc)
	j := zss.ZSSSchema{
		Kind:         zconst.TypeMap,
		Required:     toZSSRequired(s.required, zconst.TypeMap),
		DefaultValue: defaultValue,
		Processors:   processorsToZSS(rvP, zconst.TypeMap),
		Key:          s.keySchema.toZSS(ctx.withType(mapKeyType(ctx.currentType))),
		Value:        s.valueSchema.toZSS(ctx.withType(mapValueType(ctx.currentType))),
	}
	addCurrentGoType(&j, ctx)
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

func (s *StructSchema) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	rvP := reflect.ValueOf(s.processors)
	j := zss.ZSSSchema{
		Kind:       zconst.TypeStruct,
		Required:   toZSSRequired(s.required, zconst.TypeStruct),
		Processors: processorsToZSS(rvP, zconst.TypeStruct),
		Fields:     toZSSFields(s.schema, ctx),
		FieldMeta:  toZSSFieldMeta(s.schema, ctx.currentType),
	}
	addCurrentGoType(&j, ctx)
	return &j
}

func (s *Custom[T]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind: zconst.TypeCustom,
		// TODO not sure this is the right place for this info
		// TODO THIS create this
		Processors: toZSSProcessorList(&s.test, zconst.TypeCustom),
	}
	addCurrentGoType(&j, ctx)
	return &j
}

func toZSSFields(m Shape, ctx *ZSSSerializeCtx) map[string]*zss.ZSSSchema {
	out := map[string]*zss.ZSSSchema{}
	for k, v := range m {
		field, _ := fieldTypeForShapeKey(ctx.currentType, k)
		out[k] = v.toZSS(ctx.withType(field.Type))
	}
	return out
}

func toZSSFieldMeta(m Shape, typ reflect.Type) map[string]zss.ZSSFieldMeta {
	out := map[string]zss.ZSSFieldMeta{}
	for k := range m {
		field, ok := fieldTypeForShapeKey(typ, k)
		if !ok {
			continue
		}
		if field.Tag == "" {
			continue
		}
		out[k] = zss.ZSSFieldMeta{Tags: string(field.Tag)}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func (s *PreprocessSchema[F, T]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind:    zconst.TypePreprocess,
		Element: s.schema.toZSS(ctx.withType(derivedType[T](ctx))),
	}
	addCurrentGoType(&j, ctx)
	return &j
}

func (s *BoxedSchema[B, T]) toZSS(ctx *ZSSSerializeCtx) *zss.ZSSSchema {
	j := zss.ZSSSchema{
		Kind:    zconst.TypeBoxed,
		Element: s.schema.toZSS(ctx.withType(derivedType[T](ctx))),
	}
	addCurrentGoType(&j, ctx)
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
