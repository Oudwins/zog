package primitives

import (
	"fmt"
	"strconv"
	"time"
)

// takes in a data and a pointer to a destination. Attempts to coerce the data into the destination. Returns an error if the coercion fails.
type CoercerFunc = func(data any, destPtr any) error

// a map of coercer functions. The key is the type of the destination and the value is the coercer function. You may override this map to add your own coercer functions and they will affect the behaviour of all zog schemas.
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

		switch v := data.(type) {
		case int:
			*dest = v
			return nil
		case string:
			convVal, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("failed to parse int: %v", err)
			}
			*dest = convVal
			return nil
		case float64:
			*dest = int(v)
		case bool:
			if v {
				*dest = 1
			} else {
				*dest = 0
			}
			return nil
		}
		return fmt.Errorf("cannot coerce %v to int", data)
	},
	"float64": func(data any, destPtr any) error {
		dest := destPtr.(*float64)
		switch v := data.(type) {
		case int:
			*dest = float64(v)
			return nil
		case string:
			convVal, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int: %v", err)
			}
			*dest = convVal
			return nil
		case float64:
			*dest = v
		case bool:
			if v {
				*dest = 1.0
			} else {
				*dest = 0.0
			}
			return nil
		}
		return nil
	},
	"time": func(data any, destPtr any) error {
		dest := destPtr.(*time.Time)

		switch v := data.(type) {
		case time.Time:
			*dest = v
			return nil
		case string:
			tim, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return fmt.Errorf("failed to parse time: %v", err)
			}
			*dest = tim
			return nil
		}

		return fmt.Errorf("cannot coerce %v to time.Time", data)
	},
}
