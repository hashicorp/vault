// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"time"
)

// ErrUint64OutOfRange means that a uint64 exceeds the size of a int64.
var ErrUint64OutOfRange = errors.New("uint64 values with high bit set are not supported")

// ErrIntegerOutOfRange means that an integer exceeds the size of the hdb integer field.
var ErrIntegerOutOfRange = errors.New("integer out of range error")

// ErrFloatOutOfRange means that a float exceeds the size of the hdb float field.
var ErrFloatOutOfRange = errors.New("float out of range error")

// ErrDecimalOutOfRange means that a big.Rat exceeds the size of hdb decimal fields.
var ErrDecimalOutOfRange = errors.New("decimal out of range error")

// A ConvertError is returned by conversion methods if a go datatype to hdb datatype conversion fails.
type ConvertError struct {
	err error
	ft  fieldType
	v   interface{}
}

func (e *ConvertError) Error() string {
	return fmt.Sprintf("unsupported %[1]s conversion: %[2]T %[2]v", e.ft, e.v)
}

// Unwrap returns the nested error.
func (e *ConvertError) Unwrap() error { return e.err }
func newConvertError(ft fieldType, v interface{}, err error) *ConvertError {
	return &ConvertError{ft: ft, v: v, err: err}
}

/*
Conversion routines hdb parameters
- return value is interface{} to avoid allocations in case
  parameter is already of target type
*/
func convertBool(ft fieldType, v interface{}) (interface{}, error) {
	if v == nil {
		return v, nil
	}

	if v, ok := v.(bool); ok {
		return v, nil
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
		b, err := strconv.ParseBool(rv.String())
		if err != nil {
			return nil, newConvertError(ft, v, err)
		}
		return b, nil
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertBool(ft, rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(stringReflectType) {
		return convertBool(ft, rv.Convert(stringReflectType).Interface())
	}
	return nil, newConvertError(ft, v, nil)
}

func convertInteger(ft fieldType, v interface{}, min, max int64) (interface{}, error) {
	if v == nil {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	// conversions without allocations (return v)
	case reflect.Bool:
		return v, nil // return (no furhter check needed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := rv.Int()
		if i64 > max || i64 < min {
			return nil, newConvertError(ft, v, ErrIntegerOutOfRange)
		}
		return v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := rv.Uint()
		if u64 >= 1<<63 {
			return nil, newConvertError(ft, v, ErrUint64OutOfRange)
		}
		if int64(u64) > max || int64(u64) < min {
			return nil, newConvertError(ft, v, ErrIntegerOutOfRange)
		}
		return v, nil
	// conversions with allocations (return i64)
	case reflect.Float32, reflect.Float64:
		f64 := rv.Float()
		i64 := int64(f64)
		if f64 != float64(i64) { // should work for overflow, NaN, +-INF as well
			return nil, newConvertError(ft, v, nil)
		}
		if i64 > max || i64 < min {
			return nil, newConvertError(ft, v, ErrIntegerOutOfRange)
		}
		return i64, nil
	case reflect.String:
		i64, err := strconv.ParseInt(rv.String(), 10, 64)
		if err != nil {
			return nil, newConvertError(ft, v, err)
		}
		if i64 > max || i64 < min {
			return nil, newConvertError(ft, v, ErrIntegerOutOfRange)
		}
		return i64, nil
	// pointer
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertInteger(ft, rv.Elem().Interface(), min, max)
	}
	// last resort (try via string)
	if rv.Type().ConvertibleTo(stringReflectType) {
		return convertInteger(ft, rv.Convert(stringReflectType).Interface(), min, max)
	}
	return nil, newConvertError(ft, v, nil)
}

func convertFloat(ft fieldType, v interface{}, max float64) (interface{}, error) {
	if v == nil {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	// conversions without allocations (return v)
	case reflect.Float32, reflect.Float64:
		if math.Abs(rv.Float()) > max {
			return nil, newConvertError(ft, v, ErrFloatOutOfRange)
		}
		return v, nil
	// conversions with allocations (return f64)
	case reflect.String:
		f64, err := strconv.ParseFloat(rv.String(), 64)
		if err != nil {
			return nil, newConvertError(ft, v, err)
		}
		if math.Abs(f64) > max {
			return nil, newConvertError(ft, v, ErrFloatOutOfRange)
		}
		return f64, nil
	// pointer
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertFloat(ft, rv.Elem().Interface(), max)
	}
	// last resort (try via string)
	if rv.Type().ConvertibleTo(stringReflectType) {
		return convertFloat(ft, rv.Convert(stringReflectType).Interface(), max)
	}
	return nil, newConvertError(ft, v, nil)
}

func convertTime(ft fieldType, v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	if v, ok := v.(time.Time); ok {
		return v, nil
	}

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertTime(ft, rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(timeReflectType) {
		tv := rv.Convert(timeReflectType)
		return tv.Interface().(time.Time), nil
	}
	return nil, newConvertError(ft, v, nil)
}

/*
Currently the min, max check is done during encoding, as the check is expensive and
we want to avoid doing the conversion twice (convert + encode).
These checks could be done in convert only, but then we would need a
struct{m *big.Int, exp int} for decimals as intermediate format.

We would be able to accept other datatypes as well, like
int??, *big.Int, string, ...
but as the user needs to use Decimal anyway (scan), we go with
*big.Rat only for the time being.
*/
func convertDecimal(ft fieldType, v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	if v, ok := v.(*big.Rat); ok {
		return v, nil
	}
	return nil, newConvertError(ft, v, nil)
}

func convertBytes(ft fieldType, v interface{}) (interface{}, error) {
	if v == nil {
		return v, nil
	}

	switch v := v.(type) {

	case string, []byte:
		return v, nil
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {

	case reflect.String:
		return rv.String(), nil

	case reflect.Slice:
		if rv.Type() == bytesReflectType {
			return rv.Bytes(), nil
		}

	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertBytes(ft, rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(bytesReflectType) {
		bv := rv.Convert(bytesReflectType)
		return bv.Interface().([]byte), nil
	}
	return nil, newConvertError(ft, v, nil)
}

// decimals
const (
	// http://en.wikipedia.org/wiki/Decimal128_floating-point_format
	dec128Digits = 34
	// 	dec128Bias   = 6176
	dec128MinExp = -6176
	dec128MaxExp = 6111
)

var (
	natZero = big.NewInt(0)
	natOne  = big.NewInt(1)
	natTen  = big.NewInt(10)
)

const maxNatExp10 = 38 // maximal fixed decimal precision

var natExp10 = make([]*big.Int, maxNatExp10)

func init() {
	natExp10[0], natExp10[1] = natOne, natTen
	for i := 2; i < maxNatExp10; i++ {
		natExp10[i] = new(big.Int).Mul(natExp10[i-1], natTen)
	}
}

// decimal flag
const (
	dfNotExact byte = 1 << iota
	dfOverflow
	dfUnderflow
)

func convertDecimalToRat(m *big.Int, exp int) *big.Rat {
	if m == nil {
		return nil
	}

	v := new(big.Rat).SetInt(m)
	p := v.Num()
	q := v.Denom()

	switch {
	case exp < 0:
		q.Set(exp10(exp * -1))
	case exp == 0:
		q.Set(natOne)
	case exp > 0:
		p.Mul(p, exp10(exp))
		q.Set(natOne)
	}
	return v
}

func convertRatToDecimal(x *big.Rat, m *big.Int, digits, minExp, maxExp int) (int, byte) {
	if x.Num().Cmp(natZero) == 0 { // zero
		m.Set(natZero)
		return 0, 0
	}

	var tmp big.Rat

	c := (&tmp).Set(x) // copy
	a := c.Num()
	b := c.Denom()

	exp, shift := 0, 0

	if c.IsInt() {
		exp = digits10(a) - 1
	} else {
		shift = digits10(a) - digits10(b)
		switch {
		case shift < 0:
			a.Mul(a, exp10(shift*-1))
		case shift > 0:
			b.Mul(b, exp10(shift))
		}
		if a.Cmp(b) == -1 {
			exp = shift - 1
		} else {
			exp = shift
		}
	}

	var df byte

	switch {
	default:
		exp = max(exp-digits+1, minExp)
	case exp < minExp:
		df |= dfUnderflow
		exp = exp - digits + 1
	}

	if exp > maxExp {
		df |= dfOverflow
	}

	shift = exp - shift
	switch {
	case shift < 0:
		a.Mul(a, exp10(shift*-1))
	case exp > 0:
		b.Mul(b, exp10(shift))
	}

	m.QuoRem(a, b, a) // reuse a as rest
	if a.Cmp(natZero) != 0 {
		// round (business >= 0.5 up)
		df |= dfNotExact
		if a.Add(a, a).Cmp(b) >= 0 {
			m.Add(m, natOne)
			if m.Cmp(exp10(digits)) == 0 {
				shift := min(digits, maxExp-exp)
				if shift < 1 { // overflow -> shift one at minimum
					df |= dfOverflow
					shift = 1
				}
				m.Set(exp10(digits - shift))
				exp += shift
			}
		}
	}

	// norm
	for exp < maxExp {
		a.QuoRem(m, natTen, b) // reuse a, b
		if b.Cmp(natZero) != 0 {
			break
		}
		m.Set(a)
		exp++
	}

	return exp, df
}

func convertFixedToRat(m *big.Int, scale int) *big.Rat {
	if m == nil {
		return nil
	}
	if scale < 0 {
		panic(fmt.Sprintf("fixed: invalid scale: %d", scale))
	}
	q := exp10(scale)
	return new(big.Rat).SetFrac(m, q)
}

func convertRatToFixed(r *big.Rat, m *big.Int, prec, scale int) byte {
	if scale < 0 {
		panic(fmt.Sprintf("fixed: invalid scale: %d", scale))
	}

	var df byte

	m.Set(r.Num())
	m.Mul(m, exp10(scale))

	var tmp big.Rat

	c := (&tmp).SetFrac(m, r.Denom()) // norm
	a := c.Num()
	b := c.Denom()

	if b.Cmp(natZero) == 0 { //
		m.Set(a)
		return df
	}

	m.QuoRem(a, b, a) // reuse a as rest
	if a.Cmp(natZero) != 0 {
		// round (business >= 0.5 up)
		df |= dfNotExact
		if a.Add(a, a).Cmp(b) >= 0 {
			m.Add(m, natOne)
		}
	}

	max := exp10(prec)
	min := new(big.Int).Neg(max)

	if m.Cmp(min) <= 0 || m.Cmp(max) >= 0 {
		df |= dfOverflow
	}
	return df
}

// performance: tested with reference work variable
// - but int.Set is expensive, so let's live with big.Int creation for n >= len(nat)
func exp10(n int) *big.Int {
	if n < len(natExp10) {
		return natExp10[n]
	}
	r := big.NewInt(int64(n))
	return r.Exp(natTen, r, nil)
}

func digits10(p *big.Int) int {
	k := p.BitLen() // 2^k <= p < 2^(k+1) - 1
	//i := int(float64(k) / lg10) //minimal digits base 10
	i := k * 100 / 332
	if i < 1 {
		i = 1
	}

	for ; ; i++ {
		if p.Cmp(exp10(i)) < 0 {
			return i
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Longdate
func convertLongdateToTime(longdate int64) time.Time {
	const dayfactor = 10000000 * 24 * 60 * 60
	longdate--
	d := (longdate % dayfactor) * 100
	t := convertDaydateToTime((longdate / dayfactor) + 1)
	return t.Add(time.Duration(d))
}

// nanosecond: HDB - 7 digits precision (not 9 digits)
func convertTimeToLongdate(t time.Time) int64 {
	return (((((((convertTimeToDayDate(t)-1)*24)+int64(t.Hour()))*60)+int64(t.Minute()))*60)+int64(t.Second()))*1e7 + int64(t.Nanosecond()/1e2) + 1
}

// Seconddate
func convertSeconddateToTime(seconddate int64) time.Time {
	const dayfactor = 24 * 60 * 60
	seconddate--
	d := (seconddate % dayfactor) * 1e9
	t := convertDaydateToTime((seconddate / dayfactor) + 1)
	return t.Add(time.Duration(d))
}
func convertTimeToSeconddate(t time.Time) int64 {
	return (((((convertTimeToDayDate(t)-1)*24)+int64(t.Hour()))*60)+int64(t.Minute()))*60 + int64(t.Second()) + 1
}

const julianHdb = 1721423 // 1 January 0001 00:00:00 (1721424) - 1

// Daydate
func convertDaydateToTime(daydate int64) time.Time {
	return julianDayToTime(int(daydate) + julianHdb)
}
func convertTimeToDayDate(t time.Time) int64 {
	return int64(timeToJulianDay(t) - julianHdb)
}

// Secondtime
func convertSecondtimeToTime(secondtime int) time.Time {
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(int64(secondtime-1) * 1e9))
}
func convertTimeToSecondtime(t time.Time) int {
	return (t.Hour()*60+t.Minute())*60 + t.Second() + 1
}
