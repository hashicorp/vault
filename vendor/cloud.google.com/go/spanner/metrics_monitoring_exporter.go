/*
Copyright 2024 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This is a modified version of https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/blob/exporter/metric/v0.46.0/exporter/metric/metric.go

package spanner

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"cloud.google.com/go/spanner/internal"
	"github.com/googleapis/gax-go/v2/callctx"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/sdk/metric"
	otelmetricdata "go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/api/distribution"
	googlemetricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	spannerResourceType = "spanner_instance_client"

	// The number of timeserieses to send to Google Cloud Monitoring in a single request. This
	// is a hard limit in the GCM API, so we never want to exceed 200.
	// https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.timeSeries/create
	sendBatchSize = 200
)

var (
	monitoredResLabelsSet = map[string]bool{
		monitoredResLabelKeyProject:        true,
		monitoredResLabelKeyInstance:       true,
		monitoredResLabelKeyInstanceConfig: true,
		monitoredResLabelKeyLocation:       true,
		monitoredResLabelKeyClientHash:     true,
	}

	errShutdown = fmt.Errorf("exporter is shutdown")
)

type errUnexpectedAggregationKind struct {
	kind string
}

func (e errUnexpectedAggregationKind) Error() string {
	return fmt.Sprintf("the metric kind is unexpected: %v", e.kind)
}

// monitoringExporter is the implementation of OpenTelemetry metric exporter for
// Google Cloud Monitoring.
// Default exporter for built-in metrics
type monitoringExporter struct {
	shutdown     chan struct{}
	client       *monitoring.MetricClient
	shutdownOnce sync.Once
	projectID    string
	compression  string
}

func newMonitoringExporter(ctx context.Context, project, compression string, opts ...option.ClientOption) (*monitoringExporter, error) {
	client, err := monitoring.NewMetricClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &monitoringExporter{
		client:      client,
		shutdown:    make(chan struct{}),
		projectID:   project,
		compression: compression,
	}, nil
}

// ForceFlush does nothing, the exporter holds no state.
func (e *monitoringExporter) ForceFlush(ctx context.Context) error { return ctx.Err() }

// Shutdown shuts down the client connections.
func (e *monitoringExporter) Shutdown(ctx context.Context) error {
	err := errShutdown
	e.shutdownOnce.Do(func() {
		close(e.shutdown)
		err = errors.Join(ctx.Err(), e.client.Close())
	})
	return err
}

// Export exports OpenTelemetry Metrics to Google Cloud Monitoring.
func (me *monitoringExporter) Export(ctx context.Context, rm *otelmetricdata.ResourceMetrics) error {
	select {
	case <-me.shutdown:
		return errShutdown
	default:
	}

	return me.exportTimeSeries(ctx, rm)
}

// Temporality returns the Temporality to use for an instrument kind.
func (me *monitoringExporter) Temporality(ik otelmetric.InstrumentKind) otelmetricdata.Temporality {
	return otelmetricdata.CumulativeTemporality
}

// Aggregation returns the Aggregation to use for an instrument kind.
func (me *monitoringExporter) Aggregation(ik otelmetric.InstrumentKind) otelmetric.Aggregation {
	return otelmetric.DefaultAggregationSelector(ik)
}

// exportTimeSeries create TimeSeries from the records in cps.
// res should be the common resource among all TimeSeries, such as instance id, application name and so on.
func (me *monitoringExporter) exportTimeSeries(ctx context.Context, rm *otelmetricdata.ResourceMetrics) error {
	tss, err := me.recordsToTimeSeriesPbs(rm)
	if len(tss) == 0 {
		return err
	}

	name := fmt.Sprintf("projects/%s", me.projectID)
	ctx = callctx.SetHeaders(ctx, "x-goog-api-client", "gccl/"+internal.Version)
	if me.compression == gzip.Name {
		ctx = callctx.SetHeaders(ctx, requestsCompressionHeader, gzip.Name)
	}
	errs := []error{err}
	for i := 0; i < len(tss); i += sendBatchSize {
		j := i + sendBatchSize
		if j >= len(tss) {
			j = len(tss)
		}
		req := &monitoringpb.CreateTimeSeriesRequest{
			Name:       name,
			TimeSeries: tss[i:j],
		}
		err = me.client.CreateServiceTimeSeries(ctx, req)
		if err != nil {
			if status.Code(err) == codes.PermissionDenied {
				err = fmt.Errorf("%w Need monitoring metric writer permission on project=%s. Follow https://cloud.google.com/spanner/docs/view-manage-client-side-metrics#access-client-side-metrics to set up permissions",
					err, me.projectID)
			}
		}
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// recordToMetricAndMonitoredResourcePbs converts data from records to Metric and Monitored resource proto type for Cloud Monitoring.
func (me *monitoringExporter) recordToMetricAndMonitoredResourcePbs(metrics otelmetricdata.Metrics, attributes attribute.Set) (*googlemetricpb.Metric, *monitoredrespb.MonitoredResource) {
	mr := &monitoredrespb.MonitoredResource{
		Type:   spannerResourceType,
		Labels: map[string]string{},
	}
	labels := make(map[string]string)
	addAttributes := func(attr *attribute.Set) {
		iter := attr.Iter()
		for iter.Next() {
			kv := iter.Attribute()
			labelKey := string(kv.Key)

			if _, isResLabel := monitoredResLabelsSet[labelKey]; isResLabel {
				// Add labels to monitored resource
				mr.Labels[labelKey] = kv.Value.Emit()
			} else {
				// Add labels to metric
				labels[labelKey] = kv.Value.Emit()

			}
		}
	}
	addAttributes(&attributes)
	return &googlemetricpb.Metric{
		Type:   metrics.Name,
		Labels: labels,
	}, mr
}

func (me *monitoringExporter) recordsToTimeSeriesPbs(rm *otelmetricdata.ResourceMetrics) ([]*monitoringpb.TimeSeries, error) {
	var (
		tss  []*monitoringpb.TimeSeries
		errs []error
	)
	for _, scope := range rm.ScopeMetrics {
		if scope.Scope.Name != builtInMetricsMeterName {
			// Filter out metric data for instruments that are not part of the spanner builtin metrics
			continue
		}
		for _, metrics := range scope.Metrics {
			ts, err := me.recordToTimeSeriesPb(metrics)
			errs = append(errs, err)
			tss = append(tss, ts...)
		}
	}

	return tss, errors.Join(errs...)
}

// recordToTimeSeriesPb converts record to TimeSeries proto type with common resource.
// ref. https://cloud.google.com/monitoring/api/ref_v3/rest/v3/TimeSeries
func (me *monitoringExporter) recordToTimeSeriesPb(m otelmetricdata.Metrics) ([]*monitoringpb.TimeSeries, error) {
	var tss []*monitoringpb.TimeSeries
	var errs []error
	if m.Data == nil {
		return nil, nil
	}
	switch a := m.Data.(type) {
	case otelmetricdata.Histogram[float64]:
		for _, point := range a.DataPoints {
			metric, mr := me.recordToMetricAndMonitoredResourcePbs(m, point.Attributes)
			ts, err := histogramToTimeSeries(point, m, mr)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			ts.Metric = metric
			tss = append(tss, ts)
		}
	case otelmetricdata.Sum[int64]:
		for _, point := range a.DataPoints {
			metric, mr := me.recordToMetricAndMonitoredResourcePbs(m, point.Attributes)
			var ts *monitoringpb.TimeSeries
			var err error
			ts, err = sumToTimeSeries[int64](point, m, mr)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			ts.Metric = metric
			tss = append(tss, ts)
		}
	default:
		errs = append(errs, errUnexpectedAggregationKind{kind: reflect.TypeOf(m.Data).String()})
	}
	return tss, errors.Join(errs...)
}

func sumToTimeSeries[N int64 | float64](point otelmetricdata.DataPoint[N], metrics otelmetricdata.Metrics, mr *monitoredrespb.MonitoredResource) (*monitoringpb.TimeSeries, error) {
	interval, err := toNonemptyTimeIntervalpb(point.StartTime, point.Time)
	if err != nil {
		return nil, err
	}
	value, valueType := numberDataPointToValue[N](point)
	return &monitoringpb.TimeSeries{
		Resource:   mr,
		Unit:       string(metrics.Unit),
		MetricKind: googlemetricpb.MetricDescriptor_CUMULATIVE,
		ValueType:  valueType,
		Points: []*monitoringpb.Point{{
			Interval: interval,
			Value:    value,
		}},
	}, nil
}

func histogramToTimeSeries[N int64 | float64](point otelmetricdata.HistogramDataPoint[N], metrics otelmetricdata.Metrics, mr *monitoredrespb.MonitoredResource) (*monitoringpb.TimeSeries, error) {
	interval, err := toNonemptyTimeIntervalpb(point.StartTime, point.Time)
	if err != nil {
		return nil, err
	}
	distributionValue := histToDistribution(point)
	return &monitoringpb.TimeSeries{
		Resource:   mr,
		Unit:       string(metrics.Unit),
		MetricKind: googlemetricpb.MetricDescriptor_CUMULATIVE,
		ValueType:  googlemetricpb.MetricDescriptor_DISTRIBUTION,
		Points: []*monitoringpb.Point{{
			Interval: interval,
			Value: &monitoringpb.TypedValue{
				Value: &monitoringpb.TypedValue_DistributionValue{
					DistributionValue: distributionValue,
				},
			},
		}},
	}, nil
}

func toNonemptyTimeIntervalpb(start, end time.Time) (*monitoringpb.TimeInterval, error) {
	// The end time of a new interval must be at least a millisecond after the end time of the
	// previous interval, for all non-gauge types.
	// https://cloud.google.com/monitoring/api/ref_v3/rpc/google.monitoring.v3#timeinterval
	if end.Sub(start).Milliseconds() <= 1 {
		end = start.Add(time.Millisecond)
	}
	startpb := timestamppb.New(start)
	endpb := timestamppb.New(end)
	err := errors.Join(
		startpb.CheckValid(),
		endpb.CheckValid(),
	)
	if err != nil {
		return nil, err
	}

	return &monitoringpb.TimeInterval{
		StartTime: startpb,
		EndTime:   endpb,
	}, nil
}

func histToDistribution[N int64 | float64](hist otelmetricdata.HistogramDataPoint[N]) *distribution.Distribution {
	counts := make([]int64, len(hist.BucketCounts))
	for i, v := range hist.BucketCounts {
		counts[i] = int64(v)
	}
	var mean float64
	if !math.IsNaN(float64(hist.Sum)) && hist.Count > 0 { // Avoid divide-by-zero
		mean = float64(hist.Sum) / float64(hist.Count)
	}
	return &distribution.Distribution{
		Count:        int64(hist.Count),
		Mean:         mean,
		BucketCounts: counts,
		BucketOptions: &distribution.Distribution_BucketOptions{
			Options: &distribution.Distribution_BucketOptions_ExplicitBuckets{
				ExplicitBuckets: &distribution.Distribution_BucketOptions_Explicit{
					Bounds: hist.Bounds,
				},
			},
		},
	}
}

func numberDataPointToValue[N int64 | float64](
	point otelmetricdata.DataPoint[N],
) (*monitoringpb.TypedValue, googlemetricpb.MetricDescriptor_ValueType) {
	switch v := any(point.Value).(type) {
	case int64:
		return &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: v,
			}},
			googlemetricpb.MetricDescriptor_INT64
	case float64:
		return &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{
				DoubleValue: v,
			}},
			googlemetricpb.MetricDescriptor_DOUBLE
	}
	// It is impossible to reach this statement
	return nil, googlemetricpb.MetricDescriptor_INT64
}
