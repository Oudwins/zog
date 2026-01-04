package zss_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Oudwins/zog"
	zss "github.com/Oudwins/zog/zss/core"
	"github.com/stretchr/testify/assert"
)

func normalize(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", ""), "\t", "")
}

func baseZSSJson(schema string) string {
	return `{
		"version": "` + string(zss.ZSS_VERSION_LATEST) + `",
		"schema": ` + schema + `
	}`
}

func TestToJsonString(t *testing.T) {
	s := zog.String().Required().Default("Testing!").Catch("Testing2!").Min(1)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "string",
		"goType": null,
		"format": null,
		"processors": [
			{
				"kind": "test",
				"test": {
					"id": "min",
					"message": "string must contain at least 1 character(s)",
					"issuePath": null,
					"params": {
						"min": 1
					}
				},
				"transformer": null
			}
		],
		"child": null,
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": "Testing!",
		"catchValue": "Testing2!"
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPtr(t *testing.T) {
	s := zog.Ptr(zog.String().Required().Default("Testing!").Catch("Testing2!").Min(1))
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "ptr",
		"goType": null,
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goType": null,
			"format": null,
			"processors": [
				{
					"kind": "test",
					"test": {
						"id": "min",
						"message": "string must contain at least 1 character(s)",
						"issuePath": null,
						"params": {
							"min": 1
						}
					},
					"transformer": null
				}
			],
			"child": null,
			"required": {
				"id": "required",
				"message": "is required",
				"issuePath": null,
				"params": {}
			},
			"defaultValue": "Testing!",
			"catchValue": "Testing2!"
		},
		"required": null,
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

// TestToJsonStructShape tests JSON stability for struct shapes
func TestToJsonStructShape(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String().Required().Min(1),
		"age":  zog.Int().Optional(),
	})
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	// Verify JSON can be unmarshaled back into a document
	var doc zss.ZSSDocument
	err = json.Unmarshal(serialized, &doc)
	assert.Nil(t, err, "JSON should unmarshal successfully")
	assert.Equal(t, zss.ZSS_VERSION_LATEST, doc.Version)
	assert.NotNil(t, doc.Schema)
	assert.Equal(t, "struct", doc.Schema.Kind)

	// Verify child shape exists
	childShape, ok := doc.Schema.Child.(map[string]interface{})
	assert.True(t, ok, "child should be a map")
	assert.Len(t, childShape, 2, "should have 2 fields")

	// Verify name field
	nameField, nameExists := childShape["name"]
	assert.True(t, nameExists, "name field should exist")
	nameMap, ok := nameField.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "string", nameMap["kind"])
}
