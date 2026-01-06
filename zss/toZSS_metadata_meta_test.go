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
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomString", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "CustomString", doc.Schema.GoTypes[0].Display)
}

func TestIntLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomInt", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "CustomInt", doc.Schema.GoTypes[0].Display)
}

func TestFloatLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomFloat float64
	s := zog.FloatLike[CustomFloat]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomFloat", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "CustomFloat", doc.Schema.GoTypes[0].Display)
}

func TestBoolLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomBool", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "CustomBool", doc.Schema.GoTypes[0].Display)
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
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 2, "GoTypes should have two entries for Preprocess[F,T]")
	assert.Equal(t, "FromType", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "string", doc.Schema.GoTypes[1].Display)
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
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 2, "GoTypes should have two entries for Boxed[B,T]")
	assert.Equal(t, "StringBox", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "string", doc.Schema.GoTypes[1].Display)
}

func TestCustomSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "custom")
	assert.NotNil(t, doc.Schema.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Schema.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomType", doc.Schema.GoTypes[0].Name)
	assert.Equal(t, "CustomType", doc.Schema.GoTypes[0].Display)
}

func TestPrimitiveSchemasMeta_GoTypeIsNil(t *testing.T) {
	// Primitive schemas (String, Int, etc.) should not have GoTypes set
	// Only custom types via StringLike, IntLike, etc. should have GoTypes
	s := zog.String()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Schema.GoTypes, "Primitive schemas should not have GoTypes set")
}

func TestNestedSchemasMeta_GoTypePropagation(t *testing.T) {
	type CustomString string
	s := zog.Ptr(zog.StringLike[CustomString]())
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "ptr")
	assert.Nil(t, doc.Schema.GoTypes, "Ptr wrapper should not have GoTypes")

	childSchema, ok := assertChildIsSchema(t, doc.Schema)
	if assert.True(t, ok) {
		assert.NotNil(t, childSchema.GoTypes, "Child schema should have GoTypes set")
		assert.Len(t, childSchema.GoTypes, 1, "Child schema should have one GoType entry")
		assert.Equal(t, "CustomString", childSchema.GoTypes[0].Name)
	}
}
