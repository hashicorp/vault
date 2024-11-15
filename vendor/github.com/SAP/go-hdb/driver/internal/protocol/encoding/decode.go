package encoding

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"time"

	"github.com/SAP/go-hdb/driver/internal/unsafe"
	"golang.org/x/text/transform"
)

const readScratchSize = 4096

// Decoder decodes hdb protocol datatypes on basis of an io.Reader.
type Decoder struct {
	rd io.Reader
	/* err: fatal read error
	- not set by conversion errors
	- conversion errors are returned by the reader function itself
	*/
	err error
	b   []byte // scratch buffer (used for skip, CESU8Bytes - define size not too small!)
	tr  transform.Transformer
	cnt int

	// decoder options
	alphanumDfv1    bool
	emptyDateAsNull bool
}

// NewDecoder creates a new Decoder instance based on an io.Reader.
func NewDecoder(rd io.Reader, tr transform.Transformer, emptyDateAsNull bool) *Decoder {
	return &Decoder{
		rd:              rd,
		b:               make([]byte, readScratchSize),
		tr:              tr,
		emptyDateAsNull: emptyDateAsNull,
	}
}

// Transformer returns the cesu8 transformer.
func (d *Decoder) Transformer() transform.Transformer { return d.tr }

// SetAlphanumDfv1 sets the alphanum dfv1 flag decoder.
func (d *Decoder) SetAlphanumDfv1(alphanumDfv1 bool) { d.alphanumDfv1 = alphanumDfv1 }

// Cnt returns the value of the byte read counter.
func (d *Decoder) Cnt() int { return d.cnt }

// Error returns the last decoder error.
func (d *Decoder) Error() error { return d.err }

// ResetError resets reader error.
func (d *Decoder) ResetError() { d.err = nil }

// readFull reads data from reader + read counter and error handling.
func (d *Decoder) readFull(buf []byte) error {
	if d.err != nil {
		return d.err
	}
	var n int
	n, d.err = io.ReadFull(d.rd, buf)
	d.cnt += n
	return d.err
}

// Skip skips cnt bytes from reading.
func (d *Decoder) Skip(cnt int) {
	if cnt <= readScratchSize {
		d.readFull(d.b[:cnt]) //nolint: errcheck
		return
	}
	var n int64
	n, d.err = io.CopyN(io.Discard, d.rd, int64(cnt))
	d.cnt += int(n)
}

// Byte decodes a byte.
func (d *Decoder) Byte() byte {
	if err := d.readFull(d.b[:1]); err != nil {
		return 0
	}
	return d.b[0]
}

// Bytes decodes bytes.
func (d *Decoder) Bytes(p []byte) {
	d.readFull(p) //nolint:errcheck
}

// Bool decodes a boolean.
func (d *Decoder) Bool() bool {
	return d.Byte() != 0
}

// Int8 decodes an int8.
func (d *Decoder) Int8() int8 {
	return int8(d.Byte())
}

// Int16 decodes an int16.
func (d *Decoder) Int16() int16 {
	if err := d.readFull(d.b[:2]); err != nil {
		return 0
	}
	return int16(binary.LittleEndian.Uint16(d.b[:2])) //nolint: gosec
}

// Uint16 decodes an uint16.
func (d *Decoder) Uint16() uint16 {
	if err := d.readFull(d.b[:2]); err != nil {
		return 0
	}
	return binary.LittleEndian.Uint16(d.b[:2])
}

// Uint16ByteOrder decodes an uint16 in given byte order.
func (d *Decoder) Uint16ByteOrder(byteOrder binary.ByteOrder) uint16 {
	if err := d.readFull(d.b[:2]); err != nil {
		return 0
	}
	return byteOrder.Uint16(d.b[:2])
}

// Int32 decodes an int32.
func (d *Decoder) Int32() int32 {
	if err := d.readFull(d.b[:4]); err != nil {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(d.b[:4])) //nolint: gosec
}

// Uint32 decodes an uint32.
func (d *Decoder) Uint32() uint32 {
	if err := d.readFull(d.b[:4]); err != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(d.b[:4])
}

// Uint32ByteOrder decodes an uint32 in given byte order.
func (d *Decoder) Uint32ByteOrder(byteOrder binary.ByteOrder) uint32 {
	if err := d.readFull(d.b[:4]); err != nil {
		return 0
	}
	return byteOrder.Uint32(d.b[:4])
}

// Int64 decodes an int64.
func (d *Decoder) Int64() int64 {
	if err := d.readFull(d.b[:8]); err != nil {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(d.b[:8])) //nolint: gosec
}

// Uint64 decodes an uint64.
func (d *Decoder) Uint64() uint64 {
	if err := d.readFull(d.b[:8]); err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(d.b[:8])
}

// Float32 decodes a float32.
func (d *Decoder) Float32() float32 {
	if err := d.readFull(d.b[:4]); err != nil {
		return 0
	}
	bits := binary.LittleEndian.Uint32(d.b[:4])
	return math.Float32frombits(bits)
}

// Float64 decodes a float64.
func (d *Decoder) Float64() float64 {
	if err := d.readFull(d.b[:8]); err != nil {
		return 0
	}
	bits := binary.LittleEndian.Uint64(d.b[:8])
	return math.Float64frombits(bits)
}

// Decimal decodes a decimal.
// - error is only returned in case of conversion errors.
func (d *Decoder) Decimal() (*big.Int, int, error) { // m, exp
	bs := d.b[:decSize]

	if err := d.readFull(bs); err != nil {
		return nil, 0, nil
	}

	if (bs[15] & 0x70) == 0x70 { // null value (bit 4,5,6 set)
		return nil, 0, nil
	}

	if (bs[15] & 0x60) == 0x60 {
		return nil, 0, fmt.Errorf("decimal: format (infinity, nan, ...) not supported : %v", bs)
	}

	neg := (bs[15] & 0x80) != 0
	exp := int((((uint16(bs[15])<<8)|uint16(bs[14]))<<1)>>2) - dec128Bias

	// b14 := b[14]  // save b[14]
	bs[14] &= 0x01 // keep the mantissa bit (rest: sign and exp)

	// most significand byte
	msb := 14
	for msb > 0 && bs[msb] == 0 {
		msb--
	}

	// calc number of words
	numWords := (msb / _S) + 1
	ws := make([]big.Word, numWords)

	bs = bs[:msb+1]
	for i, b := range bs {
		ws[i/_S] |= (big.Word(b) << (i % _S * 8))
	}

	m := new(big.Int).SetBits(ws)
	if neg {
		m = m.Neg(m)
	}
	return m, exp, nil
}

// Fixed decodes a fixed decimal.
func (d *Decoder) Fixed(size int) *big.Int { // m, exp
	bs := d.b[:size]

	if err := d.readFull(bs); err != nil {
		return nil
	}

	neg := (bs[size-1] & 0x80) != 0 // is negative number (2s complement)

	// most significand byte
	msb := size - 1
	for msb > 0 && bs[msb] == 0 {
		msb--
	}

	// calc number of words
	numWords := (msb / _S) + 1
	ws := make([]big.Word, numWords)

	bs = bs[:msb+1]
	for i, b := range bs {
		// if negative: invert byte (2s complement)
		if neg {
			b = ^b
		}
		ws[i/_S] |= (big.Word(b) << (i % _S * 8))
	}

	m := new(big.Int).SetBits(ws)

	if neg {
		m.Add(m, natOne) // 2s complement - add 1
		m.Neg(m)         // set sign
	}
	return m
}

// CESU8Bytes decodes CESU-8 into UTF-8 bytes.
// - error is only returned in case of conversion errors.
func (d *Decoder) CESU8Bytes(size int) ([]byte, error) {
	if d.err != nil {
		return nil, nil
	}

	var p []byte
	if size > readScratchSize {
		p = make([]byte, size)
	} else {
		p = d.b[:size]
	}

	if err := d.readFull(p); err != nil {
		return nil, nil
	}

	b, _, err := transform.Bytes(d.tr, p)
	return b, err
}

// varFieldInd decodes a variable field indicator.
func (d *Decoder) varFieldInd() (n, size int, null bool) {
	ind := d.Byte() // length indicator
	switch {
	default:
		return 1, 0, false
	case ind == bytesLenIndNullValue:
		return 1, 0, true
	case ind <= bytesLenIndSmall:
		return 1, int(ind), false
	case ind == bytesLenIndMedium:
		return 3, int(d.Int16()), false
	case ind == bytesLenIndBig:
		return 5, int(d.Int32()), false
	}
}

// LIBytes decodes bytes with length indicator.
func (d *Decoder) LIBytes() (n int, b []byte) {
	n, size, null := d.varFieldInd()
	if null {
		return n, nil
	}
	b = make([]byte, size)
	d.Bytes(b)
	return n + size, b
}

// LIString decodes a string with length indicator.
func (d *Decoder) LIString() (n int, s string) {
	n, b := d.LIBytes()
	return n, unsafe.ByteSlice2String(b)
}

// CESU8LIBytes decodes CESU-8 into UTF-8 bytes with length indicator.
func (d *Decoder) CESU8LIBytes() (int, []byte, error) {
	n, size, null := d.varFieldInd()
	if null {
		return n, nil, nil
	}
	b, err := d.CESU8Bytes(size)
	return n + size, b, err
}

// CESU8LIString decodes a CESU-8 into a UTF-8 string with length indicator.
func (d *Decoder) CESU8LIString() (int, string, error) {
	n, b, err := d.CESU8LIBytes()
	return n, unsafe.ByteSlice2String(b), err
}

// Fields.

// BooleanField decodes a boolean field.
func (d *Decoder) BooleanField() (any, error) {
	b := d.Byte()
	switch b {
	case booleanNullValue:
		return nil, nil
	case booleanFalseValue:
		return false, nil
	default:
		return true, nil
	}
}

// RealField decodes a real field.
func (d *Decoder) RealField() (any, error) {
	v := d.Uint32()
	if v == realNullValue {
		return nil, nil
	}
	return float64(math.Float32frombits(v)), nil
}

// DoubleField decodes a double field.
func (d *Decoder) DoubleField() (any, error) {
	v := d.Uint64()
	if v == doubleNullValue {
		return nil, nil
	}
	return math.Float64frombits(v), nil
}

func (d *Decoder) decodeDate() (int, time.Month, int, bool) {
	// decode.
	/*
	   null values: most sig bit unset
	   year: unset second most sig bit (subtract 2^15)
	   --> read year as unsigned
	   month is 0-based
	   day is 1 byte.
	*/
	year := d.Uint16()
	null := ((year & 0x8000) == 0) // null value
	year &= 0x3fff
	month := d.Int8()
	month++
	day := d.Int8()
	return int(year), time.Month(month), int(day), null
}

// DateField decodes a date field.
func (d *Decoder) DateField() (any, error) {
	year, month, day, null := d.decodeDate()
	if null {
		return nil, nil
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

func (d *Decoder) decodeTime() (int, int, int, int, bool) {
	hour := d.Byte()
	null := (hour & 0x80) == 0 // null value
	hour &= 0x7f
	minute := d.Int8()
	msec := d.Uint16()

	sec := msec / 1000
	msec %= 1000
	nsec := int(msec) * 1000000

	return int(hour), int(minute), int(sec), nsec, null
}

// TimeField decodes a time field.
func (d *Decoder) TimeField() (any, error) {
	// time read gives only seconds (cut), no milliseconds
	hour, minute, sec, nsec, null := d.decodeTime()
	if null {
		return nil, nil
	}
	return time.Date(1, 1, 1, hour, minute, sec, nsec, time.UTC), nil
}

// TimestampField decodes a timestamp field.
func (d *Decoder) TimestampField() (any, error) {
	year, month, day, dateNull := d.decodeDate()
	hour, minute, sec, nsec, timeNull := d.decodeTime()
	if dateNull || timeNull {
		return nil, nil
	}
	return time.Date(year, month, day, hour, minute, sec, nsec, time.UTC), nil
}

// LongdateField decodes a longdate field.
func (d *Decoder) LongdateField() (any, error) {
	longdate := d.Int64()
	if longdate == longdateNullValue {
		return nil, nil
	}
	return convertLongdateToTime(longdate), nil
}

// SeconddateField decodes a seconddate field.
func (d *Decoder) SeconddateField() (any, error) {
	seconddate := d.Int64()
	if seconddate == seconddateNullValue {
		return nil, nil
	}
	return convertSeconddateToTime(seconddate), nil
}

// DaydateField decodes a daydate field.
func (d *Decoder) DaydateField() (any, error) {
	daydate := d.Int32()
	if daydate == daydateNullValue || (d.emptyDateAsNull && daydate == 0) {
		return nil, nil
	}
	return convertDaydateToTime(int64(daydate)), nil
}

// SecondtimeField decodes a secondtime field.
func (d *Decoder) SecondtimeField() (any, error) {
	secondtime := d.Int32()
	if secondtime == secondtimeNullValue {
		return nil, nil
	}
	return convertSecondtimeToTime(int(secondtime)), nil
}

// DecimalField decodes a decimal field.
func (d *Decoder) DecimalField() (any, error) {
	m, exp, err := d.Decimal()
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return convertDecimalToRat(m, exp), nil
}

func (d *Decoder) decodeFixed(size, scale int) (any, error) {
	m := d.Fixed(size)
	if m == nil { // important: return nil and not m (as m is of type *big.Int)
		return nil, nil
	}
	return convertFixedToRat(m, scale), nil
}

// Fixed8Field decodes a fixed8 field.
func (d *Decoder) Fixed8Field(scale int) (any, error) {
	if !d.Bool() { // null value
		return nil, nil
	}
	return d.decodeFixed(Fixed8FieldSize, scale)
}

// Fixed12Field decodes a fixed12 field.
func (d *Decoder) Fixed12Field(scale int) (any, error) {
	if !d.Bool() { // null value
		return nil, nil
	}
	return d.decodeFixed(Fixed12FieldSize, scale)
}

// Fixed16Field decodes a fixed16 field.
func (d *Decoder) Fixed16Field(scale int) (any, error) {
	if !d.Bool() { // null value
		return nil, nil
	}
	return d.decodeFixed(Fixed16FieldSize, scale)
}

// VarField decodes a var field.
func (d *Decoder) VarField() (any, error) {
	_, b := d.LIBytes()
	/*
	   caution:
	   - result is used as driver.Value and we do need to provide a 'real' nil value
	   - returning b == nil does not work because b is of type []byte
	*/
	if b == nil {
		return nil, nil
	}
	return b, nil
}

// AlphanumField decodes a alphanum field.
func (d *Decoder) AlphanumField() (any, error) {
	if d.alphanumDfv1 { // like VarField
		return d.VarField()
	}
	_, b := d.LIBytes()
	/*
	   caution:
	   - result is used as driver.Value and we do need to provide a 'real' nil value
	   - returning b == nil does not work because b is of type []byte
	*/
	if b == nil {
		return nil, nil
	}
	/*
	   first byte:
	   - high bit set -> numeric
	   - high bit unset -> alpha
	   - bits 0-6: field size

	   ignore first byte for now
	*/
	return b[1:], nil
}

// Cesu8Field decodes a cesu8 field.
func (d *Decoder) Cesu8Field() (any, error) {
	_, b, err := d.CESU8LIBytes()
	if err != nil {
		return nil, err
	}
	/*
	   caution:
	   - result is used as driver.Value and we do need to provide a 'real' nil value
	   - returning b == nil does not work because b is of type []byte
	*/
	if b == nil {
		return nil, nil
	}
	return b, nil
}

// HexField decodes a hex field.
func (d *Decoder) HexField() (any, error) {
	_, b := d.LIBytes()
	/*
	   caution:
	   - result is used as driver.Value and we do need to provide a 'real' nil value
	   - returning b == nil does not work because b is of type []byte
	*/
	if b == nil {
		return nil, nil
	}
	return hex.EncodeToString(b), nil
}
