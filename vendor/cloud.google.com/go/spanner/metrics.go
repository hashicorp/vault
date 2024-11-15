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

package spanner

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/spanner/internal"
)

const (
	builtInMetricsMeterName = "gax-go"

	nativeMetricsPrefix = "spanner.googleapis.com/internal/client/"

	// Monitored resource labels
	monitoredResLabelKeyProject        = "project_id"
	monitoredResLabelKeyInstance       = "instance_id"
	monitoredResLabelKeyInstanceConfig = "instance_config"
	monitoredResLabelKeyLocation       = "location"
	monitoredResLabelKeyClientHash     = "client_hash"

	// Metric labels
	metricLabelKeyClientUID         = "client_uid"
	metricLabelKeyClientName        = "client_name"
	metricLabelKeyDatabase          = "database"
	metricLabelKeyMethod            = "method"
	metricLabelKeyStatus            = "status"
	metricLabelKeyDirectPathEnabled = "directpath_enabled"
	metricLabelKeyDirectPathUsed    = "directpath_used"

	// Metric names
	metricNameOperationLatencies = "operation_latencies"
	metricNameAttemptLatencies   = "attempt_latencies"
	metricNameOperationCount     = "operation_count"
	metricNameAttemptCount       = "attempt_count"

	// Metric units
	metricUnitMS    = "ms"
	metricUnitCount = "1"
)

// These are effectively const, but for testing purposes they are mutable
var (
	// duration between two metric exports
	defaultSamplePeriod = 1 * time.Minute

	clientName = fmt.Sprintf("spanner-go/%v", internal.Version)

	bucketBounds = []float64{0.0, 0.5, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0,
		11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 19.0, 20.0,
		25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0, 300.0, 400.0, 500.0, 650.0,
		800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0, 100000.0, 200000.0,
		400000.0, 800000.0, 1600000.0, 3200000.0}

	// All the built-in metrics have same attributes except 'status' and 'streaming'
	// These attributes need to be added to only few of the metrics
	metricsDetails = map[string]metricInfo{
		metricNameOperationCount: {
			additionalAttrs: []string{
				metricLabelKeyStatus,
			},
			recordedPerAttempt: false,
		},
		metricNameOperationLatencies: {
			additionalAttrs: []string{
				metricLabelKeyStatus,
			},
			recordedPerAttempt: false,
		},
		metricNameAttemptLatencies: {
			additionalAttrs: []string{
				metricLabelKeyStatus,
			},
			recordedPerAttempt: true,
		},
		metricNameAttemptCount: {
			additionalAttrs: []string{
				metricLabelKeyStatus,
			},
			recordedPerAttempt: true,
		},
	}

	// Generates unique client ID in the format go-<random UUID>@<hostname>
	generateClientUID = func() (string, error) {
		hostname := "localhost"
		hostname, err := os.Hostname()
		if err != nil {
			return "", err
		}
		return uuid.NewString() + "@" + strconv.FormatInt(int64(os.Getpid()), 10) + "@" + hostname, nil
	}

	// generateClientHash generates a 6-digit zero-padded lowercase hexadecimal hash
	// using the 10 most significant bits of a 64-bit hash value.
	//
	// The primary purpose of this function is to generate a hash value for the `client_hash`
	// resource label using `client_uid` metric field. The range of values is chosen to be small
	// enough to keep the cardinality of the Resource targets under control. Note: If at later time
	// the range needs to be increased, it can be done by increasing the value of `kPrefixLength` to
	// up to 24 bits without changing the format of the returned value.
	generateClientHash = func(clientUID string) string {
		if clientUID == "" {
			return "000000"
		}

		// Use FNV hash function to generate a 64-bit hash
		hasher := fnv.New64()
		hasher.Write([]byte(clientUID))
		hashValue := hasher.Sum64()

		// Extract the 10 most significant bits
		// Shift right by 54 bits to get the 10 most significant bits
		kPrefixLength := 10
		tenMostSignificantBits := hashValue >> (64 - kPrefixLength)

		// Format the result as a 6-digit zero-padded hexadecimal string
		return fmt.Sprintf("%06x", tenMostSignificantBits)
	}

	detectClientLocation = func(ctx context.Context) string {
		resource, err := gcp.NewDetector().Detect(ctx)
		if err != nil {
			return "global"
		}
		for _, attr := range resource.Attributes() {
			if attr.Key == semconv.CloudRegionKey {
				return attr.Value.AsString()
			}
		}
		// If region is not found, return global
		return "global"
	}

	// GCM exporter should use the same options as Spanner client
	// createExporterOptions takes Spanner client options and returns exporter options
	// Overwritten in tests
	createExporterOptions = func(spannerOpts ...option.ClientOption) []option.ClientOption {
		return spannerOpts
	}
)

type metricInfo struct {
	additionalAttrs    []string
	recordedPerAttempt bool
}

// builtinMetricsTracerFactory is responsible for creating and managing metrics tracers.
type builtinMetricsTracerFactory struct {
	enabled             bool // Indicates if metrics tracing is enabled.
	isDirectPathEnabled bool // Indicates if DirectPath is enabled.

	// shutdown is a function to be called on client close to clean up resources.
	shutdown func()

	// clientAttributes are attributes specific to a client instance that do not change across different function calls on the client.
	clientAttributes []attribute.KeyValue

	// Metrics instruments
	operationLatencies metric.Float64Histogram // Histogram for operation latencies.
	attemptLatencies   metric.Float64Histogram // Histogram for attempt latencies.
	operationCount     metric.Int64Counter     // Counter for the number of operations.
	attemptCount       metric.Int64Counter     // Counter for the number of attempts.
}

func newBuiltinMetricsTracerFactory(ctx context.Context, dbpath string, metricsProvider metric.MeterProvider, compression string, opts ...option.ClientOption) (*builtinMetricsTracerFactory, error) {
	clientUID, err := generateClientUID()
	if err != nil {
		log.Printf("built-in metrics: generateClientUID failed: %v. Using empty string in the %v metric atteribute", err, metricLabelKeyClientUID)
	}
	project, instance, database, err := parseDatabaseName(dbpath)
	if err != nil {
		return nil, err
	}

	tracerFactory := &builtinMetricsTracerFactory{
		enabled: false,
		clientAttributes: []attribute.KeyValue{
			attribute.String(monitoredResLabelKeyProject, project),
			attribute.String(monitoredResLabelKeyInstance, instance),
			attribute.String(metricLabelKeyDatabase, database),
			attribute.String(metricLabelKeyClientUID, clientUID),
			attribute.String(metricLabelKeyClientName, clientName),
			attribute.String(monitoredResLabelKeyClientHash, generateClientHash(clientUID)),
			// Skipping instance config until we have a way to get it
			attribute.String(monitoredResLabelKeyInstanceConfig, "unknown"),
			attribute.String(monitoredResLabelKeyLocation, detectClientLocation(ctx)),
		},
		shutdown: func() {},
	}

	tracerFactory.isDirectPathEnabled = false
	tracerFactory.enabled = false
	var meterProvider *sdkmetric.MeterProvider
	if metricsProvider == nil {
		// Create default meter provider
		mpOptions, err := builtInMeterProviderOptions(project, compression, opts...)
		if err != nil {
			return tracerFactory, err
		}
		meterProvider = sdkmetric.NewMeterProvider(mpOptions...)

		tracerFactory.enabled = true
		tracerFactory.shutdown = func() { meterProvider.Shutdown(ctx) }
	} else {
		switch metricsProvider.(type) {
		case noop.MeterProvider:
			return tracerFactory, nil
		default:
			return tracerFactory, errors.New("unknown MetricsProvider type")
		}
	}

	// Create meter and instruments
	meter := meterProvider.Meter(builtInMetricsMeterName, metric.WithInstrumentationVersion(internal.Version))
	err = tracerFactory.createInstruments(meter)
	return tracerFactory, err
}

func builtInMeterProviderOptions(project, compression string, opts ...option.ClientOption) ([]sdkmetric.Option, error) {
	allOpts := createExporterOptions(opts...)
	defaultExporter, err := newMonitoringExporter(context.Background(), project, compression, allOpts...)
	if err != nil {
		return nil, err
	}

	return []sdkmetric.Option{sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(
			defaultExporter,
			sdkmetric.WithInterval(defaultSamplePeriod),
		),
	)}, nil
}

func (tf *builtinMetricsTracerFactory) createInstruments(meter metric.Meter) error {
	var err error

	// Create operation_latencies
	tf.operationLatencies, err = meter.Float64Histogram(
		nativeMetricsPrefix+metricNameOperationLatencies,
		metric.WithDescription("Total time until final operation success or failure, including retries and backoff."),
		metric.WithUnit(metricUnitMS),
		metric.WithExplicitBucketBoundaries(bucketBounds...),
	)
	if err != nil {
		return err
	}

	// Create attempt_latencies
	tf.attemptLatencies, err = meter.Float64Histogram(
		nativeMetricsPrefix+metricNameAttemptLatencies,
		metric.WithDescription("Client observed latency per RPC attempt."),
		metric.WithUnit(metricUnitMS),
		metric.WithExplicitBucketBoundaries(bucketBounds...),
	)
	if err != nil {
		return err
	}

	// Create operation_count
	tf.operationCount, err = meter.Int64Counter(
		nativeMetricsPrefix+metricNameOperationCount,
		metric.WithDescription("The count of database operations."),
		metric.WithUnit(metricUnitCount),
	)
	if err != nil {
		return err
	}

	// Create attempt_count
	tf.attemptCount, err = meter.Int64Counter(
		nativeMetricsPrefix+metricNameAttemptCount,
		metric.WithDescription("The number of attempts made for the operation, including the initial attempt."),
		metric.WithUnit(metricUnitCount),
	)
	return err
}

// builtinMetricsTracer is created one per operation.
// It is used to store metric instruments, attribute values, and other data required to obtain and record them.
type builtinMetricsTracer struct {
	ctx            context.Context // Context for the tracer.
	builtInEnabled bool            // Indicates if built-in metrics are enabled.

	// clientAttributes are attributes specific to a client instance that do not change across different operations on the client.
	clientAttributes []attribute.KeyValue

	// Metrics instruments
	instrumentOperationLatencies metric.Float64Histogram // Histogram for operation latencies.
	instrumentAttemptLatencies   metric.Float64Histogram // Histogram for attempt latencies.
	instrumentOperationCount     metric.Int64Counter     // Counter for the number of operations.
	instrumentAttemptCount       metric.Int64Counter     // Counter for the number of attempts.

	method string // The method being traced.

	currOp *opTracer // The current operation tracer.
}

// opTracer is used to record metrics for the entire operation, including retries.
// An operation is a logical unit that represents a single method invocation on the client.
// The method might require multiple attempts/RPCs and backoff logic to complete.
type opTracer struct {
	attemptCount int64     // The number of attempts made for the operation.
	startTime    time.Time // The start time of the operation.

	// status is the gRPC status code of the last completed attempt.
	status string

	directPathEnabled bool // Indicates if DirectPath is enabled for the operation.

	currAttempt *attemptTracer // The current attempt tracer.
}

// attemptTracer is used to record metrics for a single attempt within an operation.
type attemptTracer struct {
	startTime time.Time // The start time of the attempt.
	status    string    // The gRPC status code of the attempt.

	directPathUsed bool // Indicates if DirectPath was used for the attempt.
}

// setStartTime sets the start time for the operation.
func (o *opTracer) setStartTime(t time.Time) {
	o.startTime = t
}

// setStartTime sets the start time for the attempt.
func (a *attemptTracer) setStartTime(t time.Time) {
	a.startTime = t
}

// setStatus sets the status for the operation.
func (o *opTracer) setStatus(s string) {
	o.status = s
}

// setStatus sets the status for the attempt.
func (a *attemptTracer) setStatus(s string) {
	a.status = s
}

// incrementAttemptCount increments the attempt count for the operation.
func (o *opTracer) incrementAttemptCount() {
	o.attemptCount++
}

// setDirectPathUsed sets whether DirectPath was used for the attempt.
func (a *attemptTracer) setDirectPathUsed(used bool) {
	a.directPathUsed = used
}

// setDirectPathEnabled sets whether DirectPath is enabled for the operation.
func (o *opTracer) setDirectPathEnabled(enabled bool) {
	o.directPathEnabled = enabled
}

func (tf *builtinMetricsTracerFactory) createBuiltinMetricsTracer(ctx context.Context) builtinMetricsTracer {
	// Operation has started but not the attempt.
	// So, create only operation tracer and not attempt tracer
	currOpTracer := opTracer{}
	currOpTracer.setStartTime(time.Now())
	currOpTracer.setDirectPathEnabled(tf.isDirectPathEnabled)

	return builtinMetricsTracer{
		ctx:            ctx,
		builtInEnabled: tf.enabled,

		currOp:           &currOpTracer,
		clientAttributes: tf.clientAttributes,

		instrumentOperationLatencies: tf.operationLatencies,
		instrumentAttemptLatencies:   tf.attemptLatencies,
		instrumentOperationCount:     tf.operationCount,
		instrumentAttemptCount:       tf.attemptCount,
	}
}

// toOtelMetricAttrs:
// - converts metric attributes values captured throughout the operation / attempt
// to OpenTelemetry attributes format,
// - combines these with common client attributes and returns
func (mt *builtinMetricsTracer) toOtelMetricAttrs(metricName string) ([]attribute.KeyValue, error) {
	if mt.currOp == nil || mt.currOp.currAttempt == nil {
		return nil, fmt.Errorf("unable to create attributes list for unknown metric: %v", metricName)
	}
	// Create attribute key value pairs for attributes common to all metricss
	attrKeyValues := []attribute.KeyValue{
		attribute.String(metricLabelKeyMethod, strings.ReplaceAll(strings.TrimPrefix(mt.method, "/google.spanner.v1."), "/", ".")),
	}
	attrKeyValues = append(attrKeyValues, mt.clientAttributes...)

	// Get metric details
	mDetails, found := metricsDetails[metricName]
	if !found {
		return attrKeyValues, fmt.Errorf("unable to create attributes list for unknown metric: %v", metricName)
	}
	attrKeyValues = append(attrKeyValues, attribute.String(metricLabelKeyDirectPathEnabled, strconv.FormatBool(mt.currOp.directPathEnabled)))
	attrKeyValues = append(attrKeyValues, attribute.String(metricLabelKeyDirectPathUsed, strconv.FormatBool(mt.currOp.currAttempt.directPathUsed)))

	rpcStatus := mt.currOp.status
	if mDetails.recordedPerAttempt {
		rpcStatus = mt.currOp.currAttempt.status
	}

	// Add additional attributes to metrics
	for _, attrKey := range mDetails.additionalAttrs {
		switch attrKey {
		case metricLabelKeyStatus:
			attrKeyValues = append(attrKeyValues, attribute.String(metricLabelKeyStatus, rpcStatus))
		default:
			return attrKeyValues, fmt.Errorf("unknown additional attribute: %v", attrKey)
		}
	}

	return attrKeyValues, nil
}

// Convert error to grpc status error
func convertToGrpcStatusErr(err error) (codes.Code, error) {
	if err == nil {
		return codes.OK, nil
	}

	if errStatus, ok := status.FromError(err); ok {
		return errStatus.Code(), status.Error(errStatus.Code(), errStatus.Message())
	}

	ctxStatus := status.FromContextError(err)
	if ctxStatus.Code() != codes.Unknown {
		return ctxStatus.Code(), status.Error(ctxStatus.Code(), ctxStatus.Message())
	}

	return codes.Unknown, err
}

// recordAttemptCompletion records as many attempt specific metrics as it can
// Ignore errors seen while creating metric attributes since metric can still
// be recorded with rest of the attributes
func recordAttemptCompletion(mt *builtinMetricsTracer) {
	if !mt.builtInEnabled {
		return
	}

	// Calculate elapsed time
	elapsedTime := convertToMs(time.Since(mt.currOp.currAttempt.startTime))

	// Record attempt_latencies
	attemptLatAttrs, err := mt.toOtelMetricAttrs(metricNameAttemptLatencies)
	if err != nil {
		return
	}
	mt.instrumentAttemptLatencies.Record(mt.ctx, elapsedTime, metric.WithAttributes(attemptLatAttrs...))
}

// recordOperationCompletion records as many operation specific metrics as it can
// Ignores error seen while creating metric attributes since metric can still
// be recorded with rest of the attributes
func recordOperationCompletion(mt *builtinMetricsTracer) {
	if !mt.builtInEnabled {
		return
	}

	// Calculate elapsed time
	elapsedTimeMs := convertToMs(time.Since(mt.currOp.startTime))

	// Record operation_count
	opCntAttrs, err := mt.toOtelMetricAttrs(metricNameOperationCount)
	if err != nil {
		return
	}
	mt.instrumentOperationCount.Add(mt.ctx, 1, metric.WithAttributes(opCntAttrs...))

	// Record operation_latencies
	opLatAttrs, err := mt.toOtelMetricAttrs(metricNameOperationLatencies)
	if err != nil {
		return
	}
	mt.instrumentOperationLatencies.Record(mt.ctx, elapsedTimeMs, metric.WithAttributes(opLatAttrs...))

	// Record attempt_count
	attemptCntAttrs, err := mt.toOtelMetricAttrs(metricNameAttemptCount)
	if err != nil {
		return
	}
	mt.instrumentAttemptCount.Add(mt.ctx, mt.currOp.attemptCount, metric.WithAttributes(attemptCntAttrs...))
}

func convertToMs(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / float64(time.Millisecond)
}
