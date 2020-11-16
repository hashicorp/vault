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

package trace // import "go.opentelemetry.io/otel/sdk/export/trace"

import (
	"context"
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// SpanExporter handles the delivery of SpanData to external receivers. This is
// the final component in the trace export pipeline.
type SpanExporter interface {
	// ExportSpans exports a batch of SpanData.
	//
	// This function is called synchronously, so there is no concurrency
	// safety requirement. However, due to the synchronous calling pattern,
	// it is critical that all timeouts and cancellations contained in the
	// passed context must be honored.
	//
	// Any retry logic must be contained in this function. The SDK that
	// calls this function will not implement any retry logic. All errors
	// returned by this function are considered unrecoverable and will be
	// reported to a configured error Handler.
	ExportSpans(ctx context.Context, spanData []*SpanData) error
	// Shutdown notifies the exporter of a pending halt to operations. The
	// exporter is expected to preform any cleanup or synchronization it
	// requires while honoring all timeouts and cancellations contained in
	// the passed context.
	Shutdown(ctx context.Context) error
}

// SpanData contains all the information collected by a completed span.
type SpanData struct {
	SpanContext  apitrace.SpanContext
	ParentSpanID apitrace.SpanID
	SpanKind     apitrace.SpanKind
	Name         string
	StartTime    time.Time
	// The wall clock time of EndTime will be adjusted to always be offset
	// from StartTime by the duration of the span.
	EndTime                  time.Time
	Attributes               []label.KeyValue
	MessageEvents            []Event
	Links                    []apitrace.Link
	StatusCode               codes.Code
	StatusMessage            string
	HasRemoteParent          bool
	DroppedAttributeCount    int
	DroppedMessageEventCount int
	DroppedLinkCount         int

	// ChildSpanCount holds the number of child span created for this span.
	ChildSpanCount int

	// Resource contains attributes representing an entity that produced this span.
	Resource *resource.Resource

	// InstrumentationLibrary defines the instrumentation library used to
	// provide instrumentation.
	InstrumentationLibrary instrumentation.Library
}

// Event is thing that happened during a Span's lifetime.
type Event struct {
	// Name is the name of this event
	Name string

	// Attributes describe the aspects of the event.
	Attributes []label.KeyValue

	// Time is the time at which this event was recorded.
	Time time.Time
}
