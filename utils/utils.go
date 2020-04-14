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

// Precision returns the precision of a double
func Precision(s string) int {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return 0
	}

	return len(parts[1])
}

// MaxInt returns the maximum for a given set of integers
func MaxInt(args ...int) int {
	if len(args) == 0 {
		return 0
	}

	max := args[0]

	for i := 1; i < len(args); i++ {
		if args[i] > max {
			max = args[i]
		}
	}

	return max
}
