---
sidebar_position: 1
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Overview

In zog issues represent something that went wrong during any step of the [parsing execution structure](/core-concepts/parsing#parsing-execution-structure).

Most of the time, you will be interacting with Zog issues through:

- [Creating custom issue messages](/errors/custom-messages)
- [Formatting issue responses](/errors/formatting)

[Zod](https://github.com/colinhacks/zod) has both *Errors* and *Issues*. An *Issue* is the structured validation result produced by Zod, while an *Error* is the underlying exception that may have caused the issue.

In Zog, errors are never thrown directly; validation tests return issues instead. Each issue may optionally contain an underlying error that provides additional internal context.

For example, a JSON unmarshalling failure may produce a low-level error. Rather than exposing that raw error to the user, Zog wraps it in a `z.Issue`, preserving the internal cause while presenting a consistent validation interface.

## The ZogIssue

The `ZogIssue` is a struct that contains multiple fields and has some helper setter methods. It looks like this:

```go
// ZogIssue represents an issue that occurred during parsing or validation.
// When printed it looks like:
// ZogIssue{Code: coercion_issue, Params: map[], Type: number, Value: not_empty, Message: number is invalid, Error: failed to coerce string int: strconv.Atoi: parsing "not_empty": invalid syntax}
type ZogIssue struct {
	// Code is the unique identifier for the issue. Generally also the ID for the Test that caused the issue. See below for a full list
	Code zconst.ZogIssueCode
	// Path is the path to the field that caused the issue
	Path []string
	// Value is the data value that caused the issue.
	// If using Schema.Parse(data, dest) then this will be the value of data.
	Value any
	// Dtype is the destination type. i.e The zconst.ZogType of the value that was validated.
	// If using Schema.Parse(data, dest) then this will be the type of dest.
	Dtype string
	// Params is the params map for the issue. Taken from the Test that caused the issue.
	// This may be nil if Test has no params.
	Params map[string]any
	// Message is the human readable, user-friendly message for the issue.
	// This is safe to expose to the user.
	Message string
	// Err is the wrapped error or nil if none
	Err error
}
```

## ZogIssueList

All schemas return a `ZogIssueList`, which is a slice of `ZogIssue` instances:

```go
type ZogIssueList = []*ZogIssue
```

Each issue contains a `Path` field that indicates where in your data structure the issue occurred:

```go
// String schema - Path is nil for root-level primitive
errList := z.String().Min(5).Email().Parse("foo", &dest)
// errList[0].Path = nil, errList[0].Message = "min length is 5"

// Struct schema - Path contains the field name as a slice
errList := z.Struct(z.Shape{"name": z.String().Min(5)}).Parse(data, &dest)
// errList[0].Path = []string{"name"}, errList[0].Message = "min length is 5"

// Nested structs and slices - Path is a slice of path segments
errList := z.Struct(z.Shape{
    "address": z.Struct(z.Shape{
        "streets": z.Slice(z.String().Min(10)),
    }),
}).Parse(data, &dest)
// errList[0].Path = []string{"address", "streets", "[0]"}, errList[0].Message = "min length is 10"
// You can convert Path to a string using Issues.FlattenPath(issue.Path) or issue.PathString()
```

## Issue Codes

Error codes are unique identifiers for each type of issue that can occur in Zog. They are used to generate issue messages and to identify the issue in the issue formatter. A full updated list of issue codes can be found in the zconst package. But here are some common ones:

```go
type ZogIssueCode = string

const (
	IssueCodeCustom   ZogIssueCode = "custom"   // all
	IssueCodeRequired ZogIssueCode = "required" // all
	IssueCodeCoerce   ZogIssueCode = "coerce"   // all
	IssueCodeFallback ZogIssueCode = "fallback" // all. Applied when other errror code is not implemented. Required to be implemented for every zog type!

	IssueCodeEQ    ZogIssueCode = "eq"             // number, time, string
	IssueCodeOneOf ZogIssueCode = "one_of_options" // string or number

	IssueCodeMin      ZogIssueCode = "min"       // string, slice
	IssueCodeMax      ZogIssueCode = "max"       // string, slice
	IssueCodeLen      ZogIssueCode = "len"       // string, slice
	IssueCodeContains ZogIssueCode = "contained" // string, slice

	// number only
	IssueCodeLTE ZogIssueCode = "lte" // number
	IssueCodeLT  ZogIssueCode = "lt"  // number
	IssueCodeGTE ZogIssueCode = "gte" // number
	IssueCodeGT  ZogIssueCode = "gt"  // number

	// string only
	IssueCodeEmail           ZogIssueCode = "email"
	IssueCodeUUID            ZogIssueCode = "uuid"
	IssueCodeMatch           ZogIssueCode = "match"
	IssueCodeURL             ZogIssueCode = "url"
	IssueCodeHasPrefix       ZogIssueCode = "prefix"
	IssueCodeHasSuffix       ZogIssueCode = "suffix"
	IssueCodeContainsUpper   ZogIssueCode = "contains_upper"
	IssueCodeContainsLower   ZogIssueCode = "contains_lower"
	IssueCodeContainsDigit   ZogIssueCode = "contains_digit"
	IssueCodeContainsSpecial ZogIssueCode = "contains_special"
	// time only
	IssueCodeAfter  ZogIssueCode = "after"
	IssueCodeBefore ZogIssueCode = "before"
	// bool only
	IssueCodeTrue  ZogIssueCode = "true"
	IssueCodeFalse ZogIssueCode = "false"

	// ZHTTP ERRORS
	IssueCodeZHTTPInvalidJSON  ZogIssueCode = "invalid_json"  // invalid json body
	IssueCodeZHTTPInvalidForm  ZogIssueCode = "invalid_form"  // invalid form data
	IssueCodeZHTTPInvalidQuery ZogIssueCode = "invalid_query" // invalid query params
)
```
