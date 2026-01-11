package zss_test

import (
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/stretchr/testify/assert"
)

// assertChildIsNil asserts that schema.Childs is nil or empty
func assertChildIsNil(t assert.TestingT, schema *zss.ZSSSchema, msgAndArgs ...interface{}) bool {
	if schema.Childs == nil {
		return assert.Nil(t, schema.Childs, msgAndArgs...)
	}
	return assert.Len(t, schema.Childs, 0, msgAndArgs...)
}

// assertChildIsSchema asserts that schema.Childs contains a schema child and returns it
func assertChildIsSchema(t assert.TestingT, schema *zss.ZSSSchema, msgAndArgs ...interface{}) (*zss.ZSSSchema, bool) {
	if len(schema.Childs) == 0 {
		return nil, assert.Fail(t, "Childs is empty, expected at least one child", msgAndArgs...)
	}
	// Find the first child with Kind="schema"
	for _, child := range schema.Childs {
		if child.Kind == zss.ZSSSchemaChildKindSchema {
			if child.Schema == nil {
				return nil, assert.Fail(t, "Child Schema is nil", msgAndArgs...)
			}
			return child.Schema, true
		}
	}
	return nil, assert.Fail(t, "No child with Kind='schema' found in Childs", msgAndArgs...)
}

// assertChildIsShape asserts that schema.Childs contains a shape child and returns it
func assertChildIsShape(t assert.TestingT, schema *zss.ZSSSchema, msgAndArgs ...interface{}) (map[string]zss.ZSSSchema, bool) {
	if len(schema.Childs) == 0 {
		return nil, assert.Fail(t, "Childs is empty, expected at least one child", msgAndArgs...)
	}
	// Find the first child with Kind="shape"
	for _, child := range schema.Childs {
		if child.Kind == zss.ZSSSchemaChildKindShape {
			if child.Shape == nil {
				return nil, assert.Fail(t, "Child Shape is nil", msgAndArgs...)
			}
			return child.Shape, true
		}
	}
	return nil, assert.Fail(t, "No child with Kind='shape' found in Childs", msgAndArgs...)
}

// assertDocumentBasics asserts basic document invariants
func assertDocumentBasics(t assert.TestingT, doc zss.ZSSDocument) bool {
	if !assert.Equal(t, zss.ZSS_VERSION_LATEST, doc.Version, "document $schema should be latest") {
		return false
	}
	if !assert.NotNil(t, doc.Root, "document root should not be nil") {
		return false
	}
	return true
}

// assertSchemaKind asserts the schema kind matches expected
func assertSchemaKind(t assert.TestingT, schema *zss.ZSSSchema, expectedKind string) bool {
	return assert.Equal(t, expectedKind, schema.Kind, "schema kind should match")
}

// assertRequired asserts required field is nil or not nil as expected
func assertRequired(t assert.TestingT, schema *zss.ZSSSchema, shouldBeRequired bool) bool {
	if shouldBeRequired {
		return assert.NotNil(t, schema.Required, "schema should be required")
	}
	return assert.Nil(t, schema.Required, "schema should not be required")
}

// assertDefaultValue asserts default value matches expected
func assertDefaultValue(t assert.TestingT, schema *zss.ZSSSchema, expected any) bool {
	if expected == nil {
		return assert.Nil(t, schema.DefaultValue, "default value should be nil")
	}
	// Handle pointer types (string, int, bool, etc.)
	switch ptr := schema.DefaultValue.(type) {
	case *string:
		if str, ok := expected.(string); ok {
			if !assert.NotNil(t, ptr, "default value pointer should not be nil") {
				return false
			}
			return assert.Equal(t, str, *ptr)
		}
	case *int:
		if i, ok := expected.(int); ok {
			if !assert.NotNil(t, ptr, "default value pointer should not be nil") {
				return false
			}
			return assert.Equal(t, i, *ptr)
		}
	case *bool:
		if b, ok := expected.(bool); ok {
			if !assert.NotNil(t, ptr, "default value pointer should not be nil") {
				return false
			}
			return assert.Equal(t, b, *ptr)
		}
	}
	return assert.Equal(t, expected, schema.DefaultValue, "default value should match")
}

// assertCatchValue asserts catch value matches expected
func assertCatchValue(t assert.TestingT, schema *zss.ZSSSchema, expected any) bool {
	if expected == nil {
		return assert.Nil(t, schema.CatchValue, "catch value should be nil")
	}
	// Handle pointer types (string, int, bool, etc.)
	switch ptr := schema.CatchValue.(type) {
	case *string:
		if str, ok := expected.(string); ok {
			if !assert.NotNil(t, ptr, "catch value pointer should not be nil") {
				return false
			}
			return assert.Equal(t, str, *ptr)
		}
	case *int:
		if i, ok := expected.(int); ok {
			if !assert.NotNil(t, ptr, "catch value pointer should not be nil") {
				return false
			}
			return assert.Equal(t, i, *ptr)
		}
	case *bool:
		if b, ok := expected.(bool); ok {
			if !assert.NotNil(t, ptr, "catch value pointer should not be nil") {
				return false
			}
			return assert.Equal(t, b, *ptr)
		}
	}
	return assert.Equal(t, expected, schema.CatchValue, "catch value should match")
}

// assertProcessorsCount asserts the number of processors
func assertProcessorsCount(t assert.TestingT, schema *zss.ZSSSchema, expectedCount int) bool {
	if expectedCount == 0 {
		return assert.Nil(t, schema.Processors, "processors should be nil when count is 0") ||
			assert.Len(t, schema.Processors, 0, "processors should be empty when count is 0")
	}
	return assert.NotNil(t, schema.Processors, "processors should not be nil") &&
		assert.Len(t, schema.Processors, expectedCount, "processors count should match expected")
}

// assertTestProcessor asserts a processor at index is a test processor with expected properties
func assertTestProcessor(t assert.TestingT, processors []zss.ZSSProcessor, index int, expectedID string, expectedParams map[string]any) bool {
	if !assert.GreaterOrEqual(t, len(processors), index+1, "processors should have enough elements") {
		return false
	}
	proc := processors[index]
	if !assert.Equal(t, "test", string(proc.Kind), "processor kind should be test") {
		return false
	}
	if !assert.NotNil(t, proc.Test, "test processor should have test") {
		return false
	}
	if !assert.Nil(t, proc.Transformer, "test processor should not have transformer") {
		return false
	}
	if !assert.Equal(t, expectedID, string(proc.Test.ID), "test ID should match") {
		return false
	}
	if expectedParams != nil {
		return assert.Equal(t, expectedParams, proc.Test.Params, "test params should match")
	}
	return true
}

// assertTransformProcessor asserts a processor at index is a transform processor with expected ID
func assertTransformProcessor(t assert.TestingT, processors []zss.ZSSProcessor, index int, expectedID string) bool {
	if !assert.GreaterOrEqual(t, len(processors), index+1, "processors should have enough elements") {
		return false
	}
	proc := processors[index]
	if !assert.Equal(t, "transform", string(proc.Kind), "processor kind should be transform") {
		return false
	}
	if !assert.NotNil(t, proc.Transformer, "transform processor should have transformer") {
		return false
	}
	if !assert.Nil(t, proc.Test, "transform processor should not have test") {
		return false
	}
	return assert.Equal(t, expectedID, string(proc.Transformer.ID), "transformer ID should match")
}
