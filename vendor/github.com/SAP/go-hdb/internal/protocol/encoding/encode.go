// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package encoding

import (
	"encoding/binary"
	"io"
	"math"
	"math/big"

	"github.com/SAP/go-hdb/internal/unicode"
	"golang.org/x/text/transform"
)

const writeScratchSize = 4096

// Encoder encodes hdb protocol datatypes an basis of an io.Writer.
type Encoder struct {
	wr  io.Writer
	err error
	b   []byte // scratch buffer (min 15 Bytes - Decimal)
	tr  transform.Transformer
}

// NewEncoder creates a new Encoder instance.
func NewEncoder(wr io.Writer) *Encoder {
	return &Encoder{
		wr: wr,
		b:  make([]byte, writeScratchSize),
		tr: unicode.Utf8ToCesu8Transformer,
	}
}

// Zeroes writes cnt zero byte values.
func (e *Encoder) Zeroes(cnt int) {
	if e.err != nil {
		return
	}

	// zero out scratch area
	l := cnt
	if l > len(e.b) {
		l = len(e.b)
	}
	for i := 0; i < l; i++ {
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

// Bytes writes a bytes slice.
func (e *Encoder) Bytes(p []byte) {
	if e.err != nil {
		return
	}
	e.wr.Write(p)
}

// Byte writes a byte.
func (e *Encoder) Byte(b byte) { // WriteB as sig differs from WriteByte (vet issues)
	if e.err != nil {
		return
	}
	e.b[0] = b
	e.Bytes(e.b[:1])
}

// Bool writes a boolean.
func (e *Encoder) Bool(v bool) {
	if e.err != nil {
		return
	}
	if v {
		e.Byte(1)
	} else {
		e.Byte(0)
	}
}

// Int8 writes an int8.
func (e *Encoder) Int8(i int8) {
	if e.err != nil {
		return
	}
	e.Byte(byte(i))
}

// Int16 writes an int16.
func (e *Encoder) Int16(i int16) {
	if e.err != nil {
		return
	}
	binary.LittleEndian.PutUint16(e.b[:2], uint16(i))
	e.wr.Write(e.b[:2])
}

// Uint16 writes an uint16.
func (e *Encoder) Uint16(i uint16) {
	if e.err != nil {
		return
	}
	binary.LittleEndian.PutUint16(e.b[:2], i)
	e.wr.Write(e.b[:2])
}

// Int32 writes an int32.
func (e *Encoder) Int32(i int32) {
	if e.err != nil {
		return
	}
	binary.LittleEndian.PutUint32(e.b[:4], uint32(i))
	e.wr.Write(e.b[:4])
}

// Uint32 writes an uint32.
func (e *Encoder) Uint32(i uint32) {
	if e.err != nil {
		return
	}
	binary.LittleEndian.PutUint32(e.b[:4], i)
	e.wr.Write(e.b[:4])
}

// Int64 writes an int64.
func (e *Encoder) Int64(i int64) {
	if e.err != nil {
		return
	}
	binary.LittleEndian.PutUint64(e.b[:8], uint64(i))
	e.wr.Write(e.b[:8])
}

// Uint64 writes an uint64.
func (e *Encoder) Uint64(i uint64) {
	if e.err != nil {
		return
	}
	binary.LittleEndian.PutUint64(e.b[:8], i)
	e.wr.Write(e.b[:8])
}

// Float32 writes a float32.
func (e *Encoder) Float32(f float32) {
	if e.err != nil {
		return
	}
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(e.b[:4], bits)
	e.wr.Write(e.b[:4])
}

// Float64 writes a float64.
func (e *Encoder) Float64(f float64) {
	if e.err != nil {
		return
	}
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(e.b[:8], bits)
	e.wr.Write(e.b[:8])
}

// Decimal writes a decimal value.
func (e *Encoder) Decimal(m *big.Int, exp int) {
	b := e.b[:decSize]

	// little endian bigint words (significand) -> little endian db decimal format
	j := 0
	for _, d := range m.Bits() {
		for i := 0; i < _S; i++ {
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
	b[15] = byte(uint16(exp) >> 7)

	if m.Sign() == -1 {
		b[15] |= 0x80
	}

	e.wr.Write(b)
}

// Fixed writes a fixed decimal value.
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

	e.wr.Write(b)
}

// String writes a string.
func (e *Encoder) String(s string) {
	if e.err != nil {
		return
	}
	e.Bytes([]byte(s))
}

// CESU8Bytes writes an UTF-8 byte slice as CESU-8 and returns the CESU-8 bytes written.
func (e *Encoder) CESU8Bytes(p []byte) int {
	if e.err != nil {
		return 0
	}
	e.tr.Reset()
	cnt := 0
	i := 0
	for i < len(p) {
		m, n, err := e.tr.Transform(e.b, p[i:], true)
		if err != nil && err != transform.ErrShortDst {
			e.err = err
			return cnt
		}
		if m == 0 {
			e.err = transform.ErrShortDst
			return cnt
		}
		o, _ := e.wr.Write(e.b[:m])
		cnt += o
		i += n
	}
	return cnt
}

// CESU8String is like WriteCesu8 with an UTF-8 string as parameter.
func (e *Encoder) CESU8String(s string) int {
	return e.CESU8Bytes([]byte(s))
}
