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

This is quite simple to do for the supported primitive types (string, number, boolean). Here is an example:

```go
// definition in your code
type Env string

const (
	Prod Env = "prod"
	Dev  Env = "env"
)

func EnvSchema() *StringSchema[Env] {
	s := &z.StringSchema[Env]
	return s.OneOf([]Env{Prod, Dev}) // you can also just return the schema and define the tests when calling it it doesn't matter
}

// usage
type S struct {
	Environment Env
}

schema := z.Struct(
	z.Schema{
		"Environment": EnvSchema(), // All string methods will now be typed to Env type
	},
)
```

This becomes a little more complex if you need to use `Parse` instead of just `Validate` since you need to define a custom `Coercer` function. Here is what I would recommend and it is also very similar to the way Zog creates the schemas you use:

```go
// Definition
func EnvSchema(opts ...z.SchemaOption) *StringSchema[Env] {
	s := &StringSchema[Env]{}
	ops = append([]z.SchemaOption{
		// This is required if you want to use Parse since we don't use reflection to set the value you need to coerce it manually
		WithCoercer(func(x any) (any, error) {
			v, e := conf.DefaultCoercers.String(x)
			if e != nil {
				return nil, e
			}
			return Env(v.(string)), nil
		}),
	}, opts...)

	for _, op := range ops {
		op(s)
	}
	return s
}

// Usage is the same as before
```

> Why is this so verbose?
> Although we considered introducing an API that would allow you to define this types of schemas in a more concise way (and we may still do so), to keep code consistency & reusability we recommend that you make a factory function like the one above for your custom types. And we felt that providing a simpler API could lead to people just inlining the schema's which would make it impossible to reuse them.

## Quick and Dirty Custom Schema

Sometimes you may want to create a custom schema for a type that is not a primitive and you don't want to go through the process of defining everything needed to create a full schema. You just want to run a validation inside Zog. Zog supports a simple way to do this using the `CustomFunc` function which looks like this:

```go
// fn signature
func CustomFunc[T any](fn func(valPtr *T, ctx z.Ctx) bool, opts ...z.TestOption) *z.Custom[T]
```

Usage is very similar to the `schema.TestFunc()` function:

```go

user := User{
	"uuid": z.CustomFunc(func(valPtr *uuid.UUID, ctx z.Ctx) bool {
		return (*valPtr).IsValid()
	}, z.Message("invalid uuid"))
}
```

> **Limitations**
>
> - CustomFunc doesn't support type coercion yet. You can still use it with parse but it will not be able to coerce the type.
>   **Why is valPtr a pointer?**
>   Mainly for performance reasons. It is faster in almost every case to pass a pointer to the value than the value itself. This is specially true if the value is a large struct.
