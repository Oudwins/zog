package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func MyTest(opts ...TestOption) Test {
	options := []TestOption{
		Message("Default message"),
	}
	options = append(options, opts...)
	return TestFunc(zconst.ZogIssueCode("customTest"), func(val any, ctx Ctx) bool {
		return false
	}, options...)
}

func TestReusableTestFunc(t *testing.T) {
	v := "value"
	err := String().Test(MyTest()).Validate(&v)
	assert.Equal(t, "Default message", err[0].Message)
}

func TestReusableTestFuncOverride(t *testing.T) {
	s := String().Test(MyTest(
		IssueCode(zconst.ZogIssueCode("customTest3")), Message("override"),
	))
	v := "test"
	errs := s.Validate(&v)
	assert.Equal(t, "customTest3", errs[0].Code)
	assert.Equal(t, "override", errs[0].Message)
}
