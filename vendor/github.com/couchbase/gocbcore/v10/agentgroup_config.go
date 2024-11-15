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
		BucketName:           config.BucketName,
		UserAgent:            config.UserAgent,
		SeedConfig:           config.SeedConfig,
		SecurityConfig:       config.SecurityConfig,
		CompressionConfig:    config.CompressionConfig,
		ConfigPollerConfig:   config.ConfigPollerConfig,
		IoConfig:             config.IoConfig,
		KVConfig:             config.KVConfig,
		HTTPConfig:           config.HTTPConfig,
		DefaultRetryStrategy: config.DefaultRetryStrategy,
		CircuitBreakerConfig: config.CircuitBreakerConfig,
		OrphanReporterConfig: config.OrphanReporterConfig,
		MeterConfig:          config.MeterConfig,
		TracerConfig:         config.TracerConfig,
		InternalConfig:       config.InternalConfig,
	}
}
