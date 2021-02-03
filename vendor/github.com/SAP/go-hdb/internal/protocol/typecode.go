// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=typeCode

// typeCode identify the type of a field transferred to or from the database.
type typeCode byte

// null value indicator is high bit

const (
	tcNullL             typeCode = 0x00
	tcTinyint           typeCode = 0x01
	tcSmallint          typeCode = 0x02
	tcInteger           typeCode = 0x03
	tcBigint            typeCode = 0x04
	tcDecimal           typeCode = 0x05
	tcReal              typeCode = 0x06
	tcDouble            typeCode = 0x07
	tcChar              typeCode = 0x08
	tcVarchar           typeCode = 0x09 // changed from tcVarchar1 to tcVarchar (ref hdbclient)
	tcNchar             typeCode = 0x0A
	tcNvarchar          typeCode = 0x0B
	tcBinary            typeCode = 0x0C
	tcVarbinary         typeCode = 0x0D
	tcDate              typeCode = 0x0E
	tcTime              typeCode = 0x0F
	tcTimestamp         typeCode = 0x10
	tcTimetz            typeCode = 0x11
	tcTimeltz           typeCode = 0x12
	tcTimestampTz       typeCode = 0x13
	tcTimestampLtz      typeCode = 0x14
	tcIntervalYm        typeCode = 0x15
	tcIntervalDs        typeCode = 0x16
	tcRowid             typeCode = 0x17
	tcUrowid            typeCode = 0x18
	tcClob              typeCode = 0x19
	tcNclob             typeCode = 0x1A
	tcBlob              typeCode = 0x1B
	tcBoolean           typeCode = 0x1C
	tcString            typeCode = 0x1D
	tcNstring           typeCode = 0x1E
	tcLocator           typeCode = 0x1F
	tcNlocator          typeCode = 0x20
	tcBstring           typeCode = 0x21
	tcDecimalDigitArray typeCode = 0x22
	tcVarchar2          typeCode = 0x23
	tcTable             typeCode = 0x2D
	tcSmalldecimal      typeCode = 0x2f // inserted (not existent in hdbclient)
	tcAbapstream        typeCode = 0x30
	tcAbapstruct        typeCode = 0x31
	tcAarray            typeCode = 0x32
	tcText              typeCode = 0x33
	tcShorttext         typeCode = 0x34
	tcBintext           typeCode = 0x35
	tcAlphanum          typeCode = 0x37
	tcLongdate          typeCode = 0x3D
	tcSeconddate        typeCode = 0x3E
	tcDaydate           typeCode = 0x3F
	tcSecondtime        typeCode = 0x40
	tcClocator          typeCode = 0x46
	tcBlobDiskReserved  typeCode = 0x47
	tcClobDiskReserved  typeCode = 0x48
	tcNclobDiskReserved typeCode = 0x49
	tcStGeometry        typeCode = 0x4A
	tcStPoint           typeCode = 0x4B
	tcFixed16           typeCode = 0x4C
	tcAbapItab          typeCode = 0x4D
	tcRecordRowStore    typeCode = 0x4E
	tcRecordColumnStore typeCode = 0x4F
	tcFixed8            typeCode = 0x51
	tcFixed12           typeCode = 0x52
	tcCiphertext        typeCode = 0x5A

	// additional internal typecodes
	tcTableRef  typeCode = 0x7e // 126
	tcTableRows typeCode = 0x7f // 127
)

func (tc typeCode) isLob() bool {
	return tc == tcClob || tc == tcNclob || tc == tcBlob || tc == tcText || tc == tcBintext || tc == tcLocator || tc == tcNlocator
}

func (tc typeCode) isCharBased() bool {
	return tc == tcNvarchar || tc == tcNstring || tc == tcNclob || tc == tcText || tc == tcBintext
}

func (tc typeCode) isVariableLength() bool {
	return tc == tcChar || tc == tcNchar || tc == tcVarchar || tc == tcNvarchar || tc == tcBinary || tc == tcVarbinary || tc == tcShorttext || tc == tcAlphanum
}

func (tc typeCode) isIntegerType() bool {
	return tc == tcTinyint || tc == tcSmallint || tc == tcInteger || tc == tcBigint
}

func (tc typeCode) isDecimalType() bool {
	return tc == tcSmalldecimal || tc == tcDecimal || tc == tcFixed8 || tc == tcFixed12 || tc == tcFixed16
}

//
func (tc typeCode) supportNullValue() bool {
	/*
		(*1)
		HDB bug: secondtime null value cannot be set by setting high bit
		- trying so, gives:
		  SQL HdbError 1033 - error while parsing protocol: no such data type: type_code=192, index=2

		Traffic analysis of python client (https://pypi.org/project/hdbcli) resulted in:
		- set null value constant directly instead of using high bit

		(*2)
		boolean: use false =:= 0; null =:= 1; true =:= 2

	*/
	return !(tc == tcBoolean || tc == tcSecondtime)
}

// see hdbclient
func (tc typeCode) encTc() typeCode {
	switch tc {
	default:
		return tc
	case tcText, tcBintext, tcLocator:
		return tcNclob
	}
}

/*
tcBintext:
- protocol returns tcLocator for tcBintext
- see dataTypeMap and encTc
*/

func (tc typeCode) dataType() DataType {
	// performance: use switch instead of map
	switch tc {
	case tcBoolean:
		return DtBoolean
	case tcTinyint:
		return DtTinyint
	case tcSmallint:
		return DtSmallint
	case tcInteger:
		return DtInteger
	case tcBigint:
		return DtBigint
	case tcReal:
		return DtReal
	case tcDouble:
		return DtDouble
	case tcDate:
		return DtTime
	case tcTime, tcTimestamp, tcLongdate, tcSeconddate, tcDaydate, tcSecondtime:
		return DtTime
	case tcDecimal, tcFixed8, tcFixed12, tcFixed16:
		return DtDecimal
	case tcChar, tcVarchar, tcString, tcAlphanum, tcNchar, tcNvarchar, tcNstring, tcShorttext, tcTableRef:
		return DtString
	case tcBinary, tcVarbinary:
		return DtBytes
	case tcBlob, tcClob, tcNclob, tcText, tcBintext:
		return DtLob
	case tcTableRows:
		return DtRows
	default:
		panic(fmt.Sprintf("Missing DataType for typeCode %s", tc))
	}
}

// typeName returns the database type name.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypeDatabaseTypeName
func (tc typeCode) typeName() string {
	return strings.ToUpper(tc.String()[2:])
}

func (tc typeCode) fieldType(length, fraction int) fieldType {
	// performance: use switch instead of map
	switch tc {
	case tcBoolean:
		return booleanType
	case tcTinyint:
		return tinyintType
	case tcSmallint:
		return smallintType
	case tcInteger:
		return integerType
	case tcBigint:
		return bigintType
	case tcReal:
		return realType
	case tcDouble:
		return doubleType
	case tcDate:
		return dateType
	case tcTime:
		return timeType
	case tcTimestamp:
		return timestampType
	case tcLongdate:
		return longdateType
	case tcSeconddate:
		return seconddateType
	case tcDaydate:
		return daydateType
	case tcSecondtime:
		return secondtimeType
	case tcDecimal:
		return decimalType
	case tcChar, tcVarchar, tcString:
		return varType
	case tcAlphanum:
		return alphaType
	case tcNchar, tcNvarchar, tcNstring, tcShorttext:
		return cesu8Type
	case tcBinary, tcVarbinary:
		return varType
	case tcBlob, tcClob, tcLocator:
		return lobVarType
	case tcNclob, tcText, tcNlocator:
		return lobCESU8Type
	case tcBintext: // ?? lobCESU8Type
		return lobVarType
	case tcFixed8:
		return _fixed8Type{prec: length, scale: fraction} // used for decimals(x,y) 2^63 - 1 (int64)
	case tcFixed12:
		return _fixed12Type{prec: length, scale: fraction} // used for decimals(x,y) 2^96 - 1 (int96)
	case tcFixed16:
		return _fixed16Type{prec: length, scale: fraction} // used for decimals(x,y) 2^63 - 1 (int128)
	default:
		panic(fmt.Sprintf("Missing FieldType for typeCode %s", tc))
	}
}
