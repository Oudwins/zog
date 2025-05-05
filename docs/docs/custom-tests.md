---
sidebar_position: 5
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Creating Custom Tests

> Please read the [Anatomy of a Schema](/core-concepts/anatomy-of-schema) page before continuing.

## Simple Custom Tests - Aka Zod's `refine`

All schemas contain the `TestFunc()` method which can be used to create a simple custom test in a similar way to Zod's `refine` method. The `TestFunc()` method takes a `ValidateFunc` as an argument. This is a function that takes the data as input and returns a boolean indicating if it is valid or not. If you return `false` from the function Zog will create a [ZogIssue](/errors). For example:

```go
z.String().TestFunc(func(data *string, ctx z.Ctx) bool { // notice that here Zog already knows you need to pass a *string to the test.
	return *data == "test"
})
```

Test funcs for structs and slices instead receive a pointer to the data to avoid copying large data structures. For example:

```go
z.Struct(z.Shape{
	"name": z.String(),
}).TestFunc(func(dataPtr any, ctx z.Ctx) bool { // notice that here we have to cast the dataPtr because no inference for struct types
	user := dataPtr.(*User)
	return user.Name == "test"
})
```

> **Pro tip**
> It is very likely that you may want to set custom messages or paths, you can do that like with any other tests with `TestOptions`. For more on this checkout the [Anatomy of a Schema](/core-concepts/anatomy-of-schema#test-options) page.

## Complex Custom Tests - Aka Zod's `superRefine`

For complex tests you can use the `schema.TestFunc()` method but it is recommended that you use the `schema.Test()` method as it provides a bit more flexbility. This is quite simple to do using the [zog context](/context), lets look at an example that will execute a DB call to verify a user's session is valid:

```go
sessionSchema := z.String().Test(z.Test{
  Func: func (val any, ctx z.Ctx) {
    session := val.(string)
    if !sessionStore.IsValid(session) {
      // This ctx.Issue() is a shortcut to creating Zog issues that are aware of the current schema context. Basically this means that it will prefil some data like the path, value, etc. for you.
      ctx.AddIssue(ctx.Issue().SetMessage("Invalid session"))
      return
    }
    if sessionStore.HasExpired(session) {
      // But you can also just use the normal z.Issue{} struct if you want to.
      ctx.AddIssue(z.Issue{
        Message: "Session expired",
        Path: "session",
        Value: val,
      })
      return
    }
    if sessionStore.IsRevoked(session) {
      ctx.AddIssue(ctx.Issue().SetMessage("Session revoked"))
      return
    }
    // etc
  }
})
```

## Making Reusable Tests

In general I recommend you wrap you reusable tests in a function. Here are examples for both simple and complex tests:

```go

// Notice how we can pass default values to the test which can then be overriden by the function called. This is super nice if you need it!
func MySimpleTest(opts ...z.TestOption) z.Test[any] {
	options := []TestOption{
		Message("Default message, can be overriden"),
	}
	options = append(options, opts...)
  return z.TestFunc(
    func (val any, ctx z.Ctx) bool {
      user := val.(*User)
      return user.Name == "test" // or any other validation
    },
   ...options
  )
}


func MyComplexTest() z.Test[*string] {
  return z.Test{
    Func: func (valPtr *string, ctx z.Ctx) {
      // complex test here
    }
  }
}

```
