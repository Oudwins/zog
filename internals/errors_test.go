package internals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZogErrorString(t *testing.T) {
	err := ZogErr{
		C:       "test",
		ParamsM: map[string]any{},
		Typ:     "string",
		Val:     "asda",
		Msg:     "asda",
	}

	assert.Equal(t, err.Error(), err.String())
}
