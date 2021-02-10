package api

import (
	"strconv"
	"strings"
	"time"
)

// boolToPtr returns the pointer to a boolean
func boolToPtr(b bool) *bool {
	return &b
}

// int8ToPtr returns the pointer to an int8
func int8ToPtr(i int8) *int8 {
	return &i
}

// intToPtr returns the pointer to an int
func intToPtr(i int) *int {
	return &i
}

// uint64ToPtr returns the pointer to an uint64
func uint64ToPtr(u uint64) *uint64 {
	return &u
}

// int64ToPtr returns the pointer to a int64
func int64ToPtr(i int64) *int64 {
	return &i
}

// stringToPtr returns the pointer to a string
func stringToPtr(str string) *string {
	return &str
}

// timeToPtr returns the pointer to a time stamp
func timeToPtr(t time.Duration) *time.Duration {
	return &t
}

// formatFloat converts the floating-point number f to a string,
// after rounding it to the passed unit.
//
// Uses 'f' format (-ddd.dddddd, no exponent), and uses at most
// maxPrec digits after the decimal point.
func formatFloat(f float64, maxPrec int) string {
	v := strconv.FormatFloat(f, 'f', -1, 64)

	idx := strings.LastIndex(v, ".")
	if idx == -1 {
		return v
	}

	sublen := idx + maxPrec + 1
	if sublen > len(v) {
		sublen = len(v)
	}

	return v[:sublen]
}
