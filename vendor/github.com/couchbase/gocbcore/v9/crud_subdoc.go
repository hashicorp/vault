package gocbcore

import (
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
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

	// Internal: This should never be used and is not supported.
	User []byte

	// Volatile: Tracer API is subject to change.
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

	// Internal: This should never be used and is not supported.
	User []byte

	// Volatile: Tracer API is subject to change.
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
		IsDeleted bool
	}
}

// MutateInResult encapsulates the result of a MutateInEx operation.
type MutateInResult struct {
	Cas           Cas
	MutationToken MutationToken
	Ops           []SubDocResult
}
