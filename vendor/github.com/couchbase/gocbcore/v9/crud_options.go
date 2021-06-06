package gocbcore

import (
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

// GetOptions encapsulates the parameters for a GetEx operation.
type GetOptions struct {
	Key            []byte
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetAndTouchOptions encapsulates the parameters for a GetAndTouchEx operation.
type GetAndTouchOptions struct {
	Key            []byte
	Expiry         uint32
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetAndLockOptions encapsulates the parameters for a GetAndLockEx operation.
type GetAndLockOptions struct {
	Key            []byte
	LockTime       uint32
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetAnyReplicaOptions encapsulates the parameters for a GetAnyReplicaEx operation.
type GetAnyReplicaOptions struct {
	Key            []byte
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetOneReplicaOptions encapsulates the parameters for a GetOneReplicaEx operation.
type GetOneReplicaOptions struct {
	Key            []byte
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	ReplicaIdx     int
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// TouchOptions encapsulates the parameters for a TouchEx operation.
type TouchOptions struct {
	Key            []byte
	Expiry         uint32
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// UnlockOptions encapsulates the parameters for a UnlockEx operation.
type UnlockOptions struct {
	Key            []byte
	Cas            Cas
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// DeleteOptions encapsulates the parameters for a DeleteEx operation.
type DeleteOptions struct {
	Key                    []byte
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Cas                    Cas
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// AddOptions encapsulates the parameters for a AddEx operation.
type AddOptions struct {
	Key                    []byte
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Value                  []byte
	Flags                  uint32
	Datatype               uint8
	Expiry                 uint32
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

type storeOptions struct {
	Key                    []byte
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Value                  []byte
	Flags                  uint32
	Datatype               uint8
	Cas                    Cas
	Expiry                 uint32
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// SetOptions encapsulates the parameters for a SetEx operation.
type SetOptions struct {
	Key                    []byte
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Value                  []byte
	Flags                  uint32
	Datatype               uint8
	Expiry                 uint32
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// ReplaceOptions encapsulates the parameters for a ReplaceEx operation.
type ReplaceOptions struct {
	Key                    []byte
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Value                  []byte
	Flags                  uint32
	Datatype               uint8
	Cas                    Cas
	Expiry                 uint32
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// AdjoinOptions encapsulates the parameters for a AppendEx or PrependEx operation.
type AdjoinOptions struct {
	Key                    []byte
	Value                  []byte
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Cas                    Cas
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// CounterOptions encapsulates the parameters for a IncrementEx or DecrementEx operation.
type CounterOptions struct {
	Key                    []byte
	Delta                  uint64
	Initial                uint64
	Expiry                 uint32
	CollectionName         string
	ScopeName              string
	RetryStrategy          RetryStrategy
	Cas                    Cas
	DurabilityLevel        memd.DurabilityLevel
	DurabilityLevelTimeout time.Duration
	CollectionID           uint32
	Deadline               time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetRandomOptions encapsulates the parameters for a GetRandomEx operation.
type GetRandomOptions struct {
	RetryStrategy RetryStrategy
	Deadline      time.Time

	CollectionName string
	ScopeName      string
	CollectionID   uint32

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// GetMetaOptions encapsulates the parameters for a GetMetaEx operation.
type GetMetaOptions struct {
	Key            []byte
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// SetMetaOptions encapsulates the parameters for a SetMetaEx operation.
type SetMetaOptions struct {
	Key            []byte
	Value          []byte
	Extra          []byte
	Datatype       uint8
	Options        uint32
	Flags          uint32
	Expiry         uint32
	Cas            Cas
	RevNo          uint64
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}

// DeleteMetaOptions encapsulates the parameters for a DeleteMetaEx operation.
type DeleteMetaOptions struct {
	Key            []byte
	Value          []byte
	Extra          []byte
	Datatype       uint8
	Options        uint32
	Flags          uint32
	Expiry         uint32
	Cas            Cas
	RevNo          uint64
	CollectionName string
	ScopeName      string
	CollectionID   uint32
	RetryStrategy  RetryStrategy
	Deadline       time.Time

	// Volatile: Tracer API is subject to change.
	TraceContext RequestSpanContext
}
