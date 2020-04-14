package utils

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
