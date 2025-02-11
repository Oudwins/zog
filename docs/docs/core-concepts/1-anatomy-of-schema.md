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
// takes the data as input and returns the new data which will then be passed onto the next functions.
// The function may return an error or a ZogIssue. In this case all validation will be skipped and the error will be wrapped into a ZogIssue and entire execution will return.
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

`schema.Required()` is a boolean that indicates if the field is required. If it is required and the data is a zero value the schema will return a [ZogIssue](/errors).

`schema.Default(value)` sets a default value for the field. If the data is a zero value it will be replaced with this value, this takes priority over required. Tests will still run on this value.

`schema.Catch(value)` sets a catch value. If this is set it will "catch" any errors or ZogIssues with the catch value. Meaning it will set the destination value to the catch value and exit. When this is triggered, no matter what error triggers it code will automatically jump to the [PostTransforms](#posttransforms). For more information checkout the [parsing execution structure](/core-concepts/parsing#parsing-execution-structure).

## Tests

> A test is what zod calls a "validator". It is a struct that represents an individual validation. For example for the String schema `z.String()` the method `Min(3)` generates a test that checks if the string is at least 3 characters long. You can view all the default tests that come with each [schema type here.](/zog-schemas)


### Test Options
You can configure tests with `TestOptions` which modify a test in some manner. Here are some examples:

```go
z.String().Min(3, z.Message("String must be at least 3 characters long")) // This sets the message that Zogissues will have if the validation fails
z.String().Min(3, z.IssueCode("min_3")) // This sets the issue code that Zogissues will have if the validation fails
z.String().Min(3, z.IssuePath("name")) // This sets the issue path that Zogissues will have if the validation fails
```

### Creating Custom Tests

You are also free to create custom tests and pass them to the `schema.Test()` and `schema.TestFunc()` methods. For more details on this checkout the [Creating Custom Tests](/custom-tests) page.


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
