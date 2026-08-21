package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/require"
)

func TestDefaultZogMetaRegistryInitializesSchemaMetadata(t *testing.T) {
	schema := String()
	registry := NewMetaRegistry()

	require.NoError(t, registry.Set(schema, zconst.MetaKeyDescription, "description"))
	require.NoError(t, registry.SetCopy(schema, SchemaMetadata{"example": "value"}))

	description, err := registry.Get(schema, zconst.MetaKeyDescription)
	require.NoError(t, err)
	require.Equal(t, "description", description)
	example, err := registry.Get(schema, "example")
	require.NoError(t, err)
	require.Equal(t, "value", example)
}

func TestDefaultZogMetaRegistryZeroValue(t *testing.T) {
	schema := String()
	var registry DefaultZogMetaRegistry

	require.NoError(t, registry.Set(schema, zconst.MetaKeyDescription, "description"))
	require.NoError(t, registry.SetCopy(schema, SchemaMetadata{"example": "value"}))

	description, err := registry.Get(schema, zconst.MetaKeyDescription)
	require.NoError(t, err)
	require.Equal(t, "description", description)
	example, err := registry.Get(schema, "example")
	require.NoError(t, err)
	require.Equal(t, "value", example)
}
