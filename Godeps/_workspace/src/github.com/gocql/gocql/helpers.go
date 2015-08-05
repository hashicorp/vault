// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocql

import (
	"math/big"
	"reflect"
	"strings"
	"time"

	"gopkg.in/inf.v0"
)

type RowData struct {
	Columns []string
	Values  []interface{}
}

func goType(t TypeInfo) reflect.Type {
	switch t.Type() {
	case TypeVarchar, TypeAscii, TypeInet:
		return reflect.TypeOf(*new(string))
	case TypeBigInt, TypeCounter:
		return reflect.TypeOf(*new(int64))
	case TypeTimestamp:
		return reflect.TypeOf(*new(time.Time))
	case TypeBlob:
		return reflect.TypeOf(*new([]byte))
	case TypeBoolean:
		return reflect.TypeOf(*new(bool))
	case TypeFloat:
		return reflect.TypeOf(*new(float32))
	case TypeDouble:
		return reflect.TypeOf(*new(float64))
	case TypeInt:
		return reflect.TypeOf(*new(int))
	case TypeDecimal:
		return reflect.TypeOf(*new(*inf.Dec))
	case TypeUUID, TypeTimeUUID:
		return reflect.TypeOf(*new(UUID))
	case TypeList, TypeSet:
		return reflect.SliceOf(goType(t.(CollectionType).Elem))
	case TypeMap:
		return reflect.MapOf(goType(t.(CollectionType).Key), goType(t.(CollectionType).Elem))
	case TypeVarint:
		return reflect.TypeOf(*new(*big.Int))
	case TypeTuple:
		// what can we do here? all there is to do is to make a list of interface{}
		tuple := t.(TupleTypeInfo)
		return reflect.TypeOf(make([]interface{}, len(tuple.Elems)))
	default:
		return nil
	}
}

func dereference(i interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(i)).Interface()
}

func getApacheCassandraType(class string) Type {
	switch strings.TrimPrefix(class, apacheCassandraTypePrefix) {
	case "AsciiType":
		return TypeAscii
	case "LongType":
		return TypeBigInt
	case "BytesType":
		return TypeBlob
	case "BooleanType":
		return TypeBoolean
	case "CounterColumnType":
		return TypeCounter
	case "DecimalType":
		return TypeDecimal
	case "DoubleType":
		return TypeDouble
	case "FloatType":
		return TypeFloat
	case "Int32Type":
		return TypeInt
	case "DateType", "TimestampType":
		return TypeTimestamp
	case "UUIDType":
		return TypeUUID
	case "UTF8Type":
		return TypeVarchar
	case "IntegerType":
		return TypeVarint
	case "TimeUUIDType":
		return TypeTimeUUID
	case "InetAddressType":
		return TypeInet
	case "MapType":
		return TypeMap
	case "ListType":
		return TypeList
	case "SetType":
		return TypeSet
	case "TupleType":
		return TypeTuple
	default:
		return TypeCustom
	}
}

func (r *RowData) rowMap(m map[string]interface{}) {
	for i, column := range r.Columns {
		val := dereference(r.Values[i])
		if valVal := reflect.ValueOf(val); valVal.Kind() == reflect.Slice {
			valCopy := reflect.MakeSlice(valVal.Type(), valVal.Len(), valVal.Cap())
			reflect.Copy(valCopy, valVal)
			m[column] = valCopy.Interface()
		} else {
			m[column] = val
		}
	}
}

func (iter *Iter) RowData() (RowData, error) {
	if iter.err != nil {
		return RowData{}, iter.err
	}
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	for _, column := range iter.Columns() {
		val := column.TypeInfo.New()
		columns = append(columns, column.Name)
		values = append(values, val)
	}
	rowData := RowData{
		Columns: columns,
		Values:  values,
	}
	return rowData, nil
}

// SliceMap is a helper function to make the API easier to use
// returns the data from the query in the form of []map[string]interface{}
func (iter *Iter) SliceMap() ([]map[string]interface{}, error) {
	if iter.err != nil {
		return nil, iter.err
	}

	// Not checking for the error because we just did
	rowData, _ := iter.RowData()
	dataToReturn := make([]map[string]interface{}, 0)
	for iter.Scan(rowData.Values...) {
		m := make(map[string]interface{})
		rowData.rowMap(m)
		dataToReturn = append(dataToReturn, m)
	}
	if iter.err != nil {
		return nil, iter.err
	}
	return dataToReturn, nil
}

// MapScan takes a map[string]interface{} and populates it with a row
// That is returned from cassandra.
func (iter *Iter) MapScan(m map[string]interface{}) bool {
	if iter.err != nil {
		return false
	}

	// Not checking for the error because we just did
	rowData, _ := iter.RowData()

	for i, col := range rowData.Columns {
		if dest, ok := m[col]; ok {
			rowData.Values[i] = dest
		}
	}

	if iter.Scan(rowData.Values...) {
		rowData.rowMap(m)
		return true
	}
	return false
}
