package permits

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/armon/go-metrics"
)

func TestInstrumentedPermitPool(t *testing.T) {
	conf := metrics.DefaultConfig("test-permits")
	conf.EnableHostname = false
	conf.EnableHostnameLabel = false
	sink := metrics.NewInmemSink(time.Hour, time.Hour)
	metrics.NewGlobal(conf, sink)

	permitsUsedMetricName := "test-permits.a.pool.permits"
	maxPermitMetricName := "test-permits.a.pool.permits-limit"
	sinkAssert := newSinkAssert(t, sink)

	pool := NewInstrumentedPermitPool(658, "a", "pool")
	sinkAssert.GaugeEquals(permitsUsedMetricName, 0)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)

	pool.Acquire()
	sinkAssert.GaugeEquals(permitsUsedMetricName, 1)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)

	pool.Acquire()
	sinkAssert.GaugeEquals(permitsUsedMetricName, 2)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)

	pool.Release()
	sinkAssert.GaugeEquals(permitsUsedMetricName, 1)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)

	pool.Acquire()
	sinkAssert.GaugeEquals(permitsUsedMetricName, 2)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)

	pool.Release()
	sinkAssert.GaugeEquals(permitsUsedMetricName, 1)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)

	pool.Release()
	sinkAssert.GaugeEquals(permitsUsedMetricName, 0)
	sinkAssert.GaugeEquals(maxPermitMetricName, 658)
}

type sinkAssert struct {
	t    *testing.T
	sink *metrics.InmemSink
}

func newSinkAssert(t *testing.T, sink *metrics.InmemSink) *sinkAssert {
	return &sinkAssert{
		t:    t,
		sink: sink,
	}
}

func (st *sinkAssert) GaugeEquals(gaugeName string, expected interface{}) {
	data := st.sink.Data()
	require.Equal(st.t, 1, len(data), "expected 1 and only 1 interval, found %v intervals", len(data))

	interval := data[0]
	gauge, err := readGaugeFromInterval(interval, gaugeName)
	assert.NoError(st.t, err, "could not find gauge named %q", gaugeName)
	assert.InDelta(st.t, expected, gauge.Value, 0.0001, "gauge %q should be %v but was %v", gaugeName, expected, gauge.Value)
}

func readGaugeFromInterval(metric *metrics.IntervalMetrics, name string) (metrics.GaugeValue, error) {
	metric.RLock()
	defer metric.RUnlock()
	gaugesFound := []string{}
	for k, gauge := range metric.Gauges {
		if k == name {
			return gauge, nil
		}
		gaugesFound = append(gaugesFound, k)
	}
	return metrics.GaugeValue{}, fmt.Errorf("no gauge named %s found (current gauges were %+v)", name, gaugesFound)
}
