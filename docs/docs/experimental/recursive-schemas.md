---
sidebar_position: 3
---

# Recursive Schemas

> **⚠️ Experimental API**: This feature is experimental and subject to breaking changes. Use with caution.

Define schemas that reference themselves recursively, enabling validation of recursive data structures like linked lists, trees, and nested hierarchies.

## Overview

Recursive schemas solve the problem of defining schemas for data structures that contain references to themselves. Without recursive schemas, you cannot define a schema for a struct that contains a field of its own type.

## Basic Usage

Use `EXPERIMENTAL_RECURSIVE` to create a recursive schema:

```go
import z "github.com/Oudwins/zog"

type Node struct {
    Value int
    Self  *Node
}

var nodeSchema = z.EXPERIMENTAL_RECURSIVE(func(self z.RecursiveSchema[*z.PointerSchema]) *z.PointerSchema {
    return z.Ptr(z.Struct(z.Shape{
        "value": z.Int().Required(),
        "self":  self(),
    }))
})
```

The `self` parameter is a function that returns a schema representing the recursive type. Call `self()` wherever you need the recursive reference.

### Updating Recursive Schemas

You can modify the recursive schema using updater functions. This is mostly intended to be used to clone a schema and apply additional constraints to it. If you make changes to the schema without cloning it you can land on undefined behavior land. Which maybe means this API should change and clone automatically for you.

```go
var nodeSchema = z.EXPERIMENTAL_RECURSIVE(func(self z.RecursiveSchema[*z.PointerSchema]) *z.PointerSchema {
    return z.Ptr(z.Struct(z.Shape{
        "value": z.Int().Required(),
        "self": self(func(original *z.PointerSchema) *z.PointerSchema {
            // Modify the original schema
            return original
        }),
    }))
})
```

## Important Notes

- **Experimental**: API may change in future versions
- **Lazy Initialization**: The recursive schema is only materialized when first accessed, so first call will be slightly slower
- **Thread-Safe**: Safe for concurrent use across multiple goroutines
- **Type Safety**: Ensure your Go struct matches the schema structure
- **Nil Handling**: Recursive references can be `nil` (use `z.Ptr()` for optional fields)
