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

package aggregation // import "go.opentelemetry.io/otel/sdk/export/metric/aggregation"

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/api/metric"
)

// These interfaces describe the various ways to access state from an
// Aggregation.

type (
	// Aggregation is an interface returned by the Aggregator
	// containing an interval of metric data.
	Aggregation interface {
		// Kind returns a short identifying string to identify
		// the Aggregator that was used to produce the
		// Aggregation (e.g., "Sum").
		Kind() Kind
	}

	// Sum returns an aggregated sum.
	Sum interface {
		Aggregation
		Sum() (metric.Number, error)
	}

	// Count returns the number of values that were aggregated.
	Count interface {
		Aggregation
		Count() (int64, error)
	}

	// Min returns the minimum value over the set of values that were aggregated.
	Min interface {
		Aggregation
		Min() (metric.Number, error)
	}

	// Max returns the maximum value over the set of values that were aggregated.
	Max interface {
		Aggregation
		Max() (metric.Number, error)
	}

	// Quantile returns an exact or estimated quantile over the
	// set of values that were aggregated.
	Quantile interface {
		Aggregation
		Quantile(float64) (metric.Number, error)
	}

	// LastValue returns the latest value that was aggregated.
	LastValue interface {
		Aggregation
		LastValue() (metric.Number, time.Time, error)
	}

	// Points returns the raw set of values that were aggregated.
	Points interface {
		Aggregation
		Points() ([]metric.Number, error)
	}

	// Buckets represents histogram buckets boundaries and counts.
	//
	// For a Histogram with N defined boundaries, e.g, [x, y, z].
	// There are N+1 counts: [-inf, x), [x, y), [y, z), [z, +inf]
	Buckets struct {
		// Boundaries are floating point numbers, even when
		// aggregating integers.
		Boundaries []float64

		// Counts are floating point numbers to account for
		// the possibility of sampling which allows for
		// non-integer count values.
		Counts []float64
	}

	// Histogram returns the count of events in pre-determined buckets.
	Histogram interface {
		Aggregation
		Count() (int64, error)
		Sum() (metric.Number, error)
		Histogram() (Buckets, error)
	}

	// MinMaxSumCount supports the Min, Max, Sum, and Count interfaces.
	MinMaxSumCount interface {
		Aggregation
		Min() (metric.Number, error)
		Max() (metric.Number, error)
		Sum() (metric.Number, error)
		Count() (int64, error)
	}

	// Distribution supports the Min, Max, Sum, Count, and Quantile
	// interfaces.
	Distribution interface {
		Aggregation
		Min() (metric.Number, error)
		Max() (metric.Number, error)
		Sum() (metric.Number, error)
		Count() (int64, error)
		Quantile(float64) (metric.Number, error)
	}
)

type (
	// Kind is a short name for the Aggregator that produces an
	// Aggregation, used for descriptive purpose only.  Kind is a
	// string to allow user-defined Aggregators.
	//
	// When deciding how to handle an Aggregation, Exporters are
	// encouraged to decide based on conversion to the above
	// interfaces based on strength, not on Kind value, when
	// deciding how to expose metric data.  This enables
	// user-supplied Aggregators to replace builtin Aggregators.
	//
	// For example, test for a Distribution before testing for a
	// MinMaxSumCount, test for a Histogram before testing for a
	// Sum, and so on.
	Kind string
)

const (
	SumKind            Kind = "Sum"
	MinMaxSumCountKind Kind = "MinMaxSumCount"
	HistogramKind      Kind = "Histogram"
	LastValueKind      Kind = "Lastvalue"
	SketchKind         Kind = "Sketch"
	ExactKind          Kind = "Exact"
)

var (
	ErrInvalidQuantile  = fmt.Errorf("the requested quantile is out of range")
	ErrNegativeInput    = fmt.Errorf("negative value is out of range for this instrument")
	ErrNaNInput         = fmt.Errorf("NaN value is an invalid input")
	ErrInconsistentType = fmt.Errorf("inconsistent aggregator types")
	ErrNoSubtraction    = fmt.Errorf("aggregator does not subtract")

	// ErrNoData is returned when (due to a race with collection)
	// the Aggregator is check-pointed before the first value is set.
	// The aggregator should simply be skipped in this case.
	ErrNoData = fmt.Errorf("no data collected by this aggregator")
)

// String returns the string value of Kind.
func (k Kind) String() string {
	return string(k)
}
