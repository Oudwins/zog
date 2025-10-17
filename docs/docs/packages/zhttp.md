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

var userSchema = z.Struct(z.Shape{
	"name": z.String().Required(),
	"age":  z.Int().Required().GT(18),
})

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Name string
		Age  int
	}

	// if using json (i.e json Content-Type header):
	errs := userSchema.Parse(zhttp.Request(r), &user)
	// if using form data (i.e Content-Type header = application/x-www-form-urlencoded)
	errs := userSchema.Parse(zhttp.Request(r), &user)
	// if using multipart form data you are expected to parse the form yourself before using it with zhttp. See this article on why/how to correctly parse multipart form data: https://medium.com/@owlwalks/dont-parse-everything-from-client-multipart-post-golang-9280d23cd4ad
	// After that you can just use it as normal
	errs := userSchema.Parse(zhttp.Request(r), &user)
	// if using query params (i.e no http Content-Type header)
	errs := userSchema.Parse(zhttp.Request(r), &user)




	if errs != nil {
		// ...
	}

	user.Name // defined
	user.Age  // defined
}
```

> **WARNING** The `zhttp` package does NOT currently support parsing into any data type that is NOT a struct.

## Behaviour on unmarshal errors

If the json, form or query params are not valid, a top level `ZogIssue` will be generated with the `IssueCode` `IssueCodeInvalidJSON` or `IssueCodeZHTTPInvalidForm` or `IssueCodeZHTTPInvalidQuery` and the schema will not be run.

## Complex Forms

If you need to parse complex forms or query params such as those parsed by packages like [qs](https://www.npmjs.com/package/qs), for example:

```js
assert.deepEqual(qs.parse("foo[bar]=baz"), {
  foo: {
    bar: "baz",
  },
});
```

zhttp does not currently support these types of forms (see [issue #8](https://github.com/Oudwins/zog/issues/8)). However I suggest you try using the [form go package](https://github.com/go-playground/form) which supports this type of parsing. You can integrate the library with zhttp by overriding the `zhttp.Config.Parsers.Form` function.

> **WARNING**: This depends on `DataProviders` which are not yet documented and may change in the future. I encourage you to avoid doing this unless you really need to.
