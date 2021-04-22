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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/export/metric"
	exporttrace "go.opentelemetry.io/otel/sdk/export/trace"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Exporter struct {
	traceExporter
	metricExporter
}

var (
	_ metric.Exporter          = &Exporter{}
	_ exporttrace.SpanExporter = &Exporter{}
)

// NewExporter creates an Exporter with the passed options.
func NewExporter(options ...Option) (*Exporter, error) {
	config, err := NewConfig(options...)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		traceExporter:  traceExporter{config: config},
		metricExporter: metricExporter{config},
	}, nil
}

// NewExportPipeline creates a complete export pipeline with the default
// selectors, processors, and trace registration. It is the responsibility
// of the caller to stop the returned push Controller.
func NewExportPipeline(exportOpts []Option, pushOpts []controller.Option) (trace.TracerProvider, *controller.Controller, error) {
	exporter, err := NewExporter(exportOpts...)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	pusher := controller.New(
		processor.New(
			simple.NewWithInexpensiveDistribution(),
			exporter,
		),
		append(
			pushOpts,
			controller.WithExporter(exporter),
		)...,
	)
	err = pusher.Start(context.Background())

	return tp, pusher, err
}

// InstallNewPipeline creates a complete export pipelines with defaults and
// registers it globally. It is the responsibility of the caller to stop the
// returned push Controller.
//
// Typically this is called as:
//
// 	pipeline, err := stdout.InstallNewPipeline(stdout.Config{...})
// 	if err != nil {
// 		...
// 	}
// 	defer pipeline.Stop()
// 	... Done
func InstallNewPipeline(exportOpts []Option, pushOpts []controller.Option) (*controller.Controller, error) {
	tracerProvider, controller, err := NewExportPipeline(exportOpts, pushOpts)
	if err != nil {
		return controller, err
	}
	otel.SetTracerProvider(tracerProvider)
	global.SetMeterProvider(controller.MeterProvider())
	return controller, err
}
