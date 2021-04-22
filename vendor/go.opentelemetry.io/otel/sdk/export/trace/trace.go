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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// SpanExporter handles the delivery of SpanSnapshot structs to external
// receivers. This is the final component in the trace export pipeline.
type SpanExporter interface {
	// ExportSpans exports a batch of SpanSnapshots.
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
	ExportSpans(ctx context.Context, ss []*SpanSnapshot) error
	// Shutdown notifies the exporter of a pending halt to operations. The
	// exporter is expected to preform any cleanup or synchronization it
	// requires while honoring all timeouts and cancellations contained in
	// the passed context.
	Shutdown(ctx context.Context) error
}

// SpanSnapshot is a snapshot of a span which contains all the information
// collected by the span. Its main purpose is exporting completed spans.
// Although SpanSnapshot fields can be accessed and potentially modified,
// SpanSnapshot should be treated as immutable. Changes to the span from which
// the SpanSnapshot was created are NOT reflected in the SpanSnapshot.
type SpanSnapshot struct {
	SpanContext  trace.SpanContext
	ParentSpanID trace.SpanID
	SpanKind     trace.SpanKind
	Name         string
	StartTime    time.Time
	// The wall clock time of EndTime will be adjusted to always be offset
	// from StartTime by the duration of the span.
	EndTime         time.Time
	Attributes      []attribute.KeyValue
	MessageEvents   []trace.Event
	Links           []trace.Link
	StatusCode      codes.Code
	StatusMessage   string
	HasRemoteParent bool

	// DroppedAttributeCount contains dropped attributes for the span itself, events and links.
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
