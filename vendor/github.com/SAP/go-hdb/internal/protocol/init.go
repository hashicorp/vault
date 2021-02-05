// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

const (
	okEndianess int8 = 1
)

const (
	initRequestFillerSize = 4
)

var initRequestFiller uint32 = 0xffffffff

type version struct {
	major int8
	minor int16
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

type initRequest struct {
	product    version
	protocol   version
	numOptions int8
	endianess  endianess
}

func (r *initRequest) String() string {
	switch r.numOptions {
	default:
		return fmt.Sprintf("productVersion %s protocolVersion %s", r.product, r.protocol)
	case 1:
		return fmt.Sprintf("productVersion %s protocolVersion %s endianess %s", r.product, r.protocol, r.endianess)
	}
}

func (r *initRequest) decode(dec *encoding.Decoder) error {
	dec.Skip(initRequestFillerSize) //filler
	r.product.major = dec.Int8()
	r.product.minor = dec.Int16()
	r.protocol.major = dec.Int8()
	r.protocol.minor = dec.Int16()
	dec.Skip(1) //reserved filler
	r.numOptions = dec.Int8()

	switch r.numOptions {
	default:
		plog.Fatalf("invalid number of options %d", r.numOptions)

	case 0:
		dec.Skip(2)

	case 1:
		cnt := dec.Int8()
		if cnt != 1 {
			plog.Fatalf("endianess %d - 1 expected", cnt)
		}
		r.endianess = endianess(dec.Int8())
	}
	return dec.Error()
}

func (r *initRequest) encode(enc *encoding.Encoder) error {
	enc.Uint32(initRequestFiller)
	enc.Int8(r.product.major)
	enc.Int16(r.product.minor)
	enc.Int8(r.protocol.major)
	enc.Int16(r.protocol.minor)

	switch r.numOptions {
	default:
		plog.Fatalf("invalid number of options %d", r.numOptions)

	case 0:
		enc.Zeroes(4)

	case 1:
		// reserved
		enc.Zeroes(1)
		enc.Int8(r.numOptions)
		enc.Int8(okEndianess)
		enc.Int8(int8(r.endianess))
	}
	return nil
}

type initReply struct {
	product  version
	protocol version
}

func (r *initReply) String() string {
	return fmt.Sprintf("productVersion %s protocolVersion %s", r.product, r.protocol)
}

func (r *initReply) decode(dec *encoding.Decoder) error {
	r.product.major = dec.Int8()
	r.product.minor = dec.Int16()
	r.protocol.major = dec.Int8()
	r.protocol.minor = dec.Int16()
	dec.Skip(2) //commitInitReplySize
	return dec.Error()
}
