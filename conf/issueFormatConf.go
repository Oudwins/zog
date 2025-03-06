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

// Deprecated: Use DefaultIssueMsgMap instead
// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
var DefaultErrMsgMap zconst.LangMap = en.Map

// Default error messages for all schemas. Replace the text with your own messages to customize the error messages for all zog schemas
// As a general rule of thumb, if an error message only has one parameter, the parameter name will be the same as the error code
var DefaultIssueMessageMap zconst.LangMap = en.Map

const valuePlaceholder = "{{value}}"

func NewDefaultFormatter(m zconst.LangMap) p.IssueFmtFunc {
	return func(e *p.ZogIssue, c p.Ctx) {
		if e.Message != "" {
			return
		}
		// Check if the error msg is defined do nothing if it set
		t := e.Dtype
		msg, ok := m[t][e.Code]
		if !ok {
			e.SetMessage(m[t][zconst.IssueCodeFallback])
			return
		}
		for k, v := range e.Params {
			msg = strings.ReplaceAll(msg, "{{"+k+"}}", fmt.Sprintf("%v", v))
		}
		msg = strings.ReplaceAll(msg, valuePlaceholder, fmt.Sprintf("%v", e.Value))
		e.SetMessage(msg)
	}

}

// Default Issue Message formatter it uses the errors above. Please override the `IssueFormatter` variable instead of this one to customize the error messages for all zog schemas
var DefaultIssueFormatter p.IssueFmtFunc = NewDefaultFormatter(DefaultIssueMessageMap)

// Override this. This is the function use across all Zog schemas to format issue messages
/*
Usage:
```go
conf.IssueFormatter = func(e p.ZogIssue, c z.Ctx) {
     switch e.Code() {
     case zconst.IssueCodeCustom:
          e.SetMessage("Custom message")
	 default:
		conf.DefaultIssueFormatter(e, c) // fallback to default formatter
     }
}
```
*/
var IssueFormatter = DefaultIssueFormatter
