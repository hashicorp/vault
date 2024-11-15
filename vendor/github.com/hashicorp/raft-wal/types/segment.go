// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package types

import (
	"io"
	"time"
)

// SegmentInfo is the metadata describing a single WAL segment.
type SegmentInfo struct {
	// ID uniquely identifies this segment file
	ID uint64

	// BaseIndex is the raft index of the first entry that will be written to the
	// segment.
	BaseIndex uint64

	// MinIndex is the logical lowest index that still exists in the segment. It
	// may be greater than BaseIndex if a head truncation has "deleted" a prefix
	// of the segment.
	MinIndex uint64

	// MaxIndex is the logical highest index that still exists in the segment. It
	// may be lower than the actual highest index if a tail truncation has
	// "deleted" a suffix of the segment. It is zero for unsealed segments and
	// only set one seal.
	MaxIndex uint64

	// Codec identifies the codec used to encode log entries. Codec values 0 to
	// 16k (i.e. the lower 16 bits) are reserved for internal future usage. Custom
	// codecs must be registered with an identifier higher than this which the
	// caller is responsible for ensuring uniquely identifies the specific version
	// of their codec used in any given log. uint64 provides sufficient space that
	// a randomly generated identifier is almost certainly unique.
	Codec uint64

	// IndexStart is the file offset where the index can be read from it's 0 for
	// tail segments and only set after a segment is sealed.
	IndexStart uint64

	// CreateTime records when the segment was first created.
	CreateTime time.Time

	// SealTime records when the segment was sealed. Zero indicates that it's not
	// sealed yet.
	SealTime time.Time

	// SizeLimit is the soft limit for the segment's size. The segment file may be
	// pre-allocated to this size on filesystems that support it. It is a soft
	// limit in the sense that the final Append usually takes the segment file
	// past this size before it is considered full and sealed.
	SizeLimit uint32
}

// SegmentFiler is the interface that provides access to segments to the WAL. It
// encapsulated creating, and recovering segments and returning reader or writer
// interfaces to interact with them. It's main purpose is to abstract the core
// WAL logic both from the actual encoding layer of segment files. You can think
// of it as a layer of abstraction above the VFS which abstracts actual file
// system operations on files but knows nothing about the format. In tests for
// example we can implement a SegmentFiler that is way simpler than the real
// encoding/decoding layer on top of a VFS - even an in-memory VFS which makes
// tests much simpler to write and run.
type SegmentFiler interface {
	// Create adds a new segment with the given info and returns a writer or an
	// error.
	Create(info SegmentInfo) (SegmentWriter, error)

	// RecoverTail is called on an unsealed segment when re-opening the WAL it
	// will attempt to recover from a possible crash. It will either return an
	// error, or return a valid segmentWriter that is ready for further appends.
	// If the expected tail segment doesn't exist it must return an error wrapping
	// os.ErrNotExist.
	RecoverTail(info SegmentInfo) (SegmentWriter, error)

	// Open an already sealed segment for reading. Open may validate the file's
	// header and return an error if it doesn't match the expected info.
	Open(info SegmentInfo) (SegmentReader, error)

	// List returns the set of segment IDs currently stored. It's used by the WAL
	// on recovery to find any segment files that need to be deleted following a
	// unclean shutdown. The returned map is a map of ID -> BaseIndex. BaseIndex
	// is returned to allow subsequent Delete calls to be made.
	List() (map[uint64]uint64, error)

	// Delete removes the segment with given baseIndex and id if it exists. Note
	// that baseIndex is technically redundant since ID is unique on it's own. But
	// in practice we name files (or keys) with both so that they sort correctly.
	// This interface allows a  simpler implementation where we can just delete
	// the file if it exists without having to scan the underlying storage for a.
	Delete(baseIndex, ID uint64) error
}

// SegmentWriter manages appending logs to the tail segment of the WAL. It's an
// interface to make testing core WAL simpler. Every SegmentWriter will have
// either `init` or `recover` called once before any other methods. When either
// returns it must either return an error or be ready to accept new writes and
// reads.
type SegmentWriter interface {
	io.Closer
	SegmentReader

	// Append adds one or more entries. It must not return until the entries are
	// durably stored otherwise raft's guarantees will be compromised. Append must
	// not be called concurrently with any other call to Sealed, Append or
	// ForceSeal.
	Append(entries []LogEntry) error

	// Sealed returns whether the segment is sealed or not. If it is it returns
	// true and the file offset that it's index array starts at to be saved in
	// meta data. WAL will call this after every append so it should be relatively
	// cheap in the common case. This design allows the final Append to write out
	// the index or any additional data needed at seal time in the same fsync.
	// Sealed must not be called concurrently with any other call to Sealed,
	// Append or ForceSeal.
	Sealed() (bool, uint64, error)

	// ForceSeal causes the segment to become sealed by writing out an index
	// block. This is not used in the typical flow of append and rotation, but is
	// necessary during truncations where some suffix of the writer needs to be
	// truncated. Rather than manipulate what is on disk in a complex way, the WAL
	// will simply force seal it with whatever state it has already saved and then
	// open a new segment at the right offset for continued writing. ForceSeal may
	// be called on a segment that has already been sealed and should just return
	// the existing index offset in that case. (We don't actually rely on that
	// currently but it's easier not to assume we'll always call it at most once).
	// ForceSeal must not be called concurrently with any other call to Sealed,
	// Append or ForceSeal.
	ForceSeal() (uint64, error)

	// LastIndex returns the most recently persisted index in the log. It must
	// respond without blocking on Append since it's needed frequently by read
	// paths that may call it concurrently. Typically this will be loaded from an
	// atomic int. If the segment is empty lastIndex should return zero.
	LastIndex() uint64
}

// SegmentReader wraps a ReadableFile to allow lookup of logs in an existing
// segment file. It's an interface to make testing core WAL simpler. The first
// call will always be validate which passes in the ReaderAt to be used for
// subsequent reads.
type SegmentReader interface {
	io.Closer

	// GetLog returns the raw log entry bytes associated with idx. If the log
	// doesn't exist in this segment ErrNotFound must be returned.
	GetLog(idx uint64) (*PooledBuffer, error)
}
