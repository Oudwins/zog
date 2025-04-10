package zog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var custom = Custom(func(val any, ctx Ctx) any {
	i, ok := val.(int)
	if ok {
		return fmt.Sprintf("%d", i)
	}
	return "not an int"
})

func TestCustomSchemaFnSimple(t *testing.T) {
	var out string
	custom.Parse(1, &out)
	assert.Equal(t, "1", out)
	custom.Validate(&out)
	assert.Equal(t, "not an int", out)
}

type customTesting struct {
	Hello string
}

func TestCustomSchemaFnComplex(t *testing.T) {
	s := Struct(Schema{
		"Hello": custom,
	})

	var out customTesting

	s.Parse(map[string]any{
		"Hello": 1,
	}, &out)

	assert.Equal(t, "1", out.Hello)
	s.Validate(&out)
	assert.Equal(t, "not an int", out.Hello)
}

/* With custom function you can already do:
z.Custom(func(val any, ctx z.Ctx) any {
// grab the value from decimal.Decimal
// execute validate on the value you grabbed
// return nil
})


But the issue here is that now you have to merge the errors...



Also I want the API for the PublicZogSchema. Maybe CustomSchema(). Then we would have:
- z.Custom() -> for the fn
- z.Custom.String() -> for the stringlike
- z.CustomSchema() -> for actual full on schemas
*/
