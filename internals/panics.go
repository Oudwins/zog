package internals

import "fmt"

const (
	PanicTypeCast           = "Zog Panic: Type Cast Error\n Current context: %s\n Expected valPtr type to correspond with type defined in schema. But it does not. Expected type: *%T, got: %T\nFor more information see: https://zog.dev/panics#type-cast-errors"
	PanicTypeCastCoercer    = "Zog Panic: Type Cast Error\n Current context: %s\n Expected coercer return value to correspond with type defined in schema. But it does not. Expected type: *%T, got: %T\nFor more information see: https://zog.dev/panics#type-cast-errors"
	PanicMissingStructField = "Zog Panic: Struct Schema Definition Error\n Current context: %s\n Provided struct is missing expected schema key: %s.\n This means you have made a mistake in your schema definition.\nFor more information see: https://zog.dev/panics#schema-definition-errors"
)

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}
