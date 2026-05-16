//go:build zogmeta
// +build zogmeta

package zss_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

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

func TestToJsonStringLike(t *testing.T) {
	type CustomString string
	s := zog.StringLike[CustomString]().Required().Min(1)
	d := zog.EXPERIMENTAL_TO_ZSS[CustomString](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "string",
		"goTypes": [
			{
				"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
				"name": "CustomString",
				"display": "zss_test.CustomString"
			}
		],
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
		}
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonIntLike(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]().Required().GT(0)
	d := zog.EXPERIMENTAL_TO_ZSS[CustomInt](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "number",
		"goTypes": [
			{
				"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
				"name": "CustomInt",
				"display": "zss_test.CustomInt"
			}
		],
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
		}
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonBoolLike(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]().Required()
	d := zog.EXPERIMENTAL_TO_ZSS[CustomBool](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "bool",
		"goTypes": [
			{
				"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
				"name": "CustomBool",
				"display": "zss_test.CustomBool"
			}
		],
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		}
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonTimeWithFormat(t *testing.T) {
	s := zog.Time(zog.Time.Format(time.RFC3339)).Required()
	d := zog.EXPERIMENTAL_TO_ZSS[time.Time](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "time",
		"goTypes": [
			{
				"pkgPath": "time",
				"name": "Time",
				"display": "time.Time"
			}
		],
		"format": "2006-01-02T15:04:05Z07:00",
		"required": {
			"id": "required",
			"message": "is required",
			"issuePath": null,
			"params": {}
		}
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPtrWithStringLike(t *testing.T) {
	type CustomString string
	s := zog.Ptr(zog.StringLike[CustomString]().Required())
	d := zog.EXPERIMENTAL_TO_ZSS[*CustomString](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "ptr",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "",
				"display": "*zss_test.CustomString"
			}
		],
		"element": {
			"kind": "string",
			"goTypes": [
				{
					"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
					"name": "CustomString",
					"display": "zss_test.CustomString"
				}
			],
			"required": {
				"id": "required",
				"message": "is required",
				"issuePath": null,
				"params": {}
			}
		}
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPreprocessWithGoType(t *testing.T) {
	type FromType string
	s := zog.Preprocess[FromType, string](
		func(data FromType, ctx zog.Ctx) (string, error) {
			return string(data), nil
		},
		zog.String().Min(1),
	)
	d := zog.EXPERIMENTAL_TO_ZSS[FromType](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "preprocess",
		"goTypes": [
			{
				"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
				"name": "FromType",
				"display": "zss_test.FromType"
			}
		],
		"element": {
			"kind": "string",
			"goTypes": [
				{
					"pkgPath": "",
					"name": "string",
					"display": "string"
				}
			],
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

func TestToJsonBoxedWithGoType(t *testing.T) {
	type StringBox struct {
		V string
	}
	s := zog.Boxed[StringBox, string](
		zog.String().Min(1),
		func(b StringBox, ctx zog.Ctx) (string, error) { return b.V, nil },
		func(s string, ctx zog.Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)
	d := zog.EXPERIMENTAL_TO_ZSS[StringBox](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "boxed",
		"goTypes": [
			{
				"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
				"name": "StringBox",
				"display": "zss_test.StringBox"
			}
		],
		"element": {
			"kind": "string",
			"goTypes": [
				{
					"pkgPath": "",
					"name": "string",
					"display": "string"
				}
			],
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

func TestToJsonCustomWithGoType(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	d := zog.EXPERIMENTAL_TO_ZSS[CustomType](s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "custom",
		"goTypes": [
			{
				"pkgPath": "github.com/Oudwins/zog/pkgs/zss_test",
				"name": "CustomType",
				"display": "zss_test.CustomType"
			}
		],
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
