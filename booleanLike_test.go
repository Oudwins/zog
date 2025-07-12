package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Custom boolean type for testing
type CustomBool bool

func TestBoolLikeParse(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: CustomBool(true),
		},
		{
			name:     "Valid false value",
			data:     false,
			expected: CustomBool(false),
		},
		{
			name:      "Invalid value",
			data:      "invalid",
			expectErr: true,
			expected:  CustomBool(false), // Since it's an invalid value, it should default to false
		},
	}

	boolProc := BoolLike[CustomBool]()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeSchemaOption(t *testing.T) {
	s := BoolLike[CustomBool](WithCoercer(func(original any) (value any, err error) {
		return CustomBool(true), nil
	}))

	var result CustomBool
	err := s.Parse("asdasdas", &result)
	assert.Nil(t, err)
	assert.Equal(t, CustomBool(true), result)
}

func TestBoolLikeExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		boolProc := BoolLike[CustomBool]()
		var result CustomBool
		var contextPassed bool

		// Create a fake parsing option that checks if it receives a Ctx
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}

		errs := boolProc.Parse(true, &result, fakeOption)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}

		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestBoolLikeRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: CustomBool(true),
		},
		{
			name:     "Valid false value",
			data:     false,
			expected: CustomBool(false),
		},
		{
			name:      "Nil value",
			data:      nil,
			expectErr: true,
		},
	}

	boolProc := BoolLike[CustomBool]().Required(Message("test"))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: CustomBool(true),
		},
		{
			name:     "Valid false value",
			data:     false,
			expected: CustomBool(false),
		},
		{
			name:     "Nil value",
			data:     nil,
			expected: CustomBool(false), // Default value for CustomBool
		},
	}

	boolProc := BoolLike[CustomBool]().Optional()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		default_  CustomBool
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     true,
			default_: CustomBool(false),
			expected: CustomBool(true),
		},
		{
			name:     "Valid false value",
			data:     false,
			default_: CustomBool(true),
			expected: CustomBool(false),
		},
		{
			name:     "Nil value with true default",
			data:     nil,
			default_: CustomBool(true),
			expected: CustomBool(true),
		},
		{
			name:     "Nil value with false default",
			data:     nil,
			default_: CustomBool(false),
			expected: CustomBool(false),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().Default(test.default_)
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     CustomBool
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     true,
			catch:    CustomBool(false),
			expected: CustomBool(true),
		},
		{
			name:     "Valid false value",
			data:     false,
			catch:    CustomBool(true),
			expected: CustomBool(false),
		},
		{
			name:      "Invalid value with true catch",
			data:      "invalid",
			catch:     CustomBool(true),
			expectErr: false,
			expected:  CustomBool(true),
		},
		{
			name:      "Invalid value with false catch",
			data:      "invalid",
			catch:     CustomBool(false),
			expectErr: false,
			expected:  CustomBool(false),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().Catch(test.catch)
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeTrue(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: CustomBool(true),
		},
		{
			name:      "Invalid false value",
			data:      false,
			expectErr: true,
			expected:  CustomBool(false),
		},
	}

	boolProc := BoolLike[CustomBool]().True()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeFalse(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid false value",
			data:     false,
			expected: CustomBool(false),
		},
		{
			name:      "Invalid true value",
			data:      true,
			expectErr: true,
			expected:  CustomBool(true),
		},
	}

	boolProc := BoolLike[CustomBool]().False()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.Transform[*CustomBool]
		expectErr bool
		expected  CustomBool
	}{
		{
			name: "Invert boolean",
			data: true,
			transform: func(val *CustomBool, ctx Ctx) error {
				*val = !*val
				return nil
			},
			expected: CustomBool(false),
		},
		{
			name: "No change",
			data: false,
			transform: func(val *CustomBool, ctx Ctx) error {
				return nil
			},
			expected: CustomBool(false),
		},
		{
			name: "Invalid transform",
			data: true,
			transform: func(val *CustomBool, ctx Ctx) error {
				return fmt.Errorf("invalid operation")
			},
			expectErr: true,
			expected:  CustomBool(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().Transform(test.transform)
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeCustomTest(t *testing.T) {
	validator := BoolLike[CustomBool]().TestFunc(func(val *CustomBool, ctx Ctx) bool {
		// Custom test logic here
		return *val == CustomBool(true)
	}, Message("custom"))

	tests := []struct {
		name      string
		input     bool
		expectErr bool
		expected  CustomBool
	}{
		{
			name:      "valid true value",
			input:     true,
			expectErr: false,
			expected:  CustomBool(true),
		},
		{
			name:      "invalid false value",
			input:     false,
			expectErr: true,
			expected:  CustomBool(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest CustomBool
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

func TestBoolLikeEQ(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		eqValue   CustomBool
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Equal true value",
			data:     true,
			eqValue:  CustomBool(true),
			expected: CustomBool(true),
		},
		{
			name:     "Equal false value",
			data:     false,
			eqValue:  CustomBool(false),
			expected: CustomBool(false),
		},
		{
			name:      "Not equal value",
			data:      true,
			eqValue:   CustomBool(false),
			expectErr: true,
			expected:  CustomBool(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().EQ(test.eqValue)
			var result CustomBool
			errs := boolProc.Parse(test.data, &result)

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

func TestBoolLikeGetType(t *testing.T) {
	s := BoolLike[CustomBool]()
	assert.Equal(t, zconst.TypeBool, s.getType())
}

// Validation tests

func TestBoolLikeValidate(t *testing.T) {
	tests := []struct {
		name string
		data CustomBool
	}{
		{
			name: "Valid true value",
			data: CustomBool(true),
		},
		{
			name: "Valid false value",
			data: CustomBool(false),
		},
	}

	boolProc := BoolLike[CustomBool]()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := boolProc.Validate(&test.data)

			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestBoolLikeValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		boolProc := BoolLike[CustomBool]()
		var result CustomBool
		var contextPassed bool

		// Create a fake parsing option that checks if it receives a Ctx
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}

		errs := boolProc.Validate(&result, fakeOption)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}

		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestBoolLikeValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomBool
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     CustomBool(true),
			expected: CustomBool(true),
		},
		{
			name:      "Valid false value",
			data:      CustomBool(false),
			expected:  CustomBool(false),
			expectErr: true,
		},
	}

	boolProc := BoolLike[CustomBool]().Required()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := boolProc.Validate(&test.data)
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

func TestBoolLikeValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomBool
		expected  CustomBool
		proc      *BoolSchema[CustomBool]
		expectErr bool
	}{
		{
			name:     "Optional by default",
			data:     CustomBool(false),
			expected: CustomBool(false),
			proc:     BoolLike[CustomBool](),
		},
		{
			name:     "Optional Overrides Required",
			data:     CustomBool(false),
			expected: CustomBool(false),
			proc:     BoolLike[CustomBool]().Required().Optional(),
		},
		{
			name:      "required errors on zero value",
			data:      CustomBool(false),
			expected:  CustomBool(false),
			proc:      BoolLike[CustomBool]().Required(),
			expectErr: true,
		},
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

func TestBoolLikeValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomBool
		default_  CustomBool
		expectErr bool
		expected  CustomBool
	}{
		{
			name:     "Valid true value",
			data:     CustomBool(true),
			default_: CustomBool(false),
			expected: CustomBool(true),
		},
		{
			name:     "Valid false value",
			data:     CustomBool(false),
			default_: CustomBool(true),
			expected: CustomBool(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().Default(test.default_)
			errs := boolProc.Validate(&test.data)

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

func TestBoolLikeValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     CustomBool
		catch    CustomBool
		expected CustomBool
	}{
		{
			name:     "Without catch",
			data:     CustomBool(true),
			expected: CustomBool(true),
		},
		{
			name:     "With Catch",
			data:     CustomBool(false),
			catch:    CustomBool(true),
			expected: CustomBool(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().TestFunc(func(val *CustomBool, ctx Ctx) bool {
				return *val == CustomBool(true)
			}).Catch(test.catch).Required()
			errs := boolProc.Validate(&test.data)
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if test.data != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, test.data)
			}
		})
	}
}

func TestBoolLikeValidateTrue(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomBool
		expectErr bool
	}{
		{
			name: "Valid true value",
			data: CustomBool(true),
		},
		{
			name:      "Invalid false value",
			data:      CustomBool(false),
			expectErr: true,
		},
	}

	boolProc := BoolLike[CustomBool]().True().Required()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := boolProc.Validate(&test.data)
			if test.expectErr {
				assert.NotEmpty(t, errs)
				tutils.VerifyDefaultIssueMessages(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestBoolLikeValidateFalse(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomBool
		expectErr bool
	}{
		{
			name: "Valid false value",
			data: CustomBool(false),
		},
		{
			name:      "Invalid true value",
			data:      CustomBool(true),
			expectErr: true,
		},
	}

	boolProc := BoolLike[CustomBool]().False()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := boolProc.Validate(&test.data)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}
		})
	}
}

func TestBoolLikeValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      CustomBool
		transform p.Transform[*CustomBool]
		expectErr bool
		expected  CustomBool
	}{
		{
			name: "Invert boolean",
			data: CustomBool(true),
			transform: func(val *CustomBool, ctx Ctx) error {
				*val = !*val
				return nil
			},
			expected: CustomBool(false),
		},
		{
			name: "No change",
			data: CustomBool(false),
			transform: func(val *CustomBool, ctx Ctx) error {
				return nil
			},
			expected: CustomBool(false),
		},
		{
			name: "Invalid transform",
			data: CustomBool(true),
			transform: func(val *CustomBool, ctx Ctx) error {
				return fmt.Errorf("invalid operation")
			},
			expectErr: true,
			expected:  CustomBool(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := BoolLike[CustomBool]().Transform(test.transform)
			errs := boolProc.Validate(&test.data)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if test.data != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, test.data)
			}
		})
	}
}

func TestBoolLikeValidateCustomTest(t *testing.T) {
	validator := BoolLike[CustomBool]().TestFunc(func(val *CustomBool, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, CustomBool(true), *val)
		return true
	}, Message("custom"))
	dest := CustomBool(true)
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, CustomBool(true), dest)
}
