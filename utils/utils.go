package utils

import (
	"strings"
)

// InArray returns true if the needle exists in the haystack
func InArray(needle interface{}, haystack []interface{}) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// Abs returns the absolute value of an integer
func Abs(a int64) int64 {
	if a < 0 {
		return a * -1
	}

	return a
}

// Precision returns the precision of a double
func Precision(s string) int {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return 0
	}

	return len(parts[1])
}
