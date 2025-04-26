package zog

import (
	"fmt"
	"testing"

	"github.com/Oudwins/zog/conf"
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
	WithIssueFormatter(func(e *p.ZogIssue, p Ctx) {
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
	schema := Struct(Shape{
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
