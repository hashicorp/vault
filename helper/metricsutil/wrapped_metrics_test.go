package metricsutil

import (
	"testing"
	"time"

	"github.com/armon/go-metrics"
)

func isLabelPresent(toFind Label, ls []Label) bool {
	for _, l := range ls {
		if l == toFind {
			return true
		}
	}
	return false
}

func TestClusterLabelPresent(t *testing.T) {
	// With only one test method, we can check that sending a metric
	// before SetGlobalSink doesn't crash anything.
	// If we had multiple test methods, I don't know how we would check that
	// since we can't guarantee order?

	SetGaugeWithLabels([]string{"test", "gauge"},
		1.0,
		[]Label{
			{"label", "value"},
		},
	)

	testClusterName := "test-cluster"

	// Use a ridiculously long time to minimize the chance
	// that we have to deal with more than one interval.
	// InMemSink rounds down to an interval boundary rather than
	// starting one at the time of initialization.
	inmemSink := metrics.NewInmemSink(
		1000000*time.Hour,
		2000000*time.Hour)
	clusterSink := &ClusterMetricSink{
		ClusterName: testClusterName,
		Sink:        inmemSink,
	}

	SetGlobalSink(clusterSink)

	key1 := []string{"aaa", "bbb"}
	key2 := []string{"ccc", "ddd"}
	key3 := []string{"eee", "fff"}
	labels1 := []Label{{"dim1", "val1"}}
	labels2 := []Label{{"dim2", "val2"}}
	labels3 := []Label{{"dim3", "val3"}}
	clusterLabel := Label{"cluster", testClusterName}
	expectedKey1 := "aaa.bbb;dim1=val1;cluster=" + testClusterName
	expectedKey2 := "ccc.ddd;dim2=val2;cluster=" + testClusterName
	expectedKey3 := "eee.fff;dim3=val3;cluster=" + testClusterName

	SetGaugeWithLabels(key1, 1.0, labels1)
	IncrCounterWithLabels(key2, 2.0, labels2)
	AddSampleWithLabels(key3, 3.0, labels3)

	intervals := inmemSink.Data()
	if len(intervals) > 1 {
		t.Log("Skipping test, detected interval crossing.")
		return
	}

	// Check Gauge
	g, ok := intervals[0].Gauges[expectedKey1]
	if !ok {
		t.Fatal("Key", expectedKey1, "not found in map", intervals[0].Gauges)
	}
	if g.Value != 1.0 {
		t.Error("Gauge value", g.Value, "does not match", 1.0)
	}
	if !isLabelPresent(labels1[0], g.Labels) {
		t.Error("Gauge label", g.Labels, "does not include", labels1)
	}
	if !isLabelPresent(clusterLabel, g.Labels) {
		t.Error("Gauge label", g.Labels, "does not include", clusterLabel)
	}

	// Check Counter
	c, ok := intervals[0].Counters[expectedKey2]
	if !ok {
		t.Fatal("Key", expectedKey2, "not found in map", intervals[0].Counters)
	}
	if c.Sum != 2.0 {
		t.Error("Counter value", c.Sum, "does not match", 2.0)
	}
	if !isLabelPresent(labels2[0], c.Labels) {
		t.Error("Counter label", c.Labels, "does not include", labels2)
	}
	if !isLabelPresent(clusterLabel, c.Labels) {
		t.Error("Counter label", c.Labels, "does not include", clusterLabel)
	}

	// Check Sample
	s, ok := intervals[0].Samples[expectedKey3]
	if !ok {
		t.Fatal("Key", expectedKey3, "not found in map", intervals[0].Samples)
	}
	if s.Sum != 3.0 {
		t.Error("Sample value", s.Sum, "does not match", 3.0)
	}
	if !isLabelPresent(labels3[0], s.Labels) {
		t.Error("Sample label", s.Labels, "does not include", labels3)
	}
	if !isLabelPresent(clusterLabel, s.Labels) {
		t.Error("Sample label", s.Labels, "does not include", clusterLabel)
	}

}
