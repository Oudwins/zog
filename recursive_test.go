package zog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Node struct {
	Value int
	Self  *Node
}

var nodeSchema = Recursive(func(self ZogSchema) *PointerSchema {
	return Ptr(Struct(Shape{
		"value": Int().Required(),
		"self":  self,
	}))
})

func TestRecursive(t *testing.T) {

	var node *Node
	errs := nodeSchema.Parse(map[string]any{"value": 10, "self": map[string]any{"value": 20}}, &node)
	assert.Nil(t, errs)
	assert.Equal(t, 10, node.Value)
	assert.NotNil(t, node.Self)
	assert.Equal(t, 20, node.Self.Value)
}
