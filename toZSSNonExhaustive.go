//go:build !zogmeta
// +build !zogmeta

package zog

const (
	EXHAUSTIVE_METADATA = false
)

var exMetaRegistry ExMetaRegistry = nil

func RegistryAdd(_ ExMetaRegistry, _ any, _ string, _ any) {
	// no op
}

func RegistryGet(_ ExMetaRegistry, _ any, _ string) (any, bool) {
	return nil, false
}
