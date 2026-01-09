package zog

import (
	p "github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

type lazySchema struct {
	fn func() ZogSchema
}

func (l *lazySchema) process(ctx *p.SchemaCtx) {
	x := l.fn()
	x.process(ctx)
}
func (l *lazySchema) validate(ctx *p.SchemaCtx) {}
func (l *lazySchema) getType() zconst.ZogType   { return l.fn().getType() }
func (l *lazySchema) setCoercer(c CoercerFunc)  { l.fn().setCoercer(c) }

func lazy(fn func() ZogSchema) *lazySchema {
	return &lazySchema{fn: fn}
}

type RecursiveSchemaFunc[T ZogSchema] func(self T) T
type RecursiveFunc[T ZogSchema] func(optionalFunc ...RecursiveSchemaFunc[T]) T

type RecursiveBuildFunc[T ZogSchema] func(self RecursiveFunc[T]) T

func Recursive[T ZogSchema](build RecursiveBuildFunc[T]) T {
	// var self ZogSchema
	// self = lazy(func() ZogSchema { return self })
	// real := build(self)
	// self = real
	// return real
	var self T
	// x := build()
	return self
}

var x = Recursive(func(self RecursiveFunc[*StructSchema]) *StructSchema {
	return Struct(Shape{
		"self": self(),
	})
})
