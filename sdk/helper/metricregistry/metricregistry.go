// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package metricregistry is a helper that allows Vault code or plugins that are
// compiled into Vault to pre-define any metrics they will emit to go-metrics at
// init time. Metrics registered this way will always be reported by the
// go-metrics PrometheusSink if it is used so infrequently updated metrics are
// always present. It is not required to pre-register metrics to use go-metrics
// with Prometheus, but it's preferable as it makes them behave more like the
// Prometheus ecosystem expects, being always present and with a helpful
// description in the output which some systems use to help operators explore
// metrics.
//
// Note that this will not work for external Vault plugins since they are in a
// separate process and only started after Vault's metrics sink is already
// configured.
package metricregistry

import promsink "github.com/armon/go-metrics/prometheus"

var Registry definitionRegistry

// Re-export these types so that we don't have the whole of Vault depending
// directly on go-metrics prometheus sink and can buffer changes if needed
type (
	// GaugeDefinition provides the name and help text of a gauge metric that will
	// be exported via go-metrics' Prometheus sink if enabled.
	GaugeDefinition promsink.GaugeDefinition

	// CounterDefinition provides the name and help text of a counter metric that
	// will be exported via go-metrics' Prometheus sink if enabled.
	CounterDefinition promsink.CounterDefinition

	// SummaryDefinition provides the name and help text of a summary metric that
	// will be exported via go-metrics' Prometheus sink if enabled.
	SummaryDefinition promsink.SummaryDefinition
)

// definitionRegistry is a central place for packages to register their metrics
// definitions during init so that we can correctly report metrics to Prometheus
// even before they are observed. Typically there is one global instance.
type definitionRegistry struct {
	gauges    []GaugeDefinition
	counters  []CounterDefinition
	summaries []SummaryDefinition
}

// RegisterGauges is intended to be called during init. It accesses global state
// without synchronization. Statically defined definitions should be registered
// during `init` of a package read to be configured if the prometheus sink is
// enabled in configuration. Registering metrics is not mandatory but it is
// strongly preferred as it ensures they are always output even before the are
// observed which makes dashboards much easier to work with, provides helpful
// descriptions and matches Prometheus eco system expectations. It also prevents
// the metrics ever being expired which means users don't need to work around
// that quirk of go-metrics by setting long prometheus retention times. All
// registered metrics will report 0 until an actual observation is made.
func RegisterGauges(defs []GaugeDefinition) {
	Registry.gauges = append(Registry.gauges, defs...)
}

// RegisterCounters is intended to be called during init. It accesses global
// state without synchronization. Statically defined definitions should be
// registered during `init` of a package read to be configured if the prometheus
// sink is enabled in configuration. Registering metrics is not mandatory but it
// is strongly preferred as it ensures they are always output even before the
// are observed which makes dashboards much easier to work with, provides
// helpful descriptions and matches Prometheus eco system expectations. It also
// prevents the metrics ever being expired which means users don't need to work
// around that quirk of go-metrics by setting long prometheus retention times.
// All registered metrics will report 0 until an actual observation is made.
func RegisterCounters(defs []CounterDefinition) {
	Registry.counters = append(Registry.counters, defs...)
}

// RegisterSummaries is intended to be called during init. It accesses global
// state without synchronization. Statically defined definitions should be
// registered during `init` of a package read to be configured if the prometheus
// sink is enabled in configuration. Registering metrics is not mandatory but it
// is strongly preferred as it ensures they are always output even before the
// are observed which makes dashboards much easier to work with, provides
// helpful descriptions and matches Prometheus eco system expectations. It also
// prevents the metrics ever being expired which means users don't need to work
// around that quirk of go-metrics by setting long prometheus retention times.
// All registered metrics will report 0 until an actual observation is made.
func RegisterSummaries(defs []SummaryDefinition) {
	Registry.summaries = append(Registry.summaries, defs...)
}

// MergeDefinitions adds all registered metrics to any already present in `cfg`
// ready to be passed to the go-metrics prometheus sink. Note it is not safe to
// call this concurrently with registrations or other calls, it's intended this
// is called once only after all registrations (which should be in init
// functions) just before the PrometheusSink is created. Calling more than once
// could result in duplicate metrics definitions being passed unless the cfg is
// different each time for different Prometheus sinks.
func MergeDefinitions(cfg *promsink.PrometheusOpts) {
	for _, g := range Registry.gauges {
		cfg.GaugeDefinitions = append(cfg.GaugeDefinitions, promsink.GaugeDefinition(g))
	}
	for _, c := range Registry.counters {
		cfg.CounterDefinitions = append(cfg.CounterDefinitions, promsink.CounterDefinition(c))
	}
	for _, s := range Registry.summaries {
		cfg.SummaryDefinitions = append(cfg.SummaryDefinitions, promsink.SummaryDefinition(s))
	}
}
