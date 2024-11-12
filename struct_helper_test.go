package zog

import (
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/stretchr/testify/assert"
)

func TestStructMergeSimple(t *testing.T) {
	var nameSchema = Struct(Schema{
		"name": String().Contains("hello").Required(),
	})
	var ageSchema = Struct(Schema{
		"age": Int().GT(18).Required(),
	})
	var schema = nameSchema.Merge(ageSchema)

	type User struct {
		Name string
		Age  int
	}

	var o User
	errs := schema.Parse(NewMapDataProvider(map[string]any{"name": "hello", "age": 20}), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Age, 20)
}

func TestStructMergeOverride(t *testing.T) {
	var nameSchema = Struct(Schema{
		"name": String().Contains("hello").Required(),
	})
	var ageSchema = Struct(Schema{
		"name": String().Contains("world").Required(),
	})
	var schema = nameSchema.Merge(ageSchema)

	type User struct {
		Name string
	}

	var o User
	errs := schema.Parse(NewMapDataProvider(map[string]any{"name": "world"}), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "world")
}

func TestStructMergeWithPostTransforms(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	var nameSchema = Struct(Schema{
		"name": String().Contains("hello").Required(),
	}).PostTransform(func(data any, ctx p.ParseCtx) error {
		u := data.(*User)
		u.Name = u.Name + "_post"
		return nil
	})
	var ageSchema = Struct(Schema{
		"age": Int().GT(18).Required(),
	}).PostTransform(func(data any, ctx p.ParseCtx) error {
		u := data.(*User)
		u.Age = u.Age + 10
		return nil
	})
	var schema = nameSchema.Merge(ageSchema)

	var o User

	errs := schema.Parse(NewMapDataProvider(map[string]any{"name": "hello", "age": 20}), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello_post")
	assert.Equal(t, o.Age, 30)
}

func TestStructMergeWithPreTransforms(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	var nameSchema = Struct(Schema{
		"name": String().Contains("hello").Required(),
	}).PreTransform(func(data any, ctx p.ParseCtx) (any, error) {
		m := data.(map[string]any)
		m["name"] = m["name"].(string) + "_pre"
		return m, nil
	})
	var ageSchema = Struct(Schema{
		"age": Int().GT(18).Required(),
	}).PreTransform(func(data any, ctx p.ParseCtx) (any, error) {
		m := data.(map[string]any)
		m["age"] = m["age"].(int) + 5
		return m, nil
	})
	var schema = nameSchema.Merge(ageSchema)

	var o User

	errs := schema.Parse(map[string]any{"name": "hello", "age": 20}, &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello_pre")
	assert.Equal(t, o.Age, 25)
}

func TestStructMergeMultiple(t *testing.T) {
	var nameSchema = Struct(Schema{
		"name": String().Contains("hello").Required(),
	})
	var ageSchema = Struct(Schema{
		"age": Int().GT(18).Required(),
	})
	var emailSchema = Struct(Schema{
		"email": String().Email().Required(),
	})
	var schema = nameSchema.Merge(ageSchema, emailSchema)

	type User struct {
		Name  string
		Age   int
		Email string
	}

	var o User
	errs := schema.Parse(NewMapDataProvider(map[string]any{"name": "hello", "age": 20, "email": "test@test.com"}), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Age, 20)
	assert.Equal(t, o.Email, "test@test.com")
}

func TestStructPick(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	var schema = Struct(Schema{
		"name":  String().Contains("hello").Required(),
		"age":   Int().GT(18).Required(),
		"email": String().Email().Required(),
	})

	pickedSchema := schema.Pick("name", map[string]bool{
		"email": true,
	})

	var o User
	errs := pickedSchema.Parse(NewMapDataProvider(map[string]any{
		"name":  "hello",
		"email": "test@test.com",
		"age":   20, // This should be ignored
	}), &o)

	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Email, "test@test.com")
	assert.Equal(t, o.Age, 0) // Age should be zero since it was not picked
}

func TestStructOmit(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	var schema = Struct(Schema{
		"name":  String().Contains("hello").Required(),
		"age":   Int().GT(18).Required(),
		"email": String().Email().Required(),
	})

	omittedSchema := schema.Omit(map[string]bool{
		"age": true,
	})

	var o User
	errs := omittedSchema.Parse(NewMapDataProvider(map[string]any{
		"name":  "hello",
		"email": "test@test.com",
		"age":   20, // This should be ignored
	}), &o)

	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Email, "test@test.com")
	assert.Equal(t, o.Age, 0) // Age should be zero since it was omitted
}

func TestStructPickWithTransforms(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	var schema = Struct(Schema{
		"name": String().Contains("hello").Required(),
		"age":  Int().GT(18).Required(),
	}).PostTransform(func(data any, ctx p.ParseCtx) error {
		u := data.(*User)
		u.Name = u.Name + "_post"
		return nil
	})

	pickedSchema := schema.Pick(map[string]bool{
		"name": true,
	})

	var o User
	errs := pickedSchema.Parse(NewMapDataProvider(map[string]any{
		"name": "hello",
		"age":  20, // This should be ignored
	}), &o)

	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello_post") // Transform should still work
	assert.Equal(t, o.Age, 0)             // Age should be zero since it was not picked
}

func TestStructOmitWithTransforms(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	var schema = Struct(Schema{
		"name": String().Contains("hello").Required(),
		"age":  Int().GT(18).Required(),
	}).PostTransform(func(data any, ctx p.ParseCtx) error {
		u := data.(*User)
		u.Name = u.Name + "_post"
		return nil
	})

	omittedSchema := schema.Omit(map[string]bool{
		"age": true,
	})

	var o User
	errs := omittedSchema.Parse(NewMapDataProvider(map[string]any{
		"name": "hello",
		"age":  20, // This should be ignored
	}), &o)

	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello_post") // Transform should still work
	assert.Equal(t, o.Age, 0)             // Age should be zero since it was omitted
}

func TestStructPickIgnoresFalseKeys(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	var schema = Struct(Schema{
		"name":  String().Contains("hello").Required(),
		"age":   Int().GT(18).Required(),
		"email": String().Email().Required(),
	})

	pickedSchema := schema.Pick(map[string]bool{
		"name":  true,
		"email": false, // Should be ignored (not picked)
		"age":   true,
	})

	var o User
	errs := pickedSchema.Parse(NewMapDataProvider(map[string]any{
		"name":  "hello",
		"age":   20,
		"email": "test@test.com", // Should be ignored since email was not picked
	}), &o)

	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Age, 20)
	assert.Equal(t, o.Email, "") // Email should be empty since it was not picked (false)
}

func TestStructOmitIgnoresFalseKeys(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	var schema = Struct(Schema{
		"name":  String().Contains("hello").Required(),
		"age":   Int().GT(18).Required(),
		"email": String().Email().Required(),
	})

	omittedSchema := schema.Omit("age", map[string]bool{
		"email": false, // Should be ignored (not omitted)
	})

	var o User
	errs := omittedSchema.Parse(NewMapDataProvider(map[string]any{
		"name":  "hello",
		"email": "test@test.com", // Should still be processed since email omit was false
		"age":   20,              // Should be ignored since age was omitted
	}), &o)

	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Email, "test@test.com") // Email should be processed since omit was false
	assert.Equal(t, o.Age, 0)                 // Age should be zero since it was omitted
}
