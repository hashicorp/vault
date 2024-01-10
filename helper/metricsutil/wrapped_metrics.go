// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package metricsutil

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/namespace"
)

// ClusterMetricSink serves as a shim around go-metrics
// and inserts a "cluster" label.
//
// It also provides a mechanism to limit the cardinality of the labels on a gauge
// (at each reporting interval, which isn't sufficient if there is variability in which
// labels are the top N) and a backoff mechanism for gauge computation.
type ClusterMetricSink struct {
	// ClusterName is either the cluster ID, or a name provided
	// in the telemetry configuration stanza.
	//
	// Because it may be set after the Core is initialized, we need
	// to protect against concurrent access.
	ClusterName atomic.Value

	MaxGaugeCardinality int
	GaugeInterval       time.Duration

	// Sink is the go-metrics instance to send to.
	Sink metrics.MetricSink

	// Constants that are helpful for metrics within the metrics sink
	TelemetryConsts TelemetryConstConfig
}

type TelemetryConstConfig struct {
	LeaseMetricsEpsilon              time.Duration
	NumLeaseMetricsTimeBuckets       int
	LeaseMetricsNameSpaceLabels      bool
	RollbackMetricsIncludeMountPoint bool
}

type Metrics interface {
	SetGaugeWithLabels(key []string, val float32, labels []Label)
	IncrCounterWithLabels(key []string, val float32, labels []Label)
	AddSampleWithLabels(key []string, val float32, labels []Label)
	AddDurationWithLabels(key []string, d time.Duration, labels []Label)
	MeasureSinceWithLabels(key []string, start time.Time, labels []Label)
}

var _ Metrics = &ClusterMetricSink{}

// SinkWrapper implements `metricsutil.Metrics` using an instance of
// armon/go-metrics `MetricSink` as the underlying implementation.
type SinkWrapper struct {
	metrics.MetricSink
}

func (s SinkWrapper) AddDurationWithLabels(key []string, d time.Duration, labels []Label) {
	val := float32(d) / float32(time.Millisecond)
	s.MetricSink.AddSampleWithLabels(key, val, labels)
}

func (s SinkWrapper) MeasureSinceWithLabels(key []string, start time.Time, labels []Label) {
	elapsed := time.Now().Sub(start)
	val := float32(elapsed) / float32(time.Millisecond)
	s.MetricSink.AddSampleWithLabels(key, val, labels)
}

var _ Metrics = SinkWrapper{}

// Convenience alias
type Label = metrics.Label

func (m *ClusterMetricSink) SetGauge(key []string, val float32) {
	m.Sink.SetGaugeWithLabels(key, val, []Label{{"cluster", m.ClusterName.Load().(string)}})
}

func (m *ClusterMetricSink) SetGaugeWithLabels(key []string, val float32, labels []Label) {
	m.Sink.SetGaugeWithLabels(key, val,
		append(labels, Label{"cluster", m.ClusterName.Load().(string)}))
}

func (m *ClusterMetricSink) IncrCounterWithLabels(key []string, val float32, labels []Label) {
	m.Sink.IncrCounterWithLabels(key, val,
		append(labels, Label{"cluster", m.ClusterName.Load().(string)}))
}

func (m *ClusterMetricSink) AddSample(key []string, val float32) {
	m.Sink.AddSampleWithLabels(key, val, []Label{{"cluster", m.ClusterName.Load().(string)}})
}

func (m *ClusterMetricSink) AddSampleWithLabels(key []string, val float32, labels []Label) {
	m.Sink.AddSampleWithLabels(key, val,
		append(labels, Label{"cluster", m.ClusterName.Load().(string)}))
}

func (m *ClusterMetricSink) AddDurationWithLabels(key []string, d time.Duration, labels []Label) {
	val := float32(d) / float32(time.Millisecond)
	m.AddSampleWithLabels(key, val, labels)
}

func (m *ClusterMetricSink) MeasureSinceWithLabels(key []string, start time.Time, labels []Label) {
	elapsed := time.Now().Sub(start)
	val := float32(elapsed) / float32(time.Millisecond)
	m.AddSampleWithLabels(key, val, labels)
}

// BlackholeSink is a default suitable for use in unit tests.
func BlackholeSink() *ClusterMetricSink {
	conf := metrics.DefaultConfig("")
	conf.EnableRuntimeMetrics = false
	sink, _ := metrics.New(conf, &metrics.BlackholeSink{})
	cms := &ClusterMetricSink{
		ClusterName: atomic.Value{},
		Sink:        sink,
	}
	cms.ClusterName.Store("")
	return cms
}

func NewClusterMetricSink(clusterName string, sink metrics.MetricSink) *ClusterMetricSink {
	cms := &ClusterMetricSink{
		ClusterName:     atomic.Value{},
		Sink:            sink,
		TelemetryConsts: TelemetryConstConfig{},
	}
	cms.ClusterName.Store(clusterName)
	return cms
}

// SetDefaultClusterName changes the cluster name from its default value,
// if it has not previously been configured.
func (m *ClusterMetricSink) SetDefaultClusterName(clusterName string) {
	// This is not a true compare-and-swap, but it should be
	// consistent enough for normal uses
	if m.ClusterName.Load().(string) == "" {
		m.ClusterName.Store(clusterName)
	}
}

// NamespaceLabel creates a metrics label for the given
// Namespace: root is "root"; others are path with the
// final '/' removed.
func NamespaceLabel(ns *namespace.Namespace) metrics.Label {
	switch {
	case ns == nil:
		return metrics.Label{"namespace", "root"}
	case ns.ID == namespace.RootNamespaceID:
		return metrics.Label{"namespace", "root"}
	default:
		return metrics.Label{
			"namespace",
			strings.Trim(ns.Path, "/"),
		}
	}
}
