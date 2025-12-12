---
sidebar_position: 3
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Formatting error responses

Zog emphasizes completeness and correctness in its error reporting. In many cases, it's helpful to convert the `ZogIssueList` to a more useful format. Zog provides some utilities for this.

Consider this simple slice schema.

```go
schema := z.Slice(z.String().Required().HasPrefix("PREFIX_")).Min(3);
```

Attempting to parse/validate this invalid data results in an error containing three issues.

```go
s := []string{"one", "two"}
errs := schema.Parse(&s)

errs;
[
  {
    code: 'prefix',
    path: "[0]"
    message: "string must start with 'PREFIX_'"
  },
  {
    code: 'prefix',
    path: "[1]"
    message: "string must start with 'PREFIX_'"
  },
  {
	code: "min",
	path: "",
	message: "slice must contain at least 3 items"
  }
]
```

## `z.Issues.Flatten()`

This converts the errors list to a `map[path][]messages`. For the example above it would generate:

```json
{
  "$root": ["string must start with 'PREFIX_'"], // Use zconst.ISSUE_KEY_ROOT for a constant of this key!
  "[0]": ["string must start with 'PREFIX_'"],
  "[1]": ["slice must contain at least 3 items"]
}
```

> **How does flatten logic work?**
> It follows a few simple rules:
>
> 1. issues with an empty path will be assigned to `$root` reserved key
> 2. Struct/map keys are mapped to their key names and joined by `.` (For example `user.firstname` where the path is first user key inside struct/map then its firstname key)
> 3. Slices are mapped to their index and can be appended to a previous struct/map key. For example `[0]`, `[0].firstname` and `users[0].firstname` are all valid keys
