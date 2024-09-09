package zog

import (
	"testing"

	p "github.com/Oudwins/zog/primitives"
	"github.com/stretchr/testify/assert"
)

func TestWithCtxValue(t *testing.T) {
	var ctx = p.NewParseCtx(nil, nil)
	WithCtxValue("foo", "bar")(ctx)
	assert.Equal(t, "bar", ctx.Get("foo"))
}

func TestWithErrFormatter(t *testing.T) {
	var ctx = p.NewParseCtx(p.NewErrsList(), nil)
	WithErrFormatter(func(e p.ZogError, p p.ParseCtx) {
		e.SetMessage("foo")
	})(ctx)

	err := &p.ZogErr{}
	ctx.NewError(p.PathBuilder(""), err)
	assert.Equal(t, "foo", err.Message())
}
