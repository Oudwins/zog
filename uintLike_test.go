package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Custom unsigned integer type for testing
type TestCustomUint uint

func TestUintLikeParse(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid uint value", 42, false, TestCustomUint(42)},
		{"Valid zero value", 0, false, TestCustomUint(0)},
		{"Valid string number", "123", false, TestCustomUint(123)},
		{"Invalid type (bool)", true, true, TestCustomUint(1)},
		{"Invalid string", "abc", true, TestCustomUint(0)},
		{"Invalid negative value", -10, true, TestCustomUint(0)},
	}

	uintProc := UintLike[TestCustomUint]()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
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

func TestUintLikeSchemaOption(t *testing.T) {
	s := UintLike[TestCustomUint](WithCoercer(func(original any) (value any, err error) {
		return TestCustomUint(999), nil
	}))
	var result TestCustomUint
	err := s.Parse("invalid", &result)
	assert.Nil(t, err)
	assert.Equal(t, TestCustomUint(999), result)
}

func TestUintLikeExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		uintProc := UintLike[TestCustomUint]()
		var result TestCustomUint
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := uintProc.Parse(42, &result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestUintLikeRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid uint value", 42, false, TestCustomUint(42)},
		{"Valid zero value", 0, false, TestCustomUint(0)},
		{"Nil value", nil, true, TestCustomUint(0)},
	}
	uintProc := UintLike[TestCustomUint]().Required(Message("test"))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
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

func TestUintLikeOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid uint value", 42, false, TestCustomUint(42)},
		{"Valid zero value", 0, false, TestCustomUint(0)},
		{"Nil value", nil, false, TestCustomUint(0)},
	}
	uintProc := UintLike[TestCustomUint]().Optional()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
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

func TestUintLikeDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		default_  TestCustomUint
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid uint value", 42, TestCustomUint(100), false, TestCustomUint(42)},
		{"Nil value with default", nil, TestCustomUint(100), false, TestCustomUint(100)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintProc := UintLike[TestCustomUint]().Default(test.default_)
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
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

func TestUintLikeCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     TestCustomUint
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid uint value", 42, TestCustomUint(999), false, TestCustomUint(42)},
		{"Invalid type with catch", "invalid", TestCustomUint(999), false, TestCustomUint(999)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintProc := UintLike[TestCustomUint]().Catch(test.catch)
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
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

func TestUintLikeOneOf(t *testing.T) {
	enum := []TestCustomUint{1, 2, 3, 5, 8}
	uintProc := UintLike[TestCustomUint]().OneOf(enum)
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid enum value", 3, false, TestCustomUint(3)},
		{"Invalid enum value", 4, true, TestCustomUint(4)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestUintLikeComparisons(t *testing.T) {
	tests := []struct {
		name      string
		schema    *NumberSchema[TestCustomUint]
		data      any
		expectErr bool
		expected  TestCustomUint
	}{
		{"GT valid", UintLike[TestCustomUint]().GT(5), 10, false, TestCustomUint(10)},
		{"GT invalid", UintLike[TestCustomUint]().GT(5), 3, true, TestCustomUint(3)},
		{"LT valid", UintLike[TestCustomUint]().LT(10), 5, false, TestCustomUint(5)},
		{"LT invalid", UintLike[TestCustomUint]().LT(10), 15, true, TestCustomUint(15)},
		{"GTE valid equal", UintLike[TestCustomUint]().GTE(5), 5, false, TestCustomUint(5)},
		{"GTE valid greater", UintLike[TestCustomUint]().GTE(5), 10, false, TestCustomUint(10)},
		{"GTE invalid", UintLike[TestCustomUint]().GTE(5), 3, true, TestCustomUint(3)},
		{"LTE valid equal", UintLike[TestCustomUint]().LTE(10), 10, false, TestCustomUint(10)},
		{"LTE valid less", UintLike[TestCustomUint]().LTE(10), 5, false, TestCustomUint(5)},
		{"LTE invalid", UintLike[TestCustomUint]().LTE(10), 15, true, TestCustomUint(15)},
		{"EQ valid", UintLike[TestCustomUint]().EQ(42), 42, false, TestCustomUint(42)},
		{"EQ invalid", UintLike[TestCustomUint]().EQ(42), 43, true, TestCustomUint(43)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomUint
			errs := test.schema.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestUintLikeNot(t *testing.T) {
	tests := []struct {
		name      string
		schema    *NumberSchema[TestCustomUint]
		data      any
		expectErr bool
		expected  TestCustomUint
	}{
		{"Not EQ valid", UintLike[TestCustomUint]().Not().EQ(42), 43, false, TestCustomUint(43)},
		{"Not EQ invalid", UintLike[TestCustomUint]().Not().EQ(42), 42, true, TestCustomUint(42)},
		{"Not OneOf valid", UintLike[TestCustomUint]().Not().OneOf([]TestCustomUint{1, 2, 3}), 4, false, TestCustomUint(4)},
		{"Not OneOf invalid", UintLike[TestCustomUint]().Not().OneOf([]TestCustomUint{1, 2, 3}), 2, true, TestCustomUint(2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomUint
			errs := test.schema.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestUintLikeTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.Transform[*TestCustomUint]
		expectErr bool
		expected  TestCustomUint
	}{
		{"Double value", 5, func(val *TestCustomUint, ctx Ctx) error {
			*val = *val * 2
			return nil
		}, false, TestCustomUint(10)},
		{"No change", 42, func(val *TestCustomUint, ctx Ctx) error { return nil }, false, TestCustomUint(42)},
		{"Invalid transform", 42, func(val *TestCustomUint, ctx Ctx) error { return fmt.Errorf("fail") }, true, TestCustomUint(42)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintProc := UintLike[TestCustomUint]().Transform(test.transform)
			var result TestCustomUint
			errs := uintProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestUintLikeCustomTest(t *testing.T) {
	validator := UintLike[TestCustomUint]().TestFunc(func(val *TestCustomUint, ctx Ctx) bool {
		return *val%2 == 0 // Even numbers only
	}, Message("must be even"))
	tests := []struct {
		name      string
		input     uint
		expectErr bool
		expected  TestCustomUint
	}{
		{"valid even value", 4, false, TestCustomUint(4)},
		{"invalid odd value", 5, true, TestCustomUint(5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest TestCustomUint
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

func TestUintLikeGetType(t *testing.T) {
	s := UintLike[TestCustomUint]()
	assert.Equal(t, zconst.TypeNumber, s.getType())
}

// Validation tests
func TestUintLikeValidate(t *testing.T) {
	tests := []struct {
		name string
		data TestCustomUint
	}{
		{"Valid positive value", TestCustomUint(42)},
		{"Valid zero value", TestCustomUint(0)},
		{"Valid max value", TestCustomUint(1000)},
	}
	uintProc := UintLike[TestCustomUint]()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := uintProc.Validate(&test.data)
			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestUintLikeValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		uintProc := UintLike[TestCustomUint]()
		var result TestCustomUint
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := uintProc.Validate(&result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestUintLikeValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomUint
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid positive value", TestCustomUint(42), false, TestCustomUint(42)},
		{"Valid large value", TestCustomUint(1000), false, TestCustomUint(1000)},
		{"Zero value", TestCustomUint(0), true, TestCustomUint(0)},
	}
	uintProc := UintLike[TestCustomUint]().Required()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := uintProc.Validate(&test.data)
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

func TestUintLikeValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomUint
		expected  TestCustomUint
		proc      *NumberSchema[TestCustomUint]
		expectErr bool
	}{
		{"Optional by default", TestCustomUint(0), TestCustomUint(0), UintLike[TestCustomUint](), false},
		{"Optional overrides required", TestCustomUint(0), TestCustomUint(0), UintLike[TestCustomUint]().Required().Optional(), false},
		{"Required errors on zero", TestCustomUint(0), TestCustomUint(0), UintLike[TestCustomUint]().Required(), true},
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

func TestUintLikeValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomUint
		default_  TestCustomUint
		expectErr bool
		expected  TestCustomUint
	}{
		{"Valid value", TestCustomUint(42), TestCustomUint(100), false, TestCustomUint(42)},
		{"Zero value with default", TestCustomUint(0), TestCustomUint(100), false, TestCustomUint(100)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintProc := UintLike[TestCustomUint]().Default(test.default_)
			errs := uintProc.Validate(&test.data)
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

func TestUintLikeValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     TestCustomUint
		catch    TestCustomUint
		expected TestCustomUint
	}{
		{"Without catch", TestCustomUint(42), TestCustomUint(999), TestCustomUint(42)},
		{"With catch", TestCustomUint(0), TestCustomUint(999), TestCustomUint(999)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintProc := UintLike[TestCustomUint]().TestFunc(func(val *TestCustomUint, ctx Ctx) bool {
				return *val != 0
			}).Catch(test.catch).Required()
			errs := uintProc.Validate(&test.data)
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if test.data != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, test.data)
			}
		})
	}
}

func TestUintLikeValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomUint
		transform p.Transform[*TestCustomUint]
		expectErr bool
		expected  TestCustomUint
	}{
		{"Double value", TestCustomUint(5), func(val *TestCustomUint, ctx Ctx) error {
			*val = *val * 2
			return nil
		}, false, TestCustomUint(10)},
		{"No change", TestCustomUint(42), func(val *TestCustomUint, ctx Ctx) error { return nil }, false, TestCustomUint(42)},
		{"Invalid transform", TestCustomUint(42), func(val *TestCustomUint, ctx Ctx) error { return fmt.Errorf("fail") }, true, TestCustomUint(42)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintProc := UintLike[TestCustomUint]().Transform(test.transform)
			errs := uintProc.Validate(&test.data)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if test.data != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, test.data)
			}
		})
	}
}

func TestUintLikeValidateCustomTest(t *testing.T) {
	validator := UintLike[TestCustomUint]().TestFunc(func(val *TestCustomUint, ctx Ctx) bool {
		assert.Equal(t, TestCustomUint(42), *val)
		return true
	}, Message("custom"))
	dest := TestCustomUint(42)
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, TestCustomUint(42), dest)
}
