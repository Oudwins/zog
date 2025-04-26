package zog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomFunc(t *testing.T) {
	s := CustomFunc(func(ptr *int, ctx Ctx) bool {
		return *ptr > 10
	}, Message("custom error"))
	// parse
	i := 0
	errs := s.Parse(1, &i)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom error", errs[0].Message)

	// Test successful case
	errs = s.Parse(11, &i)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 11, i)

	// Validate failure
	i = 1
	errs = s.Validate(&i)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom error", errs[0].Message)

	// Validate success
	i = 11
	errs = s.Validate(&i)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 11, i)

}

func TestCustomFuncString(t *testing.T) {
	s := CustomFunc(func(ptr *string, ctx Ctx) bool {
		return len(*ptr) > 5
	}, Message("string too short"))

	// Parse failure
	str := ""
	errs := s.Parse("test", &str)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "string too short", errs[0].Message)

	// Parse success
	errs = s.Parse("testing", &str)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "testing", str)

	// Validate failure
	str = "test"
	errs = s.Validate(&str)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "string too short", errs[0].Message)

	// Validate success
	str = "testing"
	errs = s.Validate(&str)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "testing", str)
}

func TestCustomFuncSlice(t *testing.T) {
	s := CustomFunc(func(ptr *[]int, ctx Ctx) bool {
		return len(*ptr) > 3
	}, Message("custom error"))

	// Parse failure
	i := []int{}
	errs := s.Parse([]int{1, 2, 3}, &i)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom error", errs[0].Message)

	// Parse success
	errs = s.Parse([]int{1, 2, 3, 4}, &i)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, []int{1, 2, 3, 4}, i)

	// Validate failure
	i = []int{1, 2, 3}
	errs = s.Validate(&i)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom error", errs[0].Message)

	// Validate success
	i = []int{1, 2, 3, 4}
	errs = s.Validate(&i)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, []int{1, 2, 3, 4}, i)
}

func TestCustomFuncStruct(t *testing.T) {
	type User struct {
		Age  int       `zog:"min=18"`
		Name string    `zog:"min=3"`
		DOB  time.Time `zog:"format=2006-01-02"`
	}
	s := CustomFunc(func(ptr *User, ctx Ctx) bool {
		out := ptr.Age >= 18 && len(ptr.Name) >= 3
		return out
	}, Message("invalid user data"))

	// Parse failure
	u := User{}
	errs := s.Parse(User{Age: 17, Name: "Jo"}, &u)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "invalid user data", errs[0].Message)

	// Parse success
	errs = s.Parse(User{Age: 18, Name: "John"}, &u)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 18, u.Age)
	assert.Equal(t, "John", u.Name)

	// Validate failure
	u = User{Age: 17, Name: "Jo"}
	errs = s.Validate(&u)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "invalid user data", errs[0].Message)

	// Validate success
	u = User{Age: 18, Name: "John"}
	errs = s.Validate(&u)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 18, u.Age)
	assert.Equal(t, "John", u.Name)
}

func TestCustomFuncInStruct(t *testing.T) {
	type ID struct {
		ID string `zog:"custom"`
	}
	type User struct {
		ID ID
	}

	s := Struct(Shape{
		"ID": CustomFunc(func(ptr *ID, ctx Ctx) bool {
			return len(ptr.ID) > 0
		}, Message("invalid id")),
	})

	u := User{}
	errs := s.Parse(map[string]any{
		"ID": ID{ID: "123"},
	}, &u)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "123", u.ID.ID)
	errs = s.Validate(&u)
	assert.Equal(t, 0, len(errs))
}

func TestCustomFuncInStructPtr(t *testing.T) {
	type ID struct {
		ID string `zog:"custom"`
	}
	type User struct {
		ID *ID
	}

	s := Struct(Shape{
		"ID": Ptr(CustomFunc(func(ptr *ID, ctx Ctx) bool {
			return len(ptr.ID) > 0
		}, Message("invalid id"))),
	})

	u := User{}
	errs := s.Parse(map[string]any{
		"ID": ID{ID: "123"},
	}, &u)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "123", u.ID.ID)
	errs = s.Validate(&u)
	assert.Equal(t, 0, len(errs))
}
