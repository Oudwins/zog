package conf

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// takes in an original value and attempts to coerce it into another type. Returns an error if the coercion fails.
type CoercerFunc = func(original any) (value any, err error)

// The coercer functions used in zog by default.
var DefaultCoercers = struct {
	Bool    CoercerFunc
	String  CoercerFunc
	Int     CoercerFunc
	Float64 CoercerFunc
	Time    CoercerFunc
	Slice   CoercerFunc
}{
	Bool: func(data any) (any, error) {
		switch v := data.(type) {
		case bool:
			return v, nil
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
		case int:
			if v == 0 {
				return false, nil
			} else if v == 1 {
				return true, nil
			}
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to bool: %v", data)
		}

		return nil, fmt.Errorf("input data is an unsupported type to coerce to bool: %v", data)
	},
	String: func(data any) (any, error) {
		switch v := data.(type) {
		case string:
			return v, nil
		default:
			return fmt.Sprintf("%v", data), nil
		}

	},
	Int: func(data any) (any, error) {

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
	Float64: func(data any) (any, error) {
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
	Time: func(data any) (any, error) {
		switch v := data.(type) {
		case time.Time:
			return v, nil
		case string:
			tim, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time: %v", err)
			}
			return tim, nil
		case int:
			return time.Unix(int64(v), 0), nil
		case int64:
			return time.Unix(v, 0), nil
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to time.Time: %v", data)
		}
	},
	Slice: func(data any) (any, error) {
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

// Please override this variable instead of `DefaultCoercers` to add your own coercer functions.
var Coercers = DefaultCoercers
