package metricsutil

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

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

func startSimulatedTime() *SimulatedTime {
	s := &SimulatedTime{
		now:           time.Now(),
		tickerBarrier: make(chan *SimulatedTicker, 1),
	}
	timeNow = s.Now
	newTicker = s.NewTicker
	return s
}

func stopSimulatedTime() {
	timeNow = time.Now
	newTicker = time.NewTicker
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
	s := startSimulatedTime()
	defer stopSimulatedTime()

	c := newSimulatedCollector()

	sink := BlackholeSink()
	sink.GaugeInterval = 2 * time.Hour

	p, err := sink.NewGaugeCollectionProcess(
		[]string{"example", "count"},
		[]Label{{"gauge", "test"}},
		c.EmptyCollectionFunction,
		log.Default(),
	)
	if err != nil {
		t.Fatalf("Error creating collection process: %v", err)
	}

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

	p.Stop <- true

}
