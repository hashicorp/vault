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

// ErrorIntegerOutOfRange means that an integer exceeds the size of the hdb integer field.
var ErrIntegerOutOfRange = errors.New("integer out of range error")

// ErrorIntegerOutOfRange means that a float exceeds the size of the hdb float field.
var ErrFloatOutOfRange = errors.New("float out of range error")

var typeOfTime = reflect.TypeOf((*time.Time)(nil)).Elem()
var typeOfBytes = reflect.TypeOf((*[]byte)(nil)).Elem()

func columnConverter(dt p.DataType) driver.ValueConverter {

	switch dt {

	default:
		return dbUnknownType{}
	case p.DtTinyint:
		return dbTinyint
	case p.DtSmallint:
		return dbSmallint
	case p.DtInt:
		return dbInt
	case p.DtBigint:
		return dbBigint
	case p.DtReal:
		return dbReal
	case p.DtDouble:
		return dbDouble
	case p.DtTime:
		return dbTime
	case p.DtDecimal:
		return dbDecimal
	case p.DtString:
		return dbString
	case p.DtBytes:
		return dbBytes
	case p.DtLob:
		return dbLob
	}
}

// unknown type
type dbUnknownType struct{}

var _ driver.ValueConverter = dbUnknownType{} //check that type implements interface

func (t dbUnknownType) ConvertValue(v interface{}) (driver.Value, error) {
	return nil, fmt.Errorf("column converter for data %v type %T is not implemented", v, v)
}

// int types
var dbTinyint = dbIntType{min: minTinyint, max: maxTinyint}
var dbSmallint = dbIntType{min: minSmallint, max: maxSmallint}
var dbInt = dbIntType{min: minInteger, max: maxInteger}
var dbBigint = dbIntType{min: minBigint, max: maxBigint}

type dbIntType struct {
	min int64
	max int64
}

var _ driver.ValueConverter = dbIntType{} //check that type implements interface

func (i dbIntType) ConvertValue(v interface{}) (driver.Value, error) {

	if v == nil {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := rv.Int()
		if i64 > i.max || i64 < i.min {
			return nil, ErrIntegerOutOfRange
		}
		return i64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := rv.Uint()
		if u64 > uint64(i.max) {
			return nil, ErrIntegerOutOfRange
		}
		return int64(u64), nil
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return i.ConvertValue(rv.Elem().Interface())
	}

	return nil, fmt.Errorf("unsupported integer conversion type error %T %v", v, v)
}

//float types
var dbReal = dbFloatType{max: maxReal}
var dbDouble = dbFloatType{max: maxDouble}

type dbFloatType struct {
	max float64
}

var _ driver.ValueConverter = dbFloatType{} //check that type implements interface

func (f dbFloatType) ConvertValue(v interface{}) (driver.Value, error) {

	if v == nil {
		return v, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {

	case reflect.Float32, reflect.Float64:
		f64 := rv.Float()
		if math.Abs(f64) > f.max {
			return nil, ErrFloatOutOfRange
		}
		return f64, nil
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		}
		return f.ConvertValue(rv.Elem().Interface())
	}

	return nil, fmt.Errorf("unsupported float conversion type error %T %v", v, v)
}

//time
var dbTime = dbTimeType{}

type dbTimeType struct{}

var _ driver.ValueConverter = dbTimeType{} //check that type implements interface

func (t dbTimeType) ConvertValue(v interface{}) (driver.Value, error) {

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
		return t.ConvertValue(rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(typeOfTime) {
		tv := rv.Convert(typeOfTime)
		return tv.Interface().(time.Time), nil
	}

	return nil, fmt.Errorf("unsupported time conversion type error %T %v", v, v)
}

//decimal
var dbDecimal = dbDecimalType{}

type dbDecimalType struct{}

var _ driver.ValueConverter = dbDecimalType{} //check that type implements interface

func (d dbDecimalType) ConvertValue(v interface{}) (driver.Value, error) {

	if v == nil {
		return nil, nil
	}

	if v, ok := v.([]byte); ok {
		return v, nil
	}

	return nil, fmt.Errorf("unsupported decimal conversion type error %T %v", v, v)
}

//string
var dbString = dbStringType{}

type dbStringType struct{}

var _ driver.ValueConverter = dbStringType{} //check that type implements interface

func (d dbStringType) ConvertValue(v interface{}) (driver.Value, error) {

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
		return d.ConvertValue(rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(typeOfBytes) {
		bv := rv.Convert(typeOfBytes)
		return bv.Interface().([]byte), nil
	}

	return nil, fmt.Errorf("unsupported character conversion type error %T %v", v, v)
}

//bytes
var dbBytes = dbBytesType{}

type dbBytesType struct{}

var _ driver.ValueConverter = dbBytesType{} //check that type implements interface

func (d dbBytesType) ConvertValue(v interface{}) (driver.Value, error) {

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
		return d.ConvertValue(rv.Elem().Interface())
	}

	if rv.Type().ConvertibleTo(typeOfBytes) {
		bv := rv.Convert(typeOfBytes)
		return bv.Interface().([]byte), nil
	}

	return nil, fmt.Errorf("unsupported bytes conversion type error %T %v", v, v)
}

//lob
var dbLob = dbLobType{}

type dbLobType struct{}

var _ driver.ValueConverter = dbLobType{} //check that type implements interface

func (d dbLobType) ConvertValue(v interface{}) (driver.Value, error) {

	if v, ok := v.(int64); ok {
		return v, nil
	}

	return nil, fmt.Errorf("unsupported lob conversion type error %T %v", v, v)
}
