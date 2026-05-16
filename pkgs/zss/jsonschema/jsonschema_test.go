package zjsonschema_test

import (
	"testing"

	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	zjsonschema "github.com/Oudwins/zog/pkgs/zss/jsonschema"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromZSSConvertsRootAndDefs(t *testing.T) {
	ref := zsscore.ZSSRefFromKey(1)
	doc := zsscore.ZSSDocument{
		Root: &zsscore.ZSSSchema{Kind: zconst.TypeStruct, Fields: map[string]*zsscore.ZSSSchema{
			"name": {Kind: zconst.TypeString, Required: requiredTest()},
			"node": {Ref: &ref},
		}},
		Defs: map[string]*zsscore.ZSSSchema{
			zsscore.ZSSDefKeyFromKey(1): {Kind: zconst.TypeNumber},
		},
	}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)

	assert.Equal(t, string(zjsonschema.Draft2020_12), schema["$schema"])
	assert.Equal(t, "object", schema["type"])
	assert.Equal(t, []string{"name"}, schema["required"])
	assert.Equal(t, zjsonschema.Schema{"schema1": zjsonschema.Schema{"type": "number"}}, schema["$defs"])

	properties := schema["properties"].(zjsonschema.Schema)
	assert.Equal(t, zjsonschema.Schema{"type": "string"}, properties["name"])
	assert.Equal(t, zjsonschema.Schema{"$ref": "#/$defs/schema1"}, properties["node"])
}

func TestFromZSSConvertsKinds(t *testing.T) {
	tests := []struct {
		name string
		in   *zsscore.ZSSSchema
		want zjsonschema.Schema
	}{
		{name: "string", in: &zsscore.ZSSSchema{Kind: zconst.TypeString}, want: zjsonschema.Schema{"$schema": string(zjsonschema.Draft2020_12), "type": "string"}},
		{name: "number", in: &zsscore.ZSSSchema{Kind: zconst.TypeNumber}, want: zjsonschema.Schema{"$schema": string(zjsonschema.Draft2020_12), "type": "number"}},
		{name: "bool", in: &zsscore.ZSSSchema{Kind: zconst.TypeBool}, want: zjsonschema.Schema{"$schema": string(zjsonschema.Draft2020_12), "type": "boolean"}},
		{name: "time", in: &zsscore.ZSSSchema{Kind: zconst.TypeTime}, want: zjsonschema.Schema{"$schema": string(zjsonschema.Draft2020_12), "type": "string", "format": "date-time"}},
		{name: "any", in: &zsscore.ZSSSchema{Kind: zconst.TypeAny}, want: zjsonschema.Schema{"$schema": string(zjsonschema.Draft2020_12)}},
		{name: "custom", in: &zsscore.ZSSSchema{Kind: zconst.TypeCustom}, want: zjsonschema.Schema{"$schema": string(zjsonschema.Draft2020_12)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: tt.in}, zjsonschema.Options{})
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromZSSConvertsContainers(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeMap, Value: &zsscore.ZSSSchema{
		Kind:    zconst.TypeSlice,
		Element: &zsscore.ZSSSchema{Kind: zconst.TypeString},
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)

	assert.Equal(t, zjsonschema.Schema{
		"$schema": string(zjsonschema.Draft2020_12),
		"type":    "object",
		"additionalProperties": zjsonschema.Schema{
			"type":  "array",
			"items": zjsonschema.Schema{"type": "string"},
		},
	}, schema)
}

func TestFromZSSConvertsPointerNullability(t *testing.T) {
	optional, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypePtr, Element: &zsscore.ZSSSchema{Kind: zconst.TypeString}}}, zjsonschema.Options{})
	require.NoError(t, err)
	assert.Equal(t, []string{"string", "null"}, optional["type"])

	required, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypePtr, Required: requiredTest(), Element: &zsscore.ZSSSchema{Kind: zconst.TypeString}}}, zjsonschema.Options{})
	require.NoError(t, err)
	assert.Equal(t, "string", required["type"])

	ref := zsscore.ZSSRefFromKey(1)
	refPtr, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypePtr, Element: &zsscore.ZSSSchema{Ref: &ref}}}, zjsonschema.Options{})
	require.NoError(t, err)
	assert.Equal(t, []any{zjsonschema.Schema{"$ref": ref}, zjsonschema.Schema{"type": "null"}}, refPtr["anyOf"])
}

func TestFromZSSConvertsValidationProcessors(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString, Processors: []zsscore.ZSSProcessor{
		testProcessor(zconst.IssueCodeMin, map[string]any{"min": 2}),
		testProcessor(zconst.IssueCodeMax, map[string]any{"max": 5}),
		testProcessor(zconst.IssueCodeEmail, map[string]any{}),
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)

	assert.Equal(t, 2, schema["minLength"])
	assert.Equal(t, 5, schema["maxLength"])
	assert.Equal(t, "email", schema["format"])
}

func TestFromZSSUsesFieldMetaPropertyNames(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeStruct, Fields: map[string]*zsscore.ZSSSchema{
		"name":   {Kind: zconst.TypeString, Required: requiredTest()},
		"email":  {Kind: zconst.TypeString, Required: requiredTest()},
		"hidden": {Kind: zconst.TypeString, Required: requiredTest()},
	}, FieldMeta: map[string]zsscore.ZSSFieldMeta{
		"name":   {Tags: `json:"full_name,omitempty" zog:"zog_name"`},
		"email":  {Tags: `zog:"email_address"`},
		"hidden": {Tags: `json:"-"`},
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)

	properties := schema["properties"].(zjsonschema.Schema)
	assert.Contains(t, properties, "full_name")
	assert.Contains(t, properties, "email_address")
	assert.NotContains(t, properties, "hidden")
	assert.Equal(t, []string{"full_name", "email_address"}, schema["required"])
}

func TestFromZSSReturnsErrors(t *testing.T) {
	_, err := zjsonschema.FromZSS(zsscore.ZSSDocument{}, zjsonschema.Options{})
	assert.Error(t, err)

	_, err = zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString}}, zjsonschema.Options{Draft: "draft-07"})
	assert.Error(t, err)

	_, err = zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: "unknown"}}, zjsonschema.Options{})
	assert.Error(t, err)
}

func requiredTest() *zsscore.ZSSTest {
	return &zsscore.ZSSTest{ID: zconst.IssueCodeRequired, Params: map[string]any{}}
}

func testProcessor(id zconst.ZogIssueCode, params map[string]any) zsscore.ZSSProcessor {
	return zsscore.ZSSProcessor{Kind: zconst.ZogProcessorTest, Test: &zsscore.ZSSTest{ID: id, Params: params}}
}
