package gocb

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AnalyticsScanConsistency indicates the level of data consistency desired for an analytics query.
type AnalyticsScanConsistency uint

const (
	// AnalyticsScanConsistencyNotBounded indicates no data consistency is required.
	AnalyticsScanConsistencyNotBounded AnalyticsScanConsistency = iota + 1
	// AnalyticsScanConsistencyRequestPlus indicates that request-level data consistency is required.
	AnalyticsScanConsistencyRequestPlus
)

// AnalyticsOptions is the set of options available to an Analytics query.
type AnalyticsOptions struct {
	// ClientContextID provides a unique ID for this query which can be used matching up requests between connectionManager and
	// server. If not provided will be assigned a uuid value.
	ClientContextID string

	// Priority sets whether this query should be assigned as high priority by the analytics engine.
	Priority             bool
	PositionalParameters []interface{}
	NamedParameters      map[string]interface{}
	Readonly             bool
	ScanConsistency      AnalyticsScanConsistency

	// Raw provides a way to provide extra parameters in the request body for the query.
	Raw map[string]interface{}

	Timeout       time.Duration
	RetryStrategy RetryStrategy

	ParentSpan RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context

	// Internal: This should never be used and is not supported.
	Internal struct {
		User string
	}
}

func (opts *AnalyticsOptions) toMap() (map[string]interface{}, error) {
	execOpts := make(map[string]interface{})

	if opts.ClientContextID == "" {
		execOpts["client_context_id"] = uuid.New().String()
	} else {
		execOpts["client_context_id"] = opts.ClientContextID
	}

	if opts.ScanConsistency != 0 {
		if opts.ScanConsistency == AnalyticsScanConsistencyNotBounded {
			execOpts["scan_consistency"] = "not_bounded"
		} else if opts.ScanConsistency == AnalyticsScanConsistencyRequestPlus {
			execOpts["scan_consistency"] = "request_plus"
		} else {
			return nil, makeInvalidArgumentsError("unexpected consistency option")
		}
	}

	if opts.PositionalParameters != nil && opts.NamedParameters != nil {
		return nil, makeInvalidArgumentsError("positional and named parameters must be used exclusively")
	}

	if opts.PositionalParameters != nil {
		execOpts["args"] = opts.PositionalParameters
	}

	if opts.NamedParameters != nil {
		for key, value := range opts.NamedParameters {
			if !strings.HasPrefix(key, "$") {
				key = "$" + key
			}
			execOpts[key] = value
		}
	}

	if opts.Readonly {
		execOpts["readonly"] = true
	}

	if opts.Raw != nil {
		for k, v := range opts.Raw {
			execOpts[k] = v
		}
	}

	return execOpts, nil
}
