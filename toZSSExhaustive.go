//go:build zogmeta
// +build zogmeta

package zog

import "sync"

const (
	EXHAUSTIVE_METADATA = true
)

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
var exMetaRegistry ExMetaRegistry = map[any]map[string]any{}

var exMetaRegistryMu sync.RWMutex

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func RegistryAdd(r ExMetaRegistry, key any, path string, value any) {
	exMetaRegistryMu.Lock()
	defer exMetaRegistryMu.Unlock()
	if _, ok := r[key]; !ok {
		r[key] = map[string]any{}
	}
	r[key][path] = value
}

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
func RegistryGet(r ExMetaRegistry, key any, path string) (any, bool) {
	exMetaRegistryMu.RLock()
	defer exMetaRegistryMu.RUnlock()
	if m, ok := r[key]; ok {
		if val, ok := m[path]; ok {
			return val, true
		}
	}
	return nil, false
}
