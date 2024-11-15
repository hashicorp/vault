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

func initThresholdLogGroup(name string, floor time.Duration, size uint32) *thresholdLogGroup {
	return &thresholdLogGroup{
		name:  name,
		floor: floor,
		ops:   make([]*thresholdLogSpan, 0, size),
	}
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
	TotalTimeUs            uint64 `json:"total_duration_us,omitempty"`
	EncodeDurationUs       uint64 `json:"encode_duration_us,omitempty"`
	DispatchDurationUs     uint64 `json:"total_dispatch_duration_us,omitempty"`
	ServerDurationUs       uint64 `json:"total_server_duration_us,omitempty"`
	LastRemoteAddress      string `json:"last_remote_socket,omitempty"`
	LastLocalAddress       string `json:"last_local_socket,omitempty"`
	LastDispatchDurationUs uint64 `json:"last_dispatch_duration_us,omitempty"`
	LastServerDurationUs   uint64 `json:"last_server_duration_us,omitempty"`
	LastOperationID        string `json:"operation_id,omitempty"`
	LastLocalID            string `json:"last_local_id,omitempty"`
}

type thresholdLogEntry struct {
	Count uint64             `json:"total_count"`
	Top   []thresholdLogItem `json:"top_requests"`
}

type thresholdLogService map[string]thresholdLogEntry

// ThresholdLoggingOptions is the set of options available for configuring threshold logging.
type ThresholdLoggingOptions struct {
	Interval            time.Duration
	SampleSize          uint32
	KVThreshold         time.Duration
	KVScanThreshold     time.Duration
	ViewsThreshold      time.Duration
	QueryThreshold      time.Duration
	SearchThreshold     time.Duration
	AnalyticsThreshold  time.Duration
	ManagementThreshold time.Duration
}

// ThresholdLoggingTracer is a specialized Tracer implementation which will automatically
// log operations which fall outside of a set of thresholds.  Note that this tracer is
// only safe for use within the Couchbase SDK, uses by external event sources are
// likely to fail.
type ThresholdLoggingTracer struct {
	Interval            time.Duration
	SampleSize          uint32
	KVThreshold         time.Duration
	KVScanThreshold     time.Duration
	ViewsThreshold      time.Duration
	QueryThreshold      time.Duration
	SearchThreshold     time.Duration
	AnalyticsThreshold  time.Duration
	ManagementThreshold time.Duration

	killCh   chan struct{}
	refCount int32
	nextTick time.Time
	groups   map[string]*thresholdLogGroup
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
	if opts.KVScanThreshold == 0 {
		opts.KVScanThreshold = 1 * time.Second
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
		KVScanThreshold:     opts.KVScanThreshold,
		ViewsThreshold:      opts.ViewsThreshold,
		QueryThreshold:      opts.QueryThreshold,
		SearchThreshold:     opts.SearchThreshold,
		AnalyticsThreshold:  opts.AnalyticsThreshold,
		ManagementThreshold: opts.ManagementThreshold,
	}

	t.groups = map[string]*thresholdLogGroup{
		"kv":         initThresholdLogGroup("kv", t.KVThreshold, t.SampleSize),
		"kv_scan":    initThresholdLogGroup("kv_scan", t.KVScanThreshold, t.SampleSize),
		"views":      initThresholdLogGroup("views", t.ViewsThreshold, t.SampleSize),
		"query":      initThresholdLogGroup("query", t.QueryThreshold, t.SampleSize),
		"search":     initThresholdLogGroup("search", t.SearchThreshold, t.SampleSize),
		"analytics":  initThresholdLogGroup("analytics", t.AnalyticsThreshold, t.SampleSize),
		"management": initThresholdLogGroup("management", t.ManagementThreshold, t.SampleSize),
	}

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

func (t *ThresholdLoggingTracer) buildJSONData() thresholdLogService {
	// Preallocate space to copy the ops into...
	oldOps := make([]*thresholdLogSpan, t.SampleSize)

	jsonData := make(thresholdLogService)
	for _, g := range t.groups {
		g.lock.Lock()
		// Escape early if we have no ops to log...
		if len(g.ops) == 0 {
			g.lock.Unlock()
			continue
		}

		// Copy out our ops so we can cheaply print them out without blocking
		// our ops from actually being recorded in other goroutines (which would
		// effectively slow down the op pipeline for logging).

		oldOps = oldOps[0:len(g.ops)]
		copy(oldOps, g.ops)
		g.ops = g.ops[:0]

		g.lock.Unlock()

		entry := thresholdLogEntry{}

		for i := len(oldOps) - 1; i >= 0; i-- {
			op := oldOps[i]

			localAddr := op.lastDispatchLocal
			if localAddr != "" && op.lastDispatchLocalPort != "" {
				localAddr = localAddr + ":" + op.lastDispatchLocalPort
			}

			peerAddr := op.lastDispatchPeer
			if peerAddr != "" && op.lastDispatchPeerPort != "" {
				peerAddr = peerAddr + ":" + op.lastDispatchPeerPort
			}

			entry.Top = append(entry.Top, thresholdLogItem{
				OperationName:          op.opName,
				TotalTimeUs:            uint64(op.duration / time.Microsecond),
				DispatchDurationUs:     uint64(op.totalDispatchDuration / time.Microsecond),
				ServerDurationUs:       uint64(op.totalServerDuration / time.Microsecond),
				EncodeDurationUs:       uint64(op.totalEncodeDuration / time.Microsecond),
				LastLocalAddress:       localAddr,
				LastRemoteAddress:      peerAddr,
				LastDispatchDurationUs: uint64(op.lastDispatchDuration / time.Microsecond),
				LastServerDurationUs:   uint64(op.lastServerDuration / time.Microsecond),
				LastOperationID:        op.lastOperationID,
				LastLocalID:            op.lastLocalID,
			})
		}

		entry.Count = uint64(len(entry.Top))

		jsonData[g.name] = entry
	}

	return jsonData
}

func (t *ThresholdLoggingTracer) logRecordedRecords() {
	jsonData := t.buildJSONData()

	if len(jsonData) == 0 {
		// Nothing to log so make sure we don't just log empty objects.
		return
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		logDebugf("Failed to generate threshold logging service JSON: %s", err)
	}

	logInfof("Threshold Log: %s", jsonBytes)
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
		t.groups["management"].recordOp(span)
	case "kv":
		t.groups["kv"].recordOp(span)
	case "kv_scan":
		t.groups["kv_scan"].recordOp(span)
	case "views":
		t.groups["views"].recordOp(span)
	case "query":
		t.groups["query"].recordOp(span)
	case "search":
		t.groups["search"].recordOp(span)
	case "analytics":
		t.groups["analytics"].recordOp(span)
	}
}

// RequestSpan belongs to the Tracer interface.
func (t *ThresholdLoggingTracer) RequestSpan(parentContext RequestSpanContext, operationName string) RequestSpan {
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
	localAddress          string
	peerPort              string
	localPort             string
	serverDuration        time.Duration
	duration              time.Duration
	totalServerDuration   time.Duration
	totalDispatchDuration time.Duration
	totalEncodeDuration   time.Duration
	lastDispatchPeer      string
	lastDispatchLocal     string
	lastDispatchPeerPort  string
	lastDispatchLocalPort string
	lastDispatchDuration  time.Duration
	lastServerDuration    time.Duration
	lastOperationID       string
	lastLocalID           string
	lock                  sync.Mutex
}

func (n *thresholdLogSpan) Context() RequestSpanContext {
	return &thresholdLogSpanContext{n}
}

func (n *thresholdLogSpan) SetAttribute(key string, value interface{}) {
	var ok bool

	switch key {
	case spanAttribServerDurationKey:
		if n.serverDuration, ok = value.(time.Duration); !ok {
			logDebugf("Failed to cast span db.couchbase.server_duration tag")
		}
	case spanAttribServiceKey:
		if n.serviceName, ok = value.(string); !ok {
			logDebugf("Failed to cast span db.couchbase.service tag")
		}
	case spanAttribNetPeerNameKey:
		if n.peerAddress, ok = value.(string); !ok {
			logDebugf("Failed to cast span net.peer.name tag")
		}
	case spanAttribNetHostNameKey:
		if n.localAddress, ok = value.(string); !ok {
			logDebugf("Failed to cast span net.host.local tag")
		}
	case spanAttribOperationIDKey:
		if n.lastOperationID, ok = value.(string); !ok {
			logDebugf("Failed to cast span db.couchbase.operation_id tag")
		}
	case spanAttribLocalIDKey:
		if n.lastLocalID, ok = value.(string); !ok {
			logDebugf("Failed to cast span db.couchbase.local_id tag")
		}
	case spanAttribNetPeerPortKey:
		if n.peerPort, ok = value.(string); !ok {
			logDebugf("Failed to cast span net.peer.port tag")
		}
	case spanAttribNetHostPortKey:
		if n.localPort, ok = value.(string); !ok {
			logDebugf("Failed to cast span net.host.port tag")
		}
	}
}

func (n *thresholdLogSpan) AddEvent(key string, timestamp time.Time) {
}

func (n *thresholdLogSpan) End() {
	n.duration = time.Since(n.startTime)

	n.totalServerDuration += n.serverDuration
	if n.opName == spanNameDispatchToServer {
		n.totalDispatchDuration += n.duration
		n.lastDispatchPeer = n.peerAddress
		n.lastDispatchLocal = n.localAddress
		n.lastDispatchPeerPort = n.peerPort
		n.lastDispatchLocalPort = n.localPort
		n.lastDispatchDuration = n.duration
		n.lastServerDuration = n.serverDuration
	}
	if n.opName == spanNameRequestEncoding {
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
		if n.lastDispatchPeer != "" || n.lastServerDuration > 0 {
			n.parent.lastServerDuration = n.lastServerDuration
		}
		if n.lastDispatchLocal != "" {
			n.parent.lastDispatchLocal = n.lastDispatchLocal
		}
		if n.lastOperationID != "" {
			n.parent.lastOperationID = n.lastOperationID
		}
		if n.lastLocalID != "" {
			n.parent.lastLocalID = n.lastLocalID
		}
		if n.lastDispatchLocalPort != "" {
			n.parent.lastDispatchLocalPort = n.lastDispatchLocalPort
		}
		if n.lastDispatchPeerPort != "" {
			n.parent.lastDispatchPeerPort = n.lastDispatchPeerPort
		}
		if n.lastDispatchPeerPort != "" {
			n.parent.lastDispatchPeerPort = n.lastDispatchPeerPort
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
