package env

import (
	"os"
	"strconv"
)

type ConstraintType interface {
	~int | ~string | ~bool
}

func GetEnv[T ConstraintType](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T)
		}
	case string:
		return any(value).(T)
	case bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return any(boolValue).(T)
		}
	default:
		panic("Unsupported type")
	}

	return defaultValue
}
