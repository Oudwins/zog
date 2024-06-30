package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	field := String().Len(3)

	val, errs, ok := field.Parse("foo")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "foo", val)

	field = String().Min(5).Max(7)
	val, errs, ok = field.Parse("123456789")
	assert.NotEmpty(t, errs)
	assert.False(t, ok)
	assert.Equal(t, "123456789", val)

	field = String().Min(5).Max(7)
	val, errs, ok = field.Parse("1234567")
	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, "1234567", val)
}
