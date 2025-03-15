package zog

import (
	"errors"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/stretchr/testify/assert"
)

func TestZogIssueString(t *testing.T) {
	var err *p.ZogIssue = &p.ZogIssue{
		Code:    "test",
		Params:  map[string]any{},
		Dtype:   "string",
		Value:   "asda",
		Message: "asda",
	}

	assert.Equal(t, err.String(), err.Error())
}

func TestZogIssueUnwrap(t *testing.T) {
	var err *p.ZogIssue = &p.ZogIssue{
		Err: errors.New("test"),
	}

	assert.Equal(t, err.Unwrap().Error(), "test")
}
