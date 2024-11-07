package conf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBoolCoercer(t *testing.T) {
	var b any
	var err error
	b, err = Coercers.Bool(true)
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool("true")
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool("on")
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool(1)
	assert.True(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool("off")
	assert.False(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool(0)
	assert.False(t, b.(bool))
	assert.Nil(t, err)
	b, err = Coercers.Bool(false)
	assert.False(t, b.(bool))
	assert.Nil(t, err)
}

func TestStringCoercer(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{input: "hello", want: "hello"},
		{input: 123, want: "123"},
		{input: true, want: "true"},
		{input: 1.23, want: "1.23"},
		{input: []any{"hello"}, want: "[hello]"},
		{input: []int{1, 2, 3}, want: "[1 2 3]"},
		{input: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), want: "2022-01-01 00:00:00 +0000 UTC"},
	}
	var s any
	var err error
	for _, test := range tests {
		s, err = Coercers.String(test.input)
		assert.Equal(t, test.want, s)
		assert.Nil(t, err)
	}
}

func TestIntCoercer(t *testing.T) {
	var i any
	var err error
	tests := []struct {
		input any
		want  int
		err   bool
	}{
		{input: 123, want: 123},
		{input: "123", want: 123},
		{input: 1.23, want: 1},
		{input: true, want: 1},
		{input: "x", err: true},
	}
	for _, test := range tests {
		i, err = Coercers.Int(test.input)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.want, i.(int))
		}
	}
}

func TestFloat64Coercer(t *testing.T) {
	var f any
	var err error
	tests := []struct {
		input any
		want  float64
		err   bool
	}{
		{input: 123, want: 123.00},
		{input: "123", want: 123.00},
		{input: 1.23, want: 1.23},
		{input: "x", err: true},
	}

	for _, test := range tests {
		f, err = Coercers.Float64(test.input)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.want, f.(float64))
		}
	}
}

func TestTimeCoercer(t *testing.T) {
	var out any
	var err error
	now := time.Now()
	tests := []struct {
		input any
		want  time.Time
		err   bool
	}{
		{input: now, want: now},
		{input: "2024-09-09T00:00:00.000Z", want: time.Date(2024, 9, 9, 0, 0, 0, 0, time.UTC)},
		{input: 1.23, err: true},
		{input: 1733007600, want: time.Unix(1733007600, 0)},
	}

	for _, test := range tests {
		out, err = Coercers.Time(test.input)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.want, out.(time.Time))
		}
	}
}

func TestSliceCoercer(t *testing.T) {
	var out any
	var err error
	now := time.Now()
	tests := []struct {
		input any
		want  []any
		err   bool
	}{
		{input: now, want: []any{now}},
		{input: "x", want: []any{"x"}},
		{input: 1.23, want: []any{1.23}},
		{input: []any{"x"}, want: []any{"x"}},
	}
	for _, test := range tests {
		out, err = Coercers.Slice(test.input)
		if test.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.want, out)
		}
	}
}
