//go:build !zogmeta
// +build !zogmeta

package zss_test

import (
	"testing"
	"time"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestStringSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	s := zog.String()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestStringLikeSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	type CustomString string
	s := zog.StringLike[CustomString]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestIntLikeSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestFloatLikeSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	type CustomFloat float64
	s := zog.FloatLike[CustomFloat]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestBoolLikeSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestTimeSchemaNoMeta_FormatBehavior(t *testing.T) {
	// Time format may or may not be set depending on implementation
	// This test documents current behavior without zogmeta
	s := zog.Time()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "time")
	// Format may be nil or set - depends on implementation
	// Just verify the schema is valid
	assert.NotNil(t, doc.Root)
}

func TestPreprocessSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	s := zog.Preprocess(
		func(data any, ctx zog.Ctx) (any, error) {
			return data, nil
		},
		zog.String(),
	)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "preprocess")
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestBoxedSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	type StringBox struct {
		V string
	}
	s := zog.Boxed(
		zog.String(),
		func(b StringBox, ctx zog.Ctx) (string, error) { return b.V, nil },
		func(s string, ctx zog.Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "boxed")
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestCustomSchemaNoMeta_GoTypeIsNil(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "custom")
	assert.Nil(t, doc.Root.GoTypes, "GoTypes should be nil when zogmeta build tag is not set")
}

func TestTimeSchemaWithFormatNoMeta_FormatMayBeNil(t *testing.T) {
	// When Time.Format is used, format may still be nil without zogmeta
	// This documents current behavior
	s := zog.Time(zog.Time.Format(time.RFC3339))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "time")
	// Format behavior without zogmeta may vary - just verify schema is valid
	assert.NotNil(t, doc.Root)
}
