// Copyright 2022 Google LLC
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

package trace

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/api/googleapi"
)

var (
	keyInstance, _  = tag.NewKey("cloudsql_instance")
	keyDialerID, _  = tag.NewKey("cloudsql_dialer_id")
	keyErrorCode, _ = tag.NewKey("cloudsql_error_code")

	mLatencyMS = stats.Int64(
		"cloudsqlconn/latency",
		"The latency in milliseconds per Dial",
		stats.UnitMilliseconds,
	)
	mConnections = stats.Int64(
		"cloudsqlconn/connection",
		"A connect or disconnect event to Cloud SQL",
		stats.UnitDimensionless,
	)
	mDialError = stats.Int64(
		"cloudsqlconn/dial_failure",
		"A failure to dial a Cloud SQL instance",
		stats.UnitDimensionless,
	)
	mSuccessfulRefresh = stats.Int64(
		"cloudsqlconn/refresh_success",
		"A successful certificate refresh operation",
		stats.UnitDimensionless,
	)
	mFailedRefresh = stats.Int64(
		"cloudsqlconn/refresh_failure",
		"A failed certificate refresh operation",
		stats.UnitDimensionless,
	)
	mBytesSent = stats.Int64(
		"cloudsqlconn/bytes_sent",
		"The bytes sent to Cloud SQL",
		stats.UnitDimensionless,
	)
	mBytesReceived = stats.Int64(
		"cloudsqlconn/bytes_received",
		"The bytes received from Cloud SQL",
		stats.UnitDimensionless,
	)

	latencyView = &view.View{
		Name:        "cloudsqlconn/dial_latency",
		Measure:     mLatencyMS,
		Description: "The distribution of dialer latencies (ms)",
		// Latency in buckets, e.g., >=0ms, >=100ms, etc.
		Aggregation: view.Distribution(0, 5, 25, 100, 250, 500, 1000, 2000, 5000, 30000),
		TagKeys:     []tag.Key{keyInstance, keyDialerID},
	}
	connectionsView = &view.View{
		Name:        "cloudsqlconn/open_connections",
		Measure:     mConnections,
		Description: "The current number of open Cloud SQL connections",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{keyInstance, keyDialerID},
	}
	dialFailureView = &view.View{
		Name:        "cloudsqlconn/dial_failure_count",
		Measure:     mDialError,
		Description: "The number of failed dial attempts",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyInstance, keyDialerID},
	}
	refreshCountView = &view.View{
		Name:        "cloudsqlconn/refresh_success_count",
		Measure:     mSuccessfulRefresh,
		Description: "The number of successful certificate refresh operations",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyInstance, keyDialerID},
	}
	failedRefreshCountView = &view.View{
		Name:        "cloudsqlconn/refresh_failure_count",
		Measure:     mFailedRefresh,
		Description: "The number of failed certificate refresh operations",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyInstance, keyDialerID, keyErrorCode},
	}
	bytesSentView = &view.View{
		Name:        "cloudsqlconn/bytes_sent",
		Measure:     mBytesSent,
		Description: "The number of bytes sent to Cloud SQL",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{keyInstance, keyDialerID},
	}
	bytesReceivedView = &view.View{
		Name:        "cloudsqlconn/bytes_received",
		Measure:     mBytesReceived,
		Description: "The number of bytes received from Cloud SQL",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{keyInstance, keyDialerID},
	}

	registerOnce sync.Once
	registerErr  error
)

// InitMetrics registers all views once. Without registering views, metrics will
// not be reported. If any names of the registered views conflict, this function
// returns an error to indicate an internal configuration problem.
func InitMetrics() error {
	registerOnce.Do(func() {
		if rErr := view.Register(
			latencyView,
			connectionsView,
			dialFailureView,
			refreshCountView,
			failedRefreshCountView,
			bytesSentView,
			bytesReceivedView,
		); rErr != nil {
			registerErr = fmt.Errorf("failed to initialize metrics: %v", rErr)
		}
	})
	return registerErr
}

// RecordDialLatency records a latency value for a call to dial.
func RecordDialLatency(ctx context.Context, instance, dialerID string, latency int64) {
	// tag.New creates a new context and errors only if the new tag already
	// exists in the provided context. Since we're adding tags within this
	// package only, we can be confident that there were be no duplicate tags
	// and so can ignore the error.
	ctx, _ = tag.New(ctx, tag.Upsert(keyInstance, instance), tag.Upsert(keyDialerID, dialerID))
	stats.Record(ctx, mLatencyMS.M(latency))
}

// RecordOpenConnections records the number of open connections
func RecordOpenConnections(ctx context.Context, num int64, dialerID, instance string) {
	ctx, _ = tag.New(ctx, tag.Upsert(keyInstance, instance), tag.Upsert(keyDialerID, dialerID))
	stats.Record(ctx, mConnections.M(num))
}

// RecordDialError reports a failed dial attempt. If err is nil, RecordDialError
// is a no-op.
func RecordDialError(ctx context.Context, instance, dialerID string, err error) {
	if err == nil {
		return
	}
	ctx, _ = tag.New(ctx, tag.Upsert(keyInstance, instance), tag.Upsert(keyDialerID, dialerID))
	stats.Record(ctx, mDialError.M(1))
}

// RecordRefreshResult reports the result of a refresh operation, either
// successfull or failed.
func RecordRefreshResult(ctx context.Context, instance, dialerID string, err error) {
	ctx, _ = tag.New(ctx, tag.Upsert(keyInstance, instance), tag.Upsert(keyDialerID, dialerID))
	if err != nil {
		if c := errorCode(err); c != "" {
			ctx, _ = tag.New(ctx, tag.Upsert(keyErrorCode, c))
		}
		stats.Record(ctx, mFailedRefresh.M(1))
		return
	}
	stats.Record(ctx, mSuccessfulRefresh.M(1))
}

// RecordBytesSent reports the number of bytes sent to Cloud SQL
func RecordBytesSent(ctx context.Context, num int64, instance, dialerID string) {
	ctx, _ = tag.New(ctx, tag.Upsert(keyInstance, instance), tag.Upsert(keyDialerID, dialerID))
	stats.Record(ctx, mBytesSent.M(num))
}

// RecordBytesReceived reports the number of bytes received from Cloud SQL
func RecordBytesReceived(ctx context.Context, num int64, instance, dialerID string) {
	ctx, _ = tag.New(ctx, tag.Upsert(keyInstance, instance), tag.Upsert(keyDialerID, dialerID))
	stats.Record(ctx, mBytesReceived.M(num))
}

// errorCode returns an error code as given from the SQL Admin API, provided the
// error wraps a googleapi.Error type. If multiple error codes are returned from
// the API, then a comma-separated string of all codes is returned.
//
// For possible error codes and their meaning see:
// https://cloud.google.com/sql/docs/mysql/admin-api-error-messages
func errorCode(err error) string {
	var apiErr *googleapi.Error
	ok := errors.As(err, &apiErr)
	if !ok {
		return ""
	}
	var codes []string
	for _, e := range apiErr.Errors {
		codes = append(codes, e.Reason)
	}
	return strings.Join(codes, ",")
}
