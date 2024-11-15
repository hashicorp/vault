package gocbcore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/couchbase/gocbcore/v10/connstr"
)

func parseDurationOrInt(valStr string) (time.Duration, error) {
	dur, err := time.ParseDuration(valStr)
	if err != nil {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return 0, err
		}

		dur = time.Duration(val) * time.Millisecond
	}

	return dur, nil
}

// AgentConfig specifies the configuration options for creation of an Agent.
type AgentConfig struct {
	BucketName string
	UserAgent  string

	SeedConfig SeedConfig

	SecurityConfig SecurityConfig

	CompressionConfig CompressionConfig

	ConfigPollerConfig ConfigPollerConfig

	IoConfig IoConfig

	KVConfig KVConfig

	HTTPConfig HTTPConfig

	DefaultRetryStrategy RetryStrategy

	CircuitBreakerConfig CircuitBreakerConfig

	OrphanReporterConfig OrphanReporterConfig

	TracerConfig TracerConfig

	MeterConfig MeterConfig

	InternalConfig InternalConfig
}

// OrphanReporterConfig specifies options for controlling the orphan
// reporter which records when the SDK receives responses for requests
// that are no longer in the system (usually due to being timed out).
type OrphanReporterConfig struct {
	Enabled bool
	// ReportInterval is the time period used for how often a report is logged.
	ReportInterval time.Duration
	// SampleSize is the number of requests which will be reported.
	SampleSize int
}

func (config OrphanReporterConfig) fromSpec(spec connstr.ResolvedConnSpec) (OrphanReporterConfig, error) {
	if valStr, ok := fetchOption(spec, "orphaned_response_logging"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return OrphanReporterConfig{}, fmt.Errorf("orphaned_response_logging option must be a boolean")
		}
		config.Enabled = val
	}

	if valStr, ok := fetchOption(spec, "orphaned_response_logging_interval"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return OrphanReporterConfig{}, fmt.Errorf("orphaned_response_logging_interval option must be a number")
		}
		config.ReportInterval = val
	}

	if valStr, ok := fetchOption(spec, "orphaned_response_logging_sample_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return OrphanReporterConfig{}, fmt.Errorf("orphaned_response_logging_sample_size option must be a number")
		}
		config.SampleSize = int(val)
	}

	return config, nil
}

// SecurityConfig specifies options for controlling security related
// items such as TLS root certificates and verification skipping.
type SecurityConfig struct {
	UseTLS            bool
	TLSRootCAProvider func() *x509.CertPool

	// NoTLSSeedNode indicates that, even with UseTLS set to true, the SDK should always connect to the seed node
	// over a non TLS connection. This means that the seed node should ALWAYS be localhost.
	// This option must be used with the ConfigPollerConfig UseSeedPoller set to true.
	// Internal: This should never be used and is not supported.
	NoTLSSeedNode bool

	Auth AuthProvider

	// AuthMechanisms is the list of mechanisms that the SDK can use to attempt authentication.
	// Note that if you add PLAIN to the list, this will cause credential leakage on the network
	// since PLAIN sends the credentials in cleartext. It is disabled by default to prevent downgrade attacks. We
	// recommend using a TLS connection if using PLAIN.
	AuthMechanisms []AuthMechanism
}

func (config SecurityConfig) fromSpec(spec connstr.ResolvedConnSpec) (SecurityConfig, error) {
	if spec.UseSsl {
		cacertpaths := spec.Options["ca_cert_path"]

		if len(cacertpaths) > 0 {
			roots := x509.NewCertPool()

			for _, path := range cacertpaths {
				cacert, err := ioutil.ReadFile(path)
				if err != nil {
					return SecurityConfig{}, err
				}

				ok := roots.AppendCertsFromPEM(cacert)
				if !ok {
					return SecurityConfig{}, errInvalidCertificate
				}
			}

			config.TLSRootCAProvider = func() *x509.CertPool {
				return roots
			}
		}

		config.UseTLS = true
	}

	if spec.NSServerHost != nil {
		config.NoTLSSeedNode = true
	}

	return config, nil
}

// CompressionConfig specifies options for controlling compression applied to documents using KV.
type CompressionConfig struct {
	Enabled              bool
	DisableDecompression bool
	MinSize              int
	MinRatio             float64
}

func (config CompressionConfig) fromSpec(spec connstr.ResolvedConnSpec) (CompressionConfig, error) {
	if valStr, ok := fetchOption(spec, "compression"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return CompressionConfig{}, fmt.Errorf("compression option must be a boolean")
		}
		config.Enabled = val
	}

	if valStr, ok := fetchOption(spec, "compression_min_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return CompressionConfig{}, fmt.Errorf("compression_min_size option must be an int")
		}
		config.MinSize = int(val)
	}

	if valStr, ok := fetchOption(spec, "compression_min_ratio"); ok {
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return CompressionConfig{}, fmt.Errorf("compression_min_size option must be an int")
		}
		config.MinRatio = val
	}

	return config, nil
}

// ConfigPollerConfig specifies options for controlling the cluster configuration pollers.
type ConfigPollerConfig struct {
	HTTPRedialPeriod time.Duration
	HTTPRetryDelay   time.Duration
	HTTPMaxWait      time.Duration
	CccpMaxWait      time.Duration
	CccpPollPeriod   time.Duration
}

func (config ConfigPollerConfig) fromSpec(spec connstr.ResolvedConnSpec) (ConfigPollerConfig, error) {
	if valStr, ok := fetchOption(spec, "config_poll_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return ConfigPollerConfig{}, fmt.Errorf("config poll timeout option must be a duration or a number")
		}
		config.CccpMaxWait = val
	}

	if valStr, ok := fetchOption(spec, "config_poll_interval"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return ConfigPollerConfig{}, fmt.Errorf("config pool interval option must be duration or a number")
		}
		config.CccpPollPeriod = val
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "http_redial_period"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return ConfigPollerConfig{}, fmt.Errorf("http redial period option must be a duration or a number")
		}
		config.HTTPRedialPeriod = val
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "http_retry_delay"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return ConfigPollerConfig{}, fmt.Errorf("http retry delay option must be a duration or a number")
		}
		config.HTTPRetryDelay = val
	}

	if valStr, ok := fetchOption(spec, "http_config_poll_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return ConfigPollerConfig{}, fmt.Errorf("http_config_poll_timeout option must be a duration or a number")
		}
		config.HTTPMaxWait = val
	}

	return config, nil
}

// IoConfig specifies IO related configuration options such as HELLO flags and the network type to use.
type IoConfig struct {
	// NetworkType defines which network to use from the cluster config.
	NetworkType string

	UseMutationTokens           bool
	UseDurations                bool
	UseOutOfOrderResponses      bool
	DisableXErrorHello          bool
	DisableJSONHello            bool
	DisableSyncReplicationHello bool
	EnablePITRHello             bool
	UseCollections              bool

	UseClusterMapNotifications bool
}

func (config IoConfig) fromSpec(spec connstr.ResolvedConnSpec) (IoConfig, error) {
	if valStr, ok := fetchOption(spec, "network"); ok {
		config.NetworkType = valStr
	}

	if valStr, ok := fetchOption(spec, "enable_mutation_tokens"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return IoConfig{}, fmt.Errorf("enable_mutation_tokens option must be a boolean")
		}
		config.UseMutationTokens = val
	}

	if valStr, ok := fetchOption(spec, "enable_server_durations"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return IoConfig{}, fmt.Errorf("server_duration option must be a boolean")
		}
		config.UseDurations = val
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "unordered_execution_enabled"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return IoConfig{}, fmt.Errorf("unordered_execution_enabled option must be a boolean")
		}
		config.UseOutOfOrderResponses = val
	}

	if valStr, ok := fetchOption(spec, "enable_cluster_config_notifications"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return IoConfig{}, fmt.Errorf("enable_cluster_config_notifications option must be a boolean")
		}
		config.UseClusterMapNotifications = val
	}

	return config, nil
}

// TracerConfig specifies tracer related configuration options.
type TracerConfig struct {
	Tracer           RequestTracer
	NoRootTraceSpans bool
}

// MeterConfig specifies meter related configuration options.
type MeterConfig struct {
	Meter Meter
}

// HTTPConfig specifies http related configuration options.
type HTTPConfig struct {
	// MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts.
	MaxIdleConns int
	// MaxIdleConnsPerHost controls the maximum idle (keep-alive) connections to keep per-host.
	MaxIdleConnsPerHost int
	ConnectTimeout      time.Duration
	// IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing
	// itself.
	IdleConnectionTimeout time.Duration
}

func (config HTTPConfig) fromSpec(spec connstr.ResolvedConnSpec) (HTTPConfig, error) {
	if valStr, ok := fetchOption(spec, "max_idle_http_connections"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return HTTPConfig{}, fmt.Errorf("http max idle connections option must be a number")
		}
		config.MaxIdleConns = int(val)
	}

	if valStr, ok := fetchOption(spec, "max_perhost_idle_http_connections"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return HTTPConfig{}, fmt.Errorf("max_perhost_idle_http_connections option must be a number")
		}
		config.MaxIdleConnsPerHost = int(val)
	}

	if valStr, ok := fetchOption(spec, "idle_http_connection_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return HTTPConfig{}, fmt.Errorf("idle_http_connection_timeout option must be a duration or a number")
		}
		config.IdleConnectionTimeout = val
	}

	if valStr, ok := fetchOption(spec, "http_connect_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return HTTPConfig{}, fmt.Errorf("http_connect_timeout option must be a duration or a number")
		}
		config.ConnectTimeout = val
	}

	return config, nil
}

// KVConfig specifies kv related configuration options.
type KVConfig struct {
	// ConnectTimeout is the timeout value to apply when dialling tcp connections.
	ConnectTimeout time.Duration
	// ServerWaitBackoff is the period of time that the SDK will wait before reattempting connection to a node after
	// bootstrap fails against that node.
	ServerWaitBackoff time.Duration

	// The number of connections to create to each node.
	PoolSize int
	// The maximum number of requests that can be queued waiting to be sent to a node.
	MaxQueueSize int

	// Note: if you create multiple agents with different buffer sizes within the same environment then you will
	// get indeterminate behaviour, the connections may not even use the provided buffer size.
	ConnectionBufferSize uint
}

func (config KVConfig) fromSpec(spec connstr.ResolvedConnSpec) (KVConfig, error) {

	if valStr, ok := fetchOption(spec, "kv_connect_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return KVConfig{}, fmt.Errorf("kv_connect_timeout option must be a duration or a number")
		}
		config.ConnectTimeout = val
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "kv_pool_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return KVConfig{}, fmt.Errorf("kv pool size option must be a number")
		}
		config.PoolSize = int(val)
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "max_queue_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return KVConfig{}, fmt.Errorf("max queue size option must be a number")
		}
		config.MaxQueueSize = int(val)
	}

	// This option is experimental
	if valStr, ok := fetchOption(spec, "kv_buffer_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return KVConfig{}, fmt.Errorf("kv buffer size option must be a number")
		}
		config.ConnectionBufferSize = uint(val)
	}

	if valStr, ok := fetchOption(spec, "server_wait_backoff"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return KVConfig{}, fmt.Errorf("server wait backoff must be a number")
		}
		config.ServerWaitBackoff = time.Duration(val) * time.Millisecond
	}

	return config, nil
}

// SRVRecord describes the SRV record used to extract memd addresses in the SeedConfig.
type SRVRecord struct {
	Proto  string
	Scheme string
	Host   string
}

func (srv SRVRecord) redacted() interface{} {
	return fmt.Sprintf("%s %s %s", srv.Proto, srv.Scheme, redactSystemData(srv.Host))
}

// SeedConfig specifies initial seed configuration options such as addresses.
type SeedConfig struct {
	HTTPAddrs []string
	MemdAddrs []string
	SRVRecord *SRVRecord
}

func (config SeedConfig) fromSpec(spec connstr.ResolvedConnSpec) (SeedConfig, error) {
	// Grab the resolved hostnames into a set of string arrays
	var httpHosts []string
	for _, specHost := range spec.HttpHosts {
		httpHosts = append(httpHosts, fmt.Sprintf("%s:%d", specHost.Host, specHost.Port))
	}

	var memdHosts []string
	for _, specHost := range spec.MemdHosts {
		memdHosts = append(memdHosts, fmt.Sprintf("%s:%d", specHost.Host, specHost.Port))
	}

	var nsServerHost string
	if spec.NSServerHost != nil {
		nsServerHost = fmt.Sprintf("%s:%d", spec.NSServerHost.Host, spec.NSServerHost.Port)
	}

	if nsServerHost != "" {
		if len(httpHosts) > 0 || len(memdHosts) > 0 {
			return SeedConfig{}, errors.New("ns_server host cannot be used alongside http or memd hosts")
		}

		httpHosts = append(httpHosts, nsServerHost)
	}

	// Get bootstrap_on option to determine which, if any, of the bootstrap nodes should be cleared
	switch val, _ := fetchOption(spec, "bootstrap_on"); val {
	case "http":
		memdHosts = nil
		if len(httpHosts) == 0 {
			return SeedConfig{}, errors.New("bootstrap_on=http but no HTTP hosts in connection string")
		}
	case "cccp":
		httpHosts = nil
		if len(memdHosts) == 0 {
			return SeedConfig{}, errors.New("bootstrap_on=cccp but no CCCP/Memcached hosts in connection string")
		}
	case "both":
		if nsServerHost != "" {
			return SeedConfig{}, errors.New("bootstrap_on=both but ns_server host in connection string")
		}
	case "":
		// Do nothing
		break
	default:
		// Don't advertise ns_server as an option
		return SeedConfig{}, errors.New("bootstrap_on={http,cccp,both}")
	}
	config.MemdAddrs = memdHosts
	config.HTTPAddrs = httpHosts
	if spec.SrvRecord != nil {
		config.SRVRecord = &SRVRecord{
			Proto:  spec.SrvRecord.Proto,
			Scheme: spec.SrvRecord.Scheme,
			Host:   spec.SrvRecord.Host,
		}
	}

	return config, nil
}

func (config SeedConfig) redacted() SeedConfig {
	newConfig := SeedConfig{
		HTTPAddrs: config.HTTPAddrs,
		MemdAddrs: config.MemdAddrs,
		SRVRecord: config.SRVRecord,
	}
	// The slices here are still pointing at config's underlying arrays
	// so we need to make them not do that.
	newConfig.HTTPAddrs = append([]string(nil), newConfig.HTTPAddrs...)
	for i, addr := range newConfig.HTTPAddrs {
		newConfig.HTTPAddrs[i] = redactSystemData(addr)
	}
	newConfig.MemdAddrs = append([]string(nil), newConfig.MemdAddrs...)
	for i, addr := range newConfig.MemdAddrs {
		newConfig.MemdAddrs[i] = redactSystemData(addr)
	}

	return newConfig
}

// InternalConfig specifies internal configs.
// Internal: This should never be used and is not supported.
type InternalConfig struct {
	EnableResourceUnitsTrackingHello bool
}

func (config InternalConfig) fromSpec(spec connstr.ResolvedConnSpec) (InternalConfig, error) {
	if valStr, ok := fetchOption(spec, "enable_resource_units"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return InternalConfig{}, fmt.Errorf("enable_resource_units option must be a boolean")
		}
		config.EnableResourceUnitsTrackingHello = val
	}

	return config, nil
}

func (config *AgentConfig) redacted() interface{} {
	newConfig := *config
	if isLogRedactionLevelFull() {
		newConfig.SeedConfig = newConfig.SeedConfig.redacted()

		if newConfig.BucketName != "" {
			newConfig.BucketName = redactMetaData(newConfig.BucketName)
		}
	}

	return newConfig
}

func fetchOption(spec connstr.ResolvedConnSpec, name string) (string, bool) {
	optValue := spec.Options[name]
	if len(optValue) == 0 {
		return "", false
	}
	return optValue[len(optValue)-1], true
}

// FromConnStr populates the AgentConfig with information from a
// Couchbase Connection String.
// Supported options are:
//
//		bootstrap_on (bool) - Specifies what protocol to bootstrap on (cccp, http).
//		ca_cert_path (string) - Specifies the path to a CA certificate.
//		network (string) - The network type to use.
//		kv_connect_timeout (duration) - Maximum period to attempt to connect to cluster in ms.
//		config_poll_interval (duration) - Period to wait between CCCP config polling in ms.
//		config_poll_timeout (duration) - Maximum period of time to wait for a CCCP request.
//		compression (bool) - Whether to enable network-wise compression of documents.
//		compression_min_size (int) - The minimal size of the document in bytes to consider compression.
//		compression_min_ratio (float64) - The minimal compress ratio (compressed / original) for the document to be sent compressed.
//		enable_server_durations (bool) - Whether to enable fetching server operation durations.
//		max_idle_http_connections (int) - Maximum number of idle http connections in the pool.
//		max_perhost_idle_http_connections (int) - Maximum number of idle http connections in the pool per host.
//		idle_http_connection_timeout (duration) - Maximum length of time for an idle connection to stay in the pool in ms.
//		orphaned_response_logging (bool) - Whether to enable orphaned response logging.
//		orphaned_response_logging_interval (duration) - How often to print the orphan log records.
//		orphaned_response_logging_sample_size (int) - The maximum number of orphan log records to track.
//		dcp_priority (int) - Specifies the priority to request from the Cluster when connecting for DCP.
//		enable_dcp_expiry (bool) - Whether to enable the feature to distinguish between explicit delete and expired delete on DCP.
//		http_redial_period (duration) - The maximum length of time for the HTTP poller to stay connected before reconnecting.
//		http_retry_delay (duration) - The length of time to wait between HTTP poller retries if connecting fails.
//		kv_pool_size (int) - The number of connections to create to each kv node.
//		max_queue_size (int) - The maximum number of requests that can be queued for sending per connection.
//		unordered_execution_enabled (bool) - Whether to enabled the "out of order responses" feature.
//	 server_wait_backoff (duration) -The period of time waited between kv reconnect attmepts to a node after connection failure
func (config *AgentConfig) FromConnStr(connStr string) error {
	baseSpec, err := connstr.Parse(connStr)
	if err != nil {
		return err
	}

	spec, err := connstr.Resolve(baseSpec)
	if err != nil {
		return err
	}

	if spec.Bucket != "" {
		config.BucketName = spec.Bucket
	}

	config.SeedConfig, err = config.SeedConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.SecurityConfig, err = config.SecurityConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	config.OrphanReporterConfig, err = config.OrphanReporterConfig.fromSpec(spec)
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

	config.InternalConfig, err = config.InternalConfig.fromSpec(spec)
	if err != nil {
		return err
	}

	return nil
}
