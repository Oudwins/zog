package zog

import (
	"testing"
	"time"

	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

type validateObj struct {
	Str string
	In  int
	Fl  float64
	Bol bool
	Tim time.Time
}

var validateObjSchema = Struct(Schema{
	"str": String().Required(),
	"in":  Int().Required(),
	"fl":  Float().Required(),
	"bol": Bool().Required(),
	"tim": Time().Required(),
})

func TestValidateStructExample(t *testing.T) {
	o := validateObj{
		Str: "hello",
		In:  10,
		Fl:  10.5,
		Bol: true,
		Tim: time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC),
	}

	errs := validateObjSchema.Validate(&o)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", o.Str)
}

var validateNestedSchema = Struct(Schema{
	"str":    String().Required(),
	"schema": Struct(Schema{"str": String().Required()}),
})

func TestValidateStructNestedStructs(t *testing.T) {
	v := struct {
		Str    string
		Schema struct {
			Str string
		}
	}{
		Str: "hello",
		Schema: struct {
			Str string
		}{
			Str: "hello",
		},
	}

	errs := validateNestedSchema.Validate(&v)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", v.Str)
	assert.Equal(t, "hello", v.Schema.Str)

	v2 := struct {
		Str    string
		Schema struct {
			Str string
		}
	}{
		Schema: struct {
			Str string
		}{},
	}

	errs = validateNestedSchema.Validate(&v2)
	assert.NotNil(t, errs)
}

func TestValidateStructOptional(t *testing.T) {
	type TestStruct struct {
		Str string `zog:"str"`
		In  int    `zog:"in"`
		Fl  float64
		Bol bool
		Tim time.Time
	}

	var validateOptionalSchema = Ptr(Struct(Schema{
		"str": String().Required(),
		"in":  Int().Required(),
		"fl":  Float().Required(),
		"bol": Bool().Required(),
		"tim": Time().Required(),
	}))

	var o *TestStruct
	errs := validateOptionalSchema.Validate(&o)
	assert.Empty(t, errs)
}

func TestValidateStructCustomTestInSchema(t *testing.T) {
	type CustomStruct struct {
		Str string `zog:"str"`
		Num int    `zog:"num"`
	}

	// Create a custom test function
	customTest := func(val any, ctx Ctx) bool {
		// Custom test logic here
		num := val.(int)
		return num > 0
	}

	// Create a schema with a custom test
	schema := Struct(Schema{
		"str": String().Required(),
		"num": Int().Test(TestFunc("customTest", customTest)),
	})

	obj := CustomStruct{
		Str: "hello",
		Num: 10,
	}

	errs := schema.Validate(&obj)
	assert.Empty(t, errs)
	assert.Equal(t, "hello", obj.Str)
	assert.Equal(t, 10, obj.Num)

	obj = CustomStruct{
		Str: "hello",
		Num: -10,
	}

	errs = schema.Validate(&obj)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "customTest", errs["num"][0].Code)
}

func TestValidateStructCustomTest(t *testing.T) {
	type CustomStruct struct {
		Str string `zog:"str"`
	}

	schema := Struct(Schema{
		"str": String(),
	}).Test(TestFunc("customTest", func(val any, ctx Ctx) bool {
		s := val.(*CustomStruct)
		return s.Str == "valid"
	}), Message("customTest"))

	obj := CustomStruct{
		Str: "invalid",
	}

	errs := schema.Validate(&obj)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "customTest", errs["$root"][0].Code)
	assert.Equal(t, "customTest", errs["$root"][0].Message)

	obj.Str = "valid"
	errs = schema.Validate(&obj)
	assert.Empty(t, errs)
}

func TestValidateStructPreTransforms(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	preTransform := func(val any, ctx Ctx) (any, error) {
		if s, ok := val.(TestStruct); ok {
			s.Value = "transformed"
			return s, nil
		}
		return val, nil
	}

	schema := Struct(Schema{
		"value": String().Required(),
	}).PreTransform(preTransform)

	output := TestStruct{Value: "original"}

	errs := schema.Validate(&output)
	assert.Empty(t, errs)
	assert.Equal(t, "transformed", output.Value)
}

func TestValidateStructPostTransforms(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	postTransform := func(val any, ctx Ctx) error {
		if s, ok := val.(*TestStruct); ok {
			s.Value = "post_" + s.Value
		}
		return nil
	}

	schema := Struct(Schema{
		"value": String().Required(),
	}).PostTransform(postTransform)

	output := TestStruct{Value: "original"}

	errs := schema.Validate(&output)
	assert.Empty(t, errs)
	assert.Equal(t, "post_original", output.Value)
}

func TestValidateStructPassThroughRequired(t *testing.T) {
	type TestStruct struct {
		Somefield string
	}

	schema := Struct(Schema{
		"somefield": String().Required(),
	})

	output := TestStruct{
		Somefield: "someValue",
	}

	errs := schema.Validate(&output)
	assert.Empty(t, errs)
	assert.Equal(t, "someValue", output.Somefield)

	var output2 TestStruct
	errs = schema.Validate(&output2)
	assert.NotEmpty(t, errs)
	assert.Equal(t, zconst.IssueCodeRequired, errs["somefield"][0].Code)
}

func TestValidateStructGetType(t *testing.T) {
	s := Struct(Schema{
		"field": String(),
	})
	assert.Equal(t, zconst.TypeStruct, s.getType())
}
