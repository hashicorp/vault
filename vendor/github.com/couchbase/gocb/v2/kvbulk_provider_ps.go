package gocb

import (
	"context"
	"time"

	"github.com/couchbase/goprotostellar/genproto/kv_v1"
)

type kvBulkProviderPs struct {
	client kv_v1.KvServiceClient

	tracer RequestTracer
	meter  *meterWrapper
}

func (p *kvBulkProviderPs) Do(c *Collection, ops []BulkOp, opts *BulkOpOptions) error {
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

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

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
			go p.Get(ctx, i, span.Context(), c, transcoder, signal)
		case *GetAndTouchOp:
			go p.GetAndTouch(ctx, i, span.Context(), c, transcoder, signal)
		case *TouchOp:
			go p.Touch(ctx, i, span.Context(), c, signal)
		case *RemoveOp:
			go p.Remove(ctx, i, span.Context(), c, signal)
		case *UpsertOp:
			go p.Upsert(ctx, i, span.Context(), c, transcoder, signal)
		case *InsertOp:
			go p.Insert(ctx, i, span.Context(), c, transcoder, signal)
		case *ReplaceOp:
			go p.Replace(ctx, i, span.Context(), c, transcoder, signal)
		case *AppendOp:
			go p.Append(ctx, i, span.Context(), c, signal)
		case *PrependOp:
			go p.Prepend(ctx, i, span.Context(), c, signal)
		case *IncrementOp:
			go p.Increment(ctx, i, span.Context(), c, signal)
		case *DecrementOp:
			go p.Decrement(ctx, i, span.Context(), c, signal)
		}
	}

	// Wait for all the ops to complete.
	for range ops {
		item := <-signal
		item.finish()
	}

	return nil
}

func (p *kvBulkProviderPs) Get(ctx context.Context, item *GetOp, tracectx RequestSpanContext, c *Collection,
	transcoder Transcoder, signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "get", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "get", start)
	}

	request := &kv_v1.GetRequest{
		Key: item.ID,

		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		BucketName:     c.bucketName(),
	}

	res, err := p.client.Get(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, true)
		signal <- item
		return
	}

	var content []byte
	switch c := res.Content.(type) {
	case *kv_v1.GetResponse_ContentUncompressed:
		content = c.ContentUncompressed
	case *kv_v1.GetResponse_ContentCompressed:
		content = c.ContentCompressed
		logWarnf("couchbase2 does not currently support compressed content, passing through compressed value")
	}

	item.Result = &GetResult{
		Result: Result{
			cas: Cas(res.Cas),
		},
		transcoder: transcoder,
		contents:   content,
		flags:      res.ContentFlags,
	}
	signal <- item
}

func (p *kvBulkProviderPs) GetAndTouch(ctx context.Context, item *GetAndTouchOp, tracectx RequestSpanContext, c *Collection,
	transcoder Transcoder, signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "get_and_touch", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "get_and_touch", start)
	}

	reqExpiry := &kv_v1.GetAndTouchRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())}

	request := &kv_v1.GetAndTouchRequest{
		Key: item.ID,

		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		BucketName:     c.bucketName(),

		Expiry: reqExpiry,
	}

	res, err := p.client.GetAndTouch(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	var content []byte
	switch c := res.Content.(type) {
	case *kv_v1.GetAndTouchResponse_ContentUncompressed:
		content = c.ContentUncompressed
	case *kv_v1.GetAndTouchResponse_ContentCompressed:
		content = c.ContentCompressed
		logWarnf("couchbase2 does not currently support compressed content, passing through compressed value")
	}

	item.Result = &GetResult{
		Result: Result{
			cas: Cas(res.Cas),
		},
		transcoder: transcoder,
		contents:   content,
		flags:      res.ContentFlags,
	}
	signal <- item
}

func (p *kvBulkProviderPs) Touch(ctx context.Context, item *TouchOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "touch", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "touch", start)
	}

	request := &kv_v1.TouchRequest{
		Key: item.ID,

		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		BucketName:     c.bucketName(),

		Expiry: &kv_v1.TouchRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())},
	}

	res, err := p.client.Touch(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}
	signal <- item
}

func (p *kvBulkProviderPs) Remove(ctx context.Context, item *RemoveOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "remove", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "remove", start)
	}

	var cas *uint64
	if item.Cas > 0 {
		cas = (*uint64)(&item.Cas)
	}

	request := &kv_v1.RemoveRequest{
		Key: item.ID,

		CollectionName: c.name(),
		ScopeName:      c.ScopeName(),
		BucketName:     c.bucketName(),

		Cas: cas,
	}

	res, err := p.client.Remove(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)
	outCas := res.Cas

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(outCas),
		},
	}
	signal <- item
}

func (p *kvBulkProviderPs) Upsert(ctx context.Context, item *UpsertOp, tracectx RequestSpanContext, c *Collection,
	transcoder Transcoder, signal chan BulkOp) {
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

	var expiry *kv_v1.UpsertRequest_ExpirySecs
	if item.Expiry > 0 {
		expiry = &kv_v1.UpsertRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())}
	}

	content := &kv_v1.UpsertRequest_ContentUncompressed{ContentUncompressed: bytes}

	request := &kv_v1.UpsertRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Content:        content,
		ContentFlags:   flags,

		Expiry: expiry,
	}

	res, err := p.client.Upsert(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(res.Cas),
		},
	}

	signal <- item
}

func (p *kvBulkProviderPs) Insert(ctx context.Context, item *InsertOp, tracectx RequestSpanContext, c *Collection,
	transcoder Transcoder, signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "insert", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "insert", start)
	}

	etrace := p.StartKvOpTrace(c, "request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	etrace.End()
	if err != nil {
		item.Err = err
		signal <- item
		return
	}

	var expiry *kv_v1.InsertRequest_ExpirySecs
	if item.Expiry > 0 {
		expiry = &kv_v1.InsertRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())}
	}

	content := &kv_v1.InsertRequest_ContentUncompressed{ContentUncompressed: bytes}

	request := &kv_v1.InsertRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Content:        content,
		ContentFlags:   flags,

		Expiry: expiry,
	}

	res, err := p.client.Insert(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(res.Cas),
		},
	}

	signal <- item
}

func (p *kvBulkProviderPs) Replace(ctx context.Context, item *ReplaceOp, tracectx RequestSpanContext, c *Collection,
	transcoder Transcoder, signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "replace", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "replace", start)
	}

	etrace := p.StartKvOpTrace(c, "request_encoding", span.Context(), true)
	bytes, flags, err := transcoder.Encode(item.Value)
	etrace.End()
	if err != nil {
		item.Err = err
		signal <- item
		return
	}

	var cas *uint64
	if item.Cas > 0 {
		cas = (*uint64)(&item.Cas)
	}

	var expiry *kv_v1.ReplaceRequest_ExpirySecs
	if item.Expiry > 0 {
		expiry = &kv_v1.ReplaceRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())}
	}

	content := &kv_v1.ReplaceRequest_ContentUncompressed{ContentUncompressed: bytes}

	request := &kv_v1.ReplaceRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Content:        content,
		ContentFlags:   flags,
		Cas:            cas,

		Expiry: expiry,
	}

	res, err := p.client.Replace(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(res.Cas),
		},
	}

	signal <- item
}

func (p *kvBulkProviderPs) Append(ctx context.Context, item *AppendOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "append", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "append", start)
	}

	request := &kv_v1.AppendRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Content:        []byte(item.Value),
	}

	res, err := p.client.Append(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(res.Cas),
		},
	}

	signal <- item
}

func (p *kvBulkProviderPs) Prepend(ctx context.Context, item *PrependOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "prepend", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "prepend", start)
	}

	request := &kv_v1.PrependRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Content:        []byte(item.Value),
	}

	res, err := p.client.Prepend(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &MutationResult{
		mt: mt,
		Result: Result{
			cas: Cas(res.Cas),
		},
	}

	signal <- item
}

func (p *kvBulkProviderPs) Increment(ctx context.Context, item *IncrementOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "increment", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "increment", start)
	}

	var expiry *kv_v1.IncrementRequest_ExpirySecs
	if item.Expiry > 0 {
		expiry = &kv_v1.IncrementRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())}
	}

	request := &kv_v1.IncrementRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Delta:          uint64(item.Delta),
		Expiry:         expiry,
		Initial:        &item.Initial,
	}

	res, err := p.client.Increment(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &CounterResult{}
	item.Result.mt = mt
	item.Result.cas = Cas(res.Cas)
	item.Result.content = uint64(res.Content)

	signal <- item
}

func (p *kvBulkProviderPs) Decrement(ctx context.Context, item *DecrementOp, tracectx RequestSpanContext, c *Collection,
	signal chan BulkOp) {
	span := p.StartKvOpTrace(c, "decrement", tracectx, false)
	start := time.Now()
	item.bulkOp.finishFn = func() {
		span.End()
		p.meter.ValueRecord(meterValueServiceKV, "decrement", start)
	}

	var expiry *kv_v1.DecrementRequest_ExpirySecs
	if item.Expiry > 0 {
		expiry = &kv_v1.DecrementRequest_ExpirySecs{ExpirySecs: uint32(item.Expiry.Seconds())}
	}

	request := &kv_v1.DecrementRequest{
		Key:            item.ID,
		BucketName:     c.bucketName(),
		ScopeName:      c.ScopeName(),
		CollectionName: c.name(),
		Delta:          uint64(item.Delta),
		Expiry:         expiry,
		Initial:        &item.Initial,
	}

	res, err := p.client.Decrement(ctx, request)
	if err != nil {
		item.Err = mapPsErrorToGocbError(err, false)
		signal <- item
		return
	}

	mt := psMutToGoCbMut(res.MutationToken)

	item.Result = &CounterResult{}
	item.Result.mt = mt
	item.Result.cas = Cas(res.Cas)
	item.Result.content = uint64(res.Content)

	signal <- item
}

func (p *kvBulkProviderPs) StartKvOpTrace(c *Collection, operationName string, tracectx RequestSpanContext, noAttributes bool) RequestSpan {
	return c.startKvOpTrace(operationName, tracectx, p.tracer, noAttributes)
}
