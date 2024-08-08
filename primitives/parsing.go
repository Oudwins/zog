package primitives

import "context"

type ParseCtx struct {
	context.Context
	path string
}

func NewParseCtx() *ParseCtx {
	return &ParseCtx{
		path: "",
	}
}

type ErrorFunc = func(val any, ctx *ParseCtx) string

type TestFunc = func(val any, ctx *ParseCtx) bool

type Test struct {
	Name         string
	ErrorFunc    ErrorFunc
	ValidateFunc TestFunc
}
