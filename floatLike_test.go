package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Custom float type for testing
type TestCustomFloat float64

func TestFloatLikeParse(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid float value", 42.5, false, TestCustomFloat(42.5)},
		{"Valid integer value", 42, false, TestCustomFloat(42.0)},
		{"Valid zero value", 0, false, TestCustomFloat(0.0)},
		{"Valid negative value", -10.5, false, TestCustomFloat(-10.5)},
		{"Valid string number", "123.456", false, TestCustomFloat(123.456)},
		{"Invalid type (bool)", true, true, TestCustomFloat(0.0)},
		{"Invalid string", "abc", true, TestCustomFloat(0.0)},
	}

	floatProc := FloatLike[TestCustomFloat]()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
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

func TestFloatLikeSchemaOption(t *testing.T) {
	s := FloatLike[TestCustomFloat](WithCoercer(func(original any) (value any, err error) {
		return TestCustomFloat(999.999), nil
	}))
	var result TestCustomFloat
	err := s.Parse("invalid", &result)
	assert.Nil(t, err)
	assert.Equal(t, TestCustomFloat(999.999), result)
}

func TestFloatLikeExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		floatProc := FloatLike[TestCustomFloat]()
		var result TestCustomFloat
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := floatProc.Parse(42.5, &result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestFloatLikeRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid float value", 42.5, false, TestCustomFloat(42.5)},
		{"Valid zero value", 0.0, false, TestCustomFloat(0.0)},
		{"Nil value", nil, true, TestCustomFloat(0.0)},
	}
	floatProc := FloatLike[TestCustomFloat]().Required(Message("test"))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
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

func TestFloatLikeOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid float value", 42.5, false, TestCustomFloat(42.5)},
		{"Valid zero value", 0.0, false, TestCustomFloat(0.0)},
		{"Nil value", nil, false, TestCustomFloat(0.0)},
	}
	floatProc := FloatLike[TestCustomFloat]().Optional()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
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

func TestFloatLikeDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		default_  TestCustomFloat
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid float value", 42.5, TestCustomFloat(100.5), false, TestCustomFloat(42.5)},
		{"Nil value with default", nil, TestCustomFloat(100.5), false, TestCustomFloat(100.5)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatProc := FloatLike[TestCustomFloat]().Default(test.default_)
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
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

func TestFloatLikeCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     TestCustomFloat
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid float value", 42.5, TestCustomFloat(999.999), false, TestCustomFloat(42.5)},
		{"Invalid type with catch", "invalid", TestCustomFloat(999.999), false, TestCustomFloat(999.999)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatProc := FloatLike[TestCustomFloat]().Catch(test.catch)
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
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

func TestFloatLikeOneOf(t *testing.T) {
	enum := []TestCustomFloat{1.1, 2.2, 3.3, 5.5, 8.8}
	floatProc := FloatLike[TestCustomFloat]().OneOf(enum)
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid enum value", 3.3, false, TestCustomFloat(3.3)},
		{"Invalid enum value", 4.4, true, TestCustomFloat(4.4)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestFloatLikeComparisons(t *testing.T) {
	tests := []struct {
		name      string
		schema    *NumberSchema[TestCustomFloat]
		data      any
		expectErr bool
		expected  TestCustomFloat
	}{
		{"GT valid", FloatLike[TestCustomFloat]().GT(5.5), 10.5, false, TestCustomFloat(10.5)},
		{"GT invalid", FloatLike[TestCustomFloat]().GT(5.5), 3.3, true, TestCustomFloat(3.3)},
		{"LT valid", FloatLike[TestCustomFloat]().LT(10.5), 5.5, false, TestCustomFloat(5.5)},
		{"LT invalid", FloatLike[TestCustomFloat]().LT(10.5), 15.5, true, TestCustomFloat(15.5)},
		{"GTE valid equal", FloatLike[TestCustomFloat]().GTE(5.5), 5.5, false, TestCustomFloat(5.5)},
		{"GTE valid greater", FloatLike[TestCustomFloat]().GTE(5.5), 10.5, false, TestCustomFloat(10.5)},
		{"GTE invalid", FloatLike[TestCustomFloat]().GTE(5.5), 3.3, true, TestCustomFloat(3.3)},
		{"LTE valid equal", FloatLike[TestCustomFloat]().LTE(10.5), 10.5, false, TestCustomFloat(10.5)},
		{"LTE valid less", FloatLike[TestCustomFloat]().LTE(10.5), 5.5, false, TestCustomFloat(5.5)},
		{"LTE invalid", FloatLike[TestCustomFloat]().LTE(10.5), 15.5, true, TestCustomFloat(15.5)},
		{"EQ valid", FloatLike[TestCustomFloat]().EQ(42.5), 42.5, false, TestCustomFloat(42.5)},
		{"EQ invalid", FloatLike[TestCustomFloat]().EQ(42.5), 43.5, true, TestCustomFloat(43.5)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomFloat
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

func TestFloatLikeNot(t *testing.T) {
	tests := []struct {
		name      string
		schema    *NumberSchema[TestCustomFloat]
		data      any
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Not EQ valid", FloatLike[TestCustomFloat]().Not().EQ(42.5), 43.5, false, TestCustomFloat(43.5)},
		{"Not EQ invalid", FloatLike[TestCustomFloat]().Not().EQ(42.5), 42.5, true, TestCustomFloat(42.5)},
		{"Not OneOf valid", FloatLike[TestCustomFloat]().Not().OneOf([]TestCustomFloat{1.1, 2.2, 3.3}), 4.4, false, TestCustomFloat(4.4)},
		{"Not OneOf invalid", FloatLike[TestCustomFloat]().Not().OneOf([]TestCustomFloat{1.1, 2.2, 3.3}), 2.2, true, TestCustomFloat(2.2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomFloat
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

func TestFloatLikeTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.Transform[*TestCustomFloat]
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Double value", 5.5, func(val *TestCustomFloat, ctx Ctx) error {
			*val = *val * 2
			return nil
		}, false, TestCustomFloat(11.0)},
		{"No change", 42.5, func(val *TestCustomFloat, ctx Ctx) error { return nil }, false, TestCustomFloat(42.5)},
		{"Invalid transform", 42.5, func(val *TestCustomFloat, ctx Ctx) error { return fmt.Errorf("fail") }, true, TestCustomFloat(42.5)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatProc := FloatLike[TestCustomFloat]().Transform(test.transform)
			var result TestCustomFloat
			errs := floatProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestFloatLikeCustomTest(t *testing.T) {
	validator := FloatLike[TestCustomFloat]().TestFunc(func(val *TestCustomFloat, ctx Ctx) bool {
		return *val > 0 // Positive numbers only
	}, Message("must be positive"))
	tests := []struct {
		name      string
		input     float64
		expectErr bool
		expected  TestCustomFloat
	}{
		{"valid positive value", 4.5, false, TestCustomFloat(4.5)},
		{"invalid negative value", -5.5, true, TestCustomFloat(-5.5)},
		{"invalid zero value", 0.0, true, TestCustomFloat(0.0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest TestCustomFloat
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

func TestFloatLikeGetType(t *testing.T) {
	s := FloatLike[TestCustomFloat]()
	assert.Equal(t, zconst.TypeNumber, s.getType())
}

// Validation tests
func TestFloatLikeValidate(t *testing.T) {
	tests := []struct {
		name string
		data TestCustomFloat
	}{
		{"Valid positive value", TestCustomFloat(42.5)},
		{"Valid zero value", TestCustomFloat(0.0)},
		{"Valid negative value", TestCustomFloat(-10.5)},
	}
	floatProc := FloatLike[TestCustomFloat]()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := floatProc.Validate(&test.data)
			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestFloatLikeValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		floatProc := FloatLike[TestCustomFloat]()
		var result TestCustomFloat
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := floatProc.Validate(&result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestFloatLikeValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomFloat
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid positive value", TestCustomFloat(42.5), false, TestCustomFloat(42.5)},
		{"Valid negative value", TestCustomFloat(-10.5), false, TestCustomFloat(-10.5)},
		{"Zero value", TestCustomFloat(0.0), true, TestCustomFloat(0.0)},
	}
	floatProc := FloatLike[TestCustomFloat]().Required()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := floatProc.Validate(&test.data)
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

func TestFloatLikeValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomFloat
		expected  TestCustomFloat
		proc      *NumberSchema[TestCustomFloat]
		expectErr bool
	}{
		{"Optional by default", TestCustomFloat(0.0), TestCustomFloat(0.0), FloatLike[TestCustomFloat](), false},
		{"Optional overrides required", TestCustomFloat(0.0), TestCustomFloat(0.0), FloatLike[TestCustomFloat]().Required().Optional(), false},
		{"Required errors on zero", TestCustomFloat(0.0), TestCustomFloat(0.0), FloatLike[TestCustomFloat]().Required(), true},
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

func TestFloatLikeValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomFloat
		default_  TestCustomFloat
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Valid value", TestCustomFloat(42.5), TestCustomFloat(100.5), false, TestCustomFloat(42.5)},
		{"Zero value with default", TestCustomFloat(0.0), TestCustomFloat(100.5), false, TestCustomFloat(100.5)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatProc := FloatLike[TestCustomFloat]().Default(test.default_)
			errs := floatProc.Validate(&test.data)
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

func TestFloatLikeValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     TestCustomFloat
		catch    TestCustomFloat
		expected TestCustomFloat
	}{
		{"Without catch", TestCustomFloat(42.5), TestCustomFloat(999.999), TestCustomFloat(42.5)},
		{"With catch", TestCustomFloat(0.0), TestCustomFloat(999.999), TestCustomFloat(999.999)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatProc := FloatLike[TestCustomFloat]().TestFunc(func(val *TestCustomFloat, ctx Ctx) bool {
				return *val != 0.0
			}).Catch(test.catch).Required()
			errs := floatProc.Validate(&test.data)
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if test.data != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, test.data)
			}
		})
	}
}

func TestFloatLikeValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomFloat
		transform p.Transform[*TestCustomFloat]
		expectErr bool
		expected  TestCustomFloat
	}{
		{"Double value", TestCustomFloat(5.5), func(val *TestCustomFloat, ctx Ctx) error {
			*val = *val * 2
			return nil
		}, false, TestCustomFloat(11.0)},
		{"No change", TestCustomFloat(42.5), func(val *TestCustomFloat, ctx Ctx) error { return nil }, false, TestCustomFloat(42.5)},
		{"Invalid transform", TestCustomFloat(42.5), func(val *TestCustomFloat, ctx Ctx) error { return fmt.Errorf("fail") }, true, TestCustomFloat(42.5)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatProc := FloatLike[TestCustomFloat]().Transform(test.transform)
			errs := floatProc.Validate(&test.data)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if test.data != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, test.data)
			}
		})
	}
}

func TestFloatLikeValidateCustomTest(t *testing.T) {
	validator := FloatLike[TestCustomFloat]().TestFunc(func(val *TestCustomFloat, ctx Ctx) bool {
		assert.Equal(t, TestCustomFloat(42.5), *val)
		return true
	}, Message("custom"))
	dest := TestCustomFloat(42.5)
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, TestCustomFloat(42.5), dest)
}
