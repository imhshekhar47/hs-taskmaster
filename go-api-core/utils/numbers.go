package utils

type NumericType interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type NonNegativeNumericType interface {
	uint | uint8 | uint16 | uint64
}

func GetNumOrElse[T NumericType](value T, defaultValue T) T {
	if value > 0 {
		return value
	} else {
		return defaultValue
	}
}
