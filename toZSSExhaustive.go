//go:build zogmeta
// +build zogmeta

package zog

import "sync"

const (
	EXHAUSTIVE_METADATA = true
)

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
var EX_META_REGISTRY ExMetaRegistry = map[any]map[string]any{}

var exMetaRegistryMu sync.RWMutex

func registryAdd(r ExMetaRegistry, key any, path string, value any) {
	exMetaRegistryMu.Lock()
	defer exMetaRegistryMu.Unlock()
	if _, ok := r[key]; !ok {
		r[key] = map[string]any{}
	}
	r[key][path] = value
}

func registryGet(r ExMetaRegistry, key any, path string) (any, bool) {
	exMetaRegistryMu.RLock()
	defer exMetaRegistryMu.RUnlock()
	if m, ok := r[key]; ok {
		if val, ok := m[path]; ok {
			return val, true
		}
	}
	return nil, false
}
