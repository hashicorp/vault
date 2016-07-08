package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

// ReloadFunc are functions that are called when a reload is requested.
type ReloadFunc func(map[string]string) error

// Config is the configuration for the vault server.
type Config struct {
	Listeners []*Listener `hcl:"-"`
	Backend   *Backend    `hcl:"-"`
	HABackend *Backend    `hcl:"-"`

	DisableCache bool `hcl:"disable_cache"`
	DisableMlock bool `hcl:"disable_mlock"`

	Telemetry *Telemetry `hcl:"telemetry"`

	MaxLeaseTTL        time.Duration `hcl:"-"`
	MaxLeaseTTLRaw     string        `hcl:"max_lease_ttl"`
	DefaultLeaseTTL    time.Duration `hcl:"-"`
	DefaultLeaseTTLRaw string        `hcl:"default_lease_ttl"`
}

// DevConfig is a Config that is used for dev mode of Vault.
func DevConfig() *Config {
	return &Config{
		DisableCache: false,
		DisableMlock: true,

		Backend: &Backend{
			Type: "inmem",
		},

		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]string{
					"address":     "127.0.0.1:8200",
					"tls_disable": "1",
				},
			},
		},

		Telemetry: &Telemetry{},

		MaxLeaseTTL:     30 * 24 * time.Hour,
		DefaultLeaseTTL: 30 * 24 * time.Hour,
	}
}

// Listener is the listener configuration for the server.
type Listener struct {
	Type   string
	Config map[string]string
}

func (l *Listener) GoString() string {
	return fmt.Sprintf("*%#v", *l)
}

// Backend is the backend configuration for the server.
type Backend struct {
	Type          string
	AdvertiseAddr string
	Config        map[string]string
}

func (b *Backend) GoString() string {
	return fmt.Sprintf("*%#v", *b)
}

// Telemetry is the telemetry configuration for the server
type Telemetry struct {
	StatsiteAddr string `hcl:"statsite_address"`
	StatsdAddr   string `hcl:"statsd_address"`

	DisableHostname bool `hcl:"disable_hostname"`
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

	result.Backend = c.Backend
	if c2.Backend != nil {
		result.Backend = c2.Backend
	}

	result.HABackend = c.HABackend
	if c2.HABackend != nil {
		result.HABackend = c2.HABackend
	}

	result.Telemetry = c.Telemetry
	if c2.Telemetry != nil {
		result.Telemetry = c2.Telemetry
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
		return LoadConfigDir(path)
	} else {
		return LoadConfigFile(path)
	}
}

// LoadConfigFile loads the configuration from the given file.
func LoadConfigFile(path string) (*Config, error) {
	// Read the file
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(string(d))
}

func ParseConfig(d string) (*Config, error) {
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

	if result.MaxLeaseTTLRaw != "" {
		if result.MaxLeaseTTL, err = time.ParseDuration(result.MaxLeaseTTLRaw); err != nil {
			return nil, err
		}
	}
	if result.DefaultLeaseTTLRaw != "" {
		if result.DefaultLeaseTTL, err = time.ParseDuration(result.DefaultLeaseTTLRaw); err != nil {
			return nil, err
		}
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	valid := []string{
		"atlas",
		"backend",
		"ha_backend",
		"listener",
		"disable_cache",
		"disable_mlock",
		"telemetry",
		"default_lease_ttl",
		"max_lease_ttl",

		// TODO: Remove in 0.6.0
		// Deprecated keys
		"statsd_addr",
		"statsite_addr",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	// TODO: Remove in 0.6.0
	// Preflight checks for deprecated keys
	sda := list.Filter("statsd_addr")
	ssa := list.Filter("statsite_addr")
	if len(sda.Items) > 0 || len(ssa.Items) > 0 {
		log.Println("[WARN] The top-level keys 'statsd_addr' and 'statsite_addr' " +
			"have been moved into a 'telemetry' block instead. Please update your " +
			"Vault configuration as this deprecation will be removed in the next " +
			"major release. Values specified in a 'telemetry' block will take " +
			"precendence.")

		t := struct {
			StatsdAddr   string `hcl:"statsd_addr"`
			StatsiteAddr string `hcl:"statsite_addr"`
		}{}
		if err := hcl.DecodeObject(&t, list); err != nil {
			return nil, err
		}

		result.Telemetry = &Telemetry{
			StatsdAddr:   t.StatsdAddr,
			StatsiteAddr: t.StatsiteAddr,
		}
	}

	if o := list.Filter("backend"); len(o.Items) > 0 {
		if err := parseBackends(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'backend': %s", err)
		}
	}

	if o := list.Filter("ha_backend"); len(o.Items) > 0 {
		if err := parseHABackends(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'ha_backend': %s", err)
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
		config, err := LoadConfigFile(f)
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

func parseBackends(result *Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'backend' block is permitted")
	}

	// Get our item
	item := list.Items[0]

	key := "backend"
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	var m map[string]string
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("backend.%s:", key))
	}

	// Pull out the advertise address since it's common to all backends
	var advertiseAddr string
	if v, ok := m["advertise_addr"]; ok {
		advertiseAddr = v
		delete(m, "advertise_addr")
	}

	result.Backend = &Backend{
		AdvertiseAddr: advertiseAddr,
		Type:          strings.ToLower(key),
		Config:        m,
	}
	return nil
}

func parseHABackends(result *Config, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'ha_backend' block is permitted")
	}

	// Get our item
	item := list.Items[0]

	key := "backend"
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	var m map[string]string
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return multierror.Prefix(err, fmt.Sprintf("ha_backend.%s:", key))
	}

	// Pull out the advertise address since it's common to all backends
	var advertiseAddr string
	if v, ok := m["advertise_addr"]; ok {
		advertiseAddr = v
		delete(m, "advertise_addr")
	}

	result.HABackend = &Backend{
		AdvertiseAddr: advertiseAddr,
		Type:          strings.ToLower(key),
		Config:        m,
	}
	return nil
}

func parseListeners(result *Config, list *ast.ObjectList) error {
	var foundAtlas bool

	listeners := make([]*Listener, 0, len(list.Items))
	for _, item := range list.Items {
		key := "listener"
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		valid := []string{
			"address",
			"endpoint",
			"infrastructure",
			"node_id",
			"tls_disable",
			"tls_cert_file",
			"tls_key_file",
			"tls_min_version",
			"token",
		}
		if err := checkHCLKeys(item.Val, valid); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%s:", key))
		}

		var m map[string]string
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%s:", key))
		}

		lnType := strings.ToLower(key)

		if lnType == "atlas" {
			if foundAtlas {
				return multierror.Prefix(fmt.Errorf("only one listener of type 'atlas' is permitted"), fmt.Sprintf("listeners.%s", key))
			}

			foundAtlas = true
			if m["token"] == "" {
				return multierror.Prefix(fmt.Errorf("'token' must be specified for an Atlas listener"), fmt.Sprintf("listeners.%s", key))
			}
			if m["infrastructure"] == "" {
				return multierror.Prefix(fmt.Errorf("'infrastructure' must be specified for an Atlas listener"), fmt.Sprintf("listeners.%s", key))
			}
			if m["node_id"] == "" {
				return multierror.Prefix(fmt.Errorf("'node_id' must be specified for an Atlas listener"), fmt.Sprintf("listeners.%s", key))
			}
		}

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
		"statsite_address",
		"statsd_address",
		"disable_hostname",
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
