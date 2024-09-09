package zog

import (
	"testing"

	p "github.com/Oudwins/zog/primitives"
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
	assert.NotEmpty(t, errs["[0]"])
	assert.NotEmpty(t, errs["[1]"])
	assert.Empty(t, errs["[2]"])
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

type User struct {
	Name string
}

type Team struct {
	Users []User
}

func TestSliceOfStructs(t *testing.T) {

	var userSchema = Struct(Schema{
		"name": String().Required(),
	})

	var teamSchema = Struct(Schema{
		"users": Slice(userSchema),
	})

	var data = map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"name": "Jane",
			},
			map[string]interface{}{
				"name": "John",
			},
		},
	}
	var team Team

	errsMap := teamSchema.Parse(NewMapDataProvider(data), &team)
	assert.Nil(t, errsMap)
	assert.Len(t, team.Users, 2)
	assert.Equal(t, team.Users[0].Name, "Jane")
	assert.Equal(t, team.Users[1].Name, "John")

	data = map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		},
	}
	errsMap = teamSchema.Parse(NewMapDataProvider(data), &team)

	assert.Len(t, errsMap["users[0].name"], 1)
	assert.Len(t, errsMap["users[1].name"], 1)
}

func TestSliceCustomTest(t *testing.T) {
	input := []string{"abc", "defg", "hijkl"}
	s := []string{}
	schema := Slice(String()).Test("custom_test", func(val any, ctx p.ParseCtx) bool {
		// Custom test logic here
		x := val.(*[]string)
		return assert.Equal(t, input, *x)
	})
	errs := schema.Parse(input, &s)
	assert.Empty(t, errs)
}

func TestSliceInvalidData(t *testing.T) {
	input := "not a slice"
	s := []string{}
	schema := Slice(String())
	errs := schema.Parse(input, &s)
	assert.Nil(t, errs)
	assert.Equal(t, s, []string{input})
}
