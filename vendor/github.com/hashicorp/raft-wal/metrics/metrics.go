// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metrics

// Collector provides a simple abstraction for counter type metrics that
// the WAL and log verifier can use without depending on a specific metrics
// collector implementation.
type Collector interface {
	// IncrementCounter record val occurrences of the named event. Names will
	// follow prometheus conventions with lower_case_and_underscores. We don't
	// need any additional labels currently.
	IncrementCounter(name string, delta uint64)

	// SetGauge sets the value of the named gauge overriding any previous value.
	SetGauge(name string, val uint64)
}

// Definitions provides a simple description of a set of scalar metrics.
type Definitions struct {
	Counters []Descriptor
	Gauges   []Descriptor
}

// Descriptor describes a specific metric.
type Descriptor struct {
	Name string
	Desc string
}

var _ Collector = &NoOpCollector{}

// NoOpCollector is a Collector that does nothing.
type NoOpCollector struct{}

// IncrementCounter record val occurrences of the named event. Names will
// follow prometheus conventions with lower_case_and_underscores. We don't
// need any additional labels currently.
func (c *NoOpCollector) IncrementCounter(name string, delta uint64) {}

// SetGauge sets the value of the named gauge overriding any previous value.
func (c *NoOpCollector) SetGauge(name string, val uint64) {}
