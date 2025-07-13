package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Custom integer type for testing
type TestCustomInt int

func TestIntLikeParse(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid integer value", 42, false, TestCustomInt(42)},
		{"Valid zero value", 0, false, TestCustomInt(0)},
		{"Valid negative value", -10, false, TestCustomInt(-10)},
		{"Valid string number", "123", false, TestCustomInt(123)},
		{"Invalid type (bool)", true, true, TestCustomInt(1)},
		{"Invalid string", "abc", true, TestCustomInt(0)},
	}

	intProc := IntLike[TestCustomInt]()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
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

func TestIntLikeSchemaOption(t *testing.T) {
	s := IntLike[TestCustomInt](WithCoercer(func(original any) (value any, err error) {
		return TestCustomInt(999), nil
	}))
	var result TestCustomInt
	err := s.Parse("invalid", &result)
	assert.Nil(t, err)
	assert.Equal(t, TestCustomInt(999), result)
}

func TestIntLikeExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		intProc := IntLike[TestCustomInt]()
		var result TestCustomInt
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := intProc.Parse(42, &result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestIntLikeRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid integer value", 42, false, TestCustomInt(42)},
		{"Valid zero value", 0, false, TestCustomInt(0)},
		{"Nil value", nil, true, TestCustomInt(0)},
	}
	intProc := IntLike[TestCustomInt]().Required(Message("test"))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
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

func TestIntLikeOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid integer value", 42, false, TestCustomInt(42)},
		{"Valid zero value", 0, false, TestCustomInt(0)},
		{"Nil value", nil, false, TestCustomInt(0)},
	}
	intProc := IntLike[TestCustomInt]().Optional()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
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

func TestIntLikeDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		default_  TestCustomInt
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid integer value", 42, TestCustomInt(100), false, TestCustomInt(42)},
		{"Nil value with default", nil, TestCustomInt(100), false, TestCustomInt(100)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intProc := IntLike[TestCustomInt]().Default(test.default_)
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
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

func TestIntLikeCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     TestCustomInt
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid integer value", 42, TestCustomInt(999), false, TestCustomInt(42)},
		{"Invalid type with catch", "invalid", TestCustomInt(999), false, TestCustomInt(999)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intProc := IntLike[TestCustomInt]().Catch(test.catch)
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
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

func TestIntLikeOneOf(t *testing.T) {
	enum := []TestCustomInt{1, 2, 3, 5, 8}
	intProc := IntLike[TestCustomInt]().OneOf(enum)
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid enum value", 3, false, TestCustomInt(3)},
		{"Invalid enum value", 4, true, TestCustomInt(4)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestIntLikeComparisons(t *testing.T) {
	tests := []struct {
		name      string
		schema    *NumberSchema[TestCustomInt]
		data      any
		expectErr bool
		expected  TestCustomInt
	}{
		{"GT valid", IntLike[TestCustomInt]().GT(5), 10, false, TestCustomInt(10)},
		{"GT invalid", IntLike[TestCustomInt]().GT(5), 3, true, TestCustomInt(3)},
		{"LT valid", IntLike[TestCustomInt]().LT(10), 5, false, TestCustomInt(5)},
		{"LT invalid", IntLike[TestCustomInt]().LT(10), 15, true, TestCustomInt(15)},
		{"GTE valid equal", IntLike[TestCustomInt]().GTE(5), 5, false, TestCustomInt(5)},
		{"GTE valid greater", IntLike[TestCustomInt]().GTE(5), 10, false, TestCustomInt(10)},
		{"GTE invalid", IntLike[TestCustomInt]().GTE(5), 3, true, TestCustomInt(3)},
		{"LTE valid equal", IntLike[TestCustomInt]().LTE(10), 10, false, TestCustomInt(10)},
		{"LTE valid less", IntLike[TestCustomInt]().LTE(10), 5, false, TestCustomInt(5)},
		{"LTE invalid", IntLike[TestCustomInt]().LTE(10), 15, true, TestCustomInt(15)},
		{"EQ valid", IntLike[TestCustomInt]().EQ(42), 42, false, TestCustomInt(42)},
		{"EQ invalid", IntLike[TestCustomInt]().EQ(42), 43, true, TestCustomInt(43)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomInt
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

func TestIntLikeNot(t *testing.T) {
	tests := []struct {
		name      string
		schema    *NumberSchema[TestCustomInt]
		data      any
		expectErr bool
		expected  TestCustomInt
	}{
		{"Not EQ valid", IntLike[TestCustomInt]().Not().EQ(42), 43, false, TestCustomInt(43)},
		{"Not EQ invalid", IntLike[TestCustomInt]().Not().EQ(42), 42, true, TestCustomInt(42)},
		{"Not OneOf valid", IntLike[TestCustomInt]().Not().OneOf([]TestCustomInt{1, 2, 3}), 4, false, TestCustomInt(4)},
		{"Not OneOf invalid", IntLike[TestCustomInt]().Not().OneOf([]TestCustomInt{1, 2, 3}), 2, true, TestCustomInt(2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result TestCustomInt
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

func TestIntLikeTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.Transform[*TestCustomInt]
		expectErr bool
		expected  TestCustomInt
	}{
		{"Double value", 5, func(val *TestCustomInt, ctx Ctx) error {
			*val = *val * 2
			return nil
		}, false, TestCustomInt(10)},
		{"No change", 42, func(val *TestCustomInt, ctx Ctx) error { return nil }, false, TestCustomInt(42)},
		{"Invalid transform", 42, func(val *TestCustomInt, ctx Ctx) error { return fmt.Errorf("fail") }, true, TestCustomInt(42)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intProc := IntLike[TestCustomInt]().Transform(test.transform)
			var result TestCustomInt
			errs := intProc.Parse(test.data, &result)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestIntLikeCustomTest(t *testing.T) {
	validator := IntLike[TestCustomInt]().TestFunc(func(val *TestCustomInt, ctx Ctx) bool {
		return *val%2 == 0 // Even numbers only
	}, Message("must be even"))
	tests := []struct {
		name      string
		input     int
		expectErr bool
		expected  TestCustomInt
	}{
		{"valid even value", 4, false, TestCustomInt(4)},
		{"invalid odd value", 5, true, TestCustomInt(5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest TestCustomInt
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

func TestIntLikeGetType(t *testing.T) {
	s := IntLike[TestCustomInt]()
	assert.Equal(t, zconst.TypeNumber, s.getType())
}

// Validation tests
func TestIntLikeValidate(t *testing.T) {
	tests := []struct {
		name string
		data TestCustomInt
	}{
		{"Valid positive value", TestCustomInt(42)},
		{"Valid zero value", TestCustomInt(0)},
		{"Valid negative value", TestCustomInt(-10)},
	}
	intProc := IntLike[TestCustomInt]()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := intProc.Validate(&test.data)
			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestIntLikeValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		intProc := IntLike[TestCustomInt]()
		var result TestCustomInt
		var contextPassed bool
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}
		errs := intProc.Validate(&result, fakeOption)
		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}
		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestIntLikeValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomInt
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid positive value", TestCustomInt(42), false, TestCustomInt(42)},
		{"Valid negative value", TestCustomInt(-10), false, TestCustomInt(-10)},
		{"Zero value", TestCustomInt(0), true, TestCustomInt(0)},
	}
	intProc := IntLike[TestCustomInt]().Required()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := intProc.Validate(&test.data)
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

func TestIntLikeValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomInt
		expected  TestCustomInt
		proc      *NumberSchema[TestCustomInt]
		expectErr bool
	}{
		{"Optional by default", TestCustomInt(0), TestCustomInt(0), IntLike[TestCustomInt](), false},
		{"Optional overrides required", TestCustomInt(0), TestCustomInt(0), IntLike[TestCustomInt]().Required().Optional(), false},
		{"Required errors on zero", TestCustomInt(0), TestCustomInt(0), IntLike[TestCustomInt]().Required(), true},
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

func TestIntLikeValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomInt
		default_  TestCustomInt
		expectErr bool
		expected  TestCustomInt
	}{
		{"Valid value", TestCustomInt(42), TestCustomInt(100), false, TestCustomInt(42)},
		{"Zero value with default", TestCustomInt(0), TestCustomInt(100), false, TestCustomInt(100)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intProc := IntLike[TestCustomInt]().Default(test.default_)
			errs := intProc.Validate(&test.data)
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

func TestIntLikeValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     TestCustomInt
		catch    TestCustomInt
		expected TestCustomInt
	}{
		{"Without catch", TestCustomInt(42), TestCustomInt(999), TestCustomInt(42)},
		{"With catch", TestCustomInt(0), TestCustomInt(999), TestCustomInt(999)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intProc := IntLike[TestCustomInt]().TestFunc(func(val *TestCustomInt, ctx Ctx) bool {
				return *val != 0
			}).Catch(test.catch).Required()
			errs := intProc.Validate(&test.data)
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
			if test.data != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, test.data)
			}
		})
	}
}

func TestIntLikeValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      TestCustomInt
		transform p.Transform[*TestCustomInt]
		expectErr bool
		expected  TestCustomInt
	}{
		{"Double value", TestCustomInt(5), func(val *TestCustomInt, ctx Ctx) error {
			*val = *val * 2
			return nil
		}, false, TestCustomInt(10)},
		{"No change", TestCustomInt(42), func(val *TestCustomInt, ctx Ctx) error { return nil }, false, TestCustomInt(42)},
		{"Invalid transform", TestCustomInt(42), func(val *TestCustomInt, ctx Ctx) error { return fmt.Errorf("fail") }, true, TestCustomInt(42)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intProc := IntLike[TestCustomInt]().Transform(test.transform)
			errs := intProc.Validate(&test.data)
			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if test.data != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, test.data)
			}
		})
	}
}

func TestIntLikeValidateCustomTest(t *testing.T) {
	validator := IntLike[TestCustomInt]().TestFunc(func(val *TestCustomInt, ctx Ctx) bool {
		assert.Equal(t, TestCustomInt(42), *val)
		return true
	}, Message("custom"))
	dest := TestCustomInt(42)
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, TestCustomInt(42), dest)
}
