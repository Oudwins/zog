package zog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// !STRUCTS
// NOT YET SUPPORTED. Should be simple. Just need to check if its a struct processor in which case we need to build an any data provider
//
//	func TestSliceStruct(t *testing.T) {
//		type TestStruct struct {
//			T string
//		}
//		schema := Slice(Struct(Schema{"t": String().Required()}))
//		sl := []TestStruct{}
//		errs := schema.Parse([]any{map[string]any{"t": "a"}, map[string]any{"t": "b"}}, &sl)
//	}

func TestSlicePassSchema(t *testing.T) {

	s := []string{}
	schema := Slice(String().Required())

	errs := schema.Parse([]any{"a", "b", "c"}, &s)
	assert.Nil(t, errs)
	fmt.Println(s)
	assert.Len(t, s, 3)
	assert.Equal(t, s[0], "a")
	assert.Equal(t, s[1], "b")
	assert.Equal(t, s[2], "c")
}

func TestSliceErrors(t *testing.T) {
	s := []string{}
	schema := Slice(String().Required().Min(2))

	errs := schema.Parse([]any{"a", "b"}, &s)
	assert.Len(t, errs, 3)
	assert.NotEmpty(t, errs["0"])
	assert.NotEmpty(t, errs["1"])
	assert.Empty(t, errs["2"])
}

func TestSliceLen(t *testing.T) {
	s := []string{}

	els := []string{"a", "b", "c", "d", "e"}
	schema := Slice(String().Required()).Len(2)
	errs := schema.Parse(els[:2], &s)
	assert.Len(t, s, 2)
	assert.Nil(t, errs)
	errs = schema.Parse(els[:1], &s)
	assert.NotEmpty(t, errs)

	// min
	schema = Slice(String().Required()).Min(2)
	errs = schema.Parse(els[:4], &s)
	assert.Nil(t, errs)
	errs = schema.Parse(els[:1], &s)
	assert.NotEmpty(t, errs)
	// max
	schema = Slice(String().Required()).Max(3)
	errs = schema.Parse(els[:1], &s)
	assert.Nil(t, errs)
	errs = schema.Parse(els[:4], &s)
	assert.NotNil(t, errs)

}

func TestSliceContains(t *testing.T) {

	s := []string{}
	items := []string{"a", "b", "c"}

	schema := Slice(String()).Contains("a")
	errs := schema.Parse(items, &s)
	assert.Nil(t, errs)
	assert.Len(t, s, 3)

	schema = Slice(String()).Contains("d")
	errs = schema.Parse(items, &s)
	assert.NotEmpty(t, errs)
}

func TestSliceDefaultCoercing(t *testing.T) {
	s := []string{}
	schema := Slice(String())
	errs := schema.Parse("a", &s)
	assert.Nil(t, errs)
	assert.Len(t, s, 1)
	assert.Equal(t, s[0], "a")
}

func TestSliceDefault(t *testing.T) {
	schema := Slice(String()).Default([]string{"a", "b", "c"})
	s := []string{}
	err := schema.Parse(nil, &s)
	assert.Nil(t, err)
	assert.Len(t, s, 3)
	assert.Equal(t, s[0], "a")
	assert.Equal(t, s[1], "b")
	assert.Equal(t, s[2], "c")
}
