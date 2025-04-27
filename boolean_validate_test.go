package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/stretchr/testify/assert"
)

func TestBoolValidate(t *testing.T) {
	tests := []struct {
		name string
		data bool
	}{
		{
			name: "Valid true value",
			data: true,
		},
		{
			name: "Valid false value",
			data: false,
		},
	}

	boolProc := Bool()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := boolProc.Validate(&test.data)

			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestBoolValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		boolProc := Bool()
		var result bool
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

func TestBoolValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      bool
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: true,
		},
		{
			name:      "Valid false value",
			data:      false,
			expected:  false,
			expectErr: true,
		},
	}

	boolProc := Bool().Required()

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

func TestBoolValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      bool
		expected  bool
		proc      *BoolSchema[bool]
		expectErr bool
	}{
		{
			name:     "Optiona by default",
			data:     false,
			expected: false,
			proc:     Bool(),
		},
		{
			name:     "Optional Overrides Required",
			data:     false,
			expected: false,
			proc:     Bool().Required().Optional(),
		},
		{
			name:      "required errors on zero value",
			data:      false,
			expected:  false,
			proc:      Bool().Required(),
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

func TestBoolValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      bool
		default_  bool
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			default_: false,
			expected: true,
		},
		{
			name:     "Valid false value",
			data:     false,
			default_: true,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().Default(test.default_)
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

func TestBoolValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     bool
		catch    bool
		expected bool
	}{
		{
			name:     "Without catch",
			data:     true,
			expected: true,
		},
		{
			name:     "With Catch",
			data:     false,
			catch:    true,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().TestFunc(func(val *bool, ctx Ctx) bool {
				return *val == true
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

func TestBoolValidateTrue(t *testing.T) {
	tests := []struct {
		name      string
		data      bool
		expectErr bool
	}{
		{
			name: "Valid true value",
			data: true,
		},
		{
			name:      "Invalid false value",
			data:      false,
			expectErr: true,
		},
	}

	boolProc := Bool().True().Required()

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

func TestBoolValidateFalse(t *testing.T) {
	tests := []struct {
		name      string
		data      bool
		expectErr bool
	}{
		{
			name: "Valid false value",
			data: false,
		},
		{
			name:      "Invalid true value",
			data:      true,
			expectErr: true,
		},
	}

	boolProc := Bool().False()

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

func TestBoolValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      bool
		transform p.Transform[*bool]
		expectErr bool
		expected  bool
	}{
		{
			name: "Invert boolean",
			data: true,
			transform: func(val *bool, ctx Ctx) error {
				*val = !*val
				return nil
			},
			expected: false,
		},
		{
			name: "No change",
			data: false,
			transform: func(val *bool, ctx Ctx) error {
				return nil
			},
			expected: false,
		},
		{
			name: "Invalid transform",
			data: true,
			transform: func(val *bool, ctx Ctx) error {
				return fmt.Errorf("invalid operation")
			},
			expectErr: true,
			expected:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().Transform(test.transform)
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

func TestBoolValidateCustomTest(t *testing.T) {
	validator := Bool().TestFunc(func(val *bool, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, true, *val)
		return true
	}, Message("custom"))
	dest := true
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, true, dest)
}
