package internals

import "time"

// takes the data as input and returns the new data which will then be passed onto the next functions. If the function returns an error all validation will be skipped & the error will be returned. You may return a ZogIssue or an error. If its an error it will be wraped inside a ZogIssue
type PreTransform = func(data any, ctx Ctx) (out any, err error)

// type for functions called after validation & parsing is done
type PostTransform = func(dataPtr any, ctx Ctx) error

// Primitive types that can be used in Zod schemas
type ZogPrimitive interface {
	~string | ~int | ~float64 | ~bool | time.Time
}
