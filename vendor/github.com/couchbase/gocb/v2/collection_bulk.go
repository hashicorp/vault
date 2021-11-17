// nolint: unused
package gocb

import (
	"context"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type bulkOp struct {
	pendop   gocbcore.PendingOp
	finishFn func()
}

func (op *bulkOp) cancel() {
	op.pendop.Cancel()
}

func (op *bulkOp) finish() {
	op.finishFn()
}

// BulkOp represents a single operation that can be submitted (within a list of more operations) to .Do()
// You can create a bulk operation by instantiating one of the implementations of BulkOp,
// such as GetOp, UpsertOp, ReplaceOp, and more.
// UNCOMMITTED: This API may change in the future.
type BulkOp interface {
	execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
		retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan)
	markError(err error)
	cancel()
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
	if opts == nil {
		opts = &BulkOpOptions{}
	}

	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	span := c.startKvOpTrace("bulk", tracectx, false)
	defer span.End()

	timeout := opts.Timeout
	if opts.Timeout == 0 {
		timeout = c.timeoutsConfig.KVTimeout * time.Duration(len(ops))
	}

	retryWrapper := c.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryWrapper = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	if opts.Transcoder == nil {
		opts.Transcoder = c.transcoder
	}

	agent, err := c.getKvProvider()
	if err != nil {
		return err
	}

	// Make the channel big enough to hold all our ops in case
	//   we get delayed inside execute (don't want to block the
	//   individual op handlers when they dispatch their signal).
	signal := make(chan BulkOp, len(ops))
	for _, item := range ops {
		item.execute(span.Context(), c, agent, opts.Transcoder, signal, retryWrapper, time.Now().Add(timeout), c.startKvOpTrace)
	}

	for range ops {
		item := <-signal
		// We're really just clearing the pendop from this thread,
		//   since it already completed, no cancel actually occurs
		item.finish()
	}
	return nil
}

// GetOp represents a type of `BulkOp` used for Get operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type GetOp struct {
	bulkOp

	ID     string
	Result *GetResult
	Err    error
}

func (item *GetOp) markError(err error) {
	item.Err = err
}

func (item *GetOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("get", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "get", start)
	}

	op, err := provider.Get(gocbcore.GetOptions{
		Key:            []byte(item.ID),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.GetResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: transcoder,
				contents:   res.Value,
				flags:      res.Flags,
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

// GetAndTouchOp represents a type of `BulkOp` used for GetAndTouch operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type GetAndTouchOp struct {
	bulkOp

	ID     string
	Expiry time.Duration
	Result *GetResult
	Err    error
}

func (item *GetAndTouchOp) markError(err error) {
	item.Err = err
}

func (item *GetAndTouchOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("get_and_touch", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "get_and_touch", start)
	}

	op, err := provider.GetAndTouch(gocbcore.GetAndTouchOptions{
		Key:            []byte(item.ID),
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.GetAndTouchResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: transcoder,
				contents:   res.Value,
				flags:      res.Flags,
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

// TouchOp represents a type of `BulkOp` used for Touch operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type TouchOp struct {
	bulkOp

	ID     string
	Expiry time.Duration
	Result *MutationResult
	Err    error
}

func (item *TouchOp) markError(err error) {
	item.Err = err
}

func (item *TouchOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("touch", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "touch", start)
	}

	op, err := provider.Touch(gocbcore.TouchOptions{
		Key:            []byte(item.ID),
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.TouchResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

// RemoveOp represents a type of `BulkOp` used for Remove operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type RemoveOp struct {
	bulkOp

	ID     string
	Cas    Cas
	Result *MutationResult
	Err    error
}

func (item *RemoveOp) markError(err error) {
	item.Err = err
}

func (item *RemoveOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("remove", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "remove", start)
	}

	op, err := provider.Delete(gocbcore.DeleteOptions{
		Key:            []byte(item.ID),
		Cas:            gocbcore.Cas(item.Cas),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.DeleteResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

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

func (item *UpsertOp) markError(err error) {
	item.Err = err
}

func (item *UpsertOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder,
	signal chan BulkOp, retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("upsert", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "upsert", start)
	}

	etrace := c.startKvOpTrace("request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	etrace.End()
	if err != nil {
		item.Err = err
		signal <- item
		return
	}

	op, err := provider.Set(gocbcore.SetOptions{
		Key:            []byte(item.ID),
		Value:          bytes,
		Flags:          flags,
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.StoreResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)

		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

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

func (item *InsertOp) markError(err error) {
	item.Err = err
}

func (item *InsertOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("insert", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "insert", start)
	}

	etrace := c.startKvOpTrace("request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	if err != nil {
		etrace.End()
		item.Err = err
		signal <- item
		return
	}
	etrace.End()

	op, err := provider.Add(gocbcore.AddOptions{
		Key:            []byte(item.ID),
		Value:          bytes,
		Flags:          flags,
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.StoreResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

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

func (item *ReplaceOp) markError(err error) {
	item.Err = err
}

func (item *ReplaceOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("replace", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "replace", start)
	}

	etrace := c.startKvOpTrace("request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	if err != nil {
		etrace.End()
		item.Err = err
		signal <- item
		return
	}
	etrace.End()

	op, err := provider.Replace(gocbcore.ReplaceOptions{
		Key:            []byte(item.ID),
		Value:          bytes,
		Flags:          flags,
		Cas:            gocbcore.Cas(item.Cas),
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.StoreResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

// AppendOp represents a type of `BulkOp` used for Append operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type AppendOp struct {
	bulkOp

	ID     string
	Value  string
	Result *MutationResult
	Err    error
}

func (item *AppendOp) markError(err error) {
	item.Err = err
}

func (item *AppendOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("append", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "append", start)
	}

	op, err := provider.Append(gocbcore.AdjoinOptions{
		Key:            []byte(item.ID),
		Value:          []byte(item.Value),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.AdjoinResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

// PrependOp represents a type of `BulkOp` used for Prepend operations. See BulkOp.
// UNCOMMITTED: This API may change in the future.
type PrependOp struct {
	bulkOp

	ID     string
	Value  string
	Result *MutationResult
	Err    error
}

func (item *PrependOp) markError(err error) {
	item.Err = err
}

func (item *PrependOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("prepend", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "prepend", start)
	}

	op, err := provider.Prepend(gocbcore.AdjoinOptions{
		Key:            []byte(item.ID),
		Value:          []byte(item.Value),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.AdjoinResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &MutationResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

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

func (item *IncrementOp) markError(err error) {
	item.Err = err
}

func (item *IncrementOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("increment", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "increment", start)
	}

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if item.Initial > 0 {
		realInitial = uint64(item.Initial)
	}

	op, err := provider.Increment(gocbcore.CounterOptions{
		Key:            []byte(item.ID),
		Delta:          uint64(item.Delta),
		Initial:        realInitial,
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.CounterResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &CounterResult{
				MutationResult: MutationResult{
					Result: Result{
						cas: Cas(res.Cas),
					},
				},
				content: res.Value,
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}

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

func (item *DecrementOp) markError(err error) {
	item.Err = err
}

func (item *DecrementOp) execute(tracectx RequestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, RequestSpanContext, bool) RequestSpan) {
	span := startSpanFunc("decrement", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		c.meter.ValueRecord(meterValueServiceKV, "decrement", start)
	}

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if item.Initial > 0 {
		realInitial = uint64(item.Initial)
	}

	op, err := provider.Decrement(gocbcore.CounterOptions{
		Key:            []byte(item.ID),
		Delta:          uint64(item.Delta),
		Initial:        realInitial,
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.CounterResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, provider, c, item.ID)
		if item.Err == nil {
			item.Result = &CounterResult{
				MutationResult: MutationResult{
					Result: Result{
						cas: Cas(res.Cas),
					},
				},
				content: res.Value,
			}

			if res.MutationToken.VbUUID != 0 {
				mutTok := &MutationToken{
					token:      res.MutationToken,
					bucketName: c.bucketName(),
				}
				item.Result.mt = mutTok
			}
		}
		signal <- item
	})
	if err != nil {
		item.Err = err
		signal <- item
	} else {
		item.bulkOp.pendop = op
	}
}
