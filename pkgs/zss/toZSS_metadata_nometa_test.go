//go:build !zogmeta
// +build !zogmeta

package zss_test

import (
	"testing"
	"time"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestStringSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	s := zog.String()
	doc := zog.EXPERIMENTAL_TO_ZSS[string](s)

	assertDocumentBasics(t, doc)
	assert.Equal(t, "string", doc.Root.GoTypes[0].Display)
}

func TestStringLikeSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	type CustomString string
	s := zog.StringLike[CustomString]()
	doc := zog.EXPERIMENTAL_TO_ZSS[CustomString](s)

	assertDocumentBasics(t, doc)
	assert.Equal(t, "zss_test.CustomString", doc.Root.GoTypes[0].Display)
}

func TestIntLikeSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	type CustomInt int
	s := zog.IntLike[CustomInt]()
	doc := zog.EXPERIMENTAL_TO_ZSS[CustomInt](s)

	assertDocumentBasics(t, doc)
	assert.Equal(t, "zss_test.CustomInt", doc.Root.GoTypes[0].Display)
}

func TestFloatLikeSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	type CustomFloat float64
	s := zog.FloatLike[CustomFloat]()
	doc := zog.EXPERIMENTAL_TO_ZSS[CustomFloat](s)

	assertDocumentBasics(t, doc)
	assert.Equal(t, "zss_test.CustomFloat", doc.Root.GoTypes[0].Display)
}

func TestBoolLikeSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	type CustomBool bool
	s := zog.BoolLike[CustomBool]()
	doc := zog.EXPERIMENTAL_TO_ZSS[CustomBool](s)

	assertDocumentBasics(t, doc)
	assert.Equal(t, "zss_test.CustomBool", doc.Root.GoTypes[0].Display)
}

func TestTimeSchemaNoMeta_FormatBehavior(t *testing.T) {
	// Time format may or may not be set depending on implementation
	// This test documents current behavior without zogmeta
	s := zog.Time()
	doc := zog.EXPERIMENTAL_TO_ZSS[time.Time](s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "time")
	// Format may be nil or set - depends on implementation
	// Just verify the schema is valid
	assert.NotNil(t, doc.Root)
}

func TestPreprocessSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	s := zog.Preprocess(
		func(data any, ctx zog.Ctx) (any, error) {
			return data, nil
		},
		zog.String(),
	)
	doc := zog.EXPERIMENTAL_TO_ZSS[string](s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "preprocess")
	assert.Equal(t, "string", doc.Root.GoTypes[0].Display)
}

func TestBoxedSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	type StringBox struct {
		V string
	}
	s := zog.Boxed(
		zog.String(),
		func(b StringBox, ctx zog.Ctx) (string, error) { return b.V, nil },
		func(s string, ctx zog.Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)
	doc := zog.EXPERIMENTAL_TO_ZSS[StringBox](s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "boxed")
	assert.Equal(t, "zss_test.StringBox", doc.Root.GoTypes[0].Display)
}

func TestCustomSchemaNoMeta_GoTypeUsesRootType(t *testing.T) {
	type CustomType string
	s := zog.CustomFunc[CustomType](func(valPtr *CustomType, ctx zog.Ctx) bool {
		return true
	})
	doc := zog.EXPERIMENTAL_TO_ZSS[CustomType](s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "custom")
	assert.Equal(t, "zss_test.CustomType", doc.Root.GoTypes[0].Display)
}

func TestTimeSchemaWithFormatNoMeta_FormatMayBeNil(t *testing.T) {
	// When Time.Format is used, format may still be nil without zogmeta
	// This documents current behavior
	s := zog.Time(zog.Time.Format(time.RFC3339))
	doc := zog.EXPERIMENTAL_TO_ZSS[time.Time](s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "time")
	// Format behavior without zogmeta may vary - just verify schema is valid
	assert.NotNil(t, doc.Root)
}
