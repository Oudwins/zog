package zog

import (
	"regexp"
	"testing"

	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/tutils"
	"github.com/stretchr/testify/assert"
)

type Env string

const (
	A Env = "A"
	B Env = "B"
	C Env = "C"
)

func MyStringSchema() *StringSchema[Env] {
	s := &StringSchema[Env]{}
	WithCoercer(func(x any) (any, error) {
		v, e := conf.DefaultCoercers.String(x)
		if e != nil {
			return nil, e
		}
		return Env(v.(string)), e
	})(s)
	return s
}

func TestCustomStringBasics(t *testing.T) {
	// Test optional by default
	s := MyStringSchema()
	var data Env
	err := s.Parse("", &data)
	assert.Nil(t, err)
	assert.Equal(t, Env(""), data)

	// Test required
	s = MyStringSchema().Required()
	err = s.Parse("", &data)
	assert.NotNil(t, err)

	// Test default value
	s = MyStringSchema().Default(A)
	err = s.Parse("", &data)
	assert.Nil(t, err)
	assert.Equal(t, A, data)

	// Test catch
	s = MyStringSchema().Required().Min(5).Catch(A)
	err = s.Parse("x", &data)
	assert.Nil(t, err)
	assert.Equal(t, A, data)
}

func TestCustomStringTransforms(t *testing.T) {
	// Test pre-transform
	s := MyStringSchema()
	var data Env = "test"

	// Test post-transform
	s = MyStringSchema().PostTransform(func(val any, ctx Ctx) error {
		v := val.(*Env)
		*v = Env(string(*v) + "_post")
		return nil
	})
	data = "test"
	err := s.Validate(&data)
	assert.Nil(t, err)
	assert.Equal(t, Env("test_post"), data)
}

func TestCustomStringValidators(t *testing.T) {
	var data Env

	// Test length
	s := MyStringSchema().Len(3)
	err := s.Parse("foo", &data)
	assert.Nil(t, err)
	assert.Equal(t, Env("foo"), data)

	// Test min/max
	s = MyStringSchema().Min(2).Max(4)
	err = s.Parse("foo", &data)
	assert.Nil(t, err)
	err = s.Parse("toolong", &data)
	assert.NotNil(t, err)

	// Test contains
	s = MyStringSchema().Contains("test")
	err = s.Parse("testing", &data)
	assert.Nil(t, err)
	err = s.Parse("fail", &data)
	assert.NotNil(t, err)

	// Test prefix/suffix
	s = MyStringSchema().HasPrefix("pre").HasSuffix("fix")
	err = s.Parse("prefix", &data)
	assert.Nil(t, err)

	// Test regex
	s = MyStringSchema().Match(regexp.MustCompile("^[0-9]+$"))
	err = s.Parse("123", &data)
	assert.Nil(t, err)
	err = s.Validate(&data)
	assert.Nil(t, err)
	err = s.Parse("abc", &data)
	assert.NotNil(t, err)

	// Test OneOf
	s = MyStringSchema().OneOf([]Env{A, B, C})
	err = s.Parse("A", &data)
	assert.Nil(t, err)
	assert.Equal(t, A, data)
	err = s.Parse("D", &data)
	assert.NotNil(t, err)
}

func TestCustomStringSpecialValidators(t *testing.T) {
	var data Env

	// Test email
	s := MyStringSchema().Email()
	err := s.Parse("test@example.com", &data)
	assert.Nil(t, err)
	err = s.Parse("invalid", &data)
	assert.NotNil(t, err)

	// Test URL
	s = MyStringSchema().URL()
	err = s.Parse("http://example.com", &data)
	assert.Nil(t, err)
	err = s.Parse("invalid", &data)
	assert.NotNil(t, err)

	// Test UUID
	s = MyStringSchema().UUID()
	err = s.Parse("123e4567-e89b-12d3-a456-426614174000", &data)
	assert.Nil(t, err)
	err = s.Parse("invalid", &data)
	assert.NotNil(t, err)
}

func TestCustomStringNot(t *testing.T) {
	tests := map[string]struct {
		schema         *StringSchema[Env]
		strVal         string
		expectedErrMap internals.ZogIssueList
	}{
		"not len success": {
			schema:         MyStringSchema().Not().Len(10).Contains("test"),
			strVal:         "test",
			expectedErrMap: nil,
		},
		"not len fail": {
			schema: MyStringSchema().Not().Len(4).Contains("t"),
			strVal: "test",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_len",
					Params:  map[string]any{"len": 4},
					Dtype:   "string",
					Value:   tutils.PtrOf[Env]("test"),
					Message: "string must not be exactly 4 character(s)",
					Err:     nil,
				},
			},
		},
		"not email": {
			schema:         MyStringSchema().Not().Email(),
			strVal:         "not-an-email",
			expectedErrMap: nil,
		},
		"not email failure": {
			schema: MyStringSchema().Not().Email(),
			strVal: "test@test.com",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_email",
					Params:  nil,
					Dtype:   "string",
					Value:   tutils.PtrOf[Env]("test@test.com"),
					Message: "must not be a valid email",
					Err:     nil,
				},
			},
		},
		"not with empty": {
			schema: MyStringSchema().Not().Len(1),
			strVal: "a",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_len",
					Params:  map[string]any{"len": 1},
					Dtype:   "string",
					Value:   tutils.PtrOf[Env]("a"),
					Message: "string must not be exactly 1 character(s)",
					Err:     nil,
				},
			},
		},
		"not url": {
			schema:         MyStringSchema().Not().URL(),
			strVal:         "not a url",
			expectedErrMap: nil,
		},
		"not url failure": {
			schema: MyStringSchema().Not().URL(),
			strVal: "https://google.com",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_url",
					Dtype:   "string",
					Value:   tutils.PtrOf[Env]("https://google.com"),
					Message: "must not be a valid URL",
					Err:     nil,
				},
			},
		},
		"not has prefix": {
			schema:         MyStringSchema().Not().HasPrefix("test_"),
			strVal:         "value",
			expectedErrMap: nil,
		},
		"not has prefix failure": {
			schema: MyStringSchema().Not().HasPrefix("test_"),
			strVal: "test_value",
			expectedErrMap: internals.ZogIssueList{
				&internals.ZogIssue{
					Code:    "not_prefix",
					Params:  map[string]any{"prefix": Env("test_")},
					Dtype:   "string",
					Value:   tutils.PtrOf[Env]("test_value"),
					Message: "string must not start with test_",
					Err:     nil,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var data Env
			err := test.schema.Parse(test.strVal, &data)
			assert.Equal(t, test.expectedErrMap, err)
			assert.Equal(t, Env(test.strVal), data)
		})
	}
}
