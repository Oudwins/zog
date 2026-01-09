//go:build !zogmeta
// +build !zogmeta

package zog

const (
	EXHAUSTIVE_METADATA = false
)

var EX_META_REGISTRY ExMetaRegistry = nil

func registryAdd(r ExMetaRegistry, key any, path string, value any) {
	// no op
}
