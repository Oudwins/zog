package zenv

import (
	"os"
	"testing"

	z "github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestEnvParsing(t *testing.T) {
	type TestStruct struct {
		Str  string `zog:"TEST_STR"`
		PORT int
	}
	env := TestStruct{}
	schema := z.Struct(z.Schema{
		"str":  z.String().Required(),
		"PORT": z.Int().Default(8080),
	})

	os.Setenv("TEST_STR", "hello")
	err := schema.Parse(NewDataProvider(), &env)
	assert.Nil(t, err)
	assert.Equal(t, "hello", env.Str)
	assert.Equal(t, 8080, env.PORT)
	os.Setenv("TEST_STR", "")
}

func TestEnvParsingWithEnvTag(t *testing.T) {
	type TestStruct struct {
		Str  string `env:"TEST_STR"`
		PORT int
	}
	env := TestStruct{}
	schema := z.Struct(z.Schema{
		"str":  z.String().Required(),
		"PORT": z.Int().Default(8080),
	})

	os.Setenv("TEST_STR", "hello")
	err := schema.Parse(NewDataProvider(), &env)
	assert.Nil(t, err)
	assert.Equal(t, "hello", env.Str)
	assert.Equal(t, 8080, env.PORT)
	os.Setenv("TEST_STR", "")
}

func TestEnvParsingWithConflictingTags(t *testing.T) {
	type TestStruct struct {
		Str  string `zog:"TEST_STR2" env:"TEST_STR"` // env takes priority for being more specific
		PORT int
	}
	env := TestStruct{}
	schema := z.Struct(z.Schema{
		"str":  z.String().Required(),
		"PORT": z.Int().Default(8080),
	})

	os.Setenv("TEST_STR", "hello")
	err := schema.Parse(NewDataProvider(), &env)
	assert.Nil(t, err)
	assert.Equal(t, "hello", env.Str)
	assert.Equal(t, 8080, env.PORT)
	os.Setenv("TEST_STR", "")
}

// Unit tests for envDataProvider
func TestNewDataProvider(t *testing.T) {
	provider := NewDataProvider()
	assert.NotNil(t, provider)
	assert.IsType(t, &envDataProvider{}, provider)
}

func TestEnvDataProviderGet(t *testing.T) {
	provider := NewDataProvider()

	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	result := provider.Get("TEST_VAR")
	assert.Equal(t, "test_value", result)

	// Test with non-existent environment variable
	result = provider.Get("NON_EXISTENT_VAR")
	assert.Equal(t, nil, result)

	// Test with environment variable containing whitespace
	os.Setenv("WHITESPACE_VAR", "  trimmed  ")
	result = provider.Get("WHITESPACE_VAR")
	assert.Equal(t, "trimmed", result)

	// Clean up
	os.Unsetenv("TEST_VAR")
	os.Unsetenv("WHITESPACE_VAR")
}

func TestEnvDataProviderGetNestedProvider(t *testing.T) {
	provider := NewDataProvider()
	nestedProvider := provider.GetNestedProvider("any_key")
	assert.Equal(t, provider, nestedProvider)
}

func TestEnvDataProviderGetUnderlying(t *testing.T) {
	provider := NewDataProvider()
	underlying := provider.GetUnderlying()
	assert.Nil(t, underlying)
}
