package configutil

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

// SharedConfig contains some shared values
type SharedConfig struct {
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
		result.DefaultMaxRequestDurationRaw = nil
	}

	if result.DisableMlockRaw != nil {
		if result.DisableMlock, err = parseutil.ParseBool(result.DisableMlockRaw); err != nil {
			return nil, err
		}
		result.DisableMlockRaw = nil
	}

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	if o := list.Filter("hsm"); len(o.Items) > 0 {
		if err := parseKMS(&result, o, "hsm", 1); err != nil {
			return nil, errwrap.Wrapf("error parsing 'hsm': {{err}}", err)
		}
	}

	if o := list.Filter("seal"); len(o.Items) > 0 {
		if err := parseKMS(&result, o, "seal", 2); err != nil {
			return nil, errwrap.Wrapf("error parsing 'seal': {{err}}", err)
		}
	}

	if o := list.Filter("kms"); len(o.Items) > 0 {
		if err := parseKMS(&result, o, "kms", 2); err != nil {
			return nil, errwrap.Wrapf("error parsing 'kms': {{err}}", err)
		}
	}

	if o := list.Filter("entropy"); len(o.Items) > 0 {
		if err := ParseEntropy(&result, o, "entropy"); err != nil {
			return nil, errwrap.Wrapf("error parsing 'entropy': {{err}}", err)
		}
	}

	if o := list.Filter("listener"); len(o.Items) > 0 {
		if err := ParseListeners(&result, o); err != nil {
			return nil, errwrap.Wrapf("error parsing 'listener': {{err}}", err)
		}
	}

	if o := list.Filter("telemetry"); len(o.Items) > 0 {
		if err := parseTelemetry(&result, o); err != nil {
			return nil, errwrap.Wrapf("error parsing 'telemetry': {{err}}", err)
		}
	}

	entConfig := &(result.EntSharedConfig)
	if err := entConfig.ParseConfig(list); err != nil {
		return nil, errwrap.Wrapf("error parsing enterprise config: {{err}}", err)
	}

	return &result, nil
}

func parseKMS(result *SharedConfig, list *ast.ObjectList, blockName string, maxKMS int) error {
	if len(list.Items) > maxKMS {
		return fmt.Errorf("only two or less %q blocks are permitted", blockName)
	}

	seals := make([]*KMS, 0, len(list.Items))
	for _, item := range list.Items {
		key := blockName
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		var m map[string]string
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
		}
		var disabled bool
		var err error
		if v, ok := m["disabled"]; ok {
			disabled, err = strconv.ParseBool(v)
			if err != nil {
				return multierror.Prefix(err, fmt.Sprintf("%s.%s:", blockName, key))
			}
			delete(m, "disabled")
		}
		seals = append(seals, &KMS{
			Type:     strings.ToLower(key),
			Purpose:  strings.ToLower(m["purpose"]),
			Disabled: disabled,
			Config:   m,
		})
	}

	result.Seals = seals

	return nil
}

func parseTelemetry(result *SharedConfig, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'telemetry' block is permitted")
	}

	// Get our one item
	item := list.Items[0]

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

	if result.Telemetry.PrometheusRetentionTimeRaw != nil {
		var err error
		if result.Telemetry.PrometheusRetentionTime, err = parseutil.ParseDurationSecond(result.Telemetry.PrometheusRetentionTimeRaw); err != nil {
			return err
		}
		result.Telemetry.PrometheusRetentionTimeRaw = nil
	} else {
		result.Telemetry.PrometheusRetentionTime = PrometheusDefaultRetentionTime
	}

	return nil
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
				"config": ln.rawConfig,
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
		}
		result["telemetry"] = sanitizedTelemetry
	}

	return result
}
