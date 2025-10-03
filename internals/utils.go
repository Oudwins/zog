package internals

import (
	"fmt"
	"reflect"
)

const defaultString = "<nil>"

func SafeString(x any) string {
	if x == nil {
		return defaultString
	}
	return fmt.Sprintf("%v", x)
}

func SafeError(x error) string {
	if x == nil {
		return defaultString
	}
	return x.Error()
}

func UnwrapPtr(x any) any {
	refVal := reflect.ValueOf(x)
	if refVal.Kind() != reflect.Ptr {
		return x
	}
	for refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}
	return refVal.Interface()
}
