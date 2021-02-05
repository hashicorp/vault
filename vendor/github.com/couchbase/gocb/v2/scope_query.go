package gocb

import (
	"fmt"
	"time"
)

// Query executes the query statement on the server, constraining the query to the bucket and scope.
// VOLATILE: This API is subject to change at any time.
func (s *Scope) Query(statement string, opts *QueryOptions) (*QueryResult, error) {
	if opts == nil {
		opts = &QueryOptions{}
	}

	span := s.tracer.StartSpan("Query", opts.parentSpan).
		SetTag("couchbase.service", "query")
	defer span.Finish()

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = s.timeoutsConfig.QueryTimeout
	}
	deadline := time.Now().Add(timeout)

	retryStrategy := s.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	queryOpts, err := opts.toMap()
	if err != nil {
		return nil, QueryError{
			InnerError:      wrapError(err, "failed to generate query options"),
			Statement:       statement,
			ClientContextID: opts.ClientContextID,
		}
	}

	queryOpts["statement"] = statement
	queryOpts["query_context"] = fmt.Sprintf("%s.%s", s.BucketName(), s.Name())

	provider, err := s.getQueryProvider()
	if err != nil {
		return nil, QueryError{
			InnerError:      wrapError(err, "failed to get query provider"),
			Statement:       statement,
			ClientContextID: maybeGetQueryOption(queryOpts, "client_context_id"),
		}
	}

	return execN1qlQuery(span, queryOpts, deadline, retryStrategy, opts.Adhoc, provider, s.tracer)
}
