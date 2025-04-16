package internals

import (
	"fmt"
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
