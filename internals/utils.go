package internals

import (
	"fmt"
	"reflect"
)

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

func SafeString(x any) string {
	if x == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", x)
}

func SafeError(x error) string {
	if x == nil {
		return "<nil>"
	}
	return x.Error()
}
