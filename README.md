<div align="center">
    <br />
    <a href="https://github.com/Oudwins/zog">
     <img src="https://raw.githubusercontent.com/Oudwins/zog/master/assets/logo-v1.jpg" alt="Zog, a Zod-like schema parser & validator" />
    </a>
</div>

# ZOG - A Zod & Yup like Schema Parser & Validator for GO

<a href="https://pkg.go.dev/github.com/Oudwins/zog"><img src="https://pkg.go.dev/badge/github.com//github.com/Oudwins/tailwind-merge-go.svg" alt="Go Reference" /></a>
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/Oudwins/zog)](https://goreportcard.com/report/github.com/Oudwins/zog) [![Coverage Status](https://coveralls.io/repos/github/Oudwins/zog/badge.svg?branch=master)](https://coveralls.io/github/Oudwins/zog?branch=master)

Zog is a schema builder for runtime value parsing and validation. Define a schema, transform a value to match, assert the shape of an existing value, or both. Zog schemas are extremely expressive and allow modeling complex, interdependent validations, or value transformation.

Killer Features:

- Concise yet expressive schema interface, equipped to model simple to complex data models
- Support for parsing and validating Query Params and Form inputs directing from the request object with `schema.Parse(zhttp.NewRequestDataProvider(r), &dest)`
- Extensible: add your own validators, schemas and data providers
- Rich error details, make debugging a breeze
- Almost no reflection when using primitive types
- Built-in coercion support for most types
- Zero dependencies!

API Stability:

- I will consider the API stable when we reach v1.0.0
- However, I believe very little API changes will happen from the current implementation. The API that is most likely to change is the everything related to customizing error messages. However, if you use z.Message() that API will most likely not change and you won't be affected.

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

var nameSchema = z.Struct(z.Schema{
  "name": z.String().Min(3, z.Message("Override default message")).Max(10),
})

var ageSchema = z.Struct(z.Schema{
  "age": z.Int().GT(18).Required(z.Message("is required")),
})

// Merge the schemas creating a new schema
var schema = nameSchema.Merge(ageSchema)

type User struct {
  Name string `zog:"firstname"` // tag is optional. If not set zog will check for "name" field in the input data
  Age int
}
```

**3 Parse the struct**

```go
func main() {
  u := User{}
  m := map[string]string{
    "firstname": "", // won't return an error because fields are optional by default
    "age": "30", // will get casted to int
  }
  errsMap := schema.Parse(z.NewMapDataProvider(m), &u)
  if errsMap != nil {
    // handle errors -> see Errors section
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

- The struct validator always expects a `DataProvider`, which is an interface that wraps around an input like a map. This is less efficient than doing it directly but allows us to reuse the same code for all kinds of data sources (i.e json, query params, forms, etc).
- All fields optinal by default. Same as graphql
- Errors returned by you can be an instance of ZogError or an error. If you return an error, it will be wrapped in a ZogError using the error.Error() value as the message. ZogError is just a struct that wraps around an error and adds a message field which is is text that can be shown to the user.
- You should not depend on test execution order. They might run in parallel in the future
- When parsing into structs, private fields are ignored (same as stdlib json.Unmarshal)

> **A WORD OF CAUTION. ZOG & PANICS**
> Zog will never panic due to invalid input but will always panic if invalid destination is passed to the `Parse` function (i.e if the destination does not match the schema).

```go
var schema = z.Struct(z.Schema{
  "name": z.String().Required(),
})
// This struct is a valid destionation for the schema
type User struct {
  Name string `zog:"firstname"` // name will be parsed from the firstname field of the input data (i.e form, json, query params)
  Age int // age will be ignored since it is not a field in the schema
}
// this struct is not a valid destionation for the schema. It is missing the name field
type User2 struct {
  Email string,
  Age int
}

```

**Changes from zog**:

- The refine method for providing a custom validation function is renamed to `schema.Test()`
- schemas are optional by default (in zod they are required)
- The `z.Enum()` type from zod is removed in favor of `z.String().OneOf()` and is only supported for strings and numbers

## Limitations

Most of these things are issues we would like to address in future versions.

- Structs do not support pointers at the moment
- slices do not support pointers
- maps are not a supported schema type
- structs and slices don't support catch, and structs don't suppoort default values
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

// Init our typesafe env vars, panic if any envs are missing
func Init() {
  errs := envSchema.Parse(zenv.NewDataProvider(), &Env)
  if errs != nil {
    log.Fatal(errs)
  }
}
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

## Errors

> **WARNING**: The errors API is probably what is most likely to change in the future. I will try to keep it backwards compatible but I can't guarantee it.

Zog creates its own error type called `ZogError` that implements the error interface.

```go
type ZogError struct {
  Message string
  Err     error
}
```

This is what will be returned by the `Parse` function. To be precise:

- Primitive types will return a list of `ZogError` instances.
- Complex types will return a map of `ZogError` instances. Which uses the field path as the key & the list of errors as the value.

For example:

```go
errList := z.String().Min(5).Parse("foo", &dest) // can return []z.ZogError{z.ZogError{Message: "min length is 5"}} or nil
errMap := z.Struct(z.Schema{"name": z.String().Min(5)}).Parse(data, &dest) // can return map[string][]z.ZogError{"name": []z.ZogError{{Message: "min length is 5"}}} or nil

// Slice of 2 strings with min length of 5
errsMap2 := z.Slice(z.String().Min(5)).Len(2).Parse(data, &dest) // can return map[string][]z.ZogError{"$root": []z.ZogError{{Message: "slice length is not 2"}, "[0]": []z.ZogError{{Message: "min length is 5"}}}} or nil
```

Additionally, `z.ZogErrMap` will use the field path as the key. Meaning

```go
errsMap := z.Struct(z.Schema{"inner": z.Struct(z.Schema{"name": z.String().Min(5)}), "slice": z.Slice(z.String().Min(5))}).Parse(data, &dest)
errsMap["inner.name"] // will return []z.ZogError{{Message: "min length is 5"}}
errsMap["slice[0]"] // will return []z.ZogError{{Message: "min length is 5"}}
```

`$root` & `$first` are reserved keys for both Struct & Slice validation, they are used for root level errors and for the first error found in a schema, for xample:

```go
errsMap := z.Slice(z.String()).Min(2).Parse(data, &dest)
errsMap["$root"] // will return []z.ZogError{{Message: "slice length is not 2"}}
errsMap["$first"] // will return the same in this case []z.ZogError{{Message: "slice length is not 2"}}
```

### Example ways of delivering errors to users

#### Using go templ templates

**Example use case: simplified Signup form validation**
Imagine our handler looks like this:

```go
type SignupFormData struct {
  Email string
  Password string
}

schema := z.Struct(z.Schema{
  "email": z.String().Email().Required(),
  "password": z.String().Min(8).Required(),
})

func handleSignup(w http.ResponseWriter, r *http.Request) {
  var signupFormData = SignupFormData{}
  errs := schema.Parse(zhttp.NewRequestDataProvider(r), &signupFormData)

  if errs != nil {
    www.Render(signupFormTempl(&signupFormData, errs))
  }
  // handle successful signup
}

templ signupFormTempl(data *SignupFormData, errs z.ZogErrMap) {
  <input type="text" name="email" value={data.Email}>
  // display only the first error
  if e, ok := errs["email"]; ok {
    <p class="error">{e[0].Message}</p>
  }
  <input type="text" name="password" value={data.Password}>
  // display only the first error
  if e, ok := errs["password"]; ok {
    <p class="error">{e[0].Message}</p>
  }
}
```

#### REST API Responses

Zog providers a helper function called `z.Errors.SanitizeMap(errsMap)` that will return a map of strings of the error messages (stripping out the internal error). So, if you do not mind sending errors to your users in the same form zog returns them, you can do something like this:

```go
errs := schema.Parse(data, &userFormData)

if errs != nil {
  sanitized := z.Errors.SanitizeMap(errs)
  // sanitize will be map[string][]string
  // for example:
  // {"name": []string{"min length is 5", "max length is 10"}, "email": []string{"is not a valid email"}}

  // ... marshal sanitized to json and send to the user
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

## Overriding Defaults

Zog uses internal functions to handle many aspects of validation & parsing. We aim to provide a simple way for you to customize the default behaviour of Zog through simple declarative code inside your project. You can find the options you can tweak & override in the conf package (`github.com/Oudwins/zog/conf`).

Currently the only default behaviour that can be overridden are the coerce functions, in the future we will add more.

### Overriding the default coercer functions

Lets go through an example of overriding the `float64` coercer function, because we want to support floats that use a comma as the decimal separator.

```go
import (
  // import the conf package
	"github.com/Oudwins/zog/conf"
)

// we save the original to use later
var zogFloat64Coercer =  conf.Coercers["float64"];

// we override the coercer function for float64
conf.Coercers["float64"] = func(data any) (any, error) {
  str, ok := data.(string)
  // identify the case we want to override
  if !ok && strings.Contains(str, ",") {
    return MyCustomFloatCoercer(str)
  }
  // fallback to the original function
  return zogFloat64Coercer(data)
}
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
- Support for schema.Clone()
- Better support for custom error messages (including failed coercion error messages) & i18n
- support for catch & default for structs & slices
- implement errors.SanitizeMap/Slice -> Will leave only the safe error messages. No internal stuff. Optionally this could be a parsing option in the style of `schema.Parse(m, &dest, z.WithSanitizeErrors())`
- Add additional tests
- Better docs

## Acknowledgments

- Big thank you to @AlexanderArvidsson for being there to talk about architecture and design decisions. It helped a lot to have someone to bounce ideas off of
- Credit for all the inspiration goes to /colinhacks/zod & /jquense/yup
- Credit for the initial idea goes to anthony -> /anthdm/superkit he made a hacky version of this idea that I used as a starting point, I was never happy with it so I inspired me to rewrite it from scratch
- Credit for the logo goes to /colinhacks/zod

## License

This project is licensed under the MIT License -
see the [LICENSE](LICENSE) file for details.
