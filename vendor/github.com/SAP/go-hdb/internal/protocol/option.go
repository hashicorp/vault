// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

type optBooleanType bool
type optTinyintType int8
type optIntType int32
type optBigintType int64
type optDoubleType float64
type optStringType []byte
type optBinaryStringType []byte

func (t optBooleanType) String() string      { return fmt.Sprintf("%t", bool(t)) }
func (t optTinyintType) String() string      { return fmt.Sprintf("%d", int8(t)) }
func (t optIntType) String() string          { return fmt.Sprintf("%d", int(t)) }
func (t optBigintType) String() string       { return fmt.Sprintf("%d", int64(t)) }
func (t optDoubleType) String() string       { return fmt.Sprintf("%g", float64(t)) }
func (t optStringType) String() string       { return string(t) }
func (t optBinaryStringType) String() string { return fmt.Sprintf("%v", []byte(t)) }

type multiLineOptions []plainOptions

func (o multiLineOptions) size() int {
	size := 0
	for _, m := range o {
		size += m.size()
	}
	return size
}

func (o *multiLineOptions) reset(size int) {
	if o == nil || size > cap(*o) {
		*o = make(multiLineOptions, size)
	} else {
		*o = (*o)[:size]
	}
}

func (o *multiLineOptions) decode(dec *encoding.Decoder, lineCnt int) {
	o.reset(lineCnt)
	for i := 0; i < lineCnt; i++ {
		m := plainOptions{}
		(*o)[i] = m
		cnt := dec.Int16()
		m.decode(dec, int(cnt))
	}
}

func (o multiLineOptions) encode(enc *encoding.Encoder) {
	for _, m := range o {
		enc.Int16(int16(len(m)))
		m.encode(enc)
	}
}

type plainOptions map[int8]interface{}

func (o plainOptions) size() int {
	size := 2 * len(o) //option + type
	for _, v := range o {
		switch v := v.(type) {
		default:
			plog.Fatalf("type %T not implemented", v)
		case optBooleanType:
			size++
		case optTinyintType:
			size++
		case optIntType:
			size += 4
		case optBigintType:
			size += 8
		case optDoubleType:
			size += 8
		case optStringType:
			size += (2 + len(v)) //length int16 + string length
		case optBinaryStringType:
			size += (2 + len(v)) //length int16 + string length
		}
	}
	return size
}

func (o plainOptions) decode(dec *encoding.Decoder, cnt int) {

	for i := 0; i < cnt; i++ {

		k := dec.Int8()
		tc := dec.Byte()

		switch typeCode(tc) {

		default:
			plog.Fatalf("type code %s not implemented", typeCode(tc))

		case tcBoolean:
			o[k] = optBooleanType(dec.Bool())

		case tcTinyint:
			o[k] = optTinyintType(dec.Int8())

		case tcInteger:
			o[k] = optIntType(dec.Int32())

		case tcBigint:
			o[k] = optBigintType(dec.Int64())

		case tcDouble:
			o[k] = optDoubleType(dec.Float64())

		case tcString:
			size := dec.Int16()
			v := make([]byte, size)
			dec.Bytes(v)
			o[k] = optStringType(v)

		case tcBstring:
			size := dec.Int16()
			v := make([]byte, size)
			dec.Bytes(v)
			o[k] = optBinaryStringType(v)
		}
	}
}

func (o plainOptions) encode(enc *encoding.Encoder) {

	for k, v := range o {

		enc.Int8(int8(k))

		switch v := v.(type) {

		default:
			plog.Fatalf("type %T not implemented", v)

		case optBooleanType:
			enc.Int8(int8(tcBoolean))
			enc.Bool(bool(v))

		case optTinyintType:
			enc.Int8(int8(tcTinyint))
			enc.Int8(int8(v))

		case optIntType:
			enc.Int8(int8(tcInteger))
			enc.Int32(int32(v))

		case optBigintType:
			enc.Int8(int8(tcBigint))
			enc.Int64(int64(v))

		case optDoubleType:
			enc.Int8(int8(tcDouble))
			enc.Float64(float64(v))

		case optStringType:
			enc.Int8(int8(tcString))
			enc.Int16(int16(len(v)))
			enc.Bytes(v)

		case optBinaryStringType:
			enc.Int8(int8(tcBstring))
			enc.Int16(int16(len(v)))
			enc.Bytes(v)
		}
	}
}
