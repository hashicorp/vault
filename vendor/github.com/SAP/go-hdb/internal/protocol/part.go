// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"
	"math"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

const (
	partHeaderSize = 16
	// MaxNumArg is the maximum number of arguments allowed to send in a part.
	MaxNumArg = math.MaxInt16
)

type partAttributes int8

const (
	paLastPacket      partAttributes = 0x01
	paNextPacket      partAttributes = 0x02
	paFirstPacket     partAttributes = 0x04
	paRowNotFound     partAttributes = 0x08
	paResultsetClosed partAttributes = 0x10
)

var partAttributesText = map[partAttributes]string{
	paLastPacket:      "lastPacket",
	paNextPacket:      "nextPacket",
	paFirstPacket:     "firstPacket",
	paRowNotFound:     "rowNotFound",
	paResultsetClosed: "resultsetClosed",
}

func (k partAttributes) String() string {
	t := make([]string, 0, len(partAttributesText))

	for attr, text := range partAttributesText {
		if (k & attr) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

func (k partAttributes) ResultsetClosed() bool {
	return (k & paResultsetClosed) == paResultsetClosed
}

func (k partAttributes) LastPacket() bool {
	return (k & paLastPacket) == paLastPacket
}

func (k partAttributes) NoRows() bool {
	attrs := paLastPacket | paRowNotFound
	return (k & attrs) == attrs
}

// part header
type partHeader struct {
	partKind         partKind
	partAttributes   partAttributes
	argumentCount    int16
	bigArgumentCount int32
	bufferLength     int32
	bufferSize       int32
}

func (h *partHeader) String() string {
	return fmt.Sprintf("kind %s partAttributes %s argumentCount %d bigArgumentCount %d bufferLength %d bufferSize %d",
		h.partKind,
		h.partAttributes,
		h.argumentCount,
		h.bigArgumentCount,
		h.bufferLength,
		h.bufferSize,
	)
}

func (h *partHeader) setNumArg(numArg int) error {
	switch {
	default:
		return fmt.Errorf("maximum number of arguments %d exceeded", numArg)
	case numArg <= MaxNumArg:
		h.argumentCount = int16(numArg)
		h.bigArgumentCount = 0

		// TODO: seems not to work: see bulk insert test
		// case numArg <= math.MaxInt32:
		// 	s.ph.argumentCount = 0
		// 	s.ph.bigArgumentCount = int32(numArg)
		//
	}
	return nil
}

func (h *partHeader) numArg() int {
	if h.bigArgumentCount != 0 {
		panic("part header: bigArgumentCount is set")
	}
	return int(h.argumentCount)
}

func (h *partHeader) encode(enc *encoding.Encoder) error {
	enc.Int8(int8(h.partKind))
	enc.Int8(int8(h.partAttributes))
	enc.Int16(h.argumentCount)
	enc.Int32(h.bigArgumentCount)
	enc.Int32(h.bufferLength)
	enc.Int32(h.bufferSize)
	//no filler
	return nil
}

func (h *partHeader) decode(dec *encoding.Decoder) error {
	h.partKind = partKind(dec.Int8())
	h.partAttributes = partAttributes(dec.Int8())
	h.argumentCount = dec.Int16()
	h.bigArgumentCount = dec.Int32()
	h.bufferLength = dec.Int32()
	h.bufferSize = dec.Int32()
	// no filler
	return dec.Error()
}
