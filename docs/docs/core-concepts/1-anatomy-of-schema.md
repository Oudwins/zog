---
sidebar_position: 1
---

# Anatomy of the Zog Schema

A zog schema is an interface implemented by multiple custom structs that represent a set of `validation` and `transformation` logic for a variable of a given type. For example:

```go
stringSchema := z.String().Min(3).Required().Trim() // A zod schema that represents a required string string of minimum 3 characters and will be trimmed for white space
userSchema := z.Struct(z.Schema{"name": stringSchema}).Required() // a zod schema that represents a user struct. Also yes I know that z.Schema might be confusing but think of it as the schema for the struct not a ZogSchema
```

**The string schema, for example, looks something like this:**

```go
type stringSchema struct {
   preTransforms: []PreTransforms // transformations executed before the validation. For example trimming the string
   isRequired: bool // optional. Defaults to FALSE
   defaultValue: string // optional. if the input value is a "zero value" it will be replaced with this. Tests will still run on this value.
   catchValue: string // optional. If this is set it will "catch" any errors, set the destination value to this value and exit
   [tests](#tests): []Test // These are your validation checks. Such as .Min(), .Contains(), etc
   postTransforms: []PostTransform // transformations executed after the validation.
}
```

## PreTransforms

Pretransforms is a list of function that are applied to the data before the [tests](#tests) are run. You can think of it like a `pipeline` if transformations for a specific schema. **PreTransforms are PURE functions**. They take in data and return new data. This is the function signature:

```go
// takes the data as input and returns the new data which will then be passed onto the next functions. If the function returns an error all validation will be skipped & the error will be returned
type PreTransform = func(data any, ctx Ctx) (out any, err error)
```

You can use pretransforms for things like trimming whitespace, splitting strings, etc. Here is an example of splitting a string into a slice of strings:

```go
z.Slice(z.String()).PreTransform(func(data any, ctx Ctx) (any, error) {
  return strings.Split(data.(string), ","), nil
}).Parse("item1,item2,item3", &dest)
```

> **FOOTGUNS**
> *Type Coercion*: Please note that pretransforms are executed before type coercion (if using `schema.Parse()`). This means that if you are responsible for checking that the data matches your expected type. If you blinding typecast the data to, for example, an int and your user provides a string as input you will cause a panic.
> *Pure Functions*: Since pretransforms are pure functions. A copy of the data is passed in. So if you place a preTransform on a large struct it will copy the entire struct which is not efficient. Be careful!

## Required, Default and Catch

`schema.Required()` is a boolean that indicates if the field is required. If it is required and the data is a zero value the schema will return an error.

`schema.Default(value)` sets a default value for the field. If the data is a zero value it will be replaced with this value, this takes priority over required. Tests will still run on this value.

`schema.Catch(value)` sets a catch value. If this is set it will "catch" any errors with the catch value. Meaning it will set the destination value to the catch value and exit. When this is triggered, no matter what error triggers it code will automatically jump to the [PostTransforms](#posttransforms). For more information checkout the [parsing execution structure](/core-concepts/parsing#parsing-execution-structure).

## Tests

A test is what zod calls a "validator". It is a struct that represents an individual validation. For example `z.String().Min(3)` is a test that checks if the string is at least 3 characters long.

A test is a struct that looks something like this:

```go
type Test struct {
	ErrCode      zconst.ZogErrCode // the error code to use if the validation fails. This helps identify the type of error, for example ErrCodeMin identifies the Min() test
	ValidateFunc TestFunc // a function that takes the data as input and returns a boolean indicating if it is valid or not
}
type TestFunc = func(val any, ctx Ctx) bool
```

You can view all the default tests that come with each [schema type here.](/zog-schemas)

##### Creating Custom Tests

There are two ways to create custom tests:

```go
// 1. Using the `z.TestFunc()` function:
z.String().Test(z.TestFunc("my_custom_error_code", func(data any, ctx z.Ctx) bool {
  return data == "test"
}))
// 2. Using the `z.Test` struct directly:
z.String().Test(z.Test{
  ErrCode: "my_custom_error_code",
  ValidateFunc: func(data any, ctx z.Ctx) bool {
    return data == "test"
  },
})
```

## PostTransforms

PostTransforms is a list of function that are applied to the data after the [tests](#tests) are run. You can think of it like a `pipeline` if transformations for a specific schema. This is the function signature:

```go
// type for functions called after validation & parsing is done
type PostTransform = func(dataPtr any, ctx Ctx) error
```

As you can see the function takes a pointer to the data as input. This is to allow the function to modify the data.

You can use posttransforms for any transformation you want to do but don't want it to affect the validation process. For example imagine you want to validate a phone number and afterwards separate the string into the area code and the rest of the number.

```go
type User struct {
  Phone string
  AreaCode string
}
z.Struct(z.Schema{
  "phone": z.String().Test(....),
  }).PostTransform(func(dataPtr any, ctx z.Ctx) error {
  user := dataPtr.(*User)
  user.AreaCode = user.Phone[:3]
  user.Phone = user.Phone[3:]
  return nil
})
```
