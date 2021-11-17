package gocbcore

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// PingState is the current state of a endpoint used in a PingResult.
type PingState uint32

const (
	// PingStateOK indicates that an endpoint is OK.
	PingStateOK PingState = 1

	// PingStateTimeout indicates that the ping request to an endpoint timed out.
	PingStateTimeout PingState = 2

	// PingStateError indicates that the ping request to an endpoint encountered an error.
	PingStateError PingState = 3
)

// EndpointState is the current connection state of an endpoint.
type EndpointState uint32

const (
	// EndpointStateDisconnected indicates that the endpoint is disconnected.
	EndpointStateDisconnected EndpointState = 1

	// EndpointStateConnecting indicates that the endpoint is connecting.
	EndpointStateConnecting EndpointState = 2

	// EndpointStateConnected indicates that the endpoint is connected.
	EndpointStateConnected EndpointState = 3

	// EndpointStateDisconnecting indicates that the endpoint is disconnecting.
	EndpointStateDisconnecting EndpointState = 4
)

// EndpointPingResult contains the results of a ping to a single server.
type EndpointPingResult struct {
	Endpoint string
	Error    error
	Latency  time.Duration
	ID       string
	Scope    string
	State    PingState
}

type pingSubOp struct {
	op       PendingOp
	endpoint string
}

type pingOp struct {
	lock       sync.Mutex
	subops     []pingSubOp
	remaining  int32
	results    map[ServiceType][]EndpointPingResult
	callback   PingCallback
	bucketName string
	httpCancel context.CancelFunc
}

func (pop *pingOp) Cancel() {
	for _, subop := range pop.subops {
		subop.op.Cancel()
	}
	pop.httpCancel()
}

func (pop *pingOp) handledOneLocked(configRev int64) {
	remaining := atomic.AddInt32(&pop.remaining, -1)
	if remaining == 0 {
		pop.httpCancel()
		pop.callback(&PingResult{
			ConfigRev: configRev,
			Services:  pop.results,
		}, nil)
	}
}

// PingOptions encapsulates the parameters for a PingKv operation.
type PingOptions struct {
	TraceContext RequestSpanContext
	KVDeadline   time.Time
	CbasDeadline time.Time
	N1QLDeadline time.Time
	FtsDeadline  time.Time
	CapiDeadline time.Time
	MgmtDeadline time.Time
	ServiceTypes []ServiceType

	// Internal: This should never be used and is not supported.
	User string

	ignoreMissingServices bool
}

// PingResult encapsulates the result of a PingKv operation.
type PingResult struct {
	ConfigRev int64
	Services  map[ServiceType][]EndpointPingResult
}

// DiagnosticsOptions encapsulates the parameters for a Diagnostics operation.
type DiagnosticsOptions struct {
}

// MemdConnInfo represents information we know about a particular
// memcached connection reported in a diagnostics report.
type MemdConnInfo struct {
	LocalAddr    string
	RemoteAddr   string
	LastActivity time.Time
	Scope        string
	ID           string
	State        EndpointState
}

// DiagnosticInfo is returned by the Diagnostics method and includes
// information about the overall health of the clients connections.
type DiagnosticInfo struct {
	ConfigRev int64
	MemdConns []MemdConnInfo
	State     ClusterState
}

// ClusterState is used to describe the state of a cluster.
type ClusterState uint32

const (
	// ClusterStateOnline specifies that all nodes and their sockets are reachable.
	ClusterStateOnline = ClusterState(1)

	// ClusterStateDegraded specifies that at least one socket per service is reachable.
	ClusterStateDegraded = ClusterState(2)

	// ClusterStateOffline is used to specify that not even one socker per service is reachable.
	ClusterStateOffline = ClusterState(3)
)

type waitUntilOp struct {
	lock       sync.Mutex
	remaining  int32
	callback   WaitUntilReadyCallback
	stopCh     chan struct{}
	timer      *time.Timer
	httpCancel context.CancelFunc
	closed     bool

	retryLock    sync.Mutex
	retries      uint32
	retryReasons []RetryReason
	retryStrat   RetryStrategy
}

func (wuo *waitUntilOp) RetryAttempts() uint32 {
	return atomic.LoadUint32(&wuo.retries)
}

func (wuo *waitUntilOp) RetryReasons() []RetryReason {
	wuo.retryLock.Lock()
	defer wuo.retryLock.Unlock()
	return wuo.retryReasons
}

func (wuo *waitUntilOp) Identifier() string {
	return "waituntilready"
}

func (wuo *waitUntilOp) Idempotent() bool {
	return true
}

func (wuo *waitUntilOp) retryStrategy() RetryStrategy {
	return wuo.retryStrat
}

func (wuo *waitUntilOp) recordRetryAttempt(reason RetryReason) {
	atomic.AddUint32(&wuo.retries, 1)
	wuo.retryLock.Lock()
	defer wuo.retryLock.Unlock()
	idx := sort.Search(len(wuo.retryReasons), func(i int) bool {
		return wuo.retryReasons[i] == reason
	})

	// if idx is out of the range of retryReasons then it wasn't found.
	if idx > len(wuo.retryReasons)-1 {
		wuo.retryReasons = append(wuo.retryReasons, reason)
	}
}

func (wuo *waitUntilOp) cancel(err error) {
	wuo.lock.Lock()
	wuo.timer.Stop()
	if wuo.closed {
		wuo.lock.Unlock()
		return
	}
	wuo.closed = true
	wuo.lock.Unlock()
	close(wuo.stopCh)
	wuo.httpCancel()
	wuo.callback(nil, err)
}

func (wuo *waitUntilOp) Cancel() {
	wuo.cancel(errRequestCanceled)
}

func (wuo *waitUntilOp) handledOneLocked() {
	remaining := atomic.AddInt32(&wuo.remaining, -1)
	if remaining == 0 {
		wuo.timer.Stop()
		wuo.httpCancel()
		wuo.callback(&WaitUntilReadyResult{}, nil)
	}
}

// WaitUntilReadyResult encapsulates the result of a WaitUntilReady operation.
type WaitUntilReadyResult struct {
}

// WaitUntilReadyOptions encapsulates the parameters for a WaitUntilReady operation.
type WaitUntilReadyOptions struct {
	DesiredState ClusterState  // Defaults to ClusterStateOnline
	ServiceTypes []ServiceType // Defaults to all services
	// If the cluster state is offline and a connect error has been observed then fast fail and return it.
	RetryStrategy RetryStrategy
}
