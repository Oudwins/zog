package zog

import (
	"regexp"
	"testing"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestStringOptionalByDefault(t *testing.T) {
	field := String().Len(3).Contains("foo").HasPrefix("pre").HasSuffix("fix")
	var dest string

	errs := field.Parse("", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	field = field.Required()

	errs = field.Parse("", &dest)
	assert.NotEmpty(t, errs)

	field.Required().Optional()

}

func TestStringOptional(t *testing.T) {
	field := String().Required().Optional()
	var dest string

	errs := field.Parse("", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	errs = field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
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

func TestStringPostTransform(t *testing.T) {
	field := String().Required().PostTransform(func(val any, ctx ParseCtx) error {
		s := val.(*string)
		*s = *s + "_transformed"
		return nil
	})
	var dest string

	errs := field.Parse("hello", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "hello_transformed", dest)

	// Test that PostTransform is not applied when there's an error
	field = String().Required().Len(1).PostTransform(func(val any, ctx ParseCtx) error {
		s := val.(*string)
		*s = *s + "_transformed"
		return nil
	})

	errs = field.Parse("short", &dest)
	assert.NotEmpty(t, errs)
	assert.NotEqual(t, "short_transformed", dest)
}

func TestStringRequiredAborts(t *testing.T) {
	field := String().Required().Len(3)
	var dest string

	errs := field.Parse("", &dest)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 1)
}

func TestStringUserTests(t *testing.T) {

	field := String().Test(TestFunc("test", func(val any, ctx ParseCtx) bool {
		return val == "test"
	}), Message("Invalid"))

	var dest string

	errs := field.Parse("test", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test", dest)

	errs = field.Parse("not test", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Invalid", errs[0].Message())

}

func TestStringRequired(t *testing.T) {
	field := String().Required(Message("a"))
	var dest string

	errs := field.Parse("", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, errs[0].Message(), "a")

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

// VALIDATORS / Tests / Validators

func TestStringLength(t *testing.T) {
	field := String().Len(3, Message("custom length"))
	var dest string

	errs := field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)

	errs = field.Parse("foobar", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom length", errs[0].Message())

	field = String().Min(5, Message("custom min")).Max(7, Message("custom max"))
	errs = field.Parse("123456789", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom max", errs[0].Message())

	assert.Equal(t, "123456789", dest)

	field = String().Min(5, Message("custom min")).Max(7, Message("custom max"))
	errs = field.Parse("1234567", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234567", dest)
}

func TestStringEmail(t *testing.T) {
	field := String().Email(Message("custom email"))
	var dest string

	errs := field.Parse("not an email", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom email", errs[0].Message())

	errs = field.Parse("test@example.com", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test@example.com", dest)
}

func TestStringURL(t *testing.T) {
	field := String().URL(Message("custom url"))
	var dest string

	errs := field.Parse("not a url", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom url", errs[0].Message())

	errs = field.Parse("http://example.com", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "http://example.com", dest)
}

func TestStringHasPrefix(t *testing.T) {
	field := String().HasPrefix("pre", Message("custom prefix"))
	var dest string

	errs := field.Parse("not prefixed", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom prefix", errs[0].Message())

	errs = field.Parse("prefix", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "prefix", dest)
}

func TestStringHasPostfix(t *testing.T) {
	field := String().HasSuffix("fix", Message("custom suffix"))
	var dest string

	errs := field.Parse("not postfixed", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom suffix", errs[0].Message())

	errs = field.Parse("postfix", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "postfix", dest)
}

func TestStringContains(t *testing.T) {
	field := String().Contains("contains", Message("custom contains"))
	var dest string

	errs := field.Parse("does not contain", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom contains", errs[0].Message())

	errs = field.Parse("contains", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "contains", dest)
}

func TestStringContainsDigit(t *testing.T) {
	field := String().ContainsDigit(Message("custom digit"))
	var dest string

	errs := field.Parse("no digit here", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom digit", errs[0].Message())

	errs = field.Parse("1234", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234", dest)
}

func TestStringContainsUpper(t *testing.T) {
	field := String().ContainsUpper(Message("custom upper"))
	var dest string

	errs := field.Parse("no uppercase here", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom upper", errs[0].Message())

	errs = field.Parse("UPPERCASE", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "UPPERCASE", dest)
}

func TestStringContainsSpecial(t *testing.T) {
	field := String().ContainsSpecial(Message("custom special"))
	var dest string

	errs := field.Parse("no special character here", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom special", errs[0].Message())

	errs = field.Parse("!@#$%", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "!@#$%", dest)
}

func TestStringOneOf(t *testing.T) {
	field := String().OneOf([]string{"apple", "banana", "cherry"}, Message("custom one of")).Required(Message("custom required"))
	var dest string

	errs := field.Parse("orange", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom one of", errs[0].Message())

	errs = field.Parse("banana", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "banana", dest)

	// Test with non-string input
	errs = field.Parse(123, &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom one of", errs[0].Message())

	// Test with empty string
	errs = field.Parse("", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom required", errs[0].Message())

	// Test with nil
	errs = field.Parse(nil, &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom required", errs[0].Message())
}

func TestStringUUID(t *testing.T) {
	field := String().UUID(Message("custom uuid msg"))
	var dest string

	errs := field.Parse("f81d4fae-7dec-11d0-a765-00a0c91e", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom uuid msg", errs[0].Message())

	errs = field.Parse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "f81d4fae-7dec-11d0-a765-00a0c91e6bf6", dest)

	errs = field.Parse("F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6", dest)
}

func TestStringRegex(t *testing.T) {
	r := regexp.MustCompile("^[0-9]{2}$")
	field := String().Match(r, Message("custom regex msg"))
	var dest string

	errs := field.Parse("f81d4fae-7dec-11d0-a765-00a0c91e", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom regex msg", errs[0].Message())

	errs = field.Parse("00", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "00", dest)
}

func TestStringSchemaOption(t *testing.T) {
	s := String(WithCoercer(func(original any) (value any, err error) {
		return "coerced", nil
	}))

	var result string
	err := s.Parse(123, &result)
	assert.Nil(t, err)
	assert.Equal(t, "coerced", result)
}

func TestStringGetType(t *testing.T) {
	s := String()
	assert.Equal(t, zconst.TypeString, s.getType())
}
