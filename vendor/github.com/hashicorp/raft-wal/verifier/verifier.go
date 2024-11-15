// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package verifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/raft"
	"github.com/segmentio/fasthash/fnv1a"
)

const (
	// ExtensionMagicPrefix is the prefix we append to log Extensions fields to
	// disambiguate from other middleware that may use extensions. This value is
	// carefully constructed to be completely invalid as the beginning of a
	// protobuf (3) wire protocol message since the other known user of this field
	// encodes its data that way. If the first byte were 0xa8 this would be a
	// valid protobuf field encoding for an int field, however currently the 3
	// least significant bits encode the field type as 7, which is not a valid
	// type in the current spec. Even if this does change in the future, the
	// field's tag number encoded here is 123456789 so it's extremely unlikely
	// that any valid protobuf schema will ever have enough fields or arbitrarily
	// decide to assign field tags that large (though unrecognized tags would be
	// ignored on decode). Finally, the value of the field is the varint encoding
	// of the randomly chosen value 53906 so if type 7 is ever valid in the future
	// and used as a length-prefixed type, the length decoded would be way longer
	// than the buffer making it invalid.
	ExtensionMagicPrefix uint64 = 0xafd1f9d60392a503
)

// IsCheckpointFn is a function that can decide whether the contents of a raft
// log's Data represents a checkpoint message. It is called on every append so
// it must be relatively fast in the common case. If it returns true for a log,
// the log's Extra field will be used to encode verification metadata and must
// be empty - if it's not empty the append will fail and force the leader to
// step down. If an error is returned the same will happen.
type IsCheckpointFn func(*raft.Log) (bool, error)

// ReportFn is a function that will be called after every checkpoint has been
// verified. It will not be called concurrently. The VerificationReport may
// represent a failure to report so it's Err field should be checked. For
// example, if checkpoints are arriving faster than they can be calculated, some
// will be skipped and no report will be made for that range. The next report
// that is delivered will contain the range missed for logging. Note that
// ReportFn is called synchronously by the verifier so it should not block for
// long otherwise it may cause the verifier to miss later checkpoints.
type ReportFn func(VerificationReport)

// ErrRangeMismatch is the error type returned in a VerificationReport where the
// follower does not have enough logs on disk to fill the checkpoint's range and
// so is bound to fail. This is a separate type from pure failures to read a log
// because it's expected this could happen just after truncations or if the
// interval is to large for the number of logs retained etc. Implementations may
// choose to detect this and report as a warning rather than a failure as it
// indicates only an inability to report correctly not an actual error in
// processing data.
var ErrRangeMismatch = errors.New("range mismatch")

// ErrChecksumMismatch is the error type returned in a VerificationReport where
// the log range's checksum didn't match.
type ErrChecksumMismatch string

// Error implements error
func (e ErrChecksumMismatch) Error() string {
	return string(e)
}

// LogRange describes the set of logs in the range [Start, End). That is End is
// NOT inclusive.
type LogRange struct {
	Start uint64
	End   uint64
}

// String implements Stringer
func (r LogRange) String() string {
	return fmt.Sprintf("[%d, %d)", r.Start, r.End)
}

// VerificationReport describes the result of attempting to verify the contents
// of all logs in a range compared with the input the leader delivered for that
// same range.
type VerificationReport struct {
	// Range is the range of raft indexes over which the leader calculated its
	// checksum. In steady state it typically starts with the index of the
	// previous checkpoint command, but after an election it could be an arbitrary
	// point in the log. If the range is no longer in the server's log (due to not
	// seeing one yet or it being truncated too soon) this will be reported as an
	// Err - a longer log retention (`raft.Config.TrailingLogs`) or shorter
	// interval between checkpoints should be chosen if this happens often.
	Range LogRange

	// ExpectedSum is a uint64 checksum over the logs in the range as calculated
	// by the leader before appending to disk.
	ExpectedSum uint64

	// WrittenSum is the uint64 checksum calculated over the logs in the range of
	// a follower as it wrote them to it's own LogStore. It might be zero to
	// indicate that the follower has not written all the logs in Range since
	// startup and so its written sum will be invalid. Risk of collision with
	// genuine zero sum is acceptable. If zero the verifier will have ignored it
	// and not raised an error if it didn't match expected.
	WrittenSum uint64

	// ReadSum is the uint64 checksum calculated over the logs in the range as
	// read from the underlying LogStore in the range [StartIndex, EndIndex).
	ReadSum uint64

	// Err indicates any error that prevented the report from being completed or
	// the result of the report. It will be set to ErrChecksumMismatch if the
	// report was conducted correctly, but the log data written or read checksum
	// did not match the leader's write checksum. The message in the error
	// describes the nature of the failure.
	Err error

	// SkippedRange indicates the ranges of logs covered by any checkpoints that
	// we skipped due to spending too much time verifying. If this is regularly
	// non-nil it likely indicates that the checkpoint frequency is too fast.
	SkippedRange *LogRange

	// Elapsed records how long it took to read the range and generate the report.
	Elapsed time.Duration
}

func (s *LogStore) runVerifier() {
	if s.reportFn == nil {
		// Nothing to do!
		return
	}

	var lastCheckPointIdx uint64
	for {
		report, ok := <-s.verifyCh
		if !ok {
			// Close was called
			return
		}

		// Detect skipped checkpoints
		if lastCheckPointIdx > 0 && lastCheckPointIdx != report.Range.Start {
			report.SkippedRange = &LogRange{
				Start: lastCheckPointIdx,
				End:   report.Range.Start,
			}
		}
		lastCheckPointIdx = report.Range.End

		st := time.Now()
		s.verify(&report)

		// Whatever state report ended up in, deliver it!
		report.Elapsed = time.Since(st)
		s.reportFn(report)
		s.metrics.IncrementCounter("ranges_verified", 1)
	}
}

func (s *LogStore) verify(report *VerificationReport) {
	// Attempt to read all the logs in the range from underlying store.
	var log raft.Log

	// If this is a follower but it _wrote_ different data to it's log than the
	// leader in this range then there's not much point verifying that we read it
	// back OK.
	if report.WrittenSum != 0 && report.WrittenSum != report.ExpectedSum {
		s.metrics.IncrementCounter("write_checksum_failures", 1)
		report.Err = ErrChecksumMismatch(fmt.Sprintf("log verification failed for range %s: "+
			"in-flight corruption: follower wrote checksum=%08x, leader wrote checksum=%08x",
			report.Range, report.WrittenSum, report.ExpectedSum))
		return
	}

	// Do we actually have enough logs to calculate the checksum? If not indicate
	// that explicitly as its an expected case rather than a real "error". Note
	// that we may get a racey false negative here if truncation happens right
	// between this check and the GetLog call below but there's not much we can do
	// about that and hopefully is rare enough!
	first, err := s.s.FirstIndex()
	if err != nil {
		report.Err = fmt.Errorf("unable to verify log range %s: %w", report.Range, err)
		return
	}
	if first > report.Range.Start {
		// We don't have enough logs to calculate this correctly.
		report.Err = ErrRangeMismatch
		return
	}

	sum := uint64(0)
	for idx := report.Range.Start; idx < report.Range.End; idx++ {
		err := s.s.GetLog(idx, &log)
		if err != nil {
			report.Err = fmt.Errorf("unable to verify log range %s: %w", report.Range, err)
			return
		}
		sum = checksumLog(sum, &log)
	}
	report.ReadSum = sum

	if report.ReadSum != report.ExpectedSum {
		s.metrics.IncrementCounter("read_checksum_failures", 1)
		report.Err = ErrChecksumMismatch(fmt.Sprintf("log verification failed for range %s: "+
			"storage corruption: node read checksum=%08x, leader wrote checksum=%08x",
			report.Range, report.ReadSum, report.ExpectedSum))
		return
	}
}

func checksumLog(sum uint64, log *raft.Log) uint64 {
	// Special case for bootstrap config entries (index 1, type configuration)
	// since these are not replicated by raft and so may not be byte-for-byte
	// identical as long as they are logical the same on all peers. So just treat
	// them all as identical to avoid false-positives on startup.
	if log.Index == 1 && log.Type == raft.LogConfiguration {
		return 0
	}
	sum = fnv1a.AddUint64(sum, log.Index)
	sum = fnv1a.AddUint64(sum, log.Term)
	sum = fnv1a.AddUint64(sum, uint64(log.Type))
	sum = fnv1a.AddBytes64(sum, log.Data)
	if len(log.Extensions) > 0 {
		sum = fnv1a.AddBytes64(sum, log.Extensions)
	}
	return sum
}
