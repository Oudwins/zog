package zog

import "testing"

func TestToJson(t *testing.T) {
	s := String().Required().Default("Testing!").Catch("Testing2!").Min(1)
	ToJson(s)
}
