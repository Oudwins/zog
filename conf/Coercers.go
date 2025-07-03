package conf

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

func TimeCoercerFactory(format func(data string) (time.Time, error)) CoercerFunc {
	return func(data any) (any, error) {
		switch v := data.(type) {
		case time.Time:
			return v, nil
		case string:
			tim, err := format(v)
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
	}
}

// takes in an original value and attempts to coerce it into another type. Returns an error if the coercion fails.
type CoercerFunc = func(original any) (value any, err error)

// The coercer functions used in zog by default.
var DefaultCoercers = struct {
	Bool    CoercerFunc
	String  CoercerFunc
	Int     CoercerFunc
	Float64 CoercerFunc
	Uint    CoercerFunc
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
		case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
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
		case int64:
			if v > math.MaxInt || v < math.MinInt {
				return nil, fmt.Errorf("int64 value %d overflows int", v)
			}
			return int(v), nil
		case int32:
			// int32 always fits in int on 32-bit and 64-bit systems
			return int(v), nil
		case int16:
			return int(v), nil
		case int8:
			return int(v), nil
		case uint64:
			if v > math.MaxInt {
				return nil, fmt.Errorf("uint64 value %d overflows int", v)
			}
			return int(v), nil
		case uint32:
			if uint64(v) > math.MaxInt {
				return nil, fmt.Errorf("uint32 value %d overflows int", v)
			}
			return int(v), nil
		case uint16:
			if uint64(v) > math.MaxInt {
				return nil, fmt.Errorf("uint16 value %d overflows int", v)
			}
			return int(v), nil
		case uint8:
			return int(v), nil
		case float64:
			if v > math.MaxInt || v < math.MinInt {
				return nil, fmt.Errorf("float64 value %g overflows int", v)
			}
			return int(v), nil
		case float32:
			if float64(v) > math.MaxInt || float64(v) < math.MinInt {
				return nil, fmt.Errorf("float32 value %g overflows int", v)
			}
			return int(v), nil
		case string:
			convVal, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to coerce string int: %v", err)
			}
			return convVal, nil

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
	Uint: func(data any) (any, error) {
		switch v := data.(type) {
		case uint:
			return v, nil
		case int:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case int64:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case int32:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case int16:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case uint64:
			return uint(v), nil
		case uint32:
			return uint(v), nil
		case uint16:
			return uint(v), nil
		case uint8:
			return uint(v), nil
		case float64:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case float32:
			if v < 0 {
				return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
			}
			return uint(v), nil
		case bool:
			if v {
				return uint(1), nil
			} else {
				return uint(0), nil
			}
		case string:
			convVal, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to coerce string to uint: %v", err)
			}
			return uint(convVal), nil
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to uint: %v", data)
		}
	},
	Float64: func(data any) (any, error) {
		switch v := data.(type) {
		case string:
			convVal, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to coerce string to float64: %v", err)
			}
			return convVal, nil
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case int32:
			return float64(v), nil
		case int16:
			return float64(v), nil
		case int8:
			return float64(v), nil
		case uint:
			return float64(v), nil
		case uint64:
			return float64(v), nil
		case uint32:
			return float64(v), nil
		case uint16:
			return float64(v), nil
		case uint8:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("input data is an unsupported type to coerce to float64: %v", data)
		}
	},
	Time: TimeCoercerFactory(func(data string) (time.Time, error) {
		return time.Parse(time.RFC3339, data)
	}),
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
