// Package unsafe provides wrapper functions for 'unsafe' type conversions.
package unsafe

import "unsafe"

// String2ByteSlice converts a string to a byte slice.
func String2ByteSlice(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

// ByteSlice2String converts a byte slice to a string.
func ByteSlice2String(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
