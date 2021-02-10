package gocb

import (
	"fmt"
	"time"
)

// AnalyticsQuery executes the analytics query statement on the server, constraining the query to the bucket and scope.
// VOLATILE: This API is subject to change at any time.
func (s *Scope) AnalyticsQuery(statement string, opts *AnalyticsOptions) (*AnalyticsResult, error) {
	if opts == nil {
		opts = &AnalyticsOptions{}
	}

	span := s.tracer.StartSpan("Query", opts.parentSpan).
		SetTag("couchbase.service", "analytics")
	defer span.Finish()

	timeout := opts.Timeout
	if opts.Timeout == 0 {
		timeout = s.timeoutsConfig.AnalyticsTimeout
	}
	deadline := time.Now().Add(timeout)

	retryStrategy := s.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	queryOpts, err := opts.toMap()
	if err != nil {
		return nil, AnalyticsError{
			InnerError:      wrapError(err, "failed to generate query options"),
			Statement:       statement,
			ClientContextID: opts.ClientContextID,
		}
	}

	var priorityInt int32
	if opts.Priority {
		priorityInt = -1
	}

	queryOpts["statement"] = statement
	queryOpts["query_context"] = fmt.Sprintf("default:`%s`.`%s`", s.BucketName(), s.Name())

	provider, err := s.getAnalyticsProvider()
	if err != nil {
		return nil, AnalyticsError{
			InnerError:      wrapError(err, "failed to get query provider"),
			Statement:       statement,
			ClientContextID: maybeGetAnalyticsOption(queryOpts, "client_context_id"),
		}
	}

	return execAnalyticsQuery(span, queryOpts, priorityInt, deadline, retryStrategy, provider, s.tracer)
}
