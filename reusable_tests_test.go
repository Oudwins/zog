package zog

import (
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func MyTest(opts ...TestOption) Test[*string] {
	options := []TestOption{
		Message("Default message"),
	}
	options = append(options, opts...)
	return TestFunc(zconst.ZogIssueCode("customTest"), func(val *string, ctx Ctx) bool {
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

func TestReusableTestWithParams(t *testing.T) {
	s := String().Test(MyTest(
		Params(map[string]any{"min": 3}),
	))
	v := "ab"
	errs := s.Validate(&v)
	assert.Equal(t, map[string]any{"min": 3}, errs[0].Params)
}

func TestReusableTestWithPath(t *testing.T) {
	s := String().Test(MyTest(
		IssuePath([]string{"user", "name"}),
	))
	v := "test"
	errs := s.Validate(&v)
	assert.Equal(t, "user.name", errs[0].Path)
}

func TestReusableTestWithMultipleOptions(t *testing.T) {
	s := String().Test(MyTest(
		Message("custom message"),
		IssueCode("custom_code"),
		Params(map[string]any{"key": "value"}),
		IssuePath([]string{"field", "path"}),
	))
	v := "test"
	errs := s.Validate(&v)

	assert.Equal(t, "custom message", errs[0].Message)
	assert.Equal(t, "custom_code", errs[0].Code)
	assert.Equal(t, map[string]any{"key": "value"}, errs[0].Params)
	assert.Equal(t, "field.path", errs[0].Path)
}
