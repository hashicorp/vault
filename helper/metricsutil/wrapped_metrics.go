package metricsutil

import (
	"time"

	metrics "github.com/armon/go-metrics"
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
	ClusterName string

	MaxGaugeCardinality int
	GaugeInterval       time.Duration

	// Sink is the go-metrics sink to send to
	Sink metrics.MetricSink
}

// Convenience alias
type Label = metrics.Label

func (m *ClusterMetricSink) SetGaugeWithLabels(key []string, val float32, labels []Label) {
	m.Sink.SetGaugeWithLabels(key, val,
		append(labels, Label{"cluster", m.ClusterName}))
}

func (m *ClusterMetricSink) IncrCounterWithLabels(key []string, val float32, labels []Label) {
	m.Sink.IncrCounterWithLabels(key, val,
		append(labels, Label{"cluster", m.ClusterName}))
}

func (m *ClusterMetricSink) AddSampleWithLabels(key []string, val float32, labels []Label) {
	m.Sink.AddSampleWithLabels(key, val,
		append(labels, Label{"cluster", m.ClusterName}))
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
	return &ClusterMetricSink{
		ClusterName: "",
		Sink:        &metrics.BlackholeSink{},
	}
}

// SetDefaultClusterName changes the cluster name from its default value,
// if it has not previously been configured.
func (m *ClusterMetricSink) SetDefaultClusterName(clusterName string) {
	if m.ClusterName == "" {
		m.ClusterName = clusterName
	}
}
