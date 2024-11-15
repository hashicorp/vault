// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package segment

import (
	"fmt"
	"hash/crc32"
	"io"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/raft-wal/types"
)

// Writer allows appending logs to a segment file as well as reading them back.
type Writer struct {
	// commitIdx is updated after an append batch is fully persisted to disk to
	// allow readers to read the new value. Note that readers must not read values
	// larger than this even if they are available in tailIndex as they are not
	// yet committed to disk!
	commitIdx uint64

	// offsets is the index offset. The first element corresponds to the
	// BaseIndex. It is accessed concurrently by readers and the single writer
	// without locks! This is race-free via the following invariants:
	//  - the slice here is never mutated only copied though it may still refer to
	//    the same backing array.
	//  - readers only ever read up to len(offsets) in the atomically accessed
	//    slice. Those elements of the backing array are immutable and will never
	//    be modified once they are accessible to readers.
	//  - readers and writers synchronize on atomic access to the slice
	//  - serial writer will only append to the end which either mutates the
	//    shared backing array but at an index greater than the len any reader has
	//    seen, or a new backing array is allocated and the old one copied into it
	//    which also will never mutate the entries readers can already "see" via
	//    the old slice.
	offsets atomic.Value // []uint32

	// writer state is accessed only on the (serial) write path so doesn't need
	// synchronization.
	writer struct {
		// commitBuf stores the pending frames waiting to be flushed to the current
		// tail block.
		commitBuf []byte

		// crc is the rolling crc32 Castagnoli sum of all data written since the
		// last fsync.
		crc uint32

		// writeOffset is the absolute file offset up to which we've written data to
		// the file. The contents of commitBuf will be written at this offset when
		// it commits or we reach the end of the block, whichever happens first.
		writeOffset uint32

		// indexStart is set when the tail is sealed indicating the file offset at
		// which the index array was written.
		indexStart uint64
	}

	info types.SegmentInfo
	wf   types.WritableFile
	r    types.SegmentReader
}

func createFile(info types.SegmentInfo, wf types.WritableFile, bufPool *sync.Pool) (*Writer, error) {
	r, err := openReader(info, wf, bufPool)
	if err != nil {
		return nil, err
	}
	w := &Writer{
		info: info,
		wf:   wf,
		r:    r,
	}
	r.tail = w
	if err := w.initEmpty(); err != nil {
		return nil, err
	}
	return w, nil
}

func recoverFile(info types.SegmentInfo, wf types.WritableFile, bufPool *sync.Pool) (*Writer, error) {
	r, err := openReader(info, wf, bufPool)
	if err != nil {
		return nil, err
	}
	w := &Writer{
		info: info,
		wf:   wf,
		r:    r,
	}
	r.tail = w

	if err := w.recoverTail(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Writer) initEmpty() error {
	// Write header into write buffer to be written out with the first commit.
	w.writer.writeOffset = 0
	w.ensureBufCap(fileHeaderLen)
	w.writer.commitBuf = w.writer.commitBuf[:fileHeaderLen]

	if err := writeFileHeader(w.writer.commitBuf, w.info); err != nil {
		return err
	}

	w.writer.crc = crc32.Checksum(w.writer.commitBuf[:fileHeaderLen], castagnoliTable)

	// Initialize the index
	offsets := make([]uint32, 0, 32*1024)
	w.offsets.Store(offsets)
	return nil
}

func (w *Writer) recoverTail() error {
	// We need to track the last two commit frames
	type commitInfo struct {
		fh         frameHeader
		offset     int64
		crcStart   int64
		offsetsLen int
	}
	var prevCommit, finalCommit *commitInfo

	offsets := make([]uint32, 0, 32*1024)

	readInfo, err := readThroughSegment(w.wf, func(_ types.SegmentInfo, fh frameHeader, offset int64) (bool, error) {
		switch fh.typ {
		case FrameEntry:
			// Record the frame offset
			offsets = append(offsets, uint32(offset))

		case FrameIndex:
			// So this segment was sealed! (or attempted) keep track of this
			// indexStart in case it turns out the Seal actually committed completely.
			// We store the start of the actual array not the frame header.
			w.writer.indexStart = uint64(offset) + frameHeaderLen

		case FrameCommit:
			// The payload is not the length field in this case!
			prevCommit = finalCommit
			finalCommit = &commitInfo{
				fh:         fh,
				offset:     offset,
				crcStart:   0,            // First commit includes the file header
				offsetsLen: len(offsets), // Track how many entries were found up to this commit point.
			}
			if prevCommit != nil {
				finalCommit.crcStart = prevCommit.offset + frameHeaderLen
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	if finalCommit == nil {
		// There were no commit frames found at all. This segment file is
		// effectively empty. Init it that way ready for appending. This overwrites
		// the file header so it doesn't matter if it was valid or not.
		return w.initEmpty()
	}

	// Assume that the final commit is good for now and set the writer state
	w.writer.writeOffset = uint32(finalCommit.offset + frameHeaderLen)

	// Just store what we have for now to ensure the defer doesn't panic we'll
	// probably update this below.
	w.offsets.Store(offsets)

	// Whichever path we take, fix up the commitIdx before we leave
	defer func() {
		ofs := w.getOffsets()
		if len(ofs) > 0 {
			// Non atomic is OK because this file is not visible to any other threads
			// yet.
			w.commitIdx = w.info.BaseIndex + uint64(len(ofs)) - 1
		}
	}()

	if finalCommit.offsetsLen < len(offsets) {
		// Some entries were found after the last commit. Those must be a partial
		// write that was uncommitted so can be ignored. But the fact they were
		// written at all means that the last commit frame must have been completed
		// and acknowledged so we don't need to verify anything. Just truncate the
		// extra entries from index and reset the write cursor to continue appending
		// after the last commit.
		offsets = offsets[:finalCommit.offsetsLen]
		w.offsets.Store(offsets)

		// Since at least one commit was found, the header better be valid!
		return validateFileHeader(*readInfo, w.info)
	}

	// Last frame was a commit frame! Let's check that all the data written in
	// that commit frame made it to disk.
	// Verify the length first
	bufLen := finalCommit.offset - finalCommit.crcStart
	// We know bufLen can't be bigger than the whole segment file because none of
	// the values above were read from the data just from the offsets we moved
	// through.
	batchBuf := make([]byte, bufLen)

	if _, err := w.wf.ReadAt(batchBuf, finalCommit.crcStart); err != nil {
		return fmt.Errorf("failed to read last committed batch for CRC validation: %w", err)
	}

	gotCrc := crc32.Checksum(batchBuf, castagnoliTable)
	if gotCrc == finalCommit.fh.crc {
		// All is good. We already setup the state we need for writer other than
		// offsets.
		w.offsets.Store(offsets)

		// Since at least one commit was found, the header better be valid!
		return validateFileHeader(*readInfo, w.info)
	}

	// Last commit was incomplete rewind back to the previous one or start of file
	if prevCommit == nil {
		// Init wil re-write the file header so it doesn't matter if it was corrupt
		// or not!
		return w.initEmpty()
	}

	w.writer.writeOffset = uint32(prevCommit.offset + frameHeaderLen)
	offsets = offsets[:prevCommit.offsetsLen]
	w.offsets.Store(offsets)

	// Since at least one commit was found, the header better be valid!
	return validateFileHeader(*readInfo, w.info)
}

// Close implements io.Closer
func (w *Writer) Close() error {
	return w.r.Close()
}

// GetLog implements types.SegmentReader
func (w *Writer) GetLog(idx uint64) (*types.PooledBuffer, error) {
	return w.r.GetLog(idx)
}

// Append adds one or more entries. It must not return until the entries are
// durably stored otherwise raft's guarantees will be compromised.
func (w *Writer) Append(entries []types.LogEntry) error {
	if len(entries) < 1 {
		return nil
	}

	if w.writer.indexStart > 0 {
		return types.ErrSealed
	}

	flushed := false

	// Save any state we may need to rollback.
	beforeBuf := w.writer.commitBuf
	beforeCRC := w.writer.crc
	beforeIndexStart := w.writer.indexStart
	beforeWriteOffset := w.writer.writeOffset
	beforeOffsets := w.offsets.Load()

	defer func() {
		if !flushed {
			// rollback writer state on error
			w.writer.commitBuf = beforeBuf
			w.writer.crc = beforeCRC
			w.writer.indexStart = beforeIndexStart
			w.writer.writeOffset = beforeWriteOffset
			w.offsets.Store(beforeOffsets)
		}
	}()

	// Iterate entries and append each one
	for _, e := range entries {
		if err := w.appendEntry(e); err != nil {
			return err
		}
	}

	ofs := w.getOffsets()
	// Work out if we need to seal before we commit and sync.
	if (w.writer.writeOffset + uint32(len(w.writer.commitBuf)+indexFrameSize(len(ofs)))) > w.info.SizeLimit {
		// Seal the segment! We seal it by writing an index frame before we commit.
		if err := w.appendIndex(); err != nil {
			return err
		}
	}

	// Write the commit frame
	if err := w.appendCommit(); err != nil {
		return err
	}

	flushed = true

	// Commit in-memory
	atomic.StoreUint64(&w.commitIdx, entries[len(entries)-1].Index)
	return nil
}

func (w *Writer) getOffsets() []uint32 {
	return w.offsets.Load().([]uint32)
}

// OffsetForFrame implements tailWriter and allows readers to lookup entry
// frames in the tail's in-memory index.
func (w *Writer) OffsetForFrame(idx uint64) (uint32, error) {
	if idx < w.info.BaseIndex || idx < w.info.MinIndex || idx > w.LastIndex() {
		return 0, types.ErrNotFound
	}
	os := w.getOffsets()
	entryIndex := idx - w.info.BaseIndex
	// No bounds check on entryIndex since LastIndex must ensure it's in bounds.
	return os[entryIndex], nil
}

func (w *Writer) appendEntry(e types.LogEntry) error {
	offsets := w.getOffsets()

	// Check the invariant that this entry is the next one we expect otherwise our
	// index logic is incorrect and will result in panics on read.
	if e.Index != w.info.BaseIndex+uint64(len(offsets)) {
		return fmt.Errorf("non-monotonic append to segment with BaseIndex=%d. Entry index %d, expected %d",
			w.info.BaseIndex, e.Index, w.info.BaseIndex+uint64(len(offsets)))
	}

	fh := frameHeader{
		typ: FrameEntry,
		len: uint32(len(e.Data)),
	}
	bufOffset, err := w.appendFrame(fh, e.Data)
	if err != nil {
		return err
	}
	// Update the offsets index

	// Add the index entry. Note this is safe despite mutating the same backing
	// array as tail because it's beyond the limit current readers will access
	// until we do the atomic update below. Even if append re-allocates the
	// backing array, it will only read the indexes smaller than numEntries from
	// the old array to copy them into the new one and we are not mutating the
	// same memory locations. Old readers might still be looking at the old
	// array (lower than numEntries) through the current tail.offsets slice but
	// we are not touching that at least below numEntries.
	offsets = append(offsets, w.writer.writeOffset+uint32(bufOffset))

	// Now we can make it available to readers. Note that readers still
	// shouldn't read it until we actually commit to disk (and increment
	// commitIdx) but it's race free for them to now!
	w.offsets.Store(offsets)
	return nil
}

func (w *Writer) appendCommit() error {
	fh := frameHeader{
		typ: FrameCommit,
		crc: w.writer.crc,
	}
	if _, err := w.appendFrame(fh, nil); err != nil {
		return err
	}

	// Flush all writes to the file
	if err := w.sync(); err != nil {
		return err
	}

	// Finally, reset crc so that by the time we write the next trailer
	// we'll know where the append batch started.
	w.writer.crc = 0
	return nil
}

func (w *Writer) ensureBufCap(extraLen int) {
	needCap := len(w.writer.commitBuf) + extraLen
	if cap(w.writer.commitBuf) < needCap {
		newSize := minBufSize
		// Double buffer size until it's big enough to amortize cost
		for newSize < needCap {
			newSize = newSize * 2
		}
		newBuf := make([]byte, newSize)
		oldLen := len(w.writer.commitBuf)
		copy(newBuf, w.writer.commitBuf)
		w.writer.commitBuf = newBuf[:oldLen]
	}
}

func (w *Writer) appendIndex() error {
	// Append the index record before we commit (commit and flush happen later
	// generally)
	offsets := w.getOffsets()
	l := indexFrameSize(len(offsets))
	w.ensureBufCap(l)

	startOff := len(w.writer.commitBuf)

	if err := writeIndexFrame(w.writer.commitBuf[startOff:startOff+l], offsets); err != nil {
		return err
	}
	w.writer.commitBuf = w.writer.commitBuf[:startOff+l]

	// Update crc with those values
	w.writer.crc = crc32.Update(w.writer.crc, castagnoliTable, w.writer.commitBuf[startOff:startOff+l])

	// Record the file offset where the index starts (the actual index data so
	// after the frame header).
	w.writer.indexStart = uint64(w.writer.writeOffset) + uint64(startOff+frameHeaderLen)
	return nil
}

// appendFrame appends the given frame to the current block. The frame must fit
// already otherwise an error will be returned.
func (w *Writer) appendFrame(fh frameHeader, data []byte) (int, error) {
	// Encode frame header into current block buffer
	l := encodedFrameSize(len(data))
	w.ensureBufCap(l)

	bufOffset := len(w.writer.commitBuf)
	if err := writeFrame(w.writer.commitBuf[bufOffset:bufOffset+l], fh, data); err != nil {
		return 0, err
	}
	// Update len of commitBuf since we resliced it for the write
	w.writer.commitBuf = w.writer.commitBuf[:bufOffset+l]

	// Update the CRC
	w.writer.crc = crc32.Update(w.writer.crc, castagnoliTable, w.writer.commitBuf[bufOffset:bufOffset+l])
	return bufOffset, nil
}

func (w *Writer) flush() error {
	// Write to file
	n, err := w.wf.WriteAt(w.writer.commitBuf, int64(w.writer.writeOffset))
	if err == io.EOF && n == len(w.writer.commitBuf) {
		// Writer may return EOF even if it wrote all bytes if it wrote right up to
		// the end of the file. Ignore that case though.
		err = nil
	}
	if err != nil {
		return err
	}

	// Reset writer state ready for next writes
	w.writer.writeOffset += uint32(len(w.writer.commitBuf))
	w.writer.commitBuf = w.writer.commitBuf[:0]
	return nil
}

func (w *Writer) sync() error {
	// Write out current buffer to file
	if err := w.flush(); err != nil {
		return err
	}

	// Sync file
	if err := w.wf.Sync(); err != nil {
		return err
	}

	// Update commitIdx atomically
	offsets := w.getOffsets()
	commitIdx := uint64(0)
	if len(offsets) > 0 {
		// Probably not possible for the to be less, but just in case we ever flush
		// the file with only meta data written...
		commitIdx = uint64(w.info.BaseIndex) + uint64(len(offsets)) - 1
	}
	atomic.StoreUint64(&w.commitIdx, commitIdx)
	return nil
}

// Sealed returns whether the segment is sealed or not. If it is it returns
// true and the file offset that it's index array starts at to be saved in
// meta data. WAL will call this after every append so it should be relatively
// cheap in the common case. This design allows the final Append to write out
// the index or any additional data needed at seal time in the same fsync.
func (w *Writer) Sealed() (bool, uint64, error) {
	if w.writer.indexStart == 0 {
		return false, 0, nil
	}
	return true, w.writer.indexStart, nil
}

// ForceSeal forces us to seal the segment by writing out an index block
// wherever we got to in the file. After calling this it is no longer valid to
// call Append on this file.
func (w *Writer) ForceSeal() (uint64, error) {
	if w.writer.indexStart > 0 {
		// Already sealed, this is a no-op.
		return w.writer.indexStart, nil
	}

	// Seal the segment! We seal it by writing an index frame before we commit.
	if err := w.appendIndex(); err != nil {
		return 0, err
	}

	// Write the commit frame
	if err := w.appendCommit(); err != nil {
		return 0, err
	}

	return w.writer.indexStart, nil
}

// LastIndex returns the most recently persisted index in the log. It must
// respond without blocking on append since it's needed frequently by read
// paths that may call it concurrently. Typically this will be loaded from an
// atomic int. If the segment is empty lastIndex should return zero.
func (w *Writer) LastIndex() uint64 {
	return atomic.LoadUint64(&w.commitIdx)
}

func readThroughSegment(r types.ReadableFile, fn func(info types.SegmentInfo, fh frameHeader, offset int64) (bool, error)) (*types.SegmentInfo, error) {
	// First read the file header. Note we wrote it as part of the first commit so
	// it may be missing or partial written and that's OK as long as there are no
	// other later commit frames!
	var fh [fileHeaderLen]byte
	_, err := r.ReadAt(fh[:], 0)
	// EOF is ok - the file might be empty if we crashed before committing
	// anything and preallocation isn't supported.
	if err != io.EOF && err != nil {
		return nil, err
	}

	readInfo, err := readFileHeader(fh[:])
	if err == types.ErrCorrupt {
		// Header is malformed or missing, don't error yet though we'll detect it
		// later when we know if it's a problem or not.
		err = nil
	}
	if err != nil {
		return nil, err
	}
	// If header wasn't detected as corrupt, it might still be just in a way
	// that's valid since we've not verified it against the expected metadata yet.
	// We'll wait to see if the header was part of the last commit before decide
	// if we should validate it for corruption or not though. For now just make
	// sure it's not nil so we don't have to handle nil checks everywhere.
	if readInfo == nil {
		// Zero info will fail validation against the actual metadata if it was
		// corrupt when it shouldn't be later. Just prevents a nil panic.
		readInfo = &types.SegmentInfo{}
	}

	// Read through file from after header until we hit zeros, EOF or corrupt
	// frames.
	offset := int64(fileHeaderLen)
	var buf [frameHeaderLen]byte

	for {
		n, err := r.ReadAt(buf[:], offset)
		if err == io.EOF {
			if n < frameHeaderLen {
				return readInfo, nil
			}
			// This is OK! The last frame in file might be a commit frame so as long
			// as we have it all then we can ignore the EOF for this iteration.
			err = nil
		}
		if err != nil {
			return readInfo, fmt.Errorf("failed reading frame at offset=%d: %w", offset, err)
		}
		fh, err := readFrameHeader(buf[:frameHeaderLen])
		if err != nil {
			// This is not actually an error case. If we failed to decode it could be
			// because of a torn write (since we don't assume writes are atomic). We
			// assume that previously committed data is not silently corrupted by the
			// FS (see README for details). So this must be due to corruption that
			// happened due to non-atomic sector updates whilst committing the last
			// write batch.
			return readInfo, nil
		}
		if fh.typ == FrameInvalid {
			// This means we've hit zeros at the end of the file (or due to an
			// incomplete write, which we treat the same way).
			return readInfo, nil
		}

		// Call the callback
		shouldContinue, err := fn(*readInfo, fh, offset)
		if err != nil {
			return readInfo, err
		}
		if !shouldContinue {
			return readInfo, nil
		}

		// Skip to next frame
		offset += int64(encodedFrameSize(int(fh.len)))
	}
}
