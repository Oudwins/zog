package zog

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessStringValidate(t *testing.T) {
	s := Preprocess(func(data *string, ctx Ctx) (out string, err error) {
		return *data + "1", nil
	}, String().Min(1))

	str := "1"
	errs := s.Validate(&str)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "11", str)
}

func TestPreprocessBoolValidate(t *testing.T) {
	s := Preprocess(func(data *bool, ctx Ctx) (out bool, err error) {
		return !*data, nil
	}, Bool())

	b := true
	errs := s.Validate(&b)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, false, b)
}

func TestPreprocessIntValidate(t *testing.T) {
	s := Preprocess(func(data *int, ctx Ctx) (out int, err error) {
		return *data + 1, nil
	}, Int().GT(0))

	i := 123
	errs := s.Validate(&i)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 124, i)
}

func TestPreprocessFloatValidate(t *testing.T) {
	s := Preprocess(func(data *float64, ctx Ctx) (out float64, err error) {
		return *data + 1.0, nil
	}, Float64().GT(0))

	f := 123.45
	errs := s.Validate(&f)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 124.45, f)
}

func TestPreprocessTimeValidate(t *testing.T) {
	s := Preprocess(func(data *time.Time, ctx Ctx) (out time.Time, err error) {
		return data.AddDate(1, 0, 0), nil
	}, Time())

	timestamp := time.Unix(1640995200, 0) // 2022-01-01 00:00:00
	errs := s.Validate(&timestamp)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, time.Unix(1640995200, 0).AddDate(1, 0, 0), timestamp)
}

func TestPreprocessSliceValidate(t *testing.T) {
	s := Preprocess(func(data *[]string, ctx Ctx) (out []string, err error) {
		return append(*data, "!"), nil
	}, Slice(String().Min(1)))

	slice := []string{"hello", "world"}
	errs := s.Validate(&slice)
	assert.Nil(t, errs)
	assert.Equal(t, []string{"hello", "world", "!"}, slice)
}

func TestPreprocessStructValidate(t *testing.T) {
	type User struct {
		Id   string
		Name string
	}
	s := Preprocess(func(data *User, ctx Ctx) (out User, err error) {
		return User{Id: data.Id + "!", Name: data.Name + "!"}, nil
	}, Struct(
		Shape{
			"Id":   String().Min(1),
			"Name": String().Min(1),
		},
	))

	user := User{Id: "1", Name: "John Doe"}
	errs := s.Validate(&user)
	assert.Nil(t, errs)
	assert.Equal(t, User{Id: "1!", Name: "John Doe!"}, user)
}

func TestPreprocessSliceOfStructsValidate(t *testing.T) {
	type User struct {
		Id   string
		Name string
	}
	s := Preprocess(func(data *[]User, ctx Ctx) (out []User, err error) {
		result := make([]User, len(*data))
		for i, u := range *data {
			result[i] = User{Id: u.Id + "!", Name: u.Name + "!"}
		}
		return result, nil
	}, Slice(Struct(Shape{
		"Id":   String().Min(1),
		"Name": String().Min(1),
	})))

	users := []User{
		{Id: "1", Name: "John Doe"},
		{Id: "2", Name: "Jane Doe"},
	}
	errs := s.Validate(&users)
	assert.Nil(t, errs)
	assert.Equal(t, []User{
		{Id: "1!", Name: "John Doe!"},
		{Id: "2!", Name: "Jane Doe!"},
	}, users)
}

func TestPreprocessStructWithSliceValidate(t *testing.T) {
	type User struct {
		Id    string
		Names []string
	}
	s := Preprocess(func(data *User, ctx Ctx) (out User, err error) {
		names := make([]string, len(data.Names))
		for i, n := range data.Names {
			names[i] = n + "!"
		}
		return User{
			Id:    data.Id + "!",
			Names: names,
		}, nil
	}, Struct(Shape{
		"Id":    String().Min(1),
		"Names": Slice(String().Min(1)),
	}))

	user := User{
		Id:    "1",
		Names: []string{"John", "Jane", "Joe"},
	}
	errs := s.Validate(&user)
	assert.Nil(t, errs)
	assert.Equal(t, User{
		Id:    "1!",
		Names: []string{"John!", "Jane!", "Joe!"},
	}, user)
}

func TestPreprocessWithAnyValidate(t *testing.T) {
	s := Preprocess(func(data any, ctx Ctx) (out string, err error) {
		switch v := (data).(type) {
		case *string:
			return *v + "!", nil
		case int:
			return "int!", nil
		default:
			return "default!", nil
		}
	}, String().Min(1))

	var str string = "x"
	errs := s.Validate(&str)
	assert.Nil(t, errs)
	assert.Equal(t, "x!", str)
}

func TestPreprocessPtrStructValidate(t *testing.T) {
	type User struct {
		Id   string
		Name string
	}
	s := Preprocess(func(data **User, ctx Ctx) (out *User, err error) {
		return &User{Id: (*data).Id + "!", Name: (*data).Name + "!"}, nil
	}, Ptr(Struct(
		Shape{
			"Id":   String().Min(1),
			"Name": String().Min(1),
		},
	)))

	user := User{Id: "1", Name: "John Doe"}
	ptr := &user
	errs := s.Validate(&ptr)
	_ = errs
	// assert.Nil(t, errs)
	// assert.Equal(t, User{Id: "1!", Name: "John Doe!"}, user)
}

func TestPreprocessPtrStringValidate(t *testing.T) {
	s := Preprocess(func(data **string, ctx Ctx) (out *string, err error) {
		x := "!"
		return &x, nil
	}, Ptr(String().Min(1)))

	var pstr *string
	errs := s.Validate(&pstr)
	assert.Nil(t, errs)
	assert.Equal(t, "!", *pstr)
}

func TestPreprocessPtrIntValidate(t *testing.T) {
	s := Preprocess(func(data **int, ctx Ctx) (out *int, err error) {
		x := **data + 1
		return &x, nil
	}, Ptr(Int().GT(0)))

	x := 42
	px := &x
	errs := s.Validate(&px)
	assert.Nil(t, errs)
	assert.Equal(t, 43, *px)
	assert.Equal(t, 42, x)
}

func TestPreprocessPtrFloatValidate(t *testing.T) {
	s := Preprocess(func(data **float64, ctx Ctx) (out *float64, err error) {
		**data += 0.1
		return *data, nil
	}, Ptr(Float64().GT(0)))

	f := 3.14
	pf := &f
	errs := s.Validate(&pf)
	assert.Nil(t, errs)
	assert.Equal(t, 3.24, *pf)
	assert.Equal(t, 3.24, f)
}

func TestPreprocessPtrBoolValidate(t *testing.T) {
	s := Preprocess(func(data **bool, ctx Ctx) (out *bool, err error) {
		**data = !**data
		return *data, nil
	}, Ptr(Bool()))

	b := true
	pb := &b
	errs := s.Validate(&pb)
	assert.Nil(t, errs)
	assert.Equal(t, false, b)
	assert.Equal(t, false, *pb)
}

func TestPreprocessPtrTimeValidate(t *testing.T) {
	s := Preprocess(func(data **time.Time, ctx Ctx) (out *time.Time, err error) {
		x := **data
		x = x.AddDate(1, 0, 0)
		return &x, nil
	}, Ptr(Time()))

	tim := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	pt := &tim
	errs := s.Validate(&pt)
	assert.Nil(t, errs)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), *pt)
}

func TestPreprocessPtrSliceValidate(t *testing.T) {
	s := Preprocess(func(data **[]string, ctx Ctx) (out *[]string, err error) {
		s := **data
		for i, str := range s {
			s[i] = str + "!"
		}
		return &s, nil
	}, Ptr(Slice(String().Min(1))))

	slice := []string{"a", "b", "c"}
	pslice := &slice
	errs := s.Validate(&pslice)
	assert.Nil(t, errs)
	assert.Equal(t, []string{"a!", "b!", "c!"}, *pslice)
	assert.Equal(t, []string{"a!", "b!", "c!"}, slice)
}

func TestPreprocessPartOfStructValidate(t *testing.T) {
	type User struct {
		Id   string
		Name string
		Age  int
	}
	s := Struct(Shape{
		"Id": Preprocess(
			func(data *string, ctx Ctx) (out string, err error) {
				return strings.ToUpper(*data), nil
			},
			String().Min(1),
		),
		"Name": String().Min(1),
		"Age": Preprocess(func(data *int, ctx Ctx) (out int, err error) {
			return *data + 1, nil
		}, Int().GT(0)),
	})

	out := User{Id: "one", Name: "John Doe", Age: 20}
	errs := s.Validate(&out)
	assert.Nil(t, errs)
	assert.Equal(t, "ONE", out.Id)
	assert.Equal(t, "John Doe", out.Name)
	assert.Equal(t, 21, out.Age)
}
