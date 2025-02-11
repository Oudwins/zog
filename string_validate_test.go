package zog

import (
	"regexp"
	"testing"

	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestValidateStringOptionalByDefault(t *testing.T) {
	field := String().Len(3).Contains("foo").HasPrefix("pre").HasSuffix("fix")
	var dest string

	errs := field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	field = field.Required()

	errs = field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	field.Required().Optional()
}

func TestValidateStringOptional(t *testing.T) {
	field := String().Required().Optional()
	var dest string

	errs := field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	dest = "foo"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestValidateStringPreTransform(t *testing.T) {
	field := String().Required().Len(6).PreTransform(func(val any, ctx ParseCtx) (any, error) {
		if x, ok := val.(string); ok {
			return "foo" + x, nil
		}
		return "invalid", nil
	})
	dest := "bar"

	errs := field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "foobar", dest)
}

func TestValidateStringTrim(t *testing.T) {
	field := String().Required().Trim()
	var dest string = " foo "

	errs := field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "foo", dest)

	dest = "123"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "123", dest)
}

func TestValidateStringPostTransform(t *testing.T) {
	field := String().Required().PostTransform(func(val any, ctx ParseCtx) error {
		s := val.(*string)
		*s = *s + "_transformed"
		return nil
	})
	var dest string = "hello"

	errs := field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "hello_transformed", dest)

	// Test that PostTransform is not applied when there's an error
	field = String().Required().Len(1).PostTransform(func(val any, ctx ParseCtx) error {
		s := val.(*string)
		*s = *s + "_transformed"
		return nil
	})

	dest = "short"
	errs = field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
	assert.NotEqual(t, "short_transformed", dest)
}

func TestValidateStringRequiredAborts(t *testing.T) {
	field := String().Required().Len(3)
	var dest string

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 1)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestValidateStringCustomTest(t *testing.T) {
	field := String().TestFunc(func(val any, ctx ParseCtx) bool {
		return val == "test"
	}, Message("Invalid"))

	var dest string = "test"

	errs := field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test", dest)

	dest = "not test"
	errs = field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Invalid", errs[0].Message())
}

func TestValidateStringRequired(t *testing.T) {
	field := String().Required(Message("a"))
	var dest string

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, errs[0].Message(), "a")

	dest = "foo"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestValidateStringDefault(t *testing.T) {
	field := String().Default("bar")
	var dest string

	errs := field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "bar", dest)

	dest = "foo"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestValidateStringCatch(t *testing.T) {
	field := String().Required().Min(5).Catch("error")
	var dest string = "x"

	errs := field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "error", dest)

	dest = "not error"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "not error", dest)
}

// VALIDATORS / Tests / Validators

func TestValidateStringLength(t *testing.T) {
	field := String().Len(3)
	var dest string = "foo"

	errs := field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)

	dest = "foobar"
	errs = field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	field = String().Min(5).Max(7)
	dest = "123456789"
	errs = field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	assert.Equal(t, "123456789", dest)

	field = String().Min(5).Max(7)
	dest = "1234567"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234567", dest)
}

func TestValidateStringEmail(t *testing.T) {
	field := String().Email()
	var dest string = "not an email"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	dest = "test@example.com"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test@example.com", dest)
}

func TestValidateStringURL(t *testing.T) {
	field := String().URL()
	var dest string = "not a url"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	dest = "http://example.com"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "http://example.com", dest)
}

func TestValidateStringHasPrefix(t *testing.T) {
	field := String().HasPrefix("pre")
	var dest string = "not prefixed"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	dest = "prefix"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "prefix", dest)
}

func TestValidateStringHasPostfix(t *testing.T) {
	field := String().HasSuffix("fix")
	var dest string = "not postfixed"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	dest = "postfix"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "postfix", dest)
}

func TestValidateStringContains(t *testing.T) {
	field := String().Contains("contains", Message("custom contains"))
	var dest string = "does not contain"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom contains", errs[0].Message())

	dest = "contains"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "contains", dest)
}

func TestValidateStringContainsDigit(t *testing.T) {
	field := String().ContainsDigit(Message("custom digit"))
	var dest string = "no digit here"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom digit", errs[0].Message())

	dest = "1234"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234", dest)
}

func TestValidateStringContainsUpper(t *testing.T) {
	field := String().ContainsUpper(Message("custom upper"))
	var dest string = "no uppercase here"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom upper", errs[0].Message())

	dest = "UPPERCASE"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "UPPERCASE", dest)
}

func TestValidateStringContainsSpecial(t *testing.T) {
	field := String().ContainsSpecial(Message("custom special"))
	var dest string = "no special character here"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom special", errs[0].Message())

	dest = "!@#$%"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "!@#$%", dest)
}

func TestValidateStringOneOf(t *testing.T) {
	field := String().OneOf([]string{"apple", "banana", "cherry"}, Message("custom one of")).Required(Message("custom required"))
	var dest string = "orange"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom one of", errs[0].Message())

	dest = "banana"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)

	assert.Equal(t, "banana", dest)

	dest = ""
	errs = field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom required", errs[0].Message())
}

func TestValidateStringUUID(t *testing.T) {
	field := String().UUID(Message("custom uuid msg"))
	var dest string = "f81d4fae-7dec-11d0-a765-00a0c91e"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom uuid msg", errs[0].Message())

	dest = "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "f81d4fae-7dec-11d0-a765-00a0c91e6bf6", dest)

	dest = "F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6", dest)
}

func TestValidateStringRegex(t *testing.T) {
	r := regexp.MustCompile("^[0-9]{2}$")
	field := String().Match(r, Message("custom regex msg"))
	var dest string = "f81d4fae-7dec-11d0-a765-00a0c91e"

	errs := field.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom regex msg", errs[0].Message())

	dest = "00"
	errs = field.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, "00", dest)
}

func TestValidateStringSchemaOption(t *testing.T) {
	var tp zconst.ZogType
	s := String(func(s ZogSchema) {
		tp = string(s.getType())
	})

	var result string = "123"
	err := s.Validate(&result)
	assert.Nil(t, err)
	assert.Equal(t, "123", result)
	assert.Equal(t, zconst.TypeString, tp)
}

func TestValidateStringGetType(t *testing.T) {
	s := String()
	assert.Equal(t, zconst.TypeString, s.getType())
}
