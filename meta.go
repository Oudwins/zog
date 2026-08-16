package zog

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

func (r *DefaultZogMetaRegistry) Set(s ZogSchema, k zconst.ZogMetaKey, v any) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.m[s][k] = v
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
	maps.Copy(r.m[s], m)
	return nil
}
