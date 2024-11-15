//nolint:unused
package gocb

import (
	"context"
	"time"
)

type bulkOp struct {
	finishFn func()
}

func (op *bulkOp) finish() {
	op.finishFn()
}

// BulkOp represents a single operation that can be submitted (within a list of more operations) to .Do()
// You can create a bulk operation by instantiating one of the implementations of BulkOp,
// such as GetOp, UpsertOp, ReplaceOp, and more.
// UNCOMMITTED: This API may change in the future.
type BulkOp interface {
	isBulkOp()
	finish()
}

// BulkOpOptions are the set of options available when performing BulkOps using Do.
type BulkOpOptions struct {
	Timeout       time.Duration
	Transcoder    Transcoder
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// Do execute one or more `BulkOp` items in parallel.
// UNCOMMITTED: This API may change in the future.
func (c *Collection) Do(ops []BulkOp, opts *BulkOpOptions) error {
	return autoOpControlErrorOnly(c.kvBulkController(), func(agent kvBulkProvider) error {
		if opts == nil {
			opts = &BulkOpOptions{}
		}

		return agent.Do(c, ops, opts)
	})
}

// GetOp represents a type of `BulkOp` used for Get operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type GetOp struct {
	bulkOp

	ID     string
	Result *GetResult
	Err    error
}

func (item *GetOp) isBulkOp() {}

// GetAndTouchOp represents a type of `BulkOp` used for GetAndTouch operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type GetAndTouchOp struct {
	bulkOp

	ID     string
	Expiry time.Duration
	Result *GetResult
	Err    error
}

func (item *GetAndTouchOp) isBulkOp() {}

// TouchOp represents a type of `BulkOp` used for Touch operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type TouchOp struct {
	bulkOp

	ID     string
	Expiry time.Duration
	Result *MutationResult
	Err    error
}

func (item *TouchOp) isBulkOp() {}

// RemoveOp represents a type of `BulkOp` used for Remove operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type RemoveOp struct {
	bulkOp

	ID     string
	Cas    Cas
	Result *MutationResult
	Err    error
}

func (item *RemoveOp) isBulkOp() {}

// UpsertOp represents a type of `BulkOp` used for Upsert operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type UpsertOp struct {
	bulkOp

	ID     string
	Value  interface{}
	Expiry time.Duration
	Cas    Cas
	Result *MutationResult
	Err    error
}

func (item *UpsertOp) isBulkOp() {}

// InsertOp represents a type of `BulkOp` used for Insert operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type InsertOp struct {
	bulkOp

	ID     string
	Value  interface{}
	Expiry time.Duration
	Result *MutationResult
	Err    error
}

func (item *InsertOp) isBulkOp() {}

// ReplaceOp represents a type of `BulkOp` used for Replace operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type ReplaceOp struct {
	bulkOp

	ID     string
	Value  interface{}
	Expiry time.Duration
	Cas    Cas
	Result *MutationResult
	Err    error
}

func (item *ReplaceOp) isBulkOp() {}

// AppendOp represents a type of `BulkOp` used for Append operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type AppendOp struct {
	bulkOp

	ID     string
	Value  string
	Result *MutationResult
	Err    error
}

func (item *AppendOp) isBulkOp() {}

// PrependOp represents a type of `BulkOp` used for Prepend operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type PrependOp struct {
	bulkOp

	ID     string
	Value  string
	Result *MutationResult
	Err    error
}

func (item *PrependOp) isBulkOp() {}

// IncrementOp represents a type of `BulkOp` used for Increment operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type IncrementOp struct {
	bulkOp

	ID      string
	Delta   int64
	Initial int64
	Expiry  time.Duration

	Result *CounterResult
	Err    error
}

func (item *IncrementOp) isBulkOp() {}

// DecrementOp represents a type of `BulkOp` used for Decrement operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type DecrementOp struct {
	bulkOp

	ID      string
	Delta   int64
	Initial int64
	Expiry  time.Duration

	Result *CounterResult
	Err    error
}

func (item *DecrementOp) isBulkOp() {}
