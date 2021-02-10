package gocb

import (
	"encoding/json"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type thresholdLogGroup struct {
	name  string
	floor time.Duration
	ops   []*thresholdLogSpan
	lock  sync.RWMutex
}

func (g *thresholdLogGroup) init(name string, floor time.Duration, size uint32) {
	g.name = name
	g.floor = floor
	g.ops = make([]*thresholdLogSpan, 0, size)
}

func (g *thresholdLogGroup) recordOp(span *thresholdLogSpan) {
	if span.duration < g.floor {
		return
	}

	// Preemptively check that we actually need to be inserted using a read lock first
	// this is a performance improvement measure to avoid locking the mutex all the time.
	g.lock.RLock()
	if len(g.ops) == cap(g.ops) && span.duration < g.ops[0].duration {
		// we are at capacity and we are faster than the fastest slow op
		g.lock.RUnlock()
		return
	}
	g.lock.RUnlock()

	g.lock.Lock()
	if len(g.ops) == cap(g.ops) && span.duration < g.ops[0].duration {
		// we are at capacity and we are faster than the fastest slow op
		g.lock.Unlock()
		return
	}

	l := len(g.ops)
	i := sort.Search(l, func(i int) bool { return span.duration < g.ops[i].duration })

	// i represents the slot where it should be inserted

	if len(g.ops) < cap(g.ops) {
		if i == l {
			g.ops = append(g.ops, span)
		} else {
			g.ops = append(g.ops, nil)
			copy(g.ops[i+1:], g.ops[i:])
			g.ops[i] = span
		}
	} else {
		if i == 0 {
			g.ops[i] = span
		} else {
			copy(g.ops[0:i-1], g.ops[1:i])
			g.ops[i-1] = span
		}
	}

	g.lock.Unlock()
}

type thresholdLogItem struct {
	OperationName          string `json:"operation_name,omitempty"`
	TotalTimeUs            uint64 `json:"total_us,omitempty"`
	EncodeDurationUs       uint64 `json:"encode_us,omitempty"`
	DispatchDurationUs     uint64 `json:"dispatch_us,omitempty"`
	ServerDurationUs       uint64 `json:"server_us,omitempty"`
	LastRemoteAddress      string `json:"last_remote_address,omitempty"`
	LastLocalAddress       string `json:"last_local_address,omitempty"`
	LastDispatchDurationUs uint64 `json:"last_dispatch_us,omitempty"`
	LastOperationID        string `json:"last_operation_id,omitempty"`
	LastLocalID            string `json:"last_local_id,omitempty"`
	DocumentKey            string `json:"document_key,omitempty"`
}

type thresholdLogService struct {
	Service string             `json:"service"`
	Count   uint64             `json:"count"`
	Top     []thresholdLogItem `json:"top"`
}

func (g *thresholdLogGroup) logRecordedRecords(sampleSize uint32) {
	// Preallocate space to copy the ops into...
	oldOps := make([]*thresholdLogSpan, sampleSize)

	g.lock.Lock()
	// Escape early if we have no ops to log...
	if len(g.ops) == 0 {
		g.lock.Unlock()
		return
	}

	// Copy out our ops so we can cheaply print them out without blocking
	// our ops from actually being recorded in other goroutines (which would
	// effectively slow down the op pipeline for logging).

	oldOps = oldOps[0:len(g.ops)]
	copy(oldOps, g.ops)
	g.ops = g.ops[:0]

	g.lock.Unlock()

	jsonData := thresholdLogService{
		Service: g.name,
	}

	for i := len(oldOps) - 1; i >= 0; i-- {
		op := oldOps[i]

		jsonData.Top = append(jsonData.Top, thresholdLogItem{
			OperationName:          op.opName,
			TotalTimeUs:            uint64(op.duration / time.Microsecond),
			DispatchDurationUs:     uint64(op.totalDispatchDuration / time.Microsecond),
			ServerDurationUs:       uint64(op.totalServerDuration / time.Microsecond),
			EncodeDurationUs:       uint64(op.totalEncodeDuration / time.Microsecond),
			LastRemoteAddress:      op.lastDispatchPeer,
			LastDispatchDurationUs: uint64(op.lastDispatchDuration / time.Microsecond),
			LastOperationID:        op.lastOperationID,
			LastLocalID:            op.lastLocalID,
			DocumentKey:            op.documentKey,
		})
	}

	jsonData.Count = uint64(len(jsonData.Top))

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		logDebugf("Failed to generate threshold logging service JSON: %s", err)
	}

	logInfof("Threshold Log: %s", jsonBytes)
}

// ThresholdLoggingOptions is the set of options available for configuring threshold logging.
type ThresholdLoggingOptions struct {
	ServerDurationDisabled bool
	Interval               time.Duration
	SampleSize             uint32
	KVThreshold            time.Duration
	ViewsThreshold         time.Duration
	QueryThreshold         time.Duration
	SearchThreshold        time.Duration
	AnalyticsThreshold     time.Duration
	ManagementThreshold    time.Duration
}

// ThresholdLoggingTracer is a specialized Tracer implementation which will automatically
// log operations which fall outside of a set of thresholds.  Note that this tracer is
// only safe for use within the Couchbase SDK, uses by external event sources are
// likely to fail.
type ThresholdLoggingTracer struct {
	Interval            time.Duration
	SampleSize          uint32
	KVThreshold         time.Duration
	ViewsThreshold      time.Duration
	QueryThreshold      time.Duration
	SearchThreshold     time.Duration
	AnalyticsThreshold  time.Duration
	ManagementThreshold time.Duration

	killCh          chan struct{}
	refCount        int32
	nextTick        time.Time
	kvGroup         thresholdLogGroup
	viewsGroup      thresholdLogGroup
	queryGroup      thresholdLogGroup
	searchGroup     thresholdLogGroup
	analyticsGroup  thresholdLogGroup
	managementGroup thresholdLogGroup
}

func NewThresholdLoggingTracer(opts *ThresholdLoggingOptions) *ThresholdLoggingTracer {
	if opts == nil {
		opts = &ThresholdLoggingOptions{}
	}
	if opts.Interval == 0 {
		opts.Interval = 10 * time.Second
	}
	if opts.SampleSize == 0 {
		opts.SampleSize = 10
	}
	if opts.KVThreshold == 0 {
		opts.KVThreshold = 500 * time.Millisecond
	}
	if opts.ViewsThreshold == 0 {
		opts.ViewsThreshold = 1 * time.Second
	}
	if opts.QueryThreshold == 0 {
		opts.QueryThreshold = 1 * time.Second
	}
	if opts.SearchThreshold == 0 {
		opts.SearchThreshold = 1 * time.Second
	}
	if opts.AnalyticsThreshold == 0 {
		opts.AnalyticsThreshold = 1 * time.Second
	}
	if opts.ManagementThreshold == 0 {
		opts.ManagementThreshold = 1 * time.Second
	}

	t := &ThresholdLoggingTracer{
		Interval:            opts.Interval,
		SampleSize:          opts.SampleSize,
		KVThreshold:         opts.KVThreshold,
		ViewsThreshold:      opts.ViewsThreshold,
		QueryThreshold:      opts.QueryThreshold,
		SearchThreshold:     opts.SearchThreshold,
		AnalyticsThreshold:  opts.AnalyticsThreshold,
		ManagementThreshold: opts.ManagementThreshold,
	}

	t.kvGroup.init("kv", t.KVThreshold, t.SampleSize)
	t.viewsGroup.init("views", t.ViewsThreshold, t.SampleSize)
	t.queryGroup.init("query", t.QueryThreshold, t.SampleSize)
	t.searchGroup.init("search", t.SearchThreshold, t.SampleSize)
	t.analyticsGroup.init("analytics", t.AnalyticsThreshold, t.SampleSize)
	t.managementGroup.init("management", t.ManagementThreshold, t.SampleSize)

	if t.killCh == nil {
		t.killCh = make(chan struct{})
	}

	if t.nextTick.IsZero() {
		t.nextTick = time.Now().Add(t.Interval)
	}

	return t
}

// AddRef is used internally to keep track of the number of Cluster instances referring to it.
// This is used to correctly shut down the aggregation routines once there are no longer any
// instances tracing to it.
func (t *ThresholdLoggingTracer) AddRef() int32 {
	newRefCount := atomic.AddInt32(&t.refCount, 1)
	if newRefCount == 1 {
		t.startLoggerRoutine()
	}
	return newRefCount
}

// DecRef is the counterpart to AddRef (see AddRef for more information).
func (t *ThresholdLoggingTracer) DecRef() int32 {
	newRefCount := atomic.AddInt32(&t.refCount, -1)
	if newRefCount == 0 {
		t.killCh <- struct{}{}
	}
	return newRefCount
}

func (t *ThresholdLoggingTracer) logRecordedRecords() {
	t.kvGroup.logRecordedRecords(t.SampleSize)
	t.viewsGroup.logRecordedRecords(t.SampleSize)
	t.queryGroup.logRecordedRecords(t.SampleSize)
	t.searchGroup.logRecordedRecords(t.SampleSize)
	t.analyticsGroup.logRecordedRecords(t.SampleSize)
	t.managementGroup.logRecordedRecords(t.SampleSize)
}

func (t *ThresholdLoggingTracer) startLoggerRoutine() {
	go t.loggerRoutine()
}

func (t *ThresholdLoggingTracer) loggerRoutine() {
	for {
		select {
		case <-time.After(time.Until(t.nextTick)):
			t.nextTick = t.nextTick.Add(t.Interval)
			t.logRecordedRecords()
		case <-t.killCh:
			t.logRecordedRecords()
			return
		}
	}
}

func (t *ThresholdLoggingTracer) recordOp(span *thresholdLogSpan) {
	switch span.serviceName {
	case "mgmt":
		t.managementGroup.recordOp(span)
	case "kv":
		t.kvGroup.recordOp(span)
	case "views":
		t.viewsGroup.recordOp(span)
	case "query":
		t.queryGroup.recordOp(span)
	case "search":
		t.searchGroup.recordOp(span)
	case "analytics":
		t.analyticsGroup.recordOp(span)
	}
}

// StartSpan belongs to the Tracer interface.
func (t *ThresholdLoggingTracer) StartSpan(operationName string, parentContext requestSpanContext) requestSpan {
	span := &thresholdLogSpan{
		tracer:    t,
		opName:    operationName,
		startTime: time.Now(),
	}

	if context, ok := parentContext.(*thresholdLogSpanContext); ok {
		span.parent = context.span
	}

	return span
}

type thresholdLogSpan struct {
	tracer                *ThresholdLoggingTracer
	parent                *thresholdLogSpan
	opName                string
	startTime             time.Time
	serviceName           string
	peerAddress           string
	serverDuration        time.Duration
	duration              time.Duration
	totalServerDuration   time.Duration
	totalDispatchDuration time.Duration
	totalEncodeDuration   time.Duration
	lastDispatchPeer      string
	lastDispatchDuration  time.Duration
	lastOperationID       string
	lastLocalID           string
	documentKey           string
	lock                  sync.Mutex
}

func (n *thresholdLogSpan) Context() requestSpanContext {
	return &thresholdLogSpanContext{n}
}

func (n *thresholdLogSpan) SetTag(key string, value interface{}) requestSpan {
	var ok bool

	switch key {
	case "server_duration":
		if n.serverDuration, ok = value.(time.Duration); !ok {
			logDebugf("Failed to cast span server_duration tag")
		}
	case "couchbase.service":
		if n.serviceName, ok = value.(string); !ok {
			logDebugf("Failed to cast span couchbase.service tag")
		}
	case "peer.address":
		if n.peerAddress, ok = value.(string); !ok {
			logDebugf("Failed to cast span peer.address tag")
		}
	case "couchbase.operation_id":
		if n.lastOperationID, ok = value.(string); !ok {
			logDebugf("Failed to cast span couchbase.operation_id tag")
		}
	case "couchbase.document_key":
		if n.documentKey, ok = value.(string); !ok {
			logDebugf("Failed to cast span couchbase.document_key tag")
		}
	case "couchbase.local_id":
		if n.lastLocalID, ok = value.(string); !ok {
			logDebugf("Failed to cast span couchbase.local_id tag")
		}
	}
	return n
}

func (n *thresholdLogSpan) Finish() {
	n.duration = time.Since(n.startTime)

	n.totalServerDuration += n.serverDuration
	if n.opName == "dispatch" {
		n.totalDispatchDuration += n.duration
		n.lastDispatchPeer = n.peerAddress
		n.lastDispatchDuration = n.duration
	}
	if n.opName == "encode" {
		n.totalEncodeDuration += n.duration
	}

	if n.parent != nil {
		n.parent.lock.Lock()
		n.parent.totalServerDuration += n.totalServerDuration
		n.parent.totalDispatchDuration += n.totalDispatchDuration
		n.parent.totalEncodeDuration += n.totalEncodeDuration
		if n.lastDispatchPeer != "" || n.lastDispatchDuration > 0 {
			n.parent.lastDispatchPeer = n.lastDispatchPeer
			n.parent.lastDispatchDuration = n.lastDispatchDuration
		}
		if n.lastOperationID != "" {
			n.parent.lastOperationID = n.lastOperationID
		}
		if n.lastLocalID != "" {
			n.parent.lastLocalID = n.lastLocalID
		}
		if n.documentKey != "" {
			n.parent.documentKey = n.documentKey
		}
		n.parent.lock.Unlock()
	}

	if n.serviceName != "" {
		n.tracer.recordOp(n)
	}
}

type thresholdLogSpanContext struct {
	span *thresholdLogSpan
}
