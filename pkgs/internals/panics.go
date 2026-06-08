package internals

import "fmt"

const (
	PanicTypeCastCoercer = "Zog Panic: Type Cast Error\n Current context: %s\n Expected coercer return value to correspond with type defined in schema. But it does not. Expected type: *%T, got: %T\nFor more information see: https://zog.dev/panics#type-cast-errors"
)

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
