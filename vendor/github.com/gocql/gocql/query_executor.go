package gocql

import (
	"context"
	"sync"
	"time"
)

type ExecutableQuery interface {
	execute(ctx context.Context, conn *Conn) *Iter
	attempt(keyspace string, end, start time.Time, iter *Iter, host *HostInfo)
	retryPolicy() RetryPolicy
	speculativeExecutionPolicy() SpeculativeExecutionPolicy
	GetRoutingKey() ([]byte, error)
	Keyspace() string
	IsIdempotent() bool

	withContext(context.Context) ExecutableQuery

	RetryableQuery
}

type queryExecutor struct {
	pool   *policyConnPool
	policy HostSelectionPolicy
}

func (q *queryExecutor) attemptQuery(ctx context.Context, qry ExecutableQuery, conn *Conn) *Iter {
	start := time.Now()
	iter := qry.execute(ctx, conn)
	end := time.Now()

	qry.attempt(q.pool.keyspace, end, start, iter, conn.host)

	return iter
}

func (q *queryExecutor) speculate(ctx context.Context, qry ExecutableQuery, sp SpeculativeExecutionPolicy,
	hostIter NextHost, results chan *Iter) *Iter {
	ticker := time.NewTicker(sp.Delay())
	defer ticker.Stop()

	for i := 0; i < sp.Attempts(); i++ {
		select {
		case <-ticker.C:
			go q.run(ctx, qry, hostIter, results)
		case <-ctx.Done():
			return &Iter{err: ctx.Err()}
		case iter := <-results:
			return iter
		}
	}

	return nil
}

func (q *queryExecutor) executeQuery(qry ExecutableQuery) (*Iter, error) {
	hostIter := q.policy.Pick(qry)

	// check if the query is not marked as idempotent, if
	// it is, we force the policy to NonSpeculative
	sp := qry.speculativeExecutionPolicy()
	if !qry.IsIdempotent() || sp.Attempts() == 0 {
		return q.do(qry.Context(), qry, hostIter), nil
	}

	// When speculative execution is enabled, we could be accessing the host iterator from multiple goroutines below.
	// To ensure we don't call it concurrently, we wrap the returned NextHost function here to synchronize access to it.
	var mu sync.Mutex
	origHostIter := hostIter
	hostIter = func() SelectedHost {
		mu.Lock()
		defer mu.Unlock()
		return origHostIter()
	}

	ctx, cancel := context.WithCancel(qry.Context())
	defer cancel()

	results := make(chan *Iter, 1)

	// Launch the main execution
	go q.run(ctx, qry, hostIter, results)

	// The speculative executions are launched _in addition_ to the main
	// execution, on a timer. So Speculation{2} would make 3 executions running
	// in total.
	if iter := q.speculate(ctx, qry, sp, hostIter, results); iter != nil {
		return iter, nil
	}

	select {
	case iter := <-results:
		return iter, nil
	case <-ctx.Done():
		return &Iter{err: ctx.Err()}, nil
	}
}

func (q *queryExecutor) do(ctx context.Context, qry ExecutableQuery, hostIter NextHost) *Iter {
	selectedHost := hostIter()
	rt := qry.retryPolicy()

	var lastErr error
	var iter *Iter
	for selectedHost != nil {
		host := selectedHost.Info()
		if host == nil || !host.IsUp() {
			selectedHost = hostIter()
			continue
		}

		pool, ok := q.pool.getPool(host)
		if !ok {
			selectedHost = hostIter()
			continue
		}

		conn := pool.Pick()
		if conn == nil {
			selectedHost = hostIter()
			continue
		}

		iter = q.attemptQuery(ctx, qry, conn)
		iter.host = selectedHost.Info()
		// Update host
		switch iter.err {
		case context.Canceled, context.DeadlineExceeded, ErrNotFound:
			// those errors represents logical errors, they should not count
			// toward removing a node from the pool
			selectedHost.Mark(nil)
			return iter
		default:
			selectedHost.Mark(iter.err)
		}

		// Exit if the query was successful
		// or no retry policy defined or retry attempts were reached
		if iter.err == nil || rt == nil || !rt.Attempt(qry) {
			return iter
		}
		lastErr = iter.err

		// If query is unsuccessful, check the error with RetryPolicy to retry
		switch rt.GetRetryType(iter.err) {
		case Retry:
			// retry on the same host
			continue
		case Rethrow, Ignore:
			return iter
		case RetryNextHost:
			// retry on the next host
			selectedHost = hostIter()
			continue
		default:
			// Undefined? Return nil and error, this will panic in the requester
			return &Iter{err: ErrUnknownRetryType}
		}
	}

	if lastErr != nil {
		return &Iter{err: lastErr}
	}

	return &Iter{err: ErrNoConnections}
}

func (q *queryExecutor) run(ctx context.Context, qry ExecutableQuery, hostIter NextHost, results chan<- *Iter) {
	select {
	case results <- q.do(ctx, qry, hostIter):
	case <-ctx.Done():
	}
}
