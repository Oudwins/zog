---
sidebar_position: 1
---

# zhttp

For Zog provides a built in helper package called `zhttp` that helps parse JSON, Forms or Query Params. Helps parse a request into a struct by using the Content-Type header to infer the type of the request. Example usage below:

```go
import (
  z "github.com/Oudwins/zog"
  "github.com/Oudwins/zog/zhttp"
)
var userSchema = z.Struct(z.Schema{
  "name": z.String().Required(),
  "age": z.Int().Required().GT(18),
})

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
  var user := struct {
    Name string
    Age int
  }
  // if using json (i.e json Content-Type header):
  errs := userSchema.Parse(zhttp.Request(r), &user)
  // if using form data (i.e Content-Type header = application/x-www-form-urlencoded)
  errs := userSchema.Parse(zhttp.Request(r), &user)
  // if using query params (i.e no http Content-Type header)
  errs := userSchema.Parse(zhttp.Request(r), &user)
  if errs != nil {
  }
  user.Name // defined
  user.Age // defined
}
```

> **WARNING** The `zhttp` package does NOT currently support parsing into any data type that is NOT a struct.
