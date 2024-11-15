// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package wal

import (
	"sync/atomic"

	"github.com/benbjohnson/immutable"
	"github.com/hashicorp/raft-wal/types"
)

// state is an immutable snapshot of the state of the log. Modifications must be
// made by copying and modifying the copy. This is easy enough because segments
// is an immutable map so changing and re-assigning to the clone won't impact
// the original map, and tail is just a pointer that can be mutated in the
// shallow clone. Note that methods called on the tail segmentWriter may mutate
// it's state so must only be called while holding the WAL's writeLock.
type state struct {
	// refCount tracks readers that are reading segments based on this metadata.
	// It is accessed atomically nd must be 64 bit aligned (i.e. leave it at the
	// start of the struct).
	refCount int32
	// finaliser is set at most once while WAL is holding the write lock in order
	// to provide a func that must be called when all current readers are done
	// with this state. It's used for deferring closing and deleting old segments
	// until we can be sure no reads are still in progress on them.
	finalizer atomic.Value // func()

	nextSegmentID uint64

	// nextBaseIndex is used to signal which baseIndex to use next if there are no
	// segments or current tail.
	nextBaseIndex uint64
	segments      *immutable.SortedMap[uint64, segmentState]
	tail          types.SegmentWriter
}

type segmentState struct {
	types.SegmentInfo

	// r is the SegmentReader for our in-memory state.
	r types.SegmentReader
}

// Commit converts the in-memory state into a PersistentState.
func (s *state) Persistent() types.PersistentState {
	segs := make([]types.SegmentInfo, 0, s.segments.Len())
	it := s.segments.Iterator()
	for !it.Done() {
		_, s, _ := it.Next()
		segs = append(segs, s.SegmentInfo)
	}
	return types.PersistentState{
		NextSegmentID: s.nextSegmentID,
		Segments:      segs,
	}
}

func (s *state) getLog(index uint64) (*types.PooledBuffer, error) {
	// Check the tail writer first
	if s.tail != nil {
		raw, err := s.tail.GetLog(index)
		if err != nil && err != ErrNotFound {
			// Return actual errors since they might mask the fact that index really
			// is in the tail but failed to read for some other reason.
			return nil, err
		}
		if err == nil {
			// No error means we found it and just need to decode.
			return raw, nil
		}
		// Not in the tail segment, fall back to searching previous segments.
	}

	seg, err := s.findSegmentReader(index)
	if err != nil {
		return nil, err
	}

	return seg.GetLog(index)
}

// findSegmentReader searches the segment tree for the segment that contains the
// log at index idx. It may return the tail segment which may not in fact
// contain idx if idx is larger than the last written index. Typically this is
// called after already checking with the tail writer whether the log is in
// there which means the caller can be sure it's not going to return the tail
// segment.
func (s *state) findSegmentReader(idx uint64) (types.SegmentReader, error) {

	if s.segments.Len() == 0 {
		return nil, ErrNotFound
	}

	// Search for a segment with baseIndex.
	it := s.segments.Iterator()

	// The baseIndex we want is the first one lower or equal to idx. Seek gets us
	// to the first result equal or greater so we are either at it (if equal) or
	// on the one _after_ the one we need. We step back since that's most likely
	it.Seek(idx)
	// The first call to Next/Prev actually returns the node the iterator is
	// currently on (which is probably the one after the one we want) but in some
	// edge cases we might actually want this one. Rather than reversing back and
	// coming forward again, just check both this and the one before it.
	_, seg, ok := it.Prev()
	if ok && seg.BaseIndex > idx {
		_, seg, ok = it.Prev()
	}

	// We either have the right segment or it doesn't exist.
	if ok && seg.MinIndex <= idx && (seg.MaxIndex == 0 || seg.MaxIndex >= idx) {
		return seg.r, nil
	}

	return nil, ErrNotFound
}

func (s *state) getTailInfo() *segmentState {
	it := s.segments.Iterator()
	it.Last()
	_, tail, ok := it.Next()
	if !ok {
		return nil
	}
	return &tail
}

func (s *state) append(entries []types.LogEntry) error {
	return s.tail.Append(entries)
}

func (s *state) firstIndex() uint64 {
	it := s.segments.Iterator()
	_, seg, ok := it.Next()
	if !ok {
		return 0
	}
	if seg.SealTime.IsZero() {
		// First segment is unsealed so is also the tail. Check it actually has at
		// least one log in otherwise it doesn't matter what the BaseIndex/MinIndex
		// are.
		if s.tail.LastIndex() == 0 {
			// No logs in the WAL
			return 0
		}
		// At least one log exists, return the MinIndex
	}
	return seg.MinIndex
}

func (s *state) lastIndex() uint64 {
	tailIdx := s.tail.LastIndex()
	if tailIdx > 0 {
		return tailIdx
	}
	// Current tail is empty. Check there are previous sealed segments.
	it := s.segments.Iterator()
	it.Last()
	_, _, ok := it.Prev()
	if !ok {
		// No tail! shouldn't be possible but means no logs yet
		return 0
	}
	// Go back to the segment before the tail
	_, _, ok = it.Prev()
	if !ok {
		// No previous segment so the whole log is empty
		return 0
	}

	// There was a previous segment so it's MaxIndex will be one less than the
	// tail's BaseIndex.
	tailSeg := s.getTailInfo()
	if tailSeg == nil || tailSeg.BaseIndex == 0 {
		return 0
	}
	return tailSeg.BaseIndex - 1
}

func (s *state) acquire() func() {
	atomic.AddInt32(&s.refCount, 1)
	return s.release
}

func (s *state) release() {
	// decrement on release
	new := atomic.AddInt32(&s.refCount, -1)
	if new == 0 {
		// Cleanup state associated with this version now all refs have gone. Since
		// there are no more refs and we should not set a finalizer until this state
		// is no longer the active state, we can be sure this will happen only one.
		// Even still lets swap the fn to ensure we only call finalizer once ever!
		// We can't swap actual nil as it's not the same type as func() so do a
		// dance with a nilFn below.
		var nilFn func()
		fnRaw := s.finalizer.Swap(nilFn)
		if fn, ok := fnRaw.(func()); ok && fn != nil {
			fn()
		}
	}
}

// clone returns a new state which is a shallow copy of just the immutable parts
// of s. This is safer than a simple assignment copy because that "reads" the
// atomically modified state non-atomically. We never want to copy the refCount
// or finalizer anyway.
func (s *state) clone() state {
	return state{
		nextSegmentID: s.nextSegmentID,
		segments:      s.segments,
		tail:          s.tail,
	}
}
