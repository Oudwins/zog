package zog

import (
	"fmt"
	"testing"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestBoolParse(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: true,
		},
		{
			name:     "Valid false value",
			data:     false,
			expected: false,
		},
		{
			name:      "Invalid value",
			data:      "invalid",
			expectErr: true,
			expected:  false, // Since it's an invalid value, it should default to false
		},
	}

	boolProc := Bool()

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if len(errs) > 0 && !test.expectErr {
				t.Errorf("Unexpected errors i = %d: %v", i, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestBoolSchemaOption(t *testing.T) {
	s := Bool(WithCoercer(func(original any) (value any, err error) {
		return true, nil
	}))

	var result bool
	err := s.Parse("asdasdas", &result)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

func TestExecOption(t *testing.T) {
	t.Run("Parse context is passed to parsing option", func(t *testing.T) {
		boolProc := Bool()
		var result bool
		var contextPassed bool

		// Create a fake parsing option that checks if it receives a ParseCtx
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

func TestBoolRequired(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: true,
		},
		{
			name:     "Valid false value",
			data:     false,
			expected: false,
		},
		{
			name:      "Nil value",
			data:      nil,
			expectErr: true,
		},
	}

	boolProc := Bool().Required(Message("test"))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if test.expectErr && errs[0].Message() != "test" {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, "test", errs[0].Message())
			}

			if !test.expectErr && result != test.expected {
				t.Errorf("On Run %s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestBoolOptional(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: true,
		},
		{
			name:     "Valid false value",
			data:     false,
			expected: false,
		},
		{
			name:     "Nil value",
			data:     nil,
			expected: false, // Default value for bool
		},
	}

	boolProc := Bool().Optional()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestBoolDefault(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
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
			expected: false,
		},
		{
			name:     "Nil value with true default",
			data:     nil,
			default_: true,
			expected: true,
		},
		{
			name:     "Nil value with false default",
			data:     nil,
			default_: false,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().Default(test.default_)
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestBoolCatch(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		catch     bool
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			catch:    false,
			expected: true,
		},
		{
			name:     "Valid false value",
			data:     false,
			catch:    true,
			expected: false,
		},
		{
			name:      "Invalid value with true catch",
			data:      "invalid",
			catch:     true,
			expectErr: false,
			expected:  true,
		},
		{
			name:      "Invalid value with false catch",
			data:      "invalid",
			catch:     false,
			expectErr: false,
			expected:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().Catch(test.catch)
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("%s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("%s -> Expected %v, but got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestBoolTrue(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid true value",
			data:     true,
			expected: true,
		},
		{
			name:      "Invalid false value",
			data:      false,
			expectErr: true,
			expected:  false,
		},
	}

	boolProc := Bool().True()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestBoolFalse(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		expectErr bool
		expected  bool
	}{
		{
			name:     "Valid false value",
			data:     false,
			expected: false,
		},
		{
			name:      "Invalid true value",
			data:      true,
			expectErr: true,
			expected:  true,
		},
	}

	boolProc := Bool().False()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
func TestBoolPreTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.PreTransform
		expectErr bool
		expected  bool
	}{
		{
			name: "Valid transform",
			data: "true",
			transform: func(val any, ctx ParseCtx) (any, error) {
				if s, ok := val.(*string); ok {
					return *s == "true", nil
				}
				return val, nil
			},
			expected: true,
		},
		{
			name: "Invalid transform",
			data: "invalid",
			transform: func(val any, ctx ParseCtx) (any, error) {
				return nil, fmt.Errorf("invalid input")
			},
			expectErr: true,
			expected:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().PreTransform(test.transform)
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestBoolPostTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		transform p.PostTransform
		expectErr bool
		expected  bool
	}{
		{
			name: "Invert boolean",
			data: true,
			transform: func(val any, ctx ParseCtx) error {
				if b, ok := val.(*bool); ok {
					*b = !*b
				}
				return nil
			},
			expected: false,
		},
		{
			name: "No change",
			data: false,
			transform: func(val any, ctx ParseCtx) error {
				return nil
			},
			expected: false,
		},
		{
			name: "Invalid transform",
			data: true,
			transform: func(val any, ctx ParseCtx) error {
				return fmt.Errorf("invalid operation")
			},
			expectErr: true,
			expected:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boolProc := Bool().PostTransform(test.transform)
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestBoolCustomTest(t *testing.T) {
	validator := Bool().TestFunc(func(val any, ctx Ctx) bool {
		// Custom test logic here
		assert.Equal(t, true, val)
		return true
	}, Message("custom"))
	dest := false
	errs := validator.Parse(true, &dest)
	if len(errs) > 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}
	assert.Equal(t, true, dest)
}

func TestBoolGetType(t *testing.T) {
	s := Bool()
	assert.Equal(t, zconst.TypeBool, s.getType())
}
