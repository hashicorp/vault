// +build go1.8

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
	"io"
	"reflect"
	"time"

	p "github.com/SAP/go-hdb/internal/protocol"
)

/*
no result type

  the following golang 1.8 interfaces are not implemented, because no result values are provided anyway:
  - RowsColumnTypeDatabaseTypeName
  - RowsColumnTypeLength
  - RowsColumnTypeNullable
  - RowsColumnTypePrecisionScale
  - RowsColumnTypeScanType
*/

var _ driver.RowsNextResultSet = (*noResultType)(nil) // golang 1.8

func (r *noResultType) HasNextResultSet() bool { return false }
func (r *noResultType) NextResultSet() error   { return io.EOF }

/*
query result
*/

var _ driver.RowsColumnTypeDatabaseTypeName = (*queryResult)(nil)
var _ driver.RowsColumnTypeLength = (*queryResult)(nil)
var _ driver.RowsColumnTypeNullable = (*queryResult)(nil)
var _ driver.RowsColumnTypePrecisionScale = (*queryResult)(nil)
var _ driver.RowsColumnTypeScanType = (*queryResult)(nil)
var _ driver.RowsNextResultSet = (*queryResult)(nil)

func (r *queryResult) ColumnTypeDatabaseTypeName(idx int) string {
	return r.fieldSet.DatabaseTypeName(idx)
}

func (r *queryResult) ColumnTypeLength(idx int) (int64, bool) {
	return r.fieldSet.TypeLength(idx)
}

func (r *queryResult) ColumnTypePrecisionScale(idx int) (int64, int64, bool) {
	return r.fieldSet.TypePrecisionScale(idx)
}

func (r *queryResult) ColumnTypeNullable(idx int) (bool, bool) {
	return r.fieldSet.TypeNullable(idx), true
}

var (
	scanTypeUnknown  = reflect.TypeOf(new(interface{})).Elem()
	scanTypeTinyint  = reflect.TypeOf(uint8(0))
	scanTypeSmallint = reflect.TypeOf(int16(0))
	scanTypeInteger  = reflect.TypeOf(int32(0))
	scanTypeBigint   = reflect.TypeOf(int64(0))
	scanTypeReal     = reflect.TypeOf(float32(0.0))
	scanTypeDouble   = reflect.TypeOf(float64(0.0))
	scanTypeTime     = reflect.TypeOf(time.Time{})
	scanTypeString   = reflect.TypeOf(string(""))
	scanTypeBytes    = reflect.TypeOf([]byte{})
	scanTypeDecimal  = reflect.TypeOf(Decimal{})
	scanTypeLob      = reflect.TypeOf(Lob{})
)

func (r *queryResult) ColumnTypeScanType(idx int) reflect.Type {
	switch r.fieldSet.DataType(idx) {
	default:
		return scanTypeUnknown
	case p.DtTinyint:
		return scanTypeTinyint
	case p.DtSmallint:
		return scanTypeSmallint
	case p.DtInteger:
		return scanTypeInteger
	case p.DtBigint:
		return scanTypeBigint
	case p.DtReal:
		return scanTypeReal
	case p.DtDouble:
		return scanTypeDouble
	case p.DtTime:
		return scanTypeTime
	case p.DtDecimal:
		return scanTypeDecimal
	case p.DtString:
		return scanTypeString
	case p.DtBytes:
		return scanTypeBytes
	case p.DtLob:
		return scanTypeLob
	}
}

func (r *queryResult) HasNextResultSet() bool { return false }
func (r *queryResult) NextResultSet() error   { return io.EOF }
