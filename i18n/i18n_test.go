package i18n

import (
	"testing"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// Define test language maps
var enMap = zconst.LangMap{
	zconst.TypeString: {
		zconst.IssueCodeRequired: "is required",
	},
}
var esMap = zconst.LangMap{
	zconst.TypeString: {
		zconst.IssueCodeRequired: "es requerido",
	},
}
var jaMap = zconst.LangMap{
	zconst.TypeString: {
		zconst.IssueCodeRequired: "必須です",
	},
}

func TestSetLanguagesErrsMap(t *testing.T) {

	// Set up language maps
	SetLanguagesErrsMap(map[string]zconst.LangMap{
		"en": enMap,
		"es": esMap,
		"ja": jaMap,
	}, "en")

	// Define a schema for testing
	schema := zog.Struct(zog.Shape{
		"name": zog.String().Required().Min(1),
	})

	// Test cases
	testCases := []struct {
		name        string
		lang        string
		input       map[string]interface{}
		expectedErr bool
		expected    string
	}{
		{
			name:        "English error message",
			lang:        "en",
			input:       map[string]interface{}{},
			expectedErr: true,
			expected:    "is required",
		},
		{
			name:        "Spanish error message",
			lang:        "es",
			input:       map[string]interface{}{},
			expectedErr: true,
			expected:    "es requerido",
		},
		{
			name:        "Japanese error message",
			lang:        "ja",
			input:       map[string]interface{}{},
			expectedErr: true,
			expected:    "必須です",
		},
		{
			name:        "Default to English when language not found",
			lang:        "fr",
			input:       map[string]interface{}{},
			expectedErr: true,
			expected:    "is required",
		},
		{
			name:        "No error when valid input",
			lang:        "en",
			input:       map[string]interface{}{"name": "John"},
			expectedErr: false,
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var dest struct {
				Name string
			}
			n, ok := tc.input["name"]
			if !ok {
				n = ""
			}
			dest2 := struct {
				Name string
				Age  int
			}{
				Age:  1,
				Name: n.(string),
			}
			errs := schema.Parse(tc.input, &dest, zog.WithCtxValue(LangKey, tc.lang))
			errs2 := schema.Validate(&dest2, zog.WithCtxValue(LangKey, tc.lang))
			m1 := zog.Issues.Flatten(errs)
			m2 := zog.Issues.Flatten(errs2)

			if tc.expectedErr {
				assert.NotNil(t, errs, "Expected errors, got nil")
				assert.NotNil(t, errs2, "Expected errors, got nil")

				nameErrs, ok := m1["name"]
				assert.True(t, ok, "Expected errors for 'name' field")
				assert.NotEmpty(t, nameErrs, "Expected at least one error for 'name' field")
				assert.Equal(t, tc.expected, nameErrs[0], "Unexpected error message")

				nameErrs2, ok2 := m2["name"]
				assert.True(t, ok2, "Expected errors for 'name' field")
				assert.NotEmpty(t, nameErrs2, "Expected at least one error for 'name' field")
				assert.Equal(t, tc.expected, nameErrs2[0], "Unexpected error message")
			} else {
				assert.Nil(t, errs, "Expected no errors, got: %v", errs)
			}
		})
	}
}

func TestLangErrsMapWithLangKey(t *testing.T) {

	// Set up language maps
	SetLanguagesErrsMap(map[string]zconst.LangMap{
		"en": enMap,
		"es": esMap,
	}, "en", WithLangKey("customLangKey"))

	dest := ""
	destSchema := zog.String().Required()

	errs := destSchema.Parse(nil, &dest, zog.WithCtxValue("customLangKey", "es"))

	assert.NotNil(t, errs, "Expected errors, got nil")
	assert.Equal(t, "es requerido", errs[0].Message, "Unexpected error message")
}
