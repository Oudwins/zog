package zog

import (
	p "github.com/Oudwins/zog/internals"
)

type TestOption = func(test *p.Test)

// Message is a function that allows you to set a custom message for the test.
func Message(msg string) TestOption {
	return func(test *p.Test) {
		test.ErrFmt = func(e p.ZogError, p ParseCtx) {
			e.SetMessage(msg)
		}
	}
}

// MessageFunc is a function that allows you to set a custom message formatter for the test.
func MessageFunc(fn p.ErrFmtFunc) TestOption {
	return func(test *p.Test) {
		test.ErrFmt = fn
	}
}

type ParsingOption = func(p *p.ZogParseCtx)

func WithErrFormatter(fmter p.ErrFmtFunc) ParsingOption {
	return func(p *p.ZogParseCtx) {
		p.SetErrFormatter(fmter)
	}
}

func WithCtxValue(key string, val any) ParsingOption {
	return func(p *p.ZogParseCtx) {
		p.Set(key, val)
	}
}
