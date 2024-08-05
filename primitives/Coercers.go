package primitives

import (
	"fmt"
	"strconv"
)

// takes in a data and a pointer to a destination. Attempts to coerce the data into the destination. Returns an error if the coercion fails.
type CoercerFunc = func(data any, destPtr any) error

// a map of coercer functions. The key is the type of the destination and the value is the coercer function.
var Coercers = map[string]CoercerFunc{
	"bool": func(data any, destPtr any) error {
		dest := destPtr.(*bool)

		if b, ok := data.(bool); ok {
			*dest = b
			return nil
		}

		val, ok := data.(string)
		if !ok {
			return fmt.Errorf("cannot coerce %v to bool", data)
		}
		// There are cases where frontend libraries use "on" as the bool data
		// think about toggles. Hence, let's try this first.
		if data == "on" {
			*dest = true
		} else if data == "off" {
			*dest = false
		} else {
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("failed to parse bool: %v", err)
			}
			*dest = boolVal
		}
		return nil
	},
	"string": func(data any, destPtr any) error {
		dest := destPtr.(*string)

		if b, ok := data.(string); ok {
			*dest = b
			return nil
		}

		*dest = fmt.Sprintf("%v", data)

		return nil
	},
	"int": func(data any, destPtr any) error {
		dest := destPtr.(*int)
		if v, ok := data.(int); ok {
			*dest = v
			return nil
		}

		// TODO handle other types

		return fmt.Errorf("cannot coerce %v to int", data)
	},
	"float64": func(data any, destPtr any) error {
		dest := destPtr.(*float64)
		if v, ok := data.(float64); ok {
			*dest = v
			return nil
		}
		// TODO handle other types
		return nil
	},
	"time": func(data any, destPtr any) error {
		// TODO
		return nil
	},
}
