package gocb

import (
	"crypto/x509"
	"sync"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
	"github.com/pkg/errors"
)

type client interface {
	connect() error
	buildConfig(cluster *Cluster, bucket string) error
	getKvProvider() (kvProvider, error)
	getViewProvider() (viewProvider, error)
	getQueryProvider() (queryProvider, error)
	getAnalyticsProvider() (analyticsProvider, error)
	getSearchProvider() (searchProvider, error)
	getHTTPProvider() (httpProvider, error)
	getDiagnosticsProvider() (diagnosticsProvider, error)
	getWaitUntilReadyProvider() (waitUntilReadyProvider, error)
	close() error
	setBootstrapError(err error)
	supportsGCCCP() bool
	supportsCollections() bool
	connected() (bool, error)
	getBootstrapError() error
}

type stdClient struct {
	lock         sync.Mutex
	agent        *gocbcore.Agent
	bootstrapErr error
	config       *gocbcore.AgentConfig
}

func newClient() *stdClient {
	client := &stdClient{}
	return client
}

func (c *stdClient) buildConfig(cluster *Cluster, bucket string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	breakerCfg := cluster.circuitBreakerConfig

	var completionCallback func(err error) bool
	if breakerCfg.CompletionCallback != nil {
		completionCallback = func(err error) bool {
			wrappedErr := maybeEnhanceKVErr(err, bucket, "", "", "")
			return breakerCfg.CompletionCallback(wrappedErr)
		}
	}

	var tlsRootCAProvider func() *x509.CertPool
	if cluster.internalConfig.TLSRootCAProvider == nil {
		tlsRootCAProvider = func() *x509.CertPool {
			if cluster.securityConfig.TLSSkipVerify {
				return nil
			}

			return cluster.securityConfig.TLSRootCAs
		}
	} else {
		tlsRootCAProvider = cluster.internalConfig.TLSRootCAProvider
	}

	config := &gocbcore.AgentConfig{
		UserAgent:              Identifier(),
		TLSRootCAProvider:      tlsRootCAProvider,
		ConnectTimeout:         cluster.timeoutsConfig.ConnectTimeout,
		UseMutationTokens:      cluster.useMutationTokens,
		KVConnectTimeout:       7000 * time.Millisecond,
		UseDurations:           cluster.useServerDurations,
		UseCollections:         true,
		BucketName:             bucket,
		UseZombieLogger:        cluster.orphanLoggerEnabled,
		ZombieLoggerInterval:   cluster.orphanLoggerInterval,
		ZombieLoggerSampleSize: int(cluster.orphanLoggerSampleSize),
		NoRootTraceSpans:       true,
		Tracer:                 &requestTracerWrapper{cluster.tracer},
		CircuitBreakerConfig: gocbcore.CircuitBreakerConfig{
			Enabled:                  !breakerCfg.Disabled,
			VolumeThreshold:          breakerCfg.VolumeThreshold,
			ErrorThresholdPercentage: breakerCfg.ErrorThresholdPercentage,
			SleepWindow:              breakerCfg.SleepWindow,
			RollingWindow:            breakerCfg.RollingWindow,
			CanaryTimeout:            breakerCfg.CanaryTimeout,
			CompletionCallback:       completionCallback,
		},
	}

	err := config.FromConnStr(cluster.connSpec().String())
	if err != nil {
		return err
	}

	config.Auth = &coreAuthWrapper{
		auth: cluster.authenticator(),
	}

	c.config = config
	return nil
}

func (c *stdClient) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	agent, err := gocbcore.CreateAgent(c.config)
	if err != nil {
		return maybeEnhanceKVErr(err, c.config.BucketName, "", "", "")
	}

	c.agent = agent
	return nil
}

func (c *stdClient) setBootstrapError(err error) {
	c.bootstrapErr = err
}

func (c *stdClient) getBootstrapError() error {
	return c.bootstrapErr
}

func (c *stdClient) getKvProvider() (kvProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return c.agent, nil
}

func (c *stdClient) getViewProvider() (viewProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &viewProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) getQueryProvider() (queryProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &queryProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) getAnalyticsProvider() (analyticsProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &analyticsProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) getSearchProvider() (searchProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &searchProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) getHTTPProvider() (httpProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &httpProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) getDiagnosticsProvider() (diagnosticsProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &diagnosticsProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) getWaitUntilReadyProvider() (waitUntilReadyProvider, error) {
	if c.bootstrapErr != nil {
		return nil, c.bootstrapErr
	}

	if c.agent == nil {
		return nil, errors.New("cluster not yet connected")
	}
	return &waitUntilReadyProviderWrapper{provider: c.agent}, nil
}

func (c *stdClient) connected() (bool, error) {
	return c.agent.HasSeenConfig()
}

func (c *stdClient) supportsGCCCP() bool {
	return c.agent.UsingGCCCP()
}

func (c *stdClient) supportsCollections() bool {
	return c.agent.HasCollectionsSupport()
}

func (c *stdClient) close() error {
	c.lock.Lock()
	if c.agent == nil {
		c.lock.Unlock()
		return errors.New("cluster not yet connected")
	}
	c.lock.Unlock()
	return c.agent.Close()
}
