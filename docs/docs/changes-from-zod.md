---
sidebar_position: 10
hide_table_of_contents: false
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Changes from Zod

- Zog is Zod inspired, we adhere to the Zod API whenever possible but there are significant differences because:
  1. Go is statically typed and does not allow optional function params
  2. I have chosen to make Zog prioritize idiomatic Golang over the Zod API. Meaning some of the schemas & tests (validation rules) have changed names, `z.Array()` is `z.Slice()`, `z.String().StartsWith()` is `z.String().HasPrefix` (to follow the std lib). Etc.
  3. When I felt like a Zod method name would be confusing for Golang devs I changed it
- Some other changes:
  - The refine method for providing a custom validation function is renamed to `schema.Test()`
  - schemas are optional by default (in zod they are required)
  - The `z.Enum()` type from zod is removed in favor of `z.String().OneOf()` and is only supported for strings and numbers
  - `string().regex` is renamed to `z.String().Match()` as that is in line with the regexp methods from the standard library (i.e `regexp.Match` and `regexp.MatchString()`)
