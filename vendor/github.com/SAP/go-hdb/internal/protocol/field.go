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
	"database/sql/driver"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/SAP/go-hdb/internal/bufio"
	"github.com/SAP/go-hdb/internal/unicode/cesu8"
)

var test uint32

const (
	realNullValue   uint32 = ^uint32(0)
	doubleNullValue uint64 = ^uint64(0)
)

const noFieldName uint32 = 0xFFFFFFFF

type uint32Slice []uint32

func (p uint32Slice) Len() int           { return len(p) }
func (p uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p uint32Slice) sort()              { sort.Sort(p) }

type fieldNames map[uint32]string

func newFieldNames() fieldNames {
	return make(map[uint32]string)
}

func (f fieldNames) addOffset(offset uint32) {
	if offset != noFieldName {
		f[offset] = ""
	}
}

func (f fieldNames) name(offset uint32) string {
	if name, ok := f[offset]; ok {
		return name
	}
	return ""
}

func (f fieldNames) setName(offset uint32, name string) {
	f[offset] = name
}

func (f fieldNames) sortOffsets() []uint32 {
	offsets := make([]uint32, 0, len(f))
	for k := range f {
		offsets = append(offsets, k)
	}
	uint32Slice(offsets).sort()
	return offsets
}

// FieldValues contains rows read from database.
type FieldValues struct {
	rows   int
	cols   int
	values []driver.Value
}

func newFieldValues() *FieldValues {
	return &FieldValues{}
}

func (f *FieldValues) String() string {
	return fmt.Sprintf("rows %d columns %d", f.rows, f.cols)
}

func (f *FieldValues) resize(rows, cols int) {
	f.rows, f.cols = rows, cols
	f.values = make([]driver.Value, rows*cols)
}

// NumRow returns the number of rows available in FieldValues.
func (f *FieldValues) NumRow() int {
	return f.rows
}

// Row fills the dest value slice with row data at index idx.
func (f *FieldValues) Row(idx int, dest []driver.Value) {
	copy(dest, f.values[idx*f.cols:(idx+1)*f.cols])
}

const (
	tinyintFieldSize       = 1
	smallintFieldSize      = 2
	intFieldSize           = 4
	bigintFieldSize        = 8
	realFieldSize          = 4
	doubleFieldSize        = 8
	dateFieldSize          = 4
	timeFieldSize          = 4
	timestampFieldSize     = dateFieldSize + timeFieldSize
	longdateFieldSize      = 8
	seconddateFieldSize    = 8
	daydateFieldSize       = 4
	secondtimeFieldSize    = 4
	decimalFieldSize       = 16
	lobInputDescriptorSize = 9
)

func fieldSize(tc TypeCode, arg driver.NamedValue) (int, error) {
	v := arg.Value

	if v == nil { //HDB bug: secondtime null value --> see writeField
		return 0, nil
	}

	switch tc {
	case tcTinyint:
		return tinyintFieldSize, nil
	case tcSmallint:
		return smallintFieldSize, nil
	case tcInteger:
		return intFieldSize, nil
	case tcBigint:
		return bigintFieldSize, nil
	case tcReal:
		return realFieldSize, nil
	case tcDouble:
		return doubleFieldSize, nil
	case tcDate:
		return dateFieldSize, nil
	case tcTime:
		return timeFieldSize, nil
	case tcTimestamp:
		return timestampFieldSize, nil
	case tcLongdate:
		return longdateFieldSize, nil
	case tcSeconddate:
		return seconddateFieldSize, nil
	case tcDaydate:
		return daydateFieldSize, nil
	case tcSecondtime:
		return secondtimeFieldSize, nil
	case tcDecimal:
		return decimalFieldSize, nil
	case tcChar, tcVarchar, tcString:
		switch v := v.(type) {
		case []byte:
			return bytesSize(len(v))
		case string:
			return bytesSize(len(v))
		default:
			outLogger.Fatalf("data type %s mismatch %T", tc, v)
		}
	case tcNchar, tcNvarchar, tcNstring:
		switch v := v.(type) {
		case []byte:
			return bytesSize(cesu8.Size(v))
		case string:
			return bytesSize(cesu8.StringSize(v))
		default:
			outLogger.Fatalf("data type %s mismatch %T", tc, v)
		}
	case tcBinary, tcVarbinary:
		v, ok := v.([]byte)
		if !ok {
			outLogger.Fatalf("data type %s mismatch %T", tc, v)
		}
		return bytesSize(len(v))
	case tcBlob, tcClob, tcNclob:
		return lobInputDescriptorSize, nil
	}
	outLogger.Fatalf("data type %s not implemented", tc)
	return 0, nil
}

func readField(session *Session, rd *bufio.Reader, tc TypeCode) (interface{}, error) {

	switch tc {

	case tcTinyint, tcSmallint, tcInteger, tcBigint:

		if !rd.ReadBool() { //null value
			return nil, nil
		}

		switch tc {
		case tcTinyint:
			return int64(rd.ReadB()), nil
		case tcSmallint:
			return int64(rd.ReadInt16()), nil
		case tcInteger:
			return int64(rd.ReadInt32()), nil
		case tcBigint:
			return rd.ReadInt64(), nil
		}

	case tcReal:
		v := rd.ReadUint32()
		if v == realNullValue {
			return nil, nil
		}
		return float64(math.Float32frombits(v)), nil

	case tcDouble:
		v := rd.ReadUint64()
		if v == doubleNullValue {
			return nil, nil
		}
		return math.Float64frombits(v), nil

	case tcDate:
		year, month, day, null := readDate(rd)
		if null {
			return nil, nil
		}
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil

	// time read gives only seconds (cut), no milliseconds
	case tcTime:
		hour, minute, nanosecs, null := readTime(rd)
		if null {
			return nil, nil
		}
		return time.Date(1, 1, 1, hour, minute, 0, nanosecs, time.UTC), nil

	case tcTimestamp:
		year, month, day, dateNull := readDate(rd)
		hour, minute, nanosecs, timeNull := readTime(rd)
		if dateNull || timeNull {
			return nil, nil
		}
		return time.Date(year, month, day, hour, minute, 0, nanosecs, time.UTC), nil

	case tcLongdate:
		time, null := readLongdate(rd)
		if null {
			return nil, nil
		}
		return time, nil

	case tcSeconddate:
		time, null := readSeconddate(rd)
		if null {
			return nil, nil
		}
		return time, nil

	case tcDaydate:
		time, null := readDaydate(rd)
		if null {
			return nil, nil
		}
		return time, nil

	case tcSecondtime:
		time, null := readSecondtime(rd)
		if null {
			return nil, nil
		}
		return time, nil

	case tcDecimal:
		b, null := readDecimal(rd)
		if null {
			return nil, nil
		}
		return b, nil

	case tcChar, tcVarchar:
		value, null := readBytes(rd)
		if null {
			return nil, nil
		}
		return value, nil

	case tcNchar, tcNvarchar:
		value, null := readUtf8(rd)
		if null {
			return nil, nil
		}
		return value, nil

	case tcBinary, tcVarbinary:
		value, null := readBytes(rd)
		if null {
			return nil, nil
		}
		return value, nil

	case tcBlob, tcClob, tcNclob:
		null, writer, err := readLob(session, rd, tc)
		if null {
			return nil, nil
		}
		return writer, err
	}

	outLogger.Fatalf("read field: type code %s not implemented", tc)
	return nil, nil
}

func writeField(wr *bufio.Writer, tc TypeCode, arg driver.NamedValue) error {
	v := arg.Value
	//HDB bug: secondtime null value cannot be set by setting high byte
	//         trying so, gives
	//         SQL HdbError 1033 - error while parsing protocol: no such data type: type_code=192, index=2

	// null value
	//if v == nil && tc != tcSecondtime
	if v == nil {
		wr.WriteB(byte(tc) | 0x80) //set high bit
		return nil
	}

	// type code
	wr.WriteB(byte(tc))

	switch tc {

	default:
		outLogger.Fatalf("write field: type code %s not implemented", tc)

	case tcTinyint, tcSmallint, tcInteger, tcBigint:
		var i64 int64

		switch v := v.(type) {
		default:
			return fmt.Errorf("invalid argument type %T", v)

		case bool:
			if v {
				i64 = 1
			} else {
				i64 = 0
			}
		case int64:
			i64 = v
		}

		switch tc {
		case tcTinyint:
			wr.WriteB(byte(i64))
		case tcSmallint:
			wr.WriteInt16(int16(i64))
		case tcInteger:
			wr.WriteInt32(int32(i64))
		case tcBigint:
			wr.WriteInt64(i64)
		}

	case tcReal:

		f64, ok := v.(float64)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		wr.WriteFloat32(float32(f64))

	case tcDouble:

		f64, ok := v.(float64)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		wr.WriteFloat64(f64)

	case tcDate:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeDate(wr, t)

	case tcTime:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeTime(wr, t)

	case tcTimestamp:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeDate(wr, t)
		writeTime(wr, t)

	case tcLongdate:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeLongdate(wr, t)

	case tcSeconddate:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeSeconddate(wr, t)

	case tcDaydate:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeDaydate(wr, t)

	case tcSecondtime:
		// HDB bug: write null value explicite
		if v == nil {
			wr.WriteInt32(86401)
			return nil
		}
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeSecondtime(wr, t)

	case tcDecimal:
		b, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		if len(b) != 16 {
			return fmt.Errorf("invalid argument length %d of type %T - expected %d", len(b), v, 16)
		}
		wr.Write(b)

	case tcChar, tcVarchar, tcString:
		switch v := v.(type) {
		case []byte:
			writeBytes(wr, v)
		case string:
			writeString(wr, v)
		default:
			return fmt.Errorf("invalid argument type %T", v)
		}

	case tcNchar, tcNvarchar, tcNstring:
		switch v := v.(type) {
		case []byte:
			writeUtf8Bytes(wr, v)
		case string:
			writeUtf8String(wr, v)
		default:
			return fmt.Errorf("invalid argument type %T", v)
		}

	case tcBinary, tcVarbinary:
		v, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		writeBytes(wr, v)

	case tcBlob, tcClob, tcNclob:
		writeLob(wr)
	}

	return nil
}

// null values: most sig bit unset
// year: unset second most sig bit (subtract 2^15)
// --> read year as unsigned
// month is 0-based
// day is 1 byte
func readDate(rd *bufio.Reader) (int, time.Month, int, bool) {
	year := rd.ReadUint16()
	null := ((year & 0x8000) == 0) //null value
	year &= 0x3fff
	month := rd.ReadInt8()
	month++
	day := rd.ReadInt8()
	return int(year), time.Month(month), int(day), null
}

// year: set most sig bit
// month 0 based
func writeDate(wr *bufio.Writer, t time.Time) {
	//store in utc
	utc := t.In(time.UTC)

	year, month, day := utc.Date()

	wr.WriteUint16(uint16(year) | 0x8000)
	wr.WriteInt8(int8(month) - 1)
	wr.WriteInt8(int8(day))
}

func readTime(rd *bufio.Reader) (int, int, int, bool) {
	hour := rd.ReadB()
	null := (hour & 0x80) == 0 //null value
	hour &= 0x7f
	minute := rd.ReadInt8()
	millisecs := rd.ReadUint16()
	nanosecs := int(millisecs) * 1000000
	return int(hour), int(minute), nanosecs, null
}

func writeTime(wr *bufio.Writer, t time.Time) {
	//store in utc
	utc := t.UTC()

	wr.WriteB(byte(utc.Hour()) | 0x80)
	wr.WriteInt8(int8(utc.Minute()))
	millisecs := utc.Second()*1000 + utc.Round(time.Millisecond).Nanosecond()/1000000
	wr.WriteUint16(uint16(millisecs))
}

var zeroTime = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

func readLongdate(rd *bufio.Reader) (time.Time, bool) {
	longdate := rd.ReadInt64()
	if longdate == 3155380704000000001 { // null value
		return zeroTime, true
	}
	return convertLongdateToTime(longdate), false
}

func writeLongdate(wr *bufio.Writer, t time.Time) {
	wr.WriteInt64(convertTimeToLongdate(t))
}

func readSeconddate(rd *bufio.Reader) (time.Time, bool) {
	seconddate := rd.ReadInt64()
	if seconddate == 315538070401 { // null value
		return zeroTime, true
	}
	return convertSeconddateToTime(seconddate), false
}

func writeSeconddate(wr *bufio.Writer, t time.Time) {
	wr.WriteInt64(convertTimeToSeconddate(t))
}

func readDaydate(rd *bufio.Reader) (time.Time, bool) {
	daydate := rd.ReadInt32()
	if daydate == 3652062 { // null value
		return zeroTime, true
	}
	return convertDaydateToTime(int64(daydate)), false
}

func writeDaydate(wr *bufio.Writer, t time.Time) {
	wr.WriteInt32(int32(convertTimeToDayDate(t)))
}

func readSecondtime(rd *bufio.Reader) (time.Time, bool) {
	secondtime := rd.ReadInt32()
	if secondtime == 86401 { // null value
		return zeroTime, true
	}
	return convertSecondtimeToTime(int(secondtime)), false
}

func writeSecondtime(wr *bufio.Writer, t time.Time) {
	wr.WriteInt32(int32(convertTimeToSecondtime(t)))
}

// nanosecond: HDB - 7 digits precision (not 9 digits)
func convertTimeToLongdate(t time.Time) int64 {
	t = t.UTC()
	return (((((((int64(convertTimeToDayDate(t))-1)*24)+int64(t.Hour()))*60)+int64(t.Minute()))*60)+int64(t.Second()))*10000000 + int64(t.Nanosecond()/100) + 1
}

func convertLongdateToTime(longdate int64) time.Time {
	const dayfactor = 10000000 * 24 * 60 * 60
	longdate--
	d := (longdate % dayfactor) * 100
	t := convertDaydateToTime((longdate / dayfactor) + 1)
	return t.Add(time.Duration(d))
}

func convertTimeToSeconddate(t time.Time) int64 {
	t = t.UTC()
	return (((((int64(convertTimeToDayDate(t))-1)*24)+int64(t.Hour()))*60)+int64(t.Minute()))*60 + int64(t.Second()) + 1
}

func convertSeconddateToTime(seconddate int64) time.Time {
	const dayfactor = 24 * 60 * 60
	seconddate--
	d := (seconddate % dayfactor) * 1000000000
	t := convertDaydateToTime((seconddate / dayfactor) + 1)
	return t.Add(time.Duration(d))
}

const julianHdb = 1721423 // 1 January 0001 00:00:00 (1721424) - 1

func convertTimeToDayDate(t time.Time) int64 {
	return int64(timeToJulianDay(t) - julianHdb)
}

func convertDaydateToTime(daydate int64) time.Time {
	return julianDayToTime(int(daydate) + julianHdb)
}

func convertTimeToSecondtime(t time.Time) int {
	t = t.UTC()
	return (t.Hour()*60+t.Minute())*60 + t.Second() + 1
}

func convertSecondtimeToTime(secondtime int) time.Time {
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(int64(secondtime-1) * 1000000000))
}

func readDecimal(rd *bufio.Reader) ([]byte, bool) {
	b := make([]byte, 16)
	rd.ReadFull(b)
	if (b[15] & 0x70) == 0x70 { //null value (bit 4,5,6 set)
		return nil, true
	}
	return b, false
}

// string / binary length indicators
const (
	bytesLenIndNullValue byte = 255
	bytesLenIndSmall     byte = 245
	bytesLenIndMedium    byte = 246
	bytesLenIndBig       byte = 247
)

func bytesSize(size int) (int, error) { //size + length indicator
	switch {
	default:
		return 0, fmt.Errorf("max string length %d exceeded %d", math.MaxInt32, size)
	case size <= int(bytesLenIndSmall):
		return size + 1, nil
	case size <= math.MaxInt16:
		return size + 3, nil
	case size <= math.MaxInt32:
		return size + 5, nil
	}
}

func readBytesSize(rd *bufio.Reader) (int, bool) {

	ind := rd.ReadB() //length indicator

	switch {

	default:
		return 0, false

	case ind == bytesLenIndNullValue:
		return 0, true

	case ind <= bytesLenIndSmall:
		return int(ind), false

	case ind == bytesLenIndMedium:
		return int(rd.ReadInt16()), false

	case ind == bytesLenIndBig:
		return int(rd.ReadInt32()), false

	}
}

func writeBytesSize(wr *bufio.Writer, size int) error {
	switch {

	default:
		return fmt.Errorf("max argument length %d of string exceeded", size)

	case size <= int(bytesLenIndSmall):
		wr.WriteB(byte(size))
	case size <= math.MaxInt16:
		wr.WriteB(bytesLenIndMedium)
		wr.WriteInt16(int16(size))
	case size <= math.MaxInt32:
		wr.WriteB(bytesLenIndBig)
		wr.WriteInt32(int32(size))
	}
	return nil
}

func readBytes(rd *bufio.Reader) ([]byte, bool) {
	size, null := readBytesSize(rd)
	if null {
		return nil, true
	}
	b := make([]byte, size)
	rd.ReadFull(b)
	return b, false
}

func readUtf8(rd *bufio.Reader) ([]byte, bool) {
	size, null := readBytesSize(rd)
	if null {
		return nil, true
	}
	b := rd.ReadCesu8(size)
	return b, false
}

// strings with one byte length
func readShortUtf8(rd *bufio.Reader) ([]byte, int) {
	size := rd.ReadB()
	b := rd.ReadCesu8(int(size))
	return b, int(size)
}

func writeBytes(wr *bufio.Writer, b []byte) {
	writeBytesSize(wr, len(b))
	wr.Write(b)
}

func writeString(wr *bufio.Writer, s string) {
	writeBytesSize(wr, len(s))
	wr.WriteString(s)
}

func writeUtf8Bytes(wr *bufio.Writer, b []byte) {
	size := cesu8.Size(b)
	writeBytesSize(wr, size)
	wr.WriteCesu8(b)
}

func writeUtf8String(wr *bufio.Writer, s string) {
	size := cesu8.StringSize(s)
	writeBytesSize(wr, size)
	wr.WriteStringCesu8(s)
}

func readLob(s *Session, rd *bufio.Reader, tc TypeCode) (bool, lobChunkWriter, error) {
	rd.ReadInt8() // type code (is int here)
	opt := rd.ReadInt8()
	null := (lobOptions(opt) & loNullindicator) != 0
	if null {
		return true, nil, nil
	}
	eof := (lobOptions(opt) & loLastdata) != 0
	rd.Skip(2)

	charLen := rd.ReadInt64()
	byteLen := rd.ReadInt64()
	id := rd.ReadUint64()
	chunkLen := rd.ReadInt32()

	lobChunkWriter := newLobChunkWriter(tc.isCharBased(), s, locatorID(id), charLen, byteLen)
	if err := lobChunkWriter.write(rd, int(chunkLen), eof); err != nil {
		return null, lobChunkWriter, err
	}
	return null, lobChunkWriter, nil
}

// TODO: first write: add content? - actually no data transferred
func writeLob(wr *bufio.Writer) {
	wr.WriteB(0)
	wr.WriteInt32(0)
	wr.WriteInt32(0)
}
