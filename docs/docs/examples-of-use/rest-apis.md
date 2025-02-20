---
sidebar_position: 2
---

# Using Zog in a REST API

Zog providers two helper functions called `z.Issues.SanitizeMap(issueMap)` and `z.Issues.SanitizeList(issueList)` that will return a map of strings of the issue messages (stripping out the internal error). So, if you do not mind sending issue messages to your users in the same form zog returns them, you can do something like this:

```go
errs := schema.Parse(zhttp.Request(r), &userFormData)

if errs != nil {
  sanitized := z.Issues.SanitizeMap(errs)
  // sanitize will be map[string][]string
  // for example:
  // {"name": []string{"min length is 5", "max length is 10"}, "email": []string{"is not a valid email"}}
  // ... marshal sanitized to json and send to the user
}

```
