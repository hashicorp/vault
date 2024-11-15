package gocbcore

import (
	"fmt"
	"strconv"

	"github.com/couchbase/gocbcore/v10/connstr"
)

// DCPAgentConfig specifies the configuration options for creation of a DCPAgent.
type DCPAgentConfig struct {
	UserAgent  string
	BucketName string

	SeedConfig SeedConfig

	SecurityConfig SecurityConfig

	CompressionConfig CompressionConfig

	ConfigPollerConfig ConfigPollerConfig

	// EnableCCCPPoller will enable the use of the CCCP poller for the SDK.
	// By default, only HTTP polling is used and CCCP polling during a DCP stream is discouraged.
	EnableCCCPPoller bool

	IoConfig IoConfig

	KVConfig KVConfig

	HTTPConfig HTTPConfig

	DCPConfig DCPConfig
}

// DCPConfig specifies DCP specific configuration options.
type DCPConfig struct {
	AgentPriority    DcpAgentPriority
	UseChangeStreams bool
	UseExpiryOpcode  bool
	UseStreamID      bool
	UseOSOBackfill   bool
	BackfillOrder    DCPBackfillOrder

	BufferSize                   int
	DisableBufferAcknowledgement bool
}

func (config DCPConfig) fromSpec(spec connstr.ResolvedConnSpec) (DCPConfig, error) {
	// This option is experimental
	if valStr, ok := fetchOption(spec, "dcp_priority"); ok {
		var priority DcpAgentPriority
		switch valStr {
		case "":
			priority = DcpAgentPriorityLow
		case "low":
			priority = DcpAgentPriorityLow
		case "medium":
			priority = DcpAgentPriorityMed
		case "high":
			priority = DcpAgentPriorityHigh
		default:
			return DCPConfig{}, fmt.Errorf("dcp_priority must be one of low, medium or high")
		}
		config.AgentPriority = priority
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "dcp_buffer_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return DCPConfig{}, fmt.Errorf("dcp buffer size option must be a number")
		}
		config.BufferSize = int(val)
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "enable_dcp_change_streams"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return DCPConfig{}, fmt.Errorf("enable_dcp_change_streams option must be a boolean")
		}
		config.UseChangeStreams = val
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "enable_dcp_expiry"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return DCPConfig{}, fmt.Errorf("enable_dcp_expiry option must be a boolean")
		}
		config.UseExpiryOpcode = val
	}

	return config, nil
}

func (config *DCPAgentConfig) redacted() interface{} {
	newConfig := *config

	if isLogRedactionLevelFull() {
		newConfig.SeedConfig = newConfig.SeedConfig.redacted()

		if newConfig.BucketName != "" {
			newConfig.BucketName = redactMetaData(newConfig.BucketName)
		}
	}

	return newConfig
}

// FromConnStr populates the AgentConfig with information from a
// Couchbase Connection String.
// Supported options are:
//
//	ca_cert_path (string) - Specifies the path to a CA certificate.
//	network (string) - The network type to use.
//	kv_connect_timeout (duration) - Maximum period to attempt to connect to cluster in ms.
//	config_poll_interval (duration) - Period to wait between CCCP config polling in ms.
//	config_poll_timeout (duration) - Maximum period of time to wait for a CCCP request.
//	compression (bool) - Whether to enable network-wise compression of documents.
//	compression_min_size (int) - The minimal size of the document in bytes to consider compression.
//	compression_min_ratio (float64) - The minimal compress ratio (compressed / original) for the document to be sent compressed.
//	orphaned_response_logging (bool) - Whether to enable orphaned response logging.
//	orphaned_response_logging_interval (duration) - How often to print the orphan log records.
//	orphaned_response_logging_sample_size (int) - The maximum number of orphan log records to track.
//	dcp_priority (int) - Specifies the priority to request from the Cluster when connecting for DCP.
//	enable_dcp_change_streams (bool) - Enables the DCP connection to allow history snapshots in DCP streams.
//	enable_dcp_expiry (bool) - Whether to enable the feature to distinguish between explicit delete and expired delete on DCP.
//	kv_pool_size (int) - The number of connections to create to each kv node.
//	max_queue_size (int) - The maximum number of requests that can be queued for sending per connection.
//	max_idle_http_connections (int) - Maximum number of idle http connections in the pool.
//	max_perhost_idle_http_connections (int) - Maximum number of idle http connections in the pool per host.
//	idle_http_connection_timeout (duration) - Maximum length of time for an idle connection to stay in the pool in ms.
//	http_redial_period (duration) - The maximum length of time for the HTTP poller to stay connected before reconnecting.
//	http_retry_delay (duration) - The length of time to wait between HTTP poller retries if connecting fails.
func (config *DCPAgentConfig) FromConnStr(connStr string) error {
	baseSpec, err := connstr.Parse(connStr)
	if err != nil {
		return err
	}

	spec, err := connstr.Resolve(baseSpec)
	if err != nil {
		return err
	}

	config.DCPConfig, err = config.DCPConfig.fromSpec(spec)
	if err != nil {
		return err
	}
	config.SeedConfig, err = config.SeedConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.SecurityConfig, err = config.SecurityConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.CompressionConfig, err = config.CompressionConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.ConfigPollerConfig, err = config.ConfigPollerConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.IoConfig, err = config.IoConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.HTTPConfig, err = config.HTTPConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.KVConfig, err = config.KVConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	return nil
}
