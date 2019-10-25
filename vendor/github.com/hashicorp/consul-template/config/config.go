package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/consul-template/signals"
	"github.com/hashicorp/hcl"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"

	"github.com/pkg/errors"
)

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = "WARN"

	// DefaultMaxStale is the default staleness permitted. This enables stale
	// queries by default for performance reasons.
	DefaultMaxStale = 2 * time.Second

	// DefaultReloadSignal is the default signal for reload.
	DefaultReloadSignal = syscall.SIGHUP

	// DefaultKillSignal is the default signal for termination.
	DefaultKillSignal = syscall.SIGINT
)

var (
	// homePath is the location to the user's home directory.
	homePath, _ = homedir.Dir()
)

// Config is used to configure Consul Template
type Config struct {
	// Consul is the configuration for connecting to a Consul cluster.
	Consul *ConsulConfig `mapstructure:"consul"`

	// Dedup is used to configure the dedup settings
	Dedup *DedupConfig `mapstructure:"deduplicate"`

	// Exec is the configuration for exec/supervise mode.
	Exec *ExecConfig `mapstructure:"exec"`

	// KillSignal is the signal to listen for a graceful terminate event.
	KillSignal *os.Signal `mapstructure:"kill_signal"`

	// LogLevel is the level with which to log for this config.
	LogLevel *string `mapstructure:"log_level"`

	// MaxStale is the maximum amount of time for staleness from Consul as given
	// by LastContact. If supplied, Consul Template will query all servers instead
	// of just the leader.
	MaxStale *time.Duration `mapstructure:"max_stale"`

	// PidFile is the path on disk where a PID file should be written containing
	// this processes PID.
	PidFile *string `mapstructure:"pid_file"`

	// ReloadSignal is the signal to listen for a reload event.
	ReloadSignal *os.Signal `mapstructure:"reload_signal"`

	// Syslog is the configuration for syslog.
	Syslog *SyslogConfig `mapstructure:"syslog"`

	// Templates is the list of templates.
	Templates *TemplateConfigs `mapstructure:"template"`

	// Vault is the configuration for connecting to a vault server.
	Vault *VaultConfig `mapstructure:"vault"`

	// Wait is the quiescence timers.
	Wait *WaitConfig `mapstructure:"wait"`

	// Additional command line options
	// Run once, executing each template exactly once, and exit
	Once bool
}

// Copy returns a deep copy of the current configuration. This is useful because
// the nested data structures may be shared.
func (c *Config) Copy() *Config {
	if c == nil {
		return nil
	}
	var o Config

	o.Consul = c.Consul

	if c.Consul != nil {
		o.Consul = c.Consul.Copy()
	}

	if c.Dedup != nil {
		o.Dedup = c.Dedup.Copy()
	}

	if c.Exec != nil {
		o.Exec = c.Exec.Copy()
	}

	o.KillSignal = c.KillSignal

	o.LogLevel = c.LogLevel

	o.MaxStale = c.MaxStale

	o.PidFile = c.PidFile

	o.ReloadSignal = c.ReloadSignal

	if c.Syslog != nil {
		o.Syslog = c.Syslog.Copy()
	}

	if c.Templates != nil {
		o.Templates = c.Templates.Copy()
	}

	if c.Vault != nil {
		o.Vault = c.Vault.Copy()
	}

	if c.Wait != nil {
		o.Wait = c.Wait.Copy()
	}

	o.Once = c.Once

	return &o
}

// Merge merges the values in config into this config object. Values in the
// config object overwrite the values in c.
func (c *Config) Merge(o *Config) *Config {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Consul != nil {
		r.Consul = r.Consul.Merge(o.Consul)
	}

	if o.Dedup != nil {
		r.Dedup = r.Dedup.Merge(o.Dedup)
	}

	if o.Exec != nil {
		r.Exec = r.Exec.Merge(o.Exec)
	}

	if o.KillSignal != nil {
		r.KillSignal = o.KillSignal
	}

	if o.LogLevel != nil {
		r.LogLevel = o.LogLevel
	}

	if o.MaxStale != nil {
		r.MaxStale = o.MaxStale
	}

	if o.PidFile != nil {
		r.PidFile = o.PidFile
	}

	if o.ReloadSignal != nil {
		r.ReloadSignal = o.ReloadSignal
	}

	if o.Syslog != nil {
		r.Syslog = r.Syslog.Merge(o.Syslog)
	}

	if o.Templates != nil {
		r.Templates = r.Templates.Merge(o.Templates)
	}

	if o.Vault != nil {
		r.Vault = r.Vault.Merge(o.Vault)
	}

	if o.Wait != nil {
		r.Wait = r.Wait.Merge(o.Wait)
	}

	r.Once = o.Once

	return r
}

// Parse parses the given string contents as a config
func Parse(s string) (*Config, error) {
	var shadow interface{}
	if err := hcl.Decode(&shadow, s); err != nil {
		return nil, errors.Wrap(err, "error decoding config")
	}

	// Convert to a map and flatten the keys we want to flatten
	parsed, ok := shadow.(map[string]interface{})
	if !ok {
		return nil, errors.New("error converting config")
	}

	flattenKeys(parsed, []string{
		"auth",
		"consul",
		"consul.auth",
		"consul.retry",
		"consul.ssl",
		"consul.transport",
		"deduplicate",
		"env",
		"exec",
		"exec.env",
		"ssl",
		"syslog",
		"vault",
		"vault.retry",
		"vault.ssl",
		"vault.transport",
		"wait",
	})

	// FlattenFlatten keys belonging to the templates. We cannot do this above
	// because it is an array of templates.
	if templates, ok := parsed["template"].([]map[string]interface{}); ok {
		for _, template := range templates {
			flattenKeys(template, []string{
				"env",
				"exec",
				"exec.env",
				"wait",
			})
		}
	}

	// Create a new, empty config
	var c Config

	// Use mapstructure to populate the basic config fields
	var md mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ConsulStringToStructFunc(),
			StringToFileModeFunc(),
			signals.StringToSignalFunc(),
			StringToWaitDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
		),
		ErrorUnused: true,
		Metadata:    &md,
		Result:      &c,
	})
	if err != nil {
		return nil, errors.Wrap(err, "mapstructure decoder creation failed")
	}
	if err := decoder.Decode(parsed); err != nil {
		return nil, errors.Wrap(err, "mapstructure decode failed")
	}

	return &c, nil
}

// Must returns a config object that must compile. If there are any errors, this
// function will panic. This is most useful in testing or constants.
func Must(s string) *Config {
	c, err := Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// TestConfig returns a default, finalized config, with the provided
// configuration taking precedence.
func TestConfig(c *Config) *Config {
	d := DefaultConfig().Merge(c)
	d.Finalize()
	return d
}

// FromFile reads the configuration file at the given path and returns a new
// Config struct with the data populated.
func FromFile(path string) (*Config, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "from file: "+path)
	}

	config, err := Parse(string(c))
	if err != nil {
		return nil, errors.Wrap(err, "from file: "+path)
	}
	return config, nil
}

// FromPath iterates and merges all configuration files in a given
// directory, returning the resulting config.
func FromPath(path string) (*Config, error) {
	// Ensure the given filepath exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "missing file/folder: "+path)
	}

	// Check if a file was given or a path to a directory
	stat, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed stating file: "+path)
	}

	// Recursively parse directories, single load files
	if stat.Mode().IsDir() {
		// Ensure the given filepath has at least one config file
		_, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed listing dir: "+path)
		}

		// Create a blank config to merge off of
		var c *Config

		// Potential bug: Walk does not follow symlinks!
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			// If WalkFunc had an error, just return it
			if err != nil {
				return err
			}

			// Do nothing for directories
			if info.IsDir() {
				return nil
			}

			// Parse and merge the config
			newConfig, err := FromFile(path)
			if err != nil {
				return err
			}
			c = c.Merge(newConfig)

			return nil
		})

		if err != nil {
			return nil, errors.Wrap(err, "walk error")
		}

		return c, nil
	} else if stat.Mode().IsRegular() {
		return FromFile(path)
	}

	return nil, fmt.Errorf("unknown filetype: %q", stat.Mode().String())
}

// GoString defines the printable version of this struct.
func (c *Config) GoString() string {
	if c == nil {
		return "(*Config)(nil)"
	}

	return fmt.Sprintf("&Config{"+
		"Consul:%#v, "+
		"Dedup:%#v, "+
		"Exec:%#v, "+
		"KillSignal:%s, "+
		"LogLevel:%s, "+
		"MaxStale:%s, "+
		"PidFile:%s, "+
		"ReloadSignal:%s, "+
		"Syslog:%#v, "+
		"Templates:%#v, "+
		"Vault:%#v, "+
		"Wait:%#v,"+
		"Once:%#v"+
		"}",
		c.Consul,
		c.Dedup,
		c.Exec,
		SignalGoString(c.KillSignal),
		StringGoString(c.LogLevel),
		TimeDurationGoString(c.MaxStale),
		StringGoString(c.PidFile),
		SignalGoString(c.ReloadSignal),
		c.Syslog,
		c.Templates,
		c.Vault,
		c.Wait,
		c.Once,
	)
}

// Show diff between 2 Configs, useful in tests
func (expected *Config) Diff(actual *Config) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n")
	ve := reflect.ValueOf(*expected)
	va := reflect.ValueOf(*actual)
	ct := ve.Type()

	for i := 0; i < ve.NumField(); i++ {
		fc := ve.Field(i)
		fo := va.Field(i)
		if !reflect.DeepEqual(fc.Interface(), fo.Interface()) {
			fmt.Fprintf(&b, "%s:\n", ct.Field(i).Name)
			fmt.Fprintf(&b, "\texp: %#v\n", fc.Interface())
			fmt.Fprintf(&b, "\tact: %#v\n", fo.Interface())
		}
	}

	return b.String()
}

// DefaultConfig returns the default configuration struct. Certain environment
// variables may be set which control the values for the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Consul:    DefaultConsulConfig(),
		Dedup:     DefaultDedupConfig(),
		Exec:      DefaultExecConfig(),
		Syslog:    DefaultSyslogConfig(),
		Templates: DefaultTemplateConfigs(),
		Vault:     DefaultVaultConfig(),
		Wait:      DefaultWaitConfig(),
	}
}

// Finalize ensures all configuration options have the default values, so it
// is safe to dereference the pointers later down the line. It also
// intelligently tries to activate stanzas that should be "enabled" because
// data was given, but the user did not explicitly add "Enabled: true" to the
// configuration.
func (c *Config) Finalize() {
	if c == nil {
		return
	}
	if c.Consul == nil {
		c.Consul = DefaultConsulConfig()
	}
	c.Consul.Finalize()

	if c.Dedup == nil {
		c.Dedup = DefaultDedupConfig()
	}
	c.Dedup.Finalize()

	if c.Exec == nil {
		c.Exec = DefaultExecConfig()
	}
	c.Exec.Finalize()

	if c.KillSignal == nil {
		c.KillSignal = Signal(DefaultKillSignal)
	}

	if c.LogLevel == nil {
		c.LogLevel = stringFromEnv([]string{
			"CT_LOG",
			"CONSUL_TEMPLATE_LOG",
		}, DefaultLogLevel)
	}

	if c.MaxStale == nil {
		c.MaxStale = TimeDuration(DefaultMaxStale)
	}

	if c.PidFile == nil {
		c.PidFile = String("")
	}

	if c.ReloadSignal == nil {
		c.ReloadSignal = Signal(DefaultReloadSignal)
	}

	if c.Syslog == nil {
		c.Syslog = DefaultSyslogConfig()
	}
	c.Syslog.Finalize()

	if c.Templates == nil {
		c.Templates = DefaultTemplateConfigs()
	}
	c.Templates.Finalize()

	if c.Vault == nil {
		c.Vault = DefaultVaultConfig()
	}
	c.Vault.Finalize()

	if c.Wait == nil {
		c.Wait = DefaultWaitConfig()
	}
	c.Wait.Finalize()

	// disable Wait if -once was specified
	if c.Once {
		c.Wait = &WaitConfig{Enabled: Bool(false)}
	}
}

func stringFromEnv(list []string, def string) *string {
	for _, s := range list {
		if v := os.Getenv(s); v != "" {
			return String(strings.TrimSpace(v))
		}
	}
	return String(def)
}

func stringFromFile(list []string, def string) *string {
	for _, s := range list {
		c, err := ioutil.ReadFile(s)
		if err == nil {
			return String(strings.TrimSpace(string(c)))
		}
	}
	return String(def)
}

func antiboolFromEnv(list []string, def bool) *bool {
	for _, s := range list {
		if v := os.Getenv(s); v != "" {
			b, err := strconv.ParseBool(v)
			if err == nil {
				return Bool(!b)
			}
		}
	}
	return Bool(def)
}

func boolFromEnv(list []string, def bool) *bool {
	for _, s := range list {
		if v := os.Getenv(s); v != "" {
			b, err := strconv.ParseBool(v)
			if err == nil {
				return Bool(b)
			}
		}
	}
	return Bool(def)
}

// flattenKeys is a function that takes a map[string]interface{} and recursively
// flattens any keys that are a []map[string]interface{} where the key is in the
// given list of keys.
func flattenKeys(m map[string]interface{}, keys []string) {
	keyMap := make(map[string]struct{})
	for _, key := range keys {
		keyMap[key] = struct{}{}
	}

	var flatten func(map[string]interface{}, string)
	flatten = func(m map[string]interface{}, parent string) {
		for k, v := range m {
			// Calculate the map key, since it could include a parent.
			mapKey := k
			if parent != "" {
				mapKey = parent + "." + k
			}

			if _, ok := keyMap[mapKey]; !ok {
				continue
			}

			switch typed := v.(type) {
			case []map[string]interface{}:
				if len(typed) > 0 {
					last := typed[len(typed)-1]
					flatten(last, mapKey)
					m[k] = last
				} else {
					m[k] = nil
				}
			case map[string]interface{}:
				flatten(typed, mapKey)
				m[k] = typed
			default:
				m[k] = v
			}
		}
	}

	flatten(m, "")
}
