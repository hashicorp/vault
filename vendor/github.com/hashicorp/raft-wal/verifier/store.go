// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package verifier

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-wal/metrics"
)

var _ raft.LogStore = &LogStore{}
var _ raft.MonotonicLogStore = &LogStore{}

// LogStore is a raft.LogStore that acts as middleware around an underlying
// persistent store. It provides support for periodically verifying that ranges
// of logs read back from the LogStore match the values written, and the values
// read from the LogStores of other peers even though all peers will have
// different actual log ranges due to independent snapshotting and truncation.
//
// Verification of the underlying log implementation may be performed as
// follows:
//  1. The application provides an implementation of `IsCheckpoint` that is
//     able to identify whether the encoded data represents a checkpoint
//     command.
//  2. The application's raft leader then may periodically append such a
//     checkpoint log to be replicated out.
//  3. When the LogStore has a log appended for which IsCheckpoint returns true,
//     it will write the current cumulative checksum over log entries since the
//     last checkpoint into the Extra field. Since hashicorp/raft only
//     replicates to peers _after_ a trip through the LogStore, this checksum
//     will be replicated.
//  4. When a follower has a log appended for which IsCheckpoint returns true,
//     but already has non-empty Extra metadata, it will trigger a background
//     verification.
//  5. Verification happens in the background and reads all logs from the
//     underlying store since the last checkpoint, calculating their checksums
//     cumulatively before calling the configured Report func with a summary of
//     what it found.
type LogStore struct {
	checksum    uint64 // accessed atomically
	sumStartIdx uint64 // accessed atomically

	s raft.LogStore

	metrics metrics.Collector
	log     hclog.Logger

	verifyCh chan VerificationReport

	checkpointFn IsCheckpointFn
	reportFn     ReportFn
}

// NewLogStore creates a verifying LogStore. CheckpointFn and ReportFn must be
// set on the returned store _before_ it is passed to Raft, or may be left as
// nil to bypass verification. Close must be called when the log store is no
// longer useful to cleanup background verification.
func NewLogStore(store raft.LogStore, checkpointFn IsCheckpointFn, reportFn ReportFn, mc metrics.Collector) *LogStore {
	c := &LogStore{
		s:            store,
		metrics:      mc,
		verifyCh:     make(chan VerificationReport, 1),
		checkpointFn: checkpointFn,
		reportFn:     reportFn,
	}
	go c.runVerifier()
	return c
}

// FirstIndex returns the first index written. 0 for no entries.
func (s *LogStore) FirstIndex() (uint64, error) {
	return s.s.FirstIndex()
}

// LastIndex returns the last index written. 0 for no entries.
func (s *LogStore) LastIndex() (uint64, error) {
	return s.s.LastIndex()
}

// GetLog gets a log entry at a given index.
func (s *LogStore) GetLog(index uint64, log *raft.Log) error {
	return s.s.GetLog(index, log)
}

// StoreLog stores a log entry.
func (s *LogStore) StoreLog(log *raft.Log) error {
	return s.StoreLogs([]*raft.Log{log})
}

func encodeCheckpointMeta(startIdx, sum uint64) []byte {
	var buf [24]byte
	binary.LittleEndian.PutUint64(buf[0:8], ExtensionMagicPrefix)
	binary.LittleEndian.PutUint64(buf[8:16], startIdx)
	binary.LittleEndian.PutUint64(buf[16:24], sum)
	return buf[:]
}

func decodeCheckpointMeta(bs []byte) (startIdx, sum uint64, err error) {
	if len(bs) < 24 {
		return 0, 0, io.ErrShortBuffer
	}
	magic := binary.LittleEndian.Uint64(bs[0:8])
	if magic != ExtensionMagicPrefix {
		return 0, 0, errors.New("invalid extension data")
	}
	startIdx = binary.LittleEndian.Uint64(bs[8:16])
	sum = binary.LittleEndian.Uint64(bs[16:24])
	return startIdx, sum, nil
}

func (s *LogStore) updateVerifyState(log *raft.Log, checksum, startIdx uint64) (newSum, newStartIdx uint64, r *VerificationReport, err error) {
	// Check if log is a checkpoint, note we already nil-checked this function
	// before calling.
	isCP, err := s.checkpointFn(log)
	if err != nil {
		return 0, 0, nil, err
	}

	if startIdx == 0 {
		startIdx = log.Index
	}

	if isCP {
		r = &VerificationReport{
			Range:      LogRange{End: log.Index},
			WrittenSum: checksum,
		}
		if len(log.Extensions) == 0 {
			// It's a new checkpoint and we must be the leader. Set our state.
			log.Extensions = encodeCheckpointMeta(startIdx, checksum)
			r.Range.Start = startIdx
			r.ExpectedSum = checksum
		} else {
			cpStartIdx, cpSum, err := decodeCheckpointMeta(log.Extensions)
			if err != nil {
				return 0, 0, nil, err
			}
			r.Range.Start = cpStartIdx
			r.ExpectedSum = cpSum

			// If we've calculated our own checksum over a different range to the
			// leader e.g. because we just started and this is the first sum then
			// there's no point trying to verify so leave WrittenSum zero.
			if cpStartIdx != startIdx {
				r.WrittenSum = 0
			}
		}
		// Reset the checksum as we're now in the range of the next checkpoint. We
		// don't update the store state yet until we know these logs committed to
		// the underlying store.
		checksum = 0
		startIdx = log.Index
	}

	// Whether checkpoint or not, hash the entry and update return updated
	// checksum.
	checksum = checksumLog(checksum, log)
	return checksum, startIdx, r, nil
}

// StoreLogs stores multiple log entries.
func (s *LogStore) StoreLogs(logs []*raft.Log) error {
	if len(logs) < 1 {
		return nil
	}

	// Maintain a local copy of the checksum and sumStartIdx, we'll update the
	// state only once we know all these entries were stored.
	cs := atomic.LoadUint64(&s.checksum)
	startIdx := atomic.LoadUint64(&s.sumStartIdx)
	var triggeredReports []VerificationReport

	if s.checkpointFn != nil {
		var vr *VerificationReport
		var err error
		for _, log := range logs {
			cs, startIdx, vr, err = s.updateVerifyState(log, cs, startIdx)
			if err != nil {
				return fmt.Errorf("failed updating verifier state: %w", err)
			}
			if vr != nil {
				// We need to trigger a new checkpoint verification. But we can't until
				// after the logs are persisted below.
				triggeredReports = append(triggeredReports, *vr)
			}
		}
	}

	err := s.s.StoreLogs(logs)
	if err != nil {
		return err
	}

	// Update the checksum state now logs are committed.
	atomic.StoreUint64(&s.checksum, cs)
	atomic.StoreUint64(&s.sumStartIdx, startIdx)
	if len(triggeredReports) > 0 {
		s.metrics.IncrementCounter("checkpoints_written", uint64(len(triggeredReports)))
	}

	for _, r := range triggeredReports {
		s.triggerVerify(r)
	}
	return nil
}

// triggerVerify triggers a verification in the background. We won't block if
// the verifier is busy. The chan is one buffered so there can be at most one
// running and one waiting. If there is already one waiting so the chan is
// blocked, we drop r.
func (s *LogStore) triggerVerify(r VerificationReport) {
	select {
	case s.verifyCh <- r:
	default:
		s.metrics.IncrementCounter("dropped_reports", 1)
	}
}

// DeleteRange deletes a range of log entries. The range is inclusive.
func (s *LogStore) DeleteRange(min uint64, max uint64) error {
	return s.s.DeleteRange(min, max)
}

// Close cleans up the background verification routine and calls Close on the
// underlying store if it is an io.Closer.
func (s *LogStore) Close() error {
	if s.verifyCh == nil {
		return nil
	}
	close(s.verifyCh)
	// Don't set verifyCh to nil as that's racey - it's being accessed from other
	// routines.
	if closer, ok := s.s.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// IsMonotonic implements the raft.MonotonicLogStore interface. This is a shim
// to expose the underlying store as monotonically indexed or not.
func (s *LogStore) IsMonotonic() bool {
	if store, ok := s.s.(raft.MonotonicLogStore); ok {
		return store.IsMonotonic()
	}
	return false
}
