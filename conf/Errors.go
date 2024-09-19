package conf

// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
import (
	"fmt"
	"strings"

	"github.com/Oudwins/zog/i18n/en"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
var DefaultErrMsgMap zconst.LangMap = en.Map

func NewDefaultFormatter(m zconst.LangMap) p.ErrFmtFunc {
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
