package utils

import "strings"

func IsEmpty(word string) bool {
	return len(strings.Trim(word, " ")) == 0
}

func NonEmptyOr(word string, defaultValue string) string {
	if !IsEmpty(word) {
		return word
	}
	return defaultValue
}
