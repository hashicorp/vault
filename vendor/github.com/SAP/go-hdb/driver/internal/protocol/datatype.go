package protocol

import (
	"database/sql"
	"reflect"
	"time"
)

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
func RegisterScanType(dt DataType, scanType, scanNullType reflect.Type) bool {
	scanTypes[dt].scanType = scanType
	scanTypes[dt].scanNullType = scanNullType
	return true
}

var scanTypes = []struct {
	scanType     reflect.Type
	scanNullType reflect.Type
}{
	DtUnknown:  {reflect.TypeFor[any](), reflect.TypeFor[any]()},
	DtBoolean:  {reflect.TypeFor[bool](), reflect.TypeFor[sql.NullBool]()},
	DtTinyint:  {reflect.TypeFor[uint8](), reflect.TypeFor[sql.NullByte]()},
	DtSmallint: {reflect.TypeFor[int16](), reflect.TypeFor[sql.NullInt16]()},
	DtInteger:  {reflect.TypeFor[int32](), reflect.TypeFor[sql.NullInt32]()},
	DtBigint:   {reflect.TypeFor[int64](), reflect.TypeFor[sql.NullInt64]()},
	DtReal:     {reflect.TypeFor[float32](), reflect.TypeFor[sql.NullFloat64]()},
	DtDouble:   {reflect.TypeFor[float64](), reflect.TypeFor[sql.NullFloat64]()},
	DtTime:     {reflect.TypeFor[time.Time](), reflect.TypeFor[sql.NullTime]()},
	DtString:   {reflect.TypeFor[string](), reflect.TypeFor[sql.NullString]()},
	DtBytes:    {nil, nil}, // to be registered by driver
	DtDecimal:  {nil, nil}, // to be registered by driver
	DtLob:      {nil, nil}, // to be registered by driver
	DtRows:     {reflect.TypeFor[sql.Rows](), reflect.TypeFor[sql.Rows]()},
}

// ScanType return the scan type (reflect.Type) of the corresponding data type.
func (dt DataType) ScanType(nullable bool) reflect.Type {
	if nullable {
		return scanTypes[dt].scanNullType
	}
	return scanTypes[dt].scanType
}
