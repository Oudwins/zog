![logo](assets/logo-v1.jpg)

# ZOG - A Zod & Yup like Schema Parser & Validator for GO

Zog is a schema builder for runtime value parsing and validation. Define a schema, transform a value to match, assert the shape of an existing value, or both. Zog schemas are extremely expressive and allow modeling complex, interdependent validations, or value transformation.

Killer Features:

- Concise yet expressive schema interface, equipped to model simple to complex data models
- Support for parsing and validating Query Params and Form inputs directing from the request object with `z.Params()` & `z.Form()`
- Extensible: add your own validators and schemas
- Rich error details, make debugging a breeze
- Almost no reflection when using primitive types

## Introduction

**1 Install**

```bash
go get github.com/Oudwins/zog
```

**2 Create a schema & a struct**

```go
import (
  z "github.com/Oudwins/zog"
   )

var schema = z.Struct(z.Schema{
  "name": z.String().Min(3, z.Message("Override default message")).Max(10),
  "age": z.Int().GT(18).Required(z.Message("is required")),
})

type User struct {
  Name string `zog:"name"` // optional zog will use field name by default
  Age int
}
```

**3 Parse the struct**

```go
func main() {
  u := User{}
  m := map[string]string{
    "name": "", // won't return an error because fields are optional by default
    "age": "30", // will get casted to int
  }
  errsMap := schema.Parse(z.NewMapDataProvider(m), &u)
  if errsMap != nil {
    // handle errors
  }
  u.Name // ""
  u.Age // 30
}
```

**4 You can also parse individual fields**

```go
var t = time.Time
errsList := Time().Required().Parse("2020-01-01T00:00:00Z", &t)
```

**5 And do stuff before and after parsing**

```go
var dest []string
Slice(String().Email().Required()).PreTransform(func(val any, ctx *z.ParseCtx) (any, error) {
  s := val.(string)
  return strings.Split(s, ","), nil
}).PostTransform(func(val any, ctx *z.ParseCtx) error {
  s := val.(*[]string)
  for i, v := range s {
    s[i] = strings.TrimSpace(v)
  }
  return nil
}).Parse("foo@bar.com,bar@foo.com", &dest) // dest = [foo@bar.com bar@foo.com]
```

## Core Design Decisions

> **A WORD OF CAUTION. ZOG & PANICS**
> Zog will never panic due to invalid input but will always panic if invalid destination is passed to the `Parse` function (i.e if there are discrepancies between the schema and the destination).

- The struct validator always expects a `DataProvider`, which is an interface that wraps around an input like a map. This is less efficient than doing it directly but allows us to reuse the same code for all kinds of data sources (i.e json, query params, forms, etc).
- All fields optinal by default. Same as graphql
- Errors returned by you can be an instance of ZogError or an error. If you return an error, it will be wrapped in a ZogError using the error.Error() value as the message. ZogError is just a struct that wraps around an error and adds a message field which is is text that can be shown to the user.
- You should not depend on test execution order. They might run in parallel in the future
- When parsing into structs, private fields are ignored (same as stdlib json.Unmarshal)

**Changes from zog**:

- The refine method for providing a custom validation function is renamed to `schema.Test()`
- schemas are optional by default (in zod they are required)
- The `z.Enum()` type from zod is removed in favor of `z.String().OneOf()` and is only supported for strings and numbers

## Limitations

Most of these things are issues we would like to address in future versions.

- Structs do not support pointers at the moment
- slices do not support pointers or structs
- maps are not a supported schema type
- structs and slices don't support catch or default values
- You can provide custom error messages, but cannot customize coercion error messages or set global defaults
- Validations and parsing cannot be run separately
- It is not recommended to use very deeply nested schemas since that requires a lot of reflection and can have a negative impact on performance

## Zenv & ZHTTP

For convenience zog provides two helper packages:

**zenv: helps validate environment variables**

```go
import (
  z "github.com/Oudwins/zog"
  "github.com/Oudwins/zog/zenv"
)

var envSchema = z.Struct(z.Schema{
	"PORT": z.Int().GT(1000).LT(65535).Default(3000),
	"Db": z.Struct(z.Schema{
		"Host": z.String().Default("localhost"),
		"User": z.String().Default("root"),
		"Pass": z.String().Default("root"),
	}),
})
var Env = struct {
	PORT int // zog will automatically coerce the PORT env to an int
	Db   struct {
		Host string `zog:"DB_HOST"` // we specify the zog tag to tell zog to parse the field from the DB_HOST environment variable
		User string `zog:"DB_USER"`
		Pass string `zog:"DB_PASS"`
	}
}{}
// Parse will log.Fatal if there are errors because we have set panicOnError to true otherwise it will return the errors
var _ = zenv.Parse(envSchema, &Env, true)
```

**zhttp: helps parse http requests forms & query params**

```go
import (
  z "github.com/Oudwins/zog"
  "github.com/Oudwins/zog/zhttp"
)
var userSchema = z.Struct(z.Schema{
  "name": z.String().Required(),
  "age": z.Int().Required().GT(18),
})

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
  var user := struct {
    Name string
    Age int
  }
  // Parse the request into the user struct from the query params or the form data
  errs := userSchema.Parse(zhttp.NewRequestDataProvider(r), &user)
  if errs != nil {
  }
  user.Name // defined
  user.Age // defined
}

```

## Reference

### Generic Schema Methods

These are methods that can be used on most types of schemas

```go
// gets passed the destionation valiue and the context and returns a boolean. Please note for complex types you will be passed a pointer to the destination value
schema.Test("rule name", z.Message("message or function"), func(val any, ctx *z.ParseCtx) bool {})

// marks the schema as required. Remember fields are optional by default
schema.Required(z.Message("message or function"))
schema.Optional() // marks the schema as optional
// optional & required are mutually exclusive
schema.Required().Optional() // marks the schema as optional

schema.Default(val) // sets the default value. See Zog execution flow
schema.Catch(val) // sets the catch value. A value to use if the validation fails. See Zog execution flow

schema.PreTransform(func(val any, ctx *z.ParseCtx) (any, error) {}) // transforms the value before validation. returned value will override the input value. See Zog execution flow

schema.PostTransform(func(destPtr any, ctx *z.ParseCtx) error {}) // transforms the value after validation. Receives a pointer to the destination value.
```

### Types

```go
// Primtives. Calling .Parse() on these will return []ZogError
String()
Int()
Float()
Bool()
Time()

// Complex Types. Calling .Parse() on these will return map[string][]ZogError. Where the key is the field path ("user.email") & $root is the list of complex type level errors not the specific field errors
Struct(Schema{
  "name": String(),
})
Slice(String())
```

#### Strings

```go
// Validations
String().Min(5)
String().Max(10)
String().Len(5)
String().Email()
String().URL()
String().Contains(string)
String().ContainsUpper()
String().ContainsDigit()
String().ContainsSpecial()
String().HasPrefix(string)
String().HasSuffix(string)
String().OneOf([]string{"a", "b", "c"})
```

#### Numbers

```go
// Validators
Int().GT(5)
Float().GTE(5)
Int().LT(5)
Float().LTE(5)
Int().EQ(5)
Float().OneOf([]float64{1.0, 2.0, 3.0})
```

#### Booleans

```go
Bool().True()
Bool().False()
```

#### Times & Dates

Use Time to validate `time.Time` instances

```go
Time().After(time.Now())
Time().Before(time.Now())
Time().Is(time.Now())
```

#### Structs

```go
s := z.Struct(z.Schema{
  "name": String().Required(),
  "age": Int().Required(),
})
user := struct {
  Name string `zog:"firstname"` // name will be parsed from the firstname field
  Age int // since zog tag is not set, age will be parsed from the age field
}
s.Parse(NewMapDataProvider(map[string]any{"firstname": "hello", "age": 10}), &user)
```

#### Slices

```go
s := Slice(String())

Slice(Int()).Min(5)
Slice(Float()).Max(5)
Slice(Bool()).Length(5)
Slice(String()).Contains("foo")
```

## Zog Schema Parsign Execution Structure

![Zog Schema Parsign Execution Structure](/assets/parsing-workflow.png)

1. Pretransforms
   - On error all parsing and validation stops and error is returned.
   - Can be caught by catch
2. Default Check -> Assigns default value if the value is nil value
3. Optional Check -> Stops validation if the value is nil value
4. Casting -> Attempts to cast the value to the correct type
   - On error all parsing and validation stops and error is returned
   - Can be caught by catch
5. Required check ->
   - On error: aborts if the value is its nil value and returns required error.
   - Can be caught by catch
6. Tests -> Run all tests on the value (including required)
   - On error: validation errors are added to the errors. All validation functions are run even if one of them fails.
   - Can be caught by catch
7. PostTransforms -> Run all postTransforms on the value.
   - On error you return: aborts and adds your error to the list of errors
   - Only run on valid values. Won't run if an error was created before the postTransforms

## Roadmap

These are the things I want to add to zog before v1.0.0

- For structs & slices: support pointers
- Support for schema.Merge(schema2) && schema.Clone()
- Better support for custom error messages (including failed coercion error messages) & i18n
- support for catch & default for structs & slices
- implement errors.SanitizeMap/Slice -> Will leave only the safe error messages. No internal stuff. Optionally this could be a parsing option in the style of `schema.Parse(m, &dest, z.WithSanitizeErrors())`
- Add additional tests
