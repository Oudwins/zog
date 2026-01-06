//go:build !zogmeta
// +build !zogmeta

package zss_test

import (
	"testing"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestStructShapeBasic(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String().Required(),
		"age":  zog.Int().Optional(),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "struct")

	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		assert.Len(t, childShape, 2, "shape should have 2 fields")

		// Check name field
		nameSchema, nameExists := childShape["name"]
		if assert.True(t, nameExists) {
			assertSchemaKind(t, &nameSchema, "string")
			assertRequired(t, &nameSchema, true)
		}

		// Check age field
		ageSchema, ageExists := childShape["age"]
		if assert.True(t, ageExists) {
			assertSchemaKind(t, &ageSchema, "number")
			assertRequired(t, &ageSchema, false)
		}
	}
}

func TestStructShapeWithDefaultsAndCatch(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String().Default("John"),
		"age":  zog.Int().Catch(0),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		nameSchema := childShape["name"]
		assertDefaultValue(t, &nameSchema, "John")

		ageSchema := childShape["age"]
		assertCatchValue(t, &ageSchema, 0)
	}
}

func TestStructShapeWithProcessors(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"email": zog.String().Email().Min(5),
		"count": zog.Int().GT(0).LT(100),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		emailSchema := childShape["email"]
		assertProcessorsCount(t, &emailSchema, 2)
		assertTestProcessor(t, emailSchema.Processors, 0, "email", map[string]any{})
		assertTestProcessor(t, emailSchema.Processors, 1, "min", map[string]any{"min": 5})

		countSchema := childShape["count"]
		assertProcessorsCount(t, &countSchema, 2)
		assertTestProcessor(t, countSchema.Processors, 0, "gt", map[string]any{"gt": 0})
		assertTestProcessor(t, countSchema.Processors, 1, "lt", map[string]any{"lt": 100})
	}
}

func TestStructShapeNestedPtr(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"user": zog.Ptr(zog.Struct(zog.Shape{
			"name": zog.String().Required(),
		})),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		userSchema := childShape["user"]
		assertSchemaKind(t, &userSchema, "ptr")

		ptrChild, ok := assertChildIsSchema(t, &userSchema)
		if assert.True(t, ok) {
			assertSchemaKind(t, ptrChild, "struct")

			userShape, ok := assertChildIsShape(t, ptrChild)
			if assert.True(t, ok) {
				nameSchema := userShape["name"]
				assertSchemaKind(t, &nameSchema, "string")
				assertRequired(t, &nameSchema, true)
			}
		}
	}
}

func TestStructShapeNestedSlice(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"tags": zog.Slice(zog.String().Min(1)).Min(1),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		tagsSchema := childShape["tags"]
		assertSchemaKind(t, &tagsSchema, "slice")
		assertProcessorsCount(t, &tagsSchema, 1)
		assertTestProcessor(t, tagsSchema.Processors, 0, "min", map[string]any{"min": 1})

		sliceChild, ok := assertChildIsSchema(t, &tagsSchema)
		if assert.True(t, ok) {
			assertSchemaKind(t, sliceChild, "string")
			assertProcessorsCount(t, sliceChild, 1)
			assertTestProcessor(t, sliceChild.Processors, 0, "min", map[string]any{"min": 1})
		}
	}
}

func TestStructShapeDeeplyNested(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"user": zog.Ptr(zog.Struct(zog.Shape{
			"profile": zog.Struct(zog.Shape{
				"bio": zog.String().Max(500),
			}),
		})),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		userSchema := childShape["user"]
		assertSchemaKind(t, &userSchema, "ptr")

		ptrChild, ok := assertChildIsSchema(t, &userSchema)
		if assert.True(t, ok) {
			assertSchemaKind(t, ptrChild, "struct")

			userShape, ok := assertChildIsShape(t, ptrChild)
			if assert.True(t, ok) {
				profileSchema := userShape["profile"]
				assertSchemaKind(t, &profileSchema, "struct")

				profileShape, ok := assertChildIsShape(t, &profileSchema)
				if assert.True(t, ok) {
					bioSchema := profileShape["bio"]
					assertSchemaKind(t, &bioSchema, "string")
					assertProcessorsCount(t, &bioSchema, 1)
					assertTestProcessor(t, bioSchema.Processors, 0, "max", map[string]any{"max": 500})
				}
			}
		}
	}
}

func TestStructShapeWithMultipleComplexTypes(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name":   zog.String().Required(),
		"emails": zog.Slice(zog.String().Email()).Min(1),
		"meta": zog.Ptr(zog.Struct(zog.Shape{
			"key": zog.String(),
		})),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		assert.Len(t, childShape, 3)

		// Check name (primitive)
		nameSchema := childShape["name"]
		assertSchemaKind(t, &nameSchema, "string")
		assertRequired(t, &nameSchema, true)

		// Check emails (slice)
		emailsSchema := childShape["emails"]
		assertSchemaKind(t, &emailsSchema, "slice")
		assertProcessorsCount(t, &emailsSchema, 1)
		assertTestProcessor(t, emailsSchema.Processors, 0, "min", map[string]any{"min": 1})

		// Check meta (ptr to struct)
		metaSchema := childShape["meta"]
		assertSchemaKind(t, &metaSchema, "ptr")
		metaPtrChild, ok := assertChildIsSchema(t, &metaSchema)
		if assert.True(t, ok) {
			assertSchemaKind(t, metaPtrChild, "struct")
			metaShape, ok := assertChildIsShape(t, metaPtrChild)
			if assert.True(t, ok) {
				keySchema := metaShape["key"]
				assertSchemaKind(t, &keySchema, "string")
			}
		}
	}
}

func TestStructShapeEmpty(t *testing.T) {
	s := zog.Struct(zog.Shape{})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Root, "struct")
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		assert.Len(t, childShape, 0, "empty shape should have no fields")
	}
}

func TestStructShapeWithTransforms(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String().Trim().Min(1),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		nameSchema := childShape["name"]
		assertProcessorsCount(t, &nameSchema, 2)
		// Trim() ID is only set when zogmeta build tag is enabled
		// Without zogmeta, it defaults to "custom"
		// This test runs without zogmeta, so expect "custom"
		assertTransformProcessor(t, nameSchema.Processors, 0, "custom")
		assertTestProcessor(t, nameSchema.Processors, 1, "min", map[string]any{"min": 1})
	}
}

func TestStructShapeWithCustomMessages(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"email": zog.String().Email(zog.Message("Invalid email address")),
		"age":   zog.Int().GT(0, zog.Message("Age must be positive")),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	childShape, ok := assertChildIsShape(t, doc.Root)
	if assert.True(t, ok) {
		emailSchema := childShape["email"]
		if assert.NotNil(t, emailSchema.Processors) && len(emailSchema.Processors) > 0 {
			if assert.NotNil(t, emailSchema.Processors[0].Test) {
				assert.Equal(t, "Invalid email address", emailSchema.Processors[0].Test.Message)
			}
		}

		ageSchema := childShape["age"]
		if assert.NotNil(t, ageSchema.Processors) && len(ageSchema.Processors) > 0 {
			if assert.NotNil(t, ageSchema.Processors[0].Test) {
				assert.Equal(t, "Age must be positive", ageSchema.Processors[0].Test.Message)
			}
		}
	}
}
