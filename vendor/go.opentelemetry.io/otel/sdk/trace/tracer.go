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

package trace

import (
	"context"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/internal/trace/parent"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type tracer struct {
	provider               *TracerProvider
	instrumentationLibrary instrumentation.Library
}

var _ apitrace.Tracer = &tracer{}

// Start starts a Span and returns it along with a context containing it.
//
// The Span is created with the provided name and as a child of any existing
// span context found in the passed context. The created Span will be
// configured appropriately by any SpanOption passed. Any Timestamp option
// passed will be used as the start time of the Span's life-cycle.
func (tr *tracer) Start(ctx context.Context, name string, options ...apitrace.SpanOption) (context.Context, apitrace.Span) {
	config := apitrace.NewSpanConfig(options...)

	parentSpanContext, remoteParent, links := parent.GetSpanContextAndLinks(ctx, config.NewRoot)

	if p := apitrace.SpanFromContext(ctx); p != nil {
		if sdkSpan, ok := p.(*span); ok {
			sdkSpan.addChild()
		}
	}

	span := startSpanInternal(tr, name, parentSpanContext, remoteParent, config)
	for _, l := range links {
		span.addLink(l)
	}
	for _, l := range config.Links {
		span.addLink(l)
	}
	span.SetAttributes(config.Attributes...)

	span.tracer = tr

	if span.IsRecording() {
		sps, _ := tr.provider.spanProcessors.Load().(spanProcessorStates)
		for _, sp := range sps {
			sp.sp.OnStart(span.data)
		}
	}

	ctx, end := startExecutionTracerTask(ctx, name)
	span.executionTracerTaskEnd = end
	return apitrace.ContextWithSpan(ctx, span), span
}
