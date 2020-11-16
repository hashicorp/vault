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

package array // import "go.opentelemetry.io/otel/sdk/metric/aggregator/array"

import (
	"context"
	"math"
	"sort"
	"sync"
	"unsafe"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

type (
	// Aggregator aggregates events that form a distribution, keeping
	// an array with the exact set of values.
	Aggregator struct {
		lock   sync.Mutex
		sum    metric.Number
		points points
	}

	points []metric.Number
)

var _ export.Aggregator = &Aggregator{}
var _ aggregation.MinMaxSumCount = &Aggregator{}
var _ aggregation.Distribution = &Aggregator{}
var _ aggregation.Points = &Aggregator{}

// New returns a new array aggregator, which aggregates recorded
// measurements by storing them in an array.  This type uses a mutex
// for Update() and SynchronizedMove() concurrency.
func New(cnt int) []Aggregator {
	return make([]Aggregator, cnt)
}

// Aggregation returns an interface for reading the state of this aggregator.
func (c *Aggregator) Aggregation() aggregation.Aggregation {
	return c
}

// Kind returns aggregation.ExactKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.ExactKind
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	return c.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return int64(len(c.points)), nil
}

// Max returns the maximum value in the checkpoint.
func (c *Aggregator) Max() (metric.Number, error) {
	return c.points.Quantile(1)
}

// Min returns the mininum value in the checkpoint.
func (c *Aggregator) Min() (metric.Number, error) {
	return c.points.Quantile(0)
}

// Quantile returns the estimated quantile of data in the checkpoint.
// It is an error if `q` is less than 0 or greated than 1.
func (c *Aggregator) Quantile(q float64) (metric.Number, error) {
	return c.points.Quantile(q)
}

// Points returns access to the raw data set.
func (c *Aggregator) Points() ([]metric.Number, error) {
	return c.points, nil
}

// SynchronizedMove saves the current state to oa and resets the current state to
// the empty set, taking a lock to prevent concurrent Update() calls.
func (c *Aggregator) SynchronizedMove(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.lock.Lock()
	o.points, c.points = c.points, nil
	o.sum, c.sum = c.sum, 0
	c.lock.Unlock()

	// TODO: This sort should be done lazily, only when quantiles
	// are requested.  The SDK specification says you can use this
	// aggregator to simply list values in the order they were
	// received as an alternative to requesting quantile information.
	o.sort(desc.NumberKind())
	return nil
}

// Update adds the recorded measurement to the current data set.
// Update takes a lock to prevent concurrent Update() and SynchronizedMove()
// calls.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	c.lock.Lock()
	c.points = append(c.points, number)
	c.sum.AddNumber(desc.NumberKind(), number)
	c.lock.Unlock()

	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	// Note: Current assumption is that `o` was checkpointed,
	// therefore is already sorted.  See the TODO above, since
	// this is an open question.
	c.sum.AddNumber(desc.NumberKind(), o.sum)
	c.points = combine(c.points, o.points, desc.NumberKind())
	return nil
}

func (c *Aggregator) sort(kind metric.NumberKind) {
	switch kind {
	case metric.Float64NumberKind:
		sort.Float64s(*(*[]float64)(unsafe.Pointer(&c.points)))

	case metric.Int64NumberKind:
		sort.Sort(&c.points)

	default:
		// NOTE: This can't happen because the SDK doesn't
		// support uint64-kind metric instruments.
		panic("Impossible case")
	}
}

func combine(a, b points, kind metric.NumberKind) points {
	result := make(points, 0, len(a)+len(b))

	for len(a) != 0 && len(b) != 0 {
		if a[0].CompareNumber(kind, b[0]) < 0 {
			result = append(result, a[0])
			a = a[1:]
		} else {
			result = append(result, b[0])
			b = b[1:]
		}
	}
	result = append(result, a...)
	result = append(result, b...)
	return result
}

func (p *points) Len() int {
	return len(*p)
}

func (p *points) Less(i, j int) bool {
	// Note this is specialized for int64, because float64 is
	// handled by `sort.Float64s` and uint64 numbers never appear
	// in this data.
	return int64((*p)[i]) < int64((*p)[j])
}

func (p *points) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// Quantile returns the least X such that Pr(x<X)>=q, where X is an
// element of the data set.  This uses the "Nearest-Rank" definition
// of a quantile.
func (p *points) Quantile(q float64) (metric.Number, error) {
	if len(*p) == 0 {
		return 0, aggregation.ErrNoData
	}

	if q < 0 || q > 1 {
		return 0, aggregation.ErrInvalidQuantile
	}

	if q == 0 || len(*p) == 1 {
		return (*p)[0], nil
	} else if q == 1 {
		return (*p)[len(*p)-1], nil
	}

	position := float64(len(*p)-1) * q
	ceil := int(math.Ceil(position))
	return (*p)[ceil], nil
}
