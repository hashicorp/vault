package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

type optType interface {
	fmt.Stringer
	typeCode() typeCode
	size(v any) int
	encode(e *encoding.Encoder, v any)
	decode(d *encoding.Decoder) any
}

var (
	optBooleanType = _optBooleanType{}
	optTinyintType = _optTinyintType{}
	optIntegerType = _optIntegerType{}
	optBigintType  = _optBigintType{}
	optDoubleType  = _optDoubleType{}
	optStringType  = _optStringType{}
	optBstringType = _optBstringType{}
)

type (
	_optBooleanType struct{}
	_optTinyintType struct{}
	_optIntegerType struct{}
	_optBigintType  struct{}
	_optDoubleType  struct{}
	_optStringType  struct{}
	_optBstringType struct{}
)

var (
	_ optType = (*_optBooleanType)(nil)
	_ optType = (*_optTinyintType)(nil)
	_ optType = (*_optIntegerType)(nil)
	_ optType = (*_optBigintType)(nil)
	_ optType = (*_optDoubleType)(nil)
	_ optType = (*_optStringType)(nil)
	_ optType = (*_optBstringType)(nil)
)

func (_optBooleanType) String() string { return "booleanType" }
func (_optTinyintType) String() string { return "tinyintType" }
func (_optIntegerType) String() string { return "integerType" }
func (_optBigintType) String() string  { return "bigintType" }
func (_optDoubleType) String() string  { return "doubleType" }
func (_optStringType) String() string  { return "dateType" }
func (_optBstringType) String() string { return "timeType" }

func (_optBooleanType) typeCode() typeCode { return tcBoolean }
func (_optTinyintType) typeCode() typeCode { return tcTinyint }
func (_optIntegerType) typeCode() typeCode { return tcInteger }
func (_optBigintType) typeCode() typeCode  { return tcBigint }
func (_optDoubleType) typeCode() typeCode  { return tcDouble }
func (_optStringType) typeCode() typeCode  { return tcString }
func (_optBstringType) typeCode() typeCode { return tcBstring }

func (_optBooleanType) size(any) int   { return encoding.BooleanFieldSize }
func (_optTinyintType) size(any) int   { return encoding.TinyintFieldSize }
func (_optIntegerType) size(any) int   { return encoding.IntegerFieldSize }
func (_optBigintType) size(any) int    { return encoding.BigintFieldSize }
func (_optDoubleType) size(any) int    { return encoding.DoubleFieldSize }
func (_optStringType) size(v any) int  { return 2 + len(v.(string)) } // length int16 + string length
func (_optBstringType) size(v any) int { return 2 + len(v.([]byte)) } // length int16 + bytes length

func (_optBooleanType) encode(e *encoding.Encoder, v any) { e.Bool(v.(bool)) }
func (_optTinyintType) encode(e *encoding.Encoder, v any) { e.Int8(v.(int8)) }
func (_optIntegerType) encode(e *encoding.Encoder, v any) { e.Int32(v.(int32)) }
func (_optBigintType) encode(e *encoding.Encoder, v any)  { e.Int64(v.(int64)) }
func (_optDoubleType) encode(e *encoding.Encoder, v any)  { e.Float64(v.(float64)) }
func (_optStringType) encode(e *encoding.Encoder, v any) {
	s := v.(string)
	e.Int16(int16(len(s))) //nolint: gosec
	e.Bytes([]byte(s))
}
func (_optBstringType) encode(e *encoding.Encoder, v any) {
	b := v.([]byte)
	e.Int16(int16(len(b))) //nolint: gosec
	e.Bytes(b)
}

func (_optBooleanType) decode(d *encoding.Decoder) any { return d.Bool() }
func (_optTinyintType) decode(d *encoding.Decoder) any { return d.Int8() }
func (_optIntegerType) decode(d *encoding.Decoder) any { return d.Int32() }
func (_optBigintType) decode(d *encoding.Decoder) any  { return d.Int64() }
func (_optDoubleType) decode(d *encoding.Decoder) any  { return d.Float64() }
func (_optStringType) decode(d *encoding.Decoder) any {
	l := d.Int16()
	b := make([]byte, l)
	d.Bytes(b)
	return string(b)
}
func (_optBstringType) decode(d *encoding.Decoder) any {
	l := d.Int16()
	b := make([]byte, l)
	d.Bytes(b)
	return b
}

func optTypeViaType(v any) optType {
	switch v.(type) {
	case bool:
		return optBooleanType
	case int8:
		return optTinyintType
	case int32:
		return optIntegerType
	case int64:
		return optBigintType
	case float64:
		return optDoubleType
	case string:
		return optStringType
	case []byte:
		return optBstringType
	default:
		panic("type not implemented") // should never happen
	}
}

func optTypeViaTypeCode(tc typeCode) optType {
	switch tc {
	case tcBoolean:
		return optBooleanType
	case tcTinyint:
		return optTinyintType
	case tcInteger:
		return optIntegerType
	case tcBigint:
		return optBigintType
	case tcDouble:
		return optDoubleType
	case tcString:
		return optStringType
	case tcBstring:
		return optBstringType
	default:
		panic("missing optType for typeCode")
	}
}
