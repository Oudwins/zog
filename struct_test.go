package zog

import (
	"encoding/json"
	"testing"
	"time"

	p "github.com/Oudwins/zog/primitives"
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

var objSchema = Struct(Schema{
	"str": String().Required(),
	"in":  Int().Required(),
	"fl":  Float().Required(),
	"bol": Bool().Required(),
	"tim": Time().Required(),
})

type objTagged struct {
	Str string  `zog:"s"`
	In  int     `zog:"i"`
	Fl  float64 `zog:"f"`
	Bol bool    `zog:"b"`
	Tim time.Time
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
	errs := objSchema.Parse(NewMapDataProvider(data), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
}

func TestStructTags(t *testing.T) {
	var o objTagged

	data := map[string]any{
		"s":   "hello",
		"i":   10,
		"f":   10.5,
		"b":   true,
		"tim": "2024-08-06T00:00:00Z",
	}

	errs := objSchema.Parse(NewMapDataProvider(data), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
	assert.Equal(t, o.In, 10)
	assert.Equal(t, o.Fl, 10.5)
	assert.Equal(t, o.Bol, true)
	assert.Equal(t, o.Tim, time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC))
}

var nestedSchema = Struct(Schema{
	"str":    String().Required(),
	"schema": Struct(Schema{"str": String().Required()}),
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

	errs := nestedSchema.Parse(NewMapDataProvider(m), &v)
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

	var o TestStruct
	errs := objSchema.Parse(&p.EmptyDataProvider{}, &o)
	assert.Nil(t, errs)
}

func TestStructMergeSchema(t *testing.T) {
	var nameSchema = Struct(Schema{
		"name": String().Min(3, Message("Override default message")).Max(10),
	})
	var ageSchema = Struct(Schema{
		"age": Int().GT(18).Required(Message("is required")),
	})
	var schema = nameSchema.Merge(ageSchema)

	type User struct {
		Name string
		Age  int
	}

	var o User
	errs := schema.Parse(NewMapDataProvider(map[string]any{"name": "hello", "age": 20}), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Name, "hello")
	assert.Equal(t, o.Age, 20)
}

func TestStructCustomTest(t *testing.T) {
	type CustomStruct struct {
		Str string `zog:"str"`
		Num int    `zog:"num"`
	}

	// Create a custom test function
	customTest := func(val any, ctx p.ParseCtx) bool {
		// Custom test logic here
		num := val.(int)
		return num > 0
	}

	// Create a schema with a custom test
	schema := Struct(Schema{
		"str": String().Required(),
		"num": Int().Test(TestFunc("customTest", customTest)),
	})

	var obj CustomStruct
	data := map[string]any{
		"str": "hello",
		"num": 10,
	}

	errs := schema.Parse(data, &obj)
	assert.Nil(t, errs)
	assert.Equal(t, obj.Str, "hello")
	assert.Equal(t, obj.Num, 10)
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
	json.Unmarshal([]byte(s), &data)

	var output struct {
		Nombre   string `zog:"nombre"`
		Apellido string `zog:"apellido"`
		Email    string `zog:"email"`
		AluID    int    `zog:"alu_id"`
		Password string `zog:"password"`
	}
	schema := Struct(Schema{
		"nombre":   String().Required(Message("this doesn't display even if validation fails")),
		"apellido": String().Required(),
		"email":    String().Required(),
		"aluID":    Int().Required(),
		"password": String().Required(),
	})
	errs := schema.Parse(data, &output)
	assert.Nil(t, errs)
	assert.Equal(t, "Juan", output.Nombre)
	assert.Equal(t, "Perez", output.Apellido)
	assert.Equal(t, "test@test.com", output.Email)
	assert.Equal(t, 25, output.AluID)
	assert.Equal(t, "hunter1", output.Password)
}

func TestStructPanicsOnSchemaMismatch(t *testing.T) {

	var objSchema = Struct(Schema{
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
