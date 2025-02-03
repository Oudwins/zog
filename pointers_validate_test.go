package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePtrPrimitive(t *testing.T) {
	var out *int = new(int)
	*out = 10
	s := Ptr(Int())
	errs := s.Validate(&out)
	assert.Empty(t, errs)

	out = nil
	errs = s.Validate(&out)
	assert.Empty(t, errs)

	assert.Panics(t, func() {
		s.Validate(nil)
	})
}

func TestValidatePtrInStruct(t *testing.T) {
	type TestStruct struct {
		Value *int
	}

	val := 10
	out := TestStruct{Value: &val}
	s := Struct(Schema{
		"value": Ptr(Int()),
	})
	errs := s.Validate(&out)
	assert.Empty(t, errs)
	assert.Equal(t, 10, *out.Value)
}

func TestValidatePtrPtrInStruct(t *testing.T) {
	type TestStruct struct {
		Value **int
	}

	val := 10
	pval := &val
	out := TestStruct{Value: &pval}
	s := Struct(Schema{
		"value": Ptr(Ptr(Int())),
	})

	errs := s.Validate(&out)
	assert.Empty(t, errs)
	assert.Equal(t, 10, **out.Value)

	out.Value = nil
	errs = s.Validate(&out)
	assert.Empty(t, errs)
}

func TestValidatePtrNestedStructs(t *testing.T) {
	type Inner struct {
		Value *int
	}
	type Outer struct {
		Inner *Inner
	}

	val := 10
	inner := Inner{Value: &val}
	out := Outer{Inner: &inner}

	schema := Struct(Schema{
		"inner": Ptr(Struct(Schema{
			"value": Ptr(Int()),
		})),
	})

	errs := schema.Validate(&out)
	assert.Empty(t, errs)
	assert.Equal(t, 10, *out.Inner.Value)
}

func TestValidatePtrInSlice(t *testing.T) {
	schema := Slice(Ptr(Int()).NotNil(Message("Testing")))
	v1, v2, v3 := 10, 20, 30
	var v4 *int
	out := []*int{&v1, &v2, &v3, v4}

	errs := schema.Validate(&out)
	assert.NotEmpty(t, errs)
	assert.Equal(t, 10, *out[0])
	assert.Equal(t, 20, *out[1])
	assert.Equal(t, 30, *out[2])
	assert.Equal(t, "Testing", errs["[3]"][0].Message())
}

func TestValidatePtrSliceStruct(t *testing.T) {
	type TestStruct struct {
		Value int
	}

	schema := Slice(Ptr(Struct(Schema{
		"value": Int(),
	})))
	out := []*TestStruct{
		{Value: 10},
		{Value: 20},
		{Value: 30},
	}

	errs := schema.Validate(&out)
	assert.Empty(t, errs)
	assert.Equal(t, 10, out[0].Value)
	assert.Equal(t, 20, out[1].Value)
	assert.Equal(t, 30, out[2].Value)
}

func TestValidatePtrRequired(t *testing.T) {
	schema := Ptr(String()).NotNil(Message("Testing"))
	var dest *string
	errs := schema.Validate(&dest)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Testing", errs[0].Message())

	str := "test"
	dest = &str
	errs = schema.Validate(&dest)
	assert.Empty(t, errs)
}

func TestValidatePtrToStruct(t *testing.T) {
	type TestStruct struct {
		Value *int
	}

	val := 10
	dest := &TestStruct{Value: &val}
	s := Ptr(Struct(Schema{
		"value": Ptr(Int()),
	}))

	errs := s.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, 10, *dest.Value)
}

func TestValidatePtrToSlice(t *testing.T) {
	v1, v2, v3 := 10, 20, 30
	dest := &[]*int{&v1, &v2, &v3}
	s := Ptr(Slice(Ptr(Int())))

	errs := s.Validate(&dest)
	assert.Empty(t, errs)
	assert.Equal(t, 10, *(*dest)[0])
	assert.Equal(t, 20, *(*dest)[1])
	assert.Equal(t, 30, *(*dest)[2])
}
