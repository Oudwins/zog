package zog

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Custom string type for testing
// (stringer is not required, but this is a placeholder for custom string type)
//
//go:generate stringer -type=CustomString
type CustomString string

func TestStringLikeParse(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  CustomString
	}{
		{"Valid string value", "hello", false, CustomString("hello")},
		{"Valid empty string", "", false, CustomString("")},
		{"Valid type (int)", 123, false, CustomString("123")},
		{"Valid type (bool)", true, false, CustomString("true")},
		{"Invalid type nil optional so fine", nil, false, CustomString("")},
	}

	strProc := StringLike[CustomString]()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if len(errs) > 0 && !test.expectErr {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeSchemaOption(t *testing.T) {
	s := StringLike[CustomString](WithCoercer(func(original any) (value any, err error) {
		return CustomString("coerced"), nil
	}))
	var result CustomString
	err := s.Parse(123, &result)
	assert.Nil(t, err)
	assert.Equal(t, CustomString("coerced"), result)
}

func TestStringLikeExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		strProc := StringLike[CustomString]()
		var result CustomString
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := strProc.Parse("hello", &result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestStringLikeRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  CustomString
	}{
		{"Valid string value", "abc", false, CustomString("abc")},
		{"Empty string", "", false, CustomString("")},
		{"Nil value", nil, true, CustomString("")},
	}
	strProc := StringLike[CustomString]().Required(Message("test"))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}
			if test.expectErr && errs[0].Message != "test" {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, "test", errs[0].Message)
			}
			if !test.expectErr && result != test.expected {
				t.Errorf("On Run %s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestStringLikeOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  CustomString
	}{
		{"Valid string value", "abc", false, CustomString("abc")},
		{"Empty string", "", false, CustomString("")},
		{"Nil value", nil, false, CustomString("")},
	}
	strProc := StringLike[CustomString]().Optional()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		default_  CustomString
		expectErr bool
		expected  CustomString
	}{
		{"Valid string value", "abc", CustomString("def"), false, CustomString("abc")},
		{"Nil value with default", nil, CustomString("def"), false, CustomString("def")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().Default(test.default_)
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if result != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestStringLikeCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     CustomString
		expectErr bool
		expected  CustomString
	}{
		{"Valid string value", "abc", CustomString("catch"), false, CustomString("abc")},
		{"Valid type with catch", 123, CustomString("catch"), false, CustomString("123")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().Catch(test.catch)
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if result != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestStringLikeOneOf(t *testing.T) {
	enum := []CustomString{"a", "b", "c"}
	strProc := StringLike[CustomString]().OneOf(enum)
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  CustomString
	}{
		{"Valid enum value", "a", false, CustomString("a")},
		{"Invalid enum value", "z", true, CustomString("z")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeMinMaxLen(t *testing.T) {
	strProc := StringLike[CustomString]().Min(2).Max(4).Len(3)
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  CustomString
	}{
		{"Valid length", "abc", false, CustomString("abc")},
		{"Too short", "a", true, CustomString("a")},
		{"Too long", "abcde", true, CustomString("abcde")},
		{"Wrong length", "ab", true, CustomString("ab")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeEmailURL(t *testing.T) {
	emailProc := StringLike[CustomString]().Email()
	urlProc := StringLike[CustomString]().URL()
	t.Run("Valid email", func(t *testing.T) {
		var result CustomString
		errs := emailProc.Parse("test@example.com", &result)
		assert.Empty(t, errs)
		assert.Equal(t, CustomString("test@example.com"), result)
	})
	t.Run("Invalid email", func(t *testing.T) {
		var result CustomString
		errs := emailProc.Parse("not-an-email", &result)
		assert.NotEmpty(t, errs)
	})
	t.Run("Valid URL", func(t *testing.T) {
		var result CustomString
		errs := urlProc.Parse("https://example.com", &result)
		assert.Empty(t, errs)
		assert.Equal(t, CustomString("https://example.com"), result)
	})
	t.Run("Invalid URL", func(t *testing.T) {
		var result CustomString
		errs := urlProc.Parse("not-a-url", &result)
		assert.NotEmpty(t, errs)
	})
}

func TestStringLikePrefixSuffixContains(t *testing.T) {
	strProc := StringLike[CustomString]().HasPrefix("pre").HasSuffix("suf").Contains("mid")
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  CustomString
	}{
		{"Valid string", "pre-midsuf", false, CustomString("pre-midsuf")},
		{"Missing prefix", "midsuf", true, CustomString("midsuf")},
		{"Missing suffix", "pre-mid", true, CustomString("pre-mid")},
		{"Missing contains", "presuf", true, CustomString("presuf")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeContainsUpperDigitSpecial(t *testing.T) {
	upperProc := StringLike[CustomString]().ContainsUpper()
	digitProc := StringLike[CustomString]().ContainsDigit()
	specialProc := StringLike[CustomString]().ContainsSpecial()
	t.Run("Contains upper", func(t *testing.T) {
		var result CustomString
		errs := upperProc.Parse("abcD", &result)
		assert.Empty(t, errs)
	})
	t.Run("No upper", func(t *testing.T) {
		var result CustomString
		errs := upperProc.Parse("abcd", &result)
		assert.NotEmpty(t, errs)
	})
	t.Run("Contains digit", func(t *testing.T) {
		var result CustomString
		errs := digitProc.Parse("abc1", &result)
		assert.Empty(t, errs)
	})
	t.Run("No digit", func(t *testing.T) {
		var result CustomString
		errs := digitProc.Parse("abcd", &result)
		assert.NotEmpty(t, errs)
	})
	t.Run("Contains special", func(t *testing.T) {
		var result CustomString
		errs := specialProc.Parse("abc$", &result)
		assert.Empty(t, errs)
	})
	t.Run("No special", func(t *testing.T) {
		var result CustomString
		errs := specialProc.Parse("abcd", &result)
		assert.NotEmpty(t, errs)
	})
}

func TestStringLikeUUIDMatch(t *testing.T) {
	uuidProc := StringLike[CustomString]().UUID()
	matchProc := StringLike[CustomString]().Match(regexp.MustCompile(`^abc[0-9]+$`))
	t.Run("Valid UUID", func(t *testing.T) {
		var result CustomString
		errs := uuidProc.Parse("123e4567-e89b-12d3-a456-426614174000", &result)
		assert.Empty(t, errs)
	})
	t.Run("Invalid UUID", func(t *testing.T) {
		var result CustomString
		errs := uuidProc.Parse("not-a-uuid", &result)
		assert.NotEmpty(t, errs)
	})
	t.Run("Valid match", func(t *testing.T) {
		var result CustomString
		errs := matchProc.Parse("abc123", &result)
		assert.Empty(t, errs)
	})
	t.Run("Invalid match", func(t *testing.T) {
		var result CustomString
		errs := matchProc.Parse("def123", &result)
		assert.NotEmpty(t, errs)
	})
}

func TestStringLikeTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.Transform[*CustomString]
		expectErr bool
		expected  CustomString
	}{
		{"To upper", "abc", func(val *CustomString, ctx Ctx) error {
			*val = CustomString(strings.ToUpper(string(*val)))
			return nil
		}, false, CustomString("ABC")},
		{"No change", "abc", func(val *CustomString, ctx Ctx) error { return nil }, false, CustomString("abc")},
		{"Invalid transform", "abc", func(val *CustomString, ctx Ctx) error { return fmt.Errorf("fail") }, true, CustomString("abc")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().Transform(test.transform)
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeTrim(t *testing.T) {
	strProc := StringLike[CustomString]().Trim()
	var result CustomString
	errs := strProc.Parse("  abc  ", &result)
	assert.Empty(t, errs)
	assert.Equal(t, CustomString("abc"), result)
}

func TestStringLikeCustomTest(t *testing.T) {
	validator := StringLike[CustomString]().TestFunc(func(val *CustomString, ctx Ctx) bool {
		return *val == CustomString("pass")
	}, Message("custom"))
	tests := []struct {
		name      string
		input     string
		expectErr bool
		expected  CustomString
	}{
		{"valid value", "pass", false, CustomString("pass")},
		{"invalid value", "fail", true, CustomString("fail")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest CustomString
			errs := validator.Parse(tt.input, &dest)
			if (len(errs) > 0) != tt.expectErr {
				t.Errorf("got errors %v, expectErr %v", errs, tt.expectErr)
			}
			if !tt.expectErr {
				assert.Equal(t, tt.expected, dest)
			}
		})
	}
}

func TestStringLikeEQ(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		eqValue   CustomString
		expectErr bool
		expected  CustomString
	}{
		{"Equal value", "abc", CustomString("abc"), false, CustomString("abc")},
		{"Not equal value", "def", CustomString("abc"), true, CustomString("def")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().TestFunc(func(val *CustomString, ctx Ctx) bool {
				return *val == test.eqValue
			})
			var result CustomString
			errs := strProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestStringLikeGetType(t *testing.T) {
	s := StringLike[CustomString]()
	assert.Equal(t, zconst.TypeString, s.getType())
}

// Validation tests
func TestStringLikeValidate(t *testing.T) {
	tests := []struct {
		name string
		data CustomString
	}{
		{"Valid value", CustomString("abc")},
		{"Empty value", CustomString("")},
	}
	strProc := StringLike[CustomString]()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := strProc.Validate(&test.data)
			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestStringLikeValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		strProc := StringLike[CustomString]()
		var result CustomString
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := strProc.Validate(&result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestStringLikeValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomString
		expectErr bool
		expected  CustomString
	}{
		{"Valid value", CustomString("abc"), false, CustomString("abc")},
		{"Empty value", CustomString(""), true, CustomString("")},
	}
	strProc := StringLike[CustomString]().Required()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := strProc.Validate(&test.data)
			if test.expectErr {
				assert.NotEmpty(t, errs)
				tutils.VerifyDefaultIssueMessages(t, errs)
			} else {
				assert.Empty(t, errs)
			}
			assert.Equal(t, test.data, test.expected)
		})
	}
}

func TestStringLikeValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomString
		expected  CustomString
		proc      *StringSchema[CustomString]
		expectErr bool
	}{
		{"Optional by default", CustomString("abc"), CustomString("abc"), StringLike[CustomString](), false},
		{"Optional overrides required", CustomString("abc"), CustomString("abc"), StringLike[CustomString]().Required().Optional(), false},
		{"Required errors on empty", CustomString(""), CustomString(""), StringLike[CustomString]().Required(), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := test.proc.Validate(&test.data)
			if test.expectErr {
				assert.NotEmpty(t, errs)
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			assert.Equal(t, test.data, test.expected)
		})
	}
}

func TestStringLikeValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomString
		default_  CustomString
		expectErr bool
		expected  CustomString
	}{
		{"Valid value", CustomString("abc"), CustomString("def"), false, CustomString("abc")},
		{"Empty value with default", CustomString(""), CustomString("def"), false, CustomString("def")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().Default(test.default_)
			errs := strProc.Validate(&test.data)
			if test.expectErr {
				assert.NotEmpty(t, errs)
				tutils.VerifyDefaultIssueMessages(t, errs)
			} else {
				assert.Empty(t, errs)
			}
			assert.Equal(t, test.data, test.expected)
		})
	}
}

func TestStringLikeValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     CustomString
		catch    CustomString
		expected CustomString
	}{
		{"Without catch", CustomString("abc"), CustomString("catch"), CustomString("abc")},
		{"With catch", CustomString(""), CustomString("catch"), CustomString("catch")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().TestFunc(func(val *CustomString, ctx Ctx) bool {
				return *val != ""
			}).Catch(test.catch).Required()
			errs := strProc.Validate(&test.data)
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if test.data != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, test.data)
			}
		})
	}
}

func TestStringLikeValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomString
		transform p.Transform[*CustomString]
		expectErr bool
		expected  CustomString
	}{
		{"To upper", CustomString("abc"), func(val *CustomString, ctx Ctx) error {
			*val = CustomString(strings.ToUpper(string(*val)))
			return nil
		}, false, CustomString("ABC")},
		{"No change", CustomString("abc"), func(val *CustomString, ctx Ctx) error { return nil }, false, CustomString("abc")},
		{"Invalid transform", CustomString("abc"), func(val *CustomString, ctx Ctx) error { return fmt.Errorf("fail") }, true, CustomString("abc")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strProc := StringLike[CustomString]().Transform(test.transform)
			errs := strProc.Validate(&test.data)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if test.data != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, test.data)
			}
		})
	}
}

func TestStringLikeValidateCustomTest(t *testing.T) {
	validator := StringLike[CustomString]().TestFunc(func(val *CustomString, ctx Ctx) bool {
		assert.Equal(t, CustomString("abc"), *val)
		return true
	}, Message("custom"))
	dest := CustomString("abc")
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, CustomString("abc"), dest)
}
