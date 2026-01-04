package zss_test

import (
	"testing"
	"time"

	"github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestStringSchema(t *testing.T) {
	s := zog.String().Required().Default("default").Catch("catch").Min(5)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "string")
	assertRequired(t, doc.Schema, true)
	assertDefaultValue(t, doc.Schema, "default")
	assertCatchValue(t, doc.Schema, "catch")
	assertProcessorsCount(t, doc.Schema, 1)
	assertChildIsNil(t, doc.Schema)
	assertTestProcessor(t, doc.Schema.Processors, 0, "min", map[string]any{"min": 5})
}

func TestStringSchemaOptional(t *testing.T) {
	s := zog.String().Optional().Min(3)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "string")
	assertRequired(t, doc.Schema, false)
	assertProcessorsCount(t, doc.Schema, 1)
	assertTestProcessor(t, doc.Schema.Processors, 0, "min", map[string]any{"min": 3})
}

func TestNumberSchemaInt(t *testing.T) {
	s := zog.Int().Required().Default(42).Catch(100).GT(10).LT(100)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "number")
	assertRequired(t, doc.Schema, true)
	assertDefaultValue(t, doc.Schema, 42)
	assertCatchValue(t, doc.Schema, 100)
	assertProcessorsCount(t, doc.Schema, 2)
	assertChildIsNil(t, doc.Schema)
	assertTestProcessor(t, doc.Schema.Processors, 0, "gt", map[string]any{"gt": 10})
	assertTestProcessor(t, doc.Schema.Processors, 1, "lt", map[string]any{"lt": 100})
}

func TestNumberSchemaInt64(t *testing.T) {
	s := zog.Int64().EQ(100).GTE(50).LTE(200)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "number")
	assertRequired(t, doc.Schema, false)
	assertProcessorsCount(t, doc.Schema, 3)
	assertTestProcessor(t, doc.Schema.Processors, 0, "eq", map[string]any{"eq": int64(100)})
	assertTestProcessor(t, doc.Schema.Processors, 1, "gte", map[string]any{"gte": int64(50)})
	assertTestProcessor(t, doc.Schema.Processors, 2, "lte", map[string]any{"lte": int64(200)})
}

func TestNumberSchemaFloat64(t *testing.T) {
	s := zog.Float64().OneOf([]float64{1.5, 2.5, 3.5})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "number")
	assertProcessorsCount(t, doc.Schema, 1)
	assertTestProcessor(t, doc.Schema.Processors, 0, "one_of_options", map[string]any{"one_of_options": []float64{1.5, 2.5, 3.5}})
}

func TestBoolSchema(t *testing.T) {
	s := zog.Bool().Required().Default(true).True()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "bool")
	assertRequired(t, doc.Schema, true)
	assertDefaultValue(t, doc.Schema, true)
	assertProcessorsCount(t, doc.Schema, 1)
}

func TestBoolSchemaFalse(t *testing.T) {
	s := zog.Bool().False()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "bool")
	assertProcessorsCount(t, doc.Schema, 1)
	// False() uses EQ internally, so it's "eq" with false value
	assertTestProcessor(t, doc.Schema.Processors, 0, "eq", map[string]any{"eq": false})
}

func TestBoolSchemaEQ(t *testing.T) {
	s := zog.Bool().EQ(true)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "bool")
	assertProcessorsCount(t, doc.Schema, 1)
	assertTestProcessor(t, doc.Schema.Processors, 0, "eq", map[string]any{"eq": true})
}

func TestTimeSchema(t *testing.T) {
	now := time.Now()
	s := zog.Time().Required().Before(now).After(now.Add(-24 * time.Hour))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "time")
	assertRequired(t, doc.Schema, true)
	assertProcessorsCount(t, doc.Schema, 2)
	assertChildIsNil(t, doc.Schema)
	// Note: params for time tests may contain time.Time values
	assertTestProcessor(t, doc.Schema.Processors, 0, "before", nil) // params may vary
	assertTestProcessor(t, doc.Schema.Processors, 1, "after", nil)  // params may vary
}

func TestPtrSchema(t *testing.T) {
	s := zog.Ptr(zog.String().Required().Min(1))
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "ptr")
	assertRequired(t, doc.Schema, false)
	assertProcessorsCount(t, doc.Schema, 0)
	assert.Nil(t, doc.Schema.DefaultValue)
	assert.Nil(t, doc.Schema.CatchValue)

	childSchema, ok := assertChildIsSchema(t, doc.Schema)
	if assert.True(t, ok, "child should be a schema") {
		assertSchemaKind(t, childSchema, "string")
		assertRequired(t, childSchema, true)
		assertProcessorsCount(t, childSchema, 1)
		assertTestProcessor(t, childSchema.Processors, 0, "min", map[string]any{"min": 1})
	}
}

func TestPtrSchemaNotNil(t *testing.T) {
	s := zog.Ptr(zog.String()).NotNil()
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "ptr")
	assertRequired(t, doc.Schema, true)
	// NotNil() sets required, but the ID is still "required" (not_nil is the issue code in the test itself)
	if assert.NotNil(t, doc.Schema.Required) {
		// The toZSSRequired function converts it to "required" ID
		assert.Equal(t, "required", string(doc.Schema.Required.ID))
	}
}

func TestSliceSchema(t *testing.T) {
	s := zog.Slice(zog.String().Min(1)).Required().Default([]any{"a", "b"}).Min(2).Max(10)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "slice")
	assertRequired(t, doc.Schema, true)
	assertProcessorsCount(t, doc.Schema, 2)
	assertTestProcessor(t, doc.Schema.Processors, 0, "min", map[string]any{"min": 2})
	assertTestProcessor(t, doc.Schema.Processors, 1, "max", map[string]any{"max": 10})

	childSchema, ok := assertChildIsSchema(t, doc.Schema)
	if assert.True(t, ok, "child should be a schema") {
		assertSchemaKind(t, childSchema, "string")
		assertProcessorsCount(t, childSchema, 1)
		assertTestProcessor(t, childSchema.Processors, 0, "min", map[string]any{"min": 1})
	}
}

func TestSliceSchemaContains(t *testing.T) {
	s := zog.Slice(zog.Int()).Contains(42)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "slice")
	assertProcessorsCount(t, doc.Schema, 1)
	assertTestProcessor(t, doc.Schema.Processors, 0, "contained", map[string]any{"contained": 42})
}

func TestSliceSchemaLen(t *testing.T) {
	s := zog.Slice(zog.String()).Len(5)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "slice")
	assertProcessorsCount(t, doc.Schema, 1)
	assertTestProcessor(t, doc.Schema.Processors, 0, "len", map[string]any{"len": 5})
}

func TestStructSchema(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"name": zog.String().Required().Min(1),
		"age":  zog.Int().Required().GT(0),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "struct")
	assertRequired(t, doc.Schema, false)
	assertProcessorsCount(t, doc.Schema, 0)

	childShape, ok := assertChildIsShape(t, doc.Schema)
	if assert.True(t, ok, "child should be a shape map") {
		assert.Len(t, childShape, 2, "shape should have 2 fields")

		nameSchema, nameExists := childShape["name"]
		if assert.True(t, nameExists, "name field should exist") {
			assertSchemaKind(t, &nameSchema, "string")
			assertRequired(t, &nameSchema, true)
			assertProcessorsCount(t, &nameSchema, 1)
			assertTestProcessor(t, nameSchema.Processors, 0, "min", map[string]any{"min": 1})
		}

		ageSchema, ageExists := childShape["age"]
		if assert.True(t, ageExists, "age field should exist") {
			assertSchemaKind(t, &ageSchema, "number")
			assertRequired(t, &ageSchema, true)
			assertProcessorsCount(t, &ageSchema, 1)
			assertTestProcessor(t, ageSchema.Processors, 0, "gt", map[string]any{"gt": 0})
		}
	}
}

func TestStructSchemaNested(t *testing.T) {
	s := zog.Struct(zog.Shape{
		"user": zog.Ptr(zog.Struct(zog.Shape{
			"email": zog.String().Email(),
		})),
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "struct")

	childShape, ok := assertChildIsShape(t, doc.Schema)
	if assert.True(t, ok) {
		userSchema, userExists := childShape["user"]
		if assert.True(t, userExists) {
			assertSchemaKind(t, &userSchema, "ptr")

			ptrChild, ok := assertChildIsSchema(t, &userSchema)
			if assert.True(t, ok) {
				assertSchemaKind(t, ptrChild, "struct")

				userShape, ok := assertChildIsShape(t, ptrChild)
				if assert.True(t, ok) {
					emailSchema, emailExists := userShape["email"]
					if assert.True(t, emailExists) {
						assertSchemaKind(t, &emailSchema, "string")
						assertProcessorsCount(t, &emailSchema, 1)
						assertTestProcessor(t, emailSchema.Processors, 0, "email", map[string]any{})
					}
				}
			}
		}
	}
}

func TestPreprocessSchema(t *testing.T) {
	s := zog.Preprocess(
		func(data any, ctx zog.Ctx) (any, error) {
			return data, nil
		},
		zog.String().Min(1),
	)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "preprocess")
	assertRequired(t, doc.Schema, false)
	assertProcessorsCount(t, doc.Schema, 0)

	childSchema, ok := assertChildIsSchema(t, doc.Schema)
	if assert.True(t, ok) {
		assertSchemaKind(t, childSchema, "string")
		assertProcessorsCount(t, childSchema, 1)
		assertTestProcessor(t, childSchema.Processors, 0, "min", map[string]any{"min": 1})
	}
}

func TestBoxedSchema(t *testing.T) {
	type StringBox struct {
		V string
	}
	s := zog.Boxed(
		zog.String().Min(1),
		func(b StringBox, ctx zog.Ctx) (string, error) { return b.V, nil },
		func(s string, ctx zog.Ctx) (StringBox, error) { return StringBox{V: s}, nil },
	)
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "boxed")
	assertRequired(t, doc.Schema, false)
	assertProcessorsCount(t, doc.Schema, 0)

	childSchema, ok := assertChildIsSchema(t, doc.Schema)
	if assert.True(t, ok) {
		assertSchemaKind(t, childSchema, "string")
		assertProcessorsCount(t, childSchema, 1)
		assertTestProcessor(t, childSchema.Processors, 0, "min", map[string]any{"min": 1})
	}
}

func TestCustomSchema(t *testing.T) {
	s := zog.CustomFunc(func(valPtr *string, ctx zog.Ctx) bool {
		return *valPtr == "valid"
	})
	doc := zog.EXPERIMENTAL_TO_ZSS(s)

	assertDocumentBasics(t, doc)
	assertSchemaKind(t, doc.Schema, "custom")
	assertRequired(t, doc.Schema, true) // CustomFunc creates a required test
	assertProcessorsCount(t, doc.Schema, 0)
	assertChildIsNil(t, doc.Schema)
}
