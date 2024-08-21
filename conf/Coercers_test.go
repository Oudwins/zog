package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolCoercer(t *testing.T) {
	var b any
	var err error
	b, err = Coercers.Bool(true)
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool("true")
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool("on")
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool(1)
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool("off")
	assert.False(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool(0)
	assert.False(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool(false)
	assert.False(t, b.(bool))
	assert.Nil(t, err)
}
