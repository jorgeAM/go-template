package env

import (
	"os"
	"strconv"
)

type EnvType interface {
	~string | ~int | ~bool
}

func GetEnv[T EnvType](key string, defaultValue T) T {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T)
		}
	case float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return any(floatValue).(T)
		}
	case string:
		return any(value).(T)
	default:
		panic("unsupported type")
	}

	return defaultValue
}
