package zog

import (
	"testing"

	p "github.com/Oudwins/zog/primitives"
	"github.com/stretchr/testify/assert"
)

func TestStringSchemaOptionalByDefault(t *testing.T) {
	field := String().Len(3).Contains("foo").HasPrefix("pre").HasSuffix("fix")
	var dest string

	errs := field.Parse("", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	field = field.Required()

	errs = field.Parse("", &dest)
	assert.NotEmpty(t, errs)

}

func TestStringPreTransform(t *testing.T) {
	field := String().Required().Len(3).PreTransform(func(val any, ctx ParseCtx) (any, error) {
		return "foo", nil
	})
	var dest string

	errs := field.Parse("", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "foo", dest)
}

func TestStringRequiredAborts(t *testing.T) {
	field := String().Required().Len(3)
	var dest string

	errs := field.Parse("", &dest)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 1)
}

func TestStringUserTests(t *testing.T) {

	field := String().Test("test", func(val any, ctx p.ParseCtx) bool {
		return val == "test"
	}, Message("Invalid"))

	var dest string

	errs := field.Parse("test", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test", dest)

	errs = field.Parse("not test", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Invalid", errs[0].Message())

}

func TestStringMessage(t *testing.T) {
	field := String().Min(5, Message("min")).Email(Message("email"))
	var dest string
	errs := field.Parse("x", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "min", errs[0].Message())
	assert.Equal(t, "email", errs[1].Message())
}

func TestStringLength(t *testing.T) {
	field := String().Len(3)
	var dest string

	errs := field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)

	field = String().Min(5).Max(7)
	errs = field.Parse("123456789", &dest)
	assert.NotEmpty(t, errs)

	assert.Equal(t, "123456789", dest)

	field = String().Min(5).Max(7)
	errs = field.Parse("1234567", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234567", dest)
}

func TestStringRequired(t *testing.T) {
	field := String().Required()
	var dest string

	errs := field.Parse("", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestStringOptional(t *testing.T) {
	field := String().Optional()
	var dest string

	errs := field.Parse("", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	errs = field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestStringDefault(t *testing.T) {
	field := String().Default("bar")
	var dest string

	errs := field.Parse("", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "bar", dest)

	errs = field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestStringCatch(t *testing.T) {
	field := String().Required().Min(5).Catch("error")
	var dest string

	errs := field.Parse("x", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "error", dest)

	errs = field.Parse("not error", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "not error", dest)
}

func TestStringEmail(t *testing.T) {
	field := String().Email()
	var dest string

	errs := field.Parse("not an email", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("test@example.com", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test@example.com", dest)
}

func TestStringURL(t *testing.T) {
	field := String().URL()
	var dest string

	errs := field.Parse("not a url", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("http://example.com", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "http://example.com", dest)
}

func TestStringHasPrefix(t *testing.T) {
	field := String().HasPrefix("pre")
	var dest string

	errs := field.Parse("not prefixed", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("prefix", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "prefix", dest)
}

func TestStringHasPostfix(t *testing.T) {
	field := String().HasSuffix("fix")
	var dest string

	errs := field.Parse("not postfixed", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("postfix", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "postfix", dest)
}

func TestStringContains(t *testing.T) {
	field := String().Contains("contains")
	var dest string

	errs := field.Parse("does not contain", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("contains", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "contains", dest)
}

func TestStringContainsDigit(t *testing.T) {
	field := String().ContainsDigit()
	var dest string

	errs := field.Parse("no digit here", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("1234", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234", dest)
}

func TestStringContainsUpper(t *testing.T) {
	field := String().ContainsUpper()
	var dest string

	errs := field.Parse("no uppercase here", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("UPPERCASE", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "UPPERCASE", dest)
}

func TestStringContainsSpecial(t *testing.T) {
	field := String().ContainsSpecial()
	var dest string

	errs := field.Parse("no special character here", &dest)
	assert.NotEmpty(t, errs)

	errs = field.Parse("!@#$%", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "!@#$%", dest)
}
