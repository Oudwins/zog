package zog

import (
	"fmt"
	"testing"
	"time"
)

type Msgs struct {
	Name  string
	Age   int
	Time  time.Time
	Bool  bool
	Slice []string
}

func TestErrorMessages(t *testing.T) {
	// make schema
	schema := Struct(Schema{
		"name":  String().Min(3).Max(10).Required(),
		"age":   Int().GT(18).Required(),
		"time":  Time().Before(time.Now()).Required(),
		"bool":  Bool().True().Required(),
		"slice": Slice(String()).Contains("foo").Required(),
	})

	var u Msgs

	errs := schema.Parse(map[string]any{"name": "0", "age": 0, "time": "2020-01-01T00:00:00Z", "bool": false, "slice": []string{"x"}}, &u)
	sanitized := Errors.SanitizeMap(errs)
	fmt.Println(sanitized)
}
