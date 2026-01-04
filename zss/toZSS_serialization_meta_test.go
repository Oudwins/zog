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
		"goType": "CustomString",
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
		"goType": "CustomInt",
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
		"goType": "CustomBool",
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
		"goType": null,
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
		"goType": null,
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goType": "CustomString",
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
		"goType": "FromType",
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goType": "string",
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
		"goType": "StringBox",
		"format": null,
		"processors": null,
		"child": {
			"kind": "string",
			"goType": "string",
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
		"goType": "CustomType",
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
