package zog

import (
	"regexp"
	"testing"

	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestStringOptionalByDefault(t *testing.T) {
	field := String().Len(3).Contains("foo").HasPrefix("pre").HasSuffix("fix")
	var dest string

	errs := field.Parse(nil, &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "", dest)

	field = field.Required()

	errs = field.Parse(nil, &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

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

func TestStringTrim(t *testing.T) {
	field := String().Required().Trim()
	var dest string

	errs := field.Parse(" foo ", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "foo", dest)

	errs = field.Parse(123, &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "123", dest)
}

func TestStringPostTransform(t *testing.T) {
	field := String().Required().Transform(func(val any, ctx Ctx) error {
		s := val.(*string)
		*s = *s + "_transformed"
		return nil
	})
	var dest string

	errs := field.Parse("hello", &dest)
	assert.Empty(t, errs)
	assert.Equal(t, "hello_transformed", dest)
}

func TestStringRequiredAborts(t *testing.T) {
	field := String().Required().Len(3)
	var dest string

	errs := field.Parse("", &dest)
	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 1)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestStringCustomTest(t *testing.T) {

	field := String().TestFunc(func(val any, ctx Ctx) bool {
		return val == "test"
	}, Message("Invalid"))

	var dest string

	errs := field.Parse("test", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test", dest)

	errs = field.Parse("not test", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Invalid", errs[0].Message)

}

func TestStringRequired(t *testing.T) {
	field := String().Required(Message("a"))
	var dest string

	errs := field.Parse(nil, &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, errs[0].Message, "a")

	errs = field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)
}

func TestStringDefault(t *testing.T) {
	field := String().Default("bar")
	var dest string

	errs := field.Parse(nil, &dest)
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
	field := String().Len(3)
	var dest string

	errs := field.Parse("foo", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "foo", dest)

	errs = field.Parse("foobar", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	field = String().Min(5).Max(7)
	errs = field.Parse("123456789", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	assert.Equal(t, "123456789", dest)

	field = String().Min(5).Max(7)
	errs = field.Parse("1234567", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234567", dest)
}

func TestStringEmail(t *testing.T) {
	field := String().Email()
	var dest string

	errs := field.Parse("not an email", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	errs = field.Parse("test@example.com", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "test@example.com", dest)
}

func TestStringURL(t *testing.T) {
	field := String().URL()
	var dest string

	errs := field.Parse("not a url", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	errs = field.Parse("http://example.com", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "http://example.com", dest)
}

func TestStringHasPrefix(t *testing.T) {
	field := String().HasPrefix("pre")
	var dest string

	errs := field.Parse("not prefixed", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	errs = field.Parse("prefix", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "prefix", dest)
}

func TestStringHasPostfix(t *testing.T) {
	field := String().HasSuffix("fix")
	var dest string

	errs := field.Parse("not postfixed", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	errs = field.Parse("postfix", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "postfix", dest)
}

func TestStringContains(t *testing.T) {
	field := String().Contains("contains")
	var dest string

	errs := field.Parse("not containing", &dest)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	errs = field.Parse("this contains that", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "this contains that", dest)
}

func TestStringContainsDigit(t *testing.T) {
	field := String().ContainsDigit(Message("custom digit"))
	var dest string

	errs := field.Parse("no digit here", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom digit", errs[0].Message)

	errs = field.Parse("1234", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "1234", dest)
}

func TestStringContainsUpper(t *testing.T) {
	field := String().ContainsUpper(Message("custom upper"))
	var dest string

	errs := field.Parse("no uppercase here", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom upper", errs[0].Message)

	errs = field.Parse("UPPERCASE", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "UPPERCASE", dest)
}

func TestStringContainsSpecial(t *testing.T) {
	field := String().ContainsSpecial(Message("custom special"))
	var dest string

	errs := field.Parse("no special character here", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom special", errs[0].Message)

	errs = field.Parse("!@#$%", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "!@#$%", dest)
}

func TestStringOneOf(t *testing.T) {
	field := String().OneOf([]string{"apple", "banana", "cherry"}, Message("custom one of")).Required(Message("custom required"))
	var dest string

	errs := field.Parse("orange", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom one of", errs[0].Message)

	errs = field.Parse("banana", &dest)
	assert.Empty(t, errs)

	assert.Equal(t, "banana", dest)

	// Test with non-string input
	errs = field.Parse(123, &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom one of", errs[0].Message)

	// Test with empty string
	errs = field.Parse("", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom one of", errs[0].Message)

	// Test with nil
	errs = field.Parse(nil, &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom required", errs[0].Message)
}

func TestStringUUID(t *testing.T) {
	field := String().UUID(Message("custom uuid msg"))
	var dest string

	errs := field.Parse("f81d4fae-7dec-11d0-a765-00a0c91e", &dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom uuid msg", errs[0].Message)

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
	assert.Equal(t, "custom regex msg", errs[0].Message)

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

func TestStringNot(t *testing.T) {
	tests := map[string]struct {
		schema         *StringSchema[string]
		strVal         string
		expectedErrMap internals.ZogIssueList
	}{
		"not len success": {
			schema:         String().Not().Len(10).Contains("test"),
			strVal:         "test",
			expectedErrMap: nil,
		},
		"not len fail": {
			schema: String().Not().Len(4).Contains("t"),
			strVal: "test",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_len",
					Params:  map[string]any{"len": 4},
					Dtype:   "string",
					Value:   tutils.PtrOf("test"),
					Message: "string must not be exactly 4 character(s)",
					Err:     nil,
				},
			},
		},
		"not email": {
			schema:         String().Not().Email(),
			strVal:         "not-an-email",
			expectedErrMap: nil,
		},
		"not email failure": {
			schema: String().Not().Email(),
			strVal: "test@test.com",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_email",
					Params:  nil,
					Dtype:   "string",
					Value:   tutils.PtrOf("test@test.com"),
					Message: "must not be a valid email",
					Err:     nil,
				},
			},
		},
		"not with empty": {
			schema: String().Not().Len(1),
			strVal: "a",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_len",
					Params:  map[string]any{"len": 1},
					Dtype:   "string",
					Value:   tutils.PtrOf("a"),
					Message: "string must not be exactly 1 character(s)",
					Err:     nil,
				},
			},
		},
		"not url": {
			schema:         String().Not().URL(),
			strVal:         "not a url",
			expectedErrMap: nil,
		},
		"not url failure": {
			schema: String().Not().URL(),
			strVal: "https://google.com",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_url",
					Dtype:   "string",
					Value:   tutils.PtrOf("https://google.com"),
					Message: "must not be a valid URL",
					Err:     nil,
				},
			},
		},
		"not has prefix": {
			schema:         String().Not().HasPrefix("test_"),
			strVal:         "value",
			expectedErrMap: nil,
		},
		"not has prefix failure": {
			schema: String().Not().HasPrefix("test_"),
			strVal: "test_value",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_prefix",
					Params:  map[string]any{"prefix": "test_"},
					Dtype:   "string",
					Value:   tutils.PtrOf("test_value"),
					Message: "string must not start with test_",
					Err:     nil,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var dest string
			errMap := tc.schema.Parse(tc.strVal, &dest)
			assert.Equal(t, tc.expectedErrMap, errMap)
		})
	}
}
