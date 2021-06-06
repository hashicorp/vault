package metricsutil

import (
	"context"
	"math/rand"
	"sort"
	"time"

	log "github.com/hashicorp/go-hclog"
)

// This interface allows unit tests to substitute in a simulated clock.
type clock interface {
	Now() time.Time
	NewTicker(time.Duration) *time.Ticker
}

type defaultClock struct{}

func (_ defaultClock) Now() time.Time {
	return time.Now()
}

func (_ defaultClock) NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}

// GaugeLabelValues is one gauge in a set sharing a single key, that
// are measured in a batch.
type GaugeLabelValues struct {
	Labels []Label
	Value  float32
}

// GaugeCollector is a callback function that returns an unfiltered
// set of label-value pairs. It may be cancelled if it takes too long.
type GaugeCollector = func(context.Context) ([]GaugeLabelValues, error)

// collectionBound is a hard limit on how long a collection process
// may take, as a fraction of the current interval.
const collectionBound = 0.02

// collectionTarget is a soft limit; if exceeded, the collection interval
// will be doubled.
const collectionTarget = 0.01

// A GaugeCollectionProcess is responsible for one particular gauge metric.
// It handles a delay on initial startup; limiting the cardinality; and
// exponential backoff on the requested interval.
type GaugeCollectionProcess struct {
	stop    chan struct{}
	stopped chan struct{}

	// gauge name
	key []string
	// labels to use when reporting
	labels []Label

	// callback function
	collector GaugeCollector

	// destination for metrics
	sink   *ClusterMetricSink
	logger log.Logger

	// time between collections
	originalInterval time.Duration
	currentInterval  time.Duration
	ticker           *time.Ticker

	// time source
	clock clock
}

// NewGaugeCollectionProcess creates a new collection process for the callback
// function given as an argument, and starts it running.
// A label should be provided for metrics *about* this collection process.
//
// The Run() method must be called to start the process.
func (m *ClusterMetricSink) NewGaugeCollectionProcess(
	key []string,
	id []Label,
	collector GaugeCollector,
	logger log.Logger,
) (*GaugeCollectionProcess, error) {
	return m.newGaugeCollectionProcessWithClock(
		key,
		id,
		collector,
		logger,
		defaultClock{},
	)
}

// test version allows an alternative clock implementation
func (m *ClusterMetricSink) newGaugeCollectionProcessWithClock(
	key []string,
	id []Label,
	collector GaugeCollector,
	logger log.Logger,
	clock clock,
) (*GaugeCollectionProcess, error) {
	process := &GaugeCollectionProcess{
		stop:             make(chan struct{}, 1),
		stopped:          make(chan struct{}, 1),
		key:              key,
		labels:           id,
		collector:        collector,
		sink:             m,
		originalInterval: m.GaugeInterval,
		currentInterval:  m.GaugeInterval,
		logger:           logger,
		clock:            clock,
	}
	return process, nil
}

// delayStart randomly delays by up to one extra interval
// so that collection processes do not all run at the time time.
// If we knew all the procsses in advance, we could just schedule them
// evenly, but a new one could be added per secret engine.
func (p *GaugeCollectionProcess) delayStart() bool {
	randomDelay := time.Duration(rand.Int63n(int64(p.currentInterval)))
	// A Timer might be better, but then we'd have to simulate
	// one of those too?
	delayTick := p.clock.NewTicker(randomDelay)
	defer delayTick.Stop()

	select {
	case <-p.stop:
		return true
	case <-delayTick.C:
		break
	}
	return false
}

// resetTicker stops the old ticker and starts a new one at the current
// interval setting.
func (p *GaugeCollectionProcess) resetTicker() {
	if p.ticker != nil {
		p.ticker.Stop()
	}
	p.ticker = p.clock.NewTicker(p.currentInterval)
}

// collectAndFilterGauges executes the callback function,
// limits the cardinality, and streams the results to the metrics sink.
func (p *GaugeCollectionProcess) collectAndFilterGauges() {
	// Run for only an allotted amount of time.
	timeout := time.Duration(collectionBound * float64(p.currentInterval))
	ctx, cancel := context.WithTimeout(context.Background(),
		timeout)
	defer cancel()

	p.sink.AddDurationWithLabels([]string{"metrics", "collection", "interval"},
		p.currentInterval,
		p.labels)

	start := p.clock.Now()
	values, err := p.collector(ctx)
	end := p.clock.Now()
	duration := end.Sub(start)

	// Report how long it took to perform the operation.
	p.sink.AddDurationWithLabels([]string{"metrics", "collection"},
		duration,
		p.labels)

	// If over threshold, back off by doubling the measurement interval.
	// Currently a restart is the only way to bring it back down.
	threshold := time.Duration(collectionTarget * float64(p.currentInterval))
	if duration > threshold {
		p.logger.Warn("gauge collection time exceeded target", "target", threshold, "actual", duration, "id", p.labels)
		p.currentInterval *= 2
		p.resetTicker()
	}

	if err != nil {
		p.logger.Error("error collecting gauge", "id", p.labels, "error", err)
		p.sink.IncrCounterWithLabels([]string{"metrics", "collection", "error"},
			1,
			p.labels)
		return
	}

	// Filter to top N.
	// This does not guarantee total cardinality is <= N, but it does slow things down
	// a little if the cardinality *is* too high and the gauge needs to be disabled.
	if len(values) > p.sink.MaxGaugeCardinality {
		sort.Slice(values, func(a, b int) bool {
			return values[a].Value > values[b].Value
		})
		values = values[:p.sink.MaxGaugeCardinality]
	}

	p.streamGaugesToSink(values)
}

func (p *GaugeCollectionProcess) streamGaugesToSink(values []GaugeLabelValues) {
	// Dumping 500 metrics in one big chunk is somewhat unfriendly to UDP-based
	// transport, and to the rest of the metrics trying to get through.
	// Let's smooth things out over the course of a second.
	// 1 second / 500 = 2 ms each, so we can send 25 per 50 milliseconds.
	// That should be one or two packets.
	sendTick := p.clock.NewTicker(50 * time.Millisecond)
	defer sendTick.Stop()

	batchSize := 25
	for i, lv := range values {
		if i > 0 && i%batchSize == 0 {
			select {
			case <-p.stop:
				// because the channel is closed,
				// the main loop will successfully
				// read from p.stop too, and exit.
				return
			case <-sendTick.C:
				break
			}
		}
		p.sink.SetGaugeWithLabels(p.key, lv.Value, lv.Labels)
	}
}

// Run should be called as a goroutine.
func (p *GaugeCollectionProcess) Run() {
	defer close(p.stopped)

	// Wait a random amount of time
	stopReceived := p.delayStart()
	if stopReceived {
		return
	}

	// Create a ticker to start each cycle
	p.resetTicker()

	// Loop until we get a signal to stop
	for {
		select {
		case <-p.ticker.C:
			p.collectAndFilterGauges()
		case <-p.stop:
			// Can't use defer because this might
			// not be the original ticker.
			p.ticker.Stop()
			return
		}
	}
}

// Stop the collection process
func (p *GaugeCollectionProcess) Stop() {
	close(p.stop)
}
