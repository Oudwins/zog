package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/stretchr/testify/assert"
)

func TestNumberValidate(t *testing.T) {
	tests := []struct {
		name      string
		data      float64
		expectErr bool
		expected  float64
	}{
		{
			name:     "Valid value",
			data:     5,
			expected: 5,
		},
		{
			name:     "Zero value",
			data:     0,
			expected: 0,
		},
	}

	schema := Float()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errs := schema.Validate(&test.data)
			assert.Empty(t, errs)
			assert.Equal(t, test.expected, test.data)
		})
	}
}

func TestNumberValidateOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		schema := Int()
		var contextPassed bool

		// Create a fake parsing option that checks if it receives a ParseCtx
		fakeOption := func(p *p.ZogParseCtx) {
			if p != nil {
				contextPassed = true
			}
		}

		data := 5
		errs := schema.Validate(&data, fakeOption)

		assert.Empty(t, errs)
		assert.True(t, contextPassed)
	})
}

func TestNumberValidateRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      int
		expectErr bool
		expected  int
	}{
		{
			name:     "Valid value",
			data:     5,
			expected: 5,
		},
		{
			name:      "Zero value",
			data:      0,
			expectErr: true,
			expected:  0,
		},
	}

	schema := Int().Required()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result int
			errs := schema.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if !test.expectErr && result != test.expected {
				t.Errorf("On Run %s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestNumberValidateOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		proc      *numberProcessor[int]
		expectErr bool
		expected  int
	}{
		{
			name:     "Optional by default",
			data:     0,
			proc:     Int(),
			expected: 0,
		},
		{
			name:     "Optional overrides Required",
			data:     0,
			proc:     Int().Required().Optional(),
			expected: 0,
		},
		{
			name:      "Required errors on zero value",
			data:      0,
			proc:      Int().Required(),
			expectErr: true,
			expected:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result int
			errs := test.proc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestNumberValidateDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		default_  int
		expectErr bool
		expected  int
	}{
		{
			name:     "Valid value",
			data:     5,
			default_: 10,
			expected: 5,
		},
		{
			name:     "Nil value with default",
			data:     nil,
			default_: 10,
			expected: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schema := Int().Default(test.default_)
			var result int
			errs := schema.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestNumberValidateCatch(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		catch    int
		expected int
	}{
		{
			name:     "Without catch",
			data:     5,
			expected: 5,
		},
		{
			name:     "With Catch",
			data:     "invalid",
			catch:    10,
			expected: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schema := Int().Catch(test.catch).PreTransform(func(val any, ctx ParseCtx) (any, error) {
				if _, ok := val.(string); ok {
					return nil, fmt.Errorf("invalid input")
				}
				return val, nil
			})
			var result int
			_ = schema.Parse(test.data, &result)

			if result != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestNumberValidatePreTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.PreTransform
		expectErr bool
		expected  int
	}{
		{
			name: "Valid transform",
			data: 5,
			transform: func(val any, ctx ParseCtx) (any, error) {
				if v, ok := val.(int); ok {
					return v * 2, nil
				}
				return val, nil
			},
			expected: 10,
		},
		{
			name: "Invalid transform",
			data: "invalid",
			transform: func(val any, ctx ParseCtx) (any, error) {
				return nil, fmt.Errorf("invalid input")
			},
			expectErr: true,
			expected:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schema := Int().PreTransform(test.transform)
			var result int
			errs := schema.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestNumberValidatePostTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.PostTransform
		expectErr bool
		expected  int
	}{
		{
			name: "Add one",
			data: 5,
			transform: func(val any, ctx ParseCtx) error {
				if v, ok := val.(*int); ok {
					*v += 1
				}
				return nil
			},
			expected: 6,
		},
		{
			name: "No change",
			data: 5,
			transform: func(val any, ctx ParseCtx) error {
				return nil
			},
			expected: 5,
		},
		{
			name: "Invalid transform",
			data: 5,
			transform: func(val any, ctx ParseCtx) error {
				return fmt.Errorf("invalid operation")
			},
			expectErr: true,
			expected:  5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schema := Int().PostTransform(test.transform)
			var result int
			errs := schema.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
