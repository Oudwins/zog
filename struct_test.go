package zog

import (
	"fmt"
	"testing"
	"time"

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

func TestExampleStruct(t *testing.T) {
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

	errs := objSchema.Parse(NewMapDataProvider(data), &o)
	assert.Nil(t, errs)
	assert.Equal(t, o.Str, "hello")
	assert.Equal(t, o.In, 10)
	assert.Equal(t, o.Fl, 10.5)
	assert.Equal(t, o.Bol, true)
	assert.Equal(t, o.Tim, time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC))
}
