package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// Options that can be passed to a test
type TestOption = func(test *p.Test)

// Message is a function that allows you to set a custom message for the test.
func Message(msg string) TestOption {
	return func(test *p.Test) {
		test.IssueFmtFunc = func(e ZogIssue, p ParseCtx) {
			e.SetMessage(msg)
		}
	}
}

// MessageFunc is a function that allows you to set a custom message formatter for the test.
func MessageFunc(fn p.IssueFmtFunc) TestOption {
	return func(test *p.Test) {
		test.IssueFmtFunc = fn
	}
}

// IssueCode is a function that allows you to set a custom issue code for the test. Most useful for TestFuncs:
/*
z.String().TestFunc(..., z.IssueCode("just_provide_a_string" or use values in zconst))
*/
func IssueCode(code zconst.ZogIssueCode) TestOption {
	return func(test *p.Test) {
		test.IssueCode = code
	}
}

// Options that can be passed to a `schema.New()` call
type SchemaOption = func(s ZogSchema)

func WithCoercer(c conf.CoercerFunc) SchemaOption {
	return func(s ZogSchema) {
		s.setCoercer(c)
	}
}

// Options that can be passed to a `schema.Parse()` call
type ExecOption = func(p *p.ExecCtx)

// Deprecated: use ExecOption instead
type ParsingOption = ExecOption

// Deprecated: use WithIssueFormatter instead
// Deprecated for naming consistency
func WithErrFormatter(fmter p.IssueFmtFunc) ExecOption {
	return WithIssueFormatter(fmter)
}

// Sets the issue formatter for the execution context. This is used to format the issues messages during execution.
// This follows principle of most specific wins. So default formatter < execution formatter < test specific formatter (i.e MessageFunc)
func WithIssueFormatter(fmter p.IssueFmtFunc) ExecOption {
	return func(p *p.ExecCtx) {
		p.SetIssueFormatter(fmter)
	}
}

func WithCtxValue(key string, val any) ExecOption {
	return func(p *p.ExecCtx) {
		p.Set(key, val)
	}
}
