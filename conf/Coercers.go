package conf

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// takes in an original value and attempts to coerce it into another type. Returns an error if the coercion fails.
type CoercerFunc = func(original any) (value any, err error)

// a map of coercer functions. The key is the type of the destination and the value is the coercer function. You may override this map to add your own coercer functions and they will affect the behaviour of all zog schemas.
var Coercers = map[string]CoercerFunc{
	"bool": func(data any) (any, error) {
		if b, ok := data.(bool); ok {
			return b, nil
		}

		switch v := data.(type) {
		case string:
			// There are cases where frontend libraries use "on" as the bool data
			// think about toggles. Hence, let's try this first.
			if v == "on" {
				return true, nil
			} else if v == "off" {
				return false, nil
			} else {
				boolVal, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("failed to coerce string to parse bool: %v", err)
				}
				return boolVal, nil
			}
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to bool: %v", data)
		}
	},
	"string": func(data any) (any, error) {
		switch v := data.(type) {
		case string:
			return v, nil
		default:
			return fmt.Sprintf("%v", data), nil
		}

	},
	"int": func(data any) (any, error) {

		switch v := data.(type) {
		case int:
			return v, nil
		case string:
			convVal, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to coerce string int: %v", err)
			}
			return convVal, nil
		case float64:
			return int(v), nil
		case bool:
			if v {
				return 1, nil
			} else {
				return 0, nil
			}
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to int: %v", data)
		}
	},
	"float64": func(data any) (any, error) {
		switch v := data.(type) {
		case int:
			return float64(v), nil
		case string:
			convVal, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to coerce string to float64: %v", err)
			}
			return convVal, nil
		case float64:
			return v, nil
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to float64: %v", data)
		}
	},
	"time": func(data any) (any, error) {
		switch v := data.(type) {
		case time.Time:
			return v, nil
		case string:
			tim, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time: %v", err)
			}
			return tim, nil
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to time.Time: %v", data)
		}
	},
	"slice": func(data any) (any, error) {
		refVal := reflect.TypeOf(data)
		switch refVal.Kind() {
		case reflect.Slice:
			return data, nil
		// any other type we box
		default:
			return []any{data}, nil
		}
	},
}
