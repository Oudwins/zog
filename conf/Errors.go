package conf

// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
import (
	"fmt"
	"strings"

	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// Map used to format errors in Zog. Both ZogType & ZogErrCode are just strings
type LangMap = map[zconst.ZogType]map[zconst.ZogErrCode]string

// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
var DefaultErrMsgMap LangMap = map[zconst.ZogType]map[zconst.ZogErrCode]string{
	zconst.TypeString: {
		zconst.ErrCodeRequired:        "is required",
		zconst.ErrCodeMin:             "string must contain at least {{min}} character(s)",
		zconst.ErrCodeMax:             "string must contain at most {{min}} character(s)",
		zconst.ErrCodeLen:             "string must be exactly {{len}} character(s)",
		zconst.ErrCodeEmail:           "must be a valid email",
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
		zconst.ErrCodeTrue:     "must be true",
		zconst.ErrCodeFalse:    "must be false",
		zconst.ErrCodeFallback: "value is invalid",
	},
	zconst.TypeNumber: {
		zconst.ErrCodeRequired: "is required",
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
		zconst.ErrCodeAfter:    "time must be after {{after}}",
		zconst.ErrCodeBefore:   "time must be before {{before}}",
		zconst.ErrCodeEQ:       "time must be equal to {{eq}}",
		zconst.ErrCodeFallback: "time is invalid",
	},
	zconst.TypeSlice: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeMin:      "slice must contain at least {{min}} items",
		zconst.ErrCodeMax:      "slice must contain at most {{max}} items",
		zconst.ErrCodeLen:      "slice must contain exactly {{len}} items",
		zconst.ErrCodeContains: "slice must contain {{contained}}",
		zconst.ErrCodeFallback: "slice is invalid",
	},
	zconst.TypeStruct: {
		zconst.ErrCodeRequired: "is required",
		zconst.ErrCodeFallback: "struct is invalid",
	},
}

func NewDefaultFormatter(m LangMap) p.ErrFmtFunc {
	return func(e p.ZogError, p p.ParseCtx) {
		if e.Message() != "" {
			return
		}
		// Check if the error msg is defined do nothing if it set
		t := e.Dtype()
		msg, ok := m[t][e.Code()]
		if !ok {
			e.SetMessage(m[t][zconst.ErrCodeFallback])
			return
		}
		for k, v := range e.Params() {
			// TODO replace this with a string builder
			msg = strings.ReplaceAll(msg, "{{"+k+"}}", fmt.Sprintf("%v", v))
		}
		e.SetMessage(msg)
	}

}

// Default error formatter it uses the errors above. Please override the `ErrorFormatter` variable instead of this one to customize the error messages for all zog schemas
var DefaultErrorFormatter p.ErrFmtFunc = NewDefaultFormatter(DefaultErrMsgMap)

// Override this
var ErrorFormatter = DefaultErrorFormatter
