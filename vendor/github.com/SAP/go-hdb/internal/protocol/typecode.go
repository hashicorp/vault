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

package protocol

import (
	"strings"
)

//go:generate stringer -type=TypeCode

// TypeCode identify the type of a field transferred to or from the database.
type TypeCode byte

// null value indicator is high bit

const (
	tcNull      TypeCode = 0
	tcTinyint   TypeCode = 1
	tcSmallint  TypeCode = 2
	tcInteger   TypeCode = 3
	tcBigint    TypeCode = 4
	tcDecimal   TypeCode = 5
	tcReal      TypeCode = 6
	tcDouble    TypeCode = 7
	tcChar      TypeCode = 8
	tcVarchar   TypeCode = 9
	tcNchar     TypeCode = 10
	tcNvarchar  TypeCode = 11
	tcBinary    TypeCode = 12
	tcVarbinary TypeCode = 13
	// deprecated with 3 (doku) - but table 'date' field uses it
	tcDate TypeCode = 14
	// deprecated with 3 (doku) - but table 'time' field uses it
	tcTime TypeCode = 15
	// deprecated with 3 (doku) - but table 'timestamp' field uses it
	tcTimestamp TypeCode = 16
	//tcTimetz            TypeCode = 17 // reserved: do not use
	//tcTimeltz           TypeCode = 18 // reserved: do not use
	//tcTimestamptz       TypeCode = 19 // reserved: do not use
	//tcTimestampltz      TypeCode = 20 // reserved: do not use
	//tcInvervalym        TypeCode = 21 // reserved: do not use
	//tcInvervalds        TypeCode = 22 // reserved: do not use
	//tcRowid             TypeCode = 23 // reserved: do not use
	//tcUrowid            TypeCode = 24 // reserved: do not use
	tcClob     TypeCode = 25
	tcNclob    TypeCode = 26
	tcBlob     TypeCode = 27
	tcBoolean  TypeCode = 28
	tcString   TypeCode = 29
	tcNstring  TypeCode = 30
	tcBlocator TypeCode = 31
	tcNlocator TypeCode = 32
	tcBstring  TypeCode = 33
	//tcDecimaldigitarray TypeCode = 34 // reserved: do not use
	tcVarchar2   TypeCode = 35
	tcVarchar3   TypeCode = 36
	tcNvarchar3  TypeCode = 37
	tcVarbinary3 TypeCode = 38
	//tcVargroup          TypeCode = 39 // reserved: do not use
	//tcTinyintnotnull    TypeCode = 40 // reserved: do not use
	//tcSmallintnotnull   TypeCode = 41 // reserved: do not use
	//tcIntnotnull        TypeCode = 42 // reserved: do not use
	//tcBigintnotnull     TypeCode = 43 // reserved: do not use
	//tcArgument          TypeCode = 44 // reserved: do not use
	//tcTable             TypeCode = 45 // reserved: do not use
	//tcCursor            TypeCode = 46 // reserved: do not use
	tcSmalldecimal TypeCode = 47
	//tcAbapitab          TypeCode = 48 // not supported by GO hdb driver
	//tcAbapstruct        TypeCode = 49 // not supported by GO hdb driver
	tcArray     TypeCode = 50
	tcText      TypeCode = 51
	tcShorttext TypeCode = 52
	//tcFixedString       TypeCode = 53 // reserved: do not use
	//tcFixedpointdecimal TypeCode = 54 // reserved: do not use
	tcAlphanum TypeCode = 55
	//tcTlocator    TypeCode = 56 // reserved: do not use
	tcLongdate   TypeCode = 61
	tcSeconddate TypeCode = 62
	tcDaydate    TypeCode = 63
	tcSecondtime TypeCode = 64
	//tcCte      TypeCode = 65 // reserved: do not use
	//tcCstimesda      TypeCode = 66 // reserved: do not use
	//tcBlobdisk    TypeCode = 71 // reserved: do not use
	//tcClobdisk    TypeCode = 72 // reserved: do not use
	//tcNclobdisk   TypeCode = 73 // reserved: do not use
	//tcGeometry    TypeCode = 74 // reserved: do not use
	//tcPoint       TypeCode = 75 // reserved: do not use
	//tcFixed16     TypeCode = 76 // reserved: do not use
	//tcBlobhybrid  TypeCode = 77 // reserved: do not use
	//tcClobhybrid  TypeCode = 78 // reserved: do not use
	//tcNclobhybrid TypeCode = 79 // reserved: do not use
	//tcPointz      TypeCode = 80 // reserved: do not use
)

func (k TypeCode) isLob() bool {
	return k == tcClob || k == tcNclob || k == tcBlob
}

func (k TypeCode) isCharBased() bool {
	return k == tcNvarchar || k == tcNstring || k == tcNclob
}

func (k TypeCode) isVariableLength() bool {
	return k == tcChar || k == tcNchar || k == tcVarchar || k == tcNvarchar || k == tcBinary || k == tcVarbinary || k == tcShorttext || k == tcAlphanum
}

func (k TypeCode) isDecimalType() bool {
	return k == tcSmalldecimal || k == tcDecimal
}

// DataType converts a type code into one of the supported data types by the driver.
func (k TypeCode) DataType() DataType {
	switch k {
	default:
		return DtUnknown
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
	case tcDate, tcTime, tcTimestamp, tcLongdate, tcSeconddate, tcDaydate, tcSecondtime:
		return DtTime
	case tcDecimal:
		return DtDecimal
	case tcChar, tcVarchar, tcString, tcNchar, tcNvarchar, tcNstring:
		return DtString
	case tcBinary, tcVarbinary:
		return DtBytes
	case tcBlob, tcClob, tcNclob:
		return DtLob
	}
}

// TypeName returns the database type name.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypeDatabaseTypeName
func (k TypeCode) TypeName() string {
	return strings.ToUpper(k.String()[2:])
}
