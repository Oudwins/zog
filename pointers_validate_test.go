package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"

	p "github.com/Oudwins/zog/internals"
)

// TODO

func TestPtrPrimitive(t *testing.T) {
	// in := 10
	var out *int
	s := Ptr(Int().Required())
	err := s.Parse("", &out)
	assert.Empty(t, err)
	assert.Nil(t, out)

	err = s.Parse("not_empty", &out)
	assert.NotNil(t, err)
	assert.Equal(t, 0, *out)

	err = s.Parse(10, &out)
	assert.Empty(t, err)
	assert.Equal(t, 10, *out)

	// with zero value it also works
	err = s.Parse(0, &out)
	assert.Empty(t, err)
	assert.Equal(t, 0, *out)
}

func TestPtrInStruct(t *testing.T) {
	type TestStruct struct {
		Value *int
	}

	s := Struct(Schema{
		"value": Ptr(Int()),
	})
	in := map[string]any{
		"value": 10,
	}
	var out TestStruct
	err := s.Parse(in, &out)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Value)
	assert.Equal(t, 10, *out.Value)
}

func TestPtrPtrInStruct(t *testing.T) {
	type TestStruct struct {
		Value **int
	}

	s := Struct(Schema{
		"value": Ptr(Ptr(Int())),
	})
	in := map[string]any{
		"value": 10,
	}
	var out TestStruct
	// empty input
	err := s.Parse("", &out)
	assert.Empty(t, err)
	assert.Nil(t, out.Value)

	// good input
	err = s.Parse(in, &out)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Value)
	assert.NotNil(t, *out.Value)
	assert.Equal(t, 10, **out.Value)
}

func TestPtrNestedStructs(t *testing.T) {
	type Inner struct {
		Value *int
	}
	type Outer struct {
		Inner *Inner
	}

	schema := Struct(Schema{
		"inner": Ptr(Struct(Schema{
			"value": Ptr(Int()),
		})),
	})

	var out Outer
	data := map[string]any{
		"inner": map[string]any{
			"value": 10,
		},
	}

	err := schema.Parse(data, &out)
	assert.Nil(t, err)
	assert.NotNil(t, out.Inner)
	assert.NotNil(t, out.Inner.Value)
	assert.Equal(t, 10, *out.Inner.Value)
}

func TestPtrInSlice(t *testing.T) {
	schema := Slice(Ptr(Int()))
	var out []*int

	data := []any{10, 20, 30}
	err := schema.Parse(data, &out)

	assert.Nil(t, err)
	assert.Len(t, out, 3)
	assert.Equal(t, 10, *out[0])
	assert.Equal(t, 20, *out[1])
	assert.Equal(t, 30, *out[2])
}

func TestPtrSliceStruct(t *testing.T) {
	type TestStruct struct {
		Value int
	}

	schema := Slice(Ptr(Struct(Schema{
		"value": Int(),
	})))
	var out []*TestStruct

	data := []any{
		map[string]any{"value": 10},
		map[string]any{"value": 20},
		map[string]any{"value": 30},
	}
	err := schema.Parse(data, &out)

	assert.Nil(t, err)
	assert.Len(t, out, 3)
	assert.Equal(t, 10, out[0].Value)
	assert.Equal(t, 20, out[1].Value)
	assert.Equal(t, 30, out[2].Value)
}

func TestPtrRequired(t *testing.T) {
	schema := Ptr(String()).NotNil(Message("Testing"))
	var dest *string
	tests := []struct {
		Val         any
		ExpectedErr bool
	}{
		{nil, true},
		{"", true},
		{0, false},
		{false, false},
	}
	for _, test := range tests {
		err := schema.Parse(test.Val, &dest)
		if test.ExpectedErr {
			assert.NotNil(t, err)
			x := err[p.ERROR_KEY_ROOT]
			assert.Equal(t, "Testing", x[0].Message())
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestPtrToStruct(t *testing.T) {
	type TestStruct struct {
		Value *int
	}

	var dest *TestStruct
	s := Ptr(Struct(Schema{
		"value": Ptr(Int()),
	}))
	in := map[string]any{
		"value": 10,
	}
	err := s.Parse(in, &dest)
	assert.Nil(t, err)
	assert.NotNil(t, dest)
	assert.NotNil(t, dest.Value)
	assert.Equal(t, 10, *dest.Value)
}

func TestPtrToSlice(t *testing.T) {

	var dest *[]*int
	s := Ptr(Slice(Ptr(Int())))
	err := s.Parse([]any{10, 20, 30}, &dest)
	assert.Nil(t, err)
	assert.NotNil(t, dest)
	assert.Len(t, *dest, 3)
	assert.Equal(t, 10, *(*dest)[0])
	assert.Equal(t, 20, *(*dest)[1])
	assert.Equal(t, 30, *(*dest)[2])
}
