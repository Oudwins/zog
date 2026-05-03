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
// without it, so pointer schemas cannot distinguish the two cases.
type ExplicitNull struct{}

// Unexported so external packages cannot reassign the singleton.
var explicitNullMarker = &ExplicitNull{}

func ExplicitNullMarker() *ExplicitNull {
	return explicitNullMarker
}

func IsExplicitNull(val any) bool {
	_, ok := val.(*ExplicitNull)
	return ok
}
