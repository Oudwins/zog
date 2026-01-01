---
sidebar_position: 300
---

# FAQ

## My code panics / zog doesn't work

Please make sure that you are using Zog properly. The expectation is that if Zog causes a panic, you are doing something wrong.

So before reporting an issue please check these common mistakes:

### 1. ❌ Not passing a pointer to the struct you are parsing into

```go
payload := models.UserPayload{}
errMap := models.UserSchema.Parse(zhttp.Request(c.Request), payload) // ❌ Incorrect: note how the "payload" is not being passed as a pointer, you must pass `&payload` as the pointer to `payload`
```

#### ✅ Correct usage

```go
payload := models.UserPayload{}
errMap := models.UserSchema.Parse(zhttp.Request(c.Request), &payload) // ✅ Correct: passing the pointer to `payload`
```

### 2. ❌ Not defining your schema properly

```go
type Name struct {
	FirstName string `json:"first_name" zog:"first_name"` // The zog/json tags define the INPUT (e.g. JSON) key, not the schema key. The schema only matches Go struct field names.
	LastName  string `json:"last_name" zog:"last_name"`
}

data := new(Name)

var schema = z.Struct(z.Shape{
	"first_name": z.String().Required(z.Message("First name is required")), // ❌ Incorrect: schema keys must match the Go struct field name (FirstName), not the JSON/input key
	"last_name":  z.String().Required(z.Message("Last name is required")),  // ❌ Same issue: this refers to the JSON key, but the struct field is LastName
})
```

The schema operates purely on the _Go struct shape_. It does not know or care about the source of the input data (JSON, form data, etc.).  
Struct tags (`json`, `zog`) are _only used to map_ input keys into struct fields.  
Schema keys must always correspond to Go struct field names (e.g. `FirstName`, `LastName`), not input keys like `first_name`.

#### ✅ Correct usage

```go
type Name struct {
    FirstName string `json:"first_name" zog:"first_name"`
    LastName  string `json:"last_name" zog:"last_name"`
}
data := new(Name)
var schema = z.Struct(z.Shape{
    "FirstName": z.String().Required(z.Message("First name is required")), // ✅ Correct: schema key matches Go struct field name
    "LastName":  z.String().Required(z.Message("Last name is required")),  // ✅ Correct: schema key matches Go struct field name
})
```

## Error: `"Struct is missing expected schema key: {some_key}"`

This error is telling you that you are defining your schema keys incorrectly. For example:

```go
type Name struct {
	FirstName string `json:"first_name" zog:"first_name"` // zog struct tag is used to define the name of the field in the input data (i.e json key) not the name of the schema key (common mistake)
	LastName  string `json:"last_name" zog:"last_name"`
}

data := new(Name)

var schema = z.Struct(z.Shape{
	"first_name": z.String().Required(z.Message("First name is required")), // ❌ Incorrect: here you are telling zog that your struct should have a `first_name` field, but this is incorrect because the struct has a `FirstName` field. The key here should be "firstName" or "FirstName" (both are valid)
	"last_name":  z.String().Required(z.Message("Last name is required")),  // ❌ Same issue here: the struct has a `LastName` field, so the key here should be "lastName" or "LastName"
})
```

This comes from a misunderstanding of how the zog struct tag works.</br>The zog struct tag is used to _define the name of the field in the input data_ (i.e json key), **not** the name of the schema key. The schema is only aware of what the struct looks like, and **it is not** aware of the source of the input data.

## Why does zog have an internals package?

Please see the [internals](/packages/internals) page for more information.
