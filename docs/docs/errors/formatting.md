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
    path: []string{"[0]"}
    message: "string must start with 'PREFIX_'"
  },
  {
    code: 'prefix',
    path: []string{"[1]"}
    message: "string must start with 'PREFIX_'"
  },
  {
	code: "min",
	path: nil,
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
> 1. issues with a nil or empty path will be assigned to `$root` reserved key
> 2. Struct/map keys are mapped to their key names and joined by `.` (For example `user.firstname` where the path is `[]string{"user", "firstname"}`)
> 3. Slices are mapped to their index and can be appended to a previous struct/map key. For example `[]string{"[0]"}`, `[]string{"[0]", "firstname"}` and `[]string{"users", "[0]", "firstname"}` are all valid paths

## `z.Issues.Treeify()`

This converts the errors list into a nested tree structure that mirrors the original data structure. This format is particularly useful when you need to display errors in a hierarchical way that matches your form or data model.

For the example above, `Treeify` would generate:

```json
{
  "errors": ["slice must contain at least 3 items"],
  "properties": {
    "items": [
      {
        "errors": ["string must start with 'PREFIX_'"]
      },
      {
        "errors": ["string must start with 'PREFIX_'"]
      }
    ]
  }
}
```

Here's a more complex example with nested structures:

```go
schema := z.Struct{
  "user": z.Struct{
    "name": z.String().Min(3),
    "email": z.String().Email(),
  },
  "users": z.Slice(z.Struct{
    "name": z.String().Required(),
  }),
}

data := map[string]any{
  "user": map[string]any{
    "name": "ab",  // too short
    "email": "invalid",  // not an email
  },
  "users": []any{
    map[string]any{"name": ""},  // required
    map[string]any{"name": "ok"},
  },
}

errs := schema.Parse(data, &dest)
tree := z.Issues.Treeify(errs)
```

The resulting tree structure:

```json
{
  "errors": [],
  "properties": {
    "user": {
      "errors": [],
      "name": {
        "errors": ["string must be at least 3 characters"]
      },
      "email": {
        "errors": ["string must be a valid email"]
      }
    },
    "users": {
      "errors": [],
      "items": [
        {
          "errors": [],
          "name": {
            "errors": ["string is required"]
          }
        },
        null
      ]
    }
  }
}
```

> **How does treeify logic work?**
> The tree structure follows these rules:
>
> 1. Root-level errors (nil or empty path) are placed in the top-level `errors` array
> 2. Property errors create nested objects under `properties`, with each path segment becoming a nested level
> 3. Array indices create an `items` array within the parent property, with each index becoming an element in the array
> 4. Numeric string segments (like `"0"`) are treated as array indices and create `items` arrays
> 5. Each node in the tree has an `errors` array, even if empty, to maintain consistent structure

## `z.Issues.Prettify()`

This formats the errors list into a human-readable string representation. Each issue is displayed with a "✖" prefix, and issues with paths include the path information on a separate line with a "→ at" prefix. This format is ideal for displaying errors directly to users in console output or error messages.

For the example above, `Prettify` would generate:

```
✖ string must start with 'PREFIX_'
  → at [0]
✖ string must start with 'PREFIX_'
  → at [1]
✖ slice must contain at least 3 items
```

Here's a more complex example:

```go
schema := z.Struct{
  "user": z.Struct{
    "name": z.String().Min(3),
    "email": z.String().Email(),
  },
  "users": z.Slice(z.Struct{
    "name": z.String().Required(),
  }),
}

data := map[string]any{
  "user": map[string]any{
    "name": "ab",  // too short
    "email": "invalid",  // not an email
  },
  "users": []any{
    map[string]any{"name": ""},  // required
  },
}

errs := schema.Parse(data, &dest)
pretty := z.Issues.Prettify(errs)
```

The resulting formatted string:

```
✖ string must be at least 3 characters
  → at user.name
✖ string must be a valid email
  → at user.email
✖ string is required
  → at users[0].name
```

> **How does prettify logic work?**
> The formatting follows these rules:
>
> 1. Empty issue lists return an empty string
> 2. Each issue message is prefixed with "✖ " (checkmark symbol)
> 3. Issues with paths (that flatten to a non-empty string) include the path on a new line with " → at " prefix
> 4. Multiple issues are separated by newlines (`\n`)
> 5. Root-level errors (nil or empty path) are displayed without a path line
> 6. Paths are flattened using the same logic as `Flatten()`, so nested properties use dot notation (e.g., `user.name`) and array indices use bracket notation (e.g., `users[0]`)
