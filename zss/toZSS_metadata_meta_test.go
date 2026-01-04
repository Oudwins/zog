//go:build zogmeta
// +build zogmeta

package zss_test

import (
	"testing"
	"time"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestStringLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomString string
	s := zog.StringLike[CustomString]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "CustomString", *doc.Schema.GoType)
}

func TestIntLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "CustomInt", *doc.Schema.GoType)
}

func TestFloatLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomFloat float64
	s := zog.FloatLike[CustomFloat]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "CustomFloat", *doc.Schema.GoType)
}

func TestBoolLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "CustomBool", *doc.Schema.GoType)
}

func TestTimeSchemaMeta_FormatIsSet(t *testing.T) {
	format := time.RFC3339
	s := zog.Time(zog.Time.Format(format))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "time")
	assert.NotNil(t, doc.Schema.Format, "Format should be set when Time.Format is used with zogmeta")
	assert.Equal(t, format, *doc.Schema.Format)
}

func TestTimeSchemaMeta_FormatNotSetWhenNotConfigured(t *testing.T) {
	s := zog.Time()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "time")
	// Format should be nil if not configured
	assert.Nil(t, doc.Schema.Format, "Format should be nil when not configured")
}

func TestPreprocessSchemaMeta_GoTypeIsSet(t *testing.T) {
	type FromType string
	s := zog.Preprocess[FromType, string](
		func(data FromType, ctx zog.Ctx) (string, error) {
			return string(data), nil
		},
		zog.String(),
	)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "preprocess")
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "FromType", *doc.Schema.GoType)
}

func TestBoxedSchemaMeta_GoTypeIsSet(t *testing.T) {
	type StringBox struct {
		V string
	}
	s := zog.Boxed[StringBox, string](
		zog.String(),
		func(b StringBox, ctx zog.Ctx) (string, error) { return b.V, nil },
		func(s string, ctx zog.Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "boxed")
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "StringBox", *doc.Schema.GoType)
}

func TestCustomSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "custom")
	assert.NotNil(t, doc.Schema.GoType, "GoType should be set when zogmeta build tag is set")
	assert.Equal(t, "CustomType", *doc.Schema.GoType)
}

func TestPrimitiveSchemasMeta_GoTypeIsNil(t *testing.T) {
	// Primitive schemas (String, Int, etc.) should not have GoType set
	// Only custom types via StringLike, IntLike, etc. should have GoType
	s := zog.String()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Schema.GoType, "Primitive schemas should not have GoType set")
}

func TestNestedSchemasMeta_GoTypePropagation(t *testing.T) {
	type CustomString string
	s := zog.Ptr(zog.StringLike[CustomString]())
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "ptr")
	assert.Nil(t, doc.Schema.GoType, "Ptr wrapper should not have GoType")

	childSchema, ok := assertChildIsSchema(t, doc.Schema)
	if assert.True(t, ok) {
		assert.NotNil(t, childSchema.GoType, "Child schema should have GoType set")
		assert.Equal(t, "CustomString", *childSchema.GoType)
	}
}
