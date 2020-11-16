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
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

// TracerProvider provides access to instrumentation Tracers.
type TracerProvider interface {
	// Tracer creates an implementation of the Tracer interface.
	// The instrumentationName must be the name of the library providing
	// instrumentation. This name may be the same as the instrumented code
	// only if that code provides built-in instrumentation. If the
	// instrumentationName is empty, then a implementation defined default
	// name will be used instead.
	Tracer(instrumentationName string, opts ...TracerOption) Tracer
}

// TracerConfig is a group of options for a Tracer.
//
// Most users will use the tracer options instead.
type TracerConfig struct {
	// InstrumentationVersion is the version of the instrumentation library.
	InstrumentationVersion string
}

// NewTracerConfig applies all the options to a returned TracerConfig.
// The default value for all the fields of the returned TracerConfig are the
// default zero value of the type. Also, this does not perform any validation
// on the returned TracerConfig (e.g. no uniqueness checking or bounding of
// data), instead it is left to the implementations of the SDK to perform this
// action.
func NewTracerConfig(opts ...TracerOption) *TracerConfig {
	config := new(TracerConfig)
	for _, option := range opts {
		option.Apply(config)
	}
	return config
}

// TracerOption applies an options to a TracerConfig.
type TracerOption interface {
	Apply(*TracerConfig)
}

type instVersionTracerOption string

func (o instVersionTracerOption) Apply(c *TracerConfig) { c.InstrumentationVersion = string(o) }

// WithInstrumentationVersion sets the instrumentation version for a Tracer.
func WithInstrumentationVersion(version string) TracerOption {
	return instVersionTracerOption(version)
}

type Tracer interface {
	// Start a span.
	Start(ctx context.Context, spanName string, opts ...SpanOption) (context.Context, Span)
}

// ErrorConfig provides options to set properties of an error
// event at the time it is recorded.
//
// Most users will use the error options instead.
type ErrorConfig struct {
	Timestamp  time.Time
	StatusCode codes.Code
}

// ErrorOption applies changes to ErrorConfig that sets options when an error event is recorded.
type ErrorOption func(*ErrorConfig)

// WithErrorTime sets the time at which the error event should be recorded.
func WithErrorTime(t time.Time) ErrorOption {
	return func(c *ErrorConfig) {
		c.Timestamp = t
	}
}

// WithErrorStatus indicates the span status that should be set when recording an error event.
func WithErrorStatus(s codes.Code) ErrorOption {
	return func(c *ErrorConfig) {
		c.StatusCode = s
	}
}

type Span interface {
	// Tracer returns tracer used to create this span. Tracer cannot be nil.
	Tracer() Tracer

	// End completes the span. No updates are allowed to span after it
	// ends. The only exception is setting status of the span.
	End(options ...SpanOption)

	// AddEvent adds an event to the span.
	AddEvent(ctx context.Context, name string, attrs ...label.KeyValue)
	// AddEventWithTimestamp adds an event with a custom timestamp
	// to the span.
	AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue)

	// IsRecording returns true if the span is active and recording events is enabled.
	IsRecording() bool

	// RecordError records an error as a span event.
	RecordError(ctx context.Context, err error, opts ...ErrorOption)

	// SpanContext returns span context of the span. Returned SpanContext is usable
	// even after the span ends.
	SpanContext() SpanContext

	// SetStatus sets the status of the span in the form of a code
	// and a message.  SetStatus overrides the value of previous
	// calls to SetStatus on the Span.
	//
	// The default span status is OK, so it is not necessary to
	// explicitly set an OK status on successful Spans unless it
	// is to add an OK message or to override a previous status on the Span.
	SetStatus(code codes.Code, msg string)

	// SetName sets the name of the span.
	SetName(name string)

	// Set span attributes
	SetAttributes(kv ...label.KeyValue)
}

// SpanConfig is a group of options for a Span.
//
// Most users will use span options instead.
type SpanConfig struct {
	// Attributes describe the associated qualities of a Span.
	Attributes []label.KeyValue
	// Timestamp is a time in a Span life-cycle.
	Timestamp time.Time
	// Links are the associations a Span has with other Spans.
	Links []Link
	// Record is the recording state of a Span.
	Record bool
	// NewRoot identifies a Span as the root Span for a new trace. This is
	// commonly used when an existing trace crosses trust boundaries and the
	// remote parent span context should be ignored for security.
	NewRoot bool
	// SpanKind is the role a Span has in a trace.
	SpanKind SpanKind
}

// NewSpanConfig applies all the options to a returned SpanConfig.
// The default value for all the fields of the returned SpanConfig are the
// default zero value of the type. Also, this does not perform any validation
// on the returned SpanConfig (e.g. no uniqueness checking or bounding of
// data). Instead, it is left to the implementations of the SDK to perform this
// action.
func NewSpanConfig(opts ...SpanOption) *SpanConfig {
	c := new(SpanConfig)
	for _, option := range opts {
		option.Apply(c)
	}
	return c
}

// SpanOption applies an option to a SpanConfig.
type SpanOption interface {
	Apply(*SpanConfig)
}

type attributeSpanOption []label.KeyValue

func (o attributeSpanOption) Apply(c *SpanConfig) {
	c.Attributes = append(c.Attributes, []label.KeyValue(o)...)
}

// WithAttributes adds the attributes to a span. These attributes are meant to
// provide additional information about the work the Span represents. The
// attributes are added to the existing Span attributes, i.e. this does not
// overwrite.
func WithAttributes(attributes ...label.KeyValue) SpanOption {
	return attributeSpanOption(attributes)
}

type timestampSpanOption time.Time

func (o timestampSpanOption) Apply(c *SpanConfig) { c.Timestamp = time.Time(o) }

// WithTimestamp sets the time of a Span life-cycle moment (e.g. started or
// stopped).
func WithTimestamp(t time.Time) SpanOption {
	return timestampSpanOption(t)
}

type linksSpanOption []Link

func (o linksSpanOption) Apply(c *SpanConfig) { c.Links = append(c.Links, []Link(o)...) }

// WithLinks adds links to a Span. The links are added to the existing Span
// links, i.e. this does not overwrite.
func WithLinks(links ...Link) SpanOption {
	return linksSpanOption(links)
}

type recordSpanOption bool

func (o recordSpanOption) Apply(c *SpanConfig) { c.Record = bool(o) }

// WithRecord specifies that the span should be recorded. It is important to
// note that implementations may override this option, i.e. if the span is a
// child of an un-sampled trace.
func WithRecord() SpanOption {
	return recordSpanOption(true)
}

type newRootSpanOption bool

func (o newRootSpanOption) Apply(c *SpanConfig) { c.NewRoot = bool(o) }

// WithNewRoot specifies that the Span should be treated as a root Span. Any
// existing parent span context will be ignored when defining the Span's trace
// identifiers.
func WithNewRoot() SpanOption {
	return newRootSpanOption(true)
}

type spanKindSpanOption SpanKind

func (o spanKindSpanOption) Apply(c *SpanConfig) { c.SpanKind = SpanKind(o) }

// WithSpanKind sets the SpanKind of a Span.
func WithSpanKind(kind SpanKind) SpanOption {
	return spanKindSpanOption(kind)
}

// Link is used to establish relationship between two spans within the same Trace or
// across different Traces. Few examples of Link usage.
//   1. Batch Processing: A batch of elements may contain elements associated with one
//      or more traces/spans. Since there can only be one parent SpanContext, Link is
//      used to keep reference to SpanContext of all elements in the batch.
//   2. Public Endpoint: A SpanContext in incoming client request on a public endpoint
//      is untrusted from service provider perspective. In such case it is advisable to
//      start a new trace with appropriate sampling decision.
//      However, it is desirable to associate incoming SpanContext to new trace initiated
//      on service provider side so two traces (from Client and from Service Provider) can
//      be correlated.
type Link struct {
	SpanContext
	Attributes []label.KeyValue
}

// SpanKind represents the role of a Span inside a Trace. Often, this defines how a Span
// will be processed and visualized by various backends.
type SpanKind int

const (
	// As a convenience, these match the proto definition, see
	// opentelemetry/proto/trace/v1/trace.proto
	//
	// The unspecified value is not a valid `SpanKind`.  Use
	// `ValidateSpanKind()` to coerce a span kind to a valid
	// value.
	SpanKindUnspecified SpanKind = 0
	SpanKindInternal    SpanKind = 1
	SpanKindServer      SpanKind = 2
	SpanKindClient      SpanKind = 3
	SpanKindProducer    SpanKind = 4
	SpanKindConsumer    SpanKind = 5
)

// ValidateSpanKind returns a valid span kind value.  This will coerce
// invalid values into the default value, SpanKindInternal.
func ValidateSpanKind(spanKind SpanKind) SpanKind {
	switch spanKind {
	case SpanKindInternal,
		SpanKindServer,
		SpanKindClient,
		SpanKindProducer,
		SpanKindConsumer:
		// valid
		return spanKind
	default:
		return SpanKindInternal
	}
}

// String returns the specified name of the SpanKind in lower-case.
func (sk SpanKind) String() string {
	switch sk {
	case SpanKindInternal:
		return "internal"
	case SpanKindServer:
		return "server"
	case SpanKindClient:
		return "client"
	case SpanKindProducer:
		return "producer"
	case SpanKindConsumer:
		return "consumer"
	default:
		return "unspecified"
	}
}
