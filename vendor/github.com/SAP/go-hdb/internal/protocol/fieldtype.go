// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"database/sql/driver"
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"
	"time"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
	"github.com/SAP/go-hdb/internal/unicode/cesu8"
)

const (
	minTinyint  = 0
	maxTinyint  = math.MaxUint8
	minSmallint = math.MinInt16
	maxSmallint = math.MaxInt16
	minInteger  = math.MinInt32
	maxInteger  = math.MaxInt32
	minBigint   = math.MinInt64
	maxBigint   = math.MaxInt64
	maxReal     = math.MaxFloat32
	maxDouble   = math.MaxFloat64
)

// string / binary length indicators
const (
	bytesLenIndNullValue byte = 255
	bytesLenIndSmall     byte = 245
	bytesLenIndMedium    byte = 246
	bytesLenIndBig       byte = 247
)

const (
	realNullValue   uint32 = ^uint32(0)
	doubleNullValue uint64 = ^uint64(0)
)

const (
	booleanFalseValue   byte  = 0
	booleanNullValue    byte  = 1
	booleanTrueValue    byte  = 2
	longdateNullValue   int64 = 3155380704000000001
	seconddateNullValue int64 = 315538070401
	daydateNullValue    int32 = 3652062
	secondtimeNullValue int32 = 86402
)

type locatorID uint64 // byte[locatorIdSize]

var timeReflectType = reflect.TypeOf((*time.Time)(nil)).Elem()
var bytesReflectType = reflect.TypeOf((*[]byte)(nil)).Elem()
var stringReflectType = reflect.TypeOf((*string)(nil)).Elem()

var zeroTime = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

const (
	booleanFieldSize    = 1
	tinyintFieldSize    = 1
	smallintFieldSize   = 2
	integerFieldSize    = 4
	bigintFieldSize     = 8
	realFieldSize       = 4
	doubleFieldSize     = 8
	dateFieldSize       = 4
	timeFieldSize       = 4
	timestampFieldSize  = dateFieldSize + timeFieldSize
	longdateFieldSize   = 8
	seconddateFieldSize = 8
	daydateFieldSize    = 4
	secondtimeFieldSize = 4
	decimalFieldSize    = 16
	fixed8FieldSize     = 8
	fixed12FieldSize    = 12
	fixed16FieldSize    = 16

	lobInputParametersSize = 9
)

type fieldType interface {
	/*
		statements:
		- first parameter could be many
		- so the check needs to 'fail fast'
		- fmt.Errorf is too slow because contructor formats the error -> use ConvertError
	*/
	convert(v interface{}) (interface{}, error)
	prmSize(v interface{}) int
	encodePrm(e *encoding.Encoder, v interface{}) error
	decodeRes(d *encoding.Decoder) (interface{}, error)
	decodePrm(d *encoding.Decoder) (interface{}, error)
}

var (
	booleanType    = _booleanType{}
	tinyintType    = _tinyintType{}
	smallintType   = _smallintType{}
	integerType    = _integerType{}
	bigintType     = _bigintType{}
	realType       = _realType{}
	doubleType     = _doubleType{}
	dateType       = _dateType{}
	timeType       = _timeType{}
	timestampType  = _timestampType{}
	longdateType   = _longdateType{}
	seconddateType = _seconddateType{}
	daydateType    = _daydateType{}
	secondtimeType = _secondtimeType{}
	decimalType    = _decimalType{}
	varType        = _varType{}
	alphaType      = _alphaType{}
	cesu8Type      = _cesu8Type{}
	lobVarType     = _lobVarType{}
	lobCESU8Type   = _lobCESU8Type{}
)

type _booleanType struct{}
type _tinyintType struct{}
type _smallintType struct{}
type _integerType struct{}
type _bigintType struct{}
type _realType struct{}
type _doubleType struct{}
type _dateType struct{}
type _timeType struct{}
type _timestampType struct{}
type _longdateType struct{}
type _seconddateType struct{}
type _daydateType struct{}
type _secondtimeType struct{}
type _decimalType struct{}
type _fixed8Type struct{ prec, scale int }
type _fixed12Type struct{ prec, scale int }
type _fixed16Type struct{ prec, scale int }
type _varType struct{}
type _alphaType struct{}
type _cesu8Type struct{}
type _lobVarType struct{}
type _lobCESU8Type struct{}

var (
	_ fieldType = (*_booleanType)(nil)
	_ fieldType = (*_tinyintType)(nil)
	_ fieldType = (*_smallintType)(nil)
	_ fieldType = (*_integerType)(nil)
	_ fieldType = (*_bigintType)(nil)
	_ fieldType = (*_realType)(nil)
	_ fieldType = (*_doubleType)(nil)
	_ fieldType = (*_dateType)(nil)
	_ fieldType = (*_timeType)(nil)
	_ fieldType = (*_timestampType)(nil)
	_ fieldType = (*_longdateType)(nil)
	_ fieldType = (*_seconddateType)(nil)
	_ fieldType = (*_daydateType)(nil)
	_ fieldType = (*_secondtimeType)(nil)
	_ fieldType = (*_decimalType)(nil)
	_ fieldType = (*_fixed8Type)(nil)
	_ fieldType = (*_fixed12Type)(nil)
	_ fieldType = (*_fixed16Type)(nil)
	_ fieldType = (*_varType)(nil)
	_ fieldType = (*_alphaType)(nil)
	_ fieldType = (*_cesu8Type)(nil)
	_ fieldType = (*_lobVarType)(nil)
	_ fieldType = (*_lobCESU8Type)(nil)
)

// stringer
func (_booleanType) String() string    { return "booleanType" }
func (_tinyintType) String() string    { return "tinyintType" }
func (_smallintType) String() string   { return "smallintType" }
func (_integerType) String() string    { return "integerType" }
func (_bigintType) String() string     { return "bigintType" }
func (_realType) String() string       { return "realType" }
func (_doubleType) String() string     { return "doubleType" }
func (_dateType) String() string       { return "dateType" }
func (_timeType) String() string       { return "timeType" }
func (_timestampType) String() string  { return "timestampType" }
func (_longdateType) String() string   { return "longdateType" }
func (_seconddateType) String() string { return "seconddateType" }
func (_daydateType) String() string    { return "daydateType" }
func (_secondtimeType) String() string { return "secondtimeType" }
func (_decimalType) String() string    { return "decimalType" }
func (_fixed8Type) String() string     { return "fixed8Type" }
func (_fixed12Type) String() string    { return "fixed12Type" }
func (_fixed16Type) String() string    { return "fixed16Type" }
func (_varType) String() string        { return "varType" }
func (_alphaType) String() string      { return "alphaType" }
func (_cesu8Type) String() string      { return "cesu8Type" }
func (_lobVarType) String() string     { return "lobVarType" }
func (_lobCESU8Type) String() string   { return "lobCESU8Type" }

// convert
func (ft _booleanType) convert(v interface{}) (interface{}, error) { return convertBool(ft, v) }

func (ft _tinyintType) convert(v interface{}) (interface{}, error) {
	return convertInteger(ft, v, minTinyint, maxTinyint)
}
func (ft _smallintType) convert(v interface{}) (interface{}, error) {
	return convertInteger(ft, v, minSmallint, maxSmallint)
}
func (ft _integerType) convert(v interface{}) (interface{}, error) {
	return convertInteger(ft, v, minInteger, maxInteger)
}
func (ft _bigintType) convert(v interface{}) (interface{}, error) {
	return convertInteger(ft, v, minBigint, maxBigint)
}

func (ft _realType) convert(v interface{}) (interface{}, error) {
	return convertFloat(ft, v, maxReal)
}
func (ft _doubleType) convert(v interface{}) (interface{}, error) {
	return convertFloat(ft, v, maxDouble)
}

func (ft _dateType) convert(v interface{}) (interface{}, error) { return convertTime(ft, v) }
func (ft _timeType) convert(v interface{}) (interface{}, error) { return convertTime(ft, v) }
func (ft _timestampType) convert(v interface{}) (interface{}, error) {
	return convertTime(ft, v)
}
func (ft _longdateType) convert(v interface{}) (interface{}, error)   { return convertTime(ft, v) }
func (ft _seconddateType) convert(v interface{}) (interface{}, error) { return convertTime(ft, v) }
func (ft _daydateType) convert(v interface{}) (interface{}, error)    { return convertTime(ft, v) }
func (ft _secondtimeType) convert(v interface{}) (interface{}, error) {
	return convertTime(ft, v)
}

func (ft _decimalType) convert(v interface{}) (interface{}, error) { return convertDecimal(ft, v) }
func (ft _fixed8Type) convert(v interface{}) (interface{}, error)  { return convertDecimal(ft, v) }
func (ft _fixed12Type) convert(v interface{}) (interface{}, error) { return convertDecimal(ft, v) }
func (ft _fixed16Type) convert(v interface{}) (interface{}, error) { return convertDecimal(ft, v) }

func (ft _varType) convert(v interface{}) (interface{}, error)   { return convertBytes(ft, v) }
func (ft _alphaType) convert(v interface{}) (interface{}, error) { return convertBytes(ft, v) }
func (ft _cesu8Type) convert(v interface{}) (interface{}, error) { return convertBytes(ft, v) }

func (ft _lobVarType) convert(v interface{}) (interface{}, error)   { return convertLob(ft, v) }
func (ft _lobCESU8Type) convert(v interface{}) (interface{}, error) { return convertLob(ft, v) }

// ReadProvider is the interface wrapping the Reader which provides an io.Reader.
type ReadProvider interface {
	Reader() io.Reader
}

// Lob
func convertLob(ft fieldType, v interface{}) (driver.Value, error) {
	if v == nil {
		return v, nil
	}

	switch v := v.(type) {
	case io.Reader:
		return v, nil
	case ReadProvider:
		return v.Reader(), nil
	default:
		return nil, newConvertError(ft, v, nil)
	}
}

// prm size
func (_booleanType) prmSize(interface{}) int    { return booleanFieldSize }
func (_tinyintType) prmSize(interface{}) int    { return tinyintFieldSize }
func (_smallintType) prmSize(interface{}) int   { return smallintFieldSize }
func (_integerType) prmSize(interface{}) int    { return integerFieldSize }
func (_bigintType) prmSize(interface{}) int     { return bigintFieldSize }
func (_realType) prmSize(interface{}) int       { return realFieldSize }
func (_doubleType) prmSize(interface{}) int     { return doubleFieldSize }
func (_dateType) prmSize(interface{}) int       { return dateFieldSize }
func (_timeType) prmSize(interface{}) int       { return timeFieldSize }
func (_timestampType) prmSize(interface{}) int  { return timestampFieldSize }
func (_longdateType) prmSize(interface{}) int   { return longdateFieldSize }
func (_seconddateType) prmSize(interface{}) int { return seconddateFieldSize }
func (_daydateType) prmSize(interface{}) int    { return daydateFieldSize }
func (_secondtimeType) prmSize(interface{}) int { return secondtimeFieldSize }
func (_decimalType) prmSize(interface{}) int    { return decimalFieldSize }
func (_fixed8Type) prmSize(interface{}) int     { return fixed8FieldSize }
func (_fixed12Type) prmSize(interface{}) int    { return fixed12FieldSize }
func (_fixed16Type) prmSize(interface{}) int    { return fixed16FieldSize }
func (_lobVarType) prmSize(v interface{}) int   { return lobInputParametersSize }
func (_lobCESU8Type) prmSize(v interface{}) int { return lobInputParametersSize }

func (ft _varType) prmSize(v interface{}) int {
	switch v := v.(type) {
	case []byte:
		return varBytesSize(len(v))
	case string:
		return varBytesSize(len(v))
	default:
		return -1
	}
}
func (ft _alphaType) prmSize(v interface{}) int {
	return varType.prmSize(v)
}
func (ft _cesu8Type) prmSize(v interface{}) int {
	switch v := v.(type) {
	case []byte:
		return varBytesSize(cesu8.Size(v))
	case string:
		return varBytesSize(cesu8.StringSize(v))
	default:
		return -1
	}
}

func varBytesSize(size int) int {
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

// encode
func (ft _booleanType) encodePrm(e *encoding.Encoder, v interface{}) error {
	if v == nil {
		e.Byte(booleanNullValue)
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		panic("invalid bool value") // should never happen
	}
	if b {
		e.Byte(booleanTrueValue)
	} else {
		e.Byte(booleanFalseValue)
	}
	return nil
}

func (ft _tinyintType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Byte(byte(asInt64(v)))
	return nil
}
func (ft _smallintType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Int16(int16(asInt64(v)))
	return nil
}
func (ft _integerType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Int32(int32(asInt64(v)))
	return nil
}
func (ft _bigintType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Int64(asInt64(v))
	return nil
}

func asInt64(v interface{}) int64 {
	switch v := v.(type) {
	case bool:
		if v {
			return 1
		}
		return 0
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint())
	default:
		panic("invalid bool value") // should never happen
	}
}

func (ft _realType) encodePrm(e *encoding.Encoder, v interface{}) error {
	switch v := v.(type) {
	case float32:
		e.Float32(v)
	case float64:
		e.Float32(float32(v))
	default:
		panic("invalid real value") // should never happen
	}
	return nil
}

func (ft _doubleType) encodePrm(e *encoding.Encoder, v interface{}) error {
	switch v := v.(type) {
	case float32:
		e.Float64(float64(v))
	case float64:
		e.Float64(v)
	default:
		panic("invalid double value") // should never happen
	}
	return nil
}

func (ft _dateType) encodePrm(e *encoding.Encoder, v interface{}) error {
	encodeDate(e, asTime(v))
	return nil
}
func (ft _timeType) encodePrm(e *encoding.Encoder, v interface{}) error {
	encodeTime(e, asTime(v))
	return nil
}
func (ft _timestampType) encodePrm(e *encoding.Encoder, v interface{}) error {
	t := asTime(v)
	encodeDate(e, t)
	encodeTime(e, t)
	return nil
}

func encodeDate(e *encoding.Encoder, t time.Time) {
	// year: set most sig bit
	// month 0 based
	year, month, day := t.Date()
	e.Uint16(uint16(year) | 0x8000)
	e.Int8(int8(month) - 1)
	e.Int8(int8(day))
}

func encodeTime(e *encoding.Encoder, t time.Time) {
	e.Byte(byte(t.Hour()) | 0x80)
	e.Int8(int8(t.Minute()))
	msec := t.Second()*1000 + t.Nanosecond()/1000000
	e.Uint16(uint16(msec))
}

func (ft _longdateType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Int64(convertTimeToLongdate(asTime(v)))
	return nil
}
func (ft _seconddateType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Int64(convertTimeToSeconddate(asTime(v)))
	return nil
}
func (ft _daydateType) encodePrm(e *encoding.Encoder, v interface{}) error {
	e.Int32(int32(convertTimeToDayDate(asTime(v))))
	return nil
}
func (ft _secondtimeType) encodePrm(e *encoding.Encoder, v interface{}) error {
	if v == nil {
		e.Int32(secondtimeNullValue)
		return nil
	}
	e.Int32(int32(convertTimeToSecondtime(asTime(v))))
	return nil
}

func asTime(v interface{}) time.Time {
	t, ok := v.(time.Time)
	if !ok {
		panic("invalid time value") // should never happen
	}
	//store in utc
	return t.UTC()
}

func (ft _decimalType) encodePrm(e *encoding.Encoder, v interface{}) error {
	r, ok := v.(*big.Rat)
	if !ok {
		panic("invalid decimal value") // should never happen
	}

	var m big.Int
	exp, df := convertRatToDecimal(r, &m, dec128Digits, dec128MinExp, dec128MaxExp)

	if df&dfOverflow != 0 {
		return ErrDecimalOutOfRange
	}

	if df&dfUnderflow != 0 { // set to zero
		e.Decimal(natZero, 0)
	} else {
		e.Decimal(&m, exp)
	}
	return nil
}

func (ft _fixed8Type) encodePrm(e *encoding.Encoder, v interface{}) error {
	return encodeFixed(e, v, fixed8FieldSize, ft.prec, ft.scale)
}

func (ft _fixed12Type) encodePrm(e *encoding.Encoder, v interface{}) error {
	return encodeFixed(e, v, fixed12FieldSize, ft.prec, ft.scale)
}

func (ft _fixed16Type) encodePrm(e *encoding.Encoder, v interface{}) error {
	return encodeFixed(e, v, fixed16FieldSize, ft.prec, ft.scale)
}

func encodeFixed(e *encoding.Encoder, v interface{}, size, prec, scale int) error {
	r, ok := v.(*big.Rat)
	if !ok {
		panic("invalid decimal value") // should never happen
	}

	var m big.Int
	df := convertRatToFixed(r, &m, prec, scale)

	if df&dfOverflow != 0 {
		return ErrDecimalOutOfRange
	}

	e.Fixed(&m, size)
	return nil
}

func (ft _varType) encodePrm(e *encoding.Encoder, v interface{}) error {
	switch v := v.(type) {
	case []byte:
		return encodeVarBytes(e, v)
	case string:
		return encodeVarString(e, v)
	default:
		panic("invalid var value") // should never happen
	}
}
func (ft _alphaType) encodePrm(e *encoding.Encoder, v interface{}) error {
	return varType.encodePrm(e, v)
}
func encodeVarBytesSize(e *encoding.Encoder, size int) error {
	switch {
	default:
		return fmt.Errorf("max argument length %d of string exceeded", size)
	case size <= int(bytesLenIndSmall):
		e.Byte(byte(size))
	case size <= math.MaxInt16:
		e.Byte(bytesLenIndMedium)
		e.Int16(int16(size))
	case size <= math.MaxInt32:
		e.Byte(bytesLenIndBig)
		e.Int32(int32(size))
	}
	return nil
}

func encodeVarBytes(e *encoding.Encoder, p []byte) error {
	if err := encodeVarBytesSize(e, len(p)); err != nil {
		return err
	}
	e.Bytes(p)
	return nil
}

func encodeVarString(e *encoding.Encoder, s string) error {
	if err := encodeVarBytesSize(e, len(s)); err != nil {
		return err
	}
	e.String(s)
	return nil
}

func (ft _cesu8Type) encodePrm(e *encoding.Encoder, v interface{}) error {
	switch v := v.(type) {
	case []byte:
		return encodeCESU8Bytes(e, v)
	case string:
		return encodeCESU8String(e, v)
	default:
		panic("invalid cesu8 value") // should never happen
	}
}

func encodeCESU8Bytes(e *encoding.Encoder, p []byte) error {
	size := cesu8.Size(p)
	if err := encodeVarBytesSize(e, size); err != nil {
		return err
	}
	e.CESU8Bytes(p)
	return nil
}

func encodeCESU8String(e *encoding.Encoder, s string) error {
	size := cesu8.StringSize(s)
	if err := encodeVarBytesSize(e, size); err != nil {
		return err
	}
	e.CESU8String(s)
	return nil
}

func (ft _lobVarType) encodePrm(e *encoding.Encoder, v interface{}) error {
	switch v := v.(type) {
	case *lobInDescr:
		return encodeLobPrm(e, v)
	case io.Reader: //TODO check if keep
		descr := &lobInDescr{}
		return encodeLobPrm(e, descr)
	default:
		panic("invalid lob var value") // should never happen
	}
}

func (ft _lobCESU8Type) encodePrm(e *encoding.Encoder, v interface{}) error {
	switch v := v.(type) {
	case *lobInDescr:
		return encodeLobPrm(e, v)
	case io.Reader: //TODO check if keep
		descr := &lobInDescr{}
		return encodeLobPrm(e, descr)
	default:
		panic("invalid lob cesu8 value") // should never happen
	}
}

func encodeLobPrm(e *encoding.Encoder, descr *lobInDescr) error {
	e.Byte(byte(descr.opt))
	e.Int32(descr.size)
	e.Int32(descr.pos)
	return nil
}

// field types for which decodePrm is same as decodeRes
func (ft _booleanType) decodePrm(d *encoding.Decoder) (interface{}, error)    { return ft.decodeRes(d) }
func (ft _realType) decodePrm(d *encoding.Decoder) (interface{}, error)       { return ft.decodeRes(d) }
func (ft _doubleType) decodePrm(d *encoding.Decoder) (interface{}, error)     { return ft.decodeRes(d) }
func (ft _dateType) decodePrm(d *encoding.Decoder) (interface{}, error)       { return ft.decodeRes(d) }
func (ft _timeType) decodePrm(d *encoding.Decoder) (interface{}, error)       { return ft.decodeRes(d) }
func (ft _timestampType) decodePrm(d *encoding.Decoder) (interface{}, error)  { return ft.decodeRes(d) }
func (ft _longdateType) decodePrm(d *encoding.Decoder) (interface{}, error)   { return ft.decodeRes(d) }
func (ft _seconddateType) decodePrm(d *encoding.Decoder) (interface{}, error) { return ft.decodeRes(d) }
func (ft _daydateType) decodePrm(d *encoding.Decoder) (interface{}, error)    { return ft.decodeRes(d) }
func (ft _secondtimeType) decodePrm(d *encoding.Decoder) (interface{}, error) { return ft.decodeRes(d) }
func (ft _decimalType) decodePrm(d *encoding.Decoder) (interface{}, error)    { return ft.decodeRes(d) }
func (ft _fixed8Type) decodePrm(d *encoding.Decoder) (interface{}, error)     { return ft.decodeRes(d) }
func (ft _fixed12Type) decodePrm(d *encoding.Decoder) (interface{}, error)    { return ft.decodeRes(d) }
func (ft _fixed16Type) decodePrm(d *encoding.Decoder) (interface{}, error)    { return ft.decodeRes(d) }
func (ft _varType) decodePrm(d *encoding.Decoder) (interface{}, error)        { return ft.decodeRes(d) }
func (ft _alphaType) decodePrm(d *encoding.Decoder) (interface{}, error)      { return ft.decodeRes(d) }
func (ft _cesu8Type) decodePrm(d *encoding.Decoder) (interface{}, error)      { return ft.decodeRes(d) }

// decode
func (_booleanType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	b := d.Byte()
	switch b {
	case booleanNullValue:
		return nil, nil
	case booleanFalseValue:
		return false, nil
	default:
		return true, nil
	}
}

func (_tinyintType) decodePrm(d *encoding.Decoder) (interface{}, error) { return int64(d.Byte()), nil }
func (_smallintType) decodePrm(d *encoding.Decoder) (interface{}, error) {
	return int64(d.Int16()), nil
}
func (_integerType) decodePrm(d *encoding.Decoder) (interface{}, error) { return int64(d.Int32()), nil }
func (_bigintType) decodePrm(d *encoding.Decoder) (interface{}, error)  { return d.Int64(), nil }

func (ft _tinyintType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return ft.decodePrm(d)
}
func (ft _smallintType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return ft.decodePrm(d)
}
func (ft _integerType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return ft.decodePrm(d)
}
func (ft _bigintType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return ft.decodePrm(d)
}

func (_realType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	v := d.Uint32()
	if v == realNullValue {
		return nil, nil
	}
	return float64(math.Float32frombits(v)), nil
}
func (_doubleType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	v := d.Uint64()
	if v == doubleNullValue {
		return nil, nil
	}
	return math.Float64frombits(v), nil
}

func (_dateType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	year, month, day, null := decodeDate(d)
	if null {
		return nil, nil
	}
	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC), nil
}
func (_timeType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	// time read gives only seconds (cut), no milliseconds
	hour, min, sec, nsec, null := decodeTime(d)
	if null {
		return nil, nil
	}
	return time.Date(1, 1, 1, hour, min, sec, nsec, time.UTC), nil
}
func (_timestampType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	year, month, day, dateNull := decodeDate(d)
	hour, min, sec, nsec, timeNull := decodeTime(d)
	if dateNull || timeNull {
		return nil, nil
	}
	return time.Date(year, month, day, hour, min, sec, nsec, time.UTC), nil
}

// null values: most sig bit unset
// year: unset second most sig bit (subtract 2^15)
// --> read year as unsigned
// month is 0-based
// day is 1 byte
func decodeDate(d *encoding.Decoder) (int, time.Month, int, bool) {
	year := d.Uint16()
	null := ((year & 0x8000) == 0) //null value
	year &= 0x3fff
	month := d.Int8()
	month++
	day := d.Int8()
	return int(year), time.Month(month), int(day), null
}

func decodeTime(d *encoding.Decoder) (int, int, int, int, bool) {
	hour := d.Byte()
	null := (hour & 0x80) == 0 //null value
	hour &= 0x7f
	min := d.Int8()
	msec := d.Uint16()

	sec := msec / 1000
	msec %= 1000
	nsec := int(msec) * 1000000

	return int(hour), int(min), int(sec), nsec, null
}

func (_longdateType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	longdate := d.Int64()
	if longdate == longdateNullValue {
		return nil, nil
	}
	return convertLongdateToTime(longdate), nil
}
func (_seconddateType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	seconddate := d.Int64()
	if seconddate == seconddateNullValue {
		return nil, nil
	}
	return convertSeconddateToTime(seconddate), nil
}
func (_daydateType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	daydate := d.Int32()
	if daydate == daydateNullValue {
		return nil, nil
	}
	return convertDaydateToTime(int64(daydate)), nil
}
func (_secondtimeType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	secondtime := d.Int32()
	if secondtime == secondtimeNullValue {
		return nil, nil
	}
	return convertSecondtimeToTime(int(secondtime)), nil
}

func (_decimalType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	m, exp := d.Decimal()
	if m == nil {
		return nil, nil
	}
	return convertDecimalToRat(m, exp), nil
}

func (ft _fixed8Type) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return decodeFixed(d, fixed8FieldSize, ft.prec, ft.scale)
}
func (ft _fixed12Type) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return decodeFixed(d, fixed12FieldSize, ft.prec, ft.scale)
}
func (ft _fixed16Type) decodeRes(d *encoding.Decoder) (interface{}, error) {
	if !d.Bool() { //null value
		return nil, nil
	}
	return decodeFixed(d, fixed16FieldSize, ft.prec, ft.scale)
}

func decodeFixed(d *encoding.Decoder, size, prec, scale int) (interface{}, error) {
	m := d.Fixed(size)
	if m == nil {
		return nil, nil
	}
	return convertFixedToRat(m, scale), nil
}

func (_varType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	size, null := decodeVarBytesSize(d)
	if null {
		return nil, nil
	}
	b := make([]byte, size)
	d.Bytes(b)
	return b, nil
}
func (_alphaType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	size, null := decodeVarBytesSize(d)
	if null {
		return nil, nil
	}
	switch d.Dfv() {
	case dfvLevel1: // like _varType
		b := make([]byte, size)
		d.Bytes(b)
		return b, nil
	default:
		/*
			byte:
			- high bit set -> numeric
			- high bit unset -> alpha
			- bits 0-6: field size
		*/
		d.Byte() // ignore for the moment
		b := make([]byte, size-1)
		d.Bytes(b)
		return b, nil
	}
}
func (_cesu8Type) decodeRes(d *encoding.Decoder) (interface{}, error) {
	size, null := decodeVarBytesSize(d)
	if null {
		return nil, nil
	}
	return d.CESU8Bytes(size), nil
}

func decodeVarBytesSize(d *encoding.Decoder) (int, bool) {
	ind := d.Byte() //length indicator
	switch {
	default:
		return 0, false
	case ind == bytesLenIndNullValue:
		return 0, true
	case ind <= bytesLenIndSmall:
		return int(ind), false
	case ind == bytesLenIndMedium:
		return int(d.Int16()), false
	case ind == bytesLenIndBig:
		return int(d.Int32()), false
	}
}

func decodeLobPrm(d *encoding.Decoder) (interface{}, error) {
	descr := &lobInDescr{}
	descr.opt = lobOptions(d.Byte())
	descr.size = d.Int32()
	descr.pos = d.Int32()
	return nil, nil
}

func (_lobVarType) decodePrm(d *encoding.Decoder) (interface{}, error) {
	return decodeLobPrm(d)
}
func (_lobCESU8Type) decodePrm(d *encoding.Decoder) (interface{}, error) {
	return decodeLobPrm(d)
}

func decodeLobRes(d *encoding.Decoder, isCharBased bool) (interface{}, error) {
	descr := &lobOutDescr{isCharBased: isCharBased}
	descr.ltc = lobTypecode(d.Int8())
	descr.opt = lobOptions(d.Int8())
	if descr.opt.isNull() {
		return nil, nil
	}
	d.Skip(2)
	descr.numChar = d.Int64()
	descr.numByte = d.Int64()
	descr.id = locatorID(d.Uint64())
	size := int(d.Int32())
	descr.b = make([]byte, size)
	d.Bytes(descr.b)
	return descr, nil
}

func (_lobVarType) decodeRes(d *encoding.Decoder) (interface{}, error) {
	return decodeLobRes(d, false)
}
func (_lobCESU8Type) decodeRes(d *encoding.Decoder) (interface{}, error) {
	return decodeLobRes(d, true)
}
