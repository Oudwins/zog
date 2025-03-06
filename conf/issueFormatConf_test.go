package conf

import (
	"testing"

	p "github.com/Oudwins/zog/internals"
	zconst "github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestDefaultIssueFormatter(t *testing.T) {
	tests := []struct {
		input *p.ZogIssue
		want  string
	}{
		{input: &p.ZogIssue{Code: zconst.IssueCodeRequired, Dtype: zconst.TypeString}, want: DefaultErrMsgMap[zconst.TypeString][zconst.IssueCodeRequired]},
		{input: &p.ZogIssue{Code: zconst.IssueCodeRequired, Dtype: zconst.TypeString, Message: "DON'T OVERRIDE ME"}, want: "DON'T OVERRIDE ME"},
		{input: &p.ZogIssue{Code: "INVALID_ERR_CODE", Dtype: zconst.TypeString}, want: "string is invalid"},
	}

	for _, test := range tests {
		IssueFormatter(test.input, nil)
		assert.Equal(t, test.want, test.input.Message)
	}
}
