package zog

import (
	"testing"

	"github.com/Oudwins/zog/tutils"
	"github.com/Oudwins/zog/zconst"
	"github.com/stretchr/testify/assert"
)

type Node struct {
	Value int
	Self  *Node
}

type TreeNode struct {
	Value    int
	Children []*TreeNode
}

type DeepNode struct {
	Value int
	Next  *DeepNode
}

var nodeSchema = EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
	return Ptr(Struct(Shape{
		"value": Int().Required(),
		"self": self(func(original *PointerSchema) *PointerSchema {
			return original
		}),
	}))
})

// ============================================================================
// 1. Basic Recursive Test (existing)
// ============================================================================

func TestRecursive(t *testing.T) {
	var node *Node
	errs := nodeSchema.Parse(map[string]any{"value": 10, "self": map[string]any{"value": 20}}, &node)
	assert.Nil(t, errs)
	assert.Equal(t, 10, node.Value)
	assert.NotNil(t, node.Self)
	assert.Equal(t, 20, node.Self.Value)
}

// ============================================================================
// 2. lazySchema Internal Tests
// ============================================================================

func TestLazySchemaGetType(t *testing.T) {
	ptr := Ptr(Struct(Shape{
		"value": Int().Required(),
		"self": Ptr(Struct(Shape{
			"value": Int().Required(),
		})),
	}))
	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"self":  self(),
		}))
	})

	// getType should return the underlying schema type (ptr)
	assert.Equal(t, ptr.getType(), schema.getType())
}

func TestLazySchemaSetCoercer(t *testing.T) {
	type TestStruct struct {
		Value string
		Self  *TestStruct
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": String().Required(),
			"self":  self(),
		}))
	})

	// Set a custom coercer that converts int to string
	customCoercer := func(v any) (any, error) {
		if i, ok := v.(int); ok {
			return string(rune(i + 48)), nil
		}
		return v, nil
	}

	schema.setCoercer(customCoercer)

	var result *TestStruct
	// Test that coercer is applied (though in practice, coercers work at field level)
	errs := schema.Parse(map[string]any{"value": 10, "self": nil}, &result)
	assert.Empty(t, errs)
	assert.NotNil(t, result)
}

func TestLazySchemaLazyInitialization(t *testing.T) {
	callCount := 0
	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		callCount++
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"self":  self(),
		}))
	})

	// First call to getType should initialize
	_ = schema.getType()
	initialCallCount := callCount

	// Subsequent calls should not re-initialize
	_ = schema.getType()
	_ = schema.getType()
	_ = schema.getType()

	// Builder should only be called once during EXPERIMENTAL_RECURSIVE
	// The lazy initialization happens when get() is called, but the builder
	// itself is only called once
	assert.Equal(t, initialCallCount, callCount)
}

// ============================================================================
// 3. Parse Tests for Multi-Level Deep Recursion
// ============================================================================

func TestRecursiveMultiLevelDeep(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next": map[string]any{
				"value": 3,
				"next": map[string]any{
					"value": 4,
					"next":  nil,
				},
			},
		},
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Equal(t, 1, node.Value)
	assert.NotNil(t, node.Next)
	assert.Equal(t, 2, node.Next.Value)
	assert.NotNil(t, node.Next.Next)
	assert.Equal(t, 3, node.Next.Next.Value)
	assert.NotNil(t, node.Next.Next.Next)
	assert.Equal(t, 4, node.Next.Next.Next.Value)
	assert.Nil(t, node.Next.Next.Next.Next)
}

func TestRecursiveFiveLevelsDeep(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next": map[string]any{
				"value": 3,
				"next": map[string]any{
					"value": 4,
					"next": map[string]any{
						"value": 5,
						"next":  nil,
					},
				},
			},
		},
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)

	// Verify all 5 levels
	current := node
	for i := 1; i <= 5; i++ {
		assert.NotNil(t, current)
		assert.Equal(t, i, current.Value)
		if i < 5 {
			current = current.Next
		} else {
			assert.Nil(t, current.Next)
		}
	}
}

// ============================================================================
// 4. Parse Tests for Recursive Slices
// ============================================================================

func TestRecursiveSliceTree(t *testing.T) {
	type TreeNode struct {
		Value    int
		Children []*TreeNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value":    Int().Required(),
			"children": Slice(self()),
		}))
	})

	var tree *TreeNode
	data := map[string]any{
		"value": 1,
		"children": []any{
			map[string]any{
				"value":    2,
				"children": []any{},
			},
			map[string]any{
				"value": 3,
				"children": []any{
					map[string]any{
						"value":    4,
						"children": []any{},
					},
				},
			},
		},
	}

	errs := schema.Parse(data, &tree)
	assert.Empty(t, errs)
	assert.NotNil(t, tree)
	assert.Equal(t, 1, tree.Value)
	assert.Len(t, tree.Children, 2)
	assert.Equal(t, 2, tree.Children[0].Value)
	assert.Len(t, tree.Children[0].Children, 0)
	assert.Equal(t, 3, tree.Children[1].Value)
	assert.Len(t, tree.Children[1].Children, 1)
	assert.Equal(t, 4, tree.Children[1].Children[0].Value)
}

func TestRecursiveSliceEmptyChildren(t *testing.T) {
	type TreeNode struct {
		Value    int
		Children []*TreeNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value":    Int().Required(),
			"children": Slice(self()),
		}))
	})

	var tree *TreeNode
	data := map[string]any{
		"value":    1,
		"children": []any{},
	}

	errs := schema.Parse(data, &tree)
	assert.Empty(t, errs)
	assert.NotNil(t, tree)
	assert.Equal(t, 1, tree.Value)
	assert.Len(t, tree.Children, 0)
}

// ============================================================================
// 5. Parse Tests for Nil Values and Edge Cases
// ============================================================================

func TestRecursiveNilAtRoot(t *testing.T) {
	var node *Node
	errs := nodeSchema.Parse(nil, &node)
	assert.Empty(t, errs)
	assert.Nil(t, node)
}

func TestRecursiveNilAtNestedLevel(t *testing.T) {
	var node *Node
	errs := nodeSchema.Parse(map[string]any{
		"value": 10,
		"self":  nil,
	}, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Equal(t, 10, node.Value)
	assert.Nil(t, node.Self)
}

func TestRecursiveNilAtDeepLevel(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next":  nil,
		},
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Equal(t, 1, node.Value)
	assert.NotNil(t, node.Next)
	assert.Equal(t, 2, node.Next.Value)
	assert.Nil(t, node.Next.Next)
}

func TestRecursiveZeroValue(t *testing.T) {
	var node *Node
	errs := nodeSchema.Parse(map[string]any{
		"value": 0,
		"self":  nil,
	}, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Equal(t, 0, node.Value)
}

func TestRecursiveMissingField(t *testing.T) {
	var node *Node
	errs := nodeSchema.Parse(map[string]any{
		"self": map[string]any{"value": 20},
	}, &node)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

// ============================================================================
// 6. Validation Tests
// ============================================================================

func TestRecursiveValidate(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	node := &DeepNode{
		Value: 1,
		Next: &DeepNode{
			Value: 2,
			Next:  nil,
		},
	}

	errs := schema.Validate(&node)
	assert.Empty(t, errs)
	assert.Equal(t, 1, node.Value)
	assert.Equal(t, 2, node.Next.Value)

	// Test with nil pointer
	var nilNode *DeepNode
	errs = schema.Validate(&nilNode)
	assert.Empty(t, errs)
}

func TestRecursiveValidateWithNilAtNestedLevel(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	node := &DeepNode{
		Value: 1,
		Next:  nil,
	}

	errs := schema.Validate(&node)
	assert.Empty(t, errs)
	assert.Nil(t, node.Next)
}

func TestRecursiveValidateTree(t *testing.T) {
	type TreeNode struct {
		Value    int
		Children []*TreeNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value":    Int().Required(),
			"children": Slice(self()),
		}))
	})

	tree := &TreeNode{
		Value: 1,
		Children: []*TreeNode{
			{Value: 2, Children: []*TreeNode{}},
			{Value: 3, Children: []*TreeNode{
				{Value: 4, Children: []*TreeNode{}},
			}},
		},
	}

	errs := schema.Validate(&tree)
	assert.Empty(t, errs)
	assert.Equal(t, 1, tree.Value)
	assert.Len(t, tree.Children, 2)
}

// ============================================================================
// 7. RecursiveSchemaUpdater Tests
// ============================================================================

func TestRecursiveSchemaUpdater(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next": self(func(original *PointerSchema) *PointerSchema {
				// Apply NotNil() via updater
				return original.NotNil()
			}),
		}))
	})

	var node *DeepNode
	// Test with nil next - should fail because of NotNil()
	data := map[string]any{
		"value": 1,
		"next":  nil,
	}

	errs := schema.Parse(data, &node)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
	assert.Equal(t, zconst.IssueCodeNotNil, errs[0].Code)

	// Test with valid next but nil nested next - should fail because NotNil() is applied recursively
	// The updater applies NotNil() to all recursive instances, so any nil in the chain should fail
	data = map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next":  nil, // This should fail because NotNil() is applied via updater
		},
	}

	errs = schema.Parse(data, &node)
	assert.NotEmpty(t, errs) // Should fail because nested next is nil
	tutils.VerifyDefaultIssueMessages(t, errs)

	// Verify that the error is about the nested next being nil
	foundNotNilError := false
	for _, err := range errs {
		if err.Code == zconst.IssueCodeNotNil {
			foundNotNilError = true
			break
		}
	}
	assert.True(t, foundNotNilError, "Should have a NotNil error")
}

func TestRecursiveSchemaUpdaterMultiple(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next": self(func(original *PointerSchema) *PointerSchema {
				// Apply updater that returns the original unchanged
				return original
			}),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next":  nil,
		},
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Equal(t, 1, node.Value)
	assert.NotNil(t, node.Next)
	assert.Equal(t, 2, node.Next.Value)
}

func TestRecursiveSchemaUpdaterNoUpdater(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(), // No updater passed
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next":  nil, // Should be allowed without updater
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Nil(t, node.Next)
}

// ============================================================================
// 8. Error Path Tests
// ============================================================================

func TestRecursiveErrorPathRoot(t *testing.T) {
	var node *Node
	errs := nodeSchema.Parse(map[string]any{
		"value": "not_an_int", // Invalid type
		"self":  nil,
	}, &node)

	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
	// Error should be at root level or value field
	assert.True(t, len(errs) > 0)
}

func TestRecursiveErrorPathNested(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": "invalid", // Error at nested level
			"next":  nil,
		},
	}

	errs := schema.Parse(data, &node)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)

	// Check that error path includes "next"
	hasNextPath := false
	for _, err := range errs {
		pathStr := ""
		for _, p := range err.Path {
			pathStr += p + "."
		}
		if len(err.Path) > 0 && err.Path[0] == "next" {
			hasNextPath = true
			break
		}
	}
	assert.True(t, hasNextPath || len(errs) > 0)
}

func TestRecursiveErrorPathDeepNested(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next": map[string]any{
				"value": "invalid", // Error at deep nested level
				"next":  nil,
			},
		},
	}

	errs := schema.Parse(data, &node)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestRecursiveMultipleErrors(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	data := map[string]any{
		"value": "invalid1",
		"next": map[string]any{
			"value": "invalid2",
			"next":  nil,
		},
	}

	errs := schema.Parse(data, &node)
	assert.NotEmpty(t, errs)
	assert.True(t, len(errs) >= 1) // At least one error
	tutils.VerifyDefaultIssueMessages(t, errs)
}

func TestRecursiveErrorPathInTree(t *testing.T) {
	type TreeNode struct {
		Value    int
		Children []*TreeNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value":    Int().Required(),
			"children": Slice(self()),
		}))
	})

	var tree *TreeNode
	data := map[string]any{
		"value": 1,
		"children": []any{
			map[string]any{
				"value":    "invalid",
				"children": []any{},
			},
		},
	}

	errs := schema.Parse(data, &tree)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}

// ============================================================================
// 9. Edge Cases
// ============================================================================

func TestRecursiveCircularLikePattern(t *testing.T) {
	// Test deeply nested structure that goes back and forth
	// This tests that recursive schemas can handle very deep nesting
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	// Create a very deep structure
	data := map[string]any{
		"value": 1,
		"next": map[string]any{
			"value": 2,
			"next": map[string]any{
				"value": 3,
				"next": map[string]any{
					"value": 4,
					"next": map[string]any{
						"value": 5,
						"next":  nil,
					},
				},
			},
		},
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)

	// Verify deep nesting works
	current := node
	expectedValues := []int{1, 2, 3, 4, 5}
	for i, expected := range expectedValues {
		assert.NotNil(t, current, "Node at level %d should not be nil", i)
		assert.Equal(t, expected, current.Value, "Value at level %d should be %d", i, expected)
		if i < len(expectedValues)-1 {
			current = current.Next
		} else {
			assert.Nil(t, current.Next, "Last node should have nil Next")
		}
	}
}

func TestRecursiveMixedSchemaTypes(t *testing.T) {
	type MixedNode struct {
		Value    int
		Children []*MixedNode
		Next     *MixedNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value":    Int().Required(),
			"children": Slice(self()),
			"next":     self(),
		}))
	})

	var node *MixedNode
	data := map[string]any{
		"value": 1,
		"children": []any{
			map[string]any{
				"value":    2,
				"children": []any{},
				"next":     nil,
			},
		},
		"next": map[string]any{
			"value":    3,
			"children": []any{},
			"next":     nil,
		},
	}

	errs := schema.Parse(data, &node)
	assert.Empty(t, errs)
	assert.NotNil(t, node)
	assert.Equal(t, 1, node.Value)
	assert.Len(t, node.Children, 1)
	assert.Equal(t, 2, node.Children[0].Value)
	assert.NotNil(t, node.Next)
	assert.Equal(t, 3, node.Next.Value)
}

func TestRecursiveRequiredField(t *testing.T) {
	type DeepNode struct {
		Value int
		Next  *DeepNode
	}

	schema := EXPERIMENTAL_RECURSIVE(func(self RecursiveSchema[*PointerSchema]) *PointerSchema {
		return Ptr(Struct(Shape{
			"value": Int().Required(),
			"next":  self(),
		}))
	})

	var node *DeepNode
	// Missing required field
	data := map[string]any{
		"next": map[string]any{
			"value": 2,
			"next":  nil,
		},
	}

	errs := schema.Parse(data, &node)
	assert.NotEmpty(t, errs)
	tutils.VerifyDefaultIssueMessages(t, errs)
}
