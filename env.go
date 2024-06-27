package zog

import (
	"fmt"
	"os"
	"strconv"

	p "github.com/Oudwins/zog/primitives"
)

// takes a key and a validator and  returns the validated and converted environment variable
func Env[T supportedEnvTypes](key string, v fieldValidator, defailtValue ...T) T {
	str := os.Getenv(key)

	val, err := coerceString[T](str)

	if err != nil || p.IsZeroValue(val) {
		if len(defailtValue) > 0 {
			return defailtValue[0]
		} else {
			panic(fmt.Errorf("failed to parse env %s: %v", key, err))
		}
	}

	errs, ok := v.Validate(val)
	if !ok {
		panic(fmt.Errorf("failed to validate env %s: %v", key, errs))
	}

	return val
}

type supportedEnvTypes interface {
	~int | ~float64 | ~bool | ~string
}

func coerceString[T supportedEnvTypes](val string) (T, error) {
	var result T

	switch any(result).(type) {
	case int:
		var tmp int
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return result, err
		}

		result = any(tmp).(T)
	case float64:
		tmp, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return result, err
		}
		result = any(tmp).(T)

	case bool:
		tmp, err := strconv.ParseBool(val)
		if err != nil {
			return result, err
		}
		result = any(tmp).(T)
	case string:
		result = any(val).(T)
	default:
		return result, fmt.Errorf("unsupported type: %T", result)
	}

	return result, nil
}
