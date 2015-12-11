package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl"
	hclobj "github.com/hashicorp/hcl/hcl"
)

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

	// Parse!
	obj, err := hcl.Parse(string(d))
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

	if objs := obj.Get("listener", false); objs != nil {
		result.Listeners, err = loadListeners(objs)
		if err != nil {
			return nil, err
		}
	}
	if objs := obj.Get("backend", false); objs != nil {
		result.Backend, err = loadBackend(objs)
		if err != nil {
			return nil, err
		}
	}
	if objs := obj.Get("ha_backend", false); objs != nil {
		result.HABackend, err = loadBackend(objs)
		if err != nil {
			return nil, err
		}
	}

	// A little hacky but upgrades the old stats config directives to the new way
	if result.Telemetry == nil {
		statsdAddr := obj.Get("statsd_addr", false)
		statsiteAddr := obj.Get("statsite_addr", false)

		if statsdAddr != nil || statsiteAddr != nil {
			result.Telemetry = &Telemetry{
				StatsdAddr:   getString(statsdAddr),
				StatsiteAddr: getString(statsiteAddr),
			}
		}
	}

	return &result, nil
}

func getString(o *hclobj.Object) string {
	if o == nil || o.Type != hclobj.ValueTypeString {
		return ""
	}

	return o.Value.(string)
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

func loadListeners(os *hclobj.Object) ([]*Listener, error) {
	var allNames []*hclobj.Object

	// Really confusing iteration. The key is the false/true parameter
	// of whether we're expanding or not. We first iterate over all
	// the "listeners"
	for _, o1 := range os.Elem(false) {
		// Iterate expand to get the list of types
		for _, o2 := range o1.Elem(true) {
			switch o2.Type {
			case hclobj.ValueTypeList:
				// This switch is for JSON, to allow them to do this:
				//
				// "tcp": [{ ... }, { ... }]
				//
				// To configure multiple listeners of the same type.
				for _, o3 := range o2.Elem(true) {
					o3.Key = o2.Key
					allNames = append(allNames, o3)
				}
			case hclobj.ValueTypeObject:
				// This is for the standard `listener "tcp" { ... }` syntax
				allNames = append(allNames, o2)
			}
		}
	}

	if len(allNames) == 0 {
		return nil, nil
	}

	// Now go over all the types and their children in order to get
	// all of the actual resources.
	result := make([]*Listener, 0, len(allNames))
	for _, obj := range allNames {
		k := obj.Key

		var config map[string]string
		if err := hcl.DecodeObject(&config, obj); err != nil {
			return nil, fmt.Errorf(
				"Error reading config for %s: %s",
				k,
				err)
		}

		result = append(result, &Listener{
			Type:   k,
			Config: config,
		})
	}

	return result, nil
}

func loadBackend(os *hclobj.Object) (*Backend, error) {
	var allNames []*hclobj.Object

	// See loadListeners
	for _, o1 := range os.Elem(false) {
		// Iterate expand to get the list of types
		for _, o2 := range o1.Elem(true) {
			// Iterate non-expand to get the full list of types
			for _, o3 := range o2.Elem(false) {
				allNames = append(allNames, o3)
			}
		}
	}

	if len(allNames) == 0 {
		return nil, nil
	}
	if len(allNames) > 1 {
		keys := make([]string, 0, len(allNames))
		for _, o := range allNames {
			keys = append(keys, o.Key)
		}

		return nil, fmt.Errorf(
			"Multiple backends declared. Only one is allowed: %v", keys)
	}

	// Now go over all the types and their children in order to get
	// all of the actual resources.
	var result Backend
	obj := allNames[0]
	result.Type = obj.Key

	var config map[string]string
	if err := hcl.DecodeObject(&config, obj); err != nil {
		return nil, fmt.Errorf(
			"Error reading config for backend %s: %s",
			result.Type,
			err)
	}

	if v, ok := config["advertise_addr"]; ok {
		result.AdvertiseAddr = v
		delete(config, "advertise_addr")
	}

	result.Config = config
	return &result, nil
}
