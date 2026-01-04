package zog

import (
	"strings"
	"testing"

	"github.com/Oudwins/zog/zss"
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
	s := String().Required().Default("Testing!").Catch("Testing2!").Min(1)
	serialized, err := EXPERIMENTAL_TO_ZSS(s)
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
					"message": "",
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
			"message": "",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": "Testing!",
		"catchValue": "Testing2!"
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPtr(t *testing.T) {
	s := Ptr(String().Required().Default("Testing!").Catch("Testing2!").Min(1))
	serialized, err := EXPERIMENTAL_TO_ZSS(s)
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
						"message": "",
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
				"message": "",
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
