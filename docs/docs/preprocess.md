## PreTransforms

Pretransforms is a list of function that are applied to the data before the [tests](#tests) are run. You can think of it like a `pipeline` of pre validation transformations for a specific schema. These are similar to preprocess functions in zod. **PreTransforms are PURE functions**. They take in data and return new data. This is the function signature:

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

> **FOOTGUNS** > _Type Coercion_: Please note that pretransforms are executed before type coercion (if using `schema.Parse()`). This means that if you are responsible for checking that the data matches your expected type. If you blinding typecast the data to, for example, an int and your user provides a string as input you will cause a panic.
> _Pure Functions_: Since pretransforms are pure functions. A copy of the data is passed in. So if you place a preTransform on a large struct it will copy the entire struct which is not efficient. Be careful!
