---
sidebar_position: 300
---

# FAQ

## My code panics / zog doesn't work

Please make sure that you are using Zog correctly. The expectation is that if Zog causes a panic you are doing something wrong. So before reporting an issue please check this common mistakes:

**1 Not passing a pointer to the struct you are parsing into**

```go
payload := models.UserPayload{}
errMap := models.UserSchema.Parse(zhttp.Request(c.Request), payload) // note that the "payload" is not a pointer correct code here is &payload
```

**2 No defining your schema correctly**

```go
type Name struct {
	FirstName string `json:"first_name" zog:"first_name"` // zog struct tag is used to define the name of the field in the input data (i.e json key) not the name of the schema key (common mistake)
	LastName  string `json:"last_name" zog:"last_name"`
}

data := new(Name)

var schema = z.Struct(z.Schema{
	"first_name": z.String().Required(z.Message("First name is required")), // here you are telling zog that your struct should have a First_name field, but this is incorrect because the struct has a FirstName field. The key here should be "firstName" or "FirstName" (both are valid)
	"last_name":  z.String().Required(z.Message("Last name is required")),  // same issue here
})
```

This comes from a misunderstanding of how the zog struct tag works. The zog struct tag is used to define the name of the field in the input data (i.e json key) not the name of the schema key. The schema is only aware of what the struct looks like it is not aware of the source of the input data.

## Error: "Struct is missing expected schema key: \{some_key\}"

This error is telling you that you are defining your schema keys incorrectly. For example:

```go
type Name struct {
	FirstName string `json:"first_name" zog:"first_name"` // zog struct tag is used to define the name of the field in the input data (i.e json key) not the name of the schema key (common mistake)
	LastName  string `json:"last_name" zog:"last_name"`
}

data := new(Name)

var schema = z.Struct(z.Schema{
	"first_name": z.String().Required(z.Message("First name is required")), // here you are telling zog that your struct should have a First_name field, but this is incorrect because the struct has a FirstName field. The key here should be "firstName" or "FirstName" (both are valid)
	"last_name":  z.String().Required(z.Message("Last name is required")),  // same issue here
})
```

This comes from a misunderstanding of how the zog struct tag works. The zog struct tag is used to define the name of the field in the input data (i.e json key) not the name of the schema key. The schema is only aware of what the struct looks like it is not aware of the source of the input data.

## Why does zog have an internals package?

Please see the [internals](/packages/internals) page for more information.
