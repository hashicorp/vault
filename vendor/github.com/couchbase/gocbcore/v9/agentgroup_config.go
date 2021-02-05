package gocbcore

// AgentGroupConfig specifies the configuration options for creation of an AgentGroup.
type AgentGroupConfig struct {
	AgentConfig
}

func (config *AgentGroupConfig) redacted() interface{} {
	return config.AgentConfig.redacted()
}

// FromConnStr populates the AgentGroupConfig with information from a
// Couchbase Connection String. See AgentConfig for supported options.
func (config *AgentGroupConfig) FromConnStr(connStr string) error {
	return config.AgentConfig.FromConnStr(connStr)
}

func (config *AgentGroupConfig) toAgentConfig() *AgentConfig {
	return &AgentConfig{
		MemdAddrs:                 config.MemdAddrs,
		HTTPAddrs:                 config.HTTPAddrs,
		BucketName:                config.BucketName,
		UserAgent:                 config.UserAgent,
		UseTLS:                    config.UseTLS,
		NetworkType:               config.NetworkType,
		Auth:                      config.Auth,
		TLSRootCAProvider:         config.TLSRootCAProvider,
		UseMutationTokens:         config.UseMutationTokens,
		UseCompression:            config.UseCompression,
		UseDurations:              config.UseDurations,
		DisableDecompression:      config.DisableDecompression,
		UseOutOfOrderResponses:    config.UseOutOfOrderResponses,
		UseCollections:            config.UseCollections,
		CompressionMinSize:        config.CompressionMinSize,
		CompressionMinRatio:       config.CompressionMinRatio,
		HTTPRedialPeriod:          config.HTTPRedialPeriod,
		HTTPRetryDelay:            config.HTTPRetryDelay,
		CccpMaxWait:               config.CccpMaxWait,
		CccpPollPeriod:            config.CccpPollPeriod,
		ConnectTimeout:            config.ConnectTimeout,
		KVConnectTimeout:          config.KVConnectTimeout,
		KvPoolSize:                config.KvPoolSize,
		MaxQueueSize:              config.MaxQueueSize,
		HTTPMaxIdleConns:          config.HTTPMaxIdleConns,
		HTTPMaxIdleConnsPerHost:   config.HTTPMaxIdleConnsPerHost,
		HTTPIdleConnectionTimeout: config.HTTPIdleConnectionTimeout,
		Tracer:                    config.Tracer,
		NoRootTraceSpans:          config.NoRootTraceSpans,
		DefaultRetryStrategy:      config.DefaultRetryStrategy,
		CircuitBreakerConfig:      config.CircuitBreakerConfig,
		UseZombieLogger:           config.UseZombieLogger,
		ZombieLoggerInterval:      config.ZombieLoggerInterval,
		ZombieLoggerSampleSize:    config.ZombieLoggerSampleSize,
		AuthMechanisms:            config.AuthMechanisms,
	}
}
