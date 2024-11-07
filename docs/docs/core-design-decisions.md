---
sidebar_position: 9
hide_table_of_contents: false
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Core Design Decisions

- All fields optinal by default. Same as graphql
- When parsing into structs, private fields are ignored (same as stdlib json.Unmarshal)
- The struct parser expects a `DataProvider` (although if you pass something else to the data field it will try to coerce it into a `DataProvider`), which is an interface that wraps around an input like a map. This is less efficient than doing it directly but allows us to reuse the same code for all kinds of data sources (i.e json, query params, forms, etc).
- Errors returned by you can be the ZogError interface or an error. If you return an error, it will be wrapped in a ZogError. ZogError is just a struct that wraps around an error and adds a message field which is is text that can be shown to the user.
- You should not depend on test execution order. They might run in parallel in the future

> **A WORD OF CAUTION. ZOG & PANICS**
> Zog will never panic due to invalid input but will always panic if invalid destination is passed to the `Parse` function (i.e if the destination does not match the schema).

```go
var schema = z.Struct(z.Schema{
  "name": z.String().Required(),
})
// This struct is a valid destionation for the schema
type User struct {
  Name string
  Age int // age will be ignored since it is not a field in the schema
}
// this struct is not a valid destination for the schema. It is missing the name field. This will cause a panic even if the input data is map[string]any{"name": "zog"}
type User2 struct {
  Email string,
  Age int
}

```

## Limitations

Most of these things are issues we would like to address in future versions.

- Structs do not support pointers at the moment
- slices do not support pointers
- maps are not a supported schema type
- structs and slices don't support catch, and structs don't suppoort default values
- Validations and parsing cannot be run separately
- It is not recommended to use very deeply nested schemas since that requires a lot of reflection and can have a negative impact on performance
