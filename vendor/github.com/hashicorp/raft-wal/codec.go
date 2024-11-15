// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package wal

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/hashicorp/raft"
)

const (
	// FirstExternalCodecID is the lowest value an external code may use to
	// identify their codec. Values lower than this are reserved for future
	// internal use.
	FirstExternalCodecID = 1 << 16

	// Codec* constants identify internally-defined codec identifiers.
	CodecBinaryV1 uint64 = iota
)

// Codec is the interface required for encoding/decoding log entries. Callers
// can pass a custom one to manage their own serialization, or to add additional
// layers like encryption or compression of records. Each codec
type Codec interface {
	// ID returns the globally unique identifier for this codec version. This is
	// encoded into segment file headers and must remain consistent over the life
	// of the log. Values up to FirstExternalCodecID are reserved and will error
	// if specified externally.
	ID() uint64

	// Encode the log into the io.Writer. We pass a writer to allow the caller to
	// manage buffer allocation and re-use.
	Encode(l *raft.Log, w io.Writer) error

	// Decode a log from the passed byte slice into the log entry pointed to. This
	// allows the caller to manage allocation and re-use of the bytes and log
	// entry. The resulting raft.Log MUST NOT reference data in the input byte
	// slice since the input byte slice may be returned to a pool and re-used.
	Decode([]byte, *raft.Log) error
}

// BinaryCodec is a Codec that encodes raft.Log with a simple binary format. We
// test that all fields are captured using reflection.
//
// For now we assume raft.Log is not likely to change too much. If it does we'll
// use a new Codec ID for the later version and have to support decoding either.
type BinaryCodec struct{}

// ID returns the globally unique identifier for this codec version. This is
// encoded into segment file headers and must remain consistent over the life
// of the log. Values up to FirstExternalCodecID are reserved and will error
// if specified externally.
func (c *BinaryCodec) ID() uint64 {
	return CodecBinaryV1
}

// Encode the log into the io.Writer. We pass a writer to allow the caller to
// manage buffer allocation and re-use.
func (c *BinaryCodec) Encode(l *raft.Log, w io.Writer) error {
	enc := encoder{w: w}
	enc.varint(l.Index)
	enc.varint(l.Term)
	enc.varint(uint64(l.Type))
	enc.bytes(l.Data)
	enc.bytes(l.Extensions)
	enc.time(l.AppendedAt)
	return enc.err
}

// Decode a log from the passed byte slice into the log entry pointed to. This
// allows the caller to manage allocation and re-use of the bytes and log
// entry.
func (c *BinaryCodec) Decode(bs []byte, l *raft.Log) error {
	dec := decoder{buf: bs}
	l.Index = dec.varint()
	l.Term = dec.varint()
	l.Type = raft.LogType(dec.varint())
	l.Data = dec.bytes()
	l.Extensions = dec.bytes()
	l.AppendedAt = dec.time()
	return dec.err
}

type encoder struct {
	w       io.Writer
	err     error
	scratch [10]byte
}

func (e *encoder) varint(v uint64) {
	if e.err != nil {
		return
	}

	// Varint encoding might use up to 9 bytes for a uint64
	n := binary.PutUvarint(e.scratch[:], v)
	_, e.err = e.w.Write(e.scratch[:n])
}

func (e *encoder) bytes(bs []byte) {
	// Put a length prefix
	e.varint(uint64(len(bs)))
	if e.err != nil {
		return
	}
	// Copy the bytes to the writer
	_, e.err = e.w.Write(bs)
}

func (e *encoder) time(t time.Time) {
	if e.err != nil {
		return
	}
	bs, err := t.MarshalBinary()
	if err != nil {
		e.err = err
		return
	}
	_, e.err = e.w.Write(bs)
}

type decoder struct {
	buf []byte
	err error
}

func (d *decoder) varint() uint64 {
	if d.err != nil {
		return 0
	}
	v, n := binary.Uvarint(d.buf)
	d.buf = d.buf[n:]
	return v
}

func (d *decoder) bytes() []byte {
	// Get length prefix
	n := d.varint()
	if d.err != nil {
		return nil
	}
	if n == 0 {
		return nil
	}
	if n > uint64(len(d.buf)) {
		d.err = io.ErrShortBuffer
		return nil
	}
	bs := make([]byte, n)
	copy(bs, d.buf[:n])
	d.buf = d.buf[n:]
	return bs
}

func (d *decoder) time() time.Time {
	var t time.Time
	if d.err != nil {
		return t
	}
	// Note that Unmarshal Binary updates d.buf to remove the bytes it read
	// already.
	d.err = t.UnmarshalBinary(d.buf)
	return t
}
