package utils

import (
	"os"
	"strconv"
)

func GetEnvOrElse(envKey string, defaultValue string) string {
	val, found := os.LookupEnv(envKey)
	if found {
		return val
	}

	return defaultValue
}

func GetEnvOrElseInt(envKey string, defaultValue int) int {
	val, found := os.LookupEnv(envKey)
	if found {
		num, err := strconv.Atoi(val)
		if err != nil {
			return int(num)
		}
	}

	return defaultValue
}
