---
sidebar_position: 200
hide_table_of_contents: false
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Core Design Decisions

- All fields optional by default. Same as graphql
- When parsing into structs, private fields are ignored (same as stdlib json.Unmarshal)
- Errors returned by you (for example in a `Preprocess` or `Transform` function) can be the ZogIssue interface or an error. If you return an error, it will be wrapped in a ZogIssue. ZogIssue is just a struct that wraps around an error and adds a message field which is text that can be shown to the user. For more on this see [Errors](/errors)
- You should not depend on test execution order. They might run in parallel in the future

> **A WORD OF CAUTION. [ZOG & PANICS](/docs/panics)**
> In general Zog will never panic if the input data is wrong but it will panic if you configure it wrong. For example:
>
> - In parse mode Zog will never panic due to invalid input data but will always panic if invalid destination is passed to the `Parse` function. if the destination does not match the schema in terms of types or fields.
> - In validate mode Zog will panic if the expected types or fields are not present in the structure you are validating.

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
	Email string
	Age   int
}
schema.Parse(map[string]any{"name": "zog"}, &User{}) // this will panic even if input data is valid. Because the destination is not a valid structure for the schema
schema.Validate(&User2{})                            // This will panic because the structure does not match the schema
```

## Limitations

Most of these things are issues we would like to address in future versions.

- Unsupported schemas:

  - `z.Map()`

- `zhttp` does not support parsing into any data type that is not a struct
- Schema & pick, omit, etc are not really typesafe. i.e `z.Struct(z.Schema{"name"})` name is not typesafe
- structs and slices don't support catch, and structs don't support default values
- It is not recommended to use very deeply nested schemas since that requires a lot of reflection and can have a negative impact on performance
