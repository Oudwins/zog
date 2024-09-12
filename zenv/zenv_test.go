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
		"pORT": z.Int().Default(8080),
	})

	os.Setenv("TEST_STR", "hello")
	err := schema.Parse(NewDataProvider(), &env)
	assert.Nil(t, err)
	assert.Equal(t, "hello", env.Str)
	assert.Equal(t, 8080, env.PORT)
	os.Setenv("TEST_STR", "")
}
