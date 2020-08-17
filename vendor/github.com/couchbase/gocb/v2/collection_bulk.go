package gocb

import (
	"time"

	"github.com/couchbase/gocbcore/v9"
)

type bulkOp struct {
	pendop gocbcore.PendingOp
	span   requestSpan
}

func (op *bulkOp) cancel() {
	op.pendop.Cancel()
}

func (op *bulkOp) finish() {
	op.span.Finish()
}

// BulkOp represents a single operation that can be submitted (within a list of more operations) to .Do()
// You can create a bulk operation by instantiating one of the implementations of BulkOp,
// such as GetOp, UpsertOp, ReplaceOp, and more.
// UNCOMMITTED: This API may change in the future.
type BulkOp interface {
	execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
		retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan)
	markError(err error)
	cancel()
	finish()
}

// BulkOpOptions are the set of options available when performing BulkOps using Do.
type BulkOpOptions struct {
	Timeout       time.Duration
	Transcoder    Transcoder
	RetryStrategy RetryStrategy
}

// Do execute one or more `BulkOp` items in parallel.
// UNCOMMITTED: This API may change in the future.
func (c *Collection) Do(ops []BulkOp, opts *BulkOpOptions) error {
	if opts == nil {
		opts = &BulkOpOptions{}
	}

	span := c.startKvOpTrace("Do", nil)

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

func (item *GetOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("GetOp", tracectx)
	item.bulkOp.span = span

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

func (item *GetAndTouchOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("GetAndTouchOp", tracectx)
	item.bulkOp.span = span

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

func (item *TouchOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("TouchOp", tracectx)
	item.bulkOp.span = span

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

func (item *RemoveOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("RemoveOp", tracectx)
	item.bulkOp.span = span

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

func (item *UpsertOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder,
	signal chan BulkOp, retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("UpsertOp", tracectx)
	item.bulkOp.span = span

	etrace := c.startKvOpTrace("encode", span.Context())
	bytes, flags, err := transcoder.Encode(item.Value)
	etrace.Finish()
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

func (item *InsertOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("InsertOp", tracectx)
	item.bulkOp.span = span

	etrace := c.startKvOpTrace("encode", span.Context())
	bytes, flags, err := transcoder.Encode(item.Value)
	if err != nil {
		etrace.Finish()
		item.Err = err
		signal <- item
		return
	}
	etrace.Finish()

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

func (item *ReplaceOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("ReplaceOp", tracectx)
	item.bulkOp.span = span

	etrace := c.startKvOpTrace("encode", span.Context())
	bytes, flags, err := transcoder.Encode(item.Value)
	if err != nil {
		etrace.Finish()
		item.Err = err
		signal <- item
		return
	}
	etrace.Finish()

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

func (item *AppendOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("AppendOp", tracectx)
	item.bulkOp.span = span

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

func (item *PrependOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("PrependOp", tracectx)
	item.bulkOp.span = span

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

func (item *IncrementOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("IncrementOp", tracectx)
	item.bulkOp.span = span

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

func (item *DecrementOp) execute(tracectx requestSpanContext, c *Collection, provider kvProvider, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *retryStrategyWrapper, deadline time.Time, startSpanFunc func(string, requestSpanContext) requestSpan) {
	span := startSpanFunc("DecrementOp", tracectx)
	item.bulkOp.span = span

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
