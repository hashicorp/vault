package gocbcore

import (
	"sync"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type statsComponent struct {
	kvMux                *kvMux
	tracer               *tracerComponent
	defaultRetryStrategy RetryStrategy
}

func newStatsComponent(kvMux *kvMux, defaultRetry RetryStrategy, tracer *tracerComponent) *statsComponent {
	return &statsComponent{
		kvMux:                kvMux,
		tracer:               tracer,
		defaultRetryStrategy: defaultRetry,
	}
}

func (sc *statsComponent) Stats(opts StatsOptions, cb StatsCallback) (PendingOp, error) {
	tracer := sc.tracer.StartTelemeteryHandler(metricValueServiceKeyValue, "Stats", opts.TraceContext)

	iter, err := sc.kvMux.PipelineSnapshot()
	if err != nil {
		tracer.Finish()
		return nil, err
	}

	stats := make(map[string]SingleServerStats)
	var statsLock sync.Mutex

	op := new(multiPendingOp)
	op.isIdempotent = true
	var expected uint32

	pipelines := make([]*memdPipeline, 0)

	switch target := opts.Target.(type) {
	case nil:
		iter.Iterate(0, func(pipeline *memdPipeline) bool {
			pipelines = append(pipelines, pipeline)
			expected++
			return false
		})
	case VBucketIDStatsTarget:
		expected = 1

		srvIdx, err := iter.NodeByVbucket(target.VbID, 0)
		if err != nil {
			return nil, err
		}

		pipelines = append(pipelines, iter.PipelineAt(srvIdx))
	default:
		return nil, errInvalidArgument
	}

	opHandledLocked := func() {
		completed := op.IncrementCompletedOps()
		if expected-completed == 0 {
			tracer.Finish()
			cb(&StatsResult{
				Servers: stats,
			}, nil)
		}
	}

	var userFrame *memd.UserImpersonationFrame
	if len(opts.User) > 0 {
		userFrame = &memd.UserImpersonationFrame{
			User: []byte(opts.User),
		}
	}

	if opts.RetryStrategy == nil {
		opts.RetryStrategy = sc.defaultRetryStrategy
	}

	for _, pipeline := range pipelines {
		serverAddress := pipeline.Address()

		handler := func(resp *memdQResponse, req *memdQRequest, err error) {
			statsLock.Lock()
			defer statsLock.Unlock()

			// Fetch the specific stats key for this server.  Creating a new entry
			// for the server if we did not previously have one.
			curStats, ok := stats[serverAddress]
			if !ok {
				stats[serverAddress] = SingleServerStats{
					Stats: make(map[string]string),
				}
				curStats = stats[serverAddress]
			}

			if err != nil {
				// Store the first (and hopefully only) error into the Error field of this
				// server's stats entry.
				if curStats.Error == nil {
					curStats.Error = err
				} else {
					logDebugf("Got additional error for stats: %s: %v", serverAddress, err)
				}

				opHandledLocked()

				return
			}

			// Check if the key and value length is zero.  This indicates that we have reached
			// the ending of the stats listing by this server.
			if len(resp.Key) == 0 && len(resp.Value) == 0 {
				// As this is a persistent request, we must manually cancel it to remove
				// it from the pending ops list.  To ensure we do not race multiple cancels,
				// we only handle it as completed the one time cancellation succeeds.
				if req.internalCancel(err) {
					opHandledLocked()
				}

				return
			}

			curStats.StatsKeys = append(curStats.StatsKeys, resp.Key)
			curStats.StatsChunks = append(curStats.StatsChunks, resp.Value)
			if len(resp.Key) == 0 {
				// We do this for the sake of consistency.
				curStats.Stats[""] += string(resp.Value)
			} else {
				// Add the stat for this server to the list of stats.
				curStats.Stats[string(resp.Key)] += string(resp.Value)
			}
			// If we don't reassign this then we lose any values added to StatsKeys and StatsChunks.
			stats[serverAddress] = curStats
		}

		req := &memdQRequest{
			Packet: memd.Packet{
				Magic:                  memd.CmdMagicReq,
				Command:                memd.CmdStat,
				Datatype:               0,
				Cas:                    0,
				Key:                    []byte(opts.Key),
				Value:                  nil,
				UserImpersonationFrame: userFrame,
			},
			Persistent:       true,
			Callback:         handler,
			RootTraceContext: tracer.RootContext(),
			RetryStrategy:    opts.RetryStrategy,
		}

		curOp, err := sc.kvMux.DispatchDirectToAddress(req, pipeline.Address())
		if err != nil {
			statsLock.Lock()
			stats[serverAddress] = SingleServerStats{
				Error: err,
			}
			opHandledLocked()
			statsLock.Unlock()

			continue
		}

		if !opts.Deadline.IsZero() {
			start := time.Now()
			req.SetTimer(time.AfterFunc(opts.Deadline.Sub(start), func() {
				connInfo := req.ConnectionInfo()
				count, reasons := req.Retries()
				req.cancelWithCallbackAndFinishTracer(&TimeoutError{
					InnerError:         errAmbiguousTimeout,
					OperationID:        "Unlock",
					Opaque:             req.Identifier(),
					TimeObserved:       time.Since(start),
					RetryReasons:       reasons,
					RetryAttempts:      count,
					LastDispatchedTo:   connInfo.lastDispatchedTo,
					LastDispatchedFrom: connInfo.lastDispatchedFrom,
					LastConnectionID:   connInfo.lastConnectionID,
				}, tracer)
			}))
		}

		op.ops = append(op.ops, curOp)
	}

	return op, nil
}

// SingleServerStats represents the stats returned from a single server.
type SingleServerStats struct {
	Stats map[string]string
	// StatsKeys and StatsChunks provide access to the raw keys and values returned on a per packet basis.
	// This is useful for stats keys such as connections which, unlike most stats keys, return us a complex object
	// per packet. Keys and chunks maintain the same ordering for indexes.
	StatsKeys   [][]byte
	StatsChunks [][]byte
	Error       error
}

// StatsTarget is used for providing a specific target to the Stats operation.
type StatsTarget interface {
}

// VBucketIDStatsTarget indicates that a specific vbucket should be targeted by the Stats operation.
type VBucketIDStatsTarget struct {
	VbID uint16
}

// StatsOptions encapsulates the parameters for a Stats operation.
type StatsOptions struct {
	Key string
	// Target indicates that something specific should be targeted by the operation. If left nil
	// then the stats command will be sent to all servers.
	Target        StatsTarget
	RetryStrategy RetryStrategy
	Deadline      time.Time

	// Internal: This should never be used and is not supported.
	User string

	TraceContext RequestSpanContext
}

// StatsResult encapsulates the result of a Stats operation.
type StatsResult struct {
	Servers map[string]SingleServerStats
}
