package conf

import (
	"testing"

	p "github.com/Oudwins/zog/primitives"
	zconst "github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorFormatter(t *testing.T) {
	tests := []struct {
		input p.ZogError
		want  string
	}{
		{input: &p.ZogErr{C: zconst.ErrCodeRequired, Typ: zconst.TypeString}, want: DefaultErrMsgMap[zconst.TypeString][zconst.ErrCodeRequired]},
		{input: &p.ZogErr{C: zconst.ErrCodeRequired, Typ: zconst.TypeString, Msg: "DON'T OVERRIDE ME"}, want: "DON'T OVERRIDE ME"},
		{input: &p.ZogErr{C: "INVALID_ERR_CODE", Typ: zconst.TypeString}, want: "string is invalid"},
	}

	for _, test := range tests {
		ErrorFormatter(test.input, nil)
		assert.Equal(t, test.want, test.input.Message())
	}
}
