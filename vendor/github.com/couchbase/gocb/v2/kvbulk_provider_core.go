package gocb

import (
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type kvBulkProviderCore struct {
	agent kvProviderCoreProvider

	tracer RequestTracer
	meter  *meterWrapper
}

func (p *kvBulkProviderCore) Do(c *Collection, ops []BulkOp, opts *BulkOpOptions) error {
	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	span := p.StartKvOpTrace(c, "bulk", tracectx, false)
	defer span.End()

	timeout := opts.Timeout
	if opts.Timeout == 0 {
		timeout = c.timeoutsConfig.KVTimeout * time.Duration(len(ops))
	}

	retryWrapper := c.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryWrapper = newCoreRetryStrategyWrapper(opts.RetryStrategy)
	}

	transcoder := opts.Transcoder
	if transcoder == nil {
		transcoder = c.transcoder
	}

	// Make the channel big enough to hold all our ops in case
	//   we get delayed inside execute (don't want to block the
	//   individual op handlers when they dispatch their signal).
	signal := make(chan BulkOp, len(ops))
	for _, item := range ops {
		switch i := item.(type) {
		case *GetOp:
			p.Get(i, span.Context(), c, transcoder, signal, retryWrapper, time.Now().Add(timeout))
		case *GetAndTouchOp:
			p.GetAndTouch(i, span.Context(), c, transcoder, signal, retryWrapper, time.Now().Add(timeout))
		case *TouchOp:
			p.Touch(i, span.Context(), c, signal, retryWrapper, time.Now().Add(timeout))
		case *RemoveOp:
			p.Delete(i, span.Context(), c, signal, retryWrapper, time.Now().Add(timeout))
		case *UpsertOp:
			p.Set(i, span.Context(), c, transcoder, signal, retryWrapper, time.Now().Add(timeout))
		case *InsertOp:
			p.Add(i, span.Context(), c, transcoder, signal, retryWrapper, time.Now().Add(timeout))
		case *ReplaceOp:
			p.Replace(i, span.Context(), c, transcoder, signal, retryWrapper, time.Now().Add(timeout))
		case *AppendOp:
			p.Append(i, span.Context(), c, signal, retryWrapper, time.Now().Add(timeout))
		case *PrependOp:
			p.Prepend(i, span.Context(), c, signal, retryWrapper, time.Now().Add(timeout))
		case *IncrementOp:
			p.Increment(i, span.Context(), c, signal, retryWrapper, time.Now().Add(timeout))
		case *DecrementOp:
			p.Decrement(i, span.Context(), c, signal, retryWrapper, time.Now().Add(timeout))
		}
	}

	// Wait for all of the ops to complete.
	for range ops {
		item := <-signal
		item.finish()
	}

	return nil
}

func (p *kvBulkProviderCore) Get(item *GetOp, tracectx RequestSpanContext, c *Collection, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "get", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "get", start)
	}

	_, err := p.agent.Get(gocbcore.GetOptions{
		Key:            []byte(item.ID),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.GetResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
		return
	}
}

func (p *kvBulkProviderCore) GetAndTouch(item *GetAndTouchOp, tracectx RequestSpanContext, c *Collection, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "get_and_touch", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "get_and_touch", start)
	}

	_, err := p.agent.GetAndTouch(gocbcore.GetAndTouchOptions{
		Key:            []byte(item.ID),
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.GetAndTouchResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Touch(item *TouchOp, tracectx RequestSpanContext, c *Collection, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "touch", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "touch", start)
	}

	_, err := p.agent.Touch(gocbcore.TouchOptions{
		Key:            []byte(item.ID),
		Expiry:         durationToExpiry(item.Expiry),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.TouchResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Delete(item *RemoveOp, tracectx RequestSpanContext, c *Collection, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "remove", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "remove", start)
	}

	_, err := p.agent.Delete(gocbcore.DeleteOptions{
		Key:            []byte(item.ID),
		Cas:            gocbcore.Cas(item.Cas),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.DeleteResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Set(item *UpsertOp, tracectx RequestSpanContext, c *Collection, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "upsert", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "upsert", start)
	}

	etrace := p.StartKvOpTrace(c, "request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	etrace.End()
	if err != nil {
		item.Err = err
		signal <- item
		return
	}

	_, err = p.agent.Set(gocbcore.SetOptions{
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
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)

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
	}
}

func (p *kvBulkProviderCore) Add(item *InsertOp, tracectx RequestSpanContext, c *Collection, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "insert", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "insert", start)
	}

	etrace := p.StartKvOpTrace(c, "request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	if err != nil {
		etrace.End()
		item.Err = err
		signal <- item
		return
	}
	etrace.End()

	_, err = p.agent.Add(gocbcore.AddOptions{
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
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Replace(item *ReplaceOp, tracectx RequestSpanContext, c *Collection, transcoder Transcoder, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "replace", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "replace", start)
	}

	etrace := p.StartKvOpTrace(c, "request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	if err != nil {
		etrace.End()
		item.Err = err
		signal <- item
		return
	}
	etrace.End()

	_, err = p.agent.Replace(gocbcore.ReplaceOptions{
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
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Append(item *AppendOp, tracectx RequestSpanContext, c *Collection, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "append", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "append", start)
	}

	_, err := p.agent.Append(gocbcore.AdjoinOptions{
		Key:            []byte(item.ID),
		Value:          []byte(item.Value),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.AdjoinResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Prepend(item *PrependOp, tracectx RequestSpanContext, c *Collection, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "prepend", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "prepend", start)
	}

	_, err := p.agent.Prepend(gocbcore.AdjoinOptions{
		Key:            []byte(item.ID),
		Value:          []byte(item.Value),
		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		RetryStrategy:  retryWrapper,
		TraceContext:   span.Context(),
		Deadline:       deadline,
	}, func(res *gocbcore.AdjoinResult, err error) {
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Increment(item *IncrementOp, tracectx RequestSpanContext, c *Collection, signal chan BulkOp,
	retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "increment", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "increment", start)
	}

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if item.Initial > 0 {
		realInitial = uint64(item.Initial)
	}

	_, err := p.agent.Increment(gocbcore.CounterOptions{
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
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}
func (p *kvBulkProviderCore) Decrement(item *DecrementOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp, retryWrapper *coreRetryStrategyWrapper, deadline time.Time) {
	span := p.StartKvOpTrace(c, "decrement", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "decrement", start)
	}

	realInitial := uint64(0xFFFFFFFFFFFFFFFF)
	if item.Initial > 0 {
		realInitial = uint64(item.Initial)
	}

	_, err := p.agent.Decrement(gocbcore.CounterOptions{
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
		item.Err = maybeEnhanceCollKVErr(err, c, item.ID)
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
	}
}

func (p *kvBulkProviderCore) StartKvOpTrace(c *Collection, operationName string, tracectx RequestSpanContext, noAttributes bool) RequestSpan {
	return c.startKvOpTrace(operationName, tracectx, p.tracer, noAttributes)
}
