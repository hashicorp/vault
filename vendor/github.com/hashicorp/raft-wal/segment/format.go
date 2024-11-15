// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package segment

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/raft-wal/types"
)

const (
	// MaxEntrySize is the largest we allow any single raft log entry to be. This
	// is larger than our raft implementation ever allows so seems safe to encode
	// statically for now. We could make this configurable. It's main purpose it
	// to limit allocation when reading entries back if their lengths are
	// corrupted.
	MaxEntrySize = 64 * 1024 * 1024 // 64 MiB

	// minBufSize is the size we allocate read and write buffers. Setting it
	// larger wastes more memory but increases the chances that we'll read the
	// whole frame in a single shot and not need a second allocation and trip to
	// the disk.
	minBufSize = 64 * 1024

	fileHeaderLen = 32
	version       = 0
	magic         = 0x58eb6b0d

	// Note that this must remain a power of 2 to ensure aligning to this also
	// aligns to sector boundaries.
	frameHeaderLen = 8
)

const ( // Start iota from 0
	FrameInvalid uint8 = iota
	FrameEntry
	FrameIndex
	FrameCommit
)

var (
	// ErrTooBig indicates that the caller tried to write a logEntry with a
	// payload that's larger than we are prepared to support.
	ErrTooBig = errors.New("entries larger than 64MiB are not supported")
)

/*

  File Header functions

	0      1      2      3      4      5      6      7      8
	+------+------+------+------+------+------+------+------+
	| Magic                     | Reserved           | Vsn  |
	+------+------+------+------+------+------+------+------+
	| BaseIndex                                             |
	+------+------+------+------+------+------+------+------+
	| SegmentID                                             |
	+------+------+------+------+------+------+------+------+
	| Codec                                                 |
	+------+------+------+------+------+------+------+------+

*/

// writeFileHeader writes a file header into buf for the given file metadata.
func writeFileHeader(buf []byte, info types.SegmentInfo) error {
	if len(buf) < fileHeaderLen {
		return io.ErrShortBuffer
	}

	binary.LittleEndian.PutUint32(buf[0:4], magic)
	// Explicitly zero Reserved bytes just in case
	buf[4] = 0
	buf[5] = 0
	buf[6] = 0
	buf[7] = version
	binary.LittleEndian.PutUint64(buf[8:16], info.BaseIndex)
	binary.LittleEndian.PutUint64(buf[16:24], info.ID)
	binary.LittleEndian.PutUint64(buf[24:32], info.Codec)
	return nil
}

// readFileHeader reads a file header from buf.
func readFileHeader(buf []byte) (*types.SegmentInfo, error) {
	if len(buf) < fileHeaderLen {
		return nil, io.ErrShortBuffer
	}

	var i types.SegmentInfo
	m := binary.LittleEndian.Uint64(buf[0:8])
	if m != magic {
		return nil, types.ErrCorrupt
	}
	if buf[7] != version {
		return nil, types.ErrCorrupt
	}
	i.BaseIndex = binary.LittleEndian.Uint64(buf[8:16])
	i.ID = binary.LittleEndian.Uint64(buf[16:24])
	i.Codec = binary.LittleEndian.Uint64(buf[24:32])
	return &i, nil
}

func validateFileHeader(got, expect types.SegmentInfo) error {
	if expect.ID != got.ID {
		return fmt.Errorf("%w: segment header ID %x doesn't match metadata %x",
			types.ErrCorrupt, got.ID, expect.ID)
	}
	if expect.BaseIndex != got.BaseIndex {
		return fmt.Errorf("%w: segment header BaseIndex %d doesn't match metadata %d",
			types.ErrCorrupt, got.BaseIndex, expect.BaseIndex)
	}
	if expect.Codec != got.Codec {
		return fmt.Errorf("%w: segment header Codec %d doesn't match metadata %d",
			types.ErrCorrupt, got.Codec, expect.Codec)
	}

	return nil
}

/*
	Frame Functions

	0      1      2      3      4      5      6      7      8
	+------+------+------+------+------+------+------+------+
	| Type | Reserved           | Length/CRC                |
	+------+------+------+------+------+------+------+------+
*/

type frameHeader struct {
	typ uint8
	len uint32
	crc uint32
}

func writeFrame(buf []byte, h frameHeader, payload []byte) error {
	if len(buf) < encodedFrameSize(int(h.len)) {
		return io.ErrShortBuffer
	}
	if err := writeFrameHeader(buf, h); err != nil {
		return err
	}
	copy(buf[frameHeaderLen:], payload[:h.len])
	// Explicitly write null bytes for padding
	padBytes := padLen(int(h.len))
	for i := 0; i < padBytes; i++ {
		buf[frameHeaderLen+int(h.len)+i] = 0x0
	}
	return nil
}

func writeFrameHeader(buf []byte, h frameHeader) error {
	if len(buf) < frameHeaderLen {
		return io.ErrShortBuffer
	}
	buf[0] = h.typ
	buf[1] = 0
	buf[2] = 0
	buf[3] = 0
	lOrCRC := h.len
	if h.typ == FrameCommit {
		lOrCRC = h.crc
	}
	binary.LittleEndian.PutUint32(buf[4:8], lOrCRC)
	return nil
}

var zeroHeader [frameHeaderLen]byte

func readFrameHeader(buf []byte) (frameHeader, error) {
	var h frameHeader
	if len(buf) < frameHeaderLen {
		return h, io.ErrShortBuffer
	}

	switch buf[0] {
	default:
		return h, fmt.Errorf("%w: corrupt frame header with unknown type %d", types.ErrCorrupt, buf[0])

	case FrameInvalid:
		// Check if the whole header is zero and return a zero frame as this could
		// just indicate we've read right off the end of the written data during
		// recovery.
		if bytes.Equal(buf[:frameHeaderLen], zeroHeader[:]) {
			return h, nil
		}
		return h, fmt.Errorf("%w: corrupt frame header with type 0 but non-zero other fields", types.ErrCorrupt)

	case FrameEntry, FrameIndex:
		h.typ = buf[0]
		h.len = binary.LittleEndian.Uint32(buf[4:8])

	case FrameCommit:
		h.typ = buf[0]
		h.crc = binary.LittleEndian.Uint32(buf[4:8])
	}
	return h, nil
}

// padLen returns how many bytes of padding should be added to a frame of length
// n to ensure it is a multiple of headerLen. We ensure frameHeaderLen is a
// power of two so that it's always a multiple of a typical sector size (e.g.
// 512 bytes) to reduce the risk that headers are torn by being written across
// sector boundaries. It will return an int in the range [0, 7].
func padLen(n int) int {
	// This looks a bit awful but it's just doing (n % 8) and subtracting that
	// from 8 to get the number of bytes extra needed to get up to the next 8-byte
	// boundary. The extra & 7 is to handle the case where n is a multiple of 8
	// already and so n%8 is 0 and 8-0 is 8. By &ing 8 (0b1000) with 7 (0b111) we
	// effectively wrap it back around to 0. This only works as long as
	// frameHeaderLen is a power of 2 but that's necessary per comment above.
	return (frameHeaderLen - (n % frameHeaderLen)) & (frameHeaderLen - 1)
}

func encodedFrameSize(payloadLen int) int {
	return frameHeaderLen + payloadLen + padLen(payloadLen)
}

func indexFrameSize(numEntries int) int {
	// Index frames are completely unnecessary if the whole block is a
	// continuation with no new entries.
	if numEntries == 0 {
		return 0
	}
	return encodedFrameSize(numEntries * 4)
}

func writeIndexFrame(buf []byte, offsets []uint32) error {
	if len(buf) < indexFrameSize(len(offsets)) {
		return io.ErrShortBuffer
	}
	fh := frameHeader{
		typ: FrameIndex,
		len: uint32(len(offsets) * 4),
	}
	if err := writeFrameHeader(buf, fh); err != nil {
		return err
	}
	cursor := frameHeaderLen
	for _, o := range offsets {
		binary.LittleEndian.PutUint32(buf[cursor:], o)
		cursor += 4
	}
	if (len(offsets) % 2) == 1 {
		// Odd number of entries, zero pad to keep it 8-byte aligned
		binary.LittleEndian.PutUint32(buf[cursor:], 0)
	}
	return nil
}
