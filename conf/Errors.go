package conf

import (
	"fmt"
	"strings"

	p "github.com/Oudwins/zog/primitives"
)

// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
var DefaultErrMsgMap = map[p.ZogType]map[p.ZogErrCode]string{
	p.TypeString: {
		p.ErrCodeRequired:        "is required",
		p.ErrCodeMin:             "string must contain at least {{min}} character(s)",
		p.ErrCodeMax:             "string must contain at most {{min}} character(s)",
		p.ErrCodeLen:             "string must be exactly {{len}} character(s)",
		p.ErrCodeEmail:           "must be a valid email",
		p.ErrCodeURL:             "must be a valid URL",
		p.ErrCodeHasPrefix:       "string must start with {{prefix}}",
		p.ErrCodeHasSuffix:       "string must end with {{suffix}}",
		p.ErrCodeContains:        "string must contain {{contained}}",
		p.ErrCodeContainsDigit:   "string must contain at least one digit",
		p.ErrCodeContainsUpper:   "string must contain at least one uppercase letter",
		p.ErrCodeContainsLower:   "string must contain at least one lowercase letter",
		p.ErrCodeContainsSpecial: "string must contain at least one special character",
		p.ErrCodeOneOf:           "string must be one of {{one_of_options}}",
	},
	p.TypeBool: {
		p.ErrCodeRequired: "is required",
		p.ErrCodeTrue:     "must be true",
		p.ErrCodeFalse:    "must be false",
	},
	p.TypeNumber: {
		p.ErrCodeRequired: "is required",
		p.ErrCodeLTE:      "number must be less than or equal to {{lte}}",
		p.ErrCodeLT:       "number must be less than {{lt}}",
		p.ErrCodeGTE:      "number must be greater than or equal to {{gte}}",
		p.ErrCodeGT:       "number must be greater than {{gt}}",
		p.ErrCodeEQ:       "number must be equal to {{eq}}",
		p.ErrCodeOneOf:    "number must be one of {{options}}",
	},
	p.TypeTime: {
		p.ErrCodeRequired: "is required",
		p.ErrCodeAfter:    "time must be after {{after}}",
		p.ErrCodeBefore:   "time must be before {{before}}",
		p.ErrCodeEQ:       "time must be equal to {{eq}}",
	},
	p.TypeSlice: {
		p.ErrCodeRequired: "is required",
		p.ErrCodeMin:      "slice must contain at least {{min}} items",
		p.ErrCodeMax:      "slice must contain at most {{max}} items",
		p.ErrCodeLen:      "slice must contain exactly {{len}} items",
		p.ErrCodeContains: "slice must contain {{contained}}",
	},
	p.TypeStruct: {
		p.ErrCodeRequired: "is required",
	},
}

// Default error formatter it uses the errors above. Please override the `ErrorFormatter` variable instead of this one to customize the error messages for all zog schemas
var DefaultErrorFormatter p.ErrFmtFunc = func(e p.ZogError, p p.ParseCtx) {
	if e.Message() != "" {
		return
	}
	// Check if the error msg is defined do nothing if it set
	t := e.Dtype()
	msg, ok := DefaultErrMsgMap[t][e.Code()]
	if !ok {
		e.SetMessage(t + " is invalid")
		return
	}
	for k, v := range e.Params() {
		// TODO replace this with a string builder
		msg = strings.ReplaceAll(msg, "{{"+k+"}}", fmt.Sprintf("%v", v))
	}
	e.SetMessage(msg)
}

// Override this
var ErrorFormatter = DefaultErrorFormatter
