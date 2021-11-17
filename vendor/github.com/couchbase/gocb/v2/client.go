package gocb

import (
	"crypto/x509"
	"github.com/couchbase/gocbcore/v10"
	"github.com/pkg/errors"
	"sync"
)

type connectionManager interface {
	connect() error
	openBucket(bucketName string) error
	buildConfig(cluster *Cluster) error
	getKvProvider(bucketName string) (kvProvider, error)
	getKvCapabilitiesProvider(bucketName string) (kvCapabilityVerifier, error)
	getViewProvider(bucketName string) (viewProvider, error)
	getQueryProvider() (queryProvider, error)
	getAnalyticsProvider() (analyticsProvider, error)
	getSearchProvider() (searchProvider, error)
	getHTTPProvider(bucketName string) (httpProvider, error)
	getDiagnosticsProvider(bucketName string) (diagnosticsProvider, error)
	getWaitUntilReadyProvider(bucketName string) (waitUntilReadyProvider, error)
	connection(bucketName string) (*gocbcore.Agent, error)
	close() error
}

type stdConnectionMgr struct {
	lock       sync.Mutex
	agentgroup *gocbcore.AgentGroup
	config     *gocbcore.AgentGroupConfig
}

func newConnectionMgr() *stdConnectionMgr {
	client := &stdConnectionMgr{}
	return client
}

func (c *stdConnectionMgr) buildConfig(cluster *Cluster) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	breakerCfg := cluster.circuitBreakerConfig

	var completionCallback func(err error) bool
	if breakerCfg.CompletionCallback != nil {
		completionCallback = func(err error) bool {
			wrappedErr := maybeEnhanceKVErr(err, "", "", "", "")
			return breakerCfg.CompletionCallback(wrappedErr)
		}
	}

	var tlsRootCAProvider func() *x509.CertPool
	if cluster.internalConfig.TLSRootCAProvider == nil {
		tlsRootCAProvider = func() *x509.CertPool {
			if cluster.securityConfig.TLSSkipVerify {
				return nil
			}

			if cluster.securityConfig.TLSRootCAs == nil {
				return &x509.CertPool{}
			}

			return cluster.securityConfig.TLSRootCAs
		}
	} else {
		tlsRootCAProvider = cluster.internalConfig.TLSRootCAProvider
	}

	var authMechanisms []gocbcore.AuthMechanism
	for _, mech := range cluster.securityConfig.AllowedSaslMechanisms {
		authMechanisms = append(authMechanisms, gocbcore.AuthMechanism(mech))
	}

	config := &gocbcore.AgentGroupConfig{
		AgentConfig: gocbcore.AgentConfig{
			UserAgent: Identifier(),
			SecurityConfig: gocbcore.SecurityConfig{
				TLSRootCAProvider: tlsRootCAProvider,
				AuthMechanisms:    authMechanisms,
			},
			IoConfig: gocbcore.IoConfig{
				UseCollections:         true,
				UseDurations:           cluster.useServerDurations,
				UseMutationTokens:      cluster.useMutationTokens,
				UseOutOfOrderResponses: true,
			},
			KVConfig: gocbcore.KVConfig{
				ConnectTimeout: cluster.timeoutsConfig.ConnectTimeout,
			},
			DefaultRetryStrategy: cluster.retryStrategyWrapper,
			CircuitBreakerConfig: gocbcore.CircuitBreakerConfig{
				Enabled:                  !breakerCfg.Disabled,
				VolumeThreshold:          breakerCfg.VolumeThreshold,
				ErrorThresholdPercentage: breakerCfg.ErrorThresholdPercentage,
				SleepWindow:              breakerCfg.SleepWindow,
				RollingWindow:            breakerCfg.RollingWindow,
				CanaryTimeout:            breakerCfg.CanaryTimeout,
				CompletionCallback:       completionCallback,
			},
			OrphanReporterConfig: gocbcore.OrphanReporterConfig{
				Enabled:        cluster.orphanLoggerEnabled,
				ReportInterval: cluster.orphanLoggerInterval,
				SampleSize:     int(cluster.orphanLoggerSampleSize),
			},
			TracerConfig: gocbcore.TracerConfig{
				NoRootTraceSpans: true,
				Tracer:           &coreRequestTracerWrapper{tracer: cluster.tracer},
			},
			MeterConfig: gocbcore.MeterConfig{
				// At the moment we only support our own operations metric so there's no point in setting a meter for gocbcore.
				Meter: nil,
			},
		},
	}

	err := config.FromConnStr(cluster.connSpec().String())
	if err != nil {
		return err
	}

	config.SecurityConfig.Auth = &coreAuthWrapper{
		auth: cluster.authenticator(),
	}

	c.config = config
	return nil
}

func (c *stdConnectionMgr) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var err error
	c.agentgroup, err = gocbcore.CreateAgentGroup(c.config)
	if err != nil {
		return maybeEnhanceKVErr(err, "", "", "", "")
	}

	return nil
}

func (c *stdConnectionMgr) openBucket(bucketName string) error {
	if c.agentgroup == nil {
		return errors.New("cluster not yet connected")
	}

	return c.agentgroup.OpenBucket(bucketName)
}

func (c *stdConnectionMgr) getKvProvider(bucketName string) (kvProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}
	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("bucket not yet connected")
	}
	return agent, nil
}

func (c *stdConnectionMgr) getKvCapabilitiesProvider(bucketName string) (kvCapabilityVerifier, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}
	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("bucket not yet connected")
	}
	return agent.Internal(), nil
}

func (c *stdConnectionMgr) getViewProvider(bucketName string) (viewProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("bucket not yet connected")
	}
	return &viewProviderWrapper{provider: agent}, nil
}

func (c *stdConnectionMgr) getQueryProvider() (queryProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	return &queryProviderWrapper{provider: c.agentgroup}, nil
}

func (c *stdConnectionMgr) getAnalyticsProvider() (analyticsProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	return &analyticsProviderWrapper{provider: c.agentgroup}, nil
}

func (c *stdConnectionMgr) getSearchProvider() (searchProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	return &searchProviderWrapper{provider: c.agentgroup}, nil
}

func (c *stdConnectionMgr) getHTTPProvider(bucketName string) (httpProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	if bucketName == "" {
		return &httpProviderWrapper{provider: c.agentgroup}, nil
	}

	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("bucket not yet connected")
	}

	return &httpProviderWrapper{provider: agent}, nil
}

func (c *stdConnectionMgr) getDiagnosticsProvider(bucketName string) (diagnosticsProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	if bucketName == "" {
		return &diagnosticsProviderWrapper{provider: c.agentgroup}, nil
	}

	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("bucket not yet connected")
	}

	return &diagnosticsProviderWrapper{provider: agent}, nil
}

func (c *stdConnectionMgr) getWaitUntilReadyProvider(bucketName string) (waitUntilReadyProvider, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	if bucketName == "" {
		return &waitUntilReadyProviderWrapper{provider: c.agentgroup}, nil
	}

	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("provider not yet connected")
	}

	return &waitUntilReadyProviderWrapper{provider: agent}, nil
}

func (c *stdConnectionMgr) connection(bucketName string) (*gocbcore.Agent, error) {
	if c.agentgroup == nil {
		return nil, errors.New("cluster not yet connected")
	}

	agent := c.agentgroup.GetAgent(bucketName)
	if agent == nil {
		return nil, errors.New("bucket not yet connected")
	}
	return agent, nil
}

func (c *stdConnectionMgr) close() error {
	c.lock.Lock()
	if c.agentgroup == nil {
		c.lock.Unlock()
		return errors.New("cluster not yet connected")
	}
	defer c.lock.Unlock()
	return c.agentgroup.Close()
}
