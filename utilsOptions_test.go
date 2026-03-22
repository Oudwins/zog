package zog

import (
	"fmt"
	"testing"

	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/pkgs/internals/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestWithCtxValue(t *testing.T) {
	var ctx = p.NewExecCtx(nil, nil)
	WithCtxValue("foo", "bar")(ctx)
	assert.Equal(t, "bar", ctx.Get("foo"))
}

func TestWithCtxValueNotPersistedAcrossExecCtxReuse(t *testing.T) {
	ctx1 := p.NewExecCtx(nil, nil)
	ctx1.Set("secret", "should_not_leak")
	assert.Equal(t, "should_not_leak", ctx1.Get("secret"))
	ctx1.Free()

	// Prove that there is no longer a leak when pulling execution context out of
	// the pool. Previously this would not erase the map of data stored on the
	// context. It was possible to store a value on a context, free that context,
	// then see that value appear on subsequent context usages elsewhere. This
	// test should prove that that won't happen.
	ctx2 := p.NewExecCtx(nil, nil)
	assert.Nil(t, ctx2.Get("secret"))
	ctx2.Free()
}

func TestWithIssueFormatter(t *testing.T) {
	var ctx = p.NewExecCtx(p.NewErrsList(), nil)
	WithIssueFormatter(func(e *p.ZogIssue, p Ctx) {
		e.SetMessage("foo")
	})(ctx)

	err := &p.ZogIssue{
		Path: nil,
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
	schema := Struct(Shape{
		"name": String().Min(5, IssuePath([]string{"foo"}), Message("foo msg")),
	})

	// Test Parse
	err := schema.Parse(map[string]any{
		"name": "1234",
	}, &out)
	fooErrs := tutils.FindByPath(err, "foo")
	assert.NotEmpty(t, fooErrs)
	assert.Equal(t, "foo msg", fooErrs[0].Message)
	assert.Equal(t, "foo", p.FlattenPath(fooErrs[0].Path))

	// Test Validate
	out.Name = "1234"
	err = schema.Validate(&out)
	fooErrs = tutils.FindByPath(err, "foo")
	assert.NotEmpty(t, fooErrs)
	assert.Equal(t, "foo msg", fooErrs[0].Message)
	assert.Equal(t, "foo", p.FlattenPath(fooErrs[0].Path))
}

func TestWithCoercer(t *testing.T) {
	schema := String(WithCoercer(func(val any) (any, error) {
		switch v := val.(type) {
		case [32]byte:
			// Convert bytes to a proper string representation
			result := ""
			for _, b := range v {
				if b == 0 {
					break // Stop at null byte
				}
				result += fmt.Sprintf("%02x", b)
			}
			return result, nil
		default:
			// Fall back to default string coercer for other types
			return conf.DefaultCoercers.String(val)
		}
	}))
	var out string
	err := schema.Parse([32]byte{1, 2, 3}, &out)
	assert.Empty(t, err)
	assert.Equal(t, "010203", out)

	type S struct {
		Bytes string
	}

	schema2 := Struct(Shape{
		"bytes": schema,
	})

	var out2 S
	schema2.Parse(map[string]any{
		"bytes": [32]byte{1, 2, 3},
	}, &out2)
	assert.Empty(t, err)
	assert.Equal(t, "010203", out2.Bytes)
}
