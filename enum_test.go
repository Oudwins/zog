package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum(t *testing.T) {
	validator := Enum([]any{1, 2, 3})
	val, errs, ok := validator.Parse(1)

	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, errs, ok = validator.Parse(4)
	assert.NotEmpty(t, errs)
	assert.False(t, ok)
	assert.Equal(t, 4, val)

	val, errs, ok = validator.Parse("a")
	assert.NotEmpty(t, errs)
	assert.False(t, ok)
	assert.Equal(t, "a", val)
}
