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
	i := 0
	errs := s.Parse(1, &i)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom error", errs[0].Message)

	// Test successful case
	errs = s.Parse(11, &i)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 11, i)
}

func TestCustomFuncString(t *testing.T) {
	s := CustomFunc(func(ptr *string, ctx Ctx) bool {
		return len(*ptr) > 5
	}, Message("string too short"))
	str := ""
	errs := s.Parse("test", &str)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "string too short", errs[0].Message)

	errs = s.Parse("testing", &str)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, "testing", str)
}

func TestCustomFuncSlice(t *testing.T) {
	s := CustomFunc(func(ptr *[]int, ctx Ctx) bool {
		return len(*ptr) > 3
	}, Message("custom error"))
	i := []int{}
	errs := s.Parse([]int{1, 2, 3}, &i)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "custom error", errs[0].Message)
	errs = s.Parse([]int{1, 2, 3, 4}, &i)
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
		return ptr.Age >= 18 && len(ptr.Name) >= 3
	}, Message("invalid user data"))

	u := User{}
	errs := s.Parse(User{Age: 17, Name: "Jo"}, &u)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, "invalid user data", errs[0].Message)

	errs = s.Parse(User{Age: 18, Name: "John"}, &u)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 18, u.Age)
	assert.Equal(t, "John", u.Name)
}
