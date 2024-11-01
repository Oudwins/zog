package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZeroString(t *testing.T) {
	assert.True(t, DefaultParseIsZeroValue.String(""))
	assert.True(t, DefaultParseIsZeroValue.String("    "))
	assert.False(t, DefaultParseIsZeroValue.String("string"))
	assert.True(t, DefaultParseIsZeroValue.String(nil))
	assert.True(t, DefaultParseIsZeroValue.String(0))
}

func TestIsZeroValueBool(t *testing.T) {
	assert.False(t, DefaultParseIsZeroValue.Bool(false))
	assert.False(t, DefaultParseIsZeroValue.Bool(true))
	assert.True(t, DefaultParseIsZeroValue.Bool(nil))
	assert.True(t, DefaultParseIsZeroValue.Bool(0))
}
