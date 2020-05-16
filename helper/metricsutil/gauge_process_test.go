package metricsutil

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
)

// SimulatedTime maintains a virtual clock so the test isn't
// dependent upon real time.
// Unfortunately there is no way to run these tests in parallel
// since they rely on the same global timeNow function.
type SimulatedTime struct {
	now           time.Time
	tickerBarrier chan *SimulatedTicker
}

type SimulatedTicker struct {
	ticker   *time.Ticker
	duration time.Duration
	sender   chan time.Time
}

func (s *SimulatedTime) Now() time.Time {
	return s.now
}

func (s *SimulatedTime) NewTicker(d time.Duration) *time.Ticker {
	// Create a real ticker, but set its duration to an amount that will never fire for real.
	// We'll inject times into the channel directly.
	replacementChannel := make(chan time.Time)
	t := time.NewTicker(1000 * time.Hour)
	t.C = replacementChannel
	s.tickerBarrier <- &SimulatedTicker{t, d, replacementChannel}
	return t
}

func (s *SimulatedTime) waitForTicker(t *testing.T) *SimulatedTicker {
	// System under test should create a ticker within 100ms,
	// wait for it to show up or else fail the test.
	timeout := time.After(100 * time.Millisecond)
	select {
	case <-timeout:
		t.Fatal("Timeout waiting for ticker creation.")
		return nil
	case t := <-s.tickerBarrier:
		return t
	}
}

func (s *SimulatedTime) allowTickers(n int) {
	s.tickerBarrier = make(chan *SimulatedTicker, n)
}

func startSimulatedTime() *SimulatedTime {
	s := &SimulatedTime{
		now:           time.Now(),
		tickerBarrier: make(chan *SimulatedTicker, 1),
	}
	return s
}

type SimulatedCollector struct {
	numCalls    uint32
	callBarrier chan uint32
}

func newSimulatedCollector() *SimulatedCollector {
	return &SimulatedCollector{
		numCalls:    0,
		callBarrier: make(chan uint32, 1),
	}
}

func (s *SimulatedCollector) waitForCall(t *testing.T) {
	timeout := time.After(100 * time.Millisecond)
	select {
	case <-timeout:
		t.Fatal("Timeout waiting for call to collection function.")
		return
	case <-s.callBarrier:
		return
	}
}

func (s *SimulatedCollector) EmptyCollectionFunction(ctx context.Context) ([]GaugeLabelValues, error) {
	atomic.AddUint32(&s.numCalls, 1)
	s.callBarrier <- s.numCalls
	return []GaugeLabelValues{}, nil
}

func TestGauge_StartDelay(t *testing.T) {
	// Work through an entire startup sequence, up to collecting
	// the first batch of gauges.
	s := startSimulatedTime()
	c := newSimulatedCollector()

	sink := BlackholeSink()
	sink.GaugeInterval = 2 * time.Hour

	p, err := sink.newGcpWithClock(
		[]string{"example", "count"},
		[]Label{{"gauge", "test"}},
		c.EmptyCollectionFunction,
		log.Default(),
		s,
	)
	if err != nil {
		t.Fatalf("Error creating collection process: %v", err)
	}
	go p.Run()

	delayTicker := s.waitForTicker(t)
	if delayTicker.duration > sink.GaugeInterval {
		t.Errorf("Delayed start %v is more than interval %v.",
			delayTicker.duration, sink.GaugeInterval)
	}
	if c.numCalls > 0 {
		t.Error("Collection function has been called")
	}

	// Signal the end of delay, then another ticker should start
	delayTicker.sender <- time.Now()

	intervalTicker := s.waitForTicker(t)
	if intervalTicker.duration != sink.GaugeInterval {
		t.Errorf("Ticker duration is %v, expected %v",
			intervalTicker.duration, sink.GaugeInterval)
	}
	if c.numCalls > 0 {
		t.Error("Collection function has been called")
	}

	// Time's up, ensure the collection function is executed.
	intervalTicker.sender <- time.Now()
	c.waitForCall(t)
	if c.numCalls != 1 {
		t.Errorf("Collection function called %v times, expected %v.", c.numCalls, 1)
	}

	p.Stop <- struct{}{}
}

func waitForStopped(t *testing.T, p *GaugeCollectionProcess) {
	timeout := time.After(100 * time.Millisecond)
	select {
	case <-timeout:
		t.Fatal("Timeout waiting for process to stop.")
	case <-p.stopped:
		return
	}
}

func TestGauge_StoppedDuringInitialDelay(t *testing.T) {
	// Stop the process before it gets into its main loop
	s := startSimulatedTime()
	c := newSimulatedCollector()

	sink := BlackholeSink()
	sink.GaugeInterval = 2 * time.Hour

	p, err := sink.newGcpWithClock(
		[]string{"example", "count"},
		[]Label{{"gauge", "test"}},
		c.EmptyCollectionFunction,
		log.Default(),
		s,
	)
	if err != nil {
		t.Fatalf("Error creating collection process: %v", err)
	}
	go p.Run()

	// Stop during the initial delay, check that goroutine exits
	s.waitForTicker(t)
	p.Stop <- struct{}{}
	waitForStopped(t, p)
}

func TestGauge_StoppedAfterInitialDelay(t *testing.T) {
	// Stop the process during its main loop
	s := startSimulatedTime()
	c := newSimulatedCollector()

	sink := BlackholeSink()
	sink.GaugeInterval = 2 * time.Hour

	p, err := sink.newGcpWithClock(
		[]string{"example", "count"},
		[]Label{{"gauge", "test"}},
		c.EmptyCollectionFunction,
		log.Default(),
		s,
	)
	if err != nil {
		t.Fatalf("Error creating collection process: %v", err)
	}
	go p.Run()

	// Get through initial delay, wait for interval ticker
	delayTicker := s.waitForTicker(t)
	delayTicker.sender <- time.Now()

	s.waitForTicker(t)
	p.Stop <- struct{}{}
	waitForStopped(t, p)
}

func TestGauge_Backoff(t *testing.T) {
	s := startSimulatedTime()
	s.allowTickers(100)

	c := newSimulatedCollector()

	sink := BlackholeSink()
	sink.GaugeInterval = 2 * time.Hour

	threshold := time.Duration(int(sink.GaugeInterval) / 100)
	f := func(ctx context.Context) ([]GaugeLabelValues, error) {
		atomic.AddUint32(&c.numCalls, 1)
		// Move time forward by more than 1% of the gauge interval
		s.now = s.now.Add(threshold).Add(time.Second)
		c.callBarrier <- c.numCalls
		return []GaugeLabelValues{}, nil
	}

	p, err := sink.newGcpWithClock(
		[]string{"example", "count"},
		[]Label{{"gauge", "test"}},
		f,
		log.Default(),
		s,
	)
	if err != nil {
		t.Fatalf("Error creating collection process: %v", err)
	}
	// Do not run, we'll just going to call an internal function.
	p.collectAndFilterGauges()

	if p.currentInterval != 2*p.originalInterval {
		t.Errorf("Current interval is %v, should be 2x%v.",
			p.currentInterval,
			p.originalInterval)
	}
}

func waitForDone(t *testing.T,
	tick chan<- time.Time,
	done <-chan struct{},
) int {
	timeout := time.After(100 * time.Millisecond)

	numTicks := 0
	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for metrics to be sent.")
		case tick <- time.Now():
			numTicks += 1
		case <-done:
			return numTicks
		}
	}
}

func TestGauge_MaximumMeasurements(t *testing.T) {
	s := startSimulatedTime()
	c := newSimulatedCollector()

	// Long bucket time == low chance of crossing interval
	inmemSink := metrics.NewInmemSink(
		1000000*time.Hour,
		2000000*time.Hour)

	sink := &ClusterMetricSink{
		ClusterName:         "test",
		MaxGaugeCardinality: 500,
		GaugeInterval:       2 * time.Hour,
		Sink:                inmemSink,
	}

	// Create a report larger than the default limit
	excessGauges := 100
	values := make([]GaugeLabelValues, sink.MaxGaugeCardinality+excessGauges)
	for i := range values {
		values[i].Labels = []Label{
			{"test", "true"},
			{"which", fmt.Sprintf("%v", i)},
		}
		values[i].Value = float32(i + 1)
	}
	rand.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	f := func(ctx context.Context) ([]GaugeLabelValues, error) {
		atomic.AddUint32(&c.numCalls, 1)

		// Move time forward by 0.5% of the gauge interval
		timeUsed := time.Duration(int(sink.GaugeInterval) / 200)
		s.now = s.now.Add(timeUsed)

		c.callBarrier <- c.numCalls
		return values, nil
	}

	p, err := sink.newGcpWithClock(
		[]string{"example", "count"},
		[]Label{{"gauge", "test"}},
		f,
		log.Default(),
		s,
	)
	if err != nil {
		t.Fatalf("Error creating collection process: %v", err)
	}

	// This needs a ticker in order to do its thing,
	// so run it in the background and we'll send the ticks
	// from here.
	done := make(chan struct{}, 1)
	go func() {
		p.collectAndFilterGauges()
		close(done)
	}()

	sendTicker := s.waitForTicker(t)
	numTicksSent := waitForDone(t, sendTicker.sender, done)

	// 500 items, one delay after after each 25, means that
	// 19 ticks are consumed, so 19 or 20 must be sent.
	expectedTicks := sink.MaxGaugeCardinality/25 - 1
	if numTicksSent < expectedTicks || numTicksSent > expectedTicks+1 {
		t.Errorf("Number of ticks = %v, expected %v.", numTicksSent, expectedTicks)
	}

	// If we start close to the end of an interval, metrics will
	// be split across two buckets.
	intervals := inmemSink.Data()
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}

	if len(intervals[0].Gauges) != sink.MaxGaugeCardinality {
		t.Errorf("Found %v gauges, expected %v.",
			len(intervals[0].Gauges),
			sink.MaxGaugeCardinality)
	}

	minVal := float32(excessGauges)
	for _, v := range intervals[0].Gauges {
		if v.Value < minVal {
			t.Errorf("Gauge %v with value %v should not have been included.", v.Labels, v.Value)
			break
		}
	}
}
