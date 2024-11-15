// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package segment

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft-wal/types"
)

// Reader allows reading logs from a segment file.
type Reader struct {
	info types.SegmentInfo
	rf   types.ReadableFile

	bufPool *sync.Pool

	// tail optionally providers an interface to the writer state when this is an
	// unsealed segment so we can fetch from it's in-memory index.
	tail tailWriter
}

type tailWriter interface {
	OffsetForFrame(idx uint64) (uint32, error)
}

func openReader(info types.SegmentInfo, rf types.ReadableFile, bufPool *sync.Pool) (*Reader, error) {
	r := &Reader{
		info:    info,
		rf:      rf,
		bufPool: bufPool,
	}

	return r, nil
}

// Close implements io.Closer
func (r *Reader) Close() error {
	return r.rf.Close()
}

// GetLog returns the raw log entry bytes associated with idx. If the log
// doesn't exist in this segment types.ErrNotFound must be returned.
func (r *Reader) GetLog(idx uint64) (*types.PooledBuffer, error) {
	offset, err := r.findFrameOffset(idx)
	if err != nil {
		return nil, err
	}

	_, payload, err := r.readFrame(offset)
	if err != nil {
		return nil, err
	}
	return payload, err
}

func (r *Reader) readFrame(offset uint32) (frameHeader, *types.PooledBuffer, error) {
	buf := r.makeBuffer()

	n, err := r.rf.ReadAt(buf.Bs, int64(offset))
	if errors.Is(err, io.EOF) && n >= frameHeaderLen {
		// We might have hit EOF just because our read buffer (at least 64KiB) might
		// be larger than the space left in the file (say if files are tiny or if we
		// are reading a frame near the end.). So don't treat EOF as an error as
		// long as we have actually managed to read a frameHeader - we'll work out
		// if we got the whole thing or not below.
		err = nil

		// Re-slice buf.Bs so it's len() reflect only what we actually managed to
		// read. Note this doesn't impact the buffer length when it's returned to
		// the pool which will still return the whole cap.
		buf.Bs = buf.Bs[:n]
	}
	if err != nil {
		return frameHeader{}, nil, err
	}
	fh, err := readFrameHeader(buf.Bs)
	if err != nil {
		return fh, nil, err
	}

	if (frameHeaderLen + int(fh.len)) <= len(buf.Bs) {
		// We already have all we need read, just return it sliced to just include
		// the payload.
		buf.Bs = buf.Bs[frameHeaderLen : frameHeaderLen+fh.len]
		return fh, buf, nil
	}
	// Need to read again, with a bigger buffer, return this one
	buf.Close()

	// Need to read more bytes, validate that len is a sensible number
	if fh.len > MaxEntrySize {
		return fh, nil, fmt.Errorf("%w: frame header indicates a record larger than MaxEntrySize (%d bytes)", types.ErrCorrupt, MaxEntrySize)
	}

	buf = &types.PooledBuffer{
		Bs: make([]byte, fh.len),
		// No closer, let outsized buffers be GCed in case they are massive and way
		// bigger than we need again. Could reconsider this if we find we need to
		// optimize for frequent > minBufSize reads.
	}
	if _, err := r.rf.ReadAt(buf.Bs, int64(offset+frameHeaderLen)); err != nil {
		return fh, nil, err
	}
	return fh, buf, nil
}

func (r *Reader) makeBuffer() *types.PooledBuffer {
	if r.bufPool == nil {
		return &types.PooledBuffer{Bs: make([]byte, minBufSize)}
	}
	buf := r.bufPool.Get().([]byte)
	return &types.PooledBuffer{
		Bs: buf,
		CloseFn: func() {
			// Note we always return the whole allocated buf regardless of what Bs
			// ended up being sliced to.
			r.bufPool.Put(buf)
		},
	}

}

func (r *Reader) findFrameOffset(idx uint64) (uint32, error) {
	if r.tail != nil {
		// This is not a sealed segment.
		return r.tail.OffsetForFrame(idx)
	}

	// Sealed segment, read from the on-disk index block.
	if r.info.IndexStart == 0 {
		return 0, fmt.Errorf("sealed segment has no index block")
	}

	if idx < r.info.MinIndex || (r.info.MaxIndex > 0 && idx > r.info.MaxIndex) {
		return 0, types.ErrNotFound
	}

	// IndexStart is the offset to the first entry in the index array. We need to
	// find the byte offset to the Nth entry
	entryOffset := (idx - r.info.BaseIndex)
	byteOffset := r.info.IndexStart + (entryOffset * 4)

	var bs [4]byte
	n, err := r.rf.ReadAt(bs[:], int64(byteOffset))
	if err == io.EOF && n == 4 {
		// Read all of it just happened to be at end of file, ignore
		err = nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to read segment index: %w", err)
	}
	offset := binary.LittleEndian.Uint32(bs[:])
	return offset, nil
}
