// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2017 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"database/sql"
	"reflect"
)

func (mf *mysqlField) typeDatabaseName() string {
	switch mf.fieldType {
	case fieldTypeBit:
		return "BIT"
	case fieldTypeBLOB:
		if mf.charSet != binaryCollationID {
			return "TEXT"
		}
		return "BLOB"
	case fieldTypeDate:
		return "DATE"
	case fieldTypeDateTime:
		return "DATETIME"
	case fieldTypeDecimal:
		return "DECIMAL"
	case fieldTypeDouble:
		return "DOUBLE"
	case fieldTypeEnum:
		return "ENUM"
	case fieldTypeFloat:
		return "FLOAT"
	case fieldTypeGeometry:
		return "GEOMETRY"
	case fieldTypeInt24:
		if mf.flags&flagUnsigned != 0 {
			return "UNSIGNED MEDIUMINT"
		}
		return "MEDIUMINT"
	case fieldTypeJSON:
		return "JSON"
	case fieldTypeLong:
		if mf.flags&flagUnsigned != 0 {
			return "UNSIGNED INT"
		}
		return "INT"
	case fieldTypeLongBLOB:
		if mf.charSet != binaryCollationID {
			return "LONGTEXT"
		}
		return "LONGBLOB"
	case fieldTypeLongLong:
		if mf.flags&flagUnsigned != 0 {
			return "UNSIGNED BIGINT"
		}
		return "BIGINT"
	case fieldTypeMediumBLOB:
		if mf.charSet != binaryCollationID {
			return "MEDIUMTEXT"
		}
		return "MEDIUMBLOB"
	case fieldTypeNewDate:
		return "DATE"
	case fieldTypeNewDecimal:
		return "DECIMAL"
	case fieldTypeNULL:
		return "NULL"
	case fieldTypeSet:
		return "SET"
	case fieldTypeShort:
		if mf.flags&flagUnsigned != 0 {
			return "UNSIGNED SMALLINT"
		}
		return "SMALLINT"
	case fieldTypeString:
		if mf.flags&flagEnum != 0 {
			return "ENUM"
		} else if mf.flags&flagSet != 0 {
			return "SET"
		}
		if mf.charSet == binaryCollationID {
			return "BINARY"
		}
		return "CHAR"
	case fieldTypeTime:
		return "TIME"
	case fieldTypeTimestamp:
		return "TIMESTAMP"
	case fieldTypeTiny:
		if mf.flags&flagUnsigned != 0 {
			return "UNSIGNED TINYINT"
		}
		return "TINYINT"
	case fieldTypeTinyBLOB:
		if mf.charSet != binaryCollationID {
			return "TINYTEXT"
		}
		return "TINYBLOB"
	case fieldTypeVarChar:
		if mf.charSet == binaryCollationID {
			return "VARBINARY"
		}
		return "VARCHAR"
	case fieldTypeVarString:
		if mf.charSet == binaryCollationID {
			return "VARBINARY"
		}
		return "VARCHAR"
	case fieldTypeYear:
		return "YEAR"
	default:
		return ""
	}
}

var (
	scanTypeFloat32    = reflect.TypeOf(float32(0))
	scanTypeFloat64    = reflect.TypeOf(float64(0))
	scanTypeInt8       = reflect.TypeOf(int8(0))
	scanTypeInt16      = reflect.TypeOf(int16(0))
	scanTypeInt32      = reflect.TypeOf(int32(0))
	scanTypeInt64      = reflect.TypeOf(int64(0))
	scanTypeNullFloat  = reflect.TypeOf(sql.NullFloat64{})
	scanTypeNullInt    = reflect.TypeOf(sql.NullInt64{})
	scanTypeNullTime   = reflect.TypeOf(sql.NullTime{})
	scanTypeUint8      = reflect.TypeOf(uint8(0))
	scanTypeUint16     = reflect.TypeOf(uint16(0))
	scanTypeUint32     = reflect.TypeOf(uint32(0))
	scanTypeUint64     = reflect.TypeOf(uint64(0))
	scanTypeString     = reflect.TypeOf("")
	scanTypeNullString = reflect.TypeOf(sql.NullString{})
	scanTypeBytes      = reflect.TypeOf([]byte{})
	scanTypeUnknown    = reflect.TypeOf(new(any))
)

type mysqlField struct {
	tableName string
	name      string
	length    uint32
	flags     fieldFlag
	fieldType fieldType
	decimals  byte
	charSet   uint8
}

func (mf *mysqlField) scanType() reflect.Type {
	switch mf.fieldType {
	case fieldTypeTiny:
		if mf.flags&flagNotNULL != 0 {
			if mf.flags&flagUnsigned != 0 {
				return scanTypeUint8
			}
			return scanTypeInt8
		}
		return scanTypeNullInt

	case fieldTypeShort, fieldTypeYear:
		if mf.flags&flagNotNULL != 0 {
			if mf.flags&flagUnsigned != 0 {
				return scanTypeUint16
			}
			return scanTypeInt16
		}
		return scanTypeNullInt

	case fieldTypeInt24, fieldTypeLong:
		if mf.flags&flagNotNULL != 0 {
			if mf.flags&flagUnsigned != 0 {
				return scanTypeUint32
			}
			return scanTypeInt32
		}
		return scanTypeNullInt

	case fieldTypeLongLong:
		if mf.flags&flagNotNULL != 0 {
			if mf.flags&flagUnsigned != 0 {
				return scanTypeUint64
			}
			return scanTypeInt64
		}
		return scanTypeNullInt

	case fieldTypeFloat:
		if mf.flags&flagNotNULL != 0 {
			return scanTypeFloat32
		}
		return scanTypeNullFloat

	case fieldTypeDouble:
		if mf.flags&flagNotNULL != 0 {
			return scanTypeFloat64
		}
		return scanTypeNullFloat

	case fieldTypeBit, fieldTypeTinyBLOB, fieldTypeMediumBLOB, fieldTypeLongBLOB,
		fieldTypeBLOB, fieldTypeVarString, fieldTypeString, fieldTypeGeometry:
		if mf.charSet == binaryCollationID {
			return scanTypeBytes
		}
		fallthrough
	case fieldTypeDecimal, fieldTypeNewDecimal, fieldTypeVarChar,
		fieldTypeEnum, fieldTypeSet, fieldTypeJSON, fieldTypeTime:
		if mf.flags&flagNotNULL != 0 {
			return scanTypeString
		}
		return scanTypeNullString

	case fieldTypeDate, fieldTypeNewDate,
		fieldTypeTimestamp, fieldTypeDateTime:
		// NullTime is always returned for more consistent behavior as it can
		// handle both cases of parseTime regardless if the field is nullable.
		return scanTypeNullTime

	default:
		return scanTypeUnknown
	}
}
