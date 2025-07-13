package internals

import (
	"time"

	"golang.org/x/exp/constraints"
)

// type for functions called after validation & parsing is done
type PostTransform = func(dataPtr any, ctx Ctx) error

type Transform[T any] func(valPtr T, ctx Ctx) error

// Primitive types that can be used in Zod schemas
type ZogPrimitive interface {
	~string | ~bool | time.Time | constraints.Ordered
}

// Number like supported by Number schemas
type Numeric interface {
	constraints.Integer | constraints.Float
}
