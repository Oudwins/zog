package zog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func normalize(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", ""), "\t", "")
}

func TestToJsonString(t *testing.T) {
	s := String().Required().Default("Testing!").Catch("Testing2!").Min(1)
	serialized, err := EXPERIMENTAL_TO_ZSS(s)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := `{
		"Kind": "string",
		"GoType": "",
		"Format": null,
		"Processors": [
			{
				"Kind": "test",
				"Test": {
					"ID": "min",
					"Message": "",
					"IssuePath": "",
					"Params": {
						"min": 1
					}
				},
				"Transformer": null
			}
		],
		"Child": null,
		"Required": {
			"ID": "required",
			"Message": "",
			"IssuePath": "",
			"Params": {}
		},
		"DefaultValue": "Testing!",
		"CatchValue": "Testing2!"
	}`

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPtr(t *testing.T) {
	s := Ptr(String().Required().Default("Testing!").Catch("Testing2!").Min(1))
	serialized, err := EXPERIMENTAL_TO_ZSS(s)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := `{
		"Kind": "ptr",
		"GoType": "",
		"Format": null,
		"Processors": null,
		"Child": {
			"Kind": "string",
			"GoType": "",
			"Format": null,
			"Processors": [
				{
					"Kind": "test",
					"Test": {
						"ID": "min",
						"Message": "",
						"IssuePath": "",
						"Params": {
							"min": 1
						}
					},
					"Transformer": null
				}
			],
			"Child": null,
			"Required": {
				"ID": "required",
				"Message": "",
				"IssuePath": "",
				"Params": {}
			},
			"DefaultValue": "Testing!",
			"CatchValue": "Testing2!"
		},
		"Required": null,
		"DefaultValue": null,
		"CatchValue": null
	}`

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}
