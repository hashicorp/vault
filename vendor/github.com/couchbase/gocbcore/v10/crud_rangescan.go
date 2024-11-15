package gocbcore

import (
	"encoding/base64"
	"strconv"
	"time"
)

// RangeScanCreateOptions encapsulates the parameters for a RangeScanCreate operation.
type RangeScanCreateOptions struct {
	// Deadline will also be sent as a part of the payload if Snapshot is not nil.
	Deadline time.Time

	CollectionName string
	ScopeName      string

	CollectionID uint32
	// Note: if set then KeysOnly on RangeScanContinueOptions *must* also be set.
	KeysOnly bool
	Range    *RangeScanCreateRangeScanConfig
	Sampling *RangeScanCreateRandomSamplingConfig
	Snapshot *RangeScanCreateSnapshotRequirements

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

func (opts RangeScanCreateOptions) toRequest() (*rangeScanCreateRequest, error) {
	if opts.Range != nil && opts.Sampling != nil {
		return nil, wrapError(errInvalidArgument, "only one of range and sampling can be set")
	}
	if opts.Range == nil && opts.Sampling == nil {
		return nil, wrapError(errInvalidArgument, "one of range and sampling must set")
	}

	var collection string
	if opts.CollectionID != 0 {
		collection = strconv.FormatUint(uint64(opts.CollectionID), 16)
	}
	createReq := &rangeScanCreateRequest{
		Collection: collection,
		KeyOnly:    opts.KeysOnly,
	}

	if opts.Range != nil {
		if opts.Range.hasStart() && opts.Range.hasExclusiveStart() {
			return nil, wrapError(errInvalidArgument, "only one of start and exclusive start within range can be set")
		}
		if opts.Range.hasEnd() && opts.Range.hasExclusiveEnd() {
			return nil, wrapError(errInvalidArgument, "only one of end and exclusive end within range can be set")
		}
		if !(opts.Range.hasStart() || opts.Range.hasExclusiveStart()) {
			return nil, wrapError(errInvalidArgument, "one of start and exclusive start within range must both be set")
		}
		if !(opts.Range.hasEnd() || opts.Range.hasExclusiveEnd()) {
			return nil, wrapError(errInvalidArgument, "one of end and exclusive end within range must both be set")
		}

		createReq.Range = &rangeScanCreateRange{}
		if len(opts.Range.Start) > 0 {
			createReq.Range.Start = base64.StdEncoding.EncodeToString(opts.Range.Start)
		}
		if len(opts.Range.End) > 0 {
			createReq.Range.End = base64.StdEncoding.EncodeToString(opts.Range.End)
		}
		if len(opts.Range.ExclusiveStart) > 0 {
			createReq.Range.ExclusiveStart = base64.StdEncoding.EncodeToString(opts.Range.ExclusiveStart)
		}
		if len(opts.Range.ExclusiveEnd) > 0 {
			createReq.Range.ExclusiveEnd = base64.StdEncoding.EncodeToString(opts.Range.ExclusiveEnd)
		}
	}

	if opts.Sampling != nil {
		if opts.Sampling.Samples == 0 {
			return nil, wrapError(errInvalidArgument, "samples within sampling must be set")
		}

		createReq.Sampling = &rangeScanCreateSample{
			Seed:    opts.Sampling.Seed,
			Samples: opts.Sampling.Samples,
		}
	}

	if opts.Snapshot != nil {
		if opts.Snapshot.VbUUID == 0 {
			return nil, wrapError(errInvalidArgument, "vbuuid within snapshot must be set")
		}
		if opts.Snapshot.SeqNo == 0 {
			return nil, wrapError(errInvalidArgument, "seqno within snapshot must be set")
		}

		createReq.Snapshot = &rangeScanCreateSnapshot{
			VbUUID:      strconv.FormatUint(uint64(opts.Snapshot.VbUUID), 10),
			SeqNo:       uint64(opts.Snapshot.SeqNo),
			SeqNoExists: opts.Snapshot.SeqNoExists,
		}
		createReq.Snapshot.Timeout = uint64(time.Until(opts.Deadline).Milliseconds())
	}

	return createReq, nil
}

// RangeScanCreateRangeScanConfig is the configuration available for performing a range scan.
type RangeScanCreateRangeScanConfig struct {
	Start          []byte
	End            []byte
	ExclusiveStart []byte
	ExclusiveEnd   []byte
}

func (cfg *RangeScanCreateRangeScanConfig) hasStart() bool {
	return len(cfg.Start) > 0
}

func (cfg *RangeScanCreateRangeScanConfig) hasEnd() bool {
	return len(cfg.End) > 0
}
func (cfg *RangeScanCreateRangeScanConfig) hasExclusiveStart() bool {
	return len(cfg.ExclusiveStart) > 0
}

func (cfg *RangeScanCreateRangeScanConfig) hasExclusiveEnd() bool {
	return len(cfg.ExclusiveEnd) > 0
}

// RangeScanCreateRandomSamplingConfig is the configuration available for performing a random sampling.
type RangeScanCreateRandomSamplingConfig struct {
	Seed    uint64
	Samples uint64
}

// RangeScanCreateSnapshotRequirements is the set of requirements that the vbucket snapshot must meet in-order for
// the request to be successful.
type RangeScanCreateSnapshotRequirements struct {
	VbUUID      VbUUID
	SeqNo       SeqNo
	SeqNoExists bool
}

// RangeScanCreateResult encapsulates the result of a RangeScanCreate operation.
type RangeScanCreateResult interface {
	ScanUUID() []byte
	KeysOnly() bool

	RangeScanContinue(opts RangeScanContinueOptions, dataCb RangeScanContinueDataCallback,
		actionCb RangeScanContinueActionCallback) (PendingOp, error)
	RangeScanCancel(opts RangeScanCancelOptions, cb RangeScanCancelCallback) (PendingOp, error)
}

type rangeScanCreateResult struct {
	scanUUID []byte
	keysOnly bool

	vbID   uint16
	connID string

	parent *crudComponent
}

func (createRes *rangeScanCreateResult) ScanUUID() []byte {
	return createRes.scanUUID
}

func (createRes *rangeScanCreateResult) KeysOnly() bool {
	return createRes.keysOnly
}

// RangeScanContinueOptions encapsulates the parameters for a RangeScanContinue operation.
type RangeScanContinueOptions struct {
	// Deadline will also be sent as a part of the payload if not zero.
	Deadline time.Time

	MaxCount uint32
	MaxBytes uint32

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// RangeScanItem encapsulates an iterm returned during a range scan.
type RangeScanItem struct {
	Value    []byte
	Key      []byte
	Flags    uint32
	Cas      Cas
	Expiry   uint32
	SeqNo    SeqNo
	Datatype uint8
}

// RangeScanContinueResult encapsulates the result of a RangeScanContinue operation.
type RangeScanContinueResult struct {
	More     bool
	Complete bool
}

// RangeScanCancelOptions encapsulates the parameters for a RangeScanCancel operation.
type RangeScanCancelOptions struct {
	Deadline time.Time

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// RangeScanCancelResult encapsulates the result of a RangeScanCancel operation.
type RangeScanCancelResult struct{}
