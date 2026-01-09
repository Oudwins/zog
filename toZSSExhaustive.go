//go:build zogmeta
// +build zogmeta

package zog

const (
	EXHAUSTIVE_METADATA = true
)

// EXPERIMENTAL. PLEASE DO NOT USE UNLESS YOU KNOW WHAT YOU ARE DOING!
var EX_META_REGISTRY ExMetaRegistry = map[any]map[string]any{}

func registryAdd(r ExMetaRegistry, key any, path string, value any) {
	if _, ok := r[key]; !ok {
		r[key] = map[string]any{}
	}
	r[key][path] = value
}
