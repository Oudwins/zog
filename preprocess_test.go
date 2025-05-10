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

func TestPreprocessFloat(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out float64, err error) {
		return strconv.ParseFloat(data, 64)
	}, Float64().GT(0))

	var out float64
	errs := s.Parse("123.45", &out)
	assert.Equal(t, 123.45, out)
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
	}, Slice[string](String().Min(1)))

	out := []string{}
	errs := s.Parse("hello,world", &out)
	assert.Nil(t, errs)
	assert.Len(t, out, 2)
	assert.Equal(t, "hello", out[0])
	assert.Equal(t, "world", out[1])
}

func TestPreprocessStruct(t *testing.T) {
	type User struct {
		Id   string
		Name string
	}
	s := Preprocess(func(data string, ctx Ctx) (out User, err error) {
		parts := strings.Split(data, ",")
		return User{Id: parts[0], Name: parts[1]}, nil
	}, Struct(
		Shape{
			"Id":   String().Min(1),
			"Name": String().Min(1),
		},
	))
	var out User
	errs := s.Parse("1,John Doe", &out)
	assert.Nil(t, errs)
	assert.Equal(t, "1", out.Id)
	assert.Equal(t, "John Doe", out.Name)
}

func TestPreprocessSliceOfStructs(t *testing.T) {
	type User struct {
		Id   string
		Name string
	}
	s := Preprocess(func(data string, ctx Ctx) (out []User, err error) {
		rows := strings.Split(data, ";")
		result := make([]User, len(rows))
		for i, row := range rows {
			parts := strings.Split(row, ",")
			result[i] = User{Id: parts[0], Name: parts[1]}
		}
		return result, nil
	}, Slice[*Shape](Struct(Shape{
		"Id":   String().Min(1),
		"Name": String().Min(1),
	})))

	var out []User
	errs := s.Parse("1,John Doe;2,Jane Doe", &out)
	assert.Nil(t, errs)
	assert.Len(t, out, 2)
	assert.Equal(t, "1", out[0].Id)
	assert.Equal(t, "John Doe", out[0].Name)
	assert.Equal(t, "2", out[1].Id)
	assert.Equal(t, "Jane Doe", out[1].Name)
}

func TestPreprocessStructWithSlice(t *testing.T) {
	type User struct {
		Id    string
		Names []string
	}
	s := Preprocess(func(data string, ctx Ctx) (out User, err error) {
		parts := strings.Split(data, ":")
		return User{
			Id:    parts[0],
			Names: strings.Split(parts[1], ","),
		}, nil
	}, Struct(Shape{
		"Id":    String().Min(1),
		"Names": Slice[string](String().Min(1)),
	}))

	var out User
	errs := s.Parse("1:John,Jane,Joe", &out)
	assert.Nil(t, errs)
	assert.Equal(t, "1", out.Id)
	assert.Equal(t, []string{"John", "Jane", "Joe"}, out.Names)
}

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

	var in interface{} = "x"
	var out string
	errs := s.Parse(in, &out)
	assert.Nil(t, errs)
	assert.Equal(t, "x", out)
}

func TestPreprocessPtrStruct(t *testing.T) {
	type User struct {
		Id   string
		Name string
	}
	s := Preprocess(func(data string, ctx Ctx) (out *User, err error) {
		parts := strings.Split(data, ",")
		return &User{Id: parts[0], Name: parts[1]}, nil
	}, Ptr(Struct(
		Shape{
			"Id":   String().Min(1),
			"Name": String().Min(1),
		},
	)))
	var out *User
	errs := s.Parse("1,John Doe", &out)
	assert.Nil(t, errs)
	assert.Equal(t, "1", out.Id)
	assert.Equal(t, "John Doe", out.Name)
}

func TestPreprocessPtrString(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out *string, err error) {
		str := strings.ToUpper(data)
		return &str, nil
	}, Ptr(String().Min(1)))

	var out *string
	errs := s.Parse("hello", &out)
	assert.Nil(t, errs)
	assert.Equal(t, "HELLO", *out)
}

func TestPreprocessPtrInt(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out *int, err error) {
		n, _ := strconv.Atoi(data)
		return &n, nil
	}, Ptr(Int().GT(0)))

	var out *int
	errs := s.Parse("42", &out)
	assert.Nil(t, errs)
	assert.Equal(t, 42, *out)
}

func TestPreprocessPtrFloat(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out *float64, err error) {
		f, _ := strconv.ParseFloat(data, 64)
		return &f, nil
	}, Ptr(Float64().GT(0)))

	var out *float64
	errs := s.Parse("3.14", &out)
	assert.Nil(t, errs)
	assert.Equal(t, 3.14, *out)
}

func TestPreprocessPtrBool(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out *bool, err error) {
		b := data == "true"
		return &b, nil
	}, Ptr(Bool()))

	var out *bool
	errs := s.Parse("true", &out)
	assert.Nil(t, errs)
	assert.Equal(t, true, *out)
}

func TestPreprocessPtrTime(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out *time.Time, err error) {
		t, _ := time.Parse(time.RFC3339, data)
		return &t, nil
	}, Ptr(Time()))

	var out *time.Time
	errs := s.Parse("2023-01-01T00:00:00Z", &out)
	assert.Nil(t, errs)
	assert.Equal(t, "2023-01-01 00:00:00 +0000 UTC", out.String())
}

func TestPreprocessPtrSlice(t *testing.T) {
	s := Preprocess(func(data string, ctx Ctx) (out *[]string, err error) {
		slice := strings.Split(data, ",")
		return &slice, nil
	}, Ptr(Slice[string](String().Min(1))))

	var out *[]string
	errs := s.Parse("a,b,c", &out)
	assert.Nil(t, errs)
	assert.Equal(t, []string{"a", "b", "c"}, *out)
}

func TestPreprocessPartOfStruct(t *testing.T) {
	type User struct {
		Id   string
		Name string
		Age  int
	}
	s := Struct(Shape{
		"Id": Preprocess(
			func(data string, ctx Ctx) (out string, err error) {
				return strings.ToUpper(data), nil
			},
			String().Min(1),
		),
		"Name": String().Min(1),
		"Age": Preprocess(func(data string, ctx Ctx) (out int, err error) {
			return strconv.Atoi(data)
		}, Int().GT(0)),
	})

	var out User
	errs := s.Parse(map[string]string{"Id": "one", "Name": "John Doe", "Age": "20"}, &out)
	assert.Nil(t, errs)
	assert.Equal(t, "ONE", out.Id)
	assert.Equal(t, "John Doe", out.Name)
	assert.Equal(t, 20, out.Age)
}

func TestPreprocessInSlice(t *testing.T) {

	s := Slice[string](
		Preprocess(func(data string, ctx Ctx) (out int, err error) {
			return strconv.Atoi(data)
		}, Int().GT(0)),
	)

	var out []int
	errs := s.Parse([]string{"1", "2", "3"}, &out)
	assert.Nil(t, errs)
	assert.Equal(t, []int{1, 2, 3}, out)
}
