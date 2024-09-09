package conf

import (
	"testing"

	p "github.com/Oudwins/zog/primitives"
	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorFormatter(t *testing.T) {
	tests := []struct {
		input p.ZogError
		want  string
	}{
		{input: &p.ZogErr{C: p.ErrCodeRequired, Typ: p.TypeString}, want: DefaultErrMsgMap[p.TypeString][p.ErrCodeRequired]},
		{input: &p.ZogErr{C: p.ErrCodeRequired, Typ: p.TypeString, Msg: "DON'T OVERRIDE ME"}, want: "DON'T OVERRIDE ME"},
		{input: &p.ZogErr{C: "INVALID_ERR_CODE", Typ: p.TypeString}, want: "string is invalid"},
	}

	for _, test := range tests {
		ErrorFormatter(test.input, nil)
		assert.Equal(t, test.want, test.input.Message())
	}
}
