package statsd

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type (
	countsMap         map[string]*countMetric
	gaugesMap         map[string]*gaugeMetric
	setsMap           map[string]*setMetric
	bufferedMetricMap map[string]*histogramMetric
)

// bufferedMetricContexts represent the contexts for Histograms, Distributions
// and Timing. Since those 3 metric types behave the same way and are sampled
// with the same type they're represented by the same class.
type bufferedMetricContexts struct {
	nbContext int32
	mutex     sync.RWMutex
	values    bufferedMetricMap
	newMetric func(string, float64, string) *bufferedMetric
}

func newBufferedContexts(newMetric func(string, float64, string) *bufferedMetric) bufferedMetricContexts {
	return bufferedMetricContexts{
		values:    bufferedMetricMap{},
		newMetric: newMetric,
	}
}

func (bc *bufferedMetricContexts) flush(metrics []metric) []metric {
	bc.mutex.Lock()
	values := bc.values
	bc.values = bufferedMetricMap{}
	bc.mutex.Unlock()

	for _, d := range values {
		metrics = append(metrics, d.flushUnsafe())
	}
	atomic.AddInt32(&bc.nbContext, int32(len(values)))
	return metrics
}

func (bc *bufferedMetricContexts) sample(name string, value float64, tags []string) error {
	context, stringTags := getContextAndTags(name, tags)
	bc.mutex.RLock()
	if v, found := bc.values[context]; found {
		v.sample(value)
		bc.mutex.RUnlock()
		return nil
	}
	bc.mutex.RUnlock()

	bc.mutex.Lock()
	bc.values[context] = bc.newMetric(name, value, stringTags)
	bc.mutex.Unlock()
	return nil
}

func (bc *bufferedMetricContexts) resetAndGetNbContext() int32 {
	return atomic.SwapInt32(&bc.nbContext, 0)
}

type aggregator struct {
	nbContextGauge int32
	nbContextCount int32
	nbContextSet   int32

	countsM sync.RWMutex
	gaugesM sync.RWMutex
	setsM   sync.RWMutex

	gauges        gaugesMap
	counts        countsMap
	sets          setsMap
	histograms    bufferedMetricContexts
	distributions bufferedMetricContexts
	timings       bufferedMetricContexts

	closed chan struct{}
	exited chan struct{}

	client *Client
}

type aggregatorMetrics struct {
	nbContext             int32
	nbContextGauge        int32
	nbContextCount        int32
	nbContextSet          int32
	nbContextHistogram    int32
	nbContextDistribution int32
	nbContextTiming       int32
}

func newAggregator(c *Client) *aggregator {
	return &aggregator{
		client:        c,
		counts:        countsMap{},
		gauges:        gaugesMap{},
		sets:          setsMap{},
		histograms:    newBufferedContexts(newHistogramMetric),
		distributions: newBufferedContexts(newDistributionMetric),
		timings:       newBufferedContexts(newTimingMetric),
		closed:        make(chan struct{}),
		exited:        make(chan struct{}),
	}
}

func (a *aggregator) start(flushInterval time.Duration) {
	ticker := time.NewTicker(flushInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				a.sendMetrics()
			case <-a.closed:
				close(a.exited)
				return
			}
		}
	}()
}

func (a *aggregator) sendMetrics() {
	for _, m := range a.flushMetrics() {
		a.client.send(m)
	}
}

func (a *aggregator) stop() {
	close(a.closed)
	<-a.exited
	a.sendMetrics()
}

func (a *aggregator) flushTelemetryMetrics() *aggregatorMetrics {
	if a == nil {
		return nil
	}

	am := &aggregatorMetrics{
		nbContextGauge:        atomic.SwapInt32(&a.nbContextGauge, 0),
		nbContextCount:        atomic.SwapInt32(&a.nbContextCount, 0),
		nbContextSet:          atomic.SwapInt32(&a.nbContextSet, 0),
		nbContextHistogram:    a.histograms.resetAndGetNbContext(),
		nbContextDistribution: a.distributions.resetAndGetNbContext(),
		nbContextTiming:       a.timings.resetAndGetNbContext(),
	}

	am.nbContext = am.nbContextGauge + am.nbContextCount + am.nbContextSet + am.nbContextHistogram + am.nbContextDistribution + am.nbContextTiming
	return am
}

func (a *aggregator) flushMetrics() []metric {
	metrics := []metric{}

	// We reset the values to avoid sending 'zero' values for metrics not
	// sampled during this flush interval

	a.setsM.Lock()
	sets := a.sets
	a.sets = setsMap{}
	a.setsM.Unlock()

	for _, s := range sets {
		metrics = append(metrics, s.flushUnsafe()...)
	}

	a.gaugesM.Lock()
	gauges := a.gauges
	a.gauges = gaugesMap{}
	a.gaugesM.Unlock()

	for _, g := range gauges {
		metrics = append(metrics, g.flushUnsafe())
	}

	a.countsM.Lock()
	counts := a.counts
	a.counts = countsMap{}
	a.countsM.Unlock()

	for _, c := range counts {
		metrics = append(metrics, c.flushUnsafe())
	}

	metrics = a.histograms.flush(metrics)
	metrics = a.distributions.flush(metrics)
	metrics = a.timings.flush(metrics)

	atomic.AddInt32(&a.nbContextCount, int32(len(counts)))
	atomic.AddInt32(&a.nbContextGauge, int32(len(gauges)))
	atomic.AddInt32(&a.nbContextSet, int32(len(sets)))
	return metrics
}

func getContext(name string, tags []string) string {
	return name + ":" + strings.Join(tags, tagSeparatorSymbol)
}

func getContextAndTags(name string, tags []string) (string, string) {
	stringTags := strings.Join(tags, tagSeparatorSymbol)
	return name + ":" + stringTags, stringTags
}

func (a *aggregator) count(name string, value int64, tags []string) error {
	context := getContext(name, tags)
	a.countsM.RLock()
	if count, found := a.counts[context]; found {
		count.sample(value)
		a.countsM.RUnlock()
		return nil
	}
	a.countsM.RUnlock()

	a.countsM.Lock()
	a.counts[context] = newCountMetric(name, value, tags)
	a.countsM.Unlock()
	return nil
}

func (a *aggregator) gauge(name string, value float64, tags []string) error {
	context := getContext(name, tags)
	a.gaugesM.RLock()
	if gauge, found := a.gauges[context]; found {
		gauge.sample(value)
		a.gaugesM.RUnlock()
		return nil
	}
	a.gaugesM.RUnlock()

	gauge := newGaugeMetric(name, value, tags)

	a.gaugesM.Lock()
	a.gauges[context] = gauge
	a.gaugesM.Unlock()
	return nil
}

func (a *aggregator) set(name string, value string, tags []string) error {
	context := getContext(name, tags)
	a.setsM.RLock()
	if set, found := a.sets[context]; found {
		set.sample(value)
		a.setsM.RUnlock()
		return nil
	}
	a.setsM.RUnlock()

	a.setsM.Lock()
	a.sets[context] = newSetMetric(name, value, tags)
	a.setsM.Unlock()
	return nil
}

func (a *aggregator) histogram(name string, value float64, tags []string) error {
	return a.histograms.sample(name, value, tags)
}

func (a *aggregator) distribution(name string, value float64, tags []string) error {
	return a.distributions.sample(name, value, tags)
}

func (a *aggregator) timing(name string, value float64, tags []string) error {
	return a.timings.sample(name, value, tags)
}
