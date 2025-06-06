---
sidebar_position: 2
---

# zjson

A very small package for using Zog schemas to parse json into structs. It exports a single function `Decode` which takes in an `io.Reader` or an `io.ReaderCloser` and returns the necessary structures for Zog to parse the json into a struct. This package is used by the `zhttp` package.

```go
import (
	"bytes"
	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
)

var userSchema = z.Struct(z.Shape{
	"name": z.String().Required(),
	"age":  z.Int().Required().GT(18),
})

type User struct {
	Name string
	Age  int
}

func ParseJson(json []byte) {
	var user User
	errs := userSchema.Parse(zjson.Decode(bytes.NewReader(json)), &user)

	if errs != nil {
		// handle errors
	}

	user.Name // defined
	user.Age  // defined
}
```

> **WARNING** The `zjson` package does NOT currently support parsing into any data type that is NOT a struct.

## Behaviour on unmarshal errors

If the json is not valid, a top level `ZogIssue` will be generated with the `IssueCode` `IssueCodeInvalidJSON` and the schema will not be run.
