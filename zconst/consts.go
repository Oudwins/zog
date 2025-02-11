package zconst

const (
	// Key for a list of one ZogIssue that contains the first ZogIssue that occurred in a schema
	ISSUE_KEY_FIRST = "$first"
	// Key for a list of all ZogIssues that occurred in a schema at the root level for complex schemas. For example
	/*
			> Given this schema:
			z.Struct(....).TestFunc(func (v any, ctx z.Ctx) {
			   return false
			}, z.Message("test"))
			> And any input data. The output will contain:
			{
				"$root": [
					{
				     "message": "test"
					 restOfErrorFields..
					}
				]
			}
		    > This is also true for slices and even for pointers to primitive types.
	*/
	ISSUE_KEY_ROOT = "$root"
)

const (
	// Tag used for unmarshaling and finding the data to parse into a struct. Usage:
	// type TestStruct struct {
	// 	Value *int `zog:"value"`
	// }
	// Similar to `json` tag. But works with all input data sources at once (i.e query params, form data, json etc)
	ZogTag = "zog"
)

// Map used to format errors in Zog. Both ZogType & ZogErrCode are just strings
type LangMap = map[ZogType]map[ZogIssueCode]string

type ZogType = string

const (
	TypeString ZogType = "string"
	TypeNumber ZogType = "number"
	TypeBool   ZogType = "bool"
	TypeTime   ZogType = "time"
	TypeSlice  ZogType = "slice"
	TypeStruct ZogType = "struct"
	TypePtr    ZogType = "ptr"
)

// Deprecated: This will be removed in the future. Use z.ZogIssueCode instead
type ZogErrCode = string

// This is a type for the ZogIssueCode type. It is a string that represents the error code for an issue.
type ZogIssueCode = string

const (
	// Deprecated: Use IssueCodeCustom instead
	ErrCodeCustom   ZogErrCode   = "custom" // all
	IssueCodeCustom ZogIssueCode = "custom" // all

	// Deprecated: Use IssueCodeRequired instead
	ErrCodeRequired   ZogErrCode   = "required" // all
	IssueCodeRequired ZogIssueCode = "required" // all

	// Deprecated: Use IssueCodeNotNil instead
	ErrCodeNotNil   ZogErrCode   = "not_nil" // all (technically only applies to pointers)
	IssueCodeNotNil ZogIssueCode = "not_nil" // all (technically only applies to pointers)

	// Deprecated: Use IssueCodeCoerce instead
	ErrCodeCoerce   ZogErrCode   = "coerce" // all
	IssueCodeCoerce ZogIssueCode = "coerce" // all

	// Deprecated: Use IssueCodeFallback instead
	// all. Applied when other errror code is not implemented. Required to be implemented for every zog type!
	ErrCodeFallback ZogErrCode = "fallback"
	// all. Applied when other errror code is not implemented. Required to be implemented for every zog type!
	IssueCodeFallback ZogIssueCode = "fallback"

	// Deprecated: Use IssueCodeEQ instead
	ErrCodeEQ   ZogErrCode   = "eq" // number, time, string
	IssueCodeEQ ZogIssueCode = "eq" // number, time, string

	// Deprecated: Use IssueCodeOneOf instead
	ErrCodeOneOf   ZogErrCode   = "one_of_options" // string or number
	IssueCodeOneOf ZogIssueCode = "one_of_options" // string or number

	// Deprecated: Use IssueCodeMin instead
	ErrCodeMin   ZogErrCode   = "min" // string, slice
	IssueCodeMin ZogIssueCode = "min" // string, slice

	// Deprecated: Use IssueCodeMax instead
	ErrCodeMax   ZogErrCode   = "max" // string, slice
	IssueCodeMax ZogIssueCode = "max" // string, slice

	// Deprecated: Use IssueCodeLen instead
	ErrCodeLen   ZogErrCode   = "len" // string, slice
	IssueCodeLen ZogIssueCode = "len" // string, slice

	// Deprecated: Use IssueCodeContains instead
	ErrCodeContains   ZogErrCode   = "contained" // string, slice
	IssueCodeContains ZogIssueCode = "contained" // string, slice

	// number only
	// Deprecated: Use IssueCodeLTE instead
	ErrCodeLTE   ZogErrCode   = "lte" // number
	IssueCodeLTE ZogIssueCode = "lte" // number

	// Deprecated: Use IssueCodeLT instead
	ErrCodeLT   ZogErrCode   = "lt" // number
	IssueCodeLT ZogIssueCode = "lt" // number

	// Deprecated: Use IssueCodeGTE instead
	ErrCodeGTE   ZogErrCode   = "gte" // number
	IssueCodeGTE ZogIssueCode = "gte" // number

	// Deprecated: Use IssueCodeGT instead
	ErrCodeGT   ZogErrCode   = "gt" // number
	IssueCodeGT ZogIssueCode = "gt" // number

	// string only
	// Deprecated: Use IssueCodeEmail instead
	ErrCodeEmail   ZogErrCode   = "email"
	IssueCodeEmail ZogIssueCode = "email"

	// Deprecated: Use IssueCodeUUID instead
	ErrCodeUUID   ZogErrCode   = "uuid"
	IssueCodeUUID ZogIssueCode = "uuid"

	// Deprecated: Use IssueCodeMatch instead
	ErrCodeMatch   ZogErrCode   = "match"
	IssueCodeMatch ZogIssueCode = "match"

	// Deprecated: Use IssueCodeURL instead
	ErrCodeURL   ZogErrCode   = "url"
	IssueCodeURL ZogIssueCode = "url"

	// Deprecated: Use IssueCodeHasPrefix instead
	ErrCodeHasPrefix   ZogErrCode   = "prefix"
	IssueCodeHasPrefix ZogIssueCode = "prefix"

	// Deprecated: Use IssueCodeHasSuffix instead
	ErrCodeHasSuffix   ZogErrCode   = "suffix"
	IssueCodeHasSuffix ZogIssueCode = "suffix"

	// Deprecated: Use IssueCodeContainsUpper instead
	ErrCodeContainsUpper   ZogErrCode   = "contains_upper"
	IssueCodeContainsUpper ZogIssueCode = "contains_upper"

	// Deprecated: Use IssueCodeContainsLower instead
	ErrCodeContainsLower   ZogErrCode   = "contains_lower"
	IssueCodeContainsLower ZogIssueCode = "contains_lower"

	// Deprecated: Use IssueCodeContainsDigit instead
	ErrCodeContainsDigit   ZogErrCode   = "contains_digit"
	IssueCodeContainsDigit ZogIssueCode = "contains_digit"

	// Deprecated: Use IssueCodeContainsSpecial instead
	ErrCodeContainsSpecial   ZogErrCode   = "contains_special"
	IssueCodeContainsSpecial ZogIssueCode = "contains_special"

	// time only
	// Deprecated: Use IssueCodeAfter instead
	ErrCodeAfter   ZogErrCode   = "after"
	IssueCodeAfter ZogIssueCode = "after"

	// Deprecated: Use IssueCodeBefore instead
	ErrCodeBefore   ZogErrCode   = "before"
	IssueCodeBefore ZogIssueCode = "before"

	// bool only
	// Deprecated: Use IssueCodeTrue instead
	ErrCodeTrue   ZogErrCode   = "true"
	IssueCodeTrue ZogIssueCode = "true"

	// Deprecated: Use IssueCodeFalse instead
	ErrCodeFalse   ZogErrCode   = "false"
	IssueCodeFalse ZogIssueCode = "false"

	// JSON
	// Deprecated: Use IssueCodeInvalidJSON instead
	ErrCodeInvalidJSON   ZogErrCode   = "invalid_json" // invalid json body
	IssueCodeInvalidJSON ZogIssueCode = "invalid_json" // invalid json body

	// ZHTTP ERRORS
	// Deprecated: Use IssueCodeZHTTPInvalidForm instead
	ErrCodeZHTTPInvalidForm   ZogErrCode   = "invalid_form" // invalid form data
	IssueCodeZHTTPInvalidForm ZogIssueCode = "invalid_form" // invalid form data

	// Deprecated: Use IssueCodeZHTTPInvalidQuery instead
	ErrCodeZHTTPInvalidQuery   ZogErrCode   = "invalid_query" // invalid query params
	IssueCodeZHTTPInvalidQuery ZogIssueCode = "invalid_query" // invalid query params
)
