//go:build !zogmeta
// +build !zogmeta

package zss_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Oudwins/zog"
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/stretchr/testify/assert"
)

func normalize(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", ""), "\t", "")
}

func baseZSSJson(schema string) string {
	return `{
		"$schema": "` + string(zss.ZSS_VERSION_LATEST) + `",
		"root": ` + schema + `
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
		"childs": null,
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
		"format": null,
		"processors": null,
		"childs": [
			{
				"kind": "schema",
				"schema": {
					"kind": "string",
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
					"childs": null,
					"required": {
						"id": "required",
						"message": "is required",
						"issuePath": null,
						"params": {}
					},
					"defaultValue": "Testing!",
					"catchValue": "Testing2!"
				}
			}
		],
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
	assert.NotNil(t, doc.Root)
	assert.Equal(t, "struct", doc.Root.Kind)

	// Verify child shape exists
	assert.NotNil(t, doc.Root.Childs, "Childs should not be nil")
	assert.Len(t, doc.Root.Childs, 1, "should have 1 child")
	assert.Equal(t, zss.ZSSSchemaChildKindShape, doc.Root.Childs[0].Kind, "child should be a shape")

	childShape := doc.Root.Childs[0].Shape
	assert.NotNil(t, childShape, "child shape should not be nil")
	assert.Len(t, childShape, 2, "should have 2 fields")

	// Verify name field
	nameSchema, nameExists := childShape["name"]
	assert.True(t, nameExists, "name field should exist")
	assert.Equal(t, "string", nameSchema.Kind)
}

func TestToJsonNumber(t *testing.T) {
	s := zog.Int().Required().Default(42).GT(0)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "number",
		"format": null,
		"processors": [
			{
				"kind": "test",
				"test": {
					"id": "gt",
					"message": "number must be greater than 0",
					"issuePath": null,
					"params": {
						"gt": 0
					}
				},
				"transformer": null
			}
		],
		"childs": null,
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": 42,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonBool(t *testing.T) {
	s := zog.Bool().Required().Default(true)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "bool",
		"format": null,
		"processors": null,
		"childs": null,
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": true,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonTime(t *testing.T) {
	s := zog.Time().Required()
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "time",
		"format": null,
		"processors": null,
		"childs": null,
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonSlice(t *testing.T) {
	s := zog.Slice(zog.String().Min(1)).Required().Min(1)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "slice",
		"format": null,
		"processors": [
			{
				"kind": "test",
				"test": {
					"id": "min",
					"message": "slice must contain at least 1 items",
					"issuePath": null,
					"params": {
						"min": 1
					}
				},
				"transformer": null
			}
		],
		"childs": [
			{
				"kind": "schema",
				"schema": {
					"kind": "string",
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
					"childs": null,
					"required": null,
					"defaultValue": null,
					"catchValue": null
				}
			}
		],
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonStruct(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String().Required(),
		"age":  zog.Int().Optional(),
	})
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "struct",
		"format": null,
		"processors": null,
		"childs": [
			{
				"kind": "shape",
				"shape": {
					"age": {
						"kind": "number",
						"format": null,
						"processors": null,
						"childs": null,
						"required": null,
						"defaultValue": null,
						"catchValue": null
					},
					"name": {
						"kind": "string",
						"format": null,
						"processors": null,
						"childs": null,
						"required": {
							"id": "required",
							"message": "is required",
							"issuePath": null,
							"params": {}
						},
						"defaultValue": null,
						"catchValue": null
					}
				}
			}
		],
		"required": null,
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPreprocess(t *testing.T) {
	s := zog.Preprocess(
		func(data any, ctx zog.Ctx) (any, error) {
			return data, nil
		},
		zog.String().Min(1),
	)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "preprocess",
		"format": null,
		"processors": null,
		"childs": [
			{
				"kind": "schema",
				"schema": {
					"kind": "string",
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
					"childs": null,
					"required": null,
					"defaultValue": null,
					"catchValue": null
				}
			}
		],
		"required": null,
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonBoxed(t *testing.T) {
	type StringBox struct {
		V string
	}
	s := zog.Boxed(
		zog.String().Min(1),
		func(b StringBox, ctx zog.Ctx) (string, error) { return b.V, nil },
		func(s string, ctx zog.Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "boxed",
		"format": null,
		"processors": null,
		"childs": [
			{
				"kind": "schema",
				"schema": {
					"kind": "string",
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
					"childs": null,
					"required": null,
					"defaultValue": null,
					"catchValue": null
				}
			}
		],
		"required": null,
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonMap(t *testing.T) {
	s := zog.EXPERIMENTAL_MAP[string, int](zog.String().Min(1), zog.Int().GT(0)).Required().Min(2)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "map",
		"format": null,
		"processors": [
			{
				"kind": "test",
				"test": {
					"id": "min",
					"message": "must contain at least 2 entries",
					"issuePath": null,
					"params": {
						"min": 2
					}
				},
				"transformer": null
			}
		],
		"childs": [
			{
				"kind": "shape",
				"shape": {
					"key": {
						"kind": "string",
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
						"childs": null,
						"required": null,
						"defaultValue": null,
						"catchValue": null
					},
					"value": {
						"kind": "number",
						"format": null,
						"processors": [
							{
								"kind": "test",
								"test": {
									"id": "gt",
									"message": "number must be greater than 0",
									"issuePath": null,
									"params": {
										"gt": 0
									}
								},
								"transformer": null
							}
						],
						"childs": null,
						"required": null,
						"defaultValue": null,
						"catchValue": null
					}
				}
			}
		],
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonCustom(t *testing.T) {
	s := zog.CustomFunc(func(valPtr *string, ctx zog.Ctx) bool {
		return *valPtr == "valid"
	})
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "custom",
		"format": null,
		"processors": [
			{
				"kind": "test",
				"test": {
					"id": "",
					"message": "",
					"issuePath": null,
					"params": {}
				},
				"transformer": null
			}
		],
		"childs": null,
		"required": null,
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}
