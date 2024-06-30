package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlicePassSchema(t *testing.T) {
	type TestStruct struct {
		Items []any
	}

	s := TestStruct{
		Items: []any{"a", "b", "c"},
	}

	errs, ok := Validate(s, Schema{"items": Slice(String().Len(1))})
	assert.True(t, ok)
	assert.Empty(t, errs)

	s.Items = []any{"a", "b", "c", "d", 1}
	errs, ok = Validate(s, Schema{"items": Slice(String().Len(1))})
	assert.False(t, ok)
	assert.Len(t, errs, 1)
}

func TestSliceNotEmpty(t *testing.T) {
	type TestStruct struct {
		Items []any
	}
	s := TestStruct{
		Items: []any{},
	}

	errs, ok := Validate(s, Schema{"items": Slice(String()).Min(1)})
	assert.False(t, ok)
	assert.NotEmpty(t, errs)

	s.Items = []any{"a", "b", "c"}
	errs, ok = Validate(s, Schema{"items": Slice(String()).Min(1)})
	assert.True(t, ok)
	assert.Empty(t, errs)
}

func TestSliceEnum(t *testing.T) {
	type TestStruct struct {
		Items []string
	}

	s := TestStruct{
		Items: []string{"a", "b", "c"},
	}

	schema := Schema{"items": Slice(Enum([]string{"a", "b", "c"})).Optional()}

	errs, ok := Parse(s, schema)
	assert.True(t, ok)
	assert.Empty(t, errs)

	s.Items = []string{"testing", "badval", "two"}
	errs, ok = Parse(s, schema)
	assert.False(t, ok)
	assert.Len(t, errs, 1)
}

func TestSliceLength(t *testing.T) {
	type TestStruct struct {
		Items []any
	}

	s := TestStruct{
		Items: []any{"a", "b", "c"},
	}

	errs, ok := Validate(s, Schema{"items": Slice(String()).Len(3)})
	assert.True(t, ok)
	assert.Empty(t, errs)

	errs, ok = Validate(s, Schema{"items": Slice(String()).Len(2)})
	assert.False(t, ok)
	assert.Len(t, errs, 1)

	// min & max
	errs, ok = Validate(s, Schema{"items": Slice(String()).Min(2)})
	assert.True(t, ok)
	assert.Empty(t, errs)

	errs, ok = Validate(s, Schema{"items": Slice(String()).Min(4)})
	assert.False(t, ok)
	assert.Len(t, errs, 1)

	errs, ok = Validate(s, Schema{"items": Slice(String()).Max(3)})
	assert.True(t, ok)
	assert.Empty(t, errs)

	errs, ok = Validate(s, Schema{"items": Slice(String()).Max(1)})
	assert.False(t, ok)
	assert.Len(t, errs, 1)
}

func TestSliceContains(t *testing.T) {

	type TestStruct struct {
		Items []any
	}

	s := TestStruct{
		Items: []any{"a", "b", "c"},
	}

	errs, ok := Validate(s, Schema{"items": Slice(String()).Contains("a")})
	assert.True(t, ok)
	assert.Empty(t, errs)

	errs, ok = Validate(s, Schema{"items": Slice(String()).Contains("d")})
	assert.False(t, ok)
	assert.Len(t, errs, 1)
}

func TestSliceOptional(t *testing.T) {
	type TestStruct struct {
		Items []any
	}

	s := TestStruct{}

	errs, ok := Parse(s, Schema{"items": Slice(String()).Optional()})
	assert.True(t, ok)
	assert.Empty(t, errs)

	s.Items = []any{"a", "b", "c"}
	errs, ok = Parse(s, Schema{"items": Slice(String()).Optional()})
	assert.True(t, ok)
	assert.Empty(t, errs)
}
