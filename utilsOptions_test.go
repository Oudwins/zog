package zog

import (
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestWithCtxValue(t *testing.T) {
	var ctx = p.NewExecCtx(nil, nil)
	WithCtxValue("foo", "bar")(ctx)
	assert.Equal(t, "bar", ctx.Get("foo"))
}

func TestWithIssueFormatter(t *testing.T) {
	var ctx = p.NewExecCtx(p.NewErrsList(), nil)
	WithIssueFormatter(func(e *p.ZogIssue, p ParseCtx) {
		e.SetMessage("foo")
	})(ctx)

	err := &p.ZogIssue{
		Path: "",
	}
	ctx.AddIssue(err)
	assert.Equal(t, "foo", err.Message)
}

func TestWithMessageFunc(t *testing.T) {
	var out string
	err := String().Min(5, MessageFunc(func(e *p.ZogIssue, ctx Ctx) {
		e.SetMessage("HELLO WORLD")
	})).Parse("1234", &out)

	assert.Equal(t, "HELLO WORLD", err[0].Message)
}

func TestIssueCode(t *testing.T) {
	var out string
	schema := String().Min(5, IssueCode(zconst.IssueCodeCustom))

	// Test Parse
	err := schema.Parse("1234", &out)
	assert.Equal(t, zconst.IssueCodeCustom, err[0].Code)
	tutils.VerifyDefaultIssueMessages(t, err)

	// Test Validate
	out = "1234"
	err = schema.Validate(&out)
	assert.Equal(t, zconst.IssueCodeCustom, err[0].Code)
	tutils.VerifyDefaultIssueMessages(t, err)
}

func TestIssuePath(t *testing.T) {
	type User struct {
		Name string
	}
	var out User
	schema := Struct(Schema{
		"name": String().Min(5, IssuePath("foo"), Message("foo msg")),
	})

	// Test Parse
	err := schema.Parse(map[string]any{
		"name": "1234",
	}, &out)
	assert.NotEmpty(t, err["foo"])
	assert.Equal(t, "foo msg", err["foo"][0].Message)
	assert.Equal(t, "foo", err["foo"][0].Path)

	// Test Validate
	out.Name = "1234"
	err = schema.Validate(&out)
	assert.NotEmpty(t, err["foo"])
	assert.Equal(t, "foo msg", err["foo"][0].Message)
	assert.Equal(t, "foo", err["foo"][0].Path)
}
