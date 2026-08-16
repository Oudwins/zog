package zog

import "github.com/Oudwins/zog/zconst"

var GlobalMetaRegistry ZogMetaRegistry = NewMetaRegistry()

func (s *AnySchema) Meta(k zconst.ZogMetaKey, v any) *AnySchema {
	GlobalMetaRegistry.Set(s, k, v)
	return s
}

func (s *AnySchema) MetaFunc(f func() SchemaMetadata) *AnySchema {
	m := f()
	GlobalMetaRegistry.SetCopy(s, m)
	return s
}

func (s *AnySchema) Registry(r ZogMetaRegistry, k zconst.ZogMetaKey, v any) *AnySchema {
	r.Set(s, k, v)
	return s
}

func (s *AnySchema) RegistryFunc(f func() (ZogMetaRegistry, SchemaMetadata)) *AnySchema {
	r, m := f()
	r.SetCopy(s, m)
	return s
}
