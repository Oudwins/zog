---
sidebar_position: 5
---

# zconst

`zconst` is a helper package that provides constants for common use cases such as error codes, Zog Types and more. Every constant here is just a string so using `zconst` is completely optional. This is the entire code of the zconst package as of version 0.11.0:

```go

const (
	ZogTag = "zog"
)

// Map used to format errors in Zog. Both ZogType & ZogIssueCode are just strings
type LangMap = map[ZogType]map[ZogIssueCode]string

type ZogType = string

const (
	TypeString ZogType = "string"
	TypeNumber ZogType = "number"
	TypeBool   ZogType = "bool"
	TypeTime   ZogType = "time"
	TypeSlice  ZogType = "slice"
	TypeStruct ZogType = "struct"
)

type ZogIssueCode = string

const (
	IssueCodeCustom   ZogIssueCode = "custom"   // all
	IssueCodeRequired ZogIssueCode = "required" // all
	IssueCodeNotNil   ZogIssueCode = "not_nil"  // all (technically only applies to pointers)
	IssueCodeCoerce   ZogIssueCode = "coerce"   // all
	// all. Applied when other errror code is not implemented. Required to be implemented for every zog type!
	IssueCodeFallback ZogIssueCode = "fallback"
	IssueCodeEQ       ZogIssueCode = "eq"             // number, time, string
	IssueCodeOneOf    ZogIssueCode = "one_of_options" // string or number

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
