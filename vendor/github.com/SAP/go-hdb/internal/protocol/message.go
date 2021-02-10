// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

//message header (size: 32 bytes)
type messageHeader struct {
	sessionID     int64
	packetCount   int32
	varPartLength uint32
	varPartSize   uint32
	noOfSegm      int16
}

func (h *messageHeader) String() string {
	return fmt.Sprintf("session id %d packetCount %d varPartLength %d, varPartSize %d noOfSegm %d",
		h.sessionID,
		h.packetCount,
		h.varPartLength,
		h.varPartSize,
		h.noOfSegm)
}

func (h *messageHeader) encode(enc *encoding.Encoder) error {
	enc.Int64(h.sessionID)
	enc.Int32(h.packetCount)
	enc.Uint32(h.varPartLength)
	enc.Uint32(h.varPartSize)
	enc.Int16(h.noOfSegm)
	enc.Zeroes(10) // size: 32 bytes
	return nil
}

func (h *messageHeader) decode(dec *encoding.Decoder) error {
	h.sessionID = dec.Int64()
	h.packetCount = dec.Int32()
	h.varPartLength = dec.Uint32()
	h.varPartSize = dec.Uint32()
	h.noOfSegm = dec.Int16()
	dec.Skip(10) // size: 32 bytes
	return dec.Error()
}
