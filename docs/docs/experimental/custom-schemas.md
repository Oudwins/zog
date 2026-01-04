---
sidebar_position: 1
---

# Custom Schemas (Experimental)

> **⚠️ Experimental API**: This feature is experimental and subject to breaking changes. Use with caution.

Create fully custom schemas by implementing the `EXPERIMENTAL_PUBLIC_ZOG_SCHEMA` interface. This gives you complete control over parsing, validation, and type coercion.

## The Interface

Implement these four methods:

```go
type EXPERIMENTAL_PUBLIC_ZOG_SCHEMA interface {
    Process(ctx *internals.SchemaCtx)      // Called during Parse()
    Validate(ctx *internals.SchemaCtx)     // Called during Validate()
    GetType() zconst.ZogType       // Return schema type identifier
    SetCoercer(c CoercerFunc)      // Optional: set type coercion function
}
```

## Key Concepts

### Process vs Validate

- **`Process(ctx *internals.SchemaCtx)`**: Handles parsing. You must:

  - Coerce `ctx.Data` to your target type
  - Assign the value to `ctx.ValPtr` (always a pointer)
  - Run validation and add issues if validation fails

- **`Validate(ctx *internals.SchemaCtx)`**: Handles validation. The value is already at `ctx.ValPtr`, just validate it.

### SchemaCtx Essentials

- **`ctx.Data`**: Input data (only during `Process()`)
- **`ctx.ValPtr`**: Pointer to destination value (always a pointer)
- **`ctx.AddIssue(issue)`**: Add a validation error
- **`ctx.Issue()`**: Create a new issue with context pre-filled
- **`ctx.IssueFromCoerce(err)`**: Create an issue from a coercion error

## Example

Here's a complete example for a custom string schema:

```go
import (
    "fmt"

    z "github.com/Oudwins/zog"
    p "github.com/Oudwins/zog/internals"
    "github.com/Oudwins/zog/zconst"
)

type MinLengthSchema struct {
    minLength int
    errorMsg  string
    coercer   z.CoercerFunc
}

func (s *MinLengthSchema) Process(ctx *internals.SchemaCtx) {
    // Optional: handle coercion
    if s.coercer != nil {
        coerced, err := s.coercer(ctx.Data)
        if err != nil {
            ctx.AddIssue(ctx.IssueFromCoerce(err))
            return
        }
        ctx.Data = coerced
    }

    // Type assertion
    ptr, ok := ctx.ValPtr.(*string)
    if !ok {
        ctx.AddIssue(ctx.IssueFromCoerce(
            fmt.Errorf("expected *string, got %T", ctx.ValPtr)))
        return
    }

    val, ok := ctx.Data.(string)
    if !ok {
        ctx.AddIssue(ctx.IssueFromCoerce(
            fmt.Errorf("expected string, got %T", ctx.Data)))
        return
    }

    // Assign value
    *ptr = val

    // Validate
    if len(val) < s.minLength {
        issue := ctx.Issue().SetMessage(s.errorMsg)
        ctx.AddIssue(issue)
    }
}

func (s *MinLengthSchema) Validate(ctx *internals.SchemaCtx) {
    ptr, ok := ctx.ValPtr.(*string)
    if !ok {
        return
    }

    if len(*ptr) < s.minLength {
        issue := ctx.Issue().SetMessage(s.errorMsg)
        ctx.AddIssue(issue)
    }
}

func (s *MinLengthSchema) GetType() zconst.ZogType {
    return zconst.TypeString
}

func (s *MinLengthSchema) SetCoercer(c z.CoercerFunc) {
    s.coercer = c
}

// Usage
schema := z.Struct(z.Shape{
    "name": z.Use(&MinLengthSchema{
        minLength: 5,
        errorMsg:  "name must be at least 5 characters",
    }),
})
```

## Using Custom Schemas

Wrap your implementation with `z.Use()`:

```go
customSchema := &MyCustomSchema{...}
schema := z.Use(customSchema)
```

Use it anywhere a Zog schema is expected:

```go
// In structs
z.Struct(z.Shape{"id": z.Use(customSchema)})

// In slices
z.Slice(z.Use(customSchema))

// In pointers
z.Ptr(z.Use(customSchema))
```

## Type Coercion

Set a coercer to handle type conversion:

```go
type CoercerFunc func(original any) (value any, err error)

customSchema.SetCoercer(func(original any) (value any, err error) {
    if i, ok := original.(int); ok {
        return strconv.Itoa(i), nil
    }
    return nil, fmt.Errorf("cannot convert %T to string", original)
})
```

The coercer is called during `Process()` before type assertion.

## Important Notes

- **Experimental**: API may change in future versions
- **Type Safety**: You're responsible for proper type assertions
- **Error Handling**: Always use `ctx.AddIssue()`, don't panic
- **ValPtr**: Always a pointer to the destination value

For simpler validation needs, consider `z.CustomFunc()` instead.
