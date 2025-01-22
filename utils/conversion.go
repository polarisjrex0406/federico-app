package utils

import (
	"fmt"
	"strconv"
)

func Float64ToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func StringToUint(s string) (uint, error) {
	// Parse the string to uint64 first
	val, err := strconv.ParseUint(s, 10, 0) // Base 10, bit size 0 for uint
	if err != nil {
		return 0, err // Return 0 and the error if parsing fails
	}

	// Check for overflow if converting to uint
	if val > uint64(^uint(0)) {
		return 0, fmt.Errorf("value %d is too large for uint", val)
	}

	return uint(val), nil // Convert uint64 to uint and return
}
