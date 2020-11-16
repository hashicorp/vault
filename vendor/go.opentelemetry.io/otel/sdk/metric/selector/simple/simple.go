// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simple // import "go.opentelemetry.io/otel/sdk/metric/selector/simple"

import (
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type (
	selectorInexpensive struct{}
	selectorExact       struct{}
	selectorSketch      struct {
		config *ddsketch.Config
	}
	selectorHistogram struct {
		boundaries []float64
	}
)

var (
	_ export.AggregatorSelector = selectorInexpensive{}
	_ export.AggregatorSelector = selectorSketch{}
	_ export.AggregatorSelector = selectorExact{}
	_ export.AggregatorSelector = selectorHistogram{}
)

// NewWithInexpensiveDistribution returns a simple aggregation selector
// that uses counter, minmaxsumcount and minmaxsumcount aggregators
// for the three kinds of metric.  This selector is faster and uses
// less memory than the others because minmaxsumcount does not
// aggregate quantile information.
func NewWithInexpensiveDistribution() export.AggregatorSelector {
	return selectorInexpensive{}
}

// NewWithSketchDistribution returns a simple aggregation selector that
// uses counter, ddsketch, and ddsketch aggregators for the three
// kinds of metric.  This selector uses more cpu and memory than the
// NewWithInexpensiveDistribution because it uses one DDSketch per distinct
// instrument and label set.
func NewWithSketchDistribution(config *ddsketch.Config) export.AggregatorSelector {
	return selectorSketch{
		config: config,
	}
}

// NewWithExactDistribution returns a simple aggregation selector that uses
// counter, array, and array aggregators for the three kinds of metric.
// This selector uses more memory than the NewWithSketchDistribution
// because it aggregates an array of all values, therefore is able to
// compute exact quantiles.
func NewWithExactDistribution() export.AggregatorSelector {
	return selectorExact{}
}

// NewWithHistogramDistribution returns a simple aggregation selector that uses counter,
// histogram, and histogram aggregators for the three kinds of metric. This
// selector uses more memory than the NewWithInexpensiveDistribution because it
// uses a counter per bucket.
func NewWithHistogramDistribution(boundaries []float64) export.AggregatorSelector {
	return selectorHistogram{boundaries: boundaries}
}

func sumAggs(aggPtrs []*export.Aggregator) {
	aggs := sum.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func lastValueAggs(aggPtrs []*export.Aggregator) {
	aggs := lastvalue.New(len(aggPtrs))
	for i := range aggPtrs {
		*aggPtrs[i] = &aggs[i]
	}
}

func (selectorInexpensive) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderKind:
		aggs := minmaxsumcount.New(len(aggPtrs), descriptor)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}

func (s selectorSketch) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderKind:
		aggs := ddsketch.New(len(aggPtrs), descriptor, s.config)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}

func (selectorExact) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderKind:
		aggs := array.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}

func (s selectorHistogram) AggregatorFor(descriptor *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	switch descriptor.MetricKind() {
	case metric.ValueObserverKind:
		lastValueAggs(aggPtrs)
	case metric.ValueRecorderKind:
		aggs := histogram.New(len(aggPtrs), descriptor, s.boundaries)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		sumAggs(aggPtrs)
	}
}
