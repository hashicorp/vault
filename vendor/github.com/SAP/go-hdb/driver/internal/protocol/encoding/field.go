package encoding

import (
	"math"

	"github.com/SAP/go-hdb/driver/unicode/cesu8"
)

const (
	booleanFalseValue byte = 0
	booleanNullValue  byte = 1
	booleanTrueValue  byte = 2
)

const (
	realNullValue   uint32 = ^uint32(0)
	doubleNullValue uint64 = ^uint64(0)
)

const (
	longdateNullValue   int64 = 3155380704000000001
	seconddateNullValue int64 = 315538070401
	daydateNullValue    int32 = 3652062
	secondtimeNullValue int32 = 86402
)

// Field size constants.
const (
	BooleanFieldSize       = 1
	TinyintFieldSize       = 1
	SmallintFieldSize      = 2
	IntegerFieldSize       = 4
	BigintFieldSize        = 8
	RealFieldSize          = 4
	DoubleFieldSize        = 8
	DateFieldSize          = 4
	TimeFieldSize          = 4
	TimestampFieldSize     = DateFieldSize + TimeFieldSize
	LongdateFieldSize      = 8
	SeconddateFieldSize    = 8
	DaydateFieldSize       = 4
	SecondtimeFieldSize    = 4
	DecimalFieldSize       = 16
	Fixed8FieldSize        = 8
	Fixed12FieldSize       = 12
	Fixed16FieldSize       = 16
	LobInputParametersSize = 9
)

// string / binary length indicators.
const (
	bytesLenIndNullValue byte = 255
	bytesLenIndSmall     byte = 245
	bytesLenIndMedium    byte = 246
	bytesLenIndBig       byte = 247
)

// VarFieldSize returns the size of a varible field variable ([]byte, string and unicode variants).
func varSize(size int) int {
	switch {
	default:
		return -1
	case size <= int(bytesLenIndSmall):
		return size + 1
	case size <= math.MaxInt16:
		return size + 3
	case size <= math.MaxInt32:
		return size + 5
	}
}

// Cesu8FieldSize returns the size of a cesu8 field.
func Cesu8FieldSize(v any) int {
	switch v := v.(type) {
	case []byte:
		return varSize(cesu8.Size(v))
	case string:
		return varSize(cesu8.StringSize(v))
	default:
		panic("invalid type for cesu8 field") // should never happen
	}
}

// VarFieldSize returns the size of a var field.
func VarFieldSize(v any) int {
	switch v := v.(type) {
	case []byte:
		return varSize(len(v))
	case string:
		return varSize(len(v))
	default:
		panic("invalid type for var field") // should never happen
	}
}

// HexFieldSize returns the size of a hex field.
func HexFieldSize(v any) int {
	switch v := v.(type) {
	case []byte:
		l := len(v)
		if l%2 != 0 {
			panic("even hex field length required")
		}
		return varSize(l / 2)
	case string:
		l := len(v)
		if l%2 != 0 {
			panic("even hex field length required")
		}
		return varSize(l / 2)
	default:
		panic("invalid hex field type")
	}
}
