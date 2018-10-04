/*
 * tuple.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// FoundationDB Go Tuple Layer

// Package tuple provides a layer for encoding and decoding multi-element tuples
// into keys usable by FoundationDB. The encoded key maintains the same sort
// order as the original tuple: sorted first by the first element, then by the
// second element, etc. This makes the tuple layer ideal for building a variety
// of higher-level data models.
//
// For general guidance on tuple usage, see the Tuple section of Data Modeling
// (https://apple.github.io/foundationdb/data-modeling.html#tuples).
//
// FoundationDB tuples can currently encode byte and unicode strings, integers,
// floats, doubles, booleans, UUIDs, tuples, and NULL values. In Go these are
// represented as []byte (or fdb.KeyConvertible), string, int64 (or int),
// float32, float64, bool, UUID, Tuple, and nil.
package tuple

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// A TupleElement is one of the types that may be encoded in FoundationDB
// tuples. Although the Go compiler cannot enforce this, it is a programming
// error to use an unsupported types as a TupleElement (and will typically
// result in a runtime panic).
//
// The valid types for TupleElement are []byte (or fdb.KeyConvertible), string,
// int64 (or int), float, double, bool, UUID, Tuple, and nil.
type TupleElement interface{}

// Tuple is a slice of objects that can be encoded as FoundationDB tuples. If
// any of the TupleElements are of unsupported types, a runtime panic will occur
// when the Tuple is packed.
//
// Given a Tuple T containing objects only of these types, then T will be
// identical to the Tuple returned by unpacking the byte slice obtained by
// packing T (modulo type normalization to []byte and int64).
type Tuple []TupleElement

// UUID wraps a basic byte array as a UUID. We do not provide any special
// methods for accessing or generating the UUID, but as Go does not provide
// a built-in UUID type, this simple wrapper allows for other libraries
// to write the output of their UUID type as a 16-byte array into
// an instance of this type.
type UUID [16]byte

// Type codes: These prefix the different elements in a packed Tuple
// to indicate what type they are.
const nilCode = 0x00
const bytesCode = 0x01
const stringCode = 0x02
const nestedCode = 0x05
const intZeroCode = 0x14
const posIntEnd = 0x1c
const negIntStart = 0x0c
const floatCode = 0x20
const doubleCode = 0x21
const falseCode = 0x26
const trueCode = 0x27
const uuidCode = 0x30

var sizeLimits = []uint64{
	1<<(0*8) - 1,
	1<<(1*8) - 1,
	1<<(2*8) - 1,
	1<<(3*8) - 1,
	1<<(4*8) - 1,
	1<<(5*8) - 1,
	1<<(6*8) - 1,
	1<<(7*8) - 1,
	1<<(8*8) - 1,
}

func bisectLeft(u uint64) int {
	var n int
	for sizeLimits[n] < u {
		n++
	}
	return n
}

func adjustFloatBytes(b []byte, encode bool) {
	if (encode && b[0]&0x80 != 0x00) || (!encode && b[0]&0x80 == 0x00) {
		// Negative numbers: flip all of the bytes.
		for i := 0; i < len(b); i++ {
			b[i] = b[i] ^ 0xff
		}
	} else {
		// Positive number: flip just the sign bit.
		b[0] = b[0] ^ 0x80
	}
}

type packer struct {
	buf []byte
}

func (p *packer) putByte(b byte) {
	p.buf = append(p.buf, b)
}

func (p *packer) putBytes(b []byte) {
	p.buf = append(p.buf, b...)
}

func (p *packer) putBytesNil(b []byte, i int) {
	for i >= 0 {
		p.putBytes(b[:i+1])
		p.putByte(0xFF)
		b = b[i+1:]
		i = bytes.IndexByte(b, 0x00)
	}
	p.putBytes(b)
}

func (p *packer) encodeBytes(code byte, b []byte) {
	p.putByte(code)
	if i := bytes.IndexByte(b, 0x00); i >= 0 {
		p.putBytesNil(b, i)
	} else {
		p.putBytes(b)
	}
	p.putByte(0x00)
}

func (p *packer) encodeInt(i int64) {
	if i == 0 {
		p.putByte(0x14)
		return
	}

	var n int
	var scratch [8]byte

	switch {
	case i > 0:
		n = bisectLeft(uint64(i))
		p.putByte(byte(intZeroCode + n))
		binary.BigEndian.PutUint64(scratch[:], uint64(i))
	case i < 0:
		n = bisectLeft(uint64(-i))
		p.putByte(byte(0x14 - n))
		offsetEncoded := int64(sizeLimits[n]) + i
		binary.BigEndian.PutUint64(scratch[:], uint64(offsetEncoded))
	}

	p.putBytes(scratch[8-n:])
}

func (p *packer) encodeFloat(f float32) {
	var scratch [4]byte
	binary.BigEndian.PutUint32(scratch[:], math.Float32bits(f))
	adjustFloatBytes(scratch[:], true)

	p.putByte(floatCode)
	p.putBytes(scratch[:])
}

func (p *packer) encodeDouble(d float64) {
	var scratch [8]byte
	binary.BigEndian.PutUint64(scratch[:], math.Float64bits(d))
	adjustFloatBytes(scratch[:], true)

	p.putByte(doubleCode)
	p.putBytes(scratch[:])
}

func (p *packer) encodeUUID(u UUID) {
	p.putByte(uuidCode)
	p.putBytes(u[:])
}

func (p *packer) encodeTuple(t Tuple, nested bool) {
	if nested {
		p.putByte(nestedCode)
	}

	for i, e := range t {
		switch e := e.(type) {
		case Tuple:
			p.encodeTuple(e, true)
		case nil:
			p.putByte(nilCode)
			if nested {
				p.putByte(0xff)
			}
		case int64:
			p.encodeInt(e)
		case int:
			p.encodeInt(int64(e))
		case []byte:
			p.encodeBytes(bytesCode, e)
		case fdb.KeyConvertible:
			p.encodeBytes(bytesCode, []byte(e.FDBKey()))
		case string:
			p.encodeBytes(stringCode, []byte(e))
		case float32:
			p.encodeFloat(e)
		case float64:
			p.encodeDouble(e)
		case bool:
			if e {
				p.putByte(trueCode)
			} else {
				p.putByte(falseCode)
			}
		case UUID:
			p.encodeUUID(e)
		default:
			panic(fmt.Sprintf("unencodable element at index %d (%v, type %T)", i, t[i], t[i]))
		}
	}

	if nested {
		p.putByte(0x00)
	}
}

// Pack returns a new byte slice encoding the provided tuple. Pack will panic if
// the tuple contains an element of any type other than []byte,
// fdb.KeyConvertible, string, int64, int, float32, float64, bool, tuple.UUID,
// nil, or a Tuple with elements of valid types.
//
// Tuple satisfies the fdb.KeyConvertible interface, so it is not necessary to
// call Pack when using a Tuple with a FoundationDB API function that requires a
// key.
func (t Tuple) Pack() []byte {
	p := packer{buf: make([]byte, 0, 64)}
	p.encodeTuple(t, false)
	return p.buf
}

func findTerminator(b []byte) int {
	bp := b
	var length int

	for {
		idx := bytes.IndexByte(bp, 0x00)
		length += idx
		if idx+1 == len(bp) || bp[idx+1] != 0xFF {
			break
		}
		length += 2
		bp = bp[idx+2:]
	}

	return length
}

func decodeBytes(b []byte) ([]byte, int) {
	idx := findTerminator(b[1:])
	return bytes.Replace(b[1:idx+1], []byte{0x00, 0xFF}, []byte{0x00}, -1), idx + 2
}

func decodeString(b []byte) (string, int) {
	bp, idx := decodeBytes(b)
	return string(bp), idx
}

func decodeInt(b []byte) (int64, int) {
	if b[0] == intZeroCode {
		return 0, 1
	}

	var neg bool

	n := int(b[0]) - intZeroCode
	if n < 0 {
		n = -n
		neg = true
	}

	bp := make([]byte, 8)
	copy(bp[8-n:], b[1:n+1])

	var ret int64

	binary.Read(bytes.NewBuffer(bp), binary.BigEndian, &ret)

	if neg {
		ret -= int64(sizeLimits[n])
	}

	return ret, n + 1
}

func decodeFloat(b []byte) (float32, int) {
	bp := make([]byte, 4)
	copy(bp, b[1:])
	adjustFloatBytes(bp, false)
	var ret float32
	binary.Read(bytes.NewBuffer(bp), binary.BigEndian, &ret)
	return ret, 5
}

func decodeDouble(b []byte) (float64, int) {
	bp := make([]byte, 8)
	copy(bp, b[1:])
	adjustFloatBytes(bp, false)
	var ret float64
	binary.Read(bytes.NewBuffer(bp), binary.BigEndian, &ret)
	return ret, 9
}

func decodeUUID(b []byte) (UUID, int) {
	var u UUID
	copy(u[:], b[1:])
	return u, 17
}

func decodeTuple(b []byte, nested bool) (Tuple, int, error) {
	var t Tuple

	var i int

	for i < len(b) {
		var el interface{}
		var off int

		switch {
		case b[i] == nilCode:
			if !nested {
				el = nil
				off = 1
			} else if i+1 < len(b) && b[i+1] == 0xff {
				el = nil
				off = 2
			} else {
				return t, i + 1, nil
			}
		case b[i] == bytesCode:
			el, off = decodeBytes(b[i:])
		case b[i] == stringCode:
			el, off = decodeString(b[i:])
		case negIntStart <= b[i] && b[i] <= posIntEnd:
			el, off = decodeInt(b[i:])
		case b[i] == floatCode:
			if i+5 > len(b) {
				return nil, i, fmt.Errorf("insufficient bytes to decode float starting at position %d of byte array for tuple", i)
			}
			el, off = decodeFloat(b[i:])
		case b[i] == doubleCode:
			if i+9 > len(b) {
				return nil, i, fmt.Errorf("insufficient bytes to decode double starting at position %d of byte array for tuple", i)
			}
			el, off = decodeDouble(b[i:])
		case b[i] == trueCode:
			el = true
			off = 1
		case b[i] == falseCode:
			el = false
			off = 1
		case b[i] == uuidCode:
			if i+17 > len(b) {
				return nil, i, fmt.Errorf("insufficient bytes to decode UUID starting at position %d of byte array for tuple", i)
			}
			el, off = decodeUUID(b[i:])
		case b[i] == nestedCode:
			var err error
			el, off, err = decodeTuple(b[i+1:], true)
			if err != nil {
				return nil, i, err
			}
			off++
		default:
			return nil, i, fmt.Errorf("unable to decode tuple element with unknown typecode %02x", b[i])
		}

		t = append(t, el)
		i += off
	}

	return t, i, nil
}

// Unpack returns the tuple encoded by the provided byte slice, or an error if
// the key does not correctly encode a FoundationDB tuple.
func Unpack(b []byte) (Tuple, error) {
	t, _, err := decodeTuple(b, false)
	return t, err
}

// FDBKey returns the packed representation of a Tuple, and allows Tuple to
// satisfy the fdb.KeyConvertible interface. FDBKey will panic in the same
// circumstances as Pack.
func (t Tuple) FDBKey() fdb.Key {
	return t.Pack()
}

// FDBRangeKeys allows Tuple to satisfy the fdb.ExactRange interface. The range
// represents all keys that encode tuples strictly starting with a Tuple (that
// is, all tuples of greater length than the Tuple of which the Tuple is a
// prefix).
func (t Tuple) FDBRangeKeys() (fdb.KeyConvertible, fdb.KeyConvertible) {
	p := t.Pack()
	return fdb.Key(concat(p, 0x00)), fdb.Key(concat(p, 0xFF))
}

// FDBRangeKeySelectors allows Tuple to satisfy the fdb.Range interface. The
// range represents all keys that encode tuples strictly starting with a Tuple
// (that is, all tuples of greater length than the Tuple of which the Tuple is a
// prefix).
func (t Tuple) FDBRangeKeySelectors() (fdb.Selectable, fdb.Selectable) {
	b, e := t.FDBRangeKeys()
	return fdb.FirstGreaterOrEqual(b), fdb.FirstGreaterOrEqual(e)
}

func concat(a []byte, b ...byte) []byte {
	r := make([]byte, len(a)+len(b))
	copy(r, a)
	copy(r[len(a):], b)
	return r
}
