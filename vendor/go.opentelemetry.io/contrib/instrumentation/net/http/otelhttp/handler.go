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

package otelhttp

import (
	"io"
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
)

var _ http.Handler = &Handler{}

// Handler is http middleware that corresponds to the http.Handler interface and
// is designed to wrap a http.Mux (or equivalent), while individual routes on
// the mux are wrapped with WithRouteTag. A Handler will add various attributes
// to the span using the label.Keys defined in this package.
type Handler struct {
	operation string
	handler   http.Handler

	tracer            trace.Tracer
	meter             metric.Meter
	propagators       otel.TextMapPropagator
	spanStartOptions  []trace.SpanOption
	readEvent         bool
	writeEvent        bool
	filters           []Filter
	spanNameFormatter func(string, *http.Request) string
	counters          map[string]metric.Int64Counter
	valueRecorders    map[string]metric.Int64ValueRecorder
}

func defaultHandlerFormatter(operation string, _ *http.Request) string {
	return operation
}

// NewHandler wraps the passed handler, functioning like middleware, in a span
// named after the operation and with any provided Options.
func NewHandler(handler http.Handler, operation string, opts ...Option) http.Handler {
	h := Handler{
		handler:   handler,
		operation: operation,
	}

	defaultOpts := []Option{
		WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
		WithSpanNameFormatter(defaultHandlerFormatter),
	}

	c := newConfig(append(defaultOpts, opts...)...)
	h.configure(c)
	h.createMeasures()

	return &h
}

func (h *Handler) configure(c *config) {
	h.tracer = c.Tracer
	h.meter = c.Meter
	h.propagators = c.Propagators
	h.spanStartOptions = c.SpanStartOptions
	h.readEvent = c.ReadEvent
	h.writeEvent = c.WriteEvent
	h.filters = c.Filters
	h.spanNameFormatter = c.SpanNameFormatter
}

func handleErr(err error) {
	if err != nil {
		global.Handle(err)
	}
}

func (h *Handler) createMeasures() {
	h.counters = make(map[string]metric.Int64Counter)
	h.valueRecorders = make(map[string]metric.Int64ValueRecorder)

	requestBytesCounter, err := h.meter.NewInt64Counter(RequestContentLength)
	handleErr(err)

	responseBytesCounter, err := h.meter.NewInt64Counter(ResponseContentLength)
	handleErr(err)

	serverLatencyMeasure, err := h.meter.NewInt64ValueRecorder(ServerLatency)
	handleErr(err)

	h.counters[RequestContentLength] = requestBytesCounter
	h.counters[ResponseContentLength] = responseBytesCounter
	h.valueRecorders[ServerLatency] = serverLatencyMeasure
}

// ServeHTTP serves HTTP requests (http.Handler)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestStartTime := time.Now()
	for _, f := range h.filters {
		if !f(r) {
			// Simply pass through to the handler if a filter rejects the request
			h.handler.ServeHTTP(w, r)
			return
		}
	}

	opts := append([]trace.SpanOption{
		trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
		trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
		trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(h.operation, "", r)...),
	}, h.spanStartOptions...) // start with the configured options

	ctx := h.propagators.Extract(r.Context(), r.Header)
	ctx, span := h.tracer.Start(ctx, h.spanNameFormatter(h.operation, r), opts...)
	defer span.End()

	readRecordFunc := func(int64) {}
	if h.readEvent {
		readRecordFunc = func(n int64) {
			span.AddEvent(ctx, "read", ReadBytesKey.Int64(n))
		}
	}
	bw := bodyWrapper{ReadCloser: r.Body, record: readRecordFunc}
	r.Body = &bw

	writeRecordFunc := func(int64) {}
	if h.writeEvent {
		writeRecordFunc = func(n int64) {
			span.AddEvent(ctx, "write", WroteBytesKey.Int64(n))
		}
	}

	rww := &respWriterWrapper{ResponseWriter: w, record: writeRecordFunc, ctx: ctx, props: h.propagators}

	// Wrap w to use our ResponseWriter methods while also exposing
	// other interfaces that w may implement (http.CloseNotifier,
	// http.Flusher, http.Hijacker, http.Pusher, io.ReaderFrom).

	w = httpsnoop.Wrap(w, httpsnoop.Hooks{
		Header: func(httpsnoop.HeaderFunc) httpsnoop.HeaderFunc {
			return rww.Header
		},
		Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return rww.Write
		},
		WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return rww.WriteHeader
		},
	})

	labeler := &Labeler{}
	ctx = injectLabeler(ctx, labeler)

	h.handler.ServeHTTP(w, r.WithContext(ctx))

	setAfterServeAttributes(span, bw.read, rww.written, rww.statusCode, bw.err, rww.err)

	// Add request metrics

	labels := append(labeler.Get(), semconv.HTTPServerMetricAttributesFromHTTPRequest(h.operation, r)...)

	h.counters[RequestContentLength].Add(ctx, bw.read, labels...)
	h.counters[ResponseContentLength].Add(ctx, rww.written, labels...)

	elapsedTime := time.Since(requestStartTime).Microseconds()

	h.valueRecorders[ServerLatency].Record(ctx, elapsedTime, labels...)
}

func setAfterServeAttributes(span trace.Span, read, wrote int64, statusCode int, rerr, werr error) {
	labels := []label.KeyValue{}

	// TODO: Consider adding an event after each read and write, possibly as an
	// option (defaulting to off), so as to not create needlessly verbose spans.
	if read > 0 {
		labels = append(labels, ReadBytesKey.Int64(read))
	}
	if rerr != nil && rerr != io.EOF {
		labels = append(labels, ReadErrorKey.String(rerr.Error()))
	}
	if wrote > 0 {
		labels = append(labels, WroteBytesKey.Int64(wrote))
	}
	if statusCode > 0 {
		labels = append(labels, semconv.HTTPAttributesFromHTTPStatusCode(statusCode)...)
		span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(statusCode))
	}
	if werr != nil && werr != io.EOF {
		labels = append(labels, WriteErrorKey.String(werr.Error()))
	}
	span.SetAttributes(labels...)
}

// WithRouteTag annotates a span with the provided route name using the
// RouteKey Tag.
func WithRouteTag(route string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		span.SetAttributes(semconv.HTTPRouteKey.String(route))
		h.ServeHTTP(w, r)
	})
}
