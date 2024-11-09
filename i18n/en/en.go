package en

import (
	"github.com/Oudwins/zog/zconst"
)

var Map zconst.LangMap = map[zconst.ZogType]map[zconst.ZogErrCode]string{
	zconst.TypeString: {
		zconst.ErrCodeRequired:        "is required",
		zconst.ErrCodeNotNil:          "must not be empty",
		zconst.ErrCodeMin:             "string must contain at least {{min}} character(s)",
		zconst.ErrCodeMax:             "string must contain at most {{max}} character(s)",
		zconst.ErrCodeLen:             "string must be exactly {{len}} character(s)",
		zconst.ErrCodeEmail:           "must be a valid email",
		zconst.ErrCodeUUID:            "must be a valid UUID",
		zconst.ErrCodeMatch:           "string is invalid",
		zconst.ErrCodeURL:             "must be a valid URL",
		zconst.ErrCodeHasPrefix:       "string must start with {{prefix}}",
		zconst.ErrCodeHasSuffix:       "string must end with {{suffix}}",
		zconst.ErrCodeContains:        "string must contain {{contained}}",
		zconst.ErrCodeContainsDigit:   "string must contain at least one digit",
		zconst.ErrCodeContainsUpper:   "string must contain at least one uppercase letter",
		zconst.ErrCodeContainsLower:   "string must contain at least one lowercase letter",
		zconst.ErrCodeContainsSpecial: "string must contain at least one special character",
		zconst.ErrCodeOneOf:           "string must be one of {{one_of_options}}",
		zconst.ErrCodeFallback:        "string is invalid",
	},
	zconst.TypeBool: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeNotNil:   "must not be empty",
		zconst.ErrCodeTrue:     "must be true",
		zconst.ErrCodeFalse:    "must be false",
		zconst.ErrCodeFallback: "value is invalid",
	},
	zconst.TypeNumber: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeNotNil:   "must not be empty",
		zconst.ErrCodeLTE:      "number must be less than or equal to {{lte}}",
		zconst.ErrCodeLT:       "number must be less than {{lt}}",
		zconst.ErrCodeGTE:      "number must be greater than or equal to {{gte}}",
		zconst.ErrCodeGT:       "number must be greater than {{gt}}",
		zconst.ErrCodeEQ:       "number must be equal to {{eq}}",
		zconst.ErrCodeOneOf:    "number must be one of {{options}}",
		zconst.ErrCodeFallback: "number is invalid",
	},
	zconst.TypeTime: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeNotNil:   "must not be empty",
		zconst.ErrCodeAfter:    "time must be after {{after}}",
		zconst.ErrCodeBefore:   "time must be before {{before}}",
		zconst.ErrCodeEQ:       "time must be equal to {{eq}}",
		zconst.ErrCodeFallback: "time is invalid",
	},
	zconst.TypeSlice: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeNotNil:   "must not be empty",
		zconst.ErrCodeMin:      "slice must contain at least {{min}} items",
		zconst.ErrCodeMax:      "slice must contain at most {{max}} items",
		zconst.ErrCodeLen:      "slice must contain exactly {{len}} items",
		zconst.ErrCodeContains: "slice must contain {{contained}}",
		zconst.ErrCodeFallback: "slice is invalid",
	},
	zconst.TypeStruct: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeNotNil:   "must not be empty",
		zconst.ErrCodeFallback: "struct is invalid",
		// ZHTTP ERRORS
		zconst.ErrCodeZHTTPInvalidJSON:  "invalid json body",
		zconst.ErrCodeZHTTPInvalidForm:  "invalid form data",
		zconst.ErrCodeZHTTPInvalidQuery: "invalid query params",
	},
}
