package internals

import (
	"time"

	"golang.org/x/exp/constraints"
)

// type for functions called after validation & parsing is done
type PostTransform = func(dataPtr any, ctx Ctx) error

type Transform = func(valPtr any, ctx Ctx) error // TODO turn this into generic alias & upgrade to go 1.23. Maybe this can't actually be turned into a generic right? because of parse vs validate? Actually no, yes it can because it will always take a pointer to the value because we do the parsing before hand

// Primitive types that can be used in Zod schemas
type ZogPrimitive interface {
	~string | ~bool | time.Time | constraints.Ordered
}
