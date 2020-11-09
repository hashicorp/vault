package gocbcore

import (
	"crypto/x509"
	"time"
)

type clusterAgentConfig struct {
	HTTPAddrs []string
	UserAgent string
	UseTLS    bool
	Auth      AuthProvider

	TLSRootCAProvider func() *x509.CertPool

	HTTPMaxIdleConns          int
	HTTPMaxIdleConnsPerHost   int
	HTTPIdleConnectionTimeout time.Duration

	// Volatile: Tracer API is subject to change.
	Tracer           RequestTracer
	NoRootTraceSpans bool

	DefaultRetryStrategy RetryStrategy
	CircuitBreakerConfig CircuitBreakerConfig
}

func (config *clusterAgentConfig) redacted() interface{} {
	newConfig := clusterAgentConfig{}
	newConfig = *config
	if isLogRedactionLevelFull() {
		// The slices here are still pointing at config's underlying arrays
		// so we need to make them not do that.
		newConfig.HTTPAddrs = append([]string(nil), newConfig.HTTPAddrs...)
		for i, addr := range newConfig.HTTPAddrs {
			newConfig.HTTPAddrs[i] = redactSystemData(addr)
		}
	}

	return newConfig
}
