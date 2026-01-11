package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/pkgs/internals/tutils"
	"github.com/stretchr/testify/assert"
)

func TestAnyValidate(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{
			name: "Valid string value",
			data: "hello",
		},
		{
			name: "Valid int value",
			data: 42,
		},
		{
			name: "Valid bool value",
			data: true,
		},
	}

	anyProc := EXPERIMENTAL_ANY()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := anyProc.Validate(&test.data)

			if len(errs) > 0 {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
		})
	}
}

func TestAnyValidateExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		anyProc := EXPERIMENTAL_ANY()
		var result any = "test"
		var contextPassed bool

		// Create a fake parsing option that checks if it receives a Ctx
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}

		errs := anyProc.Validate(&result, fakeOption)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}

		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestAnyValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expectErr bool
		expected  any
	}{
		{
			name:     "Valid string value",
			data:     "hello",
			expected: "hello",
		},
		{
			name:      "Nil value (zero value)",
			data:      nil,
			expected:  nil,
			expectErr: true,
		},
	}

	anyProc := EXPERIMENTAL_ANY().Required()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := anyProc.Validate(&test.data)
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

func TestAnyValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		expected  any
		proc      *AnySchema
		expectErr bool
	}{
		{
			name:     "Optional by default",
			data:     nil,
			expected: nil,
			proc:     EXPERIMENTAL_ANY(),
		},
		{
			name:     "Optional Overrides Required",
			data:     nil,
			expected: nil,
			proc:     EXPERIMENTAL_ANY().Required().Optional(),
		},
		{
			name:      "required errors on zero value",
			data:      nil,
			expected:  nil,
			proc:      EXPERIMENTAL_ANY().Required(),
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

func TestAnyValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		default_  any
		expectErr bool
		expected  any
	}{
		{
			name:     "Valid string value",
			data:     "hello",
			default_: "default",
			expected: "hello",
		},
		{
			name:     "Nil value with default",
			data:     nil,
			default_: "default",
			expected: "default",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			anyProc := EXPERIMENTAL_ANY().Default(test.default_)
			errs := anyProc.Validate(&test.data)

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

func TestAnyValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		catch    any
		expected any
	}{
		{
			name:     "Without catch",
			data:     "hello",
			expected: "hello",
		},
		{
			name:     "With Catch",
			data:     nil,
			catch:    "catch",
			expected: "catch",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			anyProc := EXPERIMENTAL_ANY().TestFunc(func(val *any, ctx Ctx) bool {
				return *val != nil
			}).Catch(test.catch).Required()
			errs := anyProc.Validate(&test.data)
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, test.data) {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, test.data)
			}
		})
	}
}

func TestAnyValidateTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		transform p.Transform[*any]
		expectErr bool
		expected  any
	}{
		{
			name: "Transform string to uppercase",
			data: "hello",
			transform: func(val *any, ctx Ctx) error {
				if _, ok := (*val).(string); ok {
					*val = "HELLO"
				}
				return nil
			},
			expected: "HELLO",
		},
		{
			name: "No change",
			data: 42,
			transform: func(val *any, ctx Ctx) error {
				return nil
			},
			expected: 42,
		},
		{
			name: "Invalid transform",
			data: "test",
			transform: func(val *any, ctx Ctx) error {
				return fmt.Errorf("invalid operation")
			},
			expectErr: true,
			expected:  "test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			anyProc := EXPERIMENTAL_ANY().Transform(test.transform)
			errs := anyProc.Validate(&test.data)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, test.data) {
				t.Errorf("Expected %v, but got %v", test.expected, test.data)
			}
		})
	}
}

func TestAnyValidateCustomTest(t *testing.T) {
	validator := EXPERIMENTAL_ANY().TestFunc(func(val *any, ctx Ctx) bool {
		// Custom test logic - check if value is a string
		_, ok := (*val).(string)
		return ok
	}, Message("custom"))
	dest := any("test")
	errs := validator.Validate(&dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, "test", dest)
}
