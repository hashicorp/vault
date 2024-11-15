package protocol

import (
	"bytes"
	"errors"
	"io"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/transform"
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

var (
	timeReflectType   = reflect.TypeFor[time.Time]()
	bytesReflectType  = reflect.TypeFor[[]byte]()
	stringReflectType = reflect.TypeFor[string]()
	ratReflectType    = reflect.TypeFor[big.Rat]()
)

var (
	errConversionNotSupported = errors.New("conversion not supported")
	errUint64OutOfRange       = errors.New("uint64 values with high bit set are not supported")
	errIntegerOutOfRange      = errors.New("integer out of range")
	errFloatOutOfRange        = errors.New("float out of range")
)

/*
Conversion routines hdb parameters
  - return value is any to avoid allocations in case
    parameter is already of target type
*/

func convertBool(v any) (any, error) {
	// check needs to be done on each type individually as if combining types in one case
	// the v type stays on any and the comparison v != 0 would always be true.
	switch v := v.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0, nil
	case reflect.String:
		return strconv.ParseBool(rv.String())
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertBool(rv.Elem().Interface())
	default:
		if rv.Type().ConvertibleTo(stringReflectType) {
			return convertBool(rv.Convert(stringReflectType).Interface())
		}
		return nil, errConversionNotSupported
	}
}

var (
	i64Zero = int64(0)
	i64One  = int64(1)
)

func convertInteger(v any, minI64, maxI64 int64) (any, error) { //nolint: gocyclo
	switch v := v.(type) {
	case bool:
		if v {
			return i64One, nil
		}
		return i64Zero, nil
	case int:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case int8:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case int16:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case int32:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case int64:
		if v > maxI64 || v < minI64 {
			return nil, errIntegerOutOfRange
		}
		return v, nil
	case uint:
		u64 := uint64(v)
		if u64 > math.MaxInt64 {
			return nil, errUint64OutOfRange
		}
		i64 := int64(u64)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case uint8:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case uint16:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case uint32:
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case uint64:
		if v > math.MaxInt64 {
			return nil, errUint64OutOfRange
		}
		i64 := int64(v)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case float32:
		i64 := int64(v)
		if v != float32(i64) { // should work for overflow, NaN, +-INF as well
			return nil, errConversionNotSupported
		}
		if i64 > maxI64 || i64 < minI64 {
			return nil, errConversionNotSupported
		}
		return i64, nil
	case float64:
		i64 := int64(v)
		if v != float64(i64) { // should work for overflow, NaN, +-INF as well
			return nil, errConversionNotSupported
		}
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case string:
		i64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return i64One, nil
		}
		return i64Zero, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := rv.Int()
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		i64 := int64(rv.Uint()) //nolint: gosec
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case reflect.Uint64:
		u64 := rv.Uint()
		if u64 > math.MaxInt64 {
			return nil, errUint64OutOfRange
		}
		i64 := int64(u64)
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case reflect.Float32, reflect.Float64:
		f64 := rv.Float()
		i64 := int64(f64)
		if f64 != float64(i64) { // should work for overflow, NaN, +-INF as well
			return nil, errConversionNotSupported
		}
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case reflect.String:
		i64, err := strconv.ParseInt(rv.String(), 10, 64)
		if err != nil {
			return nil, errConversionNotSupported
		}
		if i64 > maxI64 || i64 < minI64 {
			return nil, errIntegerOutOfRange
		}
		return i64, nil
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertInteger(rv.Elem().Interface(), minI64, maxI64)
	default:
		if rv.Type().ConvertibleTo(stringReflectType) {
			return convertInteger(rv.Convert(stringReflectType).Interface(), minI64, maxI64)
		}
		return nil, errConversionNotSupported
	}
}

var (
	f64Zero = float64(0.0)
	f64One  = float64(1.0)
)

func convertFloat(v any, maxF64 float64) (any, error) { //nolint: gocyclo
	switch v := v.(type) {
	case float32:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case float64:
		if math.Abs(v) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return v, nil
	case bool:
		if v {
			return f64One, nil
		}
		return f64Zero, nil
	case int:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case int8:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case int16:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case int32:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case int64:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case uint:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case uint8:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case uint16:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case uint32:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case uint64:
		f64 := float64(v)
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case string:
		f64, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return f64One, nil
		}
		return f64Zero, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f64 := float64(rv.Int())
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f64 := float64(rv.Uint())
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case reflect.Float32, reflect.Float64:
		f64 := rv.Float()
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case reflect.String:
		f64, err := strconv.ParseFloat(rv.String(), 64)
		if err != nil {
			return nil, err
		}
		if math.Abs(f64) > maxF64 {
			return nil, errFloatOutOfRange
		}
		return f64, nil
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertFloat(rv.Elem().Interface(), maxF64)
	default:
		if rv.Type().ConvertibleTo(stringReflectType) {
			return convertFloat(rv.Convert(stringReflectType).Interface(), maxF64)
		}
		return nil, errConversionNotSupported
	}
}

func convertTime(v any) (any, error) {
	if v, ok := v.(time.Time); ok {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertTime(rv.Elem().Interface())
	default:
		if rv.Type().ConvertibleTo(timeReflectType) {
			tv := rv.Convert(timeReflectType)
			return tv.Interface().(time.Time), nil
		}
		return nil, errConversionNotSupported
	}
}

var (
	ratZero = big.NewRat(0, 1)
	ratOne  = big.NewRat(1, 1)
)

/*
Currently the min, max check is done during encoding, as the check is expensive and
we want to avoid doing the conversion twice (convert + encode).
These checks could be done in convert only, but then we would need a
struct{m *big.Int, exp int} for decimals as intermediate format.

The conversion does support other types as well (int, *big.Int, string, ...)
even though the user needs to use Decimal for scanning.
*/
func convertDecimal(v any) (any, error) { //nolint: gocyclo
	switch v := v.(type) {
	case *big.Rat:
		return v, nil
	case *big.Int:
		return new(big.Rat).SetInt(v), nil
	case *big.Float:
		r, _ := v.Rat(nil) // ignore accuracy
		return r, nil
	case bool:
		if v {
			return ratOne, nil
		}
		return ratZero, nil
	case int:
		return new(big.Rat).SetInt64(int64(v)), nil
	case int8:
		return new(big.Rat).SetInt64(int64(v)), nil
	case int16:
		return new(big.Rat).SetInt64(int64(v)), nil
	case int32:
		return new(big.Rat).SetInt64(int64(v)), nil
	case int64:
		return new(big.Rat).SetInt64(v), nil
	case uint:
		return new(big.Rat).SetUint64(uint64(v)), nil
	case uint8:
		return new(big.Rat).SetUint64(uint64(v)), nil
	case uint16:
		return new(big.Rat).SetUint64(uint64(v)), nil
	case uint32:
		return new(big.Rat).SetUint64(uint64(v)), nil
	case uint64:
		return new(big.Rat).SetUint64(v), nil
	case float32:
		r := new(big.Rat).SetFloat64(float64(v))
		if r == nil {
			return nil, errConversionNotSupported
		}
	case float64:
		r := new(big.Rat).SetFloat64(v)
		if r == nil {
			return nil, errConversionNotSupported
		}
	case string:
		r, ok := new(big.Rat).SetString(v)
		if !ok {
			return nil, errConversionNotSupported
		}
		return r, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return ratOne, nil
		}
		return ratZero, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return new(big.Rat).SetInt64(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return new(big.Rat).SetUint64(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		r := new(big.Rat).SetFloat64(rv.Float())
		if r == nil {
			return nil, errConversionNotSupported
		}
		return r, nil
	case reflect.String:
		r, ok := new(big.Rat).SetString(rv.String())
		if !ok {
			return nil, errConversionNotSupported
		}
		return r, nil
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertDecimal(rv.Elem().Interface())
	default:
		if rv.Type().ConvertibleTo(ratReflectType) {
			tv := rv.Convert(ratReflectType)
			return tv.Interface().(big.Rat), nil
		}
		return nil, errConversionNotSupported
	}
}

func convertBytes(v any) (any, error) {
	switch v := v.(type) {
	case string, []byte:
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertBytes(rv.Elem().Interface())
	case reflect.Slice:
		if rv.Type() == bytesReflectType {
			return rv.Bytes(), nil
		}
		fallthrough
	default:
		if rv.Type().ConvertibleTo(bytesReflectType) {
			bv := rv.Convert(bytesReflectType)
			return bv.Interface().([]byte), nil
		}
		return nil, errConversionNotSupported
	}
}

// readProvider is the interface wrapping the Reader which provides an io.Reader.
type readProvider interface {
	Reader() io.Reader
}

func convertToLobInDescr(tr transform.Transformer, rd io.Reader) *LobInDescr {
	return newLobInDescr(tr, rd)
}

func convertLob(v any, t transform.Transformer) (any, error) {
	switch v := v.(type) {
	case io.Reader:
		return convertToLobInDescr(t, v), nil
	case readProvider:
		return convertToLobInDescr(t, v.Reader()), nil
	default:
		// check if string or []byte
		if v, err := convertBytes(v); err == nil {
			switch v := v.(type) {
			case string:
				return convertToLobInDescr(t, strings.NewReader(v)), nil
			case []byte:
				return convertToLobInDescr(t, bytes.NewReader(v)), nil
			}
		}
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return convertLob(rv.Elem().Interface(), t)
	default:
		return nil, errConversionNotSupported
	}
}

func convertField(tc typeCode, v any, t transform.Transformer) (any, error) {
	if v == nil {
		return nil, nil
	}

	switch tc {
	case tcBoolean:
		return convertBool(v)
	case tcTinyint:
		return convertInteger(v, minTinyint, maxTinyint)
	case tcSmallint:
		return convertInteger(v, minSmallint, maxSmallint)
	case tcInteger:
		return convertInteger(v, minInteger, maxInteger)
	case tcBigint:
		return convertInteger(v, minBigint, maxBigint)
	case tcReal:
		return convertFloat(v, maxReal)
	case tcDouble:
		return convertFloat(v, maxDouble)
	case tcDate, tcTime, tcTimestamp, tcLongdate, tcSeconddate, tcDaydate, tcSecondtime:
		return convertTime(v)
	case tcDecimal, tcFixed8, tcFixed12, tcFixed16:
		return convertDecimal(v)
	case tcChar, tcVarchar, tcString, tcAlphanum, tcNchar, tcNvarchar, tcNstring, tcShorttext, tcBinary, tcVarbinary, tcStPoint, tcStGeometry:
		return convertBytes(v)
	case tcBlob, tcClob, tcLocator:
		return convertLob(v, nil)
	case tcNclob, tcText, tcNlocator:
		return convertLob(v, t)
	case tcBintext: // ?? lobCESU8Type
		return convertLob(v, nil)
	default:
		panic("invalid type code")
	}
}
