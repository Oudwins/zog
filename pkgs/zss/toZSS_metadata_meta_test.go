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
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomString", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.CustomString", doc.Root.GoTypes[0].Display)
}

func TestIntLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomInt", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.CustomInt", doc.Root.GoTypes[0].Display)
}

func TestFloatLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomFloat float64
	s := zog.FloatLike[CustomFloat]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomFloat", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.CustomFloat", doc.Root.GoTypes[0].Display)
}

func TestBoolLikeSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomBool", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.CustomBool", doc.Root.GoTypes[0].Display)
}

func TestTimeSchemaMeta_FormatIsSet(t *testing.T) {
	format := time.RFC3339
	s := zog.Time(zog.Time.Format(format))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "time")
	assert.NotNil(t, doc.Root.Format, "Format should be set when Time.Format is used with zogmeta")
	assert.Equal(t, format, *doc.Root.Format)
}

func TestTimeSchemaMeta_FormatNotSetWhenNotConfigured(t *testing.T) {
	s := zog.Time()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "time")
	// Format should be nil if not configured
	assert.Nil(t, doc.Root.Format, "Format should be nil when not configured")
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
	assertSchemaKind(t, doc.Root, "preprocess")
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 2, "GoTypes should have two entries for Preprocess[F,T]")
	assert.Equal(t, "FromType", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.FromType", doc.Root.GoTypes[0].Display)
	assert.Equal(t, "string", doc.Root.GoTypes[1].Display)
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
	assertSchemaKind(t, doc.Root, "boxed")
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 2, "GoTypes should have two entries for Boxed[B,T]")
	assert.Equal(t, "StringBox", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.StringBox", doc.Root.GoTypes[0].Display)
	assert.Equal(t, "string", doc.Root.GoTypes[1].Display)
}

func TestCustomSchemaMeta_GoTypeIsSet(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "custom")
	assert.NotNil(t, doc.Root.GoTypes, "GoTypes should be set when zogmeta build tag is set")
	assert.Len(t, doc.Root.GoTypes, 1, "GoTypes should have one entry")
	assert.Equal(t, "CustomType", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "zss_test.CustomType", doc.Root.GoTypes[0].Display)
}

func TestPrimitiveSchemasMeta_GoTypeIsSet(t *testing.T) {
	// Primitive schemas (String, Int, etc.) have GoTypes set when EXHAUSTIVE_METADATA is enabled
	s := zog.String()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.NotNil(t, doc.Root.GoTypes, "Primitive schemas should have GoTypes set when zogmeta is enabled")
	assert.Len(t, doc.Root.GoTypes, 1, "Primitive schemas should have one GoType entry")
	assert.Equal(t, "string", doc.Root.GoTypes[0].Name)
	assert.Equal(t, "string", doc.Root.GoTypes[0].Display)
}

func TestNestedSchemasMeta_GoTypePropagation(t *testing.T) {
	type CustomString string
	s := zog.Ptr(zog.StringLike[CustomString]())
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "ptr")
	assert.Nil(t, doc.Root.GoTypes, "Ptr wrapper should not have GoTypes")

	childSchema, ok := assertChildIsSchema(t, doc.Root)
	if assert.True(t, ok) {
		assert.NotNil(t, childSchema.GoTypes, "Child schema should have GoTypes set")
		assert.Len(t, childSchema.GoTypes, 1, "Child schema should have one GoType entry")
		assert.Equal(t, "CustomString", childSchema.GoTypes[0].Name)
		assert.Equal(t, "zss_test.CustomString", childSchema.GoTypes[0].Display)
	}
}
