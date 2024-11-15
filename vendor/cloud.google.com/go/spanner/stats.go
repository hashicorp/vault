// Copyright 2017 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanner

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/spanner/internal"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/grpc/metadata"
)

const statsPrefix = "cloud.google.com/go/spanner/"

// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
var (
	tagKeyClientID      = tag.MustNewKey("client_id")
	tagKeyDatabase      = tag.MustNewKey("database")
	tagKeyInstance      = tag.MustNewKey("instance_id")
	tagKeyLibVersion    = tag.MustNewKey("library_version")
	tagKeyType          = tag.MustNewKey("type")
	tagKeyIsMultiplexed = tag.MustNewKey("is_multiplexed")

	tagCommonKeys = []tag.Key{tagKeyClientID, tagKeyDatabase, tagKeyInstance, tagKeyLibVersion}

	tagNumInUseSessions = tag.Tag{Key: tagKeyType, Value: "num_in_use_sessions"}
	tagNumSessions      = tag.Tag{Key: tagKeyType, Value: "num_sessions"}

	// Deprecated: With InLine Begin transaction client won't maintain separate read or write sessions
	tagNumBeingPrepared = tag.Tag{Key: tagKeyType, Value: "num_sessions_being_prepared"}
	tagNumReadSessions  = tag.Tag{Key: tagKeyType, Value: "num_read_sessions"}
	tagNumWriteSessions = tag.Tag{Key: tagKeyType, Value: "num_write_prepared_sessions"}

	tagKeyMethod = tag.MustNewKey("grpc_client_method")
	// gfeLatencyMetricsEnabled is used to track if GFELatency and GFEHeaderMissingCount need to be recorded
	gfeLatencyMetricsEnabled = false
	// mutex to avoid data race in reading/writing the above flag
	statsMu = sync.RWMutex{}
)

func recordStat(ctx context.Context, m *stats.Int64Measure, n int64) {
	stats.Record(ctx, m.M(n))
}

var (
	// OpenSessionCount is a measure of the number of sessions currently opened.
	// It is EXPERIMENTAL and subject to change or removal without notice.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get open_session_count metrics.
	OpenSessionCount = stats.Int64(
		statsPrefix+"open_session_count",
		"Number of sessions currently opened",
		stats.UnitDimensionless,
	)

	// OpenSessionCountView is a view of the last value of OpenSessionCount.
	// It is EXPERIMENTAL and subject to change or removal without notice.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get open_session_count metrics.
	OpenSessionCountView = &view.View{
		Measure:     OpenSessionCount,
		Aggregation: view.LastValue(),
		TagKeys:     tagCommonKeys,
	}

	// MaxAllowedSessionsCount is a measure of the maximum number of sessions
	// allowed. Configurable by the user.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get max_allowed_sessions metrics.
	MaxAllowedSessionsCount = stats.Int64(
		statsPrefix+"max_allowed_sessions",
		"The maximum number of sessions allowed. Configurable by the user.",
		stats.UnitDimensionless,
	)

	// MaxAllowedSessionsCountView is a view of the last value of
	// MaxAllowedSessionsCount.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get max_allowed_sessions metrics.
	MaxAllowedSessionsCountView = &view.View{
		Measure:     MaxAllowedSessionsCount,
		Aggregation: view.LastValue(),
		TagKeys:     tagCommonKeys,
	}

	// SessionsCount is a measure of the number of sessions in the pool
	// including both in-use, idle, and being prepared.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get num_sessions_in_pool metrics.
	SessionsCount = stats.Int64(
		statsPrefix+"num_sessions_in_pool",
		"The number of sessions currently in use.",
		stats.UnitDimensionless,
	)

	// SessionsCountView is a view of the last value of SessionsCount.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get num_sessions_in_pool metrics.
	SessionsCountView = &view.View{
		Measure:     SessionsCount,
		Aggregation: view.LastValue(),
		TagKeys:     append(tagCommonKeys, tagKeyType),
	}

	// MaxInUseSessionsCount is a measure of the maximum number of sessions
	// in use during the last 10 minute interval.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get max_in_use_sessions metrics.
	MaxInUseSessionsCount = stats.Int64(
		statsPrefix+"max_in_use_sessions",
		"The maximum number of sessions in use during the last 10 minute interval.",
		stats.UnitDimensionless,
	)

	// MaxInUseSessionsCountView is a view of the last value of
	// MaxInUseSessionsCount.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get max_in_use_sessions metrics.
	MaxInUseSessionsCountView = &view.View{
		Measure:     MaxInUseSessionsCount,
		Aggregation: view.LastValue(),
		TagKeys:     tagCommonKeys,
	}

	// GetSessionTimeoutsCount is a measure of the number of get sessions
	// timeouts due to pool exhaustion.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get get_session_timeouts metrics.
	GetSessionTimeoutsCount = stats.Int64(
		statsPrefix+"get_session_timeouts",
		"The number of get sessions timeouts due to pool exhaustion.",
		stats.UnitDimensionless,
	)

	// GetSessionTimeoutsCountView is a view of the last value of
	// GetSessionTimeoutsCount.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get get_session_timeouts metrics.
	GetSessionTimeoutsCountView = &view.View{
		Measure:     GetSessionTimeoutsCount,
		Aggregation: view.Count(),
		TagKeys:     tagCommonKeys,
	}

	// AcquiredSessionsCount is the number of sessions acquired from
	// the session pool.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get num_acquired_sessions metrics.
	AcquiredSessionsCount = stats.Int64(
		statsPrefix+"num_acquired_sessions",
		"The number of sessions acquired from the session pool.",
		stats.UnitDimensionless,
	)

	// AcquiredSessionsCountView is a view of the last value of
	// AcquiredSessionsCount.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get num_acquired_sessions metrics.
	AcquiredSessionsCountView = &view.View{
		Measure:     AcquiredSessionsCount,
		Aggregation: view.Count(),
		TagKeys:     tagCommonKeys,
	}

	// ReleasedSessionsCount is the number of sessions released by the user
	// and pool maintainer.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get num_released_sessions metrics.
	ReleasedSessionsCount = stats.Int64(
		statsPrefix+"num_released_sessions",
		"The number of sessions released by the user and pool maintainer.",
		stats.UnitDimensionless,
	)

	// ReleasedSessionsCountView is a view of the last value of
	// ReleasedSessionsCount.
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get num_released_sessions metrics.
	ReleasedSessionsCountView = &view.View{
		Measure:     ReleasedSessionsCount,
		Aggregation: view.Count(),
		TagKeys:     tagCommonKeys,
	}

	// GFELatency is the latency between Google's network receiving an RPC and reading back the first byte of the response
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get gfe_latency metrics.
	GFELatency = stats.Int64(
		statsPrefix+"gfe_latency",
		"Latency between Google's network receiving an RPC and reading back the first byte of the response",
		stats.UnitMilliseconds,
	)

	// GFELatencyView is the view of distribution of GFELatency values
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get gfe_latency metrics.
	GFELatencyView = &view.View{
		Name:        "cloud.google.com/go/spanner/gfe_latency",
		Measure:     GFELatency,
		Description: "Latency between Google's network receives an RPC and reads back the first byte of the response",
		Aggregation: view.Distribution(0.0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0, 10.0, 13.0,
			16.0, 20.0, 25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0,
			300.0, 400.0, 500.0, 650.0, 800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0,
			100000.0),
		TagKeys: append(tagCommonKeys, tagKeyMethod),
	}

	// GFEHeaderMissingCount is the number of RPC responses received without the server-timing header, most likely means that the RPC never reached Google's network
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get gfe_header_missing_count metrics.
	GFEHeaderMissingCount = stats.Int64(
		statsPrefix+"gfe_header_missing_count",
		"Number of RPC responses received without the server-timing header, most likely means that the RPC never reached Google's network",
		stats.UnitDimensionless,
	)

	// GFEHeaderMissingCountView is the view of number of GFEHeaderMissingCount
	//
	// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry to get gfe_header_missing_count metrics.
	GFEHeaderMissingCountView = &view.View{
		Name:        "cloud.google.com/go/spanner/gfe_header_missing_count",
		Measure:     GFEHeaderMissingCount,
		Description: "Number of RPC responses received without the server-timing header, most likely means that the RPC never reached Google's network",
		Aggregation: view.Count(),
		TagKeys:     append(tagCommonKeys, tagKeyMethod),
	}
)

// EnableStatViews enables all views of metrics relate to session management.
//
// Deprecated: OpenCensus project is deprecated.
// Use EnableOpenTelemetryMetrics to get Session metrics through OpenTelemetry instrumentation.
func EnableStatViews() error {
	return view.Register(
		OpenSessionCountView,
		MaxAllowedSessionsCountView,
		SessionsCountView,
		MaxInUseSessionsCountView,
		GetSessionTimeoutsCountView,
		AcquiredSessionsCountView,
		ReleasedSessionsCountView,
	)
}

// EnableGfeLatencyView enables GFELatency metric
//
// Deprecated: OpenCensus project is deprecated.
// Use EnableOpenTelemetryMetrics to get GfeLatency metrics through OpenTelemetry instrumentation.
func EnableGfeLatencyView() error {
	setGFELatencyMetricsFlag(true)
	return view.Register(GFELatencyView)
}

// EnableGfeHeaderMissingCountView enables GFEHeaderMissingCount metric
//
// Deprecated: OpenCensus project is deprecated.
// Use EnableOpenTelemetryMetrics to get GfeHeaderMissingCount metrics through OpenTelemetry instrumentation.
func EnableGfeHeaderMissingCountView() error {
	setGFELatencyMetricsFlag(true)
	return view.Register(GFEHeaderMissingCountView)
}

// EnableGfeLatencyAndHeaderMissingCountViews enables GFEHeaderMissingCount and GFELatency metric
//
// Deprecated: OpenCensus project is deprecated.
// Use EnableOpenTelemetryMetrics to get GfeLatency and GfeHeaderMissingCount metrics through OpenTelemetry instrumentation.
func EnableGfeLatencyAndHeaderMissingCountViews() error {
	setGFELatencyMetricsFlag(true)
	return view.Register(
		GFELatencyView,
		GFEHeaderMissingCountView,
	)
}

// Deprecated: OpenCensus project is deprecated.
func getGFELatencyMetricsFlag() bool {
	statsMu.RLock()
	defer statsMu.RUnlock()
	return gfeLatencyMetricsEnabled
}

// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
func setGFELatencyMetricsFlag(enable bool) {
	statsMu.Lock()
	gfeLatencyMetricsEnabled = enable
	statsMu.Unlock()
}

// DisableGfeLatencyAndHeaderMissingCountViews disables GFEHeaderMissingCount and GFELatency metric
//
// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
func DisableGfeLatencyAndHeaderMissingCountViews() {
	setGFELatencyMetricsFlag(false)
	view.Unregister(
		GFELatencyView,
		GFEHeaderMissingCountView,
	)
}

// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
func captureGFELatencyStats(ctx context.Context, md metadata.MD, keyMethod string) error {
	if len(md.Get("server-timing")) == 0 {
		recordStat(ctx, GFEHeaderMissingCount, 1)
		return nil
	}
	serverTiming := md.Get("server-timing")[0]
	gfeLatency, err := strconv.Atoi(strings.TrimPrefix(serverTiming, "gfet4t7; dur="))
	if !strings.HasPrefix(serverTiming, "gfet4t7; dur=") || err != nil {
		return err
	}
	// Record GFE latency with OpenCensus.
	ctx = tag.NewContext(ctx, tag.FromContext(ctx))
	ctx, err = tag.New(ctx, tag.Insert(tagKeyMethod, keyMethod))
	if err != nil {
		return err
	}
	recordStat(ctx, GFELatency, int64(gfeLatency))
	return nil
}

// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
func createContextAndCaptureGFELatencyMetrics(ctx context.Context, ct *commonTags, md metadata.MD, keyMethod string) error {
	var ctxGFE, err = tag.New(ctx,
		tag.Upsert(tagKeyClientID, ct.clientID),
		tag.Upsert(tagKeyDatabase, ct.database),
		tag.Upsert(tagKeyInstance, ct.instance),
		tag.Upsert(tagKeyLibVersion, ct.libVersion),
	)
	if err != nil {
		return err
	}
	return captureGFELatencyStats(ctxGFE, md, keyMethod)
}

// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
func getCommonTags(sc *sessionClient) *commonTags {
	_, instance, database, err := parseDatabaseName(sc.database)
	if err != nil {
		return nil
	}
	return &commonTags{
		clientID:   sc.id,
		database:   database,
		instance:   instance,
		libVersion: internal.Version,
	}
}

// commonTags are common key-value pairs of data associated with the GFELatency measure
// Deprecated: OpenCensus project is deprecated. Use OpenTelemetry for capturing metrics.
type commonTags struct {
	// Client ID
	clientID string
	// Database Name
	database string
	// Instance ID
	instance string
	// Library Version
	libVersion string
}
