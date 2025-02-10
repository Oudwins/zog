package en

import (
	"github.com/Oudwins/zog/zconst"
)

var Map zconst.LangMap = map[zconst.ZogType]map[zconst.ZogIssueCode]string{
	zconst.TypeString: {
		zconst.IssueCodeRequired:        "is required",
		zconst.IssueCodeNotNil:          "must not be empty",
		zconst.IssueCodeMin:             "string must contain at least {{min}} character(s)",
		zconst.IssueCodeMax:             "string must contain at most {{max}} character(s)",
		zconst.IssueCodeLen:             "string must be exactly {{len}} character(s)",
		zconst.IssueCodeEmail:           "must be a valid email",
		zconst.IssueCodeUUID:            "must be a valid UUID",
		zconst.IssueCodeMatch:           "string is invalid",
		zconst.IssueCodeURL:             "must be a valid URL",
		zconst.IssueCodeHasPrefix:       "string must start with {{prefix}}",
		zconst.IssueCodeHasSuffix:       "string must end with {{suffix}}",
		zconst.IssueCodeContains:        "string must contain {{contained}}",
		zconst.IssueCodeContainsDigit:   "string must contain at least one digit",
		zconst.IssueCodeContainsUpper:   "string must contain at least one uppercase letter",
		zconst.IssueCodeContainsLower:   "string must contain at least one lowercase letter",
		zconst.IssueCodeContainsSpecial: "string must contain at least one special character",
		zconst.IssueCodeOneOf:           "string must be one of {{one_of_options}}",
		zconst.IssueCodeFallback:        "string is invalid",
	},
	zconst.TypeBool: {
		zconst.IssueCodeRequired: "is required",
		zconst.IssueCodeNotNil:   "must not be empty",
		zconst.IssueCodeTrue:     "must be true",
		zconst.IssueCodeFalse:    "must be false",
		zconst.IssueCodeFallback: "value is invalid",
	},
	zconst.TypeNumber: {
		zconst.IssueCodeRequired: "is required",
		zconst.IssueCodeNotNil:   "must not be empty",
		zconst.IssueCodeLTE:      "number must be less than or equal to {{lte}}",
		zconst.IssueCodeLT:       "number must be less than {{lt}}",
		zconst.IssueCodeGTE:      "number must be greater than or equal to {{gte}}",
		zconst.IssueCodeGT:       "number must be greater than {{gt}}",
		zconst.IssueCodeEQ:       "number must be equal to {{eq}}",
		zconst.IssueCodeOneOf:    "number must be one of {{options}}",
		zconst.IssueCodeFallback: "number is invalid",
	},
	zconst.TypeTime: {
		zconst.IssueCodeRequired: "is required",
		zconst.IssueCodeNotNil:   "must not be empty",
		zconst.IssueCodeAfter:    "time must be after {{after}}",
		zconst.IssueCodeBefore:   "time must be before {{before}}",
		zconst.IssueCodeEQ:       "time must be equal to {{eq}}",
		zconst.IssueCodeFallback: "time is invalid",
	},
	zconst.TypeSlice: {
		zconst.IssueCodeRequired: "is required",
		zconst.IssueCodeNotNil:   "must not be empty",
		zconst.IssueCodeMin:      "slice must contain at least {{min}} items",
		zconst.IssueCodeMax:      "slice must contain at most {{max}} items",
		zconst.IssueCodeLen:      "slice must contain exactly {{len}} items",
		zconst.IssueCodeContains: "slice must contain {{contained}}",
		zconst.IssueCodeFallback: "slice is invalid",
	},
	zconst.TypeStruct: {
		zconst.IssueCodeRequired: "is required",
		zconst.IssueCodeNotNil:   "must not be empty",
		zconst.IssueCodeFallback: "struct is invalid",
		// JSON
		zconst.IssueCodeInvalidJSON: "invalid json body",
		// ZHTTP ERRORS
		zconst.IssueCodeZHTTPInvalidForm:  "invalid form data",
		zconst.IssueCodeZHTTPInvalidQuery: "invalid query params",
	},
}
