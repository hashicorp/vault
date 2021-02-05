// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

//go:generate stringer -type=DataType

// DataType is the type definition for data types supported by this package.
type DataType byte

// Data type constants.
const (
	DtUnknown DataType = iota // unknown data type
	DtBoolean
	DtTinyint
	DtSmallint
	DtInteger
	DtBigint
	DtReal
	DtDouble
	DtDecimal
	DtTime
	DtString
	DtBytes
	DtLob
	DtRows
)

// RegisterScanType registers driver owned datatype scantypes (e.g. Decimal, Lob).
func RegisterScanType(dt DataType, scanType reflect.Type) {
	scanTypeMap[dt] = scanType
}

var scanTypeMap = map[DataType]reflect.Type{
	DtUnknown:  reflect.TypeOf((*interface{})(nil)).Elem(),
	DtBoolean:  reflect.TypeOf((*bool)(nil)).Elem(),
	DtTinyint:  reflect.TypeOf((*uint8)(nil)).Elem(),
	DtSmallint: reflect.TypeOf((*int16)(nil)).Elem(),
	DtInteger:  reflect.TypeOf((*int32)(nil)).Elem(),
	DtBigint:   reflect.TypeOf((*int64)(nil)).Elem(),
	DtReal:     reflect.TypeOf((*float32)(nil)).Elem(),
	DtDouble:   reflect.TypeOf((*float64)(nil)).Elem(),
	DtTime:     reflect.TypeOf((*time.Time)(nil)).Elem(),
	DtString:   reflect.TypeOf((*string)(nil)).Elem(),
	DtBytes:    reflect.TypeOf((*[]byte)(nil)).Elem(),
	DtDecimal:  nil, // to be registered by driver
	DtLob:      nil, // to be registered by driver
	DtRows:     reflect.TypeOf((*sql.Rows)(nil)).Elem(),
}

// ScanType return the scan type (reflect.Type) of the corresponding data type.
func (dt DataType) ScanType() reflect.Type {
	st, ok := scanTypeMap[dt]
	if !ok {
		panic(fmt.Sprintf("Missing ScanType for DataType %s", dt))
	}
	return st
}
