// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1
package limits

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	"github.com/platinummonkey/go-concurrency-limits/core"
	"github.com/platinummonkey/go-concurrency-limits/limit"
	"github.com/platinummonkey/go-concurrency-limits/limiter"
	"github.com/platinummonkey/go-concurrency-limits/strategy"
)

var (
	// ErrCapacity is a new error type to indicate that Vault is not accepting new
	// requests. This should be handled by callers in request paths to return
	// http.StatusServiceUnavailable to the client.
	ErrCapacity = errors.New("Vault server temporarily overloaded")

	// DefaultDebugLogger opts out of the go-concurrency-limits internal Debug
	// logger, since it's rather noisy. We're generating logs of interest in
	// Vault.
	DefaultDebugLogger limit.Logger = nil

	// DefaultMetricsRegistry opts out of the go-concurrency-limits internal
	// metrics because we're tracking what we care about in Vault.
	DefaultMetricsRegistry core.MetricRegistry = core.EmptyMetricRegistryInstance
)

const (
	// Smoothing adjusts how heavily we weight newer high-latency detection.
	// Higher values (>1) place more emphasis on recent measurements. We set
	// this below 1 to better tolerate short-lived spikes in request rate.
	DefaultSmoothing = .1

	// DefaultLongWindow is chosen as a minimum of 1000 samples. longWindow
	// defines sliding window size used for the Exponential Moving Average.
	DefaultLongWindow = 1000
)

// RequestLimiter is a thin wrapper for limiter.DefaultLimiter.
type RequestLimiter struct {
	*limiter.DefaultLimiter
	Flags LimiterFlags
}

// Acquire consults the underlying RequestLimiter to see if a new
// RequestListener can be acquired.
//
// The return values are a *RequestListener, which the caller can use to perform
// latency measurements, and a bool to indicate whether or not a RequestListener
// was acquired.
//
// The returned RequestListener is short-lived and eventually garbage-collected;
// however, the RequestLimiter keeps track of in-flight concurrency using a
// token bucket implementation. The caller must release the resulting Limiter
// token by conducting a measurement.
//
// There are three return cases:
//
// 1) If Request Limiting is disabled, we return an empty RequestListener so all
// measurements are no-ops.
//
// 2) If the request limit has been exceeded, we will not acquire a
// RequestListener and instead return nil, false. No measurement is required,
// since we immediately return from callers with ErrCapacity.
//
// 3) If we have not exceeded the request limit, the caller must call one of
// OnSuccess(), OnDropped(), or OnIgnore() to return a measurement and release
// the underlying Limiter token.
func (l *RequestLimiter) Acquire(ctx context.Context) (*RequestListener, bool) {
	// Transparently handle the case where the limiter is disabled.
	if l == nil || l.DefaultLimiter == nil {
		return &RequestListener{}, true
	}

	lsnr, ok := l.DefaultLimiter.Acquire(ctx)
	if !ok {
		metrics.IncrCounter(([]string{"limits", "concurrency", "service_unavailable"}), 1)
		// If the token acquisition fails, we've reached capacity and we won't
		// get a listener, so just return nil.
		return nil, false
	}

	return &RequestListener{
		DefaultListener: lsnr.(*limiter.DefaultListener),
		released:        new(atomic.Bool),
	}, true
}

// concurrencyChanger adjusts the current allowed concurrency with an
// exponential backoff as we approach the max limit.
func concurrencyChanger(limit int) int {
	change := math.Sqrt(float64(limit))
	if change < 1.0 {
		change = 1.0
	}
	return int(change)
}

var DefaultLimiterFlags = map[string]LimiterFlags{
	// WriteLimiter default flags have a less conservative MinLimit to prevent
	// over-optimizing the request latency, which would result in
	// under-utilization and client starvation.
	WriteLimiter: {
		MinLimit:     100,
		MaxLimit:     5000,
		InitialLimit: 100,
	},

	// SpecialPathLimiter default flags have a conservative MinLimit to allow
	// more aggressive concurrency throttling for CPU-bound workloads such as
	// `pki/issue`.
	SpecialPathLimiter: {
		MinLimit:     5,
		MaxLimit:     5000,
		InitialLimit: 5,
	},
}

// LimiterFlags establish some initial configuration for a new request limiter.
type LimiterFlags struct {
	// MinLimit defines the minimum concurrency floor to prevent over-throttling
	// requests during periods of high traffic.
	MinLimit int `json:"min_limit,omitempty" mapstructure:"min_limit,omitempty"`

	// MaxLimit defines the maximum concurrency ceiling to prevent skewing to a
	// point of no return.
	//
	// We set this to a high value (5000) with the expectation that systems with
	// high-performing specs will tolerate higher limits, while the algorithm
	// will find its own steady-state concurrency well below this threshold in
	// most cases.
	MaxLimit int `json:"max_limit,omitempty" mapstructure:"max_limit,omitempty"`

	// InitialLimit defines the starting concurrency limit prior to any
	// measurements.
	//
	// If we start this value off too high, Vault could become
	// overloaded before the algorithm has a chance to adapt. Setting the value
	// to the minimum is a safety measure which could result in early request
	// rejection; however, the adaptive nature of the algorithm will prevent
	// this from being a prolonged state as the allowed concurrency will
	// increase during normal operation.
	InitialLimit int `json:"initial_limit,omitempty" mapstructure:"initial_limit,omitempty"`
}

// NewRequestLimiter is a basic constructor for the RequestLimiter wrapper. It
// is responsible for setting up the Gradient2 Limit and instantiating a new
// wrapped DefaultLimiter.
func NewRequestLimiter(logger hclog.Logger, name string, flags LimiterFlags) (*RequestLimiter, error) {
	logger.Info("setting up new request limiter",
		"initialLimit", flags.InitialLimit,
		"maxLimit", flags.MaxLimit,
		"minLimit", flags.MinLimit,
	)

	// NewGradient2Limit is the algorithm which drives request limiting
	// decisions. It gathers latency measurements and calculates an Exponential
	// Moving Average to determine whether latency deviation warrants a change
	// in the current concurrency limit.
	lim, err := limit.NewGradient2Limit(name,
		flags.InitialLimit,
		flags.MaxLimit,
		flags.MinLimit,
		concurrencyChanger,
		DefaultSmoothing,
		DefaultLongWindow,
		DefaultDebugLogger,
		DefaultMetricsRegistry,
	)
	if err != nil {
		return &RequestLimiter{}, fmt.Errorf("failed to create gradient2 limit: %w", err)
	}

	strategy := strategy.NewSimpleStrategy(flags.InitialLimit)
	defLimiter, err := limiter.NewDefaultLimiter(lim, 1e9, 1e9, 10, 100, strategy, nil, DefaultMetricsRegistry)
	if err != nil {
		return &RequestLimiter{}, err
	}

	return &RequestLimiter{Flags: flags, DefaultLimiter: defLimiter}, nil
}
