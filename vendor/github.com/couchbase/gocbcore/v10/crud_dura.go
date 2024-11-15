package gocbcore

import (
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

// ObserveOptions encapsulates the parameters for a ObserveEx operation.
type ObserveOptions struct {
	Key            []byte
	ReplicaIdx     int
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// ObserveVbOptions encapsulates the parameters for a ObserveVbEx operation.
type ObserveVbOptions struct {
	VbID          uint16
	VbUUID        VbUUID
	ReplicaIdx    int
	RetryStrategy RetryStrategy
	Deadline      time.Time

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// ObserveResult encapsulates the result of a ObserveEx operation.
type ObserveResult struct {
	KeyState memd.KeyState
	Cas      Cas

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}

// ObserveVbResult encapsulates the result of a ObserveVbEx operation.
type ObserveVbResult struct {
	DidFailover  bool
	VbID         uint16
	VbUUID       VbUUID
	PersistSeqNo SeqNo
	CurrentSeqNo SeqNo
	OldVbUUID    VbUUID
	LastSeqNo    SeqNo

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}
