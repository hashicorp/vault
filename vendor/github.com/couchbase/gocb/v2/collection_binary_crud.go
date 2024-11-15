package gocb

import (
	"context"
	"time"
)

// BinaryCollection is a set of binary operations.
type BinaryCollection struct {
	collection *Collection
}

// AppendOptions are the options available to the Append operation.
type AppendOptions struct {
	Timeout         time.Duration
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	Cas             Cas
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Append appends a byte value to a document.
func (c *BinaryCollection) Append(id string, val []byte, opts *AppendOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.collection.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &AppendOptions{}
		}

		return agent.Append(c.collection, id, val, opts)
	})
}

// PrependOptions are the options available to the Prepend operation.
type PrependOptions struct {
	Timeout         time.Duration
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	Cas             Cas
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Prepend prepends a byte value to a document.
func (c *BinaryCollection) Prepend(id string, val []byte, opts *PrependOptions) (mutOut *MutationResult, errOut error) {
	return autoOpControl(c.collection.kvController(), func(agent kvProvider) (*MutationResult, error) {
		if opts == nil {
			opts = &PrependOptions{}
		}

		return agent.Prepend(c.collection, id, val, opts)
	})
}

// IncrementOptions are the options available to the Increment operation.
type IncrementOptions struct {
	Timeout time.Duration
	// Expiry is the length of time that the document will be stored in Couchbase.
	// A value of 0 will set the document to never expire.
	Expiry time.Duration
	// Initial, if non-negative, is the `initial` value to use for the document if it does not exist.
	// If present, this is the value that will be returned by a successful operation.
	Initial int64
	// Delta is the value to use for incrementing/decrementing if Initial is not present.
	Delta           uint64
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan

	// Deprecated: Cas is not supported by the server for Increment, and is no longer used.
	Cas Cas

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Increment performs an atomic addition for an integer document. Passing a
// non-negative `initial` value will cause the document to be created if it did not
// already exist.
func (c *BinaryCollection) Increment(id string, opts *IncrementOptions) (countOut *CounterResult, errOut error) {
	return autoOpControl(c.collection.kvController(), func(agent kvProvider) (*CounterResult, error) {
		if opts == nil {
			opts = &IncrementOptions{}
		}
		if opts.Cas > 0 {
			return nil, makeInvalidArgumentsError("cas is not supported for the Increment operation")
		}

		return agent.Increment(c.collection, id, opts)
	})
}

// DecrementOptions are the options available to the Decrement operation.
type DecrementOptions struct {
	Timeout time.Duration
	// Expiry is the length of time that the document will be stored in Couchbase.
	// A value of 0 will set the document to never expire.
	Expiry time.Duration
	// Initial, if non-negative, is the `initial` value to use for the document if it does not exist.
	// If present, this is the value that will be returned by a successful operation.
	Initial int64
	// Delta is the value to use for incrementing/decrementing if Initial is not present.
	Delta           uint64
	DurabilityLevel DurabilityLevel
	PersistTo       uint
	ReplicateTo     uint
	RetryStrategy   RetryStrategy
	ParentSpan      RequestSpan

	// Deprecated: Cas is not supported by the server for Decrement, and is no longer used.
	Cas Cas

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

// Decrement performs an atomic subtraction for an integer document. Passing a
// non-negative `initial` value will cause the document to be created if it did not
// already exist.
func (c *BinaryCollection) Decrement(id string, opts *DecrementOptions) (countOut *CounterResult, errOut error) {
	return autoOpControl(c.collection.kvController(), func(agent kvProvider) (*CounterResult, error) {
		if opts == nil {
			opts = &DecrementOptions{}
		}
		if opts.Cas > 0 {
			return nil, makeInvalidArgumentsError("cas is not supported for the Decrement operation")
		}

		return agent.Decrement(c.collection, id, opts)
	})
}
