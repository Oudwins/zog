package zog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func normalize(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\n", ""), "\t", "")
}

func TestToJsonString(t *testing.T) {
	s := String().Required().Default("Testing!").Catch("Testing2!").Min(1)
	serialized, err := ToJson(s)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	expected := `{"Type":"string","Processors":[{"Type":"test","IssueCode":"min","IssuePath":"","Params":{"min":1}}],"Child":null,"Required":{"Type":"test","IssueCode":"required","IssuePath":"","Params":{}},"DefaultValue":{},"CatchValue":{}}`

	assert.Equal(t, normalize(expected), normalize(string(serialized)))
}

func TestToJsonPtr(t *testing.T) {
	s := Ptr(String().Required().Default("Testing!").Catch("Testing2!").Min(1))
	serialized, err := ToJson(s)
	fmt.Println(string(serialized))
	assert.Nil(t, err)
	assert.NotNil(t, serialized)
}
