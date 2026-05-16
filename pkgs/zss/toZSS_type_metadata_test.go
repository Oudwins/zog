package zss_test

import (
	"testing"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type zssTypeMetaUser struct {
	Name  string         `json:"full_name,omitempty" zog:"zog_name"`
	Email *string        `zog:"email_address"`
	Tags  []string       `json:"tags"`
	Attrs map[string]int `json:"attrs"`
	Hide  string         `json:"-"`
}

func TestToZSSTypeMetadataRootPrimitive(t *testing.T) {
	doc := zog.EXPERIMENTAL_TO_ZSS[string](zog.String())

	require.Len(t, doc.Root.GoTypes, 1)
	assert.Equal(t, "string", doc.Root.GoTypes[0].Display)
}

func TestToZSSTypeMetadataStructFields(t *testing.T) {
	doc := zog.EXPERIMENTAL_TO_ZSS[zssTypeMetaUser](zog.Struct(zog.Shape{
		"name":  zog.String(),
		"email": zog.Ptr(zog.String()),
		"tags":  zog.Slice(zog.String()),
		"attrs": zog.EXPERIMENTAL_MAP[string, int](zog.String(), zog.Int()),
		"hide":  zog.String(),
	}))

	require.Len(t, doc.Root.GoTypes, 1)
	assert.Equal(t, "zss_test.zssTypeMetaUser", doc.Root.GoTypes[0].Display)
	assert.Equal(t, `json:"full_name,omitempty" zog:"zog_name"`, doc.Root.FieldMeta["name"].Tags)
	assert.Equal(t, `json:"-"`, doc.Root.FieldMeta["hide"].Tags)
	assert.Equal(t, "string", doc.Root.Fields["name"].GoTypes[0].Display)
	assert.Equal(t, "*string", doc.Root.Fields["email"].GoTypes[0].Display)
	assert.Equal(t, "string", doc.Root.Fields["email"].Element.GoTypes[0].Display)
	assert.Equal(t, "[]string", doc.Root.Fields["tags"].GoTypes[0].Display)
	assert.Equal(t, "string", doc.Root.Fields["tags"].Element.GoTypes[0].Display)
	assert.Equal(t, "map[string]int", doc.Root.Fields["attrs"].GoTypes[0].Display)
	assert.Equal(t, "string", doc.Root.Fields["attrs"].Key.GoTypes[0].Display)
	assert.Equal(t, "int", doc.Root.Fields["attrs"].Value.GoTypes[0].Display)
}
