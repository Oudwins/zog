---
sidebar_position: 2
---

# Using Zog in a REST API

Zog providers a helper function called `z.Errors.SanitizeMap(errsMap)` that will return a map of strings of the error messages (stripping out the internal error). So, if you do not mind sending errors to your users in the same form zog returns them, you can do something like this:

```go
errs := schema.Parse(zhttp.Request(r), &userFormData)

if errs != nil {
  sanitized := z.Errors.SanitizeMap(errs)
  // sanitize will be map[string][]string
  // for example:
  // {"name": []string{"min length is 5", "max length is 10"}, "email": []string{"is not a valid email"}}
  // ... marshal sanitized to json and send to the user
}

```
