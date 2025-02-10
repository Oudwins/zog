package zconst

const (
	// Key for a list of one error that contains the first error that occurred in a schema
	ERROR_KEY_FIRST = "$first"
	// Key for a list of all errors that occurred in a schema at the root level for complex schemas. For example
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
	ERROR_KEY_ROOT = "$root"
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
type LangMap = map[ZogType]map[ZogErrCode]string

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

type ZogErrCode = string

const (
	ErrCodeCustom   ZogErrCode = "custom"   // all
	ErrCodeRequired ZogErrCode = "required" // all
	ErrCodeNotNil   ZogErrCode = "not_nil"  // all (technically only applies to pointers)
	ErrCodeCoerce   ZogErrCode = "coerce"   // all
	// all. Applied when other errror code is not implemented. Required to be implemented for every zog type!
	ErrCodeFallback ZogErrCode = "fallback"
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

	// JSON
	ErrCodeInvalidJSON ZogErrCode = "invalid_json" // invalid json body
	// ZHTTP ERRORS
	ErrCodeZHTTPInvalidForm  ZogErrCode = "invalid_form"  // invalid form data
	ErrCodeZHTTPInvalidQuery ZogErrCode = "invalid_query" // invalid query params
)
