// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package metricregistry

import (
	"testing"

	promsink "github.com/hashicorp/go-metrics/compat/prometheus"
	"github.com/stretchr/testify/require"
)

var testGauges = []GaugeDefinition{
	{
		Name: []string{"test_gauge"},
		Help: "A test gauge",
	},
	{
		Name: []string{"test_gauge2"},
		Help: "Another test gauge",
	},
}

var testCounters = []CounterDefinition{
	{
		Name: []string{"test_counter"},
		Help: "A test counter",
	},
	{
		Name: []string{"test_counter2"},
		Help: "Another test counter",
	},
}

var testSummaries = []SummaryDefinition{
	{
		Name: []string{"test_summary"},
		Help: "A test summary",
	},
	{
		Name: []string{"test_summary2"},
		Help: "Another test summary",
	},
}

func TestMetricRegistry(t *testing.T) {
	// Register some metrics
	RegisterGauges(testGauges)
	RegisterCounters(testCounters)
	RegisterSummaries(testSummaries)

	var opts promsink.PrometheusOpts

	// Add some pre-existing metrics to ensure merge is really a merge
	opts.GaugeDefinitions = []promsink.GaugeDefinition{
		{
			Name: []string{"preexisting_gauge"},
			Help: "A pre-existing gauge",
		},
	}
	opts.CounterDefinitions = []promsink.CounterDefinition{
		{
			Name: []string{"preexisting_counter"},
			Help: "A pre-existing counter",
		},
	}
	opts.SummaryDefinitions = []promsink.SummaryDefinition{
		{
			Name: []string{"preexisting_summary"},
			Help: "A pre-existing summary",
		},
	}

	MergeDefinitions(&opts)

	require.Len(t, opts.GaugeDefinitions, 3)
	require.Len(t, opts.CounterDefinitions, 3)
	require.Len(t, opts.SummaryDefinitions, 3)

	wantGauges := []string{"test_gauge", "test_gauge2", "preexisting_gauge"}
	wantGaugeHelp := []string{"A test gauge", "Another test gauge", "A pre-existing gauge"}
	gotGauges := reduce(opts.GaugeDefinitions, nil, func(r []string, d promsink.GaugeDefinition) []string {
		return append(r, d.Name[0])
	})
	gotGaugeHelp := reduce(opts.GaugeDefinitions, nil, func(r []string, d promsink.GaugeDefinition) []string {
		return append(r, d.Help)
	})

	require.ElementsMatch(t, wantGauges, gotGauges)
	require.ElementsMatch(t, wantGaugeHelp, gotGaugeHelp)

	wantCounters := []string{"test_counter", "test_counter2", "preexisting_counter"}
	wantCounterHelp := []string{"A test counter", "Another test counter", "A pre-existing counter"}
	gotCounters := reduce(opts.CounterDefinitions, nil, func(r []string, d promsink.CounterDefinition) []string {
		return append(r, d.Name[0])
	})
	gotCounterHelp := reduce(opts.CounterDefinitions, nil, func(r []string, d promsink.CounterDefinition) []string {
		return append(r, d.Help)
	})

	require.ElementsMatch(t, wantCounters, gotCounters)
	require.ElementsMatch(t, wantCounterHelp, gotCounterHelp)

	wantSummaries := []string{"test_summary", "test_summary2", "preexisting_summary"}
	wantSummaryHelp := []string{"A test summary", "Another test summary", "A pre-existing summary"}
	gotSummaries := reduce(opts.SummaryDefinitions, nil, func(r []string, d promsink.SummaryDefinition) []string {
		return append(r, d.Name[0])
	})
	gotSummaryHelp := reduce(opts.SummaryDefinitions, nil, func(r []string, d promsink.SummaryDefinition) []string {
		return append(r, d.Help)
	})

	require.ElementsMatch(t, wantSummaries, gotSummaries)
	require.ElementsMatch(t, wantSummaryHelp, gotSummaryHelp)
}

func reduce[T, R any](s []T, r R, f func(R, T) R) R {
	for _, v := range s {
		r = f(r, v)
	}
	return r
}
