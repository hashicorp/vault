package configutil

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

// SharedConfig contains some shared values
type SharedConfig struct {
	FoundKeys  []string     `hcl:",decodedFields"`
	UnusedKeys UnusedKeyMap `hcl:",unusedKeyPositions"`
	Sections   map[string][]token.Pos

	EntSharedConfig

	Listeners []*Listener `hcl:"-"`

	Seals   []*KMS   `hcl:"-"`
	Entropy *Entropy `hcl:"-"`

	DisableMlock    bool        `hcl:"-"`
	DisableMlockRaw interface{} `hcl:"disable_mlock"`

	Telemetry *Telemetry `hcl:"telemetry"`

	DefaultMaxRequestDuration    time.Duration `hcl:"-"`
	DefaultMaxRequestDurationRaw interface{}   `hcl:"default_max_request_duration"`

	// LogFormat specifies the log format. Valid values are "standard" and
	// "json". The values are case-insenstive. If no log format is specified,
	// then standard format will be used.
	LogFormat string `hcl:"log_format"`
	LogLevel  string `hcl:"log_level"`

	PidFile string `hcl:"pid_file"`

	ClusterName string `hcl:"cluster_name"`
}

// LoadConfigFile loads the configuration from the given file.
func LoadConfigFile(path string) (*SharedConfig, error) {
	// Read the file
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(string(d))
}

func LoadConfigKMSes(path string) ([]*KMS, error) {
	// Read the file
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseKMSes(string(d))
}

func ParseConfig(d string) (*SharedConfig, error) {
	// Parse!
	obj, err := hcl.Parse(d)
	if err != nil {
		return nil, err
	}

	// Start building the result
	var result SharedConfig

	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	if result.DefaultMaxRequestDurationRaw != nil {
		if result.DefaultMaxRequestDuration, err = parseutil.ParseDurationSecond(result.DefaultMaxRequestDurationRaw); err != nil {
			return nil, err
		}
		result.FoundKeys = append(result.FoundKeys, "DefaultMaxRequestDuration")
		result.DefaultMaxRequestDurationRaw = nil
	}

	if result.DisableMlockRaw != nil {
		if result.DisableMlock, err = parseutil.ParseBool(result.DisableMlockRaw); err != nil {
			return nil, err
		}
		result.FoundKeys = append(result.FoundKeys, "DisableMlock")
		result.DisableMlockRaw = nil
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	if o := list.Filter("hsm"); len(o.Items) > 0 {
		result.found("hsm", "hsm")
		if err := parseKMS(&result.Seals, o, "hsm", 2); err != nil {
			return nil, fmt.Errorf("error parsing 'hsm': %w", err)
		}
	}

	if o := list.Filter("seal"); len(o.Items) > 0 {
		result.found("seal", "Seal")
		if err := parseKMS(&result.Seals, o, "seal", 3); err != nil {
			return nil, fmt.Errorf("error parsing 'seal': %w", err)
		}
	}

	if o := list.Filter("kms"); len(o.Items) > 0 {
		result.found("kms", "Seal")
		if err := parseKMS(&result.Seals, o, "kms", 3); err != nil {
			return nil, fmt.Errorf("error parsing 'kms': %w", err)
		}
	}

	if o := list.Filter("entropy"); len(o.Items) > 0 {
		result.found("entropy", "Entropy")
		if err := ParseEntropy(&result, o, "entropy"); err != nil {
			return nil, fmt.Errorf("error parsing 'entropy': %w", err)
		}
	}

	if o := list.Filter("listener"); len(o.Items) > 0 {
		result.found("listener", "Listener")
		if err := ParseListeners(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'listener': %w", err)
		}
	}

	if o := list.Filter("telemetry"); len(o.Items) > 0 {
		result.found("telemetry", "Telemetry")
		if err := parseTelemetry(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'telemetry': %w", err)
		}
	}

	entConfig := &(result.EntSharedConfig)
	if err := entConfig.ParseConfig(list); err != nil {
		return nil, fmt.Errorf("error parsing enterprise config: %w", err)
	}

	return &result, nil
}

// Sanitized returns a copy of the config with all values that are considered
// sensitive stripped. It also strips all `*Raw` values that are mainly
// used for parsing.
//
// Specifically, the fields that this method strips are:
// - KMS.Config
// - Telemetry.CirconusAPIToken
func (c *SharedConfig) Sanitized() map[string]interface{} {
	if c == nil {
		return nil
	}

	result := map[string]interface{}{
		"disable_mlock": c.DisableMlock,

		"default_max_request_duration": c.DefaultMaxRequestDuration,

		"log_level":  c.LogLevel,
		"log_format": c.LogFormat,

		"pid_file": c.PidFile,

		"cluster_name": c.ClusterName,
	}

	// Sanitize listeners
	if len(c.Listeners) != 0 {
		var sanitizedListeners []interface{}
		for _, ln := range c.Listeners {
			cleanLn := map[string]interface{}{
				"type":   ln.Type,
				"config": ln.RawConfig,
			}
			sanitizedListeners = append(sanitizedListeners, cleanLn)
		}
		result["listeners"] = sanitizedListeners
	}

	// Sanitize seals stanza
	if len(c.Seals) != 0 {
		var sanitizedSeals []interface{}
		for _, s := range c.Seals {
			cleanSeal := map[string]interface{}{
				"type":     s.Type,
				"disabled": s.Disabled,
			}
			sanitizedSeals = append(sanitizedSeals, cleanSeal)
		}
		result["seals"] = sanitizedSeals
	}

	// Sanitize telemetry stanza
	if c.Telemetry != nil {
		sanitizedTelemetry := map[string]interface{}{
			"statsite_address":                       c.Telemetry.StatsiteAddr,
			"statsd_address":                         c.Telemetry.StatsdAddr,
			"disable_hostname":                       c.Telemetry.DisableHostname,
			"metrics_prefix":                         c.Telemetry.MetricsPrefix,
			"usage_gauge_period":                     c.Telemetry.UsageGaugePeriod,
			"maximum_gauge_cardinality":              c.Telemetry.MaximumGaugeCardinality,
			"circonus_api_token":                     "",
			"circonus_api_app":                       c.Telemetry.CirconusAPIApp,
			"circonus_api_url":                       c.Telemetry.CirconusAPIURL,
			"circonus_submission_interval":           c.Telemetry.CirconusSubmissionInterval,
			"circonus_submission_url":                c.Telemetry.CirconusCheckSubmissionURL,
			"circonus_check_id":                      c.Telemetry.CirconusCheckID,
			"circonus_check_force_metric_activation": c.Telemetry.CirconusCheckForceMetricActivation,
			"circonus_check_instance_id":             c.Telemetry.CirconusCheckInstanceID,
			"circonus_check_search_tag":              c.Telemetry.CirconusCheckSearchTag,
			"circonus_check_tags":                    c.Telemetry.CirconusCheckTags,
			"circonus_check_display_name":            c.Telemetry.CirconusCheckDisplayName,
			"circonus_broker_id":                     c.Telemetry.CirconusBrokerID,
			"circonus_broker_select_tag":             c.Telemetry.CirconusBrokerSelectTag,
			"dogstatsd_addr":                         c.Telemetry.DogStatsDAddr,
			"dogstatsd_tags":                         c.Telemetry.DogStatsDTags,
			"prometheus_retention_time":              c.Telemetry.PrometheusRetentionTime,
			"stackdriver_project_id":                 c.Telemetry.StackdriverProjectID,
			"stackdriver_location":                   c.Telemetry.StackdriverLocation,
			"stackdriver_namespace":                  c.Telemetry.StackdriverNamespace,
			"stackdriver_debug_logs":                 c.Telemetry.StackdriverDebugLogs,
			"lease_metrics_epsilon":                  c.Telemetry.LeaseMetricsEpsilon,
			"num_lease_metrics_buckets":              c.Telemetry.NumLeaseMetricsTimeBuckets,
			"add_lease_metrics_namespace_labels":     c.Telemetry.LeaseMetricsNameSpaceLabels,
		}
		result["telemetry"] = sanitizedTelemetry
	}

	return result
}

func (c *SharedConfig) found(s, k string) {
	delete(c.UnusedKeys, s)
	c.FoundKeys = append(c.FoundKeys, k)
}
