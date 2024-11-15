package gocb

import (
	"context"
	"time"
)

// LookupInOptions are the set of options available to LookupIn.
type LookupInOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		DocFlags SubdocDocFlag
		User     string
	}

	noMetrics bool
}

// LookupIn performs a set of subdocument lookup operations on the document identified by id.
func (c *Collection) LookupIn(id string, ops []LookupInSpec, opts *LookupInOptions) (docOut *LookupInResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*LookupInResult, error) {
		if opts == nil {
			opts = &LookupInOptions{}
		}

		return agent.LookupIn(c, id, ops, opts)
	})
}

// LookupInAnyReplicaOptions are the set of options available to LookupInAnyReplica.
type LookupInAnyReplicaOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// UNCOMMITTED: This API may change in the future.
	ReadPreference ReadPreference

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		DocFlags SubdocDocFlag
		User     string
	}
}

// LookupInAnyReplica returns the value of a particular document from a replica server.
func (c *Collection) LookupInAnyReplica(id string, ops []LookupInSpec, opts *LookupInAnyReplicaOptions) (*LookupInReplicaResult, error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*LookupInReplicaResult, error) {
		if opts == nil {
			opts = &LookupInAnyReplicaOptions{}
		}

		return agent.LookupInAnyReplica(c, id, ops, opts)
	})
}

// LookupInAllReplicaOptions are the set of options available to LookupInAllReplicas.
type LookupInAllReplicaOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// UNCOMMITTED: This API may change in the future.
	ReadPreference ReadPreference

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		DocFlags SubdocDocFlag
		User     string
	}
}

// LookupInAllReplicas returns the value of a particular document from all replica servers. This will return an iterable
// which streams results one at a time.
func (c *Collection) LookupInAllReplicas(id string, ops []LookupInSpec, opts *LookupInAllReplicaOptions) (*LookupInAllReplicasResult, error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*LookupInAllReplicasResult, error) {
		if opts == nil {
			opts = &LookupInAllReplicaOptions{}
		}

		return agent.LookupInAllReplicas(c, id, ops, opts)
	})
}

// StoreSemantics is used to define the document level action to take during a MutateIn operation.
type StoreSemantics uint8

const (
	// StoreSemanticsReplace signifies to Replace the document, and fail if it does not exist.
	// This is the default action
	StoreSemanticsReplace StoreSemantics = iota

	// StoreSemanticsUpsert signifies to replace the document or create it if it doesn't exist.
	StoreSemanticsUpsert

	// StoreSemanticsInsert signifies to create the document, and fail if it exists.
	StoreSemanticsInsert
)

// MutateInOptions are the set of options available to MutateIn.
type MutateInOptions struct {
	Expiry          time.Duration
	Cas             Cas
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	StoreSemantic   StoreSemantics
	Timeout         time.Duration
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan
	PreserveExpiry  bool

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		DocFlags SubdocDocFlag
		User     string
	}
}

// MutateIn performs a set of subdocument mutations on the document specified by id.
func (c *Collection) MutateIn(id string, ops []MutateInSpec, opts *MutateInOptions) (mutOut *MutateInResult, errOut error) {
	return autoOpControl(c.kvController(), func(agent kvProvider) (*MutateInResult, error) {
		if opts == nil {
			opts = &MutateInOptions{}
		}

		return agent.MutateIn(c, id, ops, opts)
	})
}
