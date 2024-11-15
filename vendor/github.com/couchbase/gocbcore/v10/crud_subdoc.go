package gocbcore

import (
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

// LookupInOptions encapsulates the parameters for a LookupInEx operation.
type LookupInOptions struct {
	Key            []byte
	Flags          memd.SubdocDocFlag
	Ops            []SubDocOp
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time
	ReplicaIdx     int

	// Uncommitted: This API may change in the future.
	ServerGroup string

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// MutateInOptions encapsulates the parameters for a MutateInEx operation.
type MutateInOptions struct {
	Key                    []byte
	Flags                  memd.SubdocDocFlag
	Cas                    Cas
	Expiry                 uint32
	Ops                    []SubDocOp
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time
	PreserveExpiry         bool

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// SubDocResult encapsulates the results from a single sub-document operation.
type SubDocResult struct {
	Err   error
	Value []byte
}

// LookupInResult encapsulates the result of a LookupInEx operation.
type LookupInResult struct {
	Cas Cas
	Ops []SubDocResult

	// Internal: This should never be used and is not supported.
	Internal struct {
		IsDeleted     bool
		ResourceUnits *ResourceUnitResult
	}
}

// MutateInResult encapsulates the result of a MutateInEx operation.
type MutateInResult struct {
	Cas           Cas
	MutationToken MutationToken
	Ops           []SubDocResult

	// Internal: This should never be used and is not supported.
	Internal struct {
		ResourceUnits *ResourceUnitResult
	}
}
