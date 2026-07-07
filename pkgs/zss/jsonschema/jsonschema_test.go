package zjsonschema_test

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	zjsonschema "github.com/Oudwins/zog/pkgs/zss/jsonschema"
	"github.com/Oudwins/zog/zconst"
	googlejsonschema "github.com/google/jsonschema-go/jsonschema"
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
	requireValidJSONSchema(t, schema)

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
			requireValidJSONSchema(t, got)
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
	requireValidJSONSchema(t, schema)

	assert.Equal(t, zjsonschema.Schema{
		"$schema": string(zjsonschema.Draft2020_12),
		"type":    "object",
		"additionalProperties": zjsonschema.Schema{
			"type":  "array",
			"items": zjsonschema.Schema{"type": "string"},
		},
	}, schema)
}

func TestFromZSSConvertsUnion(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeUnion, Children: []*zsscore.ZSSSchema{
		{Kind: zconst.TypeString},
		{Kind: zconst.TypeNumber},
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)
	requireValidJSONSchema(t, schema)

	assert.Equal(t, zjsonschema.Schema{
		"$schema": string(zjsonschema.Draft2020_12),
		"anyOf": []any{
			zjsonschema.Schema{"type": "string"},
			zjsonschema.Schema{"type": "number"},
		},
	}, schema)
}

func TestFromZSSConvertsPointerNullability(t *testing.T) {
	optional, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypePtr, Element: &zsscore.ZSSSchema{Kind: zconst.TypeString}}}, zjsonschema.Options{})
	require.NoError(t, err)
	requireValidJSONSchema(t, optional)
	assert.Equal(t, []string{"string", "null"}, optional["type"])

	required, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypePtr, Required: requiredTest(), Element: &zsscore.ZSSSchema{Kind: zconst.TypeString}}}, zjsonschema.Options{})
	require.NoError(t, err)
	requireValidJSONSchema(t, required)
	assert.Equal(t, "string", required["type"])

	ref := zsscore.ZSSRefFromKey(1)
	refPtr, err := zjsonschema.FromZSS(zsscore.ZSSDocument{
		Root: &zsscore.ZSSSchema{Kind: zconst.TypePtr, Element: &zsscore.ZSSSchema{Ref: &ref}},
		Defs: map[string]*zsscore.ZSSSchema{
			zsscore.ZSSDefKeyFromKey(1): {Kind: zconst.TypeString},
		},
	}, zjsonschema.Options{})
	require.NoError(t, err)
	requireValidJSONSchema(t, refPtr)
	assert.Equal(t, []any{zjsonschema.Schema{"$ref": ref}, zjsonschema.Schema{"type": "null"}}, refPtr["anyOf"])
}

func TestFromZSSConvertsValidationProcessors(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString, Processors: []zsscore.ZSSProcessor{
		testProcessor(zconst.IssueCodeMin, map[string]any{zconst.IssueCodeMin: 2}),
		testProcessor(zconst.IssueCodeMax, map[string]any{zconst.IssueCodeMax: 5}),
		testProcessor(zconst.IssueCodeEmail, map[string]any{}),
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)
	requireValidJSONSchema(t, schema)

	assert.Equal(t, 2, schema["minLength"])
	assert.Equal(t, 5, schema["maxLength"])
	assert.Equal(t, "email", schema["format"])
}

func TestFromZSSConvertsUnknownKindWithOption(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: "uuid"}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{
		UnknownKindConverter: func(schema *zsscore.ZSSSchema) (zjsonschema.Schema, error) {
			if schema.Kind == "uuid" {
				return zjsonschema.Schema{"type": "string", "format": "uuid"}, nil
			}
			return nil, errors.New("unexpected kind")
		},
	})
	require.NoError(t, err)
	requireValidJSONSchema(t, schema)

	assert.Equal(t, zjsonschema.Schema{
		"$schema": string(zjsonschema.Draft2020_12),
		"type":    "string",
		"format":  "uuid",
	}, schema)
}

func TestFromZSSConvertsNestedUnknownKindWithOption(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeSlice, Element: &zsscore.ZSSSchema{Kind: "uuid"}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{
		UnknownKindConverter: func(schema *zsscore.ZSSSchema) (zjsonschema.Schema, error) {
			return zjsonschema.Schema{"type": "string", "format": schema.Kind}, nil
		},
	})
	require.NoError(t, err)
	requireValidJSONSchema(t, schema)

	assert.Equal(t, zjsonschema.Schema{
		"$schema": string(zjsonschema.Draft2020_12),
		"type":    "array",
		"items":   zjsonschema.Schema{"type": "string", "format": "uuid"},
	}, schema)
}

func TestFromZSSReturnsUnknownKindConverterErrors(t *testing.T) {
	wantErr := errors.New("custom kind failed")

	_, err := zjsonschema.FromZSS(zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: "custom_kind"}}, zjsonschema.Options{
		UnknownKindConverter: func(schema *zsscore.ZSSSchema) (zjsonschema.Schema, error) {
			return nil, wantErr
		},
	})

	require.ErrorIs(t, err, wantErr)
}

func TestFromZSSConvertsCustomTestsWithOption(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString, Processors: []zsscore.ZSSProcessor{
		testProcessor(zconst.IssueCodeMin, map[string]any{zconst.IssueCodeMin: 2}),
		testProcessor("starts_with", map[string]any{"starts_with": "z"}),
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{
		TestConverter: func(out zjsonschema.Schema, kind zconst.ZogType, test *zsscore.ZSSTest) error {
			if err := zjsonschema.ConvertTest(out, kind, test); err != nil {
				return err
			}
			if test.ID == "starts_with" {
				out["pattern"] = "^" + test.Params["starts_with"].(string)
			}
			return nil
		},
	})
	require.NoError(t, err)
	requireValidJSONSchema(t, schema)

	assert.Equal(t, 2, schema["minLength"])
	assert.Equal(t, "^z", schema["pattern"])
}

func TestFromZSSConvertsNamedGroupsInMatchPatterns(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString, Processors: []zsscore.ZSSProcessor{
		testProcessor(zconst.IssueCodeMatch, map[string]any{zconst.IssueCodeMatch: regexp.MustCompile(`^(?P<id>[a-z]+)/(?P<version>[0-9]+)$`).String()}),
	}}}

	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	require.NoError(t, err)
	requireValidJSONSchema(t, schema)

	assert.Equal(t, `^(?:[a-z]+)/(?:[0-9]+)$`, schema["pattern"])
}

func TestFromZSSReturnsErrorForNonStringMatchPattern(t *testing.T) {
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString, Processors: []zsscore.ZSSProcessor{
		testProcessor(zconst.IssueCodeMatch, map[string]any{zconst.IssueCodeMatch: 1}),
	}}}

	_, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})

	require.Error(t, err)
	assert.ErrorContains(t, err, `convert test "match"`)
	assert.ErrorContains(t, err, "match test param must be a string")
}

func TestFromZSSReturnsTestConverterErrors(t *testing.T) {
	wantErr := errors.New("custom test failed")
	doc := zsscore.ZSSDocument{Root: &zsscore.ZSSSchema{Kind: zconst.TypeString, Processors: []zsscore.ZSSProcessor{
		testProcessor("custom_test", map[string]any{}),
	}}}

	_, err := zjsonschema.FromZSS(doc, zjsonschema.Options{
		TestConverter: func(out zjsonschema.Schema, kind zconst.ZogType, test *zsscore.ZSSTest) error {
			return wantErr
		},
	})

	require.ErrorIs(t, err, wantErr)
	assert.ErrorContains(t, err, `convert test "custom_test"`)
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
	requireValidJSONSchema(t, schema)

	properties := schema["properties"].(zjsonschema.Schema)
	assert.Contains(t, properties, "full_name")
	assert.Contains(t, properties, "email_address")
	assert.NotContains(t, properties, "hidden")
	assert.Equal(t, []string{"email_address", "full_name"}, schema["required"])
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

func requireValidJSONSchema(t *testing.T, schema zjsonschema.Schema) {
	t.Helper()

	data, err := json.Marshal(schema)
	require.NoError(t, err)

	var parsed googlejsonschema.Schema
	require.NoError(t, json.Unmarshal(data, &parsed))

	_, err = parsed.Resolve(nil)
	require.NoError(t, err)
}
