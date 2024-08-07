package zog

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// structs with pointers
// maps with additional values
// errors are correct
// panics are correct

var objSchema = Struct(Schema{
	"str":  String().Required(),
	"in":   Int().Required(),
	"fl":   Float().Required(),
	"bol":  Bool().Required(),
	"slic": Slice(String().Required()),
})

type obj struct {
	Str  string
	In   int
	Fl   float64
	Bol  bool
	Slic []string
}

type objTagged struct {
	Str  string   `zog:"s"`
	In   int      `zog:"i"`
	Fl   float64  `zog:"f"`
	Bol  bool     `zog:"b"`
	Slic []string `zog:"sl"`
	Tim  time.Time
}

func TestExampleStruct(t *testing.T) {
	var o obj

	data := map[string]any{
		"str": "hello",
		"in":  10,
		"fl":  10.5,
		"bol": true,
	}

	// parse the data
	errs := objSchema.Parse(data, &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
}

func TestTaggedStruct(t *testing.T) {
	var o objTagged
	fmt.Println(o.Tim)

	data := map[string]any{
		"s":   "hello",
		"i":   10,
		"f":   10.5,
		"b":   true,
		"tim": "2024-08-06T00:00:00Z",
	}

	errs := objSchema.Parse(data, &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
	assert.Equal(t, o.In, 10)
	assert.Equal(t, o.Fl, 10.5)
	assert.Equal(t, o.Bol, true)
	// fmt.Println(o.Tim.Format(time.RFC3339))

	// assert.Equal(t, o.Tim.Format(time.RFC3339), "2024-08-06T00:00:00Z")
	v, _ := time.Parse(time.RFC3339, "2024-08-06T00:00:00Z")

	tschema := Time().Required()

	var tim time.Time
	tschema.Parse("2024-08-06T00:00:00Z", &tim)

	fmt.Println(v, tim)
}
