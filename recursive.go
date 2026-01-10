package zog

import (
	"sync"

	p "github.com/Oudwins/zog/pkgs/internals"
	zss "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
)

type lazySchema struct {
	innerSchema ZogSchema
	fn          func() ZogSchema
	once        sync.Once
}

var _ ZogSchema = &lazySchema{}

func (l *lazySchema) get() ZogSchema {
	l.once.Do(func() {
		l.innerSchema = l.fn()
	})
	return l.innerSchema
}

func (l *lazySchema) process(ctx *p.SchemaCtx) {
	l.get().process(ctx)
}
func (l *lazySchema) validate(ctx *p.SchemaCtx) {
	l.get().validate(ctx)
}
func (l *lazySchema) getType() zconst.ZogType  { return l.get().getType() }
func (l *lazySchema) setCoercer(c CoercerFunc) { l.get().setCoercer(c) }
func (l *lazySchema) toZSS() *zss.ZSSSchema    { return l.get().toZSS() }

func lazy(fn func() ZogSchema) *lazySchema {
	return &lazySchema{fn: fn}
}

type RecursiveSchemaUpdater[T ZogSchema] func(self T) T
type RecursiveSchema[T ZogSchema] func(updaters ...RecursiveSchemaUpdater[T]) ZogSchema
type RecursiveSchemaBuilder[T ZogSchema] func(self RecursiveSchema[T]) T

// Experimental API.
// Do not use unless you know what you are doing.
func EXPERIMENTAL_RECURSIVE[T ZogSchema](build RecursiveSchemaBuilder[T]) T {
	var self T
	var lazyBuilder RecursiveSchema[T] = func(updaters ...RecursiveSchemaUpdater[T]) ZogSchema {
		return lazy(func() ZogSchema {
			if len(updaters) > 0 {
				return updaters[0](self)
			}
			return self
		})
	}
	self = build(lazyBuilder)
	return self
}
