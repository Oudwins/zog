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
