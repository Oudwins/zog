//go:build !zogmeta
// +build !zogmeta

package zog

const (
	EXHAUSTIVE_METADATA = false
)

var EX_META_REGISTRY ExMetaRegistry = nil

func registryAdd(_ ExMetaRegistry, _ any, _ string, _ any) {
	// no op
}

func registryGet(_ ExMetaRegistry, _ any, _ string) (any, bool) {
	return nil, false
}
