package encoding

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"time"

	"github.com/SAP/go-hdb/driver/internal/unsafe"
	"github.com/SAP/go-hdb/driver/unicode/cesu8"
	"golang.org/x/text/transform"
)

const writeScratchSize = 4096

// Encoder encodes hdb protocol datatypes on basis of an io.Writer.
type Encoder struct {
	wr io.Writer
	b  []byte // scratch buffer (min 15 Bytes - Decimal)
	tr transform.Transformer
}

// NewEncoder creates a new Encoder instance.
func NewEncoder(wr io.Writer, tr transform.Transformer) *Encoder {
	return &Encoder{
		wr: wr,
		b:  make([]byte, writeScratchSize),
		tr: tr,
	}
}

// Zeroes encodes cnt zero byte values.
func (e *Encoder) Zeroes(cnt int) {
	// zero out scratch area
	l := cnt
	if l > len(e.b) {
		l = len(e.b)
	}
	for i := range l {
		e.b[i] = 0
	}

	for i := 0; i < cnt; {
		j := cnt - i
		if j > len(e.b) {
			j = len(e.b)
		}
		n, _ := e.wr.Write(e.b[:j])
		if n != j {
			return
		}
		i += n
	}
}

// Bytes encodes bytes.
func (e *Encoder) Bytes(p []byte) {
	e.wr.Write(p) //nolint:errcheck
}

// Byte encodes a byte.
func (e *Encoder) Byte(b byte) { // WriteB as sig differs from WriteByte (vet issues)
	e.b[0] = b
	e.Bytes(e.b[:1])
}

// Bool encodes a boolean.
func (e *Encoder) Bool(v bool) {
	if v {
		e.Byte(1)
	} else {
		e.Byte(0)
	}
}

// Int8 encodes an int8.
func (e *Encoder) Int8(i int8) {
	e.Byte(byte(i))
}

// Int16 encodes an int16.
func (e *Encoder) Int16(i int16) {
	binary.LittleEndian.PutUint16(e.b[:2], uint16(i)) //nolint: gosec
	e.wr.Write(e.b[:2])                               //nolint:errcheck
}

// Uint16 encodes an uint16.
func (e *Encoder) Uint16(i uint16) {
	binary.LittleEndian.PutUint16(e.b[:2], i)
	e.wr.Write(e.b[:2]) //nolint:errcheck
}

// Uint16ByteOrder encodes an uint16 in given byte order.
func (e *Encoder) Uint16ByteOrder(i uint16, byteOrder binary.ByteOrder) {
	byteOrder.PutUint16(e.b[:2], i)
	e.wr.Write(e.b[:2]) //nolint:errcheck
}

// Int32 encodes an int32.
func (e *Encoder) Int32(i int32) {
	binary.LittleEndian.PutUint32(e.b[:4], uint32(i)) //nolint:gosec
	e.wr.Write(e.b[:4])                               //nolint:errcheck
}

// Uint32 encodes an uint32.
func (e *Encoder) Uint32(i uint32) {
	binary.LittleEndian.PutUint32(e.b[:4], i)
	e.wr.Write(e.b[:4]) //nolint:errcheck
}

// Int64 encodes an int64.
func (e *Encoder) Int64(i int64) {
	binary.LittleEndian.PutUint64(e.b[:8], uint64(i)) //nolint:gosec
	e.wr.Write(e.b[:8])                               //nolint:errcheck
}

// Uint64 encodes an uint64.
func (e *Encoder) Uint64(i uint64) {
	binary.LittleEndian.PutUint64(e.b[:8], i)
	e.wr.Write(e.b[:8]) //nolint:errcheck
}

// Float32 encodes a float32.
func (e *Encoder) Float32(f float32) {
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(e.b[:4], bits)
	e.wr.Write(e.b[:4]) //nolint:errcheck
}

// Float64 encodes a float64.
func (e *Encoder) Float64(f float64) {
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(e.b[:8], bits)
	e.wr.Write(e.b[:8]) //nolint:errcheck
}

// Decimal encodes a decimal value.
func (e *Encoder) Decimal(m *big.Int, exp int) {
	b := e.b[:decSize]

	// little endian bigint words (significand) -> little endian db decimal format
	j := 0
	for _, d := range m.Bits() {
		for range _S {
			b[j] = byte(d)
			d >>= 8
			j++
		}
	}

	// clear scratch buffer
	for i := j; i < decSize; i++ {
		b[i] = 0
	}

	exp += dec128Bias
	b[14] |= (byte(exp) << 1)
	b[15] = byte(uint16(exp) >> 7) //nolint: gosec

	if m.Sign() == -1 {
		b[15] |= 0x80
	}

	e.wr.Write(b) //nolint:errcheck
}

// Fixed encodes a fixed decimal value.
func (e *Encoder) Fixed(m *big.Int, size int) {
	b := e.b[:size]

	neg := m.Sign() == -1
	fill := byte(0)

	if neg {
		// make positive
		m.Neg(m)
		// 2s complement
		bits := m.Bits()
		// - invert all bits
		for i := 0; i < len(bits); i++ {
			bits[i] = ^bits[i]
		}
		// - add 1
		m.Add(m, natOne)
		fill = 0xff
	}

	// little endian bigint words (significand) -> little endian db decimal format
	j := 0
	for _, d := range m.Bits() {
		/*
			check j < size as number of bytes in m.Bits words can exceed number of fixed size bytes
			e.g. 64 bit architecture:
			- two words equals 16 bytes but fixed size might be 12 bytes
			- invariant: all 'skipped' bytes in most significant word are zero
		*/
		for i := 0; i < _S && j < size; i++ {
			b[j] = byte(d)
			d >>= 8
			j++
		}
	}

	// clear scratch buffer
	for i := j; i < size; i++ {
		b[i] = fill
	}

	e.wr.Write(b) //nolint:errcheck
}

// String encodes a string.
func (e *Encoder) String(s string) { e.Bytes(unsafe.String2ByteSlice(s)) }

// CESU8Bytes encodes UTF-8 bytes into CESU-8 and returns the CESU-8 bytes written.
func (e *Encoder) CESU8Bytes(p []byte) (int, error) {
	e.tr.Reset()
	cnt := 0
	for i := 0; i < len(p); {
		nDst, nSrc, err := e.tr.Transform(e.b, p[i:], true)
		if nDst != 0 {
			n, _ := e.wr.Write(e.b[:nDst])
			cnt += n
		}
		if err != nil && !errors.Is(err, transform.ErrShortDst) {
			return cnt, err
		}
		i += nSrc
	}
	return cnt, nil
}

// CESU8String encodes an UTF-8 string into CESU-8 and returns the CESU-8 bytes written.
func (e *Encoder) CESU8String(s string) (int, error) { return e.CESU8Bytes(unsafe.String2ByteSlice(s)) }

// varFieldInd encodes a variable field indicator.
func (e *Encoder) varFieldInd(size int) error {
	switch {
	default:
		return fmt.Errorf("max argument length %d of string exceeded", size)
	case size <= int(bytesLenIndSmall):
		e.Byte(byte(size))
	case size <= math.MaxInt16:
		e.Byte(bytesLenIndMedium)
		e.Int16(int16(size))
	case size <= math.MaxInt32:
		e.Byte(bytesLenIndBig)
		e.Int32(int32(size))
	}
	return nil
}

// LIBytes encodes bytes with length indicator.
func (e *Encoder) LIBytes(p []byte) error {
	if err := e.varFieldInd(len(p)); err != nil {
		return err
	}
	e.Bytes(p)
	return nil
}

// LIString encodes a string with length indicator.
func (e *Encoder) LIString(s string) error {
	if err := e.varFieldInd(len(s)); err != nil {
		return err
	}
	e.String(s)
	return nil
}

// CESU8LIBytes encodes UTF-8 into CESU-8 bytes with length indicator.
func (e *Encoder) CESU8LIBytes(p []byte) error {
	size := cesu8.Size(p)
	if err := e.varFieldInd(size); err != nil {
		return err
	}
	_, err := e.CESU8Bytes(p)
	return err
}

// CESU8LIString encodes an UTF-8 into a CESU-8 string with length indicator.
func (e *Encoder) CESU8LIString(s string) error {
	size := cesu8.StringSize(s)
	if err := e.varFieldInd(size); err != nil {
		return err
	}
	_, err := e.CESU8String(s)
	return err
}

// Fields.
func asInt[E byte | int16 | int32 | int64](v any) E {
	i64, ok := v.(int64)
	if !ok {
		panic("invalid integer") // should never happen
	}
	return E(i64)
}

func asTime(v any) time.Time {
	t, ok := v.(time.Time)
	if !ok {
		panic("invalid time") // should never happen
	}
	// store in utc
	return t.UTC()
}

// BooleanField encodes a boolean field.
func (e *Encoder) BooleanField(v any) error {
	if v == nil {
		e.Byte(booleanNullValue)
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		panic("invalid boolean") // should never happen
	}
	if b {
		e.Byte(booleanTrueValue)
	} else {
		e.Byte(booleanFalseValue)
	}
	return nil
}

// TinyintField encodes a tinyint field.
func (e *Encoder) TinyintField(v any) error {
	e.Byte(asInt[byte](v))
	return nil
}

// SmallintField encodes a smallint field.
func (e *Encoder) SmallintField(v any) error {
	e.Int16(asInt[int16](v))
	return nil
}

// IntegerField encodes a integer field.
func (e *Encoder) IntegerField(v any) error {
	e.Int32(asInt[int32](v))
	return nil
}

// BigintField encodes a bigint field.
func (e *Encoder) BigintField(v any) error {
	e.Int64(asInt[int64](v))
	return nil
}

// RealField encodes a real field.
func (e *Encoder) RealField(v any) error {
	f64, ok := v.(float64)
	if !ok {
		panic("invalid real") // should never happen
	}
	e.Float32(float32(f64))
	return nil
}

// DoubleField encodes a double field.
func (e *Encoder) DoubleField(v any) error {
	f64, ok := v.(float64)
	if !ok {
		panic("invalid double") // should never happen
	}
	e.Float64(f64)
	return nil
}

func (e *Encoder) encodeDate(t time.Time) {
	// year: set most sig bit
	// month 0 based
	year, month, day := t.Date()
	e.Uint16(uint16(year) | 0x8000) //nolint: gosec
	e.Int8(int8(month) - 1)         //nolint: gosec
	e.Int8(int8(day))               //nolint: gosec
}

// DateField encodes a dayte field.
func (e *Encoder) DateField(v any) error {
	e.encodeDate(asTime(v))
	return nil
}

func (e *Encoder) encodeTime(t time.Time) {
	e.Byte(byte(t.Hour()) | 0x80)
	e.Int8(int8(t.Minute())) //nolint: gosec
	msec := t.Second()*1000 + t.Nanosecond()/1000000
	e.Uint16(uint16(msec)) //nolint: gosec
}

// TimeField encodes a time field.
func (e *Encoder) TimeField(v any) error {
	e.encodeTime(asTime(v))
	return nil
}

// TimestampField encodes a timestamp field.
func (e *Encoder) TimestampField(v any) error {
	t := asTime(v)
	e.encodeDate(t)
	e.encodeTime(t)
	return nil
}

// LongdateField encodea a longdate field.
func (e *Encoder) LongdateField(v any) error {
	e.Int64(convertTimeToLongdate(asTime(v)))
	return nil
}

// SeconddateField encodes a seconddate field.
func (e *Encoder) SeconddateField(v any) error {
	e.Int64(convertTimeToSeconddate(asTime(v)))
	return nil
}

// DaydateField encodes a daydate field.
func (e *Encoder) DaydateField(v any) error {
	e.Int32(int32(convertTimeToDayDate(asTime(v)))) //nolint: gosec
	return nil
}

// SecondtimeField encodes a secondtime field.
func (e *Encoder) SecondtimeField(v any) error {
	if v == nil {
		e.Int32(secondtimeNullValue)
		return nil
	}
	e.Int32(int32(convertTimeToSecondtime(asTime(v)))) //nolint: gosec
	return nil
}

func (e *Encoder) encodeFixed(v any, size, prec, scale int) error {
	r, ok := v.(*big.Rat)
	if !ok {
		panic("invalid fixed") // should never happen
	}

	var m big.Int
	df := convertRatToFixed(r, &m, prec, scale)

	if df&dfOverflow != 0 {
		return ErrDecimalOutOfRange
	}

	e.Fixed(&m, size)
	return nil
}

// DecimalField encodes a decimal field.
func (e *Encoder) DecimalField(v any) error {
	r, ok := v.(*big.Rat)
	if !ok {
		panic("invalid decimal") // should never happen
	}

	var m big.Int
	exp, df := convertRatToDecimal(r, &m, dec128Digits, dec128MinExp, dec128MaxExp)

	if df&dfOverflow != 0 {
		return ErrDecimalOutOfRange
	}

	if df&dfUnderflow != 0 { // set to zero
		e.Decimal(natZero, 0)
	} else {
		e.Decimal(&m, exp)
	}
	return nil
}

// Fixed8Field encodes a fixed8 field.
func (e *Encoder) Fixed8Field(v any, prec, scale int) error {
	return e.encodeFixed(v, Fixed8FieldSize, prec, scale)
}

// Fixed12Field encodes a fixed12 field.
func (e *Encoder) Fixed12Field(v any, prec, scale int) error {
	return e.encodeFixed(v, Fixed12FieldSize, prec, scale)
}

// Fixed16Field encodes a fixed16 field.
func (e *Encoder) Fixed16Field(v any, prec, scale int) error {
	return e.encodeFixed(v, Fixed16FieldSize, prec, scale)
}

// VarField encodes a var field.
func (e *Encoder) VarField(v any) error {
	switch v := v.(type) {
	case []byte:
		return e.LIBytes(v)
	case string:
		return e.LIString(v)
	default:
		panic("invalid var value") // should never happen
	}
}

// Cesu8Field encodes a cesu8 field.
func (e *Encoder) Cesu8Field(v any) error {
	switch v := v.(type) {
	case []byte:
		return e.CESU8LIBytes(v)
	case string:
		return e.CESU8LIString(v)
	default:
		panic("invalid cesu8 value") // should never happen
	}
}

// HexField encodes a hex field.
func (e *Encoder) HexField(v any) error {
	switch v := v.(type) {
	case []byte:
		b, err := hex.DecodeString(string(v))
		if err != nil {
			return err
		}
		return e.LIBytes(b)
	case string:
		b, err := hex.DecodeString(v)
		if err != nil {
			return err
		}
		return e.LIBytes(b)
	default:
		panic("invalid hex value") // should never happen
	}
}
