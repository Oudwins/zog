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
		"element": {
					"kind": "string",
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
					"required": {
						"id": "required",
						"message": "is required",
						"issuePath": null,
						"params": {}
					},
					"defaultValue": "Testing!",
					"catchValue": "Testing2!"
		}
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

	// Verify fields shape exists
	childShape := doc.Root.Fields
	assert.NotNil(t, childShape, "child shape should not be nil")
	assert.Len(t, childShape, 2, "should have 2 fields")

	// Verify name field
	nameSchema, nameExists := childShape["name"]
	assert.True(t, nameExists, "name field should exist")
	assert.Equal(t, "string", string(nameSchema.Kind))
}

func TestToJsonNumber(t *testing.T) {
	s := zog.Int().Required().Default(42).GT(0)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "number",
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
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": 42
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
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": true
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
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		}
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
		"element": {
					"kind": "string",
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
					]
		},
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		}
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
		"fields": {
					"age": {
						"kind": "number"
					},
					"name": {
						"kind": "string",
						"required": {
							"id": "required",
							"message": "is required",
							"issuePath": null,
							"params": {}
						}
					}
		}
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
		"element": {
					"kind": "string",
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
					]
		}
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
		"element": {
					"kind": "string",
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
					]
		}
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
		"key": {
			"kind": "string",
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
			]
		},
		"value": {
			"kind": "number",
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
			]
		},
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		}
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
		]
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonRecursiveUsesRefs(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"value": zog.Int().Required(),
			"self":  self(),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)
	assert.NotEmpty(t, d.Defs)

	expected := `{
		"$schema": "` + string(zss.ZSS_VERSION_LATEST) + `",
		"root": {
			"kind": "ptr",
			"element": {
				"kind": "struct",
				"fields": {
					"self": {"$ref": "#/defs/schema1"},
					"value": {
						"kind": "number",
						"required": {
							"id": "required",
							"message": "is required",
							"issuePath": null,
							"params": {}
						}
					}
				}
			}
		},
		"defs": {
			"schema1": {
				"kind": "ptr",
				"element": {
					"kind": "struct",
					"fields": {
						"self": {"$ref": "#/defs/schema1"},
						"value": {
							"kind": "number",
							"required": {
								"id": "required",
								"message": "is required",
								"issuePath": null,
								"params": {}
							}
						}
					}
				}
			}
		}
	}`

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}
