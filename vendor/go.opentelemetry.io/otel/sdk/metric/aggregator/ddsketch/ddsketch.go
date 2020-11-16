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
package ddsketch // import "go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"

import (
	"context"
	"math"
	"sync"

	sdk "github.com/DataDog/sketches-go/ddsketch"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

// Config is an alias for the underlying DDSketch config object.
type Config = sdk.Config

// Aggregator aggregates events into a distribution.
type Aggregator struct {
	lock   sync.Mutex
	cfg    *Config
	kind   metric.NumberKind
	sketch *sdk.DDSketch
}

var _ export.Aggregator = &Aggregator{}
var _ aggregation.MinMaxSumCount = &Aggregator{}
var _ aggregation.Distribution = &Aggregator{}

// New returns a new DDSketch aggregator.
func New(cnt int, desc *metric.Descriptor, cfg *Config) []Aggregator {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}
	aggs := make([]Aggregator, cnt)
	for i := range aggs {
		aggs[i] = Aggregator{
			cfg:    cfg,
			kind:   desc.NumberKind(),
			sketch: sdk.NewDDSketch(cfg),
		}
	}
	return aggs
}

// Aggregation returns an interface for reading the state of this aggregator.
func (c *Aggregator) Aggregation() aggregation.Aggregation {
	return c
}

// Kind returns aggregation.SketchKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.SketchKind
}

// NewDefaultConfig returns a new, default DDSketch config.
func NewDefaultConfig() *Config {
	return sdk.NewDefaultConfig()
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	return c.toNumber(c.sketch.Sum()), nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return c.sketch.Count(), nil
}

// Max returns the maximum value in the checkpoint.
func (c *Aggregator) Max() (metric.Number, error) {
	return c.Quantile(1)
}

// Min returns the minimum value in the checkpoint.
func (c *Aggregator) Min() (metric.Number, error) {
	return c.Quantile(0)
}

// Quantile returns the estimated quantile of data in the checkpoint.
// It is an error if `q` is less than 0 or greated than 1.
func (c *Aggregator) Quantile(q float64) (metric.Number, error) {
	if c.sketch.Count() == 0 {
		return 0, aggregation.ErrNoData
	}
	f := c.sketch.Quantile(q)
	if math.IsNaN(f) {
		return 0, aggregation.ErrInvalidQuantile
	}
	return c.toNumber(f), nil
}

func (c *Aggregator) toNumber(f float64) metric.Number {
	if c.kind == metric.Float64NumberKind {
		return metric.NewFloat64Number(f)
	}
	return metric.NewInt64Number(int64(f))
}

// SynchronizedMove saves the current state into oa and resets the current state to
// a new sketch, taking a lock to prevent concurrent Update() calls.
func (c *Aggregator) SynchronizedMove(oa export.Aggregator, _ *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}
	replace := sdk.NewDDSketch(c.cfg)

	c.lock.Lock()
	o.sketch, c.sketch = c.sketch, replace
	c.lock.Unlock()

	return nil
}

// Update adds the recorded measurement to the current data set.
// Update takes a lock to prevent concurrent Update() and SynchronizedMove()
// calls.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sketch.Add(number.CoerceToFloat64(desc.NumberKind()))
	return nil
}

// Merge combines two sketches into one.
func (c *Aggregator) Merge(oa export.Aggregator, d *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.sketch.Merge(o.sketch)
	return nil
}
