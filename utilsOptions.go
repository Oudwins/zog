package zog

import (
	p "github.com/Oudwins/zog/internals"
)

type TestOption = func(test *p.Test)

func Message(msg string) TestOption {
	return func(test *p.Test) {
		test.ErrFmt = func(e p.ZogError, p ParseCtx) {
			e.SetMessage(msg)
		}
	}
}

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
