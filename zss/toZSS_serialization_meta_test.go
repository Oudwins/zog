//go:build zogmeta
// +build zogmeta

package zss_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestToJsonStringLike(t *testing.T) {
	type CustomString string
	s := zog.StringLike[CustomString]().Required().Min(1)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "string",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "CustomString",
				"display": "CustomString"
			}
		],
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
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonIntLike(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]().Required().GT(0)
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "number",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "CustomInt",
				"display": "CustomInt"
			}
		],
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
		"child": null,
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

func TestToJsonBoolLike(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]().Required()
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "bool",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "CustomBool",
				"display": "CustomBool"
			}
		],
		"format": null,
		"processors": null,
		"child": null,
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

func TestToJsonTimeWithFormat(t *testing.T) {
	s := zog.Time(zog.Time.Format(time.RFC3339)).Required()
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "time",
		"format": "2006-01-02T15:04:05Z07:00",
		"processors": null,
		"child": null,
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

func TestToJsonPtrWithStringLike(t *testing.T) {
	type CustomString string
	s := zog.Ptr(zog.StringLike[CustomString]().Required())
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "ptr",
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goTypes": [
				{
					"pkgPath": "",
					"name": "CustomString",
					"display": "CustomString"
				}
			],
			"format": null,
			"processors": null,
			"child": null,
			"required": {
				"id": "required",
				"message": "is required",
				"issuePath": null,
				"params": {}
			},
			"defaultValue": null,
			"catchValue": null
		},
		"required": null,
		"defaultValue": null,
		"catchValue": null
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
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "preprocess",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "FromType",
				"display": "FromType"
			},
			{
				"pkgPath": "",
				"name": "string",
				"display": "string"
			}
		],
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goTypes": [
				{
					"pkgPath": "",
					"name": "string",
					"display": "string"
				}
			],
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
			"required": null,
			"defaultValue": null,
			"catchValue": null
		},
		"required": null,
		"defaultValue": null,
		"catchValue": null
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
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "boxed",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "StringBox",
				"display": "StringBox"
			},
			{
				"pkgPath": "",
				"name": "string",
				"display": "string"
			}
		],
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goTypes": [
				{
					"pkgPath": "",
					"name": "string",
					"display": "string"
				}
			],
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
			"required": null,
			"defaultValue": null,
			"catchValue": null
		},
		"required": null,
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonCustomWithGoType(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	d := zog.EXPERIMENTAL_TO_ZSS(s)
	serialized, err := json.Marshal(d)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := baseZSSJson(`{
		"kind": "custom",
		"goTypes": [
			{
				"pkgPath": "",
				"name": "CustomType",
				"display": "CustomType"
			}
		],
		"format": null,
		"processors": null,
		"child": null,
		"required": {
			"id": "required",
			"message": "",
			"issuePath": null,
			"params": {}
		},
		"defaultValue": null,
		"catchValue": null
	}`)

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}
