---
sidebar_position: 199
hide_table_of_contents: false
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Zog Panics

## When does Zog panic?

Zog follows [TigerStyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md) asserts. It panics when something in its fundamental assumptions is broken.

In practice this means that Zog will never panic if the input data is wrong but it will panic if you configure it wrong. Most of the time "configured it wrong" means that you have made a mistake in your schema definition which puts Zog into an invalid state and results in a schema that can never succeed.

## Types of Panics

> If you find a panic that is not listed here, please report it as a bug!

### Schema Definition Errors

```go
var schema = z.Struct(z.Schema{
	"name": z.String().Required(),
})

// This struct is a valid destination for the schema
type User struct {
	Name string
	Age  int // age will be ignored since it is not a field in the schema
}

// this struct is not a valid structure for the schema. It is missing the name field.
// This will cause Zog to panic in both Parse and Validate mode
type User2 struct {
	Email string `zog:"name"` // using struct tag here DOES NOT WORK. This is not the purpose of the struct tag.
	Age   int
}
schema.Parse(map[string]any{"name": "zog"}, &User{}) // this will panic even if input data is valid. Because the destination is not a valid structure for the schema
schema.Validate(&User2{})                            // This will panic because the structure does not match the schema
```

### Type Cast Errors

There are multiple ways in which a type cast error can occur. For example:

###### 1 Destination/Validation value is not a pointer

```go
var schema = z.Struct(z.Schema{
	"name": z.String().Required(),
})

var dest User
schema.Parse(map[string]any{"name": "zog"}, dest) // This will panic because dest is not a pointer
schema.Validate(dest)                             // This will panic because dest is not a pointer
// Fix this by using a pointer
schema.Parse(map[string]any{"name": "zog"}, &dest)
schema.Validate(&dest)
```

> This can only really happen on complex schemas since those are not fully typesafe. Primitive schemas are typesafe and won't let you pass a non-pointer value.

###### 2 Destination/Validation value is not a valid type for the schema

```go
type MyString string
type User struct {
	Age MyString
}
var schema = z.Struct(z.Schema{
	"age": z.String().Required(),
})

val := User{
	Age: MyString("1"),
}
schema.Validate(&val) // This will panic because the schema is expecting a string but the value is of type MyString
```

Same thing will happen if you incorrectly set the type in a z.Custom schema:

```go

type User struct {
	ID uuid.UUID
}

var schema = z.Struct(z.Schema{
	"id": z.Custom(func (ptr *string, ctx z.Ctx) bool { // Zog can't convert a UUID to a string so this will panic
		return true
	}),
})

val := User{
	ID: uuid.New(),
}

schema.Validate(&val) // This will panic because the schema is expecting a string but the value is of type uuid.UUID
```

Another common example is when you forget to use z.Ptr.

```go
type User struct {
	Friends *[]Friend
}

// This is incorrect!
var schema = z.Struct(z.Schema{
	"friends": z.Slice(z.Struct(z.Schema{
		"name": z.String().Required(),
	})),
})

// This is correct!
var schema2 = z.Struct(z.Schema{
	"friends": z.Ptr(z.Slice(z.Struct(z.Schema{
		"name": z.String().Required(),
	}))),
})
```

###### 3 The coercer returns a value of the wrong type

> Only applicable to `schema.Parse()`

```go
var schema = z.Struct(z.Schema{
	"name": z.String(z.WithCoercer(func (v any, ctx z.Ctx) (any, error) {
		return 1, nil // we are returning an int but the schema is expecting a string
	})).Required(),
})

val := User{
	Name: "zog",
}
schema.Parse(map[string]any{"name": "zog"}, &val) // This will panic because the coercer is returning an int but the schema is expecting a string
```
