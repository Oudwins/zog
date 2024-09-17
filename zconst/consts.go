package zconst

const (
	ZogTag = "zog"
)

type ZogType = string

const (
	TypeString ZogType = "string"
	TypeNumber ZogType = "number"
	TypeBool   ZogType = "bool"
	TypeTime   ZogType = "time"
	TypeSlice  ZogType = "slice"
	TypeStruct ZogType = "struct"
)

type ZogErrCode = string

const (
	ErrCodeCustom   ZogErrCode = "custom"         // all
	ErrCodeRequired ZogErrCode = "required"       // all
	ErrCodeCoerce   ZogErrCode = "coerce"         // all
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
)
