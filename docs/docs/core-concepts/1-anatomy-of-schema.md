---
sidebar_position: 1
---

# Anatomy of the Zog Schema

A zog schema is an interface implemented by multiple custom structs that represent a set of `validation` and `transformation` logic for a variable of a given type. For example:

```go
stringSchema := z.String().Trim().Min(3).Required()    // A zog schema that represents a required string which will first be trimmed then a test to ensure it has 3+ characters will be ran.
userSchema := z.Struct(z.Schema{"name": stringSchema}) // a zog schema that represents a user struct. Also yes I know that z.Schema might be confusing but think of it as the schema for the struct not a ZogSchema
```

**The string schema, for example, looks something like this:**

```go
type stringSchema struct {
	isRequired     bool            // optional. Defaults to FALSE
	defaultValue   string          // optional. if the input value is a "zero value" it will be replaced with this. Tests will still run on this value.
	catchValue     string          // optional. If this is set it will "catch" any errors, set the destination value to this value and exit
	testOrTransformation TestOrTransformation // This is the test or transformation that will be applied to the data.
}
```

## Required, Default and Catch

`schema.Required()` is a boolean that indicates if the field is required. If it is required and the data is a zero value the schema will return a [ZogIssue](/errors).

`schema.Default(value)` sets a default value for the field. If the data is a zero value it will be replaced with this value, this takes priority over required. Tests will still run on this value.

`schema.Catch(value)` sets a catch value. If this is set it will "catch" any errors or ZogIssues with the catch value. Meaning it will set the destination value to the catch value and exit. When this is triggered, no matter what error triggers it code will automatically exit. For more information checkout the [parsing execution structure](/core-concepts/parsing#parsing-execution-structure).

## Tests

> A test is what zog calls a "validator". It is a struct that represents an individual validation. For example for the String schema `z.String()` the method `Min(3)` generates a test that checks if the string is at least 3 characters long. You can view all the default tests that come with each [schema type here.](/reference)

### Test Options

You can configure tests with `TestOptions` which modify a test in some manner. Here are some examples:

```go
z.String().Min(3, z.Message("String must be at least 3 characters long")) // This sets the message that Zogissues will have if the validation fails
z.String().Min(3, z.IssueCode("min_3"))                                   // This sets the issue code that Zogissues will have if the validation fails
z.String().Min(3, z.IssuePath("name"))                                    // This sets the issue path that Zogissues will have if the validation fails
```

### Creating Custom Tests

You are also free to create custom tests and pass them to the `schema.Test()` and `schema.TestFunc()` methods. For more details on this checkout the [Creating Custom Tests](/custom-tests) page.

## Transforms

Transforms is a list of function that are applied to the data at any point. You can think of it like a `pipeline` of transformations for a specific schema. This is the function signature:

```go
// Transforms are generic functions that take a pointer to the data as input. For primitive types you won't have to typecast but for complex types it will just be a any type and you will have to manually typecast it.
type Transform[T any] func(dataPtr T, ctx Ctx) error
```

As you can see the function takes a pointer to the data as input. This is to allow the function to modify the data.

```go
type User struct {
	Phone    string
	AreaCode string
}

z.Struct(z.Schema{
	"phone": z.String().Test(...).Transform(func (valPtr *string, ctx z.Ctx) error{
		*valPtr = strings.ReplaceAll(*valPtr, " ", "") // remove all spaces
		return nil
	}),
}).Transform(func(dataPtr any, ctx z.Ctx) error {
	user := dataPtr.(*User)
	user.AreaCode = user.Phone[:3]
	user.Phone = user.Phone[3:]
	return nil
})
```
