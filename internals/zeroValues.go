package internals

import (
	"reflect"
	"strings"
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
	s, ok := val.(string)
	if ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}
