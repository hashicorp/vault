/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"

	p "github.com/SAP/go-hdb/internal/protocol"
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

// ErrIntegerOutOfRange means that an integer exceeds the size of the hdb integer field.
var ErrIntegerOutOfRange = errors.New("integer out of range error")

// ErrFloatOutOfRange means that a float exceeds the size of the hdb float field.
var ErrFloatOutOfRange = errors.New("float out of range error")

var typeOfTime = reflect.TypeOf((*time.Time)(nil)).Elem()
var typeOfBytes = reflect.TypeOf((*[]byte)(nil)).Elem()

func checkNamedValue(prmFieldSet *p.ParameterFieldSet, nv *driver.NamedValue) error {
	idx := nv.Ordinal - 1

	if idx >= prmFieldSet.NumInputField() {
		return nil
	}

	f := prmFieldSet.Field(idx)
	dt := f.TypeCode().DataType()

	value, err := convertNamedValue(idx, f, dt, nv.Value)

	if err != nil {
		return err
	}

	nv.Value = value
	return nil
}

func convertNamedValue(idx int, f *p.ParameterField, dt p.DataType, v driver.Value) (driver.Value, error) {
	var err error

	// let fields with own Value converter convert themselves first (e.g. NullInt64, ...)
	if _, ok := v.(driver.Valuer); ok {
		if v, err = driver.DefaultParameterConverter.ConvertValue(v); err != nil {
			return nil, err
		}
	}

	switch dt {

	default:
		return nil, fmt.Errorf("convert named value datatype error: %[1]d - %[1]s", dt)

	case p.DtTinyint:
		return convertNvInteger(v, minTinyint, maxTinyint)

	case p.DtSmallint:
		return convertNvInteger(v, minSmallint, maxSmallint)

	case p.DtInteger:
		return convertNvInteger(v, minInteger, maxInteger)

	case p.DtBigint:
		return convertNvInteger(v, minBigint, maxBigint)

	case p.DtReal:
		return convertNvFloat(v, maxReal)

	case p.DtDouble:
		return convertNvFloat(v, maxDouble)

	case p.DtTime:
		return convertNvTime(v)

	case p.DtDecimal:
		return convertNvDecimal(v)

	case p.DtString:
		return convertNvString(v)

	case p.DtBytes:
		return convertNvBytes(v)

	case p.DtLob:
		return convertNvLob(idx, f, v)

	}
}

// integer types
func convertNvInteger(v interface{}, min, max int64) (driver.Value, error) {

	if v == nil {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {

	// bool is represented in HDB as tinyint
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := rv.Int()
		if i64 > max || i64 < min {
			return nil, ErrIntegerOutOfRange
		}
		return i64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := rv.Uint()
		if u64 > uint64(max) {
			return nil, ErrIntegerOutOfRange
		}
		return int64(u64), nil
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertNvInteger(rv.Elem().Interface(), min, max)
	}

	return nil, fmt.Errorf("unsupported integer conversion type error %[1]T %[1]v", v)
}

// float types
func convertNvFloat(v interface{}, max float64) (driver.Value, error) {

	if v == nil {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {

	case reflect.Float32, reflect.Float64:
		f64 := rv.Float()
		if math.Abs(f64) > max {
			return nil, ErrFloatOutOfRange
		}
		return f64, nil
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertNvFloat(rv.Elem().Interface(), max)
	}

	return nil, fmt.Errorf("unsupported float conversion type error %[1]T %[1]v", v)
}

// time
func convertNvTime(v interface{}) (driver.Value, error) {

	if v == nil {
		return nil, nil
	}

	switch v := v.(type) {

	case time.Time:
		return v, nil
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {

	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertNvTime(rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(typeOfTime) {
		tv := rv.Convert(typeOfTime)
		return tv.Interface().(time.Time), nil
	}

	return nil, fmt.Errorf("unsupported time conversion type error %[1]T %[1]v", v)
}

// decimal
func convertNvDecimal(v interface{}) (driver.Value, error) {

	if v == nil {
		return nil, nil
	}

	if v, ok := v.([]byte); ok {
		return v, nil
	}

	return nil, fmt.Errorf("unsupported decimal conversion type error %[1]T %[1]v", v)
}

// string
func convertNvString(v interface{}) (driver.Value, error) {

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
		if rv.Type() == typeOfBytes {
			return rv.Bytes(), nil
		}

	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertNvString(rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(typeOfBytes) {
		bv := rv.Convert(typeOfBytes)
		return bv.Interface().([]byte), nil
	}

	return nil, fmt.Errorf("unsupported character conversion type error %[1]T %[1]v", v)
}

// bytes
func convertNvBytes(v interface{}) (driver.Value, error) {

	if v == nil {
		return v, nil
	}

	if v, ok := v.([]byte); ok {
		return v, nil
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {

	case reflect.Slice:
		if rv.Type() == typeOfBytes {
			return rv.Bytes(), nil
		}

	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertNvBytes(rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(typeOfBytes) {
		bv := rv.Convert(typeOfBytes)
		return bv.Interface().([]byte), nil
	}

	return nil, fmt.Errorf("unsupported bytes conversion type error %[1]T %[1]v", v)
}

// Lob
func convertNvLob(idx int, f *p.ParameterField, v interface{}) (driver.Value, error) {

	if v == nil {
		return v, nil
	}

	switch v := v.(type) {
	case Lob:
		if v.rd == nil {
			return nil, fmt.Errorf("lob error: initial reader %[1]T %[1]v", v)
		}
		f.SetLobReader(v.rd)
		return fmt.Sprintf("<lob %d", idx), nil
	case *Lob:
		if v.rd == nil {
			return nil, fmt.Errorf("lob error: initial reader %[1]T %[1]v", v)
		}
		f.SetLobReader(v.rd)
		return fmt.Sprintf("<lob %d", idx), nil
	case NullLob:
		if !v.Valid {
			return nil, nil
		}
		if v.Lob.rd == nil {
			return nil, fmt.Errorf("lob error: initial reader %[1]T %[1]v", v)
		}
		f.SetLobReader(v.Lob.rd)
		return fmt.Sprintf("<lob %d", idx), nil
	case *NullLob:
		if !v.Valid {
			return nil, nil
		}
		if v.Lob.rd == nil {
			return nil, fmt.Errorf("lob error: initial reader %[1]T %[1]v", v)
		}
		f.SetLobReader(v.Lob.rd)
		return fmt.Sprintf("<lob %d", idx), nil
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {

	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return convertNvLob(idx, f, rv.Elem().Interface())
	}

	return nil, fmt.Errorf("unsupported lob conversion type error %[1]T %[1]v", v)
}
