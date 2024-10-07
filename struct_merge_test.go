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
