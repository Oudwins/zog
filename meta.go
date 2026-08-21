package zog

//go:generate go run ./cmd/gen/metadata

import (
	"maps"
	"sync"

	"github.com/Oudwins/zog/zconst"
)

type ZogMetaRegistry interface {
	Set(s ZogSchema, k zconst.ZogMetaKey, v any) error
	SetCopy(s ZogSchema, m SchemaMetadata) error
	Get(s ZogSchema, k zconst.ZogMetaKey) (any, error)
}

type SchemaMetadata = map[string]any

var GlobalMetaRegistry ZogMetaRegistry = NewMetaRegistry()

var _ ZogMetaRegistry = &DefaultZogMetaRegistry{}

func NewMetaRegistry() ZogMetaRegistry {
	return &DefaultZogMetaRegistry{
		m:   map[any]SchemaMetadata{},
		mux: sync.RWMutex{},
	}
}

type DefaultZogMetaRegistry struct {
	m   map[any]SchemaMetadata
	mux sync.RWMutex
}

func (r *DefaultZogMetaRegistry) metadataFor(s ZogSchema) SchemaMetadata {
	if r.m == nil {
		r.m = map[any]SchemaMetadata{}
	}
	if r.m[s] == nil {
		r.m[s] = SchemaMetadata{}
	}
	return r.m[s]
}

func (r *DefaultZogMetaRegistry) Set(s ZogSchema, k zconst.ZogMetaKey, v any) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.metadataFor(s)[k] = v
	return nil
}

func (r *DefaultZogMetaRegistry) Get(s ZogSchema, k zconst.ZogMetaKey) (any, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.m[s][k], nil
}

func (r *DefaultZogMetaRegistry) SetCopy(s ZogSchema, m SchemaMetadata) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	maps.Copy(r.metadataFor(s), m)
	return nil
}
