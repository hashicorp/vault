// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/circonus"
	"github.com/armon/go-metrics/datadog"
	"github.com/armon/go-metrics/prometheus"
	stackdriver "github.com/google/go-metrics-stackdriver"
	stackdrivervault "github.com/google/go-metrics-stackdriver/vault"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/mitchellh/cli"
	"google.golang.org/api/option"
)

const (
	PrometheusDefaultRetentionTime    = 24 * time.Hour
	UsageGaugeDefaultPeriod           = 10 * time.Minute
	MaximumGaugeCardinalityDefault    = 500
	LeaseMetricsEpsilonDefault        = time.Hour
	NumLeaseMetricsTimeBucketsDefault = 168
)

// Telemetry is the telemetry configuration for the server
type Telemetry struct {
	FoundKeys    []string     `hcl:",decodedFields"`
	UnusedKeys   UnusedKeyMap `hcl:",unusedKeyPositions"`
	StatsiteAddr string       `hcl:"statsite_address"`
	StatsdAddr   string       `hcl:"statsd_address"`

	DisableHostname     bool   `hcl:"disable_hostname"`
	EnableHostnameLabel bool   `hcl:"enable_hostname_label"`
	MetricsPrefix       string `hcl:"metrics_prefix"`
	UsageGaugePeriod    time.Duration
	UsageGaugePeriodRaw interface{} `hcl:"usage_gauge_period,alias:UsageGaugePeriod"`

	MaximumGaugeCardinality int `hcl:"maximum_gauge_cardinality"`

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
	// behavior is to add new metrics as they are encountered. If the metric already exists in the
	// check, it will *NOT* be activated. This setting overrides that behavior.
	// Default: "false"
	CirconusCheckForceMetricActivation string `hcl:"circonus_check_force_metric_activation"`
	// CirconusCheckInstanceID serves to uniquely identify the metrics coming from this "instance".
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
	CirconusCheckTags string `hcl:"circonus_check_tags"`
	// CirconusCheckDisplayName is the name for the check which will be displayed in the Circonus UI.
	// Default: value of CirconusCheckInstanceID
	CirconusCheckDisplayName string `hcl:"circonus_check_display_name"`
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

	// Prometheus:
	// PrometheusRetentionTime is the retention time for prometheus metrics if greater than 0.
	// Default: 24h
	PrometheusRetentionTime    time.Duration `hcl:"-"`
	PrometheusRetentionTimeRaw interface{}   `hcl:"prometheus_retention_time"`

	// Stackdriver:
	// StackdriverProjectID is the project to publish stackdriver metrics to.
	StackdriverProjectID string `hcl:"stackdriver_project_id"`
	// StackdriverLocation is the GCP or AWS region of the monitored resource.
	StackdriverLocation string `hcl:"stackdriver_location"`
	// StackdriverNamespace is the namespace identifier, such as a cluster name.
	StackdriverNamespace string `hcl:"stackdriver_namespace"`
	// StackdriverDebugLogs will write additional stackdriver related debug logs to stderr.
	StackdriverDebugLogs bool `hcl:"stackdriver_debug_logs"`

	// How often metrics for lease expiry will be aggregated
	LeaseMetricsEpsilon    time.Duration
	LeaseMetricsEpsilonRaw interface{} `hcl:"lease_metrics_epsilon"`

	// Number of buckets by time that will be used in lease aggregation
	NumLeaseMetricsTimeBuckets int `hcl:"num_lease_metrics_buckets"`

	// Whether or not telemetry should add labels for namespaces
	LeaseMetricsNameSpaceLabels bool `hcl:"add_lease_metrics_namespace_labels"`

	// FilterDefault is the default for whether to allow a metric that's not
	// covered by the prefix filter.
	FilterDefault *bool `hcl:"filter_default"`

	// PrefixFilter is a list of filter rules to apply for allowing
	// or blocking metrics by prefix.
	PrefixFilter []string `hcl:"prefix_filter"`

	// Whether or not telemetry should include the mount point in the rollback
	// metrics
	RollbackMetricsIncludeMountPoint bool `hcl:"add_mount_point_rollback_metrics"`
}

func (t *Telemetry) Validate(source string) []ConfigError {
	return ValidateUnusedFields(t.UnusedKeys, source)
}

func (t *Telemetry) GoString() string {
	return fmt.Sprintf("*%#v", *t)
}

func parseTelemetry(result *SharedConfig, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'telemetry' block is permitted")
	}

	// Get our one item
	item := list.Items[0]

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

	if result.Telemetry.UsageGaugePeriodRaw != nil {
		if result.Telemetry.UsageGaugePeriodRaw == "none" {
			result.Telemetry.UsageGaugePeriod = 0
		} else {
			var err error
			if result.Telemetry.UsageGaugePeriod, err = parseutil.ParseDurationSecond(result.Telemetry.UsageGaugePeriodRaw); err != nil {
				return err
			}
			result.Telemetry.UsageGaugePeriodRaw = nil
		}
	} else {
		result.Telemetry.UsageGaugePeriod = UsageGaugeDefaultPeriod
	}

	if result.Telemetry.MaximumGaugeCardinality == 0 {
		result.Telemetry.MaximumGaugeCardinality = MaximumGaugeCardinalityDefault
	}

	if result.Telemetry.LeaseMetricsEpsilonRaw != nil {
		if result.Telemetry.LeaseMetricsEpsilonRaw == "none" {
			result.Telemetry.LeaseMetricsEpsilonRaw = 0
		} else {
			var err error
			if result.Telemetry.LeaseMetricsEpsilon, err = parseutil.ParseDurationSecond(result.Telemetry.LeaseMetricsEpsilonRaw); err != nil {
				return err
			}
			result.Telemetry.LeaseMetricsEpsilonRaw = nil
		}
	} else {
		result.Telemetry.LeaseMetricsEpsilon = LeaseMetricsEpsilonDefault
	}

	if result.Telemetry.NumLeaseMetricsTimeBuckets == 0 {
		result.Telemetry.NumLeaseMetricsTimeBuckets = NumLeaseMetricsTimeBucketsDefault
	}

	return nil
}

type SetupTelemetryOpts struct {
	Config      *Telemetry
	Ui          cli.Ui
	ServiceName string
	DisplayName string
	UserAgent   string
	ClusterName string
}

// SetupTelemetry is used to setup the telemetry sub-systems and returns the
// in-memory sink to be used in http configuration
func SetupTelemetry(opts *SetupTelemetryOpts) (*metrics.InmemSink, *metricsutil.ClusterMetricSink, bool, error) {
	if opts == nil {
		return nil, nil, false, errors.New("nil opts passed into SetupTelemetry")
	}

	if opts.Config == nil {
		opts.Config = &Telemetry{}
	}

	/* Setup telemetry
	Aggregate on 10 second intervals for 1 minute. Expose the
	metrics over stderr when there is a SIGUSR1 received.
	*/
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	if opts.Config.MetricsPrefix != "" {
		opts.ServiceName = opts.Config.MetricsPrefix
	}

	metricsConf := metrics.DefaultConfig(opts.ServiceName)
	metricsConf.EnableHostname = !opts.Config.DisableHostname
	metricsConf.EnableHostnameLabel = opts.Config.EnableHostnameLabel
	if opts.Config.FilterDefault != nil {
		metricsConf.FilterDefault = *opts.Config.FilterDefault
	}

	// Configure the statsite sink
	var fanout metrics.FanoutSink
	var prometheusEnabled bool

	// Configure the Prometheus sink
	if opts.Config.PrometheusRetentionTime != 0 {
		prometheusEnabled = true
		prometheusOpts := prometheus.PrometheusOpts{
			Expiration: opts.Config.PrometheusRetentionTime,
		}

		sink, err := prometheus.NewPrometheusSinkFrom(prometheusOpts)
		if err != nil {
			return nil, nil, false, err
		}
		fanout = append(fanout, sink)
	}

	if opts.Config.StatsiteAddr != "" {
		sink, err := metrics.NewStatsiteSink(opts.Config.StatsiteAddr)
		if err != nil {
			return nil, nil, false, err
		}
		fanout = append(fanout, sink)
	}

	// Configure the statsd sink
	if opts.Config.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(opts.Config.StatsdAddr)
		if err != nil {
			return nil, nil, false, err
		}
		fanout = append(fanout, sink)
	}

	// Configure the Circonus sink
	if opts.Config.CirconusAPIToken != "" || opts.Config.CirconusCheckSubmissionURL != "" {
		cfg := &circonus.Config{}
		cfg.Interval = opts.Config.CirconusSubmissionInterval
		cfg.CheckManager.API.TokenKey = opts.Config.CirconusAPIToken
		cfg.CheckManager.API.TokenApp = opts.Config.CirconusAPIApp
		cfg.CheckManager.API.URL = opts.Config.CirconusAPIURL
		cfg.CheckManager.Check.SubmissionURL = opts.Config.CirconusCheckSubmissionURL
		cfg.CheckManager.Check.ID = opts.Config.CirconusCheckID
		cfg.CheckManager.Check.ForceMetricActivation = opts.Config.CirconusCheckForceMetricActivation
		cfg.CheckManager.Check.InstanceID = opts.Config.CirconusCheckInstanceID
		cfg.CheckManager.Check.SearchTag = opts.Config.CirconusCheckSearchTag
		cfg.CheckManager.Check.DisplayName = opts.Config.CirconusCheckDisplayName
		cfg.CheckManager.Check.Tags = opts.Config.CirconusCheckTags
		cfg.CheckManager.Broker.ID = opts.Config.CirconusBrokerID
		cfg.CheckManager.Broker.SelectTag = opts.Config.CirconusBrokerSelectTag

		if cfg.CheckManager.API.TokenApp == "" {
			cfg.CheckManager.API.TokenApp = opts.ServiceName
		}

		if cfg.CheckManager.Check.DisplayName == "" {
			cfg.CheckManager.Check.DisplayName = opts.DisplayName
		}

		if cfg.CheckManager.Check.SearchTag == "" {
			cfg.CheckManager.Check.SearchTag = fmt.Sprintf("service:%s", opts.ServiceName)
		}

		sink, err := circonus.NewCirconusSink(cfg)
		if err != nil {
			return nil, nil, false, err
		}
		sink.Start()
		fanout = append(fanout, sink)
	}

	if opts.Config.DogStatsDAddr != "" {
		var tags []string

		if opts.Config.DogStatsDTags != nil {
			tags = opts.Config.DogStatsDTags
		}

		sink, err := datadog.NewDogStatsdSink(opts.Config.DogStatsDAddr, metricsConf.HostName)
		if err != nil {
			return nil, nil, false, fmt.Errorf("failed to start DogStatsD sink: %w", err)
		}
		sink.SetTags(tags)
		fanout = append(fanout, sink)
	}

	// Configure the stackdriver sink
	if opts.Config.StackdriverProjectID != "" {
		client, err := monitoring.NewMetricClient(context.Background(), option.WithUserAgent(opts.UserAgent))
		if err != nil {
			return nil, nil, false, fmt.Errorf("Failed to create stackdriver client: %v", err)
		}
		sink := stackdriver.NewSink(client, &stackdriver.Config{
			LabelExtractor: stackdrivervault.Extractor,
			Bucketer:       stackdrivervault.Bucketer,
			ProjectID:      opts.Config.StackdriverProjectID,
			Location:       opts.Config.StackdriverLocation,
			Namespace:      opts.Config.StackdriverNamespace,
			DebugLogs:      opts.Config.StackdriverDebugLogs,
		})
		fanout = append(fanout, sink)
	}

	// Initialize the global sink
	if len(fanout) > 1 {
		// Hostname enabled will create poor quality metrics name for prometheus
		if !opts.Config.DisableHostname {
			opts.Ui.Warn("telemetry.disable_hostname has been set to false. Recommended setting is true for Prometheus to avoid poorly named metrics.")
		}
	} else {
		metricsConf.EnableHostname = false
	}
	fanout = append(fanout, inm)
	globalMetrics, err := metrics.NewGlobal(metricsConf, fanout)
	if err != nil {
		return nil, nil, false, err
	}

	// Intialize a wrapper around the global sink; this will be passed to Core
	// and to any backend.
	wrapper := metricsutil.NewClusterMetricSink(opts.ClusterName, globalMetrics)
	wrapper.MaxGaugeCardinality = opts.Config.MaximumGaugeCardinality
	wrapper.GaugeInterval = opts.Config.UsageGaugePeriod
	wrapper.TelemetryConsts.LeaseMetricsEpsilon = opts.Config.LeaseMetricsEpsilon
	wrapper.TelemetryConsts.LeaseMetricsNameSpaceLabels = opts.Config.LeaseMetricsNameSpaceLabels
	wrapper.TelemetryConsts.NumLeaseMetricsTimeBuckets = opts.Config.NumLeaseMetricsTimeBuckets
	wrapper.TelemetryConsts.RollbackMetricsIncludeMountPoint = opts.Config.RollbackMetricsIncludeMountPoint

	// Parse the metric filters
	telemetryAllowedPrefixes, telemetryBlockedPrefixes, err := parsePrefixFilter(opts.Config.PrefixFilter)
	if err != nil {
		return nil, nil, false, err
	}

	metrics.UpdateFilter(telemetryAllowedPrefixes, telemetryBlockedPrefixes)
	return inm, wrapper, prometheusEnabled, nil
}

func parsePrefixFilter(prefixFilters []string) ([]string, []string, error) {
	var telemetryAllowedPrefixes, telemetryBlockedPrefixes []string

	for _, rule := range prefixFilters {
		if rule == "" {
			return nil, nil, fmt.Errorf("Cannot have empty filter rule in prefix_filter")
		}
		switch rule[0] {
		case '+':
			telemetryAllowedPrefixes = append(telemetryAllowedPrefixes, rule[1:])
		case '-':
			telemetryBlockedPrefixes = append(telemetryBlockedPrefixes, rule[1:])
		default:
			return nil, nil, fmt.Errorf("Filter rule must begin with either '+' or '-': %q", rule)
		}
	}
	return telemetryAllowedPrefixes, telemetryBlockedPrefixes, nil
}
