package zog

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessString(t *testing.T) {
	s := Preprocess(func(data int, ctx Ctx) (out string, err error) {
		return strconv.Itoa(data), nil
	}, String().Min(1))

	var out string
	errs := s.Parse(1, &out)
	assert.Equal(t, "1", out)
	assert.Equal(t, 0, len(errs))
}

func TestPreprocessBool(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out bool, err error) {
		return data == "true", nil
	}, Bool())

	var out bool
	errs := s.Parse("true", &out)
	assert.Equal(t, true, out)
	assert.Equal(t, 0, len(errs))
}

func TestPreprocessInt(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out int, err error) {
		return strconv.Atoi(data)
	}, Int().GT(0))

	var out int
	errs := s.Parse("123", &out)
	assert.Equal(t, 123, out)
	assert.Equal(t, 0, len(errs))
}

func TestPreprocessTime(t *testing.T) {
	s := Preprocess(func(data int64, ctx Ctx) (out time.Time, err error) {
		return time.Unix(data, 0), nil
	}, Time())

	var out time.Time
	timestamp := int64(1640995200) // 2022-01-01 00:00:00
	errs := s.Parse(timestamp, &out)
	assert.Equal(t, time.Unix(timestamp, 0), out)
	assert.Equal(t, 0, len(errs))
}

func TestPreprocessSlice(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out []string, err error) {
		return strings.Split(data, ","), nil
	}, Slice(String().Min(1)))

	out := []string{}
	errs := s.Parse("hello,world", &out)
	assert.Nil(t, errs)
	assert.Len(t, out, 2)
	assert.Equal(t, "hello", out[0])
	assert.Equal(t, "world", out[1])
}

// func TestPreprocessStruct(t *testing.T) {
// 	type User struct {
// 		Id   string
// 		Name string
// 	}
// 	s := Preprocess(func(data string, ctx Ctx) (out User, err error) {
// 		parts := strings.Split(data, ",")
// 		return User{Id: parts[0], Name: parts[1]}, nil
// 	}, Struct(
// 		Schema{
// 			"Id":   String().Min(1),
// 			"Name": String().Min(1),
// 		},
// 	))
// 	var out User
// 	errs := s.Parse("1,John Doe", &out)
// 	assert.Nil(t, errs)
// 	assert.Equal(t, "1", out.Id)
// 	assert.Equal(t, "John Doe", out.Name)
// }

func TestPreprocessWithAny(t *testing.T) {
	s := Preprocess(func(data any, ctx Ctx) (out string, err error) {
		switch v := data.(type) {
		case string:
			return v, nil
		case int:
			return "int", nil
		default:
			return "default", nil
		}
	}, String().Min(1))

	var in any
	in = "x"
	var out string
	errs := s.Parse(in, &out)
	assert.Nil(t, errs)
	assert.Equal(t, "x", out)
}

// func TestSliceDefaultCoercing(t *testing.T) {
// 	s := []string{}
// 	schema := Slice(String())
// 	errs := schema.Parse("a", &s)
// 	assert.Nil(t, errs)
// 	assert.Len(t, s, 1)
// 	assert.Equal(t, s[0], "a")
// }

// func TestSliceSchemaOption(t *testing.T) {
// 	s := Slice(String(), WithCoercer(func(original any) (value any, err error) {
// 		return []string{"coerced"}, nil
// 	}))

// 	var result []string
// 	err := s.Parse(123, &result)
// 	assert.Nil(t, err)
// 	assert.Equal(t, []string{"coerced"}, result)
// }

// func TestTimePreTransform(t *testing.T) {
// 	var now time.Time
// 	schema := Time().PreTransform(func(data any, ctx Ctx) (any, error) {
// 		// Add 1 hour to the input time
// 		t, ok := data.(time.Time)
// 		if !ok {
// 			return nil, nil
// 		}
// 		return t.Add(time.Hour), nil
// 	})

// 	input := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
// 	expected := input.Add(time.Hour)

// 	errs := schema.Parse(input, &now)
// 	assert.Nil(t, errs)
// 	assert.Equal(t, expected, now)
// }

// func TestTimeSchemaOption(t *testing.T) {
// 	s := Time(WithCoercer(func(original any) (value any, err error) {
// 		return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil
// 	}))

// 	var result time.Time
// 	err := s.Parse("invalid-date", &result)
// 	assert.Nil(t, err)
// 	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), result)
// }

// func TestTimeFormat(t *testing.T) {
// 	s := Time(Time.Format(time.RFC1123))
// 	var result time.Time
// 	err := s.Parse("Mon, 01 Jan 2024 00:00:00 UTC", &result)
// 	assert.Nil(t, err)
// 	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), result)
// }

// func TestTimeFormatFunc(t *testing.T) {
// 	s := Time(Time.FormatFunc(func(data string) (time.Time, error) {
// 		return time.Parse(time.RFC1123, data)
// 	}))
// 	var result time.Time
// 	err := s.Parse("Mon, 01 Jan 2024 00:00:00 UTC", &result)
// 	assert.Nil(t, err)
// 	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), result)
// }

// func TestBoolPreTransform(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		data      interface{}
// 		transform p.PreTransform
// 		expectErr bool
// 		expected  bool
// 	}{
// 		{
// 			name: "Valid transform",
// 			data: "true",
// 			transform: func(val any, ctx Ctx) (any, error) {
// 				if s, ok := val.(*string); ok {
// 					return *s == "true", nil
// 				}
// 				return val, nil
// 			},
// 			expected: true,
// 		},
// 		{
// 			name: "Invalid transform",
// 			data: "invalid",
// 			transform: func(val any, ctx Ctx) (any, error) {
// 				return nil, fmt.Errorf("invalid input")
// 			},
// 			expectErr: true,
// 			expected:  false,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			boolProc := Bool().PreTransform(test.transform)
// 			var result bool
// 			errs := boolProc.Parse(test.data, &result)

// 			if (len(errs) > 0) != test.expectErr {
// 				t.Errorf("Expected error: %v, got: %v", test.expectErr, errs)
// 			}

// 			if len(errs) > 0 {
// 				tutils.VerifyDefaultIssueMessages(t, errs)
// 			}

// 			if result != test.expected {
// 				t.Errorf("Expected %v, but got %v", test.expected, result)
// 			}
// 		})
// 	}
// }
