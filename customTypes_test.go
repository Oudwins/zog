package zog

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AString string
type Aint int

const (
	a = AString("A")
	b = AString("B")
	c = Aint(1)
	d = Aint(2)
)

func CustomInt() *NumberSchema[Aint] {
	return &NumberSchema[Aint]{}
}

// This works very nicely for numbers! At least for now. Since we are not using any functions that expect anything other than a comparable
func TestCustomIntWithoutConversion(t *testing.T) {
	val := Aint(1)
	schema := CustomInt().OneOf([]Aint{c, d}).LTE(1).GT(0)
	errs := schema.Validate(&val)
	assert.Nil(t, errs)
	assert.Equal(t, c, val)
}

func TestCustomString(t *testing.T) {
	var v AString = " ABB "

	schema := (&StringSchema[AString]{}).Contains(a).HasPrefix(a).HasSuffix(b).Len(3).Trim()
	errs := schema.Validate(&v)
	assert.Nil(t, errs)
	assert.Equal(t, AString("ABB"), v)
}

func TestCustomType(t *testing.T) {
	// This approach will work doing this before and then after and just treating the thing as a string for the rest
	var v AString = "A"
	var x any = &v
	refval := reflect.ValueOf(x)
	val := refval.Elem().String()
	fmt.Println(val)
	refval.Elem().SetString("B")
	val = refval.Elem().String()
	fmt.Println(val)
}
