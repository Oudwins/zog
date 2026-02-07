package zog

import (
	"testing"
	"time"

	"github.com/Oudwins/zog/pkgs/internals/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

func TestMapBasic(t *testing.T) {
	m := map[string]int{}
	schema := EXPERIMENTAL_MAP[string, int](String(), Int())

	data := map[string]any{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, 1, m["one"])
	assert.Equal(t, 2, m["two"])
	assert.Equal(t, 3, m["three"])
}

func TestMapOfStructs(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	var userSchema = Struct(Shape{
		"name": String().Required(),
		"age":  Int().Required(),
	})

	var usersSchema = Struct(Shape{
		"users": EXPERIMENTAL_MAP[string, User](String(), userSchema),
	})

	var data = map[string]interface{}{
		"users": map[string]interface{}{
			"user1": map[string]interface{}{
				"name": "Jane",
				"age":  25,
			},
			"user2": map[string]interface{}{
				"name": "John",
				"age":  30,
			},
		},
	}

	type Team struct {
		Users map[string]User
	}
	var team Team

	errs := usersSchema.Parse(data, &team)
	assert.Empty(t, errs)
	assert.Len(t, team.Users, 2)
	assert.Equal(t, "Jane", team.Users["user1"].Name)
	assert.Equal(t, 25, team.Users["user1"].Age)
	assert.Equal(t, "John", team.Users["user2"].Name)
	assert.Equal(t, 30, team.Users["user2"].Age)
}

func TestMapOptional(t *testing.T) {
	m := map[string]int{}
	schema := EXPERIMENTAL_MAP[string, int](String(), Int()) // should be optional by default

	errs := schema.Parse(nil, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 0)

	schema.Required().Optional() // should override required
	errs = schema.Parse(nil, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 0)
}

func TestMapRequired(t *testing.T) {
	m := map[string]int{}
	customMsg := "This map is required and cannot be empty"
	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).Required(Message(customMsg))

	// Test with nil value
	errs := schema.Parse(nil, &m)
	assert.NotNil(t, errs)
	assert.Equal(t, customMsg, errs[0].Message)

	// Test with empty map
	errs = schema.Parse(map[string]int{}, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 0)
}

func TestMapDefault(t *testing.T) {
	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).Default(map[string]int{"a": 1, "b": 2})
	m := map[string]int{}
	err := schema.Parse(nil, &m)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, 2, m["b"])
}

func TestMapDefaultFuncParse(t *testing.T) {
	defaultCalls := 0
	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).DefaultFunc(func() map[string]int {
		defaultCalls++
		return map[string]int{"a": 1, "b": 2}
	})
	m := map[string]int{}
	err := schema.Parse(nil, &m)
	assert.Nil(t, err)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, m)
	assert.Equal(t, 1, defaultCalls)
}

func TestMapErrors(t *testing.T) {
	m := map[string]int{}
	schema := EXPERIMENTAL_MAP[string, int](String(), Int().GTE(5))

	data := map[string]any{
		"key1": 10,
		"key2": 3, // should fail
		"key3": 7,
	}

	errs := schema.Parse(data, &m)
	assert.NotEmpty(t, errs)
	errsM := Issues.Flatten(errs)
	assert.Len(t, errsM[`["key2"]`], 1)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestMapTransform(t *testing.T) {
	m := map[string]int{}
	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).Transform(func(dataPtr any, ctx Ctx) error {
		mp := dataPtr.(*map[string]int)
		for k, v := range *mp {
			(*mp)[k] = v * 2
		}
		return nil
	})

	data := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, 2, m["a"])
	assert.Equal(t, 4, m["b"])
	assert.Equal(t, 6, m["c"])
}

func TestMapLen(t *testing.T) {
	m := map[string]int{}

	data := map[string]any{
		"a": 1,
		"b": 2,
	}

	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).Len(2)
	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 2)

	data = map[string]any{
		"a": 1,
	}
	errs = schema.Parse(data, &m)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestMapMin(t *testing.T) {
	m := map[string]int{}

	data := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).Min(2)
	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)

	data = map[string]any{
		"a": 1,
	}
	errs = schema.Parse(data, &m)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestMapMax(t *testing.T) {
	m := map[string]int{}

	data := map[string]any{
		"a": 1,
		"b": 2,
	}

	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).Max(3)
	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)

	data = map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	errs = schema.Parse(data, &m)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestMapCustomTest(t *testing.T) {
	input := map[string]int{"a": 1, "b": 2}
	m := map[string]int{}
	schema := EXPERIMENTAL_MAP[string, int](String(), Int()).TestFunc(func(val any, ctx Ctx) bool {
		mp := val.(*map[string]int)
		return len(*mp) == 2
	}, Message("custom"))
	errs := schema.Parse(input, &m)
	assert.Nil(t, errs)
	assert.Equal(t, input, m)

	errs = schema.Parse(map[string]int{"a": 1}, &m)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "custom", errs[0].Message)
}

func TestMapGetType(t *testing.T) {
	s := EXPERIMENTAL_MAP[string, int](String(), Int())
	assert.Equal(t, zconst.TypeMap, s.getType())
}

func TestMapIntKeys(t *testing.T) {
	m := map[int]string{}
	schema := EXPERIMENTAL_MAP[int, string](Int(), String())

	data := map[string]any{
		"1": "one",
		"2": "two",
		"3": "three",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, "one", m[1])
	assert.Equal(t, "two", m[2])
	assert.Equal(t, "three", m[3])
}

func TestMapNestedMaps(t *testing.T) {
	m := map[string]map[string]int{}
	schema := EXPERIMENTAL_MAP[string, map[string]int](String(), EXPERIMENTAL_MAP[string, int](String(), Int()))

	data := map[string]any{
		"outer1": map[string]any{
			"inner1": 1,
			"inner2": 2,
		},
		"outer2": map[string]any{
			"inner3": 3,
		},
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 2)
	assert.Len(t, m["outer1"], 2)
	assert.Equal(t, 1, m["outer1"]["inner1"])
	assert.Equal(t, 2, m["outer1"]["inner2"])
	assert.Equal(t, 3, m["outer2"]["inner3"])
}

// Test different key types
func TestMapFloat64Keys(t *testing.T) {
	m := map[float64]string{}
	schema := EXPERIMENTAL_MAP[float64, string](Float64(), String())

	data := map[string]any{
		"1.5":  "one point five",
		"2.7":  "two point seven",
		"3.14": "pi",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, "one point five", m[1.5])
	assert.Equal(t, "two point seven", m[2.7])
	assert.Equal(t, "pi", m[3.14])
}

func TestMapBoolKeys(t *testing.T) {
	m := map[bool]string{}
	schema := EXPERIMENTAL_MAP[bool, string](Bool(), String())

	data := map[string]any{
		"true":  "yes",
		"false": "no",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 2)
	assert.Equal(t, "yes", m[true])
	assert.Equal(t, "no", m[false])
}

func TestMapInt64Keys(t *testing.T) {
	m := map[int64]string{}
	schema := EXPERIMENTAL_MAP[int64, string](Int64(), String())

	data := map[string]any{
		"100": "hundred",
		"200": "two hundred",
		"300": "three hundred",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, "hundred", m[100])
	assert.Equal(t, "two hundred", m[200])
	assert.Equal(t, "three hundred", m[300])
}

func TestMapUintKeys(t *testing.T) {
	m := map[uint]string{}
	schema := EXPERIMENTAL_MAP[uint, string](Uint(), String())

	data := map[string]any{
		"10": "ten",
		"20": "twenty",
		"30": "thirty",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, "ten", m[10])
	assert.Equal(t, "twenty", m[20])
	assert.Equal(t, "thirty", m[30])
}

// Test different value types
func TestMapBoolValues(t *testing.T) {
	m := map[string]bool{}
	schema := EXPERIMENTAL_MAP[string, bool](String(), Bool())

	data := map[string]any{
		"flag1": true,
		"flag2": false,
		"flag3": "true",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, true, m["flag1"])
	assert.Equal(t, false, m["flag2"])
	assert.Equal(t, true, m["flag3"])
}

func TestMapFloatValues(t *testing.T) {
	m := map[string]float64{}
	schema := EXPERIMENTAL_MAP[string, float64](String(), Float64())

	data := map[string]any{
		"pi":    3.14,
		"e":     2.718,
		"sqrt2": "1.414",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.InDelta(t, 3.14, m["pi"], 0.001)
	assert.InDelta(t, 2.718, m["e"], 0.001)
	assert.InDelta(t, 1.414, m["sqrt2"], 0.001)
}

func TestMapTimeValues(t *testing.T) {
	m := map[string]time.Time{}
	schema := EXPERIMENTAL_MAP[string, time.Time](String(), Time())

	data := map[string]any{
		"date1": "2024-01-01T00:00:00Z",
		"date2": "2024-12-31T23:59:59Z",
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 2)
	expected1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expected2 := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	assert.Equal(t, expected1, m["date1"])
	assert.Equal(t, expected2, m["date2"])
}

func TestMapInt64Values(t *testing.T) {
	m := map[string]int64{}
	schema := EXPERIMENTAL_MAP[string, int64](String(), Int64())

	data := map[string]any{
		"large1": int64(9223372036854775807),
		"large2": "9223372036854775806",
		"large3": 100,
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, int64(9223372036854775807), m["large1"])
	assert.Equal(t, int64(9223372036854775806), m["large2"])
	assert.Equal(t, int64(100), m["large3"])
}

func TestMapUintValues(t *testing.T) {
	m := map[string]uint{}
	schema := EXPERIMENTAL_MAP[string, uint](String(), Uint())

	data := map[string]any{
		"val1": uint(100),
		"val2": "200",
		"val3": 300,
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.Equal(t, uint(100), m["val1"])
	assert.Equal(t, uint(200), m["val2"])
	assert.Equal(t, uint(300), m["val3"])
}

// Test maps with complex value types
func TestMapPtrValues(t *testing.T) {
	m := map[string]*int{}
	schema := EXPERIMENTAL_MAP[string, *int](String(), Ptr(Int()))

	data := map[string]any{
		"val1": 10,
		"val2": 20,
		"val3": nil,
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 3)
	assert.NotNil(t, m["val1"])
	assert.Equal(t, 10, *m["val1"])
	assert.NotNil(t, m["val2"])
	assert.Equal(t, 20, *m["val2"])
	assert.Nil(t, m["val3"])
}

func TestMapSliceValues(t *testing.T) {
	m := map[string][]string{}
	schema := EXPERIMENTAL_MAP[string, []string](String(), Slice(String()))

	data := map[string]any{
		"fruits": []any{"apple", "banana", "cherry"},
		"colors": []any{"red", "blue"},
		"empty":  []any{},
		"single": []any{"one"},
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 4)
	assert.Equal(t, []string{"apple", "banana", "cherry"}, m["fruits"])
	assert.Equal(t, []string{"red", "blue"}, m["colors"])
	assert.Equal(t, []string{}, m["empty"])
	assert.Equal(t, []string{"one"}, m["single"])
}

func TestMapNestedComplexValues(t *testing.T) {
	// Map with slice of structs as values
	type Person struct {
		Name string
		Age  int
	}

	personSchema := Struct(Shape{
		"name": String().Required(),
		"age":  Int().Required(),
	})

	m := map[string][]Person{}
	schema := EXPERIMENTAL_MAP[string, []Person](String(), Slice(personSchema))

	data := map[string]any{
		"team1": []any{
			map[string]any{"name": "Alice", "age": 30},
			map[string]any{"name": "Bob", "age": 25},
		},
		"team2": []any{
			map[string]any{"name": "Charlie", "age": 35},
		},
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 2)
	assert.Len(t, m["team1"], 2)
	assert.Equal(t, "Alice", m["team1"][0].Name)
	assert.Equal(t, 30, m["team1"][0].Age)
	assert.Equal(t, "Bob", m["team1"][1].Name)
	assert.Equal(t, 25, m["team1"][1].Age)
	assert.Len(t, m["team2"], 1)
	assert.Equal(t, "Charlie", m["team2"][0].Name)
	assert.Equal(t, 35, m["team2"][0].Age)
}

// Test maps inside Slice
func TestSliceOfMaps(t *testing.T) {
	s := []map[string]int{}
	schema := Slice(EXPERIMENTAL_MAP[string, int](String(), Int()))

	data := []any{
		map[string]any{"a": 1, "b": 2},
		map[string]any{"c": 3, "d": 4},
		map[string]any{"e": 5},
	}

	errs := schema.Parse(data, &s)
	assert.Nil(t, errs)
	assert.Len(t, s, 3)
	assert.Len(t, s[0], 2)
	assert.Equal(t, 1, s[0]["a"])
	assert.Equal(t, 2, s[0]["b"])
	assert.Len(t, s[1], 2)
	assert.Equal(t, 3, s[1]["c"])
	assert.Equal(t, 4, s[1]["d"])
	assert.Len(t, s[2], 1)
	assert.Equal(t, 5, s[2]["e"])
}

func TestSliceOfMapsWithValidation(t *testing.T) {
	s := []map[string]int{}
	schema := Slice(EXPERIMENTAL_MAP[string, int](String(), Int().GTE(5)))

	data := []any{
		map[string]any{"a": 10, "b": 20},
		map[string]any{"c": 3}, // should fail validation
		map[string]any{"d": 7},
	}

	errs := schema.Parse(data, &s)
	assert.NotEmpty(t, errs)
	errsM := Issues.Flatten(errs)
	// Error should be at [1]["c"]
	assert.True(t, len(errsM[`[1]["c"]`]) > 0 || len(errsM[`[1].c`]) > 0)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestSliceOfMapsEmpty(t *testing.T) {
	s := []map[string]int{}
	schema := Slice(EXPERIMENTAL_MAP[string, int](String(), Int()))

	errs := schema.Parse([]any{}, &s)
	assert.Nil(t, errs)
	assert.Len(t, s, 0)
}

// Test maps inside Ptr
func TestPtrToMap(t *testing.T) {
	var m *map[string]int
	schema := Ptr(EXPERIMENTAL_MAP[string, int](String(), Int()))

	data := map[string]any{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.NotNil(t, m)
	assert.Len(t, *m, 3)
	assert.Equal(t, 1, (*m)["one"])
	assert.Equal(t, 2, (*m)["two"])
	assert.Equal(t, 3, (*m)["three"])
}

func TestPtrToMapNil(t *testing.T) {
	var m *map[string]int
	schema := Ptr(EXPERIMENTAL_MAP[string, int](String(), Int()))

	errs := schema.Parse(nil, &m)
	assert.Nil(t, errs)
	assert.Nil(t, m)
}

func TestPtrToMapRequired(t *testing.T) {
	var m *map[string]int
	schema := Ptr(EXPERIMENTAL_MAP[string, int](String(), Int())).NotNil(Message("map is required"))

	errs := schema.Parse(nil, &m)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "map is required", errs[0].Message)
	assert.Nil(t, m)
}

func TestPtrToMapWithValidation(t *testing.T) {
	var m *map[string]int
	schema := Ptr(EXPERIMENTAL_MAP[string, int](String(), Int().GTE(5)))

	data := map[string]any{
		"valid":   10,
		"invalid": 3, // should fail
	}

	errs := schema.Parse(data, &m)
	assert.NotEmpty(t, errs)
	errsM := Issues.Flatten(errs)
	assert.NotEmpty(t, errsM[`["invalid"]`])
	tutils.VerifyDefaultIssueMessages(t, errs)
}

// Test maps inside Boxed schemas
type MapBox struct {
	Value map[string]int
}

func TestBoxedMap(t *testing.T) {
	s := Boxed(
		EXPERIMENTAL_MAP[string, int](String(), Int()),
		func(b MapBox, ctx Ctx) (map[string]int, error) { return b.Value, nil },
		func(m map[string]int, ctx Ctx) (MapBox, error) { return MapBox{Value: m}, nil },
	)

	var box MapBox
	input := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	errs := s.Parse(input, &box)
	assert.Empty(t, errs)
	assert.Len(t, box.Value, 3)
	assert.Equal(t, 1, box.Value["a"])
	assert.Equal(t, 2, box.Value["b"])
	assert.Equal(t, 3, box.Value["c"])
}

func TestBoxedMapWithValidation(t *testing.T) {
	s := Boxed(
		EXPERIMENTAL_MAP[string, int](String(), Int().GTE(5)),
		func(b MapBox, ctx Ctx) (map[string]int, error) { return b.Value, nil },
		func(m map[string]int, ctx Ctx) (MapBox, error) { return MapBox{Value: m}, nil },
	)

	var box MapBox
	input := map[string]any{
		"valid":   10,
		"invalid": 3, // should fail
	}
	errs := s.Parse(input, &box)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestBoxedMapFromBox(t *testing.T) {
	s := Boxed(
		EXPERIMENTAL_MAP[string, int](String(), Int()),
		func(b MapBox, ctx Ctx) (map[string]int, error) { return b.Value, nil },
		func(m map[string]int, ctx Ctx) (MapBox, error) { return MapBox{Value: m}, nil },
	)

	var box MapBox
	inputBox := MapBox{Value: map[string]int{"x": 1, "y": 2}}
	errs := s.Parse(inputBox, &box)
	assert.Empty(t, errs)
	assert.Len(t, box.Value, 2)
	assert.Equal(t, 1, box.Value["x"])
	assert.Equal(t, 2, box.Value["y"])
}

// Test maps with Any schema as value type
func TestMapWithAnyValueParse(t *testing.T) {
	m := map[string]any{}
	schema := EXPERIMENTAL_MAP[string, any](String(), EXPERIMENTAL_ANY())

	data := map[string]any{
		"number": 42,
		"string": "hello",
		"bool":   true,
		"nested": map[string]any{"inner": "value"},
	}

	errs := schema.Parse(data, &m)
	assert.Nil(t, errs)
	assert.Len(t, m, 4)
	assert.Equal(t, 42, m["number"])
	assert.Equal(t, "hello", m["string"])
	assert.Equal(t, true, m["bool"])
	assert.Equal(t, map[string]any{"inner": "value"}, m["nested"])
}

func TestMapWithAnyValueValidate(t *testing.T) {
	m := map[string]any{
		"number": 42,
		"string": "hello",
		"bool":   true,
		"nested": map[string]any{"inner": "value"},
	}
	schema := EXPERIMENTAL_MAP[string, any](String(), EXPERIMENTAL_ANY())

	errs := schema.Validate(&m)
	assert.Nil(t, errs)
	assert.Len(t, m, 4)
	assert.Equal(t, 42, m["number"])
	assert.Equal(t, "hello", m["string"])
	assert.Equal(t, true, m["bool"])
	assert.Equal(t, map[string]any{"inner": "value"}, m["nested"])
}

func TestMapWithAnyValueRequired(t *testing.T) {
	m := map[string]any{}
	schema := EXPERIMENTAL_MAP[string, any](String(), EXPERIMENTAL_ANY().Required())

	data := map[string]any{
		"present": "value",
		"nil":     nil,
	}

	errs := schema.Parse(data, &m)
	assert.NotNil(t, errs)
	errsM := Issues.Flatten(errs)
	assert.NotEmpty(t, errsM[`["nil"]`])
}

func TestMapWithAnyValueValidateRequired(t *testing.T) {
	m := map[string]any{
		"present": "value",
		"nil":     nil,
	}
	schema := EXPERIMENTAL_MAP[string, any](String(), EXPERIMENTAL_ANY().Required())

	errs := schema.Validate(&m)
	assert.NotNil(t, errs)
	errsM := Issues.Flatten(errs)
	assert.NotEmpty(t, errsM[`["nil"]`])
}
