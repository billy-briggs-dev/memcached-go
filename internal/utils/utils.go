package utils

import (
	"strconv"
)

// StringToInt converts a string to an integer. 
// It returns an error if the conversion fails.
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// IntToString converts an integer to a string.
func IntToString(i int) string {
	return strconv.Itoa(i)
}