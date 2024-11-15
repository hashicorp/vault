// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package segment

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/raft-wal/types"
)

const (
	segmentFileSuffix      = ".wal"
	segmentFileNamePattern = "%020d-%016x" + segmentFileSuffix
)

// Filer implements the abstraction for managing a set of segment files in a
// directory. It uses a VFS to abstract actual file system operations for easier
// testing.
type Filer struct {
	dir     string
	vfs     types.VFS
	bufPool sync.Pool
}

// NewFiler creates a Filer ready for use.
func NewFiler(dir string, vfs types.VFS) *Filer {
	f := &Filer{
		dir: dir,
		vfs: vfs,
	}
	f.bufPool.New = func() interface{} {
		return make([]byte, minBufSize)
	}
	return f
}

// FileName returns the formatted file name expected for this segment.
// SegmentFiler implementations could choose to ignore this but it's here to
func FileName(i types.SegmentInfo) string {
	return fmt.Sprintf(segmentFileNamePattern, i.BaseIndex, i.ID)
}

// Create adds a new segment with the given info and returns a writer or an
// error.
func (f *Filer) Create(info types.SegmentInfo) (types.SegmentWriter, error) {
	if info.BaseIndex == 0 {
		return nil, fmt.Errorf("BaseIndex must be greater than zero")
	}
	fname := FileName(info)

	wf, err := f.vfs.Create(f.dir, fname, uint64(info.SizeLimit))
	if err != nil {
		return nil, err
	}

	return createFile(info, wf, &f.bufPool)
}

// RecoverTail is called on an unsealed segment when re-opening the WAL it will
// attempt to recover from a possible crash. It will either return an error, or
// return a valid segmentWriter that is ready for further appends. If the
// expected tail segment doesn't exist it must return an error wrapping
// os.ErrNotExist.
func (f *Filer) RecoverTail(info types.SegmentInfo) (types.SegmentWriter, error) {
	fname := FileName(info)

	wf, err := f.vfs.OpenWriter(f.dir, fname)
	if err != nil {
		return nil, err
	}

	return recoverFile(info, wf, &f.bufPool)
}

// Open an already sealed segment for reading. Open may validate the file's
// header and return an error if it doesn't match the expected info.
func (f *Filer) Open(info types.SegmentInfo) (types.SegmentReader, error) {
	fname := FileName(info)

	rf, err := f.vfs.OpenReader(f.dir, fname)
	if err != nil {
		return nil, err
	}

	// Validate header here since openReader is re-used by writer where it's valid
	// for the file header not to be committed yet after a crash so we can't check
	// it there.
	var hdr [fileHeaderLen]byte

	if _, err := rf.ReadAt(hdr[:], 0); err != nil {
		if errors.Is(err, io.EOF) {
			// Treat failure to read a header as corruption since a sealed file should
			// never not have a valid header. (I.e. even if crashes happen it should
			// be impossible to seal a segment with no header written so this
			// indicates that something truncated the file after the fact)
			return nil, fmt.Errorf("%w: failed to read header: %s", types.ErrCorrupt, err)
		}
		return nil, err
	}

	gotInfo, err := readFileHeader(hdr[:])
	if err != nil {
		return nil, err
	}

	if err := validateFileHeader(*gotInfo, info); err != nil {
		return nil, err
	}

	return openReader(info, rf, &f.bufPool)
}

// List returns the set of segment IDs currently stored. It's used by the WAL
// on recovery to find any segment files that need to be deleted following a
// unclean shutdown. The returned map is a map of ID -> BaseIndex. BaseIndex
// is returned to allow subsequent Delete calls to be made.
func (f *Filer) List() (map[uint64]uint64, error) {
	segs, _, err := f.listInternal()
	return segs, err
}

func (f *Filer) listInternal() (map[uint64]uint64, []uint64, error) {
	files, err := f.vfs.ListDir(f.dir)
	if err != nil {
		return nil, nil, err
	}

	segs := make(map[uint64]uint64)
	sorted := make([]uint64, 0)
	for _, file := range files {
		if !strings.HasSuffix(file, segmentFileSuffix) {
			continue
		}
		// Parse BaseIndex and ID from the file name
		var bIdx, id uint64
		n, err := fmt.Sscanf(file, segmentFileNamePattern, &bIdx, &id)
		if err != nil {
			return nil, nil, types.ErrCorrupt
		}
		if n != 2 {
			// Misnamed segment files with the right suffix indicates a bug or
			// tampering, we can't be sure what's happened to the data.
			return nil, nil, types.ErrCorrupt
		}
		segs[id] = bIdx
		sorted = append(sorted, id)
	}

	return segs, sorted, nil
}

// Delete removes the segment with given baseIndex and id if it exists. Note
// that baseIndex is technically redundant since ID is unique on it's own. But
// in practice we name files (or keys) with both so that they sort correctly.
// This interface allows a  simpler implementation where we can just delete
// the file if it exists without having to scan the underlying storage for a.
func (f *Filer) Delete(baseIndex uint64, ID uint64) error {
	fname := fmt.Sprintf(segmentFileNamePattern, baseIndex, ID)
	return f.vfs.Delete(f.dir, fname)
}

// DumpSegment attempts to read the segment file specified by the baseIndex and
// ID. It's intended purpose is for debugging the contents of segment files and
// unlike the SegmentFiler interface, it doesn't assume the caller has access to
// the correct metadata. This allows dumping log segments in a WAL that is still
// being written to by another process. Without metadata we don't know if the
// file is sealed so always recover by reading through the whole file. If after
// or before are non-zero, the specify a exclusive lower or upper bound on which
// log entries should be emitted. No error checking is done on the read data. fn
// is called for each entry passing the raft info read from the file header (so
// that the caller knows which codec to use for example) the raft index of the
// entry and the raw bytes of the entry itself. The callback must return true to
// continue reading. The data slice is only valid for the lifetime of the call.
func (f *Filer) DumpSegment(baseIndex uint64, ID uint64, after, before uint64, fn func(info types.SegmentInfo, e types.LogEntry) (bool, error)) error {
	fname := fmt.Sprintf(segmentFileNamePattern, baseIndex, ID)

	rf, err := f.vfs.OpenReader(f.dir, fname)
	if err != nil {
		return err
	}

	buf := make([]byte, 64*1024)
	idx := baseIndex

	type frameInfo struct {
		Index  uint64
		Offset int64
		Len    uint32
	}
	var batch []frameInfo

	_, err = readThroughSegment(rf, func(info types.SegmentInfo, fh frameHeader, offset int64) (bool, error) {
		if fh.typ == FrameCommit {
			// All the previous entries have been committed. Read them and send up to
			// caller.
			for _, frame := range batch {
				// Check the header is reasonable
				if frame.Len > MaxEntrySize {
					return false, fmt.Errorf("failed to read entry idx=%d, frame header length (%d) is too big: %w",
						frame.Index, frame.Len, err)
				}

				if frame.Len > uint32(len(buf)) {
					buf = make([]byte, frame.Len)
				}

				n, err := rf.ReadAt(buf[:frame.Len], frame.Offset+frameHeaderLen)
				if err != nil {
					return false, err
				}
				if uint32(n) < frame.Len {
					return false, io.ErrUnexpectedEOF
				}

				ok, err := fn(info, types.LogEntry{Index: frame.Index, Data: buf[:n]})
				if !ok || err != nil {
					return ok, err
				}
			}
			// Reset batch
			batch = batch[:0]
			return true, nil
		}

		if fh.typ != FrameEntry {
			return true, nil
		}

		if idx <= after {
			// Not in the range we care about, skip reading the entry.
			idx++
			return true, nil
		}
		if before > 0 && idx >= before {
			// We're done
			return false, nil
		}

		batch = append(batch, frameInfo{idx, offset, fh.len})
		idx++
		return true, nil
	})

	return err
}

// DumpLogs attempts to read all log entries from segment files in the directory
// for debugging purposes. It does _not_ use the metadata and so may output log
// entries that are uncommitted or already truncated as far as the writing
// process is concerned. As such it should not be used for replication of data.
// It is useful though to debug the contents of the log even while the writing
// application is still running. After and before if non-zero specify exclusive
// bounds on the logs that should be returned which may allow the implementation
// to skip reading entire segment files that are not in the range.
func (f *Filer) DumpLogs(after, before uint64, fn func(info types.SegmentInfo, e types.LogEntry) (bool, error)) error {
	baseIndexes, segIDsSorted, err := f.listInternal()
	if err != nil {
		return err
	}

	for i, id := range segIDsSorted {
		baseIndex := baseIndexes[id]
		nextBaseIndex := uint64(0)
		if i+1 < len(segIDsSorted) {
			// This is not the last segment, peek at the base index of that one and
			// assume that this segment won't contain indexes that high.
			nextBaseIndex = baseIndexes[segIDsSorted[i+1]]
		}
		// See if this file contains any indexes in the range
		if after > 0 && nextBaseIndex > 0 && after >= nextBaseIndex {
			// This segment is all indexes before the lower bound we care about
			continue
		}
		if before > 0 && before <= baseIndex {
			// This segment is all indexes higher than the upper bound. We've output
			// every log in the range at this point (barring edge cases where we race
			// with a truncation which leaves multiple generations of segment files on
			// disk which we are going to ignore for now).
			return nil
		}

		// We probably care about at least some of the entries in this segment
		err := f.DumpSegment(baseIndex, id, after, before, fn)
		if err != nil {
			return err
		}
	}

	return nil
}
