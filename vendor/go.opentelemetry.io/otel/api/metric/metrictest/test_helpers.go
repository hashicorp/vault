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

package metrictest

import (
	"testing"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/label"
)

// Measured is the helper struct which provides flat representation of recorded measurements
// to simplify testing
type Measured struct {
	Name                   string
	InstrumentationName    string
	InstrumentationVersion string
	Labels                 map[label.Key]label.Value
	Number                 metric.Number
}

// LabelsToMap converts label set to keyValue map, to be easily used in tests
func LabelsToMap(kvs ...label.KeyValue) map[label.Key]label.Value {
	m := map[label.Key]label.Value{}
	for _, label := range kvs {
		m[label.Key] = label.Value
	}
	return m
}

// AsStructs converts recorded batches to array of flat, readable Measured helper structures
func AsStructs(batches []Batch) []Measured {
	var r []Measured
	for _, batch := range batches {
		for _, m := range batch.Measurements {
			r = append(r, Measured{
				Name:                   m.Instrument.Descriptor().Name(),
				InstrumentationName:    m.Instrument.Descriptor().InstrumentationName(),
				InstrumentationVersion: m.Instrument.Descriptor().InstrumentationVersion(),
				Labels:                 LabelsToMap(batch.Labels...),
				Number:                 m.Number,
			})
		}
	}
	return r
}

// ResolveNumberByKind takes defined metric descriptor creates a concrete typed metric number
func ResolveNumberByKind(t *testing.T, kind metric.NumberKind, value float64) metric.Number {
	t.Helper()
	switch kind {
	case metric.Int64NumberKind:
		return metric.NewInt64Number(int64(value))
	case metric.Float64NumberKind:
		return metric.NewFloat64Number(value)
	}
	panic("invalid number kind")
}
