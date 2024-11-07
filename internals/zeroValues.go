package internals

import (
	"reflect"
	"strings"
)

type IsZeroValueFunc = func(val any, ctx ParseCtx) bool

// checks that the value is the zero value for its type
func IsZeroValue(x any) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	if !v.IsValid() {
		return true
	}

	// Check if the value is the zero value for its type
	zeroValue := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zeroValue.Interface())
}

// checks if the value is the zero value but only for parsing purposes (i.e the parse function)
func IsParseZeroValue(val any, ctx ParseCtx) bool {
	if val == nil {
		return true
	}
	s, ok := val.(string)
	if ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}
