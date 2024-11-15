package gocb

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/couchbase/gocbcore/v10"
)

const (
	rangeScanDefaultItemLimit        = 50
	rangeScanDefaultBytesLimit       = 15000
	rangeScanDefaultConcurrency      = 1
	rangeScanDefaultResultBufferSize = 1024
)

type rangeScanOpManager struct {
	err error

	span        RequestSpan
	transcoder  Transcoder
	timeout     time.Duration
	deadline    time.Time
	impersonate string

	cancelCh chan struct{}

	agent       kvProviderCoreProvider
	createdTime time.Time
	meter       *meterWrapper
	tracer      RequestTracer

	defaultRetryStrategy *coreRetryStrategyWrapper
	defaultTranscoder    Transcoder
	defaultTimeout       time.Duration

	cid        uint32
	bucketName string

	// These are used only for creating errors.
	scopeName      string
	collectionName string

	rangeOptions          *gocbcore.RangeScanCreateRangeScanConfig
	samplingOptions       *gocbcore.RangeScanCreateRandomSamplingConfig
	vBucketToSnapshotOpts map[uint16]gocbcore.RangeScanCreateSnapshotRequirements

	numVbuckets        int
	serverToVbucketMap map[int][]uint16
	keysOnly           bool
	itemLimit          uint32
	byteLimit          uint32
	maxConcurrency     uint16

	result *ScanResult

	cancelled uint32
}

type rangeScanVbucket struct {
	id     uint16
	server int
}

func (p *kvProviderCore) newRangeScanOpManager(c *Collection, scanType ScanType, agent kvProviderCoreProvider,
	parentSpan RequestSpan, consistentWith *MutationState, keysOnly bool) (*rangeScanOpManager, error) {
	var tracectx RequestSpanContext
	if parentSpan != nil {
		tracectx = parentSpan.Context()
	}

	span := p.tracer.RequestSpan(tracectx, "range_scan")
	span.SetAttribute(spanAttribDBNameKey, c.bucket.Name())
	span.SetAttribute(spanAttribDBCollectionNameKey, c.Name())
	span.SetAttribute(spanAttribDBScopeNameKey, c.ScopeName())
	span.SetAttribute(spanAttribServiceKey, "kv_scan")
	span.SetAttribute(spanAttribOperationKey, "range_scan")
	span.SetAttribute(spanAttribDBSystemKey, spanAttribDBSystemValue)
	span.SetAttribute("without_content", keysOnly)

	var rangeOptions *gocbcore.RangeScanCreateRangeScanConfig
	var samplingOptions *gocbcore.RangeScanCreateRandomSamplingConfig

	setRangeScanOpts := func(st RangeScan) error {
		if st.To == nil {
			st.To = ScanTermMaximum()
		}
		if st.From == nil {
			st.From = ScanTermMinimum()
		}

		span.SetAttribute("scan_type", "range")
		span.SetAttribute("from_term", st.From.Term)
		span.SetAttribute("to_term", st.To.Term)
		var err error
		rangeOptions, err = st.toCore()
		if err != nil {
			return err
		}

		return nil
	}

	setSamplingScanOpts := func(st SamplingScan) error {
		if st.Seed == 0 {
			st.Seed = rand.Uint64() // #nosec G404
		}
		span.SetAttribute("scan_type", "sampling")
		span.SetAttribute("limit", st.Limit)
		span.SetAttribute("seed", st.Seed)
		var err error
		samplingOptions, err = st.toCore()
		if err != nil {
			return err
		}

		return nil
	}

	var err error
	switch st := scanType.(type) {
	case RangeScan:
		if err := setRangeScanOpts(st); err != nil {
			return nil, err
		}
	case *RangeScan:
		if err := setRangeScanOpts(*st); err != nil {
			return nil, err
		}
	case SamplingScan:
		if err := setSamplingScanOpts(st); err != nil {
			return nil, err
		}
	case *SamplingScan:
		if err := setSamplingScanOpts(*st); err != nil {
			return nil, err
		}
	default:
		err = makeInvalidArgumentsError("only RangeScan and SamplingScan are supported for ScanType")
	}

	vBucketToSnapshotOpts := make(map[uint16]gocbcore.RangeScanCreateSnapshotRequirements)
	if consistentWith != nil {
		for _, token := range consistentWith.tokens {
			entry, ok := vBucketToSnapshotOpts[uint16(token.PartitionID())]
			if ok {
				if uint64(entry.VbUUID) != token.PartitionUUID() {
					return nil, makeInvalidArgumentsError("mutation state contained two token with same " +
						"partition id but different partition uuids")
				}
				// Only replace if the token seqno > existing seqno
				seqno := gocbcore.SeqNo(token.SequenceNumber())
				if seqno > entry.SeqNo {
					vBucketToSnapshotOpts[uint16(token.PartitionID())] = gocbcore.RangeScanCreateSnapshotRequirements{
						VbUUID: gocbcore.VbUUID(token.PartitionUUID()),
						SeqNo:  seqno,
					}
				}
			} else {
				vBucketToSnapshotOpts[uint16(token.PartitionID())] = gocbcore.RangeScanCreateSnapshotRequirements{
					VbUUID: gocbcore.VbUUID(token.PartitionUUID()),
					SeqNo:  gocbcore.SeqNo(token.SequenceNumber()),
				}
			}
		}
	}

	m := &rangeScanOpManager{
		err: err,

		span:        span,
		createdTime: time.Now(),
		meter:       p.meter,
		tracer:      p.tracer,

		cancelCh: make(chan struct{}),

		agent:                agent,
		defaultTimeout:       c.timeoutsConfig.KVScanTimeout,
		defaultTranscoder:    c.transcoder,
		defaultRetryStrategy: c.retryStrategyWrapper,
		bucketName:           c.Bucket().Name(),

		scopeName:      c.ScopeName(),
		collectionName: c.Name(),

		rangeOptions:          rangeOptions,
		samplingOptions:       samplingOptions,
		vBucketToSnapshotOpts: vBucketToSnapshotOpts,
		keysOnly:              keysOnly,
	}

	return m, nil
}

func (m *rangeScanOpManager) getTimeout() time.Duration {
	if m.timeout > 0 {
		return m.timeout
	}

	return m.defaultTimeout
}

func (m *rangeScanOpManager) SetCollectionID(cid uint32) {
	m.cid = cid
}

func (m *rangeScanOpManager) SetNumVbuckets(numVbuckets int) {
	m.numVbuckets = numVbuckets
	m.span.SetAttribute("num_partitions", numVbuckets)
}

func (m *rangeScanOpManager) SetServerToVbucketMap(serverVbucketMap map[int][]uint16) {
	m.serverToVbucketMap = serverVbucketMap
}

func (m *rangeScanOpManager) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *rangeScanOpManager) SetItemLimit(limit *uint32) {
	if limit == nil {
		m.itemLimit = rangeScanDefaultItemLimit
	} else {
		m.itemLimit = *limit
	}
}

func (m *rangeScanOpManager) SetByteLimit(limit *uint32) {
	if limit == nil {
		m.byteLimit = rangeScanDefaultBytesLimit
	} else {
		m.byteLimit = *limit
	}
}

func (m *rangeScanOpManager) SetMaxConcurrency(max uint16) {
	if max == 0 {
		max = rangeScanDefaultConcurrency
	}
	m.maxConcurrency = max
}

func (m *rangeScanOpManager) SetResult(result *ScanResult) {
	m.result = result
}

func (m *rangeScanOpManager) SetTranscoder(transcoder Transcoder) {
	if transcoder == nil {
		transcoder = m.defaultTranscoder
	}
	m.transcoder = transcoder
}

func (m *rangeScanOpManager) SetImpersonate(user string) {
	m.impersonate = user
}

func (m *rangeScanOpManager) Finish() {
	m.span.End()

	m.meter.ValueRecord(meterValueServiceKV, "range_scan", m.createdTime)
}

func (m *rangeScanOpManager) TraceSpanContext() RequestSpanContext {
	return m.span.Context()
}

func (m *rangeScanOpManager) TraceSpan() RequestSpan {
	return m.span
}

func (m *rangeScanOpManager) CID() uint32 {
	return m.cid
}

func (m *rangeScanOpManager) BucketName() string {
	return m.bucketName
}

func (m *rangeScanOpManager) Transcoder() Transcoder {
	return m.transcoder
}

func (m *rangeScanOpManager) RangeOptions() *gocbcore.RangeScanCreateRangeScanConfig {
	return m.rangeOptions
}

func (m *rangeScanOpManager) SamplingOptions() *gocbcore.RangeScanCreateRandomSamplingConfig {
	return m.samplingOptions
}

func (m *rangeScanOpManager) SnapshotOptions(vbID uint16) *gocbcore.RangeScanCreateSnapshotRequirements {
	opts, ok := m.vBucketToSnapshotOpts[vbID]
	if !ok {
		return nil
	}

	return &opts
}

func (m *rangeScanOpManager) KeysOnly() bool {
	return m.keysOnly
}

func (m *rangeScanOpManager) CheckReadyForOp() error {
	if m.err != nil {
		return m.err
	}

	timeout := m.getTimeout()
	if timeout == 0 {
		return errors.New("range scan op manager had no timeout specified")
	}

	m.deadline = time.Now().Add(timeout)

	if m.numVbuckets == 0 {
		return errors.New("range scan op manager had no number of partitions specified")
	}

	return nil
}

func (m *rangeScanOpManager) EnhanceErr(err error) error {
	return maybeEnhanceKVErr(err, m.bucketName, m.scopeName, m.collectionName, "scan")
}

func (m *rangeScanOpManager) Deadline() time.Time {
	return m.deadline
}

func (m *rangeScanOpManager) Timeout() time.Duration {
	return m.getTimeout()
}

func (m *rangeScanOpManager) Impersonate() string {
	return m.impersonate
}

// Cancel will trigger all underlying streams to cancel themselves.
func (m *rangeScanOpManager) Cancel(err error) {
	m.cancelScan(err)
}

func (m *rangeScanOpManager) IsRangeScan() bool {
	return m.rangeOptions != nil
}

func (m *rangeScanOpManager) cancelScan(err error) {
	if atomic.CompareAndSwapUint32(&m.cancelled, 0, 1) {
		if err != nil {
			m.result.setErr(err)
		}
		close(m.cancelCh)
	}
}

func (m *rangeScanOpManager) Scan(ctx context.Context) (*ScanResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var limit uint64
	if m.SamplingOptions() != nil {
		limit = m.SamplingOptions().Samples
	}

	resultCh := make(chan *ScanResultItem, rangeScanDefaultResultBufferSize)
	r := &ScanResult{
		resultChan: resultCh,
		cancelFn:   m.Cancel,

		limit: limit,
	}
	m.SetResult(r)

	balancer := m.createLoadBalancer()

	var complete uint32
	var seenData uint32

	// We keep separate counts of running and completed to simplify shutdown of the scan.
	scansRunning := int32(m.maxConcurrency)

	isRangeScan := m.IsRangeScan()

	var i uint16
	for i = 0; i < m.maxConcurrency; i++ {
		go func() {
			defer func() {
				if atomic.AddUint32(&complete, 1) == uint32(m.maxConcurrency) {
					m.Finish()
					balancer.close()
					close(resultCh)
				}
			}()

			for vbucket, ok := balancer.selectVbucket(); ok; vbucket, ok = balancer.selectVbucket() {
				if atomic.LoadUint32(&m.cancelled) == 1 {
					return
				}

				deadline := time.Now().Add(m.Timeout())
				failPoint, err := m.scanPartition(ctx, deadline, vbucket.id, resultCh)
				balancer.scanEnded(vbucket)
				if err != nil {
					err = m.EnhanceErr(err)
					if failPoint == scanFailPointCreate {
						if errors.Is(err, gocbcore.ErrDocumentNotFound) {
							logDebugf("Ignoring vbid %d as no documents exist for that vbucket", vbucket.id)
							continue
						}

						if errors.Is(err, ErrTemporaryFailure) || errors.Is(err, gocbcore.ErrBusy) {
							// Put the vbucket back into the channel to be retried later.
							balancer.retryScan(vbucket)

							if errors.Is(err, gocbcore.ErrBusy) {
								// Busy indicates that the server is reporting too many active scans.
								// Shut ourselves down if we're not the only runner remaining.
								running := atomic.AddInt32(&scansRunning, -1)
								if running >= 1 {
									// Shutdown this worker.
									logDebugf("Shutting down scan runner, remaining %d", running)
									return
								}
							}

							continue
						}

						if !m.IsRangeScan() {
							continue
						}

						// All other errors are fatal.
						m.cancelScan(err)
						return
					}
					// For range scan these are fatal.
					var retErr error
					if errors.Is(err, ErrDocumentNotFound) {
						if isRangeScan {
							retErr = err
						}
					} else if errors.Is(err, ErrAuthenticationFailure) {
						if isRangeScan {
							retErr = err
						}
					} else if errors.Is(err, ErrCollectionNotFound) {
						if isRangeScan {
							retErr = err
						}
					} else if errors.Is(err, gocbcore.ErrRangeScanCancelled) {
						if isRangeScan {
							var kvError *KeyValueError
							if errors.As(err, &kvError) {
								kvError.InnerError = ErrRequestCanceled
								retErr = kvError
							} else {
								retErr = ErrRequestCanceled
							}
						}
					} else {
						// Any other error is fatal.
						retErr = err
					}
					if retErr != nil {
						m.cancelScan(retErr)
						return
					}

					continue
				}
			}
		}()
	}
	// Block waiting for any errors on the first scan(s) so that we can immediately return that error.
	select {
	case <-m.cancelCh:
		return nil, r.Err()
	case item, more := <-resultCh:
		// more could be false if no sampling scans returned any data, but that isn't an error case.
		if more {
			atomic.StoreUint32(&seenData, 1)
			r.peeked = unsafe.Pointer(item)
		}
	}

	return r, nil
}

type scanFailPoint uint8

const (
	scanFailPointCreate scanFailPoint = iota + 1
	scanFailPointContinue
)

func (m *rangeScanOpManager) scanPartition(ctx context.Context, deadline time.Time, vbID uint16, resultCh chan *ScanResultItem) (scanFailPoint, error) {
	span := m.tracer.RequestSpan(m.span.Context(), "range_scan_partition")
	span.SetAttribute("partition_id", vbID)
	defer span.End()
	var lastTermSeen []byte
	rangeOpts := m.RangeOptions()
	samplingOpts := m.SamplingOptions()

	var createRes gocbcore.RangeScanCreateResult
	for {
		if rangeOpts != nil && len(lastTermSeen) > 0 {
			// Make a copy of the range options so that we don't affect the manager level ones.
			newRangeOpts := *rangeOpts
			newRangeOpts.Start = lastTermSeen
			rangeOpts = &newRangeOpts
		}

		var err error
		createRes, err = m.createStream(ctx, span.Context(), deadline, vbID, rangeOpts, samplingOpts)
		if err != nil {
			err = m.EnhanceErr(err)
			return scanFailPointCreate, err
		}

		// We only apply context to the initial create stream request, after that we consider the stream active
		// and context cancellation no longer applies.
		ctx = context.Background()

		// We've created the stream so now loop continue until the stream is complete or cancelled.
		for {
			items, isComplete, err := m.continueStream(ctx, span.Context(), createRes)
			if err != nil {
				err = m.EnhanceErr(err)
				// If the error is NMV or EOF then we should recreate the stream from the last known item.
				// Breaking here without calling cancel will trigger us to reloop rather than call Cancel on
				// the stream and then return.
				if errors.Is(err, gocbcore.ErrNotMyVBucket) || errors.Is(err, io.EOF) {
					logInfof("Received NotMyVbucket or EOF, will retry")
					break
				}
				return scanFailPointContinue, err
			}
			if len(items) > 0 {
				for _, item := range items {
					var expiry time.Time
					if item.Expiry > 0 {
						expiry = time.Unix(int64(item.Expiry), 0)
					}
					select {
					case <-m.cancelCh:
						if !isComplete {
							m.cancelStream(ctx, span.Context(), deadline, createRes)
						}
						return 0, nil
					case resultCh <- &ScanResultItem{
						Result: Result{
							cas: Cas(item.Cas),
						},
						transcoder: m.Transcoder(),
						id:         string(item.Key),
						flags:      item.Flags,
						contents:   item.Value,
						expiryTime: expiry,
						keysOnly:   m.KeysOnly(),
					}:
					}
				}
				lastTermSeen = items[len(items)-1].Key
			}
			if isComplete {
				return 0, nil
			}
		}
		if atomic.LoadUint32(&m.cancelled) == 1 {
			m.cancelStream(ctx, span.Context(), deadline, createRes)
			return 0, nil
		}
	}
}

func (m *rangeScanOpManager) createStream(ctx context.Context, spanCtx RequestSpanContext, deadline time.Time, vbID uint16,
	rangeOpts *gocbcore.RangeScanCreateRangeScanConfig, samplingOpts *gocbcore.RangeScanCreateRandomSamplingConfig) (gocbcore.RangeScanCreateResult, error) {
	span := m.tracer.RequestSpan(spanCtx, "range_scan_create")
	defer span.End()
	span.SetAttribute("without_content", m.KeysOnly())
	if samplingOpts != nil {
		span.SetAttribute("scan_type", "sampling")
		span.SetAttribute("limit", samplingOpts.Samples)
		span.SetAttribute("seed", samplingOpts.Seed)
	} else if rangeOpts != nil {
		span.SetAttribute("scan_type", "range")
		span.SetAttribute("from_term", string(rangeOpts.Start))
		span.SetAttribute("to_term", string(rangeOpts.End))
		span.SetAttribute("from_exclusive", len(rangeOpts.ExclusiveStart) > 0)
		span.SetAttribute("to_exclusive", len(rangeOpts.ExclusiveEnd) > 0)
	}

	opMan := newAsyncOpManager(ctx)
	opMan.SetCancelCh(m.cancelCh)

	var createResOut gocbcore.RangeScanCreateResult
	var errOut error
	err := opMan.Wait(m.agent.RangeScanCreate(vbID, gocbcore.RangeScanCreateOptions{
		Deadline:     deadline,
		CollectionID: m.cid,
		KeysOnly:     m.KeysOnly(),
		Range:        rangeOpts,
		Sampling:     samplingOpts,
		Snapshot:     m.SnapshotOptions(vbID),
		User:         m.Impersonate(),
		TraceContext: span.Context(),
	}, func(result gocbcore.RangeScanCreateResult, err error) {
		if err != nil {
			errOut = err
			opMan.Reject()
			return
		}

		createResOut = result
		opMan.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return createResOut, errOut
}

func (m *rangeScanOpManager) continueStream(ctx context.Context, spanCtx RequestSpanContext, createRes gocbcore.RangeScanCreateResult) ([]gocbcore.RangeScanItem, bool, error) {
	span := m.tracer.RequestSpan(spanCtx, "range_scan_continue")
	defer span.End()

	span.SetAttribute("item_limit", m.itemLimit)
	span.SetAttribute("byte_limit", m.byteLimit)
	span.SetAttribute("time_limit", 0)

	opm := newAsyncOpManager(ctx)
	opm.SetCancelCh(m.cancelCh)

	var items []gocbcore.RangeScanItem
	span.SetAttribute("range_scan_id", "0x"+hex.EncodeToString(createRes.ScanUUID()))

	var itemsOut []gocbcore.RangeScanItem
	var completeOut bool
	var errOut error

	err := opm.Wait(createRes.RangeScanContinue(gocbcore.RangeScanContinueOptions{
		User:         m.Impersonate(),
		TraceContext: span.Context(),
		MaxCount:     m.itemLimit,
		MaxBytes:     m.byteLimit,
	}, func(coreItems []gocbcore.RangeScanItem) {
		items = append(items, coreItems...)
	}, func(result *gocbcore.RangeScanContinueResult, err error) {
		if err != nil {
			errOut = err
			opm.Reject()
			return
		}

		itemsOut = items
		if result.Complete {
			completeOut = true
			opm.Resolve()
			return
		}
		if result.More {
			opm.Resolve()
			return
		}
		logInfof("Received a range scan action that did not meet what we expected")
		opm.Resolve()
	}))
	if err != nil {
		errOut = err
	}

	return itemsOut, completeOut, errOut
}

func (m *rangeScanOpManager) cancelStream(ctx context.Context, spanCtx RequestSpanContext, deadline time.Time, createRes gocbcore.RangeScanCreateResult) {
	opMan := newAsyncOpManager(ctx)
	span := m.tracer.RequestSpan(spanCtx, "range_scan_cancel")
	defer span.End()

	span.SetAttribute("range_scan_id", "0x"+hex.EncodeToString(createRes.ScanUUID()))

	err := opMan.Wait(createRes.RangeScanCancel(gocbcore.RangeScanCancelOptions{
		Deadline:     deadline,
		User:         m.Impersonate(),
		TraceContext: span.Context(),
	}, func(result *gocbcore.RangeScanCancelResult, err error) {
		if err != nil {
			logDebugf("Failed to cancel scan 0x%s: %v", hex.EncodeToString(createRes.ScanUUID()), err)
			opMan.Reject()
			return
		}

		opMan.Resolve()
	}))
	if err != nil {
		return
	}
}

func (m *rangeScanOpManager) createLoadBalancer() *rangeScanLoadBalancer {
	var seed int64
	if m.SamplingOptions() != nil && m.SamplingOptions().Seed != 0 {
		// Using the sampling scan seed for the load balancer ensures that when concurrency is 1 the vbuckets are
		// always scanned in the same order for a given seed
		seed = int64(m.SamplingOptions().Seed)
	} else {
		seed = time.Now().UnixNano()
	}

	return newRangeScanLoadBalancer(m.serverToVbucketMap, seed)
}

type rangeScanLoadBalancer struct {
	vbucketChannels    map[int]chan uint16
	servers            []int
	activeScansPerNode sync.Map
	selectLock         sync.Mutex
}

func newRangeScanLoadBalancer(serverToVbucketMap map[int][]uint16, seed int64) *rangeScanLoadBalancer {
	b := &rangeScanLoadBalancer{
		vbucketChannels:    make(map[int]chan uint16),
		activeScansPerNode: sync.Map{},
	}

	for server, vbuckets := range serverToVbucketMap {
		b.servers = append(b.servers, server)

		b.vbucketChannels[server] = make(chan uint16, len(vbuckets))

		r := rand.New(rand.NewSource(seed)) // #nosec G404
		r.Shuffle(len(vbuckets), func(i, j int) {
			vbuckets[i], vbuckets[j] = vbuckets[j], vbuckets[i]
		})

		for _, vbucket := range vbuckets {
			b.vbucketChannels[server] <- vbucket
		}
	}

	return b
}

func (b *rangeScanLoadBalancer) retryScan(vbucket rangeScanVbucket) {
	b.vbucketChannels[vbucket.server] <- vbucket.id
}

func (b *rangeScanLoadBalancer) scanEnded(vbucket rangeScanVbucket) {
	zeroVal := uint32(0)
	val, _ := b.activeScansPerNode.LoadOrStore(vbucket.server, &zeroVal)
	atomic.AddUint32(val.(*uint32), ^uint32(0))
}

func (b *rangeScanLoadBalancer) scanStarting(vbucket rangeScanVbucket) {
	zeroVal := uint32(0)
	val, _ := b.activeScansPerNode.LoadOrStore(vbucket.server, &zeroVal)
	atomic.AddUint32(val.(*uint32), uint32(1))
}

// close closes all the vbucket channels. This should only be called if no more vbucket scans will happen, i.e. selectVbucket should not be called after close.
func (b *rangeScanLoadBalancer) close() {
	for _, ch := range b.vbucketChannels {
		close(ch)
	}
}

// selectVbucket returns the vbucket id, alongside the corresponding node index for a vbucket that is on the node with
// the smallest number of active scans. The boolean return value is false if there are no more vbuckets to scan.
func (b *rangeScanLoadBalancer) selectVbucket() (rangeScanVbucket, bool) {
	b.selectLock.Lock()
	defer b.selectLock.Unlock()

	var selectedServer int
	selected := false
	min := uint32(math.MaxUint32)

	for s := range b.servers {
		if len(b.vbucketChannels[s]) == 0 {
			continue
		}
		zeroVal := uint32(0)
		val, _ := b.activeScansPerNode.LoadOrStore(s, &zeroVal)
		activeScans := *val.(*uint32)
		if activeScans < min {
			min = activeScans
			selectedServer = s
			selected = true
		}
	}

	if !selected {
		return rangeScanVbucket{}, false
	}

	selectedVbucket, ok := <-b.vbucketChannels[selectedServer]
	if !ok {
		// This should be unreachable. selectVbucket should not be called after close.
		logWarnf("Vbucket channel has been closed before the range scan has finished")
		return rangeScanVbucket{}, false
	}
	vbucket := rangeScanVbucket{
		id:     selectedVbucket,
		server: selectedServer,
	}
	b.scanStarting(vbucket)
	return vbucket, true
}
