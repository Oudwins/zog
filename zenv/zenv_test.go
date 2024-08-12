package zenv

import (
	"os"
	"testing"

	z "github.com/Oudwins/zog"
	"github.com/stretchr/testify/assert"
)

func TestEnvParsing(t *testing.T) {
	type TestStruct struct {
		Str string `zog:"TEST_STR"`
	}
	env := TestStruct{}
	schema := z.Struct(z.Schema{
		"str": z.String().Required(),
	})

	os.Setenv("TEST_STR", "hello")
	err := schema.Parse(NewDataProvider(), &env)
	assert.Nil(t, err)
	assert.Equal(t, env.Str, "hello")
	os.Setenv("TEST_STR", "")
}
