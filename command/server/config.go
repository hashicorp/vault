// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/experiments"
	"github.com/hashicorp/vault/helper/osutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/testcluster"
	"github.com/mitchellh/mapstructure"
)

const (
	VaultDevCAFilename   = "vault-ca.pem"
	VaultDevCertFilename = "vault-cert.pem"
	VaultDevKeyFilename  = "vault-key.pem"
)

var (
	entConfigValidate = func(_ *Config, _ string) []configutil.ConfigError {
		return nil
	}

	// Modified internally for testing.
	validExperiments = experiments.ValidExperiments()
)

// Config is the configuration for the vault server.
type Config struct {
	UnusedKeys configutil.UnusedKeyMap `hcl:",unusedKeyPositions"`
	FoundKeys  []string                `hcl:",decodedFields"`
	entConfig

	*configutil.SharedConfig `hcl:"-"`

	Storage   *Storage `hcl:"-"`
	HAStorage *Storage `hcl:"-"`

	ServiceRegistration *ServiceRegistration `hcl:"-"`

	Experiments []string `hcl:"experiments"`

	CacheSize                int         `hcl:"cache_size"`
	DisableCache             bool        `hcl:"-"`
	DisableCacheRaw          interface{} `hcl:"disable_cache"`
	DisablePrintableCheck    bool        `hcl:"-"`
	DisablePrintableCheckRaw interface{} `hcl:"disable_printable_check"`

	EnableUI    bool        `hcl:"-"`
	EnableUIRaw interface{} `hcl:"ui"`

	MaxLeaseTTL        time.Duration `hcl:"-"`
	MaxLeaseTTLRaw     interface{}   `hcl:"max_lease_ttl,alias:MaxLeaseTTL"`
	DefaultLeaseTTL    time.Duration `hcl:"-"`
	DefaultLeaseTTLRaw interface{}   `hcl:"default_lease_ttl,alias:DefaultLeaseTTL"`

	ClusterCipherSuites string `hcl:"cluster_cipher_suites"`

	PluginDirectory string `hcl:"plugin_directory"`

	PluginFileUid int `hcl:"plugin_file_uid"`

	PluginFilePermissions    int         `hcl:"-"`
	PluginFilePermissionsRaw interface{} `hcl:"plugin_file_permissions,alias:PluginFilePermissions"`

	EnableIntrospectionEndpoint    bool        `hcl:"-"`
	EnableIntrospectionEndpointRaw interface{} `hcl:"introspection_endpoint,alias:EnableIntrospectionEndpoint"`

	EnableRawEndpoint    bool        `hcl:"-"`
	EnableRawEndpointRaw interface{} `hcl:"raw_storage_endpoint,alias:EnableRawEndpoint"`

	APIAddr              string      `hcl:"api_addr"`
	ClusterAddr          string      `hcl:"cluster_addr"`
	DisableClustering    bool        `hcl:"-"`
	DisableClusteringRaw interface{} `hcl:"disable_clustering,alias:DisableClustering"`

	DisablePerformanceStandby    bool        `hcl:"-"`
	DisablePerformanceStandbyRaw interface{} `hcl:"disable_performance_standby,alias:DisablePerformanceStandby"`

	DisableSealWrap    bool        `hcl:"-"`
	DisableSealWrapRaw interface{} `hcl:"disable_sealwrap,alias:DisableSealWrap"`

	DisableIndexing    bool        `hcl:"-"`
	DisableIndexingRaw interface{} `hcl:"disable_indexing,alias:DisableIndexing"`

	DisableSentinelTrace    bool        `hcl:"-"`
	DisableSentinelTraceRaw interface{} `hcl:"disable_sentinel_trace,alias:DisableSentinelTrace"`

	EnableResponseHeaderHostname    bool        `hcl:"-"`
	EnableResponseHeaderHostnameRaw interface{} `hcl:"enable_response_header_hostname"`

	LogRequestsLevel    string      `hcl:"-"`
	LogRequestsLevelRaw interface{} `hcl:"log_requests_level"`

	DetectDeadlocks string `hcl:"detect_deadlocks"`

	EnableResponseHeaderRaftNodeID    bool        `hcl:"-"`
	EnableResponseHeaderRaftNodeIDRaw interface{} `hcl:"enable_response_header_raft_node_id"`

	License          string `hcl:"-"`
	LicensePath      string `hcl:"license_path"`
	DisableSSCTokens bool   `hcl:"-"`
}

const (
	sectionSeal = "Seal"
)

func (c *Config) Validate(sourceFilePath string) []configutil.ConfigError {
	results := configutil.ValidateUnusedFields(c.UnusedKeys, sourceFilePath)
	if c.Telemetry != nil {
		results = append(results, c.Telemetry.Validate(sourceFilePath)...)
	}
	if c.ServiceRegistration != nil {
		results = append(results, c.ServiceRegistration.Validate(sourceFilePath)...)
	}
	for _, l := range c.Listeners {
		results = append(results, l.Validate(sourceFilePath)...)
	}
	results = append(results, c.validateEnt(sourceFilePath)...)
	return results
}

func (c *Config) validateEnt(sourceFilePath string) []configutil.ConfigError {
	return entConfigValidate(c, sourceFilePath)
}

// DevConfig is a Config that is used for dev mode of Vault.
func DevConfig(storageType string) (*Config, error) {
	hclStr := `
disable_mlock = true

listener "tcp" {
	address = "127.0.0.1:8200"
	tls_disable = true
	proxy_protocol_behavior = "allow_authorized"
	proxy_protocol_authorized_addrs = "127.0.0.1:8200"
}

telemetry {
	prometheus_retention_time = "24h"
	disable_hostname = true
}
enable_raw_endpoint = true

storage "%s" {
}

ui = true
`

	hclStr = fmt.Sprintf(hclStr, storageType)
	parsed, err := ParseConfig(hclStr, "")
	if err != nil {
		return nil, fmt.Errorf("error parsing dev config: %w", err)
	}
	return parsed, nil
}

// DevTLSConfig is a Config that is used for dev tls mode of Vault.
func DevTLSConfig(storageType, certDir string, extraSANs []string) (*Config, error) {
	ca, err := GenerateCA()
	if err != nil {
		return nil, err
	}

	cert, key, err := generateCert(ca.Template, ca.Signer, extraSANs)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(fmt.Sprintf("%s/%s", certDir, VaultDevCAFilename), []byte(ca.PEM), 0o444); err != nil {
		return nil, err
	}

	if err := os.WriteFile(fmt.Sprintf("%s/%s", certDir, VaultDevCertFilename), []byte(cert), 0o400); err != nil {
		return nil, err
	}

	if err := os.WriteFile(fmt.Sprintf("%s/%s", certDir, VaultDevKeyFilename), []byte(key), 0o400); err != nil {
		return nil, err
	}
	return parseDevTLSConfig(storageType, certDir)
}

func parseDevTLSConfig(storageType, certDir string) (*Config, error) {
	hclStr := `
disable_mlock = true

listener "tcp" {
	address = "[::]:8200"
	tls_cert_file = "%s/vault-cert.pem"
	tls_key_file = "%s/vault-key.pem"
	proxy_protocol_behavior = "allow_authorized"
	proxy_protocol_authorized_addrs = "[::]:8200"
}

telemetry {
	prometheus_retention_time = "24h"
	disable_hostname = true
}
enable_raw_endpoint = true

storage "%s" {
}

ui = true
`
	certDirEscaped := strings.Replace(certDir, "\\", "\\\\", -1)
	hclStr = fmt.Sprintf(hclStr, certDirEscaped, certDirEscaped, storageType)
	parsed, err := ParseConfig(hclStr, "")
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// Storage is the underlying storage configuration for the server.
type Storage struct {
	Type              string
	RedirectAddr      string
	ClusterAddr       string
	DisableClustering bool
	Config            map[string]string
}

func (b *Storage) GoString() string {
	return fmt.Sprintf("*%#v", *b)
}

// ServiceRegistration is the optional service discovery for the server.
type ServiceRegistration struct {
	UnusedKeys configutil.UnusedKeyMap `hcl:",unusedKeyPositions"`
	Type       string
	Config     map[string]string
}

func (b *ServiceRegistration) Validate(source string) []configutil.ConfigError {
	return configutil.ValidateUnusedFields(b.UnusedKeys, source)
}

func (b *ServiceRegistration) GoString() string {
	return fmt.Sprintf("*%#v", *b)
}

func NewConfig() *Config {
	return &Config{
		SharedConfig: new(configutil.SharedConfig),
	}
}

// Merge merges two configurations.
func (c *Config) Merge(c2 *Config) *Config {
	if c2 == nil {
		return c
	}

	result := NewConfig()

	result.SharedConfig = c.SharedConfig
	if c2.SharedConfig != nil {
		result.SharedConfig = c.SharedConfig.Merge(c2.SharedConfig)
	}

	result.Storage = c.Storage
	if c2.Storage != nil {
		result.Storage = c2.Storage
	}

	result.HAStorage = c.HAStorage
	if c2.HAStorage != nil {
		result.HAStorage = c2.HAStorage
	}

	result.ServiceRegistration = c.ServiceRegistration
	if c2.ServiceRegistration != nil {
		result.ServiceRegistration = c2.ServiceRegistration
	}

	result.CacheSize = c.CacheSize
	if c2.CacheSize != 0 {
		result.CacheSize = c2.CacheSize
	}

	// merging these booleans via an OR operation
	result.DisableCache = c.DisableCache
	if c2.DisableCache {
		result.DisableCache = c2.DisableCache
	}

	result.DisableSentinelTrace = c.DisableSentinelTrace
	if c2.DisableSentinelTrace {
		result.DisableSentinelTrace = c2.DisableSentinelTrace
	}

	result.DisablePrintableCheck = c.DisablePrintableCheck
	if c2.DisablePrintableCheckRaw != nil {
		result.DisablePrintableCheck = c2.DisablePrintableCheck
	}

	// merge these integers via a MAX operation
	result.MaxLeaseTTL = c.MaxLeaseTTL
	if c2.MaxLeaseTTL > result.MaxLeaseTTL {
		result.MaxLeaseTTL = c2.MaxLeaseTTL
	}

	result.DefaultLeaseTTL = c.DefaultLeaseTTL
	if c2.DefaultLeaseTTL > result.DefaultLeaseTTL {
		result.DefaultLeaseTTL = c2.DefaultLeaseTTL
	}

	result.ClusterCipherSuites = c.ClusterCipherSuites
	if c2.ClusterCipherSuites != "" {
		result.ClusterCipherSuites = c2.ClusterCipherSuites
	}

	result.EnableUI = c.EnableUI
	if c2.EnableUI {
		result.EnableUI = c2.EnableUI
	}

	result.EnableRawEndpoint = c.EnableRawEndpoint
	if c2.EnableRawEndpoint {
		result.EnableRawEndpoint = c2.EnableRawEndpoint
	}

	result.EnableIntrospectionEndpoint = c.EnableIntrospectionEndpoint
	if c2.EnableIntrospectionEndpoint {
		result.EnableIntrospectionEndpoint = c2.EnableIntrospectionEndpoint
	}

	result.APIAddr = c.APIAddr
	if c2.APIAddr != "" {
		result.APIAddr = c2.APIAddr
	}

	result.ClusterAddr = c.ClusterAddr
	if c2.ClusterAddr != "" {
		result.ClusterAddr = c2.ClusterAddr
	}

	// Retain raw value so that it can be assigned to storage objects
	result.DisableClustering = c.DisableClustering
	result.DisableClusteringRaw = c.DisableClusteringRaw
	if c2.DisableClusteringRaw != nil {
		result.DisableClustering = c2.DisableClustering
		result.DisableClusteringRaw = c2.DisableClusteringRaw
	}

	result.PluginDirectory = c.PluginDirectory
	if c2.PluginDirectory != "" {
		result.PluginDirectory = c2.PluginDirectory
	}

	result.PluginFileUid = c.PluginFileUid
	if c2.PluginFileUid != 0 {
		result.PluginFileUid = c2.PluginFileUid
	}

	result.PluginFilePermissions = c.PluginFilePermissions
	if c2.PluginFilePermissionsRaw != nil {
		result.PluginFilePermissions = c2.PluginFilePermissions
		result.PluginFilePermissionsRaw = c2.PluginFilePermissionsRaw
	}

	result.DisablePerformanceStandby = c.DisablePerformanceStandby
	if c2.DisablePerformanceStandby {
		result.DisablePerformanceStandby = c2.DisablePerformanceStandby
	}

	result.DisableSealWrap = c.DisableSealWrap
	if c2.DisableSealWrap {
		result.DisableSealWrap = c2.DisableSealWrap
	}

	result.DisableIndexing = c.DisableIndexing
	if c2.DisableIndexing {
		result.DisableIndexing = c2.DisableIndexing
	}

	result.EnableResponseHeaderHostname = c.EnableResponseHeaderHostname
	if c2.EnableResponseHeaderHostname {
		result.EnableResponseHeaderHostname = c2.EnableResponseHeaderHostname
	}

	result.LogRequestsLevel = c.LogRequestsLevel
	if c2.LogRequestsLevel != "" {
		result.LogRequestsLevel = c2.LogRequestsLevel
	}

	result.DetectDeadlocks = c.DetectDeadlocks
	if c2.DetectDeadlocks != "" {
		result.DetectDeadlocks = c2.DetectDeadlocks
	}

	result.EnableResponseHeaderRaftNodeID = c.EnableResponseHeaderRaftNodeID
	if c2.EnableResponseHeaderRaftNodeID {
		result.EnableResponseHeaderRaftNodeID = c2.EnableResponseHeaderRaftNodeID
	}

	result.LicensePath = c.LicensePath
	if c2.LicensePath != "" {
		result.LicensePath = c2.LicensePath
	}

	// Use values from top-level configuration for storage if set
	if storage := result.Storage; storage != nil {
		if result.APIAddr != "" {
			storage.RedirectAddr = result.APIAddr
		}
		if result.ClusterAddr != "" {
			storage.ClusterAddr = result.ClusterAddr
		}
		if result.DisableClusteringRaw != nil {
			storage.DisableClustering = result.DisableClustering
		}
	}

	if haStorage := result.HAStorage; haStorage != nil {
		if result.APIAddr != "" {
			haStorage.RedirectAddr = result.APIAddr
		}
		if result.ClusterAddr != "" {
			haStorage.ClusterAddr = result.ClusterAddr
		}
		if result.DisableClusteringRaw != nil {
			haStorage.DisableClustering = result.DisableClustering
		}
	}

	result.AdministrativeNamespacePath = c.AdministrativeNamespacePath
	if c2.AdministrativeNamespacePath != "" {
		result.AdministrativeNamespacePath = c2.AdministrativeNamespacePath
	}

	result.entConfig = c.entConfig.Merge(c2.entConfig)

	result.Experiments = mergeExperiments(c.Experiments, c2.Experiments)

	return result
}

// LoadConfig loads the configuration at the given path, regardless if
// its a file or directory.
func LoadConfig(path string) (*Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		// check permissions on the config directory
		var enableFilePermissionsCheck bool
		if enableFilePermissionsCheckEnv := os.Getenv(consts.VaultEnableFilePermissionsCheckEnv); enableFilePermissionsCheckEnv != "" {
			var err error
			enableFilePermissionsCheck, err = strconv.ParseBool(enableFilePermissionsCheckEnv)
			if err != nil {
				return nil, errors.New("Error parsing the environment variable VAULT_ENABLE_FILE_PERMISSIONS_CHECK")
			}
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if enableFilePermissionsCheck {
			err = osutil.OwnerPermissionsMatchFile(f, 0, 0)
			if err != nil {
				return nil, err
			}
		}
		return CheckConfig(LoadConfigDir(path))
	}
	return CheckConfig(LoadConfigFile(path))
}

func CheckConfig(c *Config, e error) (*Config, error) {
	if e != nil {
		return c, e
	}

	if err := c.checkSealConfig(); err != nil {
		return nil, err
	}

	sealMap := make(map[string]*configutil.KMS)
	for _, seal := range c.Seals {
		if seal.Name == "" {
			return nil, errors.New("seals: seal name is empty")
		}

		if _, ok := sealMap[seal.Name]; ok {
			return nil, errors.New("seals: seal names must be unique")
		}

		sealMap[seal.Name] = seal
	}

	return c, nil
}

// LoadConfigFile loads the configuration from the given file.
func LoadConfigFile(path string) (*Config, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Read the file
	d, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	conf, err := ParseConfig(string(d), path)
	if err != nil {
		return nil, err
	}

	var enableFilePermissionsCheck bool
	if enableFilePermissionsCheckEnv := os.Getenv(consts.VaultEnableFilePermissionsCheckEnv); enableFilePermissionsCheckEnv != "" {
		var err error
		enableFilePermissionsCheck, err = strconv.ParseBool(enableFilePermissionsCheckEnv)
		if err != nil {
			return nil, errors.New("Error parsing the environment variable VAULT_ENABLE_FILE_PERMISSIONS_CHECK")
		}
	}

	if enableFilePermissionsCheck {
		// check permissions of the config file
		err = osutil.OwnerPermissionsMatchFile(f, 0, 0)
		if err != nil {
			return nil, err
		}
		// check permissions of the plugin directory
		if conf.PluginDirectory != "" {

			err = osutil.OwnerPermissionsMatch(conf.PluginDirectory, conf.PluginFileUid, conf.PluginFilePermissions)
			if err != nil {
				return nil, err
			}
		}
	}
	return conf, nil
}

func ParseConfig(d, source string) (*Config, error) {
	// Parse!
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, err
	}

	// Start building the result
	result := NewConfig()
	if err := hcl.DecodeObject(result, obj); err != nil {
		return nil, err
	}

	if rendered, err := configutil.ParseSingleIPTemplate(result.APIAddr); err != nil {
		return nil, err
	} else {
		result.APIAddr = rendered
	}
	if rendered, err := configutil.ParseSingleIPTemplate(result.ClusterAddr); err != nil {
		return nil, err
	} else {
		result.ClusterAddr = rendered
	}

	sharedConfig, err := configutil.ParseConfig(d)
	if err != nil {
		return nil, err
	}
	result.SharedConfig = sharedConfig

	if result.MaxLeaseTTLRaw != nil {
		if result.MaxLeaseTTL, err = parseutil.ParseDurationSecond(result.MaxLeaseTTLRaw); err != nil {
			return nil, err
		}
	}
	if result.DefaultLeaseTTLRaw != nil {
		if result.DefaultLeaseTTL, err = parseutil.ParseDurationSecond(result.DefaultLeaseTTLRaw); err != nil {
			return nil, err
		}
	}

	if result.EnableUIRaw != nil {
		if result.EnableUI, err = parseutil.ParseBool(result.EnableUIRaw); err != nil {
			return nil, err
		}
	}

	if result.DisableCacheRaw != nil {
		if result.DisableCache, err = parseutil.ParseBool(result.DisableCacheRaw); err != nil {
			return nil, err
		}
	}

	if result.DisablePrintableCheckRaw != nil {
		if result.DisablePrintableCheck, err = parseutil.ParseBool(result.DisablePrintableCheckRaw); err != nil {
			return nil, err
		}
	}

	if result.EnableRawEndpointRaw != nil {
		if result.EnableRawEndpoint, err = parseutil.ParseBool(result.EnableRawEndpointRaw); err != nil {
			return nil, err
		}
	}

	if result.EnableIntrospectionEndpointRaw != nil {
		if result.EnableIntrospectionEndpoint, err = parseutil.ParseBool(result.EnableIntrospectionEndpointRaw); err != nil {
			return nil, err
		}
	}

	if result.DisableClusteringRaw != nil {
		if result.DisableClustering, err = parseutil.ParseBool(result.DisableClusteringRaw); err != nil {
			return nil, err
		}
	}

	if result.PluginFilePermissionsRaw != nil {
		octalPermissionsString, err := parseutil.ParseString(result.PluginFilePermissionsRaw)
		if err != nil {
			return nil, err
		}
		pluginFilePermissions, err := strconv.ParseInt(octalPermissionsString, 8, 64)
		if err != nil {
			return nil, err
		}
		if pluginFilePermissions < math.MinInt || pluginFilePermissions > math.MaxInt {
			return nil, fmt.Errorf("file permission value %v cannot be safely cast to int: exceeds bounds (%v, %v)", pluginFilePermissions, math.MinInt, math.MaxInt)
		}
		result.PluginFilePermissions = int(pluginFilePermissions)
	}

	if result.DisableSentinelTraceRaw != nil {
		if result.DisableSentinelTrace, err = parseutil.ParseBool(result.DisableSentinelTraceRaw); err != nil {
			return nil, err
		}
	}

	if result.DisablePerformanceStandbyRaw != nil {
		if result.DisablePerformanceStandby, err = parseutil.ParseBool(result.DisablePerformanceStandbyRaw); err != nil {
			return nil, err
		}
	}

	if result.DisableSealWrapRaw != nil {
		if result.DisableSealWrap, err = parseutil.ParseBool(result.DisableSealWrapRaw); err != nil {
			return nil, err
		}
	}

	if result.DisableIndexingRaw != nil {
		if result.DisableIndexing, err = parseutil.ParseBool(result.DisableIndexingRaw); err != nil {
			return nil, err
		}
	}

	if result.EnableResponseHeaderHostnameRaw != nil {
		if result.EnableResponseHeaderHostname, err = parseutil.ParseBool(result.EnableResponseHeaderHostnameRaw); err != nil {
			return nil, err
		}
	}

	if result.LogRequestsLevelRaw != nil {
		result.LogRequestsLevel = strings.ToLower(strings.TrimSpace(result.LogRequestsLevelRaw.(string)))
		result.LogRequestsLevelRaw = ""
	}

	if result.EnableResponseHeaderRaftNodeIDRaw != nil {
		if result.EnableResponseHeaderRaftNodeID, err = parseutil.ParseBool(result.EnableResponseHeaderRaftNodeIDRaw); err != nil {
			return nil, err
		}
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	// Look for storage but still support old backend
	if o := list.Filter("storage"); len(o.Items) > 0 {
		delete(result.UnusedKeys, "storage")
		if err := ParseStorage(result, o, "storage"); err != nil {
			return nil, fmt.Errorf("error parsing 'storage': %w", err)
		}
		result.found(result.Storage.Type, result.Storage.Type)
	} else {
		delete(result.UnusedKeys, "backend")
		if o := list.Filter("backend"); len(o.Items) > 0 {
			if err := ParseStorage(result, o, "backend"); err != nil {
				return nil, fmt.Errorf("error parsing 'backend': %w", err)
			}
		}
	}

	if o := list.Filter("ha_storage"); len(o.Items) > 0 {
		delete(result.UnusedKeys, "ha_storage")
		if err := parseHAStorage(result, o, "ha_storage"); err != nil {
			return nil, fmt.Errorf("error parsing 'ha_storage': %w", err)
		}
	} else {
		if o := list.Filter("ha_backend"); len(o.Items) > 0 {
			delete(result.UnusedKeys, "ha_backend")
			if err := parseHAStorage(result, o, "ha_backend"); err != nil {
				return nil, fmt.Errorf("error parsing 'ha_backend': %w", err)
			}
		}
	}

	// Parse service discovery
	if o := list.Filter("service_registration"); len(o.Items) > 0 {
		delete(result.UnusedKeys, "service_registration")
		if err := parseServiceRegistration(result, o, "service_registration"); err != nil {
			return nil, fmt.Errorf("error parsing 'service_registration': %w", err)
		}
	}

	if err := validateExperiments(result.Experiments); err != nil {
		return nil, fmt.Errorf("error validating experiment(s) from config: %w", err)
	}

	if err := result.parseConfig(list); err != nil {
		return nil, fmt.Errorf("error parsing enterprise config: %w", err)
	}

	// Remove all unused keys from Config that were satisfied by SharedConfig.
	result.UnusedKeys = configutil.UnusedFieldDifference(result.UnusedKeys, nil, append(result.FoundKeys, sharedConfig.FoundKeys...))
	// Assign file info
	for _, v := range result.UnusedKeys {
		for i := range v {
			v[i].Filename = source
		}
	}

	return result, nil
}

func ExperimentsFromEnvAndCLI(config *Config, envKey string, flagExperiments []string) error {
	if envExperimentsRaw := os.Getenv(envKey); envExperimentsRaw != "" {
		envExperiments := strings.Split(envExperimentsRaw, ",")
		err := validateExperiments(envExperiments)
		if err != nil {
			return fmt.Errorf("error validating experiment(s) from environment variable %q: %w", envKey, err)
		}

		config.Experiments = mergeExperiments(config.Experiments, envExperiments)
	}

	if len(flagExperiments) != 0 {
		err := validateExperiments(flagExperiments)
		if err != nil {
			return fmt.Errorf("error validating experiment(s) from command line flag: %w", err)
		}

		config.Experiments = mergeExperiments(config.Experiments, flagExperiments)
	}

	return nil
}

// Validate checks each experiment is a known experiment.
func validateExperiments(experiments []string) error {
	var invalid []string

	for _, experiment := range experiments {
		if !strutil.StrListContains(validExperiments, experiment) {
			invalid = append(invalid, experiment)
		}
	}

	if len(invalid) != 0 {
		return fmt.Errorf("valid experiment(s) are %s, but received the following invalid experiment(s): %s",
			strings.Join(validExperiments, ", "),
			strings.Join(invalid, ", "))
	}

	return nil
}

// mergeExperiments returns the logical OR of the two sets.
func mergeExperiments(left, right []string) []string {
	processed := map[string]struct{}{}
	var result []string
	for _, l := range left {
		if _, seen := processed[l]; !seen {
			result = append(result, l)
		}
		processed[l] = struct{}{}
	}

	for _, r := range right {
		if _, seen := processed[r]; !seen {
			result = append(result, r)
			processed[r] = struct{}{}
		}
	}

	return result
}

// LoadConfigDir loads all the configurations in the given directory
// in alphabetical order.
func LoadConfigDir(dir string) (*Config, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("configuration path must be a directory: %q", dir)
	}

	var files []string
	err = nil
	for err != io.EOF {
		var fis []os.FileInfo
		fis, err = f.Readdir(128)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for _, fi := range fis {
			// Ignore directories
			if fi.IsDir() {
				continue
			}

			// Only care about files that are valid to load.
			name := fi.Name()
			skip := true
			if strings.HasSuffix(name, ".hcl") {
				skip = false
			} else if strings.HasSuffix(name, ".json") {
				skip = false
			}
			if skip || isTemporaryFile(name) {
				continue
			}

			path := filepath.Join(dir, name)
			files = append(files, path)
		}
	}

	result := NewConfig()
	for _, f := range files {
		config, err := LoadConfigFile(f)
		if err != nil {
			return nil, fmt.Errorf("error loading %q: %w", f, err)
		}

		if result == nil {
			result = config
		} else {
			result = result.Merge(config)
		}
	}

	return result, nil
}

// isTemporaryFile returns true or false depending on whether the
// provided file name is a temporary file for the following editors:
// emacs or vim.
func isTemporaryFile(name string) bool {
	return strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, ".#") || // emacs
		(strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#")) // emacs
}

func ParseStorage(result *Config, list *ast.ObjectList, name string) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block is permitted", name)
	}

	// Get our item
	item := list.Items[0]

	key := name
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	var config map[string]interface{}
	if err := hcl.DecodeObject(&config, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("%s.%s:", name, key))
	}

	m := make(map[string]string)
	for key, val := range config {
		valStr, ok := val.(string)
		if ok {
			m[key] = valStr
			continue
		}
		valBytes, err := json.Marshal(val)
		if err != nil {
			return err
		}
		m[key] = string(valBytes)
	}

	// Pull out the redirect address since it's common to all backends
	var redirectAddr string
	if v, ok := m["redirect_addr"]; ok {
		redirectAddr = v
		delete(m, "redirect_addr")
	} else if v, ok := m["advertise_addr"]; ok {
		redirectAddr = v
		delete(m, "advertise_addr")
	}

	// Pull out the cluster address since it's common to all backends
	var clusterAddr string
	if v, ok := m["cluster_addr"]; ok {
		clusterAddr = v
		delete(m, "cluster_addr")
	}

	var disableClustering bool
	var err error
	if v, ok := m["disable_clustering"]; ok {
		disableClustering, err = strconv.ParseBool(v)
		if err != nil {
			return multierror.Prefix(err, fmt.Sprintf("%s.%s:", name, key))
		}
		delete(m, "disable_clustering")
	}

	// Override with top-level values if they are set
	if result.APIAddr != "" {
		redirectAddr = result.APIAddr
	}

	if result.ClusterAddr != "" {
		clusterAddr = result.ClusterAddr
	}

	if result.DisableClusteringRaw != nil {
		disableClustering = result.DisableClustering
	}

	result.Storage = &Storage{
		RedirectAddr:      redirectAddr,
		ClusterAddr:       clusterAddr,
		DisableClustering: disableClustering,
		Type:              strings.ToLower(key),
		Config:            m,
	}
	return nil
}

func parseHAStorage(result *Config, list *ast.ObjectList, name string) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block is permitted", name)
	}

	// Get our item
	item := list.Items[0]

	key := name
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	var config map[string]interface{}
	if err := hcl.DecodeObject(&config, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("%s.%s:", name, key))
	}

	m := make(map[string]string)
	for key, val := range config {
		valStr, ok := val.(string)
		if ok {
			m[key] = valStr
			continue
		}
		valBytes, err := json.Marshal(val)
		if err != nil {
			return err
		}
		m[key] = string(valBytes)
	}

	// Pull out the redirect address since it's common to all backends
	var redirectAddr string
	if v, ok := m["redirect_addr"]; ok {
		redirectAddr = v
		delete(m, "redirect_addr")
	} else if v, ok := m["advertise_addr"]; ok {
		redirectAddr = v
		delete(m, "advertise_addr")
	}

	// Pull out the cluster address since it's common to all backends
	var clusterAddr string
	if v, ok := m["cluster_addr"]; ok {
		clusterAddr = v
		delete(m, "cluster_addr")
	}

	var disableClustering bool
	var err error
	if v, ok := m["disable_clustering"]; ok {
		disableClustering, err = strconv.ParseBool(v)
		if err != nil {
			return multierror.Prefix(err, fmt.Sprintf("%s.%s:", name, key))
		}
		delete(m, "disable_clustering")
	}

	// Override with top-level values if they are set
	if result.APIAddr != "" {
		redirectAddr = result.APIAddr
	}

	if result.ClusterAddr != "" {
		clusterAddr = result.ClusterAddr
	}

	if result.DisableClusteringRaw != nil {
		disableClustering = result.DisableClustering
	}

	result.HAStorage = &Storage{
		RedirectAddr:      redirectAddr,
		ClusterAddr:       clusterAddr,
		DisableClustering: disableClustering,
		Type:              strings.ToLower(key),
		Config:            m,
	}
	return nil
}

func parseServiceRegistration(result *Config, list *ast.ObjectList, name string) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block is permitted", name)
	}

	// Get our item
	item := list.Items[0]
	key := name
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	var m map[string]string
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("%s.%s:", name, key))
	}

	result.ServiceRegistration = &ServiceRegistration{
		Type:   strings.ToLower(key),
		Config: m,
	}
	return nil
}

// Sanitized returns a copy of the config with all values that are considered
// sensitive stripped. It also strips all `*Raw` values that are mainly
// used for parsing.
//
// Specifically, the fields that this method strips are:
// - Storage.Config
// - HAStorage.Config
// - Seals.Config
// - Telemetry.CirconusAPIToken
func (c *Config) Sanitized() map[string]interface{} {
	// Create shared config if it doesn't exist (e.g. in tests) so that map
	// keys are actually populated
	if c.SharedConfig == nil {
		c.SharedConfig = new(configutil.SharedConfig)
	}
	sharedResult := c.SharedConfig.Sanitized()
	result := map[string]interface{}{
		"cache_size":              c.CacheSize,
		"disable_sentinel_trace":  c.DisableSentinelTrace,
		"disable_cache":           c.DisableCache,
		"disable_printable_check": c.DisablePrintableCheck,

		"enable_ui": c.EnableUI,

		"max_lease_ttl":     c.MaxLeaseTTL / time.Second,
		"default_lease_ttl": c.DefaultLeaseTTL / time.Second,

		"cluster_cipher_suites": c.ClusterCipherSuites,

		"plugin_directory": c.PluginDirectory,

		"plugin_file_uid": c.PluginFileUid,

		"plugin_file_permissions": c.PluginFilePermissions,

		"raw_storage_endpoint": c.EnableRawEndpoint,

		"introspection_endpoint": c.EnableIntrospectionEndpoint,

		"api_addr":           c.APIAddr,
		"cluster_addr":       c.ClusterAddr,
		"disable_clustering": c.DisableClustering,

		"disable_performance_standby": c.DisablePerformanceStandby,

		"disable_sealwrap": c.DisableSealWrap,

		"disable_indexing": c.DisableIndexing,

		"enable_response_header_hostname": c.EnableResponseHeaderHostname,

		"enable_response_header_raft_node_id": c.EnableResponseHeaderRaftNodeID,

		"log_requests_level": c.LogRequestsLevel,
		"experiments":        c.Experiments,

		"detect_deadlocks": c.DetectDeadlocks,
	}
	for k, v := range sharedResult {
		result[k] = v
	}

	// Sanitize storage stanza
	if c.Storage != nil {
		storageType := c.Storage.Type
		sanitizedStorage := map[string]interface{}{
			"type":               storageType,
			"redirect_addr":      c.Storage.RedirectAddr,
			"cluster_addr":       c.Storage.ClusterAddr,
			"disable_clustering": c.Storage.DisableClustering,
		}

		if storageType == "raft" {
			sanitizedStorage["raft"] = map[string]interface{}{
				"max_entry_size": c.Storage.Config["max_entry_size"],
			}
		}

		result["storage"] = sanitizedStorage
	}

	// Sanitize HA storage stanza
	if c.HAStorage != nil {
		haStorageType := c.HAStorage.Type
		sanitizedHAStorage := map[string]interface{}{
			"type":               haStorageType,
			"redirect_addr":      c.HAStorage.RedirectAddr,
			"cluster_addr":       c.HAStorage.ClusterAddr,
			"disable_clustering": c.HAStorage.DisableClustering,
		}

		if haStorageType == "raft" {
			sanitizedHAStorage["raft"] = map[string]interface{}{
				"max_entry_size": c.HAStorage.Config["max_entry_size"],
			}
		}

		result["ha_storage"] = sanitizedHAStorage
	}

	// Sanitize service_registration stanza
	if c.ServiceRegistration != nil {
		sanitizedServiceRegistration := map[string]interface{}{
			"type": c.ServiceRegistration.Type,
		}
		result["service_registration"] = sanitizedServiceRegistration
	}

	entConfigResult := c.entConfig.Sanitized()
	for k, v := range entConfigResult {
		result[k] = v
	}

	return result
}

func (c *Config) Prune() {
	for _, l := range c.Listeners {
		l.RawConfig = nil
		l.UnusedKeys = nil
	}
	c.FoundKeys = nil
	c.UnusedKeys = nil
	c.SharedConfig.FoundKeys = nil
	c.SharedConfig.UnusedKeys = nil
	if c.Telemetry != nil {
		c.Telemetry.FoundKeys = nil
		c.Telemetry.UnusedKeys = nil
	}
}

func (c *Config) found(s, k string) {
	delete(c.UnusedKeys, s)
	c.FoundKeys = append(c.FoundKeys, k)
}

func (c *Config) ToVaultNodeConfig() (*testcluster.VaultNodeConfig, error) {
	var vnc testcluster.VaultNodeConfig
	err := mapstructure.Decode(c, &vnc)
	if err != nil {
		return nil, err
	}
	return &vnc, nil
}
