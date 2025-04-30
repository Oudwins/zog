---
sidebar_position: 4
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Preprocess

Just like in Zod you can use preprocess function to transform the data before it is validated. This is also useful for things like type coercion. **Preprocess functions are PURE functions**. They take in data and return new data. This is the function signature:

```go
func Preprocess[F any, T any](fn func(data F, ctx Ctx) (out T, err error), schema ZogSchema) *PreprocessSchema[F, T]

// Usage:
// Note that even if preprocess takes a generic [F]rom type, my recommendation is to always set that to any unless you are 100% sure that the input data will always be of a specific type. Since if you are using this schema with schema.Parse() the input data can be anything.
z.Preprocess(func(data any, ctx z.ctx) (any, error) {
	s, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("expected string but got %T", data)
	}
	return strings.split(s, ","), nil
}, z.Slice(z.String())))
```

You can use `z.Preprocess` for things like trimming whitespace, splitting strings, etc. Here is an example of splitting a string into a slice of strings:

> **FOOTGUNS**.
> _parse vs validate_: z.Preprocess can run for both `schema.Parse()` and `schema.Validate()`, in each case the data argument will be different!. For `schema.Parse()` the data argument is the value you are parsing (i.e the input data). For `schema.Validate()` the data argument is the pointer to the value you are validating.
> _Pure Functions_: Since preprocess functions are pure functions. They create copies of the data. So be careful when using them with large data structures if you are concerned about performance.
