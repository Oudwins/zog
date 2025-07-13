---
sidebar_position: 5
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Creating Custom Schemas

> Please read the [Anatomy of a Schema](/core-concepts/anatomy-of-schema) page before continuing.

Currently Zog plans to support three different ways of creating custom schemas. Although this is subject to change and some of these are not yet implemented so keep an eye out on this page as it gets updated, [more details on my thoughts here](https://github.com/Oudwins/zog/discussions/132).

1. Generics on Primitive Schemas for custom strings, numbers, booleans, etc...
2. Custom Schemas Interface you can implement to create a 100% custom schema (Not yet implemented)
3. A system by which you can define a schema for a custom type or interface that after some transformation can become a normal zog schema.

## Creating Custom Schemas for Primitive Types

This is quite simple to do for the supported primitive types (string, number, boolean). For complete list of options see the [reference](/reference) page. Here is an example:

```go
// definition in your code
type Env string

const (
	Prod Env = "prod"
	Dev  Env = "env"
)

type ActiveInactive int
const (
	Active int = 1
	Inactive int = 0
)

func EnvSchema() *StringSchema[Env] {
	return z.StringLike[Env]().OneOf([]Env{Prod, Dev})
}

// usage
type S struct {
	Environment Env
}

schema := z.Struct(
	z.Shape{
		"Environment": EnvSchema(), // All string methods will now be typed to Env type
		"active": z.IntLike[ActiveInactive]().OneOf([]ActiveInactive{Active, Inactive}),
	},
)
```

## Quick and Dirty Custom Schema

Sometimes you may want to create a custom schema for a type that is not a primitive and you don't want to go through the process of defining everything needed to create a full schema. You just want to run a validation inside Zog. Zog supports a simple way to do this using the `CustomFunc` function which looks like this:

```go
// fn signature
func CustomFunc[T any](fn func(valPtr *T, ctx z.Ctx) bool, opts ...z.TestOption) *z.Custom[T]
```

Usage is very similar to the `schema.TestFunc()` function:

```go

user := z.Struct(z.Shape{
	"uuid": z.CustomFunc(func(valPtr *uuid.UUID, ctx z.Ctx) bool {
		return (*valPtr).IsValid()
	}, z.Message("invalid uuid"))
})
```

> **Limitations**
>
> - CustomFunc doesn't support type coercion yet. You can still use it with parse but it will not be able to coerce the type.
>   **Why is valPtr a pointer?**
>   Mainly for performance reasons. It is faster in almost every case to pass a pointer to the value than the value itself. This is specially true if the value is a large struct.
