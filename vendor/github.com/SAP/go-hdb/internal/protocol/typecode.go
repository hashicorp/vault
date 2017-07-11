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

//go:generate stringer -type=typeCode

// null value indicator is high bit
type typeCode byte

const (
	tcNull      typeCode = 0
	tcTinyint   typeCode = 1
	tcSmallint  typeCode = 2
	tcInt       typeCode = 3
	tcBigint    typeCode = 4
	tcDecimal   typeCode = 5
	tcReal      typeCode = 6
	tcDouble    typeCode = 7
	tcChar      typeCode = 8
	tcVarchar   typeCode = 9
	tcNchar     typeCode = 10
	tcNvarchar  typeCode = 11
	tcBinary    typeCode = 12
	tcVarbinary typeCode = 13
	// depricated with 3 (doku) - but table 'date' field uses it
	tcDate typeCode = 14
	// depricated with 3 (doku) - but table 'time' field uses it
	tcTime typeCode = 15
	// depricated with 3 (doku) - but table 'timestamp' field uses it
	tcTimestamp typeCode = 16
	//tcTimetz            typeCode = 17 // reserved: do not use
	//tcTimeltz           typeCode = 18 // reserved: do not use
	//tcTimestamptz       typeCode = 19 // reserved: do not use
	//tcTimestampltz      typeCode = 20 // reserved: do not use
	//tcInvervalym        typeCode = 21 // reserved: do not use
	//tcInvervalds        typeCode = 22 // reserved: do not use
	//tcRowid             typeCode = 23 // reserved: do not use
	//tcUrowid            typeCode = 24 // reserved: do not use
	tcClob     typeCode = 25
	tcNclob    typeCode = 26
	tcBlob     typeCode = 27
	tcBoolean  typeCode = 28
	tcString   typeCode = 29
	tcNstring  typeCode = 30
	tcBlocator typeCode = 31
	tcNlocator typeCode = 32
	tcBstring  typeCode = 33
	//tcDecimaldigitarray typeCode = 34 // reserved: do not use
	tcVarchar2   typeCode = 35
	tcVarchar3   typeCode = 36
	tcNvarchar3  typeCode = 37
	tcVarbinary3 typeCode = 38
	//tcVargroup          typeCode = 39 // reserved: do not use
	//tcTinyintnotnull    typeCode = 40 // reserved: do not use
	//tcSmallintnotnull   typeCode = 41 // reserved: do not use
	//tcIntnotnull        typeCode = 42 // reserved: do not use
	//tcBigintnotnull     typeCode = 43 // reserved: do not use
	//tcArgument          typeCode = 44 // reserved: do not use
	//tcTable             typeCode = 45 // reserved: do not use
	//tcCursor            typeCode = 46 // reserved: do not use
	tcSmalldecimal typeCode = 47
	//tcAbapitab          typeCode = 48 // not supported by GO hdb driver
	//tcAbapstruct        typeCode = 49 // not supported by GO hdb driver
	//tcArray             typeCode = 50 // reserved: do not use
	tcText      typeCode = 51
	tcShorttext typeCode = 52
	tcBintext   typeCode = 53
	//tcFixedpointdecimal typeCode = 54 // reserved: do not use
	tcAlphanum typeCode = 55
	//tcTlocator    typeCode = 56 // reserved: do not use
	tcLongdate   typeCode = 61
	tcSeconddate typeCode = 62
	tcDaydate    typeCode = 63
	tcSecondtime typeCode = 64
	//tcCsdate      typeCode = 65 // reserved: do not use
	//tcCstime      typeCode = 66 // reserved: do not use
	//tcBlobdisk    typeCode = 71 // reserved: do not use
	//tcClobdisk    typeCode = 72 // reserved: do not use
	//tcNclobdisk   typeCode = 73 // reserved: do not use
	tcGeometry typeCode = 74
	tcPoint    typeCode = 75
	//tcFixed16     typeCode = 76 // reserved: do not use
	//tcBlobhybrid  typeCode = 77 // reserved: do not use
	//tcClobhybrid  typeCode = 78 // reserved: do not use
	//tcNclobhybrid typeCode = 79 // reserved: do not use
	tcPointz typeCode = 80
)

func (k typeCode) isLob() bool {
	return k == tcClob || k == tcNclob || k == tcBlob
}

func (k typeCode) isCharBased() bool {
	return k == tcNvarchar || k == tcNstring || k == tcNclob
}

func (k typeCode) dataType() DataType {
	switch k {
	default:
		return DtUnknown
	case tcTinyint:
		return DtTinyint
	case tcSmallint:
		return DtSmallint
	case tcInt:
		return DtInt
	case tcBigint:
		return DtBigint
	case tcReal:
		return DtReal
	case tcDouble:
		return DtDouble
	case tcDate, tcTime, tcTimestamp:
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
