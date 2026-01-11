package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/pkgs/internals"
	"github.com/Oudwins/zog/pkgs/internals/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestAnyParse(t *testing.T) {
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
			name:     "Valid int value",
			data:     42,
			expected: 42,
		},
		{
			name:     "Valid bool value",
			data:     true,
			expected: true,
		},
		{
			name:     "Valid map value",
			data:     map[string]int{"key": 1},
			expected: map[string]int{"key": 1},
		},
		{
			name:     "Valid slice value",
			data:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Nil value",
			data:     nil,
			expected: nil,
		},
	}

	anyProc := Any()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result any
			errs := anyProc.Parse(test.data, &result)

			if len(errs) > 0 && !test.expectErr {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}
			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, result) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestAnySchemaOption(t *testing.T) {
	s := Any(WithCoercer(func(original any) (value any, err error) {
		return "coerced", nil
	}))

	var result any
	err := s.Parse("asdasdas", &result)
	assert.Nil(t, err)
	// Note: WithCoercer is a no-op for Any schema, so it should just pass through
	assert.Equal(t, "asdasdas", result)
}

func TestAnyExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		anyProc := Any()
		var result any
		var contextPassed bool

		// Create a fake parsing option that checks if it receives a Ctx
		fakeOption := func(p *p.ExecCtx) {
			if p != nil {
				contextPassed = true
			}
		}

		errs := anyProc.Parse("test", &result, fakeOption)

		if len(errs) > 0 {
			t.Errorf("Unexpected errors: %v", errs)
		}

		if !contextPassed {
			t.Error("Parse context was not passed to the parsing option")
		}
	})
}

func TestAnyRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  any
	}{
		{
			name:     "Valid string value",
			data:     "hello",
			expected: "hello",
		},
		{
			name:     "Valid int value",
			data:     42,
			expected: 42,
		},
		{
			name:      "Nil value",
			data:      nil,
			expectErr: true,
		},
	}

	anyProc := Any().Required(Message("test"))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result any
			errs := anyProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if test.expectErr && len(errs) > 0 && errs[0].Message != "test" {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, "test", errs[0].Message)
			}

			if !test.expectErr && !assert.Equal(t, test.expected, result) {
				t.Errorf("On Run %s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestAnyOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  any
	}{
		{
			name:     "Valid string value",
			data:     "hello",
			expected: "hello",
		},
		{
			name:     "Valid int value",
			data:     42,
			expected: 42,
		},
		{
			name:     "Nil value",
			data:     nil,
			expected: nil,
		},
	}

	anyProc := Any().Optional()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result any
			errs := anyProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, result) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestAnyDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
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
			name:     "Nil value with string default",
			data:     nil,
			default_: "default",
			expected: "default",
		},
		{
			name:     "Nil value with int default",
			data:     nil,
			default_: 42,
			expected: 42,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			anyProc := Any().Default(test.default_)
			var result any
			errs := anyProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, result) {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestAnyCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     any
		expectErr bool
		expected  any
	}{
		{
			name:     "Valid string value",
			data:     "hello",
			catch:    "catch",
			expected: "hello",
		},
		{
			name:      "Nil value with catch (required)",
			data:      nil,
			catch:     "catch",
			expectErr: false,
			expected:  "catch",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			anyProc := Any().Required().Catch(test.catch)
			var result any
			errs := anyProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, result) {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestAnyTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
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
			anyProc := Any().Transform(test.transform)
			var result any
			errs := anyProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if len(errs) > 0 {
				tutils.VerifyDefaultIssueMessages(t, errs)
			}

			if !assert.Equal(t, test.expected, result) {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestAnyCustomTest(t *testing.T) {
	validator := Any().TestFunc(func(val *any, ctx Ctx) bool {
		// Custom test logic - check if value is a string
		_, ok := (*val).(string)
		return ok
	}, Message("custom"))

	tests := []struct {
		name      string
		input     any
		expectErr bool
	}{
		{
			name:      "valid string value",
			input:     "hello",
			expectErr: false,
		},
		{
			name:      "invalid int value",
			input:     42,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest any
			errs := validator.Parse(tt.input, &dest)
			if (len(errs) > 0) != tt.expectErr {
				t.Errorf("got errors %v, expectErr %v", errs, tt.expectErr)
			}
			if !tt.expectErr {
				assert.Equal(t, tt.input, dest)
			}
		})
	}
}

func TestAnyGetType(t *testing.T) {
	s := Any()
	assert.Equal(t, zconst.TypeAny, s.getType())
}
