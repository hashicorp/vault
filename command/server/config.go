package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

var entConfigValidate = func(_ *Config, _ string) []configutil.ConfigError {
	return nil
}

// Config is the configuration for the vault server.
type Config struct {
	UnusedKeys configutil.UnusedKeyMap `hcl:",unusedKeyPositions"`
	FoundKeys  []string                `hcl:",decodedFields"`
	entConfig

	*configutil.SharedConfig `hcl:"-"`

	Storage   *Storage `hcl:"-"`
	HAStorage *Storage `hcl:"-"`

	ServiceRegistration *ServiceRegistration `hcl:"-"`

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

	EnableResponseHeaderRaftNodeID    bool        `hcl:"-"`
	EnableResponseHeaderRaftNodeIDRaw interface{} `hcl:"enable_response_header_raft_node_id"`

	License     string `hcl:"-"`
	LicensePath string `hcl:"license_path"`
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

	result.entConfig = c.entConfig.Merge(c2.entConfig)

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
		return CheckConfig(LoadConfigDir(path))
	}
	return CheckConfig(LoadConfigFile(path))
}

func CheckConfig(c *Config, e error) (*Config, error) {
	if e != nil {
		return c, e
	}

	if len(c.Seals) == 2 {
		switch {
		case c.Seals[0].Disabled && c.Seals[1].Disabled:
			return nil, errors.New("seals: two seals provided but both are disabled")
		case !c.Seals[0].Disabled && !c.Seals[1].Disabled:
			return nil, errors.New("seals: two seals provided but neither is disabled")
		}
	}

	return c, nil
}

// LoadConfigFile loads the configuration from the given file.
func LoadConfigFile(path string) (*Config, error) {
	// Read the file
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf, err := ParseConfig(string(d), path)
	if err != nil {
		return nil, err
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

	if result.DisableClusteringRaw != nil {
		if result.DisableClustering, err = parseutil.ParseBool(result.DisableClusteringRaw); err != nil {
			return nil, err
		}
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

	entConfig := &(result.entConfig)
	if err := entConfig.parseConfig(list); err != nil {
		return nil, fmt.Errorf("error parsing enterprise config: %w", err)
	}

	// Remove all unused keys from Config that were satisfied by SharedConfig.
	result.UnusedKeys = configutil.UnusedFieldDifference(result.UnusedKeys, nil, append(result.FoundKeys, sharedConfig.FoundKeys...))
	// Assign file info
	for _, v := range result.UnusedKeys {
		for _, p := range v {
			p.Filename = source
		}
	}

	return result, nil
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

	var m map[string]string
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("%s.%s:", name, key))
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

		"max_lease_ttl":     c.MaxLeaseTTL,
		"default_lease_ttl": c.DefaultLeaseTTL,

		"cluster_cipher_suites": c.ClusterCipherSuites,

		"plugin_directory": c.PluginDirectory,

		"raw_storage_endpoint": c.EnableRawEndpoint,

		"api_addr":           c.APIAddr,
		"cluster_addr":       c.ClusterAddr,
		"disable_clustering": c.DisableClustering,

		"disable_performance_standby": c.DisablePerformanceStandby,

		"disable_sealwrap": c.DisableSealWrap,

		"disable_indexing": c.DisableIndexing,

		"enable_response_header_hostname": c.EnableResponseHeaderHostname,

		"enable_response_header_raft_node_id": c.EnableResponseHeaderRaftNodeID,
	}
	for k, v := range sharedResult {
		result[k] = v
	}

	// Sanitize storage stanza
	if c.Storage != nil {
		sanitizedStorage := map[string]interface{}{
			"type":               c.Storage.Type,
			"redirect_addr":      c.Storage.RedirectAddr,
			"cluster_addr":       c.Storage.ClusterAddr,
			"disable_clustering": c.Storage.DisableClustering,
		}
		result["storage"] = sanitizedStorage
	}

	// Sanitize HA storage stanza
	if c.HAStorage != nil {
		sanitizedHAStorage := map[string]interface{}{
			"type":               c.HAStorage.Type,
			"redirect_addr":      c.HAStorage.RedirectAddr,
			"cluster_addr":       c.HAStorage.ClusterAddr,
			"disable_clustering": c.HAStorage.DisableClustering,
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
