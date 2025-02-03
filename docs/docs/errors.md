---
sidebar_position: 4
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Errors

## Errors in Zog

In zog errors errors represent something that went wrong during any step of the [parsing execution structure](/core-concepts/parsing#parsing-execution-structure). Based on the schema you are using the returned errors will be in a different format:

**For Primitive Types**
Zog returns a list of `ZogError` instances.

```go
// will return []z.ZogError{z.ZogError{Message: "min length is 5"}, z.ZogError{Message: "invalid email"}}
errList := z.String().Min(5).Email().Parse("foo", &dest)
```

**For Complex Types**
Zog returns a map of `ZogError` instances. Which uses the field path as the key & the list of errors as the value.

```go
// will return map[string][]z.ZogError{"name": []z.ZogError{z.ZogError{Message: "min length is 5"}}}
errMap := z.Struct(z.Schema{"name": z.String().Min(5)}).Parse(data, &dest)

// will return map[string][]z.ZogError{"$root": []z.ZogError{{Message: "slice length is not 2"}, "[0]": []z.ZogError{{Message: "min length is 10"}}}}
errsMap2 := z.Slice(z.String().Min(10)).Len(2).Parse([]string{"only_one"}, &dest)

// nested schemas will use the . or the [] notation to access the errors
errsMap3 := z.Struct(z.Schema{"name": z.String().Min(5), "address": z.Struct(z.Schema{"streets": z.Slice(z.String().Min(10))})}).Parse(data, &dest)
errsMap3["address.streets[0]"] // will return []z.ZogError{{Message: "min length is 10"}}
```

`$root` & `$first` are reserved keys for complex type validation, they are used for root level errors and for the first error found in a schema, for example:

```go
errsMap := z.Slice(z.String()).Min(2).Parse([]string{"only_one"}, &dest)
errsMap["$root"] // will return []z.ZogError{{Message: "slice length should at least be 2"}}
errsMap["$first"] // will return the same in this case []z.ZogError{{Message: "slice length should at least be 2"}}
```

## The ZogError interface

The `ZogError` is actually an interface which also implements the error interface so it can be used with the `errors` package. The error interface is as follows:

```go
// Error interface returned from all schemas
type ZogError interface {
	// returns the error code for the error. This is a unique identifier for the error. Generally also the ID for the Test that caused the error.
	Code() zconst.ZogErrCode
	// returns the data value that caused the error.
	// if using Schema.Parse(data, dest) then this will be the value of data.
	Value() any
	// Returns destination type. i.e The zconst.ZogType of the value that was validated.
	// if Using Schema.Parse(data, dest) then this will be the type of dest.
	Dtype() string
	// returns the params map for the error. Taken from the Test that caused the error. This may be nil if Test has no params.
	Params() map[string]any
	// returns the human readable, user-friendly message for the error. This is safe to expose to the user.
	Message() string
	// sets the human readable, user-friendly message for the error. This is safe to expose to the user.
	SetMessage(string)
	// returns the string representation of the ZogError (same as String())
	Error() string
	// returns the wrapped error or nil if none
	Unwrap() error
	// returns the string representation of the ZogError (same as Error())
	String() string
}
// When printed it looks like this:
// ZogError{Code: coercion_error, Params: map[], Type: number, Value: not_empty, Message: number is invalid, Error: failed to coerce string int: strconv.Atoi: parsing "not_empty": invalid syntax}
```

## Error Codes

Error codes are unique identifiers for each type of error that can occur in Zog. They are used to generate error messages and to identify the error in the error formatter. A full updated list of error codes can be found in the zconst package. But here are some common ones:

```go
type ZogErrCode = string

const (
	ErrCodeCustom   ZogErrCode = "custom"   // all
	ErrCodeRequired ZogErrCode = "required" // all
	ErrCodeCoerce   ZogErrCode = "coerce"   // all
	ErrCodeFallback ZogErrCode = "fallback" // all. Applied when other errror code is not implemented. Required to be implemented for every zog type!

	ErrCodeEQ       ZogErrCode = "eq"             // number, time, string
	ErrCodeOneOf    ZogErrCode = "one_of_options" // string or number

	ErrCodeMin      ZogErrCode = "min"       // string, slice
	ErrCodeMax      ZogErrCode = "max"       // string, slice
	ErrCodeLen      ZogErrCode = "len"       // string, slice
	ErrCodeContains ZogErrCode = "contained" // string, slice

	// number only
	ErrCodeLTE ZogErrCode = "lte" // number
	ErrCodeLT  ZogErrCode = "lt"  // number
	ErrCodeGTE ZogErrCode = "gte" // number
	ErrCodeGT  ZogErrCode = "gt"  // number

	// string only
	ErrCodeEmail           ZogErrCode = "email"
	ErrCodeUUID            ZogErrCode = "uuid"
	ErrCodeMatch           ZogErrCode = "match"
	ErrCodeURL             ZogErrCode = "url"
	ErrCodeHasPrefix       ZogErrCode = "prefix"
	ErrCodeHasSuffix       ZogErrCode = "suffix"
	ErrCodeContainsUpper   ZogErrCode = "contains_upper"
	ErrCodeContainsLower   ZogErrCode = "contains_lower"
	ErrCodeContainsDigit   ZogErrCode = "contains_digit"
	ErrCodeContainsSpecial ZogErrCode = "contains_special"
	// time only
	ErrCodeAfter  ZogErrCode = "after"
	ErrCodeBefore ZogErrCode = "before"
	// bool only
	ErrCodeTrue  ZogErrCode = "true"
	ErrCodeFalse ZogErrCode = "false"

	// ZHTTP ERRORS
	ErrCodeZHTTPInvalidJSON  ZogErrCode = "invalid_json"  // invalid json body
	ErrCodeZHTTPInvalidForm  ZogErrCode = "invalid_form"  // invalid form data
	ErrCodeZHTTPInvalidQuery ZogErrCode = "invalid_query" // invalid query params
)
```

## Custom Error Messages

Zog has multiple ways of customizing error messages as well as support for [i18n](/packages/i18n). Here is a list of the ways you can customize error messages:

#### **1. Using the z.Message() function**

This is a function available for all tests, it allows you to set a custom message for the test.

```go
err := z.String().Min(5, z.Message("string must be at least 5 characters long")).Parse("bad", &dest)
// err = []ZogError{{Message: "string must be at least 5 characters long"}}
```

#### **2. Using the z.MessageFunc() function**

This is a function available for all tests, it allows you to set a custom message for the test.

This function takes in an `ErrFmtFunc` which is the function used to format error messages in Zog. It has the following signature:

```go
type ErrFmtFunc = func(e ZogError, p Ctx)
```

```go
err := z.String().Min(5, z.MessageFunc(func(e z.ZogError, p z.Ctx) {
  e.SetMessage("string must be at least 5 characters long")
})).Parse("bad", &dest)
// err = []ZogError{{Message: "string must be at least 5 characters long"}}
```

#### **3. Using the WithErrFormatter() ParsingOption**

This allows you to set a custom error formatter for the entire parsing operation. Beware you must handle all error codes & types or you may get unexpected messages.

```go
err := z.String().Min(5).Email().Parse("zog", &dest, z.WithErrFormatter(func(e z.ZogError, p z.Ctx) {
  e.SetMessage("override message")
}))
// err = []ZogError{{Code: min_length_error, Message: "override message"}, {Code: email_error, Message: "override message"}}
```

See how our error messages were overridden? Be careful when using this!

#### **4. Iterate over the returned errors and create custom messages**

```go
errs := userSchema.Parse(data, &user)
msgs := formatZogErrors(errs)

func FormatZogErrors(errs z.ZogErrMap) map[string][]string {
  // iterate over errors and create custom messages based on the error code, the params and destination type
}
```

#### **5. Configure error messages globally**

Zog provides a `conf` package where you can override the error messages for specific error codes. You will have to do a little digging to be able to do this. But here is an example:

```go
import (
  conf "github.com/Oudwins/zog/zconf"
  zconst "github.com/Oudwins/zog/zconst"
)

// override specific error messages
// For this I recommend you import `zod/zconst` which contains zog constants
conf.DefaultErrMsgMap[zconst.TypeString]["my_custom_error_code"] = "my custom error message"
conf.DefaultErrMsgMap[zconst.TypeString][zconst.ErrCodeRequired] = "Now all required errors will get this message"
```

But you can also outright override the error formatter and ignore the errors map completely:

```go
// override the error formatter function - CAREFUL with this you can set every error message to the same thing!
conf.ErrorFormatter = func(e p.ZogError, p z.Ctx) {
  // do something with the error
  ...
  // fallback to the default error formatter
  conf.DefaultErrorFormatter(e, p) // this uses the DefaultErrMsgMap to format the error messages
}
```

#### **6. Use the [i18n](/packages/i18n) package**

Really this only makes sense if you are doing i18n. Please please check out the [i18n](/packages/i18n) section for more information.

## Sanitizing Errors

If you want to return errors to the user without the possibility of exposing internal confidential information, you can use the Zog sanitizer functions `z.Errors.SanitizeMap(errsMap)` or `z.Errors.SanitizeSlice(errsSlice)`. These functions will return a map or slice of strings of the error messages (stripping out the internal error).

```go
errs := userSchema.Parse(data, &user)
// errs = map[string][]ZogError{"name": []ZogError{{Message: "min length is 5"}, {Message: "max length is 10"}}, "email": []ZogError{{Message: "is not a valid email"}}}
sanitized := z.Errors.SanitizeMap(errs)
// sanitized = {"name": []string{"min length is 5", "max length is 10"}, "email": []string{"is not a valid email"}}
```
