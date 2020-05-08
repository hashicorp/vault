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

var globalClusterMetrics *ClusterMetricSink

func init() {
	// Default to a black-hole sink.
	// This will be changed during server init, but there might be unit
	// tests or other cases that didn't initialize metrics.
	globalClusterMetrics = &ClusterMetricSink{
		ClusterName: "",
		Sink:        &metrics.BlackholeSink{},
	}
}

func SetGlobalSink(newGlobal *ClusterMetricSink) {
	globalClusterMetrics = newGlobal
}

// SetDefaultClusterName sets ClusterName if it was not specified in the configuration.
// At the time the metrics are set up, the cluster name may not yet have been read from
// storage (or generated for the first time.)
func SetDefaultClusterName(clusterName string) {
	if globalClusterMetrics.ClusterName == "" {
		globalClusterMetrics.ClusterName = clusterName
	}
}

func SetGaugeWithLabels(key []string, val float32, labels []Label) {
	globalClusterMetrics.SetGaugeWithLabels(key, val, labels)
}

func IncrCounterWithLabels(key []string, val float32, labels []Label) {
	globalClusterMetrics.IncrCounterWithLabels(key, val, labels)
}

func AddSampleWithLabels(key []string, val float32, labels []Label) {
	globalClusterMetrics.AddSampleWithLabels(key, val, labels)
}
