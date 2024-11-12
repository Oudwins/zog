package zog

import (
	"github.com/Oudwins/zog/conf"
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// The ZogSchema is the interface all schemas must implement
type ZogSchema interface {
	process(val any, dest any, path p.PathBuilder, ctx ParseCtx)
	setCoercer(c conf.CoercerFunc)
	getType() zconst.ZogType
}

// ! Passing Types through

// ParseCtx is the context passed through the parser
type ParseCtx = p.ParseCtx

// ZogError is the ZogError interface
type ZogError = p.ZogError

// ZogErrList is a []ZogError returned from parsing primitive schemas
type ZogErrList = p.ZogErrList

// ZogErrMap is a map[string][]ZogError returned from parsing complex schemas
type ZogErrMap = p.ZogErrMap

// ! TESTS

// Test is the test object
type Test = p.Test
