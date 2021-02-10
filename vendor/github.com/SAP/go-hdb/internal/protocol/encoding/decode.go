// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/SAP/go-hdb/internal/unicode"
	"golang.org/x/text/transform"
)

const readScratchSize = 512 // used for skip as well - size not too small!

// Decoder decodes hdb protocol datatypes an basis of an io.Reader.
type Decoder struct {
	rd  io.Reader
	err error
	b   [readScratchSize]byte // scratch buffer
	tr  transform.Transformer
	cnt int
	dfv int
}

// NewDecoder creates a new Decoder instance based on an io.Reader.
func NewDecoder(rd io.Reader) *Decoder {
	return &Decoder{
		rd: rd,
		tr: unicode.Cesu8ToUtf8Transformer,
	}
}

// Dfv returns the data format version.
func (d *Decoder) Dfv() int {
	return d.dfv
}

// SetDfv sets the data format version.
func (d *Decoder) SetDfv(dfv int) {
	d.dfv = dfv
}

// ResetCnt resets the byte read counter.
func (d *Decoder) ResetCnt() {
	d.cnt = 0
}

// Cnt returns the value of the byte read counter.
func (d *Decoder) Cnt() int {
	return d.cnt
}

// Error returns the reader error.
func (d *Decoder) Error() error {
	return d.err
}

// ResetError return and resets reader error.
func (d *Decoder) ResetError() error {
	err := d.err
	d.err = nil
	return err
}

// Skip skips cnt bytes from reading.
func (d *Decoder) Skip(cnt int) {
	var n int
	for n < cnt {
		if d.err != nil {
			return
		}
		to := cnt - n
		if to > readScratchSize {
			to = readScratchSize
		}
		var m int
		m, d.err = io.ReadFull(d.rd, d.b[:to])
		n += m
		d.cnt += m
	}
}

// Byte reads and returns a byte.
func (d *Decoder) Byte() byte { // ReadB as sig differs from ReadByte (vet issues)
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:1])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return d.b[0]
}

// Bytes reads into a byte slice.
func (d *Decoder) Bytes(p []byte) {
	if d.err != nil {
		return
	}
	var n int
	n, d.err = io.ReadFull(d.rd, p)
	d.cnt += n
}

// Bool reads and returns a boolean.
func (d *Decoder) Bool() bool {
	if d.err != nil {
		return false
	}
	return d.Byte() != 0
}

// Int8 reads and returns an int8.
func (d *Decoder) Int8() int8 {
	return int8(d.Byte())
}

// Int16 reads and returns an int16.
func (d *Decoder) Int16() int16 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:2])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return int16(binary.LittleEndian.Uint16(d.b[:2]))
}

// Uint16 reads and returns an uint16.
func (d *Decoder) Uint16() uint16 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:2])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint16(d.b[:2])
}

// Int32 reads and returns an int32.
func (d *Decoder) Int32() int32 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:4])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(d.b[:4]))
}

// Uint32 reads and returns an uint32.
func (d *Decoder) Uint32() uint32 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:4])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(d.b[:4])
}

// Uint32ByteOrder reads and returns an uint32 in given byte order.
func (d *Decoder) Uint32ByteOrder(byteOrder binary.ByteOrder) uint32 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:4])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return byteOrder.Uint32(d.b[:4])
}

// Int64 reads and returns an int64.
func (d *Decoder) Int64() int64 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:8])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(d.b[:8]))
}

// Uint64 reads and returns an uint64.
func (d *Decoder) Uint64() uint64 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:8])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(d.b[:8])
}

// Float32 reads and returns a float32.
func (d *Decoder) Float32() float32 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:4])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	bits := binary.LittleEndian.Uint32(d.b[:4])
	return math.Float32frombits(bits)
}

// Float64 reads and returns a float64.
func (d *Decoder) Float64() float64 {
	if d.err != nil {
		return 0
	}
	var n int
	n, d.err = io.ReadFull(d.rd, d.b[:8])
	d.cnt += n
	if d.err != nil {
		return 0
	}
	bits := binary.LittleEndian.Uint64(d.b[:8])
	return math.Float64frombits(bits)
}

// Decimal reads and returns a decimal.
func (d *Decoder) Decimal() (*big.Int, int) { // m, exp
	bs := d.b[:decSize]

	var n int
	n, d.err = io.ReadFull(d.rd, bs)
	d.cnt += n
	if d.err != nil {
		return nil, 0
	}

	if (bs[15] & 0x70) == 0x70 { //null value (bit 4,5,6 set)
		return nil, 0
	}

	if (bs[15] & 0x60) == 0x60 {
		d.err = fmt.Errorf("decimal: format (infinity, nan, ...) not supported : %v", bs)
		return nil, 0
	}

	neg := (bs[15] & 0x80) != 0
	exp := int((((uint16(bs[15])<<8)|uint16(bs[14]))<<1)>>2) - dec128Bias

	// b14 := b[14]  // save b[14]
	bs[14] &= 0x01 // keep the mantissa bit (rest: sign and exp)

	//most significand byte
	msb := 14
	for msb > 0 && bs[msb] == 0 {
		msb--
	}

	//calc number of words
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
	return m, exp
}

// Fixed reads and returns a fixed decimal.
func (d *Decoder) Fixed(size int) *big.Int { // m, exp
	bs := d.b[:size]

	var n int
	n, d.err = io.ReadFull(d.rd, bs)
	d.cnt += n
	if d.err != nil {
		return nil
	}

	neg := (bs[size-1] & 0x80) != 0 // is negative number (2s complement)

	//most significand byte
	msb := size - 1
	for msb > 0 && bs[msb] == 0 {
		msb--
	}

	//calc number of words
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

// CESU8Bytes reads a size CESU-8 encoded byte sequence and returns an UTF-8 byte slice.
func (d *Decoder) CESU8Bytes(size int) []byte {
	if d.err != nil {
		return nil
	}
	p := make([]byte, size)
	var n int
	n, d.err = io.ReadFull(d.rd, p)
	d.cnt += n
	if d.err != nil {
		return nil
	}
	d.tr.Reset()
	if n, _, d.err = d.tr.Transform(p, p, true); d.err != nil { // inplace transformation
		return nil
	}
	return p[:n]
}

// // ShortCESU8Bytes reads a CESU-8 encoded byte sequence and returns an UTF-8 byte slice.
// // Size is encoded in one byte.
// func (d *Decoder) ShortCESU8Bytes() ([]byte, int) {
// 	size := d.Byte()
// 	return d.CESU8Bytes(int(size)), int(size)
// }

// // ShortCESU8String reads a CESU-8 encoded byte sequence and returns an UTF-8 string.
// // Size is encoded in one byte.
// func (d *Decoder) ShortCESU8Bytes() (string, int) {
// 	b, n := d.ShortCESU8Bytes()
// 	return b, n
// }

// // ShortBytes reads a byte sequence and returns a byte slice.
// // Size is encoded in one byte.
// func (d *Decoder) ShortBytes() ([]byte, int) {
// 	size := d.Byte()
// 	b := make([]byte, size)
// 	d.Bytes(b)
// 	return b
// }
