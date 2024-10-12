package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZeroStringEmpty(t *testing.T) {
	assert.True(t, DefaultParseIsZeroValue.String(""))
	assert.True(t, DefaultParseIsZeroValue.String("    "))
	assert.False(t, DefaultParseIsZeroValue.String("string"))
}
