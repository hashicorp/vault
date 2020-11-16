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

type noopSpan struct {
}

var _ Span = noopSpan{}

// SpanContext returns an invalid span context.
func (noopSpan) SpanContext() SpanContext {
	return EmptySpanContext()
}

// IsRecording always returns false for NoopSpan.
func (noopSpan) IsRecording() bool {
	return false
}

// SetStatus does nothing.
func (noopSpan) SetStatus(status codes.Code, msg string) {
}

// SetError does nothing.
func (noopSpan) SetError(v bool) {
}

// SetAttributes does nothing.
func (noopSpan) SetAttributes(attributes ...label.KeyValue) {
}

// End does nothing.
func (noopSpan) End(options ...SpanOption) {
}

// RecordError does nothing.
func (noopSpan) RecordError(ctx context.Context, err error, opts ...ErrorOption) {
}

// Tracer returns noop implementation of Tracer.
func (noopSpan) Tracer() Tracer {
	return noopTracer{}
}

// AddEvent does nothing.
func (noopSpan) AddEvent(ctx context.Context, name string, attrs ...label.KeyValue) {
}

// AddEventWithTimestamp does nothing.
func (noopSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue) {
}

// SetName does nothing.
func (noopSpan) SetName(name string) {
}
