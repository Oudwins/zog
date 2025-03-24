---
sidebar_position: 4
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Creating Custom Tests

> Please read the [Anatomy of a Schema](/core-concepts/anatomy-of-schema) page before continuing.

## Simple Custom Tests - Aka Zod's `refine`

All schemas contain the `TestFunc()` method which can be used to create a simple custom test in a similar way to Zod's `refine` method. The `TestFunc()` method takes a `ValidateFunc` as an argument. This is a function that takes the data as input and returns a boolean indicating if it is valid or not. If you return `false` from the function Zog will create a [ZogIssue](/errors). For example:

```go
z.String().TestFunc(func(data any, ctx z.Ctx) bool {
	return data == "test"
})
```

Test funcs for structs and slices instead receive a pointer to the data to avoid copying large data structures. For example:

```go
z.Struct(z.Schema{
	"name": z.String(),
}).TestFunc(func(dataPtr any, ctx z.Ctx) bool {
	user := dataPtr.(*User)
	return user.Name == "test"
})
```

> **Pro tip**
> It is very likely that you may want to set custom messages or paths, you can do that like with any other tests with `TestOptions`. For more on this checkout the [Anatomy of a Schema](/core-concepts/anatomy-of-schema#test-options) page.
