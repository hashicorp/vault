package configutil

import (
	"context"
	"fmt"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/circonus"
	"github.com/armon/go-metrics/datadog"
	"github.com/armon/go-metrics/prometheus"
	stackdriver "github.com/google/go-metrics-stackdriver"
	"github.com/hashicorp/errwrap"
	"github.com/mitchellh/cli"
	"google.golang.org/api/option"
)

// SetupTelemetry is used to setup the telemetry sub-systems and returns the
// in-memory sink to be used in http configuration
func SetupTelemetry(telConfig *Telemetry, ui cli.Ui, serviceName, displayName, useragent string) (*metrics.InmemSink, bool, error) {
	if telConfig == nil {
		telConfig = &Telemetry{}
	}

	/* Setup telemetry
	Aggregate on 10 second intervals for 1 minute. Expose the
	metrics over stderr when there is a SIGUSR1 received.
	*/
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	if telConfig.MetricsPrefix != "" {
		serviceName = telConfig.MetricsPrefix
	}

	metricsConf := metrics.DefaultConfig(serviceName)
	metricsConf.EnableHostname = !telConfig.DisableHostname
	metricsConf.EnableHostnameLabel = telConfig.EnableHostnameLabel

	// Configure the statsite sink
	var fanout metrics.FanoutSink
	var prometheusEnabled bool

	// Configure the Prometheus sink
	if telConfig.PrometheusRetentionTime != 0 {
		prometheusEnabled = true
		prometheusOpts := prometheus.PrometheusOpts{
			Expiration: telConfig.PrometheusRetentionTime,
		}

		sink, err := prometheus.NewPrometheusSinkFrom(prometheusOpts)
		if err != nil {
			return nil, false, err
		}
		fanout = append(fanout, sink)
	}

	if telConfig.StatsiteAddr != "" {
		sink, err := metrics.NewStatsiteSink(telConfig.StatsiteAddr)
		if err != nil {
			return nil, false, err
		}
		fanout = append(fanout, sink)
	}

	// Configure the statsd sink
	if telConfig.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(telConfig.StatsdAddr)
		if err != nil {
			return nil, false, err
		}
		fanout = append(fanout, sink)
	}

	// Configure the Circonus sink
	if telConfig.CirconusAPIToken != "" || telConfig.CirconusCheckSubmissionURL != "" {
		cfg := &circonus.Config{}
		cfg.Interval = telConfig.CirconusSubmissionInterval
		cfg.CheckManager.API.TokenKey = telConfig.CirconusAPIToken
		cfg.CheckManager.API.TokenApp = telConfig.CirconusAPIApp
		cfg.CheckManager.API.URL = telConfig.CirconusAPIURL
		cfg.CheckManager.Check.SubmissionURL = telConfig.CirconusCheckSubmissionURL
		cfg.CheckManager.Check.ID = telConfig.CirconusCheckID
		cfg.CheckManager.Check.ForceMetricActivation = telConfig.CirconusCheckForceMetricActivation
		cfg.CheckManager.Check.InstanceID = telConfig.CirconusCheckInstanceID
		cfg.CheckManager.Check.SearchTag = telConfig.CirconusCheckSearchTag
		cfg.CheckManager.Check.DisplayName = telConfig.CirconusCheckDisplayName
		cfg.CheckManager.Check.Tags = telConfig.CirconusCheckTags
		cfg.CheckManager.Broker.ID = telConfig.CirconusBrokerID
		cfg.CheckManager.Broker.SelectTag = telConfig.CirconusBrokerSelectTag

		if cfg.CheckManager.API.TokenApp == "" {
			cfg.CheckManager.API.TokenApp = serviceName
		}

		if cfg.CheckManager.Check.DisplayName == "" {
			cfg.CheckManager.Check.DisplayName = displayName
		}

		if cfg.CheckManager.Check.SearchTag == "" {
			cfg.CheckManager.Check.SearchTag = fmt.Sprintf("service:%s", serviceName)
		}

		sink, err := circonus.NewCirconusSink(cfg)
		if err != nil {
			return nil, false, err
		}
		sink.Start()
		fanout = append(fanout, sink)
	}

	if telConfig.DogStatsDAddr != "" {
		var tags []string

		if telConfig.DogStatsDTags != nil {
			tags = telConfig.DogStatsDTags
		}

		sink, err := datadog.NewDogStatsdSink(telConfig.DogStatsDAddr, metricsConf.HostName)
		if err != nil {
			return nil, false, errwrap.Wrapf("failed to start DogStatsD sink: {{err}}", err)
		}
		sink.SetTags(tags)
		fanout = append(fanout, sink)
	}

	// Configure the stackdriver sink
	if telConfig.StackdriverProjectID != "" {
		client, err := monitoring.NewMetricClient(context.Background(), option.WithUserAgent(useragent))
		if err != nil {
			return nil, false, fmt.Errorf("Failed to create stackdriver client: %v", err)
		}
		sink := stackdriver.NewSink(client, &stackdriver.Config{
			ProjectID: telConfig.StackdriverProjectID,
			Location:  telConfig.StackdriverLocation,
			Namespace: telConfig.StackdriverNamespace,
		})
		fanout = append(fanout, sink)
	}

	// Initialize the global sink
	if len(fanout) > 1 {
		// Hostname enabled will create poor quality metrics name for prometheus
		if !telConfig.DisableHostname {
			ui.Warn("telemetry.disable_hostname has been set to false. Recommended setting is true for Prometheus to avoid poorly named metrics.")
		}
	} else {
		metricsConf.EnableHostname = false
	}
	fanout = append(fanout, inm)
	_, err := metrics.NewGlobal(metricsConf, fanout)

	if err != nil {
		return nil, false, err
	}

	return inm, prometheusEnabled, nil
}
