package internals

import (
	"reflect"
)

type IsZeroValueFunc = func(val any, ctx Ctx) bool

// checks that the value is the zero value for its type
func IsZeroValue(x any) bool {
	v := reflect.ValueOf(x)
	return !v.IsValid() || v.IsZero()
}

// checks if the value is the zero value but only for parsing purposes (i.e the parse function)
func IsParseZeroValue(val any, ctx Ctx) bool {
	if val == nil {
		return true
	}
	return IsExplicitNull(val)
}

// ExplicitNull is a sentinel for "key present with nil value" (e.g. JSON null)
// as opposed to "key absent". MapDataProvider.Get collapses both to bare nil
// without it, so PointerSchema.Nullable() has nothing to act on.
type ExplicitNull struct{}

// Unexported so external code cannot reassign or nil it out. Detection is
// type-based via IsExplicitNull, so any *ExplicitNull instance would still
// match, but emission paths need a stable singleton to return.
var explicitNullMarker = &ExplicitNull{}

// ExplicitNullMarker returns the canonical sentinel instance for emission sites
// that need to signal "present with nil value" to downstream schemas.
func ExplicitNullMarker() *ExplicitNull {
	return explicitNullMarker
}

func IsExplicitNull(val any) bool {
	_, ok := val.(*ExplicitNull)
	return ok
}
