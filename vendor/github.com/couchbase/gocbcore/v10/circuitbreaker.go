package gocbcore

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	circuitBreakerStateDisabled uint32 = iota
	circuitBreakerStateClosed
	circuitBreakerStateHalfOpen
	circuitBreakerStateOpen
)

type circuitBreaker interface {
	AllowsRequest() bool
	MarkSuccessful()
	MarkFailure()
	State() uint32
	Reset()
	CanaryTimeout() time.Duration
	CompletionCallback(error) bool
}

// CircuitBreakerCallback is the callback used by the circuit breaker to determine if an error should count toward
// the circuit breaker failure count.
type CircuitBreakerCallback func(error) bool

// CircuitBreakerConfig is the set of configuration settings for configuring circuit breakers.
// If Disabled is set to true then a noop circuit breaker will be used, otherwise a lazy circuit
// breaker.
type CircuitBreakerConfig struct {
	Enabled                  bool
	VolumeThreshold          int64
	ErrorThresholdPercentage float64
	SleepWindow              time.Duration
	RollingWindow            time.Duration
	CompletionCallback       CircuitBreakerCallback
	CanaryTimeout            time.Duration
}

type noopCircuitBreaker struct {
}

func newNoopCircuitBreaker() *noopCircuitBreaker {
	return &noopCircuitBreaker{}
}

func (ncb *noopCircuitBreaker) AllowsRequest() bool {
	return true
}

func (ncb *noopCircuitBreaker) MarkSuccessful() {
}

func (ncb *noopCircuitBreaker) MarkFailure() {
}

func (ncb *noopCircuitBreaker) State() uint32 {
	return circuitBreakerStateDisabled
}

func (ncb *noopCircuitBreaker) Reset() {
}

func (ncb *noopCircuitBreaker) CompletionCallback(error) bool {
	return true
}

func (ncb *noopCircuitBreaker) CanaryTimeout() time.Duration {
	return 0
}

type lazyCircuitBreaker struct {
	state                    uint32
	windowStart              int64
	sleepWindow              int64
	rollingWindow            int64
	volumeThreshold          int64
	errorPercentageThreshold float64
	canaryTimeout            time.Duration
	total                    int64
	failed                   int64
	openedAt                 int64
	sendCanaryFn             func()
	completionCallback       CircuitBreakerCallback
}

func newLazyCircuitBreaker(config CircuitBreakerConfig, canaryFn func()) *lazyCircuitBreaker {
	if config.VolumeThreshold == 0 {
		config.VolumeThreshold = 20
	}
	if config.ErrorThresholdPercentage == 0 {
		config.ErrorThresholdPercentage = 50
	}
	if config.SleepWindow == 0 {
		config.SleepWindow = 5 * time.Second
	}
	if config.RollingWindow == 0 {
		config.RollingWindow = 1 * time.Minute
	}
	if config.CanaryTimeout == 0 {
		config.CanaryTimeout = 5 * time.Second
	}
	if config.CompletionCallback == nil {
		config.CompletionCallback = func(err error) bool {
			return !errors.Is(err, ErrTimeout)
		}
	}

	breaker := &lazyCircuitBreaker{
		sleepWindow:              int64(config.SleepWindow * time.Nanosecond),
		rollingWindow:            int64(config.RollingWindow * time.Nanosecond),
		volumeThreshold:          config.VolumeThreshold,
		errorPercentageThreshold: config.ErrorThresholdPercentage,
		canaryTimeout:            config.CanaryTimeout,
		sendCanaryFn:             canaryFn,
		completionCallback:       config.CompletionCallback,
	}
	breaker.Reset()

	return breaker
}

func (lcb *lazyCircuitBreaker) Reset() {
	now := time.Now().UnixNano()
	atomic.StoreUint32(&lcb.state, circuitBreakerStateClosed)
	atomic.StoreInt64(&lcb.total, 0)
	atomic.StoreInt64(&lcb.failed, 0)
	atomic.StoreInt64(&lcb.openedAt, 0)
	atomic.StoreInt64(&lcb.windowStart, now)
}

func (lcb *lazyCircuitBreaker) State() uint32 {
	return atomic.LoadUint32(&lcb.state)
}

func (lcb *lazyCircuitBreaker) AllowsRequest() bool {
	state := lcb.State()
	if state == circuitBreakerStateClosed {
		return true
	}

	elapsed := (time.Now().UnixNano() - atomic.LoadInt64(&lcb.openedAt)) > lcb.sleepWindow
	if elapsed && atomic.CompareAndSwapUint32(&lcb.state, circuitBreakerStateOpen, circuitBreakerStateHalfOpen) {
		// If we're outside of the sleep window and the circuit is open then send a canary.
		go lcb.sendCanaryFn()
	}
	return false
}

func (lcb *lazyCircuitBreaker) MarkSuccessful() {
	if atomic.CompareAndSwapUint32(&lcb.state, circuitBreakerStateHalfOpen, circuitBreakerStateClosed) {
		logDebugf("Moving circuit breaker to closed")
		lcb.Reset()
		return
	}

	lcb.maybeResetRollingWindow()
	atomic.AddInt64(&lcb.total, 1)
}

func (lcb *lazyCircuitBreaker) MarkFailure() {
	now := time.Now().UnixNano()
	if atomic.CompareAndSwapUint32(&lcb.state, circuitBreakerStateHalfOpen, circuitBreakerStateOpen) {
		logDebugf("Moving circuit breaker from half open to open")
		atomic.StoreInt64(&lcb.openedAt, now)
		return
	}

	lcb.maybeResetRollingWindow()
	atomic.AddInt64(&lcb.total, 1)
	atomic.AddInt64(&lcb.failed, 1)
	lcb.maybeOpenCircuit()
}

func (lcb *lazyCircuitBreaker) CanaryTimeout() time.Duration {
	return lcb.canaryTimeout
}

func (lcb *lazyCircuitBreaker) CompletionCallback(err error) bool {
	return lcb.completionCallback(err)
}

func (lcb *lazyCircuitBreaker) maybeOpenCircuit() {
	if atomic.LoadInt64(&lcb.total) < lcb.volumeThreshold {
		return
	}

	currentPercentage := (float64(atomic.LoadInt64(&lcb.failed)) / float64(atomic.LoadInt64(&lcb.total))) * 100
	if currentPercentage >= lcb.errorPercentageThreshold {
		logDebugf("Moving circuit breaker to open")
		atomic.StoreUint32(&lcb.state, circuitBreakerStateOpen)
		atomic.StoreInt64(&lcb.openedAt, time.Now().UnixNano())
	}
}

func (lcb *lazyCircuitBreaker) maybeResetRollingWindow() {
	now := time.Now().UnixNano()
	if (now - atomic.LoadInt64(&lcb.windowStart)) <= lcb.rollingWindow {
		return
	}

	atomic.StoreInt64(&lcb.windowStart, now)
	atomic.StoreInt64(&lcb.total, 0)
	atomic.StoreInt64(&lcb.failed, 0)
}
