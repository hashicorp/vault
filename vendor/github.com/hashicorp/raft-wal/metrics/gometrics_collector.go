// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metrics

import gometrics "github.com/armon/go-metrics"

// GoMetricsCollector implements a Collector that passes through observations to
// a go-metrics instance. The zero value works, writing metrics to the default
// global instance however to set a prefix or a static set of labels to add to
// each metric observed, or to use a non-global metrics instance use
// NewGoMetricsCollector.
type GoMetricsCollector struct {
	gm     *gometrics.Metrics
	prefix []string
	labels []gometrics.Label
}

// NewGoMetricsCollector returns a GoMetricsCollector that will attach the
// specified name prefix and/or labels to each observation. If gm is nil the
// global metrics instance is used.
func NewGoMetricsCollector(prefix []string, labels []gometrics.Label, gm *gometrics.Metrics) *GoMetricsCollector {
	if gm == nil {
		gm = gometrics.Default()
	}
	return &GoMetricsCollector{
		gm:     gm,
		prefix: prefix,
		labels: labels,
	}
}

// IncrementCounter record val occurrences of the named event. Names will
// follow prometheus conventions with lower_case_and_underscores. We don't
// need any additional labels currently.
func (c *GoMetricsCollector) IncrementCounter(name string, delta uint64) {
	c.gm.IncrCounterWithLabels(c.name(name), float32(delta), c.labels)
}

// SetGauge sets the value of the named gauge overriding any previous value.
func (c *GoMetricsCollector) SetGauge(name string, val uint64) {
	c.gm.SetGaugeWithLabels(c.name(name), float32(val), c.labels)
}

// name returns the metric name as a slice we don't want to risk modifying the
// prefix slice backing array since this might be called concurrently so we
// always allocate a new slice.
func (c *GoMetricsCollector) name(name string) []string {
	var ss []string
	return append(append(ss, c.prefix...), name)
}
