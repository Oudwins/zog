package zog

import (
	"errors"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/stretchr/testify/assert"
)

func TestZogErrorString(t *testing.T) {
	var err p.ZogError = &p.ZogErr{
		C:       "test",
		ParamsM: map[string]any{},
		Typ:     "string",
		Val:     "asda",
		Msg:     "asda",
	}

	assert.Equal(t, err.String(), err.Error())
}

func TestZogErrorUnwrap(t *testing.T) {
	var err p.ZogError = &p.ZogErr{
		Err: errors.New("test"),
	}

	assert.Equal(t, err.Unwrap().Error(), "test")
}
