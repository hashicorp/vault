package gocb

import (
	"fmt"
	"time"
)

// Query executes the query statement on the server, constraining the query to the bucket and scope.
func (s *Scope) Query(statement string, opts *QueryOptions) (*QueryResult, error) {
	if opts == nil {
		opts = &QueryOptions{}
	}

	start := time.Now()
	defer s.meter.ValueRecord(meterValueServiceQuery, "query", start)

	span := createSpan(s.tracer, opts.ParentSpan, "query", "query")
	span.SetAttribute("db.statement", statement)
	span.SetAttribute("db.name", s.BucketName())
	span.SetAttribute("db.couchbase.scope", s.Name())
	defer span.End()

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

	return execN1qlQuery(opts.Context, span, queryOpts, deadline, retryStrategy, opts.Adhoc, provider, s.tracer,
		opts.Internal.User)
}
