package internals

import "time"

// takes the data as input and returns the new data which will then be passed onto the next functions. If the function returns an error all validation will be skipped & the error will be returned
type PreTransform = func(data any, ctx ParseCtx) (out any, err error)

// type for functions called after validation & parsing is done
type PostTransform = func(dataPtr any, ctx ParseCtx) error

// Primitive types that can be used in Zod schemas
type ZogPrimitive interface {
	~string | ~int | ~float64 | ~bool | time.Time
}
