package zog

import (
	"testing"
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

	boolProc := Bool().Required()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result bool
			errs := boolProc.Parse(test.data, &result)

			if (len(errs) > 0) != test.expectErr {
				t.Errorf("On Run %s -> Expected error: %v, got: %v", test.name, test.expectErr, errs)
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
