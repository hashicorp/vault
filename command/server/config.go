package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/parseutil"
)

// Config is the configuration for the vault server.
type Config struct {
	Listeners []*Listener `hcl:"-"`
	Storage   *Storage    `hcl:"-"`
	HAStorage *Storage    `hcl:"-"`

	Seal *Seal `hcl:"-"`

	CacheSize       int         `hcl:"cache_size"`
	DisableCache    bool        `hcl:"-"`
	DisableCacheRaw interface{} `hcl:"disable_cache"`
	DisableMlock    bool        `hcl:"-"`
	DisableMlockRaw interface{} `hcl:"disable_mlock"`

	EnableUI    bool        `hcl:"-"`
	EnableUIRaw interface{} `hcl:"ui"`

	Telemetry *Telemetry `hcl:"telemetry"`

	MaxLeaseTTL        time.Duration `hcl:"-"`
	MaxLeaseTTLRaw     interface{}   `hcl:"max_lease_ttl"`
	DefaultLeaseTTL    time.Duration `hcl:"-"`
	DefaultLeaseTTLRaw interface{}   `hcl:"default_lease_ttl"`

	ClusterName         string `hcl:"cluster_name"`
	ClusterCipherSuites string `hcl:"cluster_cipher_suites"`

	PluginDirectory string `hcl:"plugin_directory"`

	PidFile              string      `hcl:"pid_file"`
	EnableRawEndpoint    bool        `hcl:"-"`
	EnableRawEndpointRaw interface{} `hcl:"raw_storage_endpoint"`

	APIAddr              string      `hcl:"api_addr"`
	ClusterAddr          string      `hcl:"cluster_addr"`
	DisableClustering    bool        `hcl:"-"`
	DisableClusteringRaw interface{} `hcl:"disable_clustering"`
}

// DevConfig is a Config that is used for dev mode of Vault.
func DevConfig(ha, transactional bool) *Config {
	ret := &Config{
		DisableCache:      false,
		DisableMlock:      true,
		EnableRawEndpoint: true,

		Storage: &Storage{
			Type: "inmem",
		},

		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":                         "127.0.0.1:8200",
					"tls_disable":                     true,
					"proxy_protocol_behavior":         "allow_authorized",
					"proxy_protocol_authorized_addrs": "127.0.0.1:8200",
				},
			},
		},

		EnableUI: true,

		Telemetry: &Telemetry{},
	}

	switch {
	case ha && transactional:
		ret.Storage.Type = "inmem_transactional_ha"
	case !ha && transactional:
		ret.Storage.Type = "inmem_transactional"
	case ha && !transactional:
		ret.Storage.Type = "inmem_ha"
	}

	return ret
}

// Listener is the listener configuration for the server.
type Listener struct {
	Type   string
	Config map[string]interface{}
}

func (l *Listener) GoString() string {
	return fmt.Sprintf("*%#v", *l)
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

// Seal contains Seal configuration for the server
type Seal struct {
	Type   string
	Config map[string]string
}

func (h *Seal) GoString() string {
	return fmt.Sprintf("*%#v", *h)
}

// Telemetry is the telemetry configuration for the server
type Telemetry struct {
	StatsiteAddr string `hcl:"statsite_address"`
	StatsdAddr   string `hcl:"statsd_address"`

	DisableHostname bool `hcl:"disable_hostname"`

	// Circonus: see https://github.com/circonus-labs/circonus-gometrics
	// for more details on the various configuration options.
	// Valid configuration combinations:
	//    - CirconusAPIToken
	//      metric management enabled (search for existing check or create a new one)
	//    - CirconusSubmissionUrl
	//      metric management disabled (use check with specified submission_url,
	//      broker must be using a public SSL certificate)
	//    - CirconusAPIToken + CirconusCheckSubmissionURL
	//      metric management enabled (use check with specified submission_url)
	//    - CirconusAPIToken + CirconusCheckID
	//      metric management enabled (use check with specified id)

	// CirconusAPIToken is a valid API Token used to create/manage check. If provided,
	// metric management is enabled.
	// Default: none
	CirconusAPIToken string `hcl:"circonus_api_token"`
	// CirconusAPIApp is an app name associated with API token.
	// Default: "consul"
	CirconusAPIApp string `hcl:"circonus_api_app"`
	// CirconusAPIURL is the base URL to use for contacting the Circonus API.
	// Default: "https://api.circonus.com/v2"
	CirconusAPIURL string `hcl:"circonus_api_url"`
	// CirconusSubmissionInterval is the interval at which metrics are submitted to Circonus.
	// Default: 10s
	CirconusSubmissionInterval string `hcl:"circonus_submission_interval"`
	// CirconusCheckSubmissionURL is the check.config.submission_url field from a
	// previously created HTTPTRAP check.
	// Default: none
	CirconusCheckSubmissionURL string `hcl:"circonus_submission_url"`
	// CirconusCheckID is the check id (not check bundle id) from a previously created
	// HTTPTRAP check. The numeric portion of the check._cid field.
	// Default: none
	CirconusCheckID string `hcl:"circonus_check_id"`
	// CirconusCheckForceMetricActivation will force enabling metrics, as they are encountered,
	// if the metric already exists and is NOT active. If check management is enabled, the default
	// behavior is to add new metrics as they are encoutered. If the metric already exists in the
	// check, it will *NOT* be activated. This setting overrides that behavior.
	// Default: "false"
	CirconusCheckForceMetricActivation string `hcl:"circonus_check_force_metric_activation"`
	// CirconusCheckInstanceID serves to uniquely identify the metrics comming from this "instance".
	// It can be used to maintain metric continuity with transient or ephemeral instances as
	// they move around within an infrastructure.
	// Default: hostname:app
	CirconusCheckInstanceID string `hcl:"circonus_check_instance_id"`
	// CirconusCheckSearchTag is a special tag which, when coupled with the instance id, helps to
	// narrow down the search results when neither a Submission URL or Check ID is provided.
	// Default: service:app (e.g. service:consul)
	CirconusCheckSearchTag string `hcl:"circonus_check_search_tag"`
	// CirconusCheckTags is a comma separated list of tags to apply to the check. Note that
	// the value of CirconusCheckSearchTag will always be added to the check.
	// Default: none
	CirconusCheckTags string `mapstructure:"circonus_check_tags"`
	// CirconusCheckDisplayName is the name for the check which will be displayed in the Circonus UI.
	// Default: value of CirconusCheckInstanceID
	CirconusCheckDisplayName string `mapstructure:"circonus_check_display_name"`
	// CirconusBrokerID is an explicit broker to use when creating a new check. The numeric portion
	// of broker._cid. If metric management is enabled and neither a Submission URL nor Check ID
	// is provided, an attempt will be made to search for an existing check using Instance ID and
	// Search Tag. If one is not found, a new HTTPTRAP check will be created.
	// Default: use Select Tag if provided, otherwise, a random Enterprise Broker associated
	// with the specified API token or the default Circonus Broker.
	// Default: none
	CirconusBrokerID string `hcl:"circonus_broker_id"`
	// CirconusBrokerSelectTag is a special tag which will be used to select a broker when
	// a Broker ID is not provided. The best use of this is to as a hint for which broker
	// should be used based on *where* this particular instance is running.
	// (e.g. a specific geo location or datacenter, dc:sfo)
	// Default: none
	CirconusBrokerSelectTag string `hcl:"circonus_broker_select_tag"`

	// Dogstats:
	// DogStatsdAddr is the address of a dogstatsd instance. If provided,
	// metrics will be sent to that instance
	DogStatsDAddr string `hcl:"dogstatsd_addr"`

	// DogStatsdTags are the global tags that should be sent with each packet to dogstatsd
	// It is a list of strings, where each string looks like "my_tag_name:my_tag_value"
	DogStatsDTags []string `hcl:"dogstatsd_tags"`
}

func (s *Telemetry) GoString() string {
	return fmt.Sprintf("*%#v", *s)
}

// Merge merges two configurations.
func (c *Config) Merge(c2 *Config) *Config {
	if c2 == nil {
		return c
	}

	result := new(Config)
	for _, l := range c.Listeners {
		result.Listeners = append(result.Listeners, l)
	}
	for _, l := range c2.Listeners {
		result.Listeners = append(result.Listeners, l)
	}

	result.Storage = c.Storage
	if c2.Storage != nil {
		result.Storage = c2.Storage
	}

	result.HAStorage = c.HAStorage
	if c2.HAStorage != nil {
		result.HAStorage = c2.HAStorage
	}

	result.Seal = c.Seal
	if c2.Seal != nil {
		result.Seal = c2.Seal
	}

	result.Telemetry = c.Telemetry
	if c2.Telemetry != nil {
		result.Telemetry = c2.Telemetry
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

	result.DisableMlock = c.DisableMlock
	if c2.DisableMlock {
		result.DisableMlock = c2.DisableMlock
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

	result.ClusterName = c.ClusterName
	if c2.ClusterName != "" {
		result.ClusterName = c2.ClusterName
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

	result.PluginDirectory = c.PluginDirectory
	if c2.PluginDirectory != "" {
		result.PluginDirectory = c2.PluginDirectory
	}

	result.PidFile = c.PidFile
	if c2.PidFile != "" {
		result.PidFile = c2.PidFile
	}

	return result
}

// LoadConfig loads the configuration at the given path, regardless if
// its a file or directory.
func LoadConfig(path string, logger log.Logger) (*Config, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return LoadConfigDir(path, logger)
	}
	return LoadConfigFile(path, logger)
}

// LoadConfigFile loads the configuration from the given file.
func LoadConfigFile(path string, logger log.Logger) (*Config, error) {
	// Read the file
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(string(d), logger)
}

func ParseConfig(d string, logger log.Logger) (*Config, error) {
	// Parse!
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, err
	}

	// Start building the result
	var result Config
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

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

	if result.DisableMlockRaw != nil {
		if result.DisableMlock, err = parseutil.ParseBool(result.DisableMlockRaw); err != nil {
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

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	valid := []string{
		"storage",
		"ha_storage",
		"backend",
		"ha_backend",
		"hsm",
		"seal",
		"listener",
		"cache_size",
		"disable_cache",
		"disable_mlock",
		"ui",
		"telemetry",
		"default_lease_ttl",
		"max_lease_ttl",
		"cluster_name",
		"cluster_cipher_suites",
		"plugin_directory",
		"pid_file",
		"raw_storage_endpoint",
		"api_addr",
		"cluster_addr",
		"disable_clustering",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	// Look for storage but still support old backend
	if o := list.Filter("storage"); len(o.Items) > 0 {
		if err := parseStorage(&result, o, "storage"); err != nil {
			return nil, fmt.Errorf("error parsing 'storage': %s", err)
		}
	} else {
		if o := list.Filter("backend"); len(o.Items) > 0 {
			if err := parseStorage(&result, o, "backend"); err != nil {
				return nil, fmt.Errorf("error parsing 'backend': %s", err)
			}
		}
	}

	if o := list.Filter("ha_storage"); len(o.Items) > 0 {
		if err := parseHAStorage(&result, o, "ha_storage"); err != nil {
			return nil, fmt.Errorf("error parsing 'ha_storage': %s", err)
		}
	} else {
		if o := list.Filter("ha_backend"); len(o.Items) > 0 {
			if err := parseHAStorage(&result, o, "ha_backend"); err != nil {
				return nil, fmt.Errorf("error parsing 'ha_backend': %s", err)
			}
		}
	}

	if o := list.Filter("hsm"); len(o.Items) > 0 {
		if err := parseSeal(&result, o, "hsm"); err != nil {
			return nil, fmt.Errorf("error parsing 'hsm': %s", err)
		}
	}

	if o := list.Filter("seal"); len(o.Items) > 0 {
		if err := parseSeal(&result, o, "seal"); err != nil {
			return nil, fmt.Errorf("error parsing 'seal': %s", err)
		}
	}

	if o := list.Filter("listener"); len(o.Items) > 0 {
		if err := parseListeners(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'listener': %s", err)
		}
	}

	if o := list.Filter("telemetry"); len(o.Items) > 0 {
		if err := parseTelemetry(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'telemetry': %s", err)
		}
	}

	return &result, nil
}

// LoadConfigDir loads all the configurations in the given directory
// in alphabetical order.
func LoadConfigDir(dir string, logger log.Logger) (*Config, error) {
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
		return nil, fmt.Errorf(
			"configuration path must be a directory: %s",
			dir)
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

	var result *Config
	for _, f := range files {
		config, err := LoadConfigFile(f, logger)
		if err != nil {
			return nil, fmt.Errorf("Error loading %s: %s", f, err)
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

func parseStorage(result *Config, list *ast.ObjectList, name string) error {
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

func parseSeal(result *Config, list *ast.ObjectList, blockName string) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q block is permitted", blockName)
	}

	// Get our item
	item := list.Items[0]

	key := blockName
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	var valid []string
	// Valid parameter for the Seal types
	switch key {
	case "pkcs11":
		valid = []string{
			"lib",
			"slot",
			"pin",
			"mechanism",
			"hmac_mechanism",
			"key_label",
			"hmac_key_label",
			"generate_key",
			"regenerate_key",
			"max_parallel",
		}
	case "awskms":
		valid = []string{
			"region",
			"access_key",
			"secret_key",
			"kms_key_id",
			"max_parallel",
		}
	default:
		return fmt.Errorf("invalid seal type %q", key)
	}

	if err := checkHCLKeys(item.Val, valid); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
	}

	var m map[string]string
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
	}

	result.Seal = &Seal{
		Type:   strings.ToLower(key),
		Config: m,
	}

	return nil
}

func parseListeners(result *Config, list *ast.ObjectList) error {
	listeners := make([]*Listener, 0, len(list.Items))
	for _, item := range list.Items {
		key := "listener"
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		valid := []string{
			"address",
			"cluster_address",
			"endpoint",
			"infrastructure",
			"node_id",
			"proxy_protocol_behavior",
			"proxy_protocol_authorized_addrs",
			"tls_disable",
			"tls_cert_file",
			"tls_key_file",
			"tls_min_version",
			"tls_cipher_suites",
			"tls_prefer_server_cipher_suites",
			"tls_require_and_verify_client_cert",
			"tls_disable_client_certs",
			"tls_client_ca_file",
			"token",
		}
		if err := checkHCLKeys(item.Val, valid); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%s:", key))
		}

		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%s:", key))
		}

		lnType := strings.ToLower(key)

		listeners = append(listeners, &Listener{
			Type:   lnType,
			Config: m,
		})
	}

	result.Listeners = listeners
	return nil
}

func parseTelemetry(result *Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'telemetry' block is permitted")
	}

	// Get our one item
	item := list.Items[0]

	// Check for invalid keys
	valid := []string{
		"circonus_api_token",
		"circonus_api_app",
		"circonus_api_url",
		"circonus_submission_interval",
		"circonus_submission_url",
		"circonus_check_id",
		"circonus_check_force_metric_activation",
		"circonus_check_instance_id",
		"circonus_check_search_tag",
		"circonus_check_display_name",
		"circonus_check_tags",
		"circonus_broker_id",
		"circonus_broker_select_tag",
		"disable_hostname",
		"dogstatsd_addr",
		"dogstatsd_tags",
		"statsd_address",
		"statsite_address",
	}
	if err := checkHCLKeys(item.Val, valid); err != nil {
		return multierror.Prefix(err, "telemetry:")
	}

	var t Telemetry
	if err := hcl.DecodeObject(&t, item.Val); err != nil {
		return multierror.Prefix(err, "telemetry:")
	}

	if result.Telemetry == nil {
		result.Telemetry = &Telemetry{}
	}

	if err := hcl.DecodeObject(&result.Telemetry, item.Val); err != nil {
		return multierror.Prefix(err, "telemetry:")
	}
	return nil
}

func checkHCLKeys(node ast.Node, valid []string) error {
	var list *ast.ObjectList
	switch n := node.(type) {
	case *ast.ObjectList:
		list = n
	case *ast.ObjectType:
		list = n.List
	default:
		return fmt.Errorf("cannot check HCL keys of type %T", n)
	}

	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := validMap[key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invalid key '%s' on line %d", key, item.Assign.Line))
		}
	}

	return result
}
