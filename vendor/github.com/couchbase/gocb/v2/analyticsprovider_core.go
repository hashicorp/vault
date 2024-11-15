package gocb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/couchbase/gocbcore/v10"
)

type analyticsProviderCoreProvider interface {
	AnalyticsQuery(ctx context.Context, opts gocbcore.AnalyticsQueryOptions) (analyticsRowReader, error)
}

type analyticsProviderCore struct {
	provider     analyticsProviderCoreProvider
	mgmtProvider mgmtProvider

	retryStrategyWrapper *coreRetryStrategyWrapper
	transcoder           Transcoder
	analyticsTimeout     time.Duration
	tracer               RequestTracer
	meter                *meterWrapper
}

type jsonAnalyticsMetrics struct {
	ElapsedTime      string `json:"elapsedTime"`
	ExecutionTime    string `json:"executionTime"`
	ResultCount      uint64 `json:"resultCount"`
	ResultSize       uint64 `json:"resultSize"`
	MutationCount    uint64 `json:"mutationCount,omitempty"`
	SortCount        uint64 `json:"sortCount,omitempty"`
	ErrorCount       uint64 `json:"errorCount,omitempty"`
	WarningCount     uint64 `json:"warningCount,omitempty"`
	ProcessedObjects uint64 `json:"processedObjects,omitempty"`
}

type jsonAnalyticsWarning struct {
	Code    uint32 `json:"code"`
	Message string `json:"msg"`
}

type jsonAnalyticsResponse struct {
	RequestID       string                 `json:"requestID"`
	ClientContextID string                 `json:"clientContextID"`
	Status          string                 `json:"status"`
	Warnings        []jsonAnalyticsWarning `json:"warnings"`
	Metrics         jsonAnalyticsMetrics   `json:"metrics"`
	Signature       interface{}            `json:"signature"`
}

func (ap *analyticsProviderCore) AnalyticsQuery(statement string, scope *Scope, opts *AnalyticsOptions) (*AnalyticsResult, error) {
	if opts == nil {
		opts = &AnalyticsOptions{}
	}

	start := time.Now()
	defer ap.meter.ValueRecord(meterValueServiceAnalytics, "analytics", start)

	span := createSpan(ap.tracer, opts.ParentSpan, "analytics", "analytics")
	span.SetAttribute("db.statement", statement)
	if scope != nil {
		span.SetAttribute("db.name", scope.BucketName())
		span.SetAttribute("db.couchbase.scope", scope.Name())
	}
	defer span.End()

	timeout := opts.Timeout
	if opts.Timeout == 0 {
		timeout = ap.analyticsTimeout
	}
	deadline := time.Now().Add(timeout)

	retryStrategy := ap.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newCoreRetryStrategyWrapper(opts.RetryStrategy)
	}

	queryOpts, err := opts.toMap()
	if err != nil {
		return nil, &AnalyticsError{
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
	if scope != nil {
		queryOpts["query_context"] = fmt.Sprintf("default:`%s`.`%s`", scope.BucketName(), scope.Name())
	}

	eSpan := createSpan(ap.tracer, span, "request_encoding", "")
	reqBytes, err := json.Marshal(queryOpts)
	eSpan.End()
	if err != nil {
		return nil, &AnalyticsError{
			InnerError:      wrapError(err, "failed to marshall query body"),
			Statement:       maybeGetAnalyticsOption(queryOpts, "statement"),
			ClientContextID: maybeGetAnalyticsOption(queryOpts, "client_context_id"),
		}
	}

	res, err := ap.provider.AnalyticsQuery(opts.Context, gocbcore.AnalyticsQueryOptions{
		Payload:       reqBytes,
		Priority:      int(priorityInt),
		RetryStrategy: retryStrategy,
		Deadline:      deadline,
		TraceContext:  span.Context(),
		User:          opts.Internal.User,
	})
	if err != nil {
		return nil, maybeEnhanceAnalyticsError(err)
	}

	return newAnalyticsResult(res), nil
}
