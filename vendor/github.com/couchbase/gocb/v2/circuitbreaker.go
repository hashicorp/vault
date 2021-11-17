package gocb

import "time"

// CircuitBreakerCallback is the callback used by the circuit breaker to determine if an error should count toward
// the circuit breaker failure count.
type CircuitBreakerCallback func(error) bool

// CircuitBreakerConfig are the settings for configuring circuit breakers.
type CircuitBreakerConfig struct {
	Disabled                 bool
	VolumeThreshold          int64
	ErrorThresholdPercentage float64
	SleepWindow              time.Duration
	RollingWindow            time.Duration
	CompletionCallback       CircuitBreakerCallback
	CanaryTimeout            time.Duration
}
