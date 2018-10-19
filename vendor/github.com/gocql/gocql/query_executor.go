package gocql

import (
	"sync"
	"time"
)

type ExecutableQuery interface {
	execute(conn *Conn) *Iter
	attempt(keyspace string, end, start time.Time, iter *Iter, host *HostInfo)
	retryPolicy() RetryPolicy
	speculativeExecutionPolicy() SpeculativeExecutionPolicy
	GetRoutingKey() ([]byte, error)
	Keyspace() string
	Cancel()
	IsIdempotent() bool
	RetryableQuery
}

type queryExecutor struct {
	pool   *policyConnPool
	policy HostSelectionPolicy
}

type queryResponse struct {
	iter *Iter
	err  error
}

func (q *queryExecutor) attemptQuery(qry ExecutableQuery, conn *Conn) *Iter {
	start := time.Now()
	iter := qry.execute(conn)
	end := time.Now()

	qry.attempt(q.pool.keyspace, end, start, iter, conn.host)

	return iter
}

func (q *queryExecutor) executeQuery(qry ExecutableQuery) (*Iter, error) {

	// check if the query is not marked as idempotent, if
	// it is, we force the policy to NonSpeculative
	sp := qry.speculativeExecutionPolicy()
	if !qry.IsIdempotent() {
		sp = NonSpeculativeExecution{}
	}

	results := make(chan queryResponse, 1)
	stop := make(chan struct{})
	defer close(stop)
	var specWG sync.WaitGroup

	// Launch the main execution
	specWG.Add(1)
	go q.run(qry, &specWG, results, stop)

	// The speculative executions are launched _in addition_ to the main
	// execution, on a timer. So Speculation{2} would make 3 executions running
	// in total.
	go func() {
		// Handle the closing of the resources. We do it here because it's
		// right after we finish launching executions. Otherwise clearing the
		// wait group is complicated.
		defer func() {
			specWG.Wait()
			close(results)
		}()

		// setup a ticker
		ticker := time.NewTicker(sp.Delay())
		defer ticker.Stop()

		for i := 0; i < sp.Attempts(); i++ {
			select {
			case <-ticker.C:
				// Launch the additional execution
				specWG.Add(1)
				go q.run(qry, &specWG, results, stop)
			case <-qry.GetContext().Done():
				// not starting additional executions
				return
			case <-stop:
				// not starting additional executions
				return
			}
		}
	}()

	res := <-results
	if res.iter == nil && res.err == nil {
		// if we're here, the results channel was closed, so no more hosts
		return nil, ErrNoConnections
	}
	return res.iter, res.err
}

func (q *queryExecutor) run(qry ExecutableQuery, specWG *sync.WaitGroup, results chan queryResponse, stop chan struct{}) {
	// Handle the wait group
	defer specWG.Done()

	hostIter := q.policy.Pick(qry)
	selectedHost := hostIter()
	rt := qry.retryPolicy()

	var iter *Iter
	for selectedHost != nil {
		host := selectedHost.Info()
		if host == nil || !host.IsUp() {
			continue
		}

		pool, ok := q.pool.getPool(host)
		if !ok {
			continue
		}

		conn := pool.Pick()
		if conn == nil {
			continue
		}

		select {
		case <-stop:
			// stop this execution and return
			return
		default:
			// Run the query
			iter = q.attemptQuery(qry, conn)
			iter.host = selectedHost.Info()
			// Update host
			selectedHost.Mark(iter.err)

			// Exit if the query was successful
			// or no retry policy defined or retry attempts were reached
			if iter.err == nil || rt == nil || !rt.Attempt(qry) {
				results <- queryResponse{iter: iter}
				return
			}

			// If query is unsuccessful, check the error with RetryPolicy to retry
			switch rt.GetRetryType(iter.err) {
			case Retry:
				// retry on the same host
				continue
			case Rethrow:
				results <- queryResponse{err: iter.err}
				return
			case Ignore:
				results <- queryResponse{iter: iter}
				return
			case RetryNextHost:
				// retry on the next host
				selectedHost = hostIter()
				if selectedHost == nil {
					results <- queryResponse{iter: iter}
					return
				}
				continue
			default:
				// Undefined? Return nil and error, this will panic in the requester
				results <- queryResponse{iter: nil, err: ErrUnknownRetryType}
				return
			}
		}

	}
	// All hosts are exhausted, return nothing
}
