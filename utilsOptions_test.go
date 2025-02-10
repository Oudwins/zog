package zog

import (
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/stretchr/testify/assert"
)

func TestWithCtxValue(t *testing.T) {
	var ctx = p.NewExecCtx(nil, nil)
	WithCtxValue("foo", "bar")(ctx)
	assert.Equal(t, "bar", ctx.Get("foo"))
}

func TestWithErrFormatter(t *testing.T) {
	var ctx = p.NewExecCtx(p.NewErrsList(), nil)
	WithErrFormatter(func(e p.ZogError, p ParseCtx) {
		e.SetMessage("foo")
	})(ctx)

	err := &p.ZogErr{
		EPath: "",
	}
	ctx.AddIssue(err)
	assert.Equal(t, "foo", err.Message())
}

func TestWithMessageFunc(t *testing.T) {
	var out string
	err := String().Min(5, MessageFunc(func(e p.ZogError, ctx Ctx) {
		e.SetMessage("HELLO WORLD")
	})).Parse("1234", &out)

	assert.Equal(t, "HELLO WORLD", err[0].Message())
}
