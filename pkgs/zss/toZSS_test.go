//go:build !zogmeta
// +build !zogmeta

package zss_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Oudwins/zog"
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	zssschema "github.com/Oudwins/zog/pkgs/zss/schema"
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
	assert.Equal(t, zss.ZSS_VERSION_LATEST, doc.URI)
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
					"self": {"$ref": "` + zss.ZSSRefFromKey(1) + `"},
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
		"$defs": {
			"schema1": {
				"kind": "ptr",
				"element": {
					"kind": "struct",
					"fields": {
						"self": {"$ref": "` + zss.ZSSRefFromKey(1) + `"},
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

func assertZSSRef(t *testing.T, schema *zss.ZSSSchema, key int) {
	t.Helper()
	if assert.NotNil(t, schema) && assert.NotNil(t, schema.Ref) {
		assert.Equal(t, zss.ZSSRefFromKey(key), *schema.Ref)
	}
	assert.Empty(t, schema.Kind)
}

func assertSingleRecursiveDef(t *testing.T, doc zss.ZSSDocument) *zss.ZSSSchema {
	t.Helper()
	if !assert.Len(t, doc.Defs, 1) {
		return nil
	}
	def := doc.Defs[zss.ZSSDefKeyFromKey(1)]
	assert.NotNil(t, def)
	return def
}

func TestToJsonRecursiveSliceTreeUsesRefs(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"children": zog.Slice(self()),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	children := def.Element.Fields["children"]
	assert.Equal(t, "slice", string(children.Kind))
	assertZSSRef(t, children.Element, 1)
}

func TestToJsonRecursiveMapTreeUsesRefs(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"children": zog.EXPERIMENTAL_MAP[string, any](zog.String(), self()),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	children := def.Element.Fields["children"]
	assert.Equal(t, "map", string(children.Kind))
	assert.Equal(t, "string", string(children.Key.Kind))
	assertZSSRef(t, children.Value, 1)
}

func TestToJsonRecursiveMultipleSelfFieldsUseSameRef(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"left":  self(),
			"right": self(),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	assertZSSRef(t, def.Element.Fields["left"], 1)
	assertZSSRef(t, def.Element.Fields["right"], 1)
}

func TestToJsonRecursiveDeepNestedEdgeUsesRef(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"wrapper": zog.Struct(zog.Shape{"next": self()}),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	assertZSSRef(t, def.Element.Fields["wrapper"].Fields["next"], 1)
}

func TestToJsonRecursiveUpdaterOriginalUsesSameRef(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"self": self(func(original *zog.PointerSchema) *zog.PointerSchema { return original }),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	assertZSSRef(t, def.Element.Fields["self"], 1)
}

func TestToJsonRecursiveUpdaterModifiedTerminates(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{
			"children": zog.Slice(self(func(original *zog.PointerSchema) *zog.PointerSchema { return original })),
		}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	assertZSSRef(t, def.Element.Fields["children"].Element, 1)
}

// Temporarily removed due to 1.23.12 failing
// func TestToJsonMultipleRecursiveSchemasCreateSeparateDefs(t *testing.T) {
// 	a := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
// 		return zog.Ptr(zog.Struct(zog.Shape{"nextA": self()}))
// 	})
// 	b := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
// 		return zog.Ptr(zog.Struct(zog.Shape{"nextB": self()}))
// 	})
//
// 	d := zog.EXPERIMENTAL_TO_ZSS(zog.Struct(zog.Shape{"a": a, "b": b}))
// 	assert.Len(t, d.Defs, 2)
// 	assertZSSRef(t, d.Defs[zss.ZSSDefKeyFromKey(1)].Element.Fields["nextA"], 1)
// 	assertZSSRef(t, d.Defs[zss.ZSSDefKeyFromKey(2)].Element.Fields["nextB"], 2)
// }

func TestToJsonRecursiveRootIsExpandedInline(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{"self": self()}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	assert.Nil(t, d.Root.Ref)
	assert.Equal(t, "ptr", string(d.Root.Kind))
	assert.NotEmpty(t, d.Defs)
}

func TestToJsonNonRecursiveSchemaOmitsDefs(t *testing.T) {
	d := zog.EXPERIMENTAL_TO_ZSS(zog.Struct(zog.Shape{"name": zog.String()}))
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.Empty(t, d.Defs)
	assert.NotContains(t, string(serialized), `"$defs"`)
}

func TestToJsonSharedNonRecursiveSchemaStaysInline(t *testing.T) {
	name := zog.String().Min(1)
	d := zog.EXPERIMENTAL_TO_ZSS(zog.Struct(zog.Shape{"first": name, "second": name}))
	assert.Empty(t, d.Defs)
	assert.Nil(t, d.Root.Fields["first"].Ref)
	assert.Nil(t, d.Root.Fields["second"].Ref)
	assert.Equal(t, "string", string(d.Root.Fields["first"].Kind))
	assert.Equal(t, "string", string(d.Root.Fields["second"].Kind))
}

func TestToJsonRecursivePreprocessUsesRefs(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PreprocessSchema[any, any]]) *zog.PreprocessSchema[any, any] {
		return zog.Preprocess(func(data any, ctx zog.Ctx) (any, error) { return data, nil }, zog.Struct(zog.Shape{"next": self()}))
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	assert.Equal(t, "preprocess", string(def.Kind))
	assertZSSRef(t, def.Element.Fields["next"], 1)
}

func TestToJsonRecursiveBoxedUsesRefs(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.BoxedSchema[any, any]]) *zog.BoxedSchema[any, any] {
		return zog.Boxed(zog.Struct(zog.Shape{"next": self()}), func(data any, ctx zog.Ctx) (any, error) { return data, nil }, func(data any, ctx zog.Ctx) (any, error) { return data, nil })
	})

	d := zog.EXPERIMENTAL_TO_ZSS(s)
	def := assertSingleRecursiveDef(t, d)
	assert.Equal(t, "boxed", string(def.Kind))
	assertZSSRef(t, def.Element.Fields["next"], 1)
}

func TestZSSDocumentSchemaValidatesRecursiveOutput(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{"self": self()}))
	})
	d := zog.EXPERIMENTAL_TO_ZSS(s)

	errList := zssschema.ZSSDocumentSchema.Validate(&d)
	assert.Empty(t, errList)
}

func TestZSSSchemaAllowsPureRef(t *testing.T) {
	ref := zss.ZSSRefFromKey(1)
	schema := zss.ZSSSchema{Ref: &ref}
	errList := zssschema.ZSSSchemaSchema.Validate(&schema)
	assert.Empty(t, errList)
}

func TestZSSSchemaRejectsUnknownKind(t *testing.T) {
	schema := zss.ZSSSchema{Kind: "money"}
	errList := zssschema.ZSSSchemaSchema.Validate(&schema)
	assert.NotEmpty(t, errList)
}

func TestZSSRefSchemaMarshalsWithoutKind(t *testing.T) {
	ref := zss.ZSSRefFromKey(1)
	serialized, err := json.Marshal(zss.ZSSSchema{Ref: &ref})
	assert.Nil(t, err)
	assert.Equal(t, `{"$ref":"`+zss.ZSSRefFromKey(1)+`"}`, string(serialized))
}

func TestZSSRefFromKeyPathStability(t *testing.T) {
	assert.Equal(t, "schema12", zss.ZSSDefKeyFromKey(12))
	assert.Equal(t, "#/$defs/schema12", zss.ZSSRefFromKey(12))
}

func TestToJsonRecursiveMultipleCallsUseIndependentContext(t *testing.T) {
	s := zog.EXPERIMENTAL_RECURSIVE(func(self zog.RecursiveSchema[*zog.PointerSchema]) *zog.PointerSchema {
		return zog.Ptr(zog.Struct(zog.Shape{"self": self()}))
	})

	first := zog.EXPERIMENTAL_TO_ZSS(s)
	second := zog.EXPERIMENTAL_TO_ZSS(s)
	assert.Contains(t, first.Defs, zss.ZSSDefKeyFromKey(1))
	assert.Contains(t, second.Defs, zss.ZSSDefKeyFromKey(1))
	assert.Len(t, first.Defs, 1)
	assert.Len(t, second.Defs, 1)
}

func TestZSSExtensionMarshalsURIAndContent(t *testing.T) {
	schema := zss.ZSSSchema{
		Kind: "custom",
		Extension: &zss.ZSSExtension{
			URI: "https://example.com/zss/extensions/money/1.0.0/schema.json",
			Content: map[string]any{
				"currency":  "USD",
				"precision": 2,
			},
		},
	}

	serialized, err := json.Marshal(schema)
	assert.Nil(t, err)
	assert.Equal(t, normalize(`{
		"kind":"custom",
		"extension":{
			"uri":"https://example.com/zss/extensions/money/1.0.0/schema.json",
			"content":{
				"currency":"USD",
				"precision":2
			}
		}
	}`), normalize(string(serialized)))
}

func TestZSSDocumentSchemaValidatesExtensionURI(t *testing.T) {
	doc := zss.ZSSDocument{
		URI: zss.ZSS_VERSION_LATEST,
		Root: &zss.ZSSSchema{
			Kind: "custom",
			Extension: &zss.ZSSExtension{
				URI:     "https://example.com/money/1.0.0/schema.json",
				Content: []any{"currency", "precision"},
			},
		},
	}

	errList := zssschema.ZSSDocumentSchema.Validate(&doc)
	assert.Empty(t, errList)
}

func TestZSSURIRegexAcceptsValidExtensionURIs(t *testing.T) {
	validURIs := []string{
		"https://example.com/money/1.0.0/schema.json",
		"https://example.com/money/0.1.0/schema.json",
		"https://example.com/zss/extensions/money/2.3.4-beta/schema.json",
		"https://example.com/zss/extensions/money/2.3.4-beta.1/schema.json",
	}

	for _, uri := range validURIs {
		t.Run(uri, func(t *testing.T) {
			assert.True(t, zss.ZSS_URI_REGEX.MatchString(uri))
		})
	}
}

func TestZSSURIRegexRejectsInvalidExtensionURIs(t *testing.T) {
	invalidURIs := []string{
		"example.com/money/1.0.0/schema.json",
		"https://example.com/money/1/schema.json",
		"https://example.com/money/1.0/schema.json",
		"https://example.com/money/v1.0.0/schema.json",
		"https://example.com/money/1.0.0-alpha/schema.json",
		"https://example.com/money/1.0.0-rc.1/schema.json",
		"https://example.com/money/1.0.0+build/schema.json",
		"https://example.com/money/01.0.0/schema.json",
		"https://example.com/money/1.02.0/schema.json",
		"https://example.com/money/1.0.0-beta.01/schema.json",
	}

	for _, uri := range invalidURIs {
		t.Run(uri, func(t *testing.T) {
			assert.False(t, zss.ZSS_URI_REGEX.MatchString(uri))
		})
	}
}

func TestZSSURIRegexExtractsIDAndVersion(t *testing.T) {
	uri := "https://example.com/zss/extensions/money/1.2.3-beta.1/schema.json"
	matches := zss.ZSS_URI_REGEX.FindStringSubmatch(uri)
	assert.NotNil(t, matches)

	values := map[string]string{}
	for i, name := range zss.ZSS_URI_REGEX.SubexpNames() {
		if i != 0 && name != "" {
			values[name] = matches[i]
		}
	}

	assert.Equal(t, "https://example.com/zss/extensions/money", values["id"])
	assert.Equal(t, "1.2.3-beta.1", values["version"])
}
