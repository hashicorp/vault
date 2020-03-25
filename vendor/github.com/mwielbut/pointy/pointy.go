// Package pointy is a set of simple helper functions to provide a shorthand to
// get a pointer to a variable holding a constant.
package pointy

// Bool returns a pointer to a variable holding the supplied bool constant
func Bool(x bool) *bool {
	return &x
}

// BoolValue returns the bool value pointed to by p or fallback if p is nil
func BoolValue(p *bool, fallback bool) bool {
	if p == nil {
		return fallback
	}
	return *p
}

// Byte returns a pointer to a variable holding the supplied byte constant
func Byte(x byte) *byte {
	return &x
}

// ByteValue returns the byte value pointed to by p or fallback if p is nil
func ByteValue(p *byte, fallback byte) byte {
	if p == nil {
		return fallback
	}
	return *p
}

// Complex128 returns a pointer to a variable holding the supplied complex128 constant
func Complex128(x complex128) *complex128 {
	return &x
}

// Complex128Value returns the complex128 value pointed to by p or fallback if p is nil
func Complex128Value(p *complex128, fallback complex128) complex128 {
	if p == nil {
		return fallback
	}
	return *p
}

// Complex64 returns a pointer to a variable holding the supplied complex64 constant
func Complex64(x complex64) *complex64 {
	return &x
}

// Complex64Value returns the complex64 value pointed to by p or fallback if p is nil
func Complex64Value(p *complex64, fallback complex64) complex64 {
	if p == nil {
		return fallback
	}
	return *p
}

// Float32 returns a pointer to a variable holding the supplied float32 constant
func Float32(x float32) *float32 {
	return &x
}

// Float32Value returns the float32 value pointed to by p or fallback if p is nil
func Float32Value(p *float32, fallback float32) float32 {
	if p == nil {
		return fallback
	}
	return *p
}

// Float64 returns a pointer to a variable holding the supplied float64 constant
func Float64(x float64) *float64 {
	return &x
}

// Float64Value returns the float64 value pointed to by p or fallback if p is nil
func Float64Value(p *float64, fallback float64) float64 {
	if p == nil {
		return fallback
	}
	return *p
}

// Int returns a pointer to a variable holding the supplied int constant
func Int(x int) *int {
	return &x
}

// IntValue returns the int value pointed to by p or fallback if p is nil
func IntValue(p *int, fallback int) int {
	if p == nil {
		return fallback
	}
	return *p
}

// Int8 returns a pointer to a variable holding the supplied int8 constant
func Int8(x int8) *int8 {
	return &x
}

// Int8Value returns the int8 value pointed to by p or fallback if p is nil
func Int8Value(p *int8, fallback int8) int8 {
	if p == nil {
		return fallback
	}
	return *p
}

// Int16 returns a pointer to a variable holding the supplied int16 constant
func Int16(x int16) *int16 {
	return &x
}

// Int16Value returns the int16 value pointed to by p or fallback if p is nil
func Int16Value(p *int16, fallback int16) int16 {
	if p == nil {
		return fallback
	}
	return *p
}

// Int32 returns a pointer to a variable holding the supplied int32 constant
func Int32(x int32) *int32 {
	return &x
}

// Int32Value returns the int32 value pointed to by p or fallback if p is nil
func Int32Value(p *int32, fallback int32) int32 {
	if p == nil {
		return fallback
	}
	return *p
}

// Int64 returns a pointer to a variable holding the supplied int64 constant
func Int64(x int64) *int64 {
	return &x
}

// Int64Value returns the int64 value pointed to by p or fallback if p is nil
func Int64Value(p *int64, fallback int64) int64 {
	if p == nil {
		return fallback
	}
	return *p
}

// Uint returns a pointer to a variable holding the supplied uint constant
func Uint(x uint) *uint {
	return &x
}

// UintValue returns the uint value pointed to by p or fallback if p is nil
func UintValue(p *uint, fallback uint) uint {
	if p == nil {
		return fallback
	}
	return *p
}

// Uint8 returns a pointer to a variable holding the supplied uint8 constant
func Uint8(x uint8) *uint8 {
	return &x
}

// Uint8Value returns the uint8 value pointed to by p or fallback if p is nil
func Uint8Value(p *uint8, fallback uint8) uint8 {
	if p == nil {
		return fallback
	}
	return *p
}

// Uint16 returns a pointer to a variable holding the supplied uint16 constant
func Uint16(x uint16) *uint16 {
	return &x
}

// Uint16Value returns the uint16 value pointed to by p or fallback if p is nil
func Uint16Value(p *uint16, fallback uint16) uint16 {
	if p == nil {
		return fallback
	}
	return *p
}

// Uint32 returns a pointer to a variable holding the supplied uint32 constant
func Uint32(x uint32) *uint32 {
	return &x
}

// Uint32Value returns the uint32 value pointed to by p or fallback if p is nil
func Uint32Value(p *uint32, fallback uint32) uint32 {
	if p == nil {
		return fallback
	}
	return *p
}

// Uint64 returns a pointer to a variable holding the supplied uint64 constant
func Uint64(x uint64) *uint64 {
	return &x
}

// Uint64Value returns the uint64 value pointed to by p or fallback if p is nil
func Uint64Value(p *uint64, fallback uint64) uint64 {
	if p == nil {
		return fallback
	}
	return *p
}

// String returns a pointer to a variable holding the supplied string constant
func String(x string) *string {
	return &x
}

// StringValue returns the string value pointed to by p or fallback if p is nil
func StringValue(p *string, fallback string) string {
	if p == nil {
		return fallback
	}
	return *p
}

// Rune returns a pointer to a variable holding the supplied rune constant
func Rune(x rune) *rune {
	return &x
}

// RuneValue returns the rune value pointed to by p or fallback if p is nil
func RuneValue(p *rune, fallback rune) rune {
	if p == nil {
		return fallback
	}
	return *p
}
