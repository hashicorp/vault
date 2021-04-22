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

package stdout // import "go.opentelemetry.io/otel/exporters/stdout"

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

type metricExporter struct {
	config Config
}

var _ exportmetric.Exporter = &metricExporter{}

type line struct {
	Name      string      `json:"Name"`
	Min       interface{} `json:"Min,omitempty"`
	Max       interface{} `json:"Max,omitempty"`
	Sum       interface{} `json:"Sum,omitempty"`
	Count     interface{} `json:"Count,omitempty"`
	LastValue interface{} `json:"Last,omitempty"`

	// Note: this is a pointer because omitempty doesn't work when time.IsZero()
	Timestamp *time.Time `json:"Timestamp,omitempty"`
}

func (e *metricExporter) ExportKindFor(desc *metric.Descriptor, kind aggregation.Kind) exportmetric.ExportKind {
	return exportmetric.StatelessExportKindSelector().ExportKindFor(desc, kind)
}

func (e *metricExporter) Export(_ context.Context, checkpointSet exportmetric.CheckpointSet) error {
	if e.config.DisableMetricExport {
		return nil
	}
	var aggError error
	var batch []line
	aggError = checkpointSet.ForEach(e, func(record exportmetric.Record) error {
		desc := record.Descriptor()
		agg := record.Aggregation()
		kind := desc.NumberKind()
		encodedResource := record.Resource().Encoded(e.config.LabelEncoder)

		var instLabels []attribute.KeyValue
		if name := desc.InstrumentationName(); name != "" {
			instLabels = append(instLabels, attribute.String("instrumentation.name", name))
			if version := desc.InstrumentationVersion(); version != "" {
				instLabels = append(instLabels, attribute.String("instrumentation.version", version))
			}
		}
		instSet := attribute.NewSet(instLabels...)
		encodedInstLabels := instSet.Encoded(e.config.LabelEncoder)

		var expose line

		if sum, ok := agg.(aggregation.Sum); ok {
			value, err := sum.Sum()
			if err != nil {
				return err
			}
			expose.Sum = value.AsInterface(kind)
		}

		if mmsc, ok := agg.(aggregation.MinMaxSumCount); ok {
			count, err := mmsc.Count()
			if err != nil {
				return err
			}
			expose.Count = count

			max, err := mmsc.Max()
			if err != nil {
				return err
			}
			expose.Max = max.AsInterface(kind)

			min, err := mmsc.Min()
			if err != nil {
				return err
			}
			expose.Min = min.AsInterface(kind)
		} else if lv, ok := agg.(aggregation.LastValue); ok {
			value, timestamp, err := lv.LastValue()
			if err != nil {
				return err
			}
			expose.LastValue = value.AsInterface(kind)

			if e.config.Timestamps {
				expose.Timestamp = &timestamp
			}
		}

		var encodedLabels string
		iter := record.Labels().Iter()
		if iter.Len() > 0 {
			encodedLabels = record.Labels().Encoded(e.config.LabelEncoder)
		}

		var sb strings.Builder

		sb.WriteString(desc.Name())

		if len(encodedLabels) > 0 || len(encodedResource) > 0 || len(encodedInstLabels) > 0 {
			sb.WriteRune('{')
			sb.WriteString(encodedResource)
			if len(encodedInstLabels) > 0 && len(encodedResource) > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(encodedInstLabels)
			if len(encodedLabels) > 0 && (len(encodedInstLabels) > 0 || len(encodedResource) > 0) {
				sb.WriteRune(',')
			}
			sb.WriteString(encodedLabels)
			sb.WriteRune('}')
		}

		expose.Name = sb.String()

		batch = append(batch, expose)
		return nil
	})
	if len(batch) == 0 {
		return aggError
	}

	data, err := e.marshal(batch)
	if err != nil {
		return err
	}
	fmt.Fprintln(e.config.Writer, string(data))

	return aggError
}

// marshal v with approriate indentation.
func (e *metricExporter) marshal(v interface{}) ([]byte, error) {
	if e.config.PrettyPrint {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}
