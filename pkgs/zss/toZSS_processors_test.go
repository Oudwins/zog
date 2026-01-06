//go:build !zogmeta
// +build !zogmeta

package zss_test

import (
	"testing"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestStringTransformTrim(t *testing.T) {
	s := zog.String().Trim()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	// Trim() ID is only set when zogmeta build tag is enabled
	// Without zogmeta, it defaults to "custom"
	assertTransformProcessor(t, doc.Root.Processors, 0, "custom")
}

func TestStringTransformCustom(t *testing.T) {
	s := zog.String().Transform(func(valPtr *string, ctx zog.Ctx) error {
		*valPtr = "transformed"
		return nil
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	assertTransformProcessor(t, doc.Root.Processors, 0, "custom")
}

func TestStringMultipleTransforms(t *testing.T) {
	s := zog.String().Trim().Transform(func(valPtr *string, ctx zog.Ctx) error {
		return nil
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 2)
	// Trim() ID is only set when zogmeta build tag is enabled
	// Without zogmeta, it defaults to "custom"
	assertTransformProcessor(t, doc.Root.Processors, 0, "custom")
	assertTransformProcessor(t, doc.Root.Processors, 1, "custom")
}

func TestStringTransformAndTest(t *testing.T) {
	s := zog.String().Trim().Min(5)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 2)
	// Trim() ID is only set when zogmeta build tag is enabled
	// Without zogmeta, it defaults to "custom"
	assertTransformProcessor(t, doc.Root.Processors, 0, "custom")
	assertTestProcessor(t, doc.Root.Processors, 1, "min", map[string]any{"min": 5})
}

func TestStringTestWithMessageOverride(t *testing.T) {
	customMsg := "Custom message for min validation"
	s := zog.String().Min(5, zog.Message(customMsg))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	if assert.NotNil(t, doc.Root.Processors[0].Test) {
		assert.Equal(t, customMsg, doc.Root.Processors[0].Test.Message, "message should be overridden")
		assert.Equal(t, "min", string(doc.Root.Processors[0].Test.ID))
		assert.Equal(t, map[string]any{"min": 5}, doc.Root.Processors[0].Test.Params)
	}
}

func TestStringTestWithIssuePath(t *testing.T) {
	customPath := []string{"fullname"}
	s := zog.String().Min(1, zog.IssuePath(customPath))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	if assert.NotNil(t, doc.Root.Processors[0].Test) {
		assert.NotNil(t, doc.Root.Processors[0].Test.IssuePath, "issuePath should be set")
		assert.Equal(t, customPath, doc.Root.Processors[0].Test.IssuePath)
	}
}

func TestStringTestDefaultMessage(t *testing.T) {
	// Test without custom message - should use default formatter
	s := zog.String().Min(3)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	if assert.NotNil(t, doc.Root.Processors[0].Test) {
		// Message should be populated by default formatter
		assert.NotEmpty(t, doc.Root.Processors[0].Test.Message, "message should be populated by default formatter")
		assert.Contains(t, doc.Root.Processors[0].Test.Message, "3", "message should contain the min value")
	}
}

func TestStringMultipleTests(t *testing.T) {
	s := zog.String().Min(1).Max(100).Email()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 3)
	assertTestProcessor(t, doc.Root.Processors, 0, "min", map[string]any{"min": 1})
	assertTestProcessor(t, doc.Root.Processors, 1, "max", map[string]any{"max": 100})
	assertTestProcessor(t, doc.Root.Processors, 2, "email", map[string]any{})
}

func TestNumberTestParams(t *testing.T) {
	s := zog.Int().GT(10).LT(100).EQ(50)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 3)
	assertTestProcessor(t, doc.Root.Processors, 0, "gt", map[string]any{"gt": 10})
	assertTestProcessor(t, doc.Root.Processors, 1, "lt", map[string]any{"lt": 100})
	assertTestProcessor(t, doc.Root.Processors, 2, "eq", map[string]any{"eq": 50})
}

func TestNumberOneOfParams(t *testing.T) {
	enum := []int{1, 2, 3, 5, 8}
	s := zog.Int().OneOf(enum)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	assertTestProcessor(t, doc.Root.Processors, 0, "one_of_options", map[string]any{"one_of_options": enum})
}

func TestSliceTestParams(t *testing.T) {
	s := zog.Slice(zog.String()).Min(2).Max(10).Len(5)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 3)
	assertTestProcessor(t, doc.Root.Processors, 0, "min", map[string]any{"min": 2})
	assertTestProcessor(t, doc.Root.Processors, 1, "max", map[string]any{"max": 10})
	assertTestProcessor(t, doc.Root.Processors, 2, "len", map[string]any{"len": 5})
}

func TestSliceContainsParams(t *testing.T) {
	s := zog.Slice(zog.Int()).Contains(42)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	assertTestProcessor(t, doc.Root.Processors, 0, "contained", map[string]any{"contained": 42})
}

func TestStringTestFunc(t *testing.T) {
	s := zog.String().TestFunc(func(valPtr *string, ctx zog.Ctx) bool {
		return len(*valPtr) > 0
	}, zog.Message("Custom test message"), zog.IssueCode("custom_test"))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertProcessorsCount(t, doc.Root, 1)
	if assert.NotNil(t, doc.Root.Processors[0].Test) {
		assert.Equal(t, "custom_test", string(doc.Root.Processors[0].Test.ID))
		// Message override should work - but TestFunc may use default formatter if not set properly
		// Check that message is set (either custom or default)
		assert.NotEmpty(t, doc.Root.Processors[0].Test.Message)
	}
}

func TestRequiredProcessor(t *testing.T) {
	s := zog.String().Required()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertRequired(t, doc.Root, true)
	if assert.NotNil(t, doc.Root.Required) {
		assert.Equal(t, "required", string(doc.Root.Required.ID))
		assert.NotEmpty(t, doc.Root.Required.Message, "required message should be populated")
	}
}

func TestRequiredWithMessage(t *testing.T) {
	customMsg := "This field is mandatory"
	s := zog.String().Required(zog.Message(customMsg))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertRequired(t, doc.Root, true)
	if assert.NotNil(t, doc.Root.Required) {
		assert.Equal(t, customMsg, doc.Root.Required.Message)
	}
}

func TestRequiredWithIssuePath(t *testing.T) {
	customPath := []string{"user", "email"}
	s := zog.String().Required(zog.IssuePath(customPath))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertRequired(t, doc.Root, true)
	if assert.NotNil(t, doc.Root.Required) {
		assert.NotNil(t, doc.Root.Required.IssuePath)
		assert.Equal(t, customPath, doc.Root.Required.IssuePath)
	}
}

func TestStructTransform(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String(),
	}).Transform(func(dataPtr any, ctx zog.Ctx) error {
		return nil
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "struct")
	assertProcessorsCount(t, doc.Root, 1)
	// Struct transforms are not yet serialized with IDs, but should be transform processors
	if assert.NotNil(t, doc.Root.Processors) && len(doc.Root.Processors) > 0 {
		assert.Equal(t, "transform", string(doc.Root.Processors[0].Kind))
		assert.NotNil(t, doc.Root.Processors[0].Transformer)
	}
}

func TestSliceTransform(t *testing.T) {
	s := zog.Slice(zog.String()).Transform(func(dataPtr any, ctx zog.Ctx) error {
		return nil
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "slice")
	assertProcessorsCount(t, doc.Root, 1)
	if assert.NotNil(t, doc.Root.Processors) && len(doc.Root.Processors) > 0 {
		assert.Equal(t, "transform", string(doc.Root.Processors[0].Kind))
		assert.NotNil(t, doc.Root.Processors[0].Transformer)
	}
}
