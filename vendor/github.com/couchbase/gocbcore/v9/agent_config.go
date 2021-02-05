package gocbcore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/couchbase/gocbcore/v9/connstr"
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
	MemdAddrs   []string
	HTTPAddrs   []string
	BucketName  string
	UserAgent   string
	UseTLS      bool
	NetworkType string
	Auth        AuthProvider

	TLSRootCAProvider func() *x509.CertPool

	UseMutationTokens      bool
	UseCompression         bool
	UseDurations           bool
	DisableDecompression   bool
	UseOutOfOrderResponses bool
	DisableXErrors         bool

	DisableJSONHello            bool
	DisableSyncReplicationHello bool

	UseCollections bool

	CompressionMinSize  int
	CompressionMinRatio float64

	HTTPRedialPeriod time.Duration
	HTTPRetryDelay   time.Duration
	CccpMaxWait      time.Duration
	CccpPollPeriod   time.Duration

	ConnectTimeout   time.Duration
	KVConnectTimeout time.Duration

	KvPoolSize   int
	MaxQueueSize int

	HTTPMaxIdleConns          int
	HTTPMaxIdleConnsPerHost   int
	HTTPIdleConnectionTimeout time.Duration

	// Volatile: Tracer API is subject to change.
	Tracer           RequestTracer
	NoRootTraceSpans bool

	DefaultRetryStrategy RetryStrategy
	CircuitBreakerConfig CircuitBreakerConfig

	UseZombieLogger        bool
	ZombieLoggerInterval   time.Duration
	ZombieLoggerSampleSize int

	// AuthMechanisms is the list of mechanisms that the SDK can use to attempt authentication.
	AuthMechanisms []AuthMechanism
}

func (config *AgentConfig) redacted() interface{} {
	newConfig := AgentConfig{}
	newConfig = *config
	if isLogRedactionLevelFull() {
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

		if newConfig.BucketName != "" {
			newConfig.BucketName = redactMetaData(newConfig.BucketName)
		}
	}

	return newConfig
}

// FromConnStr populates the AgentConfig with information from a
// Couchbase Connection String.
// Supported options are:
//   bootstrap_on (bool) - Specifies what protocol to bootstrap on (cccp, http).
//   ca_cert_path (string) - Specifies the path to a CA certificate.
//   network (string) - The network type to use.
//   kv_connect_timeout (duration) - Maximum period to attempt to connect to cluster in ms.
//   config_poll_interval (duration) - Period to wait between CCCP config polling in ms.
//   config_poll_timeout (duration) - Maximum period of time to wait for a CCCP request.
//   compression (bool) - Whether to enable network-wise compression of documents.
//   compression_min_size (int) - The minimal size of the document in bytes to consider compression.
//   compression_min_ratio (float64) - The minimal compress ratio (compressed / original) for the document to be sent compressed.
//   enable_server_durations (bool) - Whether to enable fetching server operation durations.
//   max_idle_http_connections (int) - Maximum number of idle http connections in the pool.
//   max_perhost_idle_http_connections (int) - Maximum number of idle http connections in the pool per host.
//   idle_http_connection_timeout (duration) - Maximum length of time for an idle connection to stay in the pool in ms.
//   orphaned_response_logging (bool) - Whether to enable orphaned response logging.
//   orphaned_response_logging_interval (duration) - How often to print the orphan log records.
//   orphaned_response_logging_sample_size (int) - The maximum number of orphan log records to track.
//   dcp_priority (int) - Specifies the priority to request from the Cluster when connecting for DCP.
//   enable_dcp_expiry (bool) - Whether to enable the feature to distinguish between explicit delete and expired delete on DCP.
//   http_redial_period (duration) - The maximum length of time for the HTTP poller to stay connected before reconnecting.
//   http_retry_delay (duration) - The length of time to wait between HTTP poller retries if connecting fails.
//   kv_pool_size (int) - The number of connections to create to each kv node.
//   max_queue_size (int) - The maximum number of requests that can be queued for sending per connection.
func (config *AgentConfig) FromConnStr(connStr string) error {
	baseSpec, err := connstr.Parse(connStr)
	if err != nil {
		return err
	}

	spec, err := connstr.Resolve(baseSpec)
	if err != nil {
		return err
	}

	fetchOption := func(name string) (string, bool) {
		optValue := spec.Options[name]
		if len(optValue) == 0 {
			return "", false
		}
		return optValue[len(optValue)-1], true
	}

	// Grab the resolved hostnames into a set of string arrays
	var httpHosts []string
	for _, specHost := range spec.HttpHosts {
		httpHosts = append(httpHosts, fmt.Sprintf("%s:%d", specHost.Host, specHost.Port))
	}

	var memdHosts []string
	for _, specHost := range spec.MemdHosts {
		memdHosts = append(memdHosts, fmt.Sprintf("%s:%d", specHost.Host, specHost.Port))
	}

	// Get bootstrap_on option to determine which, if any, of the bootstrap nodes should be cleared
	switch val, _ := fetchOption("bootstrap_on"); val {
	case "http":
		memdHosts = nil
		if len(httpHosts) == 0 {
			return errors.New("bootstrap_on=http but no HTTP hosts in connection string")
		}
	case "cccp":
		httpHosts = nil
		if len(memdHosts) == 0 {
			return errors.New("bootstrap_on=cccp but no CCCP/Memcached hosts in connection string")
		}
	case "both":
	case "":
		// Do nothing
		break
	default:
		return errors.New("bootstrap_on={http,cccp,both}")
	}
	config.MemdAddrs = memdHosts
	config.HTTPAddrs = httpHosts

	if spec.UseSsl {
		cacertpaths := spec.Options["ca_cert_path"]

		if len(cacertpaths) > 0 {
			roots := x509.NewCertPool()

			for _, path := range cacertpaths {
				cacert, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				ok := roots.AppendCertsFromPEM(cacert)
				if !ok {
					return errInvalidCertificate
				}
			}

			config.TLSRootCAProvider = func() *x509.CertPool {
				return roots
			}
		}

		config.UseTLS = true
	}

	if spec.Bucket != "" {
		config.BucketName = spec.Bucket
	}

	if valStr, ok := fetchOption("network"); ok {
		config.NetworkType = valStr
	}

	if valStr, ok := fetchOption("kv_connect_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("kv_connect_timeout option must be a duration or a number")
		}
		config.KVConnectTimeout = val
	}

	if valStr, ok := fetchOption("config_poll_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("config poll timeout option must be a duration or a number")
		}
		config.CccpMaxWait = val
	}

	if valStr, ok := fetchOption("config_poll_interval"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("config pool interval option must be duration or a number")
		}
		config.CccpPollPeriod = val
	}

	if valStr, ok := fetchOption("enable_mutation_tokens"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return fmt.Errorf("enable_mutation_tokens option must be a boolean")
		}
		config.UseMutationTokens = val
	}

	if valStr, ok := fetchOption("compression"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return fmt.Errorf("compression option must be a boolean")
		}
		config.UseCompression = val
	}

	if valStr, ok := fetchOption("compression_min_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("compression_min_size option must be an int")
		}
		config.CompressionMinSize = int(val)
	}

	if valStr, ok := fetchOption("compression_min_ratio"); ok {
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return fmt.Errorf("compression_min_size option must be an int")
		}
		config.CompressionMinRatio = val
	}

	if valStr, ok := fetchOption("enable_server_durations"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return fmt.Errorf("server_duration option must be a boolean")
		}
		config.UseDurations = val
	}

	if valStr, ok := fetchOption("max_idle_http_connections"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("http max idle connections option must be a number")
		}
		config.HTTPMaxIdleConns = int(val)
	}

	if valStr, ok := fetchOption("max_perhost_idle_http_connections"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("max_perhost_idle_http_connections option must be a number")
		}
		config.HTTPMaxIdleConnsPerHost = int(val)
	}

	if valStr, ok := fetchOption("idle_http_connection_timeout"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("idle_http_connection_timeout option must be a duration or a number")
		}
		config.HTTPIdleConnectionTimeout = val
	}

	if valStr, ok := fetchOption("orphaned_response_logging"); ok {
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return fmt.Errorf("orphaned_response_logging option must be a boolean")
		}
		config.UseZombieLogger = val
	}

	if valStr, ok := fetchOption("orphaned_response_logging_interval"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("orphaned_response_logging_interval option must be a number")
		}
		config.ZombieLoggerInterval = val
	}

	if valStr, ok := fetchOption("orphaned_response_logging_sample_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("orphaned_response_logging_sample_size option must be a number")
		}
		config.ZombieLoggerSampleSize = int(val)
	}

	// This option is experimental
	if valStr, ok := fetchOption("http_redial_period"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("http redial period option must be a duration or a number")
		}
		config.HTTPRedialPeriod = val
	}

	// This option is experimental
	if valStr, ok := fetchOption("http_retry_delay"); ok {
		val, err := parseDurationOrInt(valStr)
		if err != nil {
			return fmt.Errorf("http retry delay option must be a duration or a number")
		}
		config.HTTPRetryDelay = val
	}

	// This option is experimental
	if valStr, ok := fetchOption("kv_pool_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("kv pool size option must be a number")
		}
		config.KvPoolSize = int(val)
	}

	// This option is experimental
	if valStr, ok := fetchOption("max_queue_size"); ok {
		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return fmt.Errorf("max queue size option must be a number")
		}
		config.MaxQueueSize = int(val)
	}

	return nil
}
