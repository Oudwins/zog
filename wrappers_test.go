package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDefault(t *testing.T) {
	field := String().Email().Default("foo@bar.com")
	val, errs, ok := field.Parse("")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "foo@bar.com", val)
}

func TestOptional(t *testing.T) {
	field := String().Email().Optional()
	_, errs, ok := field.Parse("")
	assert.Empty(t, errs)
	assert.True(t, ok)
}

func TestCatch(t *testing.T) {
	field := String().Email().Catch("foo@bar.com")
	val, errs, ok := field.Parse("")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "foo@bar.com", val)
}

func TestTransform(t *testing.T) {
	field := String().Transform(func(val any) (any, bool) {
		return val.(string) + "foo", true
	})

	val, errs, ok := field.Parse("bar")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "barfoo", val)
}

func TestDefaultCatch(t *testing.T) {
	field := String().Email().Default("foo@bar.com").Catch("bar@baz.com")
	val, errs, ok := field.Parse("")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "foo@bar.com", val)
	val, errs, ok = field.Parse("bar")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "bar@baz.com", val)
}

func TestOptionalCatch(t *testing.T) {
	field := String().Email().Optional().Catch("bar@baz.com")
	val, errs, ok := field.Parse("")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "", val)
	val, errs, ok = field.Parse("bar")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "bar@baz.com", val)
}
