// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package stackdriver

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/compute/metadata"
	monitoring "cloud.google.com/go/monitoring/apiv3"
	metrics "github.com/armon/go-metrics"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	distributionpb "google.golang.org/genproto/googleapis/api/distribution"
	"google.golang.org/genproto/googleapis/api/metric"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// Sink conforms to the metrics.MetricSink interface and is used to transmit
// metrics information to stackdriver.
//
// Sink performs in-process aggregation of metrics to limit calls to
// stackdriver.
type Sink struct {
	client    *monitoring.MetricClient
	interval  time.Duration
	firstTime time.Time

	gauges     map[string]*gauge
	counters   map[string]*counter
	histograms map[string]*histogram

	bucketer  BucketFn
	extractor ExtractLabelsFn
	taskInfo  *taskInfo

	mu        sync.Mutex
	debugLogs bool
}

// Config options for the stackdriver Sink.
type Config struct {
	// The Google Cloud Project ID to publish metrics to.
	// Optional. GCP instance metadata is used to determine the ProjectID if
	// not set.
	ProjectID string
	// The label extractor provides a way to rewrite metric key into a new key
	// with multiple labels. This is useful if the instrumented code includes
	// variable parameters within a metric name.
	// Optional. Defaults to DefaultLabelExtractor.
	LabelExtractor ExtractLabelsFn
	// The bucketer is used to determine histogram bucket boundaries
	// for the sampled metrics. This will execute before the LabelExtractor.
	// Optional. Defaults to DefaultBucketer.
	Bucketer BucketFn
	// The interval between sampled metric points. Must be > 1 minute.
	// https://cloud.google.com/monitoring/custom-metrics/creating-metrics#writing-ts
	// Optional. Defaults to 1 minute.
	ReportingInterval time.Duration

	// The location of the running task. See:
	// https://cloud.google.com/monitoring/api/resources#tag_generic_task
	// Optional. GCP instance metadata is used to determine the location,
	// otherwise it defaults to 'global'.
	Location string
	// The namespace for the running task. See:
	// https://cloud.google.com/monitoring/api/resources#tag_generic_task
	// Optional. Defaults to 'default'.
	Namespace string
	// The job name for the running task. See:
	// https://cloud.google.com/monitoring/api/resources#tag_generic_task
	// Optional. Defaults to the running program name.
	Job string
	// The task ID for the running task. See:
	// https://cloud.google.com/monitoring/api/resources#tag_generic_task
	// Optional. Defaults to a combination of hostname+pid.
	TaskID string

	// Debug logging. Errors will be logged to stderr, but setting this to true
	// will log additional information that is helpful when debugging errors.
	// Optional. Defaults to false.
	DebugLogs bool
}

type taskInfo struct {
	ProjectID string
	Location  string
	Namespace string
	Job       string
	TaskID    string
}

// BucketFn returns the histogram bucket thresholds based on the given metric
// name.
type BucketFn func([]string) []float64

// ExtractLabelsFn converts a given metric name and type into a new metric
// name and optionally additional labels. Errors will prevent the metric from
// writing to stackdriver.
type ExtractLabelsFn func([]string, string) ([]string, []metrics.Label, error)

// DefaultBucketer is the default BucketFn used to determing bucketing values
// for metrics.
func DefaultBucketer(key []string) []float64 {
	return []float64{10.0, 25.0, 50.0, 100.0, 150.0, 200.0, 250.0, 300.0, 500.0, 1000.0, 1500.0, 2000.0, 3000.0, 4000.0, 5000.0}
}

// DefaultLabelExtractor is the default ExtractLabelsFn and is a direct
// passthrough. Counter and Gauge metrics are renamed to include the type in
// their name to avoid duplicate metrics with the same name but different
// types (which is allowed by go-metrics but not by Stackdriver).
func DefaultLabelExtractor(key []string, kind string) ([]string, []metrics.Label, error) {
	switch kind {
	case "counter":
		return append(key, "counter"), nil, nil
	case "gauge":
		return append(key, "gauge"), nil, nil
	case "histogram":
		return key, nil, nil
	}
	return nil, nil, fmt.Errorf("Unknown metric kind: %s", kind)
}

// NewSink creates a Sink to flush metrics to stackdriver every interval. The
// interval should be greater than 1 minute.
func NewSink(client *monitoring.MetricClient, config *Config) *Sink {
	s := &Sink{
		client:    client,
		extractor: config.LabelExtractor,
		bucketer:  config.Bucketer,
		interval:  config.ReportingInterval,
		taskInfo: &taskInfo{
			ProjectID: config.ProjectID,
			Location:  config.Location,
			Namespace: config.Namespace,
			Job:       config.Job,
			TaskID:    config.TaskID,
		},
		debugLogs: config.DebugLogs,
	}

	// apply defaults if not configured explicitly
	if s.extractor == nil {
		s.extractor = DefaultLabelExtractor
	}
	if s.bucketer == nil {
		s.bucketer = DefaultBucketer
	}
	if s.interval < 60*time.Second {
		s.interval = 60 * time.Second
	}
	if s.taskInfo.ProjectID == "" {
		id, err := metadata.ProjectID()
		if err != nil {
			log.Printf("could not configure go-metrics stackdriver ProjectID: %s", err)
		}
		s.taskInfo.ProjectID = id
	}
	if s.taskInfo.Location == "" {
		// attempt to detect
		zone, err := metadata.Zone()
		if err != nil {
			log.Printf("could not configure go-metric stackdriver location: %s", err)
			zone = "global"
		}
		s.taskInfo.Location = zone
	}
	if s.taskInfo.Namespace == "" {
		s.taskInfo.Namespace = "default"
	}
	if s.taskInfo.Job == "" {
		s.taskInfo.Job = path.Base(os.Args[0])
	}
	if s.taskInfo.TaskID == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
		s.taskInfo.TaskID = "go-" + strconv.Itoa(os.Getpid()) + "@" + hostname
	}

	s.reset()

	// run cancelable goroutine that reports on interval
	go s.flushMetrics(context.Background())

	return s
}

func (s *Sink) flushMetrics(ctx context.Context) {
	if s.interval == 0*time.Second {
		return
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("stopped flushing metrics")
			return
		case <-ticker.C:
			s.report(ctx)
		}
	}
}

func (s *Sink) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.firstTime = time.Now()
	s.gauges = make(map[string]*gauge)
	s.counters = make(map[string]*counter)
	s.histograms = make(map[string]*histogram)
}

func (s *Sink) deep() (time.Time, map[string]*gauge, map[string]*counter, map[string]*histogram) {
	rGauges := make(map[string]*gauge, len(s.gauges))
	rCounters := make(map[string]*counter, len(s.counters))
	rHistograms := make(map[string]*histogram, len(s.histograms))

	s.mu.Lock()
	end := time.Now()
	for k, v := range s.gauges {
		rGauges[k] = &gauge{
			name:  v.name,
			value: v.value,
		}
	}
	for k, v := range s.counters {
		rCounters[k] = &counter{
			name:  v.name,
			value: v.value,
		}
	}
	for k, v := range s.histograms {
		r := &histogram{
			name:    v.name,
			buckets: v.buckets,
			counts:  make([]int64, len(v.counts)),
		}
		copy(r.counts, v.counts)
		rHistograms[k] = r
	}
	s.mu.Unlock()

	return end, rGauges, rCounters, rHistograms
}

func (s *Sink) report(ctx context.Context) {
	end, rGauges, rCounters, rHistograms := s.deep()

	// https://cloud.google.com/monitoring/api/resources
	resource := &monitoredrespb.MonitoredResource{
		Type: "generic_task",
		Labels: map[string]string{
			"project_id": s.taskInfo.ProjectID,
			"location":   s.taskInfo.Location,
			"namespace":  s.taskInfo.Namespace,
			"job":        s.taskInfo.Job,
			"task_id":    s.taskInfo.TaskID,
		},
	}

	ts := []*monitoringpb.TimeSeries{}

	for _, v := range rCounters {
		name, labels, err := v.name.labelMap(s.extractor, "counter")
		if err != nil {
			log.Printf("Could not extract labels from %s: %v", v.name.hash, err)
			continue
		}
		if s.debugLogs {
			log.Printf("%v is now %s + (%v)\n", v.name.key, name, labels)
		}
		ts = append(ts, &monitoringpb.TimeSeries{
			Metric: &metricpb.Metric{
				Type:   path.Join("custom.googleapis.com", "go-metrics", name),
				Labels: labels,
			},
			MetricKind: metric.MetricDescriptor_GAUGE,
			Resource:   resource,
			Points: []*monitoringpb.Point{
				&monitoringpb.Point{
					Interval: &monitoringpb.TimeInterval{
						EndTime: &googlepb.Timestamp{
							Seconds: end.Unix(),
						},
					},
					Value: &monitoringpb.TypedValue{
						Value: &monitoringpb.TypedValue_DoubleValue{
							DoubleValue: v.value,
						},
					},
				},
			},
		})
	}

	for _, v := range rGauges {
		name, labels, err := v.name.labelMap(s.extractor, "gauge")
		if err != nil {
			log.Printf("Could not extract labels from %s: %v", v.name.hash, err)
			continue
		}
		if s.debugLogs {
			log.Printf("%v is now %s + (%v)\n", v.name.key, name, labels)
		}
		ts = append(ts, &monitoringpb.TimeSeries{
			Metric: &metricpb.Metric{
				Type:   path.Join("custom.googleapis.com", "go-metrics", name),
				Labels: labels,
			},
			MetricKind: metric.MetricDescriptor_GAUGE,
			Resource:   resource,
			Points: []*monitoringpb.Point{
				&monitoringpb.Point{
					Interval: &monitoringpb.TimeInterval{
						EndTime: &googlepb.Timestamp{
							Seconds: end.Unix(),
						},
					},
					Value: &monitoringpb.TypedValue{
						Value: &monitoringpb.TypedValue_DoubleValue{
							DoubleValue: float64(v.value),
						},
					},
				},
			},
		})
	}

	for _, v := range rHistograms {
		name, labels, err := v.name.labelMap(s.extractor, "histogram")
		if err != nil {
			log.Printf("Could not extract labels from %s: %v", v.name.hash, err)
			continue
		}
		if s.debugLogs {
			log.Printf("%v is now %s + (%v)\n", v.name.key, name, labels)
		}

		var count int64
		count = 0
		for _, i := range v.counts {
			count += int64(i)
		}

		ts = append(ts, &monitoringpb.TimeSeries{
			Metric: &metricpb.Metric{
				Type:   path.Join("custom.googleapis.com", "go-metrics", name),
				Labels: labels,
			},
			MetricKind: metric.MetricDescriptor_CUMULATIVE,
			Resource:   resource,
			Points: []*monitoringpb.Point{
				&monitoringpb.Point{
					Interval: &monitoringpb.TimeInterval{
						StartTime: &googlepb.Timestamp{
							Seconds: s.firstTime.Unix(),
						},
						EndTime: &googlepb.Timestamp{
							Seconds: end.Unix(),
						},
					},
					Value: &monitoringpb.TypedValue{
						Value: &monitoringpb.TypedValue_DistributionValue{
							DistributionValue: &distributionpb.Distribution{
								BucketOptions: &distributionpb.Distribution_BucketOptions{
									Options: &distributionpb.Distribution_BucketOptions_ExplicitBuckets{
										ExplicitBuckets: &distributionpb.Distribution_BucketOptions_Explicit{
											Bounds: v.buckets,
										},
									},
								},
								BucketCounts: v.counts,
								Count:        count,
							},
						},
					},
				},
			},
		})
	}

	if s.client == nil {
		return
	}

	for i := 0; i < len(ts); i += 200 {
		end := i + 200

		if end > len(ts) {
			end = len(ts)
		}

		req := &monitoringpb.CreateTimeSeriesRequest{
			Name:       fmt.Sprintf("projects/%s", s.taskInfo.ProjectID),
			TimeSeries: ts[i:end],
		}
		err := s.client.CreateTimeSeries(ctx, req)

		if err != nil {
			log.Printf("Failed to write time series data: %v", err)
			if s.debugLogs {
				for i, a := range req.TimeSeries {
					log.Printf("request timeseries[%d]: %v", i, a)
				}
			}
		}
	}
}

// A Gauge should retain the last value it is set to.
func (s *Sink) SetGauge(key []string, val float32) {
	s.SetGaugeWithLabels(key, val, nil)
}

// A Gauge should retain the last value it is set to.
func (s *Sink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	n := newSeries(key, labels)

	g := &gauge{
		name:  n,
		value: val,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.gauges[n.hash] = g
}

// Should emit a Key/Value pair for each call.
func (s *Sink) EmitKey(key []string, val float32) {
	// EmitKey is not implemented for stackdriver
}

// Counters should accumulate values.
func (s *Sink) IncrCounter(key []string, val float32) {
	s.IncrCounterWithLabels(key, val, nil)
}

// Counters should accumulate values.
func (s *Sink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	n := newSeries(key, labels)

	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.counters[n.hash]
	if ok {
		c.value += float64(val)
	} else {
		s.counters[n.hash] = &counter{
			name:  n,
			value: float64(val),
		}
	}
}

// Samples are for timing information, where quantiles are used.
func (s *Sink) AddSample(key []string, val float32) {
	s.AddSampleWithLabels(key, val, nil)
}

// Samples are for timing information, where quantiles are used.
func (s *Sink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	n := newSeries(key, labels)

	s.mu.Lock()
	defer s.mu.Unlock()

	h, ok := s.histograms[n.hash]
	if ok {
		h.addSample(val)
	} else {
		h = newHistogram(n, key, s.bucketer)
		h.addSample(val)
		s.histograms[n.hash] = h
	}
}

var _ metrics.MetricSink = (*Sink)(nil)

// Series holds the naming for a timeseries metric.
type series struct {
	key    []string
	labels []metrics.Label
	hash   string
}

var forbiddenChars = regexp.MustCompile("[ .=\\-/]")

func newSeries(key []string, labels []metrics.Label) *series {
	hash := strings.Join(key, "_")
	hash = forbiddenChars.ReplaceAllString(hash, "_")
	for _, label := range labels {
		hash += fmt.Sprintf(";%s=%s", label.Name, label.Value)
	}

	return &series{
		key:    key,
		labels: labels,
		hash:   hash,
	}
}

func (s *series) labelMap(extractor ExtractLabelsFn, kind string) (string, map[string]string, error) {
	// extract labels from the key
	key, extractedLabels, err := extractor(s.key, kind)
	if err != nil {
		return "", nil, err
	}

	name := strings.Join(key, "_")
	name = forbiddenChars.ReplaceAllString(name, "_")

	o := make(map[string]string, len(s.labels))
	for _, v := range s.labels {
		o[v.Name] = v.Value
	}
	for _, v := range extractedLabels {
		o[v.Name] = v.Value
	}
	return name, o, nil
}

// https://cloud.google.com/monitoring/api/ref_v3/rest/v3/TimeSeries#point
type gauge struct {
	name  *series
	value float32
}

// https://cloud.google.com/monitoring/api/ref_v3/rest/v3/TimeSeries#point
type counter struct {
	name  *series
	value float64
}

// https://cloud.google.com/monitoring/api/ref_v3/rest/v3/TimeSeries#distribution
type histogram struct {
	name    *series
	buckets []float64
	counts  []int64
}

func newHistogram(name *series, key []string, bucketer BucketFn) *histogram {
	h := &histogram{
		name:    name,
		buckets: bucketer(key),
	}
	h.counts = make([]int64, len(h.buckets))
	return h
}

func (h *histogram) addSample(val float32) {
	for i := range h.buckets {
		if float64(val) <= h.buckets[i] {
			h.counts[i]++
			return
		}
	}

	h.counts[len(h.buckets)-1]++
}
