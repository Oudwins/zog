package zog

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

// structs with pointers
//maps with additional values
// errors are correct
// panics are correct

type obj struct {
	Str string
	In  int
	Fl  float64
	Bol bool
	Tim time.Time
}

var objSchema = Struct(Shape{
	"str": String().Required(),
	"in":  Int().Required(),
	"fl":  Float().Required(),
	"bol": Bool().Required(),
	"tim": Time().Required(),
})

type objTagged struct {
	Str string    `zog:"s"`
	In  int       `zog:"i"`
	Fl  float64   `zog:"f"`
	Bol bool      `zog:"b"`
	Tim time.Time `zog:"tim-1"`
}

func TestStructExample(t *testing.T) {
	var o obj

	data := map[string]any{
		"str": "hello",
		"in":  10,
		"fl":  10.5,
		"bol": true,
		"tim": "2024-08-06T00:00:00Z",
	}

	// parse the data
	errs := objSchema.Parse(data, &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
}

func TestStructTags(t *testing.T) {
	var o objTagged

	data := map[string]any{
		"s":     "hello",
		"i":     10,
		"f":     10.5,
		"b":     true,
		"tim-1": "2024-08-06T00:00:00Z",
	}

	errs := objSchema.Parse(data, &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
	assert.Equal(t, o.In, 10)
	assert.Equal(t, o.Fl, 10.5)
	assert.Equal(t, o.Bol, true)
	assert.Equal(t, o.Tim, time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC))
}

var nestedSchema = Struct(Shape{
	"str":    String().Required(),
	"schema": Struct(Shape{"str": String().Required()}),
})

func TestStructNestedStructs(t *testing.T) {

	v := struct {
		Str    string
		Schema struct {
			Str string
		}
	}{
		Str: "hello",
		Schema: struct {
			Str string
		}{},
	}

	m := map[string]any{
		"str":    "hello",
		"schema": map[string]any{"str": "hello"},
	}

	errs := nestedSchema.Parse(m, &v)
	assert.Nil(t, errs)
	assert.Equal(t, v.Str, "hello")
	assert.Equal(t, v.Schema.Str, "hello")

}

func TestStructOptional(t *testing.T) {
	type TestStruct struct {
		Str string `zog:"str"`
		In  int    `zog:"in"`
		Fl  float64
		Bol bool
		Tim time.Time
	}

	var objSchema = Ptr(Struct(Shape{
		"str": String().Required(),
		"in":  Int().Required(),
		"fl":  Float().Required(),
		"bol": Bool().Required(),
		"tim": Time().Required(),
	}))

	var o TestStruct
	errs := objSchema.Parse(nil, &o)
	assert.Nil(t, errs)
}

func TestStructOptionalFields(t *testing.T) {
	type TestStruct struct {
		Str string `zog:"str"`
		In  int    `zog:"in"`
		Fl  float64
		Bol bool
		Tim time.Time
	}

	var objSchema = Struct(Shape{
		"str": String(),
		"in":  Int().Required(),
		"fl":  Float(),
		"bol": Bool().Required(),
		"tim": Time(),
	})

	var o TestStruct
	errs := objSchema.Parse(map[string]any{
		"in":  10,
		"bol": false,
	}, &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.In, 10)
	assert.Equal(t, o.Str, "")
	assert.Equal(t, o.Fl, 0.0)
	assert.Equal(t, o.Bol, false)
	assert.Equal(t, o.Tim, time.Time{})
	errs = objSchema.Parse(map[string]any{}, &o)
	assert.NotEmpty(t, errs)
	inErrs := tutils.FindByPath(errs, "in")
	bolErrs := tutils.FindByPath(errs, "bol")
	assert.NotEmpty(t, inErrs)
	assert.NotEmpty(t, bolErrs)
	assert.Equal(t, zconst.IssueCodeRequired, inErrs[0].Code)
	assert.Equal(t, zconst.IssueCodeRequired, bolErrs[0].Code)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestStructCustomTestInSchema(t *testing.T) {
	type CustomStruct struct {
		Str string `zog:"str"`
		Num int    `zog:"num"`
	}

	// Create a custom test function
	customTest := func(val *int, ctx Ctx) bool {
		// Custom test logic here
		return *val > 0
	}

	// Create a schema with a custom test
	schema := Struct(Shape{
		"str": String().Required(),
		"num": Int().TestFunc(customTest),
	})

	var obj CustomStruct
	data := map[string]any{
		"str": "hello",
		"num": -1,
	}

	errs := schema.Parse(data, &obj)
	assert.NotNil(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	data["num"] = 10
	errs = schema.Parse(data, &obj)
	assert.Nil(t, errs)
	assert.Equal(t, obj.Str, "hello")
	assert.Equal(t, obj.Num, 10)
}

func TestStructCustomTest(t *testing.T) {
	type CustomStruct struct {
		Str string `zog:"str"`
	}

	schema := Struct(Shape{
		"str": String(),
	}).TestFunc(func(val any, ctx Ctx) bool {
		s := val.(*CustomStruct)
		return s.Str == "valid"
	}, Message("customTest"))

	var obj CustomStruct
	data := map[string]any{
		"str": "invalid",
	}

	errs := schema.Parse(data, &obj)
	assert.NotEmpty(t, errs)
	// assert.Equal(t, "customTest", rootErrs[0].Code())
	assert.Equal(t, "customTest", errs[0].Message)
	data["str"] = "valid"
	errs = schema.Parse(data, &obj)
	assert.Nil(t, errs)
}

func TestStructFromIssue(t *testing.T) {
	s := `{
  "nombre": "Juan",
  "apellido": "Perez",
  "email": "test@test.com",
  "alu_id": 25,
  "password": "hunter1"
}`
	var data map[string]any
	err := json.Unmarshal([]byte(s), &data)
	assert.Nil(t, err)

	var output struct {
		Nombre   string `zog:"nombre"`
		Apellido string `zog:"apellido"`
		Email    string `zog:"email"`
		AluID    int    `zog:"alu_id"`
		Password string `zog:"password"`
	}
	schema := Struct(Shape{
		"nombre":   String().Required(),
		"apellido": String().Required(),
		"email":    String().Required(),
		"aluID":    Int().Required(),
		"password": String().Required(),
	})

	// Test with missing fields
	errs := schema.Parse(map[string]any{}, &output)
	assert.NotNil(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	// Test with valid data
	errs = schema.Parse(data, &output)
	assert.Nil(t, errs)
	assert.Equal(t, "Juan", output.Nombre)
	assert.Equal(t, "Perez", output.Apellido)
	assert.Equal(t, "test@test.com", output.Email)
	assert.Equal(t, 25, output.AluID)
	assert.Equal(t, "hunter1", output.Password)
}

func TestStructPanicsOnSchemaMismatch(t *testing.T) {

	var objSchema = Struct(Shape{
		"str":         String().Required(),
		"in":          Int().Required(),
		"fl":          Float().Required(),
		"bol":         Bool().Required(),
		"tim":         Time().Required(),
		"cause_panic": String(),
	})
	var o obj
	data := map[string]any{
		"str": "hello",
		"in":  10,
		"fl":  10.5,
		"bol": true,
		"tim": "2024-08-06T00:00:00Z",
	}
	assert.Panics(t, func() {
		objSchema.Parse(data, &o)
	})
}

func TestStructPostTransforms(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	postTransform := func(val any, ctx Ctx) error {
		if s, ok := val.(*TestStruct); ok {
			s.Value = "post_" + s.Value
		}
		return nil
	}

	schema := Struct(Shape{
		"value": String().Required(),
	}).Transform(postTransform)

	var output TestStruct
	data := map[string]any{"value": "original"}

	errs := schema.Parse(data, &output)
	assert.Nil(t, errs)
	assert.Equal(t, "post_original", output.Value)
}

func TestStructPassThroughRequired(t *testing.T) {
	type TestStruct struct {
		Somefield string
	}

	schema := Struct(Shape{
		"somefield": String().Required(),
	})

	var output TestStruct
	data := map[string]any{
		"somefield": "someValue",
	}

	errs := schema.Parse(data, &output)
	assert.Nil(t, errs)
	assert.Equal(t, "someValue", output.Somefield)
	var output2 TestStruct
	errs = schema.Parse(nil, &output2)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
	assert.NotEmpty(t, tutils.FindByPath(errs, "somefield"))
}

type CustomType int

const (
	Custom1 CustomType = 1
	Custom2 CustomType = 2
)

func TestStructCustomType(t *testing.T) {
	s := struct {
		Custom int
	}{}
	schema := Struct(Shape{
		"custom": Int().OneOf([]int{int(Custom1), int(Custom2)}),
	})
	errs := schema.Parse(map[string]any{"custom": int(Custom1)}, &s)
	assert.Nil(t, errs)
	assert.Equal(t, int(Custom1), s.Custom)
}

type Customs = int

const (
	Customs1 Customs = 1
	Customs2 Customs = 2
)

func TestStructCustomType2(t *testing.T) {
	s := struct {
		Custom Customs
	}{}
	schema := Struct(Shape{
		"custom": Int().OneOf([]int{Customs1, Customs2}),
	})
	errs := schema.Parse(map[string]any{"custom": Customs1}, &s)
	assert.Nil(t, errs)
	assert.Equal(t, Customs1, s.Custom)
}

func TestStructGetType(t *testing.T) {
	s := Struct(Shape{
		"field": String(),
	})
	assert.Equal(t, zconst.TypeStruct, s.getType())
}
