package primitives

import "context"

type ZogParseCtx struct {
	context.Context
}

func (p *ZogParseCtx) Value() any {
	return nil
}

type ParseCtx interface {
	Value() any
}

func NewParseCtx() ParseCtx {
	return &ZogParseCtx{}
}

type ErrorFunc = func(val any, ctx ParseCtx) string

type TestFunc = func(val any, ctx ParseCtx) bool

type Test struct {
	Name         string
	ErrorFunc    ErrorFunc
	ValidateFunc TestFunc
}
