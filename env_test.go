package zog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fmt.Printf("VALUE: %v | err: %v", val, err)
func TestCoerceToString(t *testing.T) {
	val, err := coerceString[string]("123")
	assert.Nil(t, err)
	assert.Equal(t, "123", val)
}

func TestCoerceToInt(t *testing.T) {
	val, err := coerceString[int]("123")
	assert.Nil(t, err)
	assert.Equal(t, 123, val)
}

func TestCoerceToFloat(t *testing.T) {
	val, err := coerceString[float64]("123.25")
	assert.Nil(t, err)
	assert.Equal(t, 123.25, val)
}

func TestCoerceToBool(t *testing.T) {
	val, err := coerceString[bool]("true")
	assert.Nil(t, err)
	assert.Equal(t, true, val)
}

func TestEmptyEnv(t *testing.T) {
	assert.Panics(t, func() {
		Env[string]("Test", String().Required())
	})

	assert.Panics(t, func() {
		os.Setenv("TEST", "")
		Env[string]("Test", String().Required())
	})
}

func TestReturnsValue(t *testing.T) {

	os.Setenv("TEST", "value")

	val := Env[string]("TEST", String().Required())

	assert.Equal(t, val, "value")
}

func TestDefault(t *testing.T) {
	val := Env[string]("TEST2", String().Required(), "hello")
	assert.Equal(t, "hello", val)

	os.Setenv("TEST2", "world")
	val = Env[string]("TEST2", String().Required(), "hello")
	assert.Equal(t, "world", val)

	assert.Panics(t, func() {
		os.Setenv("TEST2", "1")
		_ = Env[string]("TEST2", String().Min(4))
	})
}

func TestInt(t *testing.T) {
	os.Setenv("TEST", "1")
	val := Env[int]("TEST", Int().LT(2))
	assert.Equal(t, 1, val)
}

func TestBool(t *testing.T) {
	os.Setenv("TEST", "true")
	val := Env[bool]("TEST", Bool().True())
	assert.Equal(t, true, val)
}

func TestFloat(t *testing.T) {
	os.Setenv("TEST", "1.1")
	val := Env[float64]("TEST", Float().GT(1.0))
	assert.Equal(t, 1.1, val)
}
