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

const (
	realNullValue   uint32 = ^uint32(0)
	doubleNullValue uint64 = ^uint64(0)
)

type uint32Slice []uint32

func (p uint32Slice) Len() int           { return len(p) }
func (p uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p uint32Slice) sort()              { sort.Sort(p) }

type field interface {
	typeCode() typeCode
	in() bool
	out() bool
	name(map[uint32]string) string
	nameOffsets() []uint32
	String() string
}

// FieldSet contains database field metadata.
type FieldSet struct {
	fields []field
	names  map[uint32]string
}

func newFieldSet(size int) *FieldSet {
	return &FieldSet{
		fields: make([]field, size),
		names:  make(map[uint32]string),
	}
}

// String implements the Stringer interface.
func (f *FieldSet) String() string {
	a := make([]string, len(f.fields))
	for i, f := range f.fields {
		a[i] = f.String()
	}
	return fmt.Sprintf("%v", a)
}

func (f *FieldSet) nameOffsets() []uint32 {
	for _, field := range f.fields {
		for _, offset := range field.nameOffsets() {
			if offset != 0xFFFFFFFF {
				f.names[offset] = ""
			}
		}
	}
	// sort offsets (not sure if offsets are monotonically increasing in any case)
	offsets := make([]uint32, len(f.names))
	i := 0
	for offset := range f.names {
		offsets[i] = offset
		i++
	}
	uint32Slice(offsets).sort()
	return offsets
}

// NumInputField returns the number of input fields in a database statement.
func (f *FieldSet) NumInputField() int {
	cnt := 0
	for _, field := range f.fields {
		if field.in() {
			cnt++
		}
	}
	return cnt
}

// NumOutputField returns the number of output fields of a query or stored procedure.
func (f *FieldSet) NumOutputField() int {
	cnt := 0
	for _, field := range f.fields {
		if field.out() {
			cnt++
		}
	}
	return cnt
}

// DataType returns the datatype of the field at index idx.
func (f *FieldSet) DataType(idx int) DataType {
	return f.fields[idx].typeCode().dataType()
}

// OutputNames fills the names parameter with field names of all output fields. The size of the names slice must be at least
// NumOutputField big.
func (f *FieldSet) OutputNames(names []string) error {
	i := 0
	for _, field := range f.fields {
		if field.out() {
			if i >= len(names) { // assert names size
				return fmt.Errorf("names size too short %d - expected min %d", len(names), i)
			}
			names[i] = field.name(f.names)
			i++
		}
	}
	return nil
}

// FieldValues contains rows read from database.
type FieldValues struct {
	s *Session

	rows    int
	cols    int
	lobCols int
	values  []driver.Value

	descrs  []*LobReadDescr // Caution: store descriptor to guarantee valid addresses
	writers []lobWriter
}

func newFieldValues(s *Session) *FieldValues {
	return &FieldValues{s: s}
}

func (f *FieldValues) String() string {
	return fmt.Sprintf("rows %d columns %d lob columns %d", f.rows, f.cols, f.lobCols)
}

func (f *FieldValues) read(rows int, fieldSet *FieldSet, rd *bufio.Reader) error {
	f.rows = rows
	f.descrs = make([]*LobReadDescr, 0)

	f.cols, f.lobCols = 0, 0
	for _, field := range fieldSet.fields {
		if field.out() {
			if field.typeCode().isLob() {
				f.descrs = append(f.descrs, &LobReadDescr{col: f.cols})
				f.lobCols++
			}
			f.cols++
		}
	}
	f.values = make([]driver.Value, f.rows*f.cols)
	f.writers = make([]lobWriter, f.lobCols)

	for i := 0; i < f.rows; i++ {
		j := 0
		for _, field := range fieldSet.fields {

			if !field.out() {
				continue
			}

			var err error
			if f.values[i*f.cols+j], err = readField(rd, field.typeCode()); err != nil {
				return err
			}
			j++
		}
	}
	return nil
}

// NumRow returns the number of rows available in FieldValues.
func (f *FieldValues) NumRow() int {
	return f.rows
}

// Row fills the dest value slice with row data at index idx.
func (f *FieldValues) Row(idx int, dest []driver.Value) {
	copy(dest, f.values[idx*f.cols:(idx+1)*f.cols])

	if f.lobCols == 0 {
		return
	}

	for i, descr := range f.descrs {
		col := descr.col
		writer := dest[col].(lobWriter)
		f.writers[i] = writer
		descr.w = writer
		dest[col] = lobReadDescrToPointer(descr)
	}

	// last descriptor triggers lob read
	f.descrs[f.lobCols-1].fn = func() error {
		return f.s.readLobStream(f.writers)
	}
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
	decimalFieldSize       = 16
	lobInputDescriptorSize = 9
)

func fieldSize(tc typeCode, v driver.Value) (int, error) {

	if v == nil {
		return 0, nil
	}

	switch tc {
	case tcTinyint:
		return tinyintFieldSize, nil
	case tcSmallint:
		return smallintFieldSize, nil
	case tcInt:
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

func readField(rd *bufio.Reader, tc typeCode) (interface{}, error) {

	switch tc {

	case tcTinyint, tcSmallint, tcInt, tcBigint:

		valid, err := rd.ReadBool()
		if err != nil {
			return nil, err
		}
		if !valid { //null value
			return nil, nil
		}

		switch tc {

		case tcTinyint:
			if v, err := rd.ReadByte(); err == nil {
				return int64(v), nil
			}
			return nil, err

		case tcSmallint:
			if v, err := rd.ReadInt16(); err == nil {
				return int64(v), nil
			}
			return nil, err

		case tcInt:
			if v, err := rd.ReadInt32(); err == nil {
				return int64(v), nil
			}
			return nil, err

		case tcBigint:
			if v, err := rd.ReadInt64(); err == nil {
				return v, nil
			}
			return nil, err
		}

	case tcReal:
		v, err := rd.ReadUint32()
		if err != nil {
			return nil, err
		}
		if v == realNullValue {
			return nil, nil
		}
		return float64(math.Float32frombits(v)), nil

	case tcDouble:
		v, err := rd.ReadUint64()
		if err != nil {
			return nil, err
		}
		if v == doubleNullValue {
			return nil, nil
		}
		return math.Float64frombits(v), nil

	case tcDate:

		year, month, day, null, err := readDate(rd)
		if err != nil {
			return nil, err
		}
		if null {
			return nil, nil
		}

		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil

	// time read gives only seconds (cut), no milliseconds
	case tcTime:

		hour, minute, nanosecs, null, err := readTime(rd)
		if err != nil {
			return nil, err
		}
		if null {
			return nil, nil
		}

		return time.Date(1, 1, 1, hour, minute, 0, nanosecs, time.UTC), nil

	case tcTimestamp:

		year, month, day, dateNull, err := readDate(rd)
		if err != nil {
			return nil, err
		}

		hour, minute, nanosecs, timeNull, err := readTime(rd)
		if err != nil {
			return nil, err
		}

		if dateNull || timeNull {
			return nil, nil
		}

		return time.Date(year, month, day, hour, minute, 0, nanosecs, time.UTC), nil

	case tcDecimal:

		b, null, err := readDecimal(rd)
		switch {
		case err != nil:
			return nil, err
		case null:
			return nil, nil
		default:
			return b, nil
		}

	case tcChar, tcVarchar:
		value, null, err := readBytes(rd)
		if err != nil {
			return nil, err
		}
		if null {
			return nil, nil
		}
		return value, nil

	case tcNchar, tcNvarchar:
		value, null, err := readUtf8(rd)
		if err != nil {
			return nil, err
		}
		if null {
			return nil, nil
		}
		return value, nil

	case tcBinary, tcVarbinary:
		value, null, err := readBytes(rd)
		if err != nil {
			return nil, err
		}
		if null {
			return nil, nil
		}
		return value, nil

	case tcBlob, tcClob, tcNclob:
		null, writer, err := readLob(rd, tc)
		if err != nil {
			return nil, err
		}
		if null {
			return nil, nil
		}
		return writer, nil
	}

	outLogger.Fatalf("read field: type code %s not implemented", tc)
	return nil, nil
}

func writeField(wr *bufio.Writer, tc typeCode, v driver.Value) error {

	// null value
	if v == nil {
		if err := wr.WriteByte(byte(tc) | 0x80); err != nil { //set high bit
			return err
		}
		return nil
	}

	// type code
	if err := wr.WriteByte(byte(tc)); err != nil {
		return err
	}

	switch tc {

	// TODO: char, ...

	case tcTinyint, tcSmallint, tcInt, tcBigint:

		i64, ok := v.(int64)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}

		switch tc {
		case tcTinyint:
			return wr.WriteByte(byte(i64))
		case tcSmallint:
			return wr.WriteInt16(int16(i64))
		case tcInt:
			return wr.WriteInt32(int32(i64))
		case tcBigint:
			return wr.WriteInt64(i64)
		}

	case tcReal:

		f64, ok := v.(float64)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		return wr.WriteFloat32(float32(f64))

	case tcDouble:

		f64, ok := v.(float64)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		return wr.WriteFloat64(f64)

	case tcDate:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		return writeDate(wr, t)

	case tcTime:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		return writeTime(wr, t)

	case tcTimestamp:
		t, ok := v.(time.Time)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		if err := writeDate(wr, t); err != nil {
			return err
		}
		return writeTime(wr, t)

	case tcDecimal:
		b, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		if len(b) != 16 {
			return fmt.Errorf("invalid argument length %d of type %T - expected %d", len(b), v, 16)
		}
		_, err := wr.Write(b)
		return err

	case tcChar, tcVarchar, tcString:
		switch v := v.(type) {
		case []byte:
			return writeBytes(wr, v)
		case string:
			return writeString(wr, v)
		default:
			return fmt.Errorf("invalid argument type %T", v)
		}

	case tcNchar, tcNvarchar, tcNstring:
		switch v := v.(type) {
		case []byte:
			return writeUtf8Bytes(wr, v)
		case string:
			return writeUtf8String(wr, v)
		default:
			return fmt.Errorf("invalid argument type %T", v)
		}

	case tcBinary, tcVarbinary:
		v, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("invalid argument type %T", v)
		}
		return writeBytes(wr, v)

	case tcBlob, tcClob, tcNclob:
		return writeLob(wr)
	}

	outLogger.Fatalf("write field: type code %s not implemented", tc)
	return nil
}

// null values: most sig bit unset
// year: unset second most sig bit (subtract 2^15)
// --> read year as unsigned
// month is 0-based
// day is 1 byte
func readDate(rd *bufio.Reader) (int, time.Month, int, bool, error) {

	year, err := rd.ReadUint16()
	if err != nil {
		return 0, 0, 0, false, err
	}
	if (year & 0x8000) == 0 { //null value
		if err := rd.Skip(2); err != nil {
			return 0, 0, 0, false, err
		}
		return 0, 0, 0, true, nil
	}
	year &= 0x3fff
	month, err := rd.ReadInt8()
	if err != nil {
		return 0, 0, 0, false, err
	}
	month++
	day, err := rd.ReadInt8()
	if err != nil {
		return 0, 0, 0, false, err
	}
	return int(year), time.Month(month), int(day), false, nil
}

// year: set most sig bit
// month 0 based
func writeDate(wr *bufio.Writer, t time.Time) error {

	//store in utc
	utc := t.In(time.UTC)

	year, month, day := utc.Date()

	if err := wr.WriteUint16(uint16(year) | 0x8000); err != nil {
		return err
	}
	if err := wr.WriteInt8(int8(month) - 1); err != nil {
		return err
	}
	if err := wr.WriteInt8(int8(day)); err != nil {
		return err
	}
	return nil
}

func readTime(rd *bufio.Reader) (int, int, int, bool, error) {

	hour, err := rd.ReadByte()
	if err != nil {
		return 0, 0, 0, false, err
	}
	if (hour & 0x80) == 0 { //null value
		if err := rd.Skip(3); err != nil {
			return 0, 0, 0, false, err
		}
		return 0, 0, 0, true, nil
	}
	hour &= 0x7f
	minute, err := rd.ReadInt8()
	if err != nil {
		return 0, 0, 0, false, err
	}
	millisecs, err := rd.ReadUint16()
	if err != nil {
		return 0, 0, 0, false, err
	}

	nanosecs := int(millisecs) * 1000000

	return int(hour), int(minute), nanosecs, false, nil
}

func writeTime(wr *bufio.Writer, t time.Time) error {

	//store in utc
	utc := t.UTC()

	if err := wr.WriteByte(byte(utc.Hour()) | 0x80); err != nil {
		return err
	}
	if err := wr.WriteInt8(int8(utc.Minute())); err != nil {
		return err
	}

	millisecs := utc.Second()*1000 + utc.Round(time.Millisecond).Nanosecond()/1000000

	if err := wr.WriteUint16(uint16(millisecs)); err != nil {
		return err
	}

	return nil
}

func readDecimal(rd *bufio.Reader) ([]byte, bool, error) {
	b := make([]byte, 16)
	if err := rd.ReadFull(b); err != nil {
		return nil, false, err
	}
	if (b[15] & 0x70) == 0x70 { //null value (bit 4,5,6 set)
		return nil, true, nil
	}
	return b, false, nil
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

func readBytesSize(rd *bufio.Reader) (int, bool, error) {

	ind, err := rd.ReadByte() //length indicator
	if err != nil {
		return 0, false, err
	}

	switch {

	default:
		return 0, false, fmt.Errorf("invalid length indicator %d", ind)

	case ind == bytesLenIndNullValue:
		return 0, true, nil

	case ind <= bytesLenIndSmall:
		return int(ind), false, nil

	case ind == bytesLenIndMedium:
		if size, err := rd.ReadInt16(); err == nil {
			return int(size), false, nil
		}
		return 0, false, err

	case ind == bytesLenIndBig:
		if size, err := rd.ReadInt32(); err == nil {
			return int(size), false, nil
		}
		return 0, false, err
	}
}

func writeBytesSize(wr *bufio.Writer, size int) error {
	switch {

	default:
		return fmt.Errorf("max argument length %d of string exceeded", size)

	case size <= int(bytesLenIndSmall):
		if err := wr.WriteByte(byte(size)); err != nil {
			return err
		}
	case size <= math.MaxInt16:
		if err := wr.WriteByte(bytesLenIndMedium); err != nil {
			return err
		}
		if err := wr.WriteInt16(int16(size)); err != nil {
			return err
		}
	case size <= math.MaxInt32:
		if err := wr.WriteByte(bytesLenIndBig); err != nil {
			return err
		}
		if err := wr.WriteInt32(int32(size)); err != nil {
			return err
		}
	}
	return nil
}

func readBytes(rd *bufio.Reader) ([]byte, bool, error) {
	size, null, err := readBytesSize(rd)
	if err != nil {
		return nil, false, err
	}

	if null {
		return nil, true, nil
	}

	b := make([]byte, size)
	if err := rd.ReadFull(b); err != nil {
		return nil, false, err
	}
	return b, false, nil
}

func readUtf8(rd *bufio.Reader) ([]byte, bool, error) {
	size, null, err := readBytesSize(rd)
	if err != nil {
		return nil, false, err
	}

	if null {
		return nil, true, nil
	}

	b, err := rd.ReadCesu8(size)
	if err != nil {
		return nil, false, err
	}

	return b, false, nil
}

// strings with one byte length
func readShortUtf8(rd *bufio.Reader) ([]byte, int, error) {
	size, err := rd.ReadByte()
	if err != nil {
		return nil, 0, err
	}

	b, err := rd.ReadCesu8(int(size))
	if err != nil {
		return nil, 0, err
	}

	return b, int(size), nil
}

func writeBytes(wr *bufio.Writer, b []byte) error {
	if err := writeBytesSize(wr, len(b)); err != nil {
		return err
	}
	_, err := wr.Write(b)
	return err
}

func writeString(wr *bufio.Writer, s string) error {
	if err := writeBytesSize(wr, len(s)); err != nil {
		return err
	}
	_, err := wr.WriteString(s)
	return err
}

func writeUtf8Bytes(wr *bufio.Writer, b []byte) error {
	size := cesu8.Size(b)
	if err := writeBytesSize(wr, size); err != nil {
		return err
	}
	_, err := wr.WriteCesu8(b)
	return err
}

func writeUtf8String(wr *bufio.Writer, s string) error {
	size := cesu8.StringSize(s)
	if err := writeBytesSize(wr, size); err != nil {
		return err
	}
	_, err := wr.WriteStringCesu8(s)
	return err
}

func readLob(rd *bufio.Reader, tc typeCode) (bool, lobWriter, error) {

	if _, err := rd.ReadInt8(); err != nil { // type code (is int here)
		return false, nil, err
	}

	opt, err := rd.ReadInt8()
	if err != nil {
		return false, nil, err
	}

	if err := rd.Skip(2); err != nil {
		return false, nil, err
	}

	charLen, err := rd.ReadInt64()
	if err != nil {
		return false, nil, err
	}
	byteLen, err := rd.ReadInt64()
	if err != nil {
		return false, nil, err
	}
	id, err := rd.ReadUint64()
	if err != nil {
		return false, nil, err
	}
	chunkLen, err := rd.ReadInt32()
	if err != nil {
		return false, nil, err
	}

	null := (lobOptions(opt) & loNullindicator) != 0
	eof := (lobOptions(opt) & loLastdata) != 0

	var writer lobWriter
	if tc.isCharBased() {
		writer = newCharLobWriter(locatorID(id), charLen, byteLen)
	} else {
		writer = newBinaryLobWriter(locatorID(id), charLen, byteLen)
	}
	if err := writer.write(rd, int(chunkLen), eof); err != nil {
		return null, writer, err
	}
	return null, writer, nil
}

// TODO: first write: add content? - actually no data transferred
func writeLob(wr *bufio.Writer) error {

	if err := wr.WriteByte(0); err != nil {
		return err
	}
	if err := wr.WriteInt32(0); err != nil {
		return err
	}
	if err := wr.WriteInt32(0); err != nil {
		return err
	}
	return nil
}
