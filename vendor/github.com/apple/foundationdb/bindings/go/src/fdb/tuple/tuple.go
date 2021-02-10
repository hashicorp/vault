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
// large integers, floats, doubles, booleans, UUIDs, tuples, and NULL values.
// In Go these are represented as []byte (or fdb.KeyConvertible), string, int64
// (or int, uint, uint64), *big.Int (or big.Int), float32, float64, bool,
// UUID, Tuple, and nil.
package tuple

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// A TupleElement is one of the types that may be encoded in FoundationDB
// tuples. Although the Go compiler cannot enforce this, it is a programming
// error to use an unsupported types as a TupleElement (and will typically
// result in a runtime panic).
//
// The valid types for TupleElement are []byte (or fdb.KeyConvertible), string,
// int64 (or int, uint, uint64), *big.Int (or big.Int), float, double, bool,
// UUID, Tuple, and nil.
type TupleElement interface{}

// Tuple is a slice of objects that can be encoded as FoundationDB tuples. If
// any of the TupleElements are of unsupported types, a runtime panic will occur
// when the Tuple is packed.
//
// Given a Tuple T containing objects only of these types, then T will be
// identical to the Tuple returned by unpacking the byte slice obtained by
// packing T (modulo type normalization to []byte, uint64, and int64).
type Tuple []TupleElement

// String implements the fmt.Stringer interface and returns human-readable
// string representation of this tuple. For most elements, we use the
// object's default string representation.
func (tuple Tuple) String() string {
	sb := strings.Builder{}
	printTuple(tuple, &sb)
	return sb.String()
}

func printTuple(tuple Tuple, sb *strings.Builder) {
	sb.WriteString("(")

	for i, t := range tuple {
		switch t := t.(type) {
		case Tuple:
			printTuple(t, sb)
		case nil:
			sb.WriteString("<nil>")
		case string:
			sb.WriteString(strconv.Quote(t))
		case UUID:
			sb.WriteString("UUID(")
			sb.WriteString(t.String())
			sb.WriteString(")")
		case []byte:
			sb.WriteString("b\"")
			sb.WriteString(fdb.Printable(t))
			sb.WriteString("\"")
		default:
			// For user-defined and standard types, we use standard Go
			// printer, which itself uses Stringer interface.
			fmt.Fprintf(sb, "%v", t)
		}

		if (i < len(tuple) - 1) {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(")")
}

// UUID wraps a basic byte array as a UUID. We do not provide any special
// methods for accessing or generating the UUID, but as Go does not provide
// a built-in UUID type, this simple wrapper allows for other libraries
// to write the output of their UUID type as a 16-byte array into
// an instance of this type.
type UUID [16]byte

func (uuid UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// Versionstamp is struct for a FoundationDB verionstamp. Versionstamps are
// 12 bytes long composed of a 10 byte transaction version and a 2 byte user
// version. The transaction version is filled in at commit time and the user
// version is provided by the application to order results within a transaction.
type Versionstamp struct {
	TransactionVersion [10]byte
	UserVersion        uint16
}

// Returns a human-readable string for this Versionstamp.
func (vs Versionstamp) String() string {
	return fmt.Sprintf("Versionstamp(%s, %d)", fdb.Printable(vs.TransactionVersion[:]), vs.UserVersion)
}

var incompleteTransactionVersion = [10]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

const versionstampLength = 12

// IncompleteVersionstamp is the constructor you should use to make
// an incomplete versionstamp to use in a tuple.
func IncompleteVersionstamp(userVersion uint16) Versionstamp {
	return Versionstamp{
		TransactionVersion: incompleteTransactionVersion,
		UserVersion:        userVersion,
	}
}

// Bytes converts a Versionstamp struct to a byte slice for encoding in a tuple.
func (v Versionstamp) Bytes() []byte {
	var scratch [versionstampLength]byte

	copy(scratch[:], v.TransactionVersion[:])

	binary.BigEndian.PutUint16(scratch[10:], v.UserVersion)

	return scratch[:]
}

// Type codes: These prefix the different elements in a packed Tuple
// to indicate what type they are.
const nilCode = 0x00
const bytesCode = 0x01
const stringCode = 0x02
const nestedCode = 0x05
const intZeroCode = 0x14
const posIntEnd = 0x1d
const negIntStart = 0x0b
const floatCode = 0x20
const doubleCode = 0x21
const falseCode = 0x26
const trueCode = 0x27
const uuidCode = 0x30
const versionstampCode = 0x33

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

var minInt64BigInt = big.NewInt(math.MinInt64)

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
	versionstampPos int32
	buf             []byte
}

func newPacker() *packer {
	return &packer{
		versionstampPos: -1,
		buf:             make([]byte, 0, 64),
	}
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

func (p *packer) encodeUint(i uint64) {
	if i == 0 {
		p.putByte(intZeroCode)
		return
	}

	n := bisectLeft(i)
	var scratch [8]byte

	p.putByte(byte(intZeroCode + n))
	binary.BigEndian.PutUint64(scratch[:], i)

	p.putBytes(scratch[8-n:])
}

func (p *packer) encodeInt(i int64) {
	if i >= 0 {
		p.encodeUint(uint64(i))
		return
	}

	n := bisectLeft(uint64(-i))
	var scratch [8]byte

	p.putByte(byte(intZeroCode - n))
	offsetEncoded := int64(sizeLimits[n]) + i
	binary.BigEndian.PutUint64(scratch[:], uint64(offsetEncoded))

	p.putBytes(scratch[8-n:])
}

func (p *packer) encodeBigInt(i *big.Int) {
	length := len(i.Bytes())
	if length > 0xff {
		panic(fmt.Sprintf("Integer magnitude is too large (more than 255 bytes)"))
	}

	if i.Sign() >= 0 {
		intBytes := i.Bytes()
		if length > 8 {
			p.putByte(byte(posIntEnd))
			p.putByte(byte(len(intBytes)))
		} else {
			p.putByte(byte(intZeroCode + length))
		}

		p.putBytes(intBytes)
	} else {
		add := new(big.Int).Lsh(big.NewInt(1), uint(length*8))
		add.Sub(add, big.NewInt(1))
		transformed := new(big.Int)
		transformed.Add(i, add)

		intBytes := transformed.Bytes()
		if length > 8 {
			p.putByte(byte(negIntStart))
			p.putByte(byte(length ^ 0xff))
		} else {
			p.putByte(byte(intZeroCode - length))
		}

		// For large negative numbers whose absolute value begins with 0xff bytes,
		// the transformed bytes may begin with 0x00 bytes. However, intBytes
		// will only contain the non-zero suffix, so this loop is needed to make
		// the value written be the correct length.
		for i := len(intBytes); i < length; i++ {
			p.putByte(0x00)
		}

		p.putBytes(intBytes)
	}
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

func (p *packer) encodeVersionstamp(v Versionstamp) {
	p.putByte(versionstampCode)

	isIncomplete := v.TransactionVersion == incompleteTransactionVersion
	if isIncomplete {
		if p.versionstampPos != -1 {
			panic(fmt.Sprintf("Tuple can only contain one incomplete versionstamp"))
		}

		p.versionstampPos = int32(len(p.buf))
	}

	p.putBytes(v.Bytes())
}

func (p *packer) encodeTuple(t Tuple, nested bool, versionstamps bool) {
	if nested {
		p.putByte(nestedCode)
	}

	for i, e := range t {
		switch e := e.(type) {
		case Tuple:
			p.encodeTuple(e, true, versionstamps)
		case nil:
			p.putByte(nilCode)
			if nested {
				p.putByte(0xff)
			}
		case int:
			p.encodeInt(int64(e))
		case int64:
			p.encodeInt(e)
		case uint:
			p.encodeUint(uint64(e))
		case uint64:
			p.encodeUint(e)
		case *big.Int:
			p.encodeBigInt(e)
		case big.Int:
			p.encodeBigInt(&e)
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
		case Versionstamp:
			if versionstamps == false && e.TransactionVersion == incompleteTransactionVersion {
				panic(fmt.Sprintf("Incomplete Versionstamp included in vanilla tuple pack"))
			}

			p.encodeVersionstamp(e)
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
// fdb.KeyConvertible, string, int64, int, uint64, uint, *big.Int, big.Int, float32,
// float64, bool, tuple.UUID, tuple.Versionstamp, nil, or a Tuple with elements of
// valid types. It will also panic if an integer is specified with a value outside
// the range [-2**2040+1, 2**2040-1]
//
// Tuple satisfies the fdb.KeyConvertible interface, so it is not necessary to
// call Pack when using a Tuple with a FoundationDB API function that requires a
// key.
//
// This method will panic if it contains an incomplete Versionstamp. Use
// PackWithVersionstamp instead.
//
func (t Tuple) Pack() []byte {
	p := newPacker()
	p.encodeTuple(t, false, false)
	return p.buf
}

// PackWithVersionstamp packs the specified tuple into a key for versionstamp
// operations. See Pack for more information. This function will return an error
// if you attempt to pack a tuple with more than one versionstamp. This function will
// return an error if you attempt to pack a tuple with a versionstamp position larger
// than an uint16 if the API version is less than 520.
func (t Tuple) PackWithVersionstamp(prefix []byte) ([]byte, error) {
	hasVersionstamp, err := t.HasIncompleteVersionstamp()
	if err != nil {
		return nil, err
	}

	apiVersion, err := fdb.GetAPIVersion()
	if err != nil {
		return nil, err
	}

	if hasVersionstamp == false {
		return nil, errors.New("No incomplete versionstamp included in tuple pack with versionstamp")
	}

	p := newPacker()

	if prefix != nil {
		p.putBytes(prefix)
	}

	p.encodeTuple(t, false, true)

	if hasVersionstamp {
		var scratch [4]byte
		var offsetIndex int
		if apiVersion < 520 {
			if p.versionstampPos > math.MaxUint16 {
				return nil, errors.New("Versionstamp position too large")
			}

			offsetIndex = 2
			binary.LittleEndian.PutUint16(scratch[:], uint16(p.versionstampPos))
		} else {
			offsetIndex = 4
			binary.LittleEndian.PutUint32(scratch[:], uint32(p.versionstampPos))
		}

		p.putBytes(scratch[0:offsetIndex])
	}

	return p.buf, nil
}

// HasIncompleteVersionstamp determines if there is at least one incomplete
// versionstamp in a tuple. This function will return an error this tuple has
// more than one versionstamp.
func (t Tuple) HasIncompleteVersionstamp() (bool, error) {
	incompleteCount := t.countIncompleteVersionstamps()

	var err error
	if incompleteCount > 1 {
		err = errors.New("Tuple can only contain one incomplete versionstamp")
	}

	return incompleteCount >= 1, err
}

func (t Tuple) countIncompleteVersionstamps() int {
	incompleteCount := 0

	for _, el := range t {
		switch e := el.(type) {
		case Versionstamp:
			if e.TransactionVersion == incompleteTransactionVersion {
				incompleteCount++
			}
		case Tuple:
			incompleteCount += e.countIncompleteVersionstamps()
		}
	}

	return incompleteCount
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

func decodeInt(b []byte) (interface{}, int) {
	if b[0] == intZeroCode {
		return int64(0), 1
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
		return ret - int64(sizeLimits[n]), n + 1
	}

	if ret > 0 {
		return ret, n + 1
	}

	// The encoded value claimed to be positive yet when put in an int64
	// produced a negative value. This means that the number must be a positive
	// 64-bit value that uses the most significant bit. This can be fit in a
	// uint64, so return that. Note that this is the *only* time we return
	// a uint64.
	return uint64(ret), n + 1
}

func decodeBigInt(b []byte) (interface{}, int) {
	val := new(big.Int)
	offset := 1
	var length int

	if b[0] == negIntStart || b[0] == posIntEnd {
		length = int(b[1])
		if b[0] == negIntStart {
			length ^= 0xff
		}

		offset += 1
	} else {
		// Must be a negative 8 byte integer
		length = 8
	}

	val.SetBytes(b[offset : length+offset])

	if b[0] < intZeroCode {
		sub := new(big.Int).Lsh(big.NewInt(1), uint(length)*8)
		sub.Sub(sub, big.NewInt(1))
		val.Sub(val, sub)
	}

	// This is the only value that fits in an int64 or uint64 that is decoded with this function
	if val.Cmp(minInt64BigInt) == 0 {
		return val.Int64(), length + offset
	}

	return val, length + offset
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

func decodeVersionstamp(b []byte) (Versionstamp, int) {
	var transactionVersion [10]byte
	var userVersion uint16

	copy(transactionVersion[:], b[1:11])

	userVersion = binary.BigEndian.Uint16(b[11:])

	return Versionstamp{
		TransactionVersion: transactionVersion,
		UserVersion:        userVersion,
	}, versionstampLength + 1
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
		case negIntStart+1 < b[i] && b[i] < posIntEnd:
			el, off = decodeInt(b[i:])
		case negIntStart+1 == b[i] && (b[i+1]&0x80 != 0):
			el, off = decodeInt(b[i:])
		case negIntStart <= b[i] && b[i] <= posIntEnd:
			el, off = decodeBigInt(b[i:])
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
		case b[i] == versionstampCode:
			if i+versionstampLength+1 > len(b) {
				return nil, i, fmt.Errorf("insufficient bytes to decode Versionstamp starting at position %d of byte array for tuple", i)
			}
			el, off = decodeVersionstamp(b[i:])
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
