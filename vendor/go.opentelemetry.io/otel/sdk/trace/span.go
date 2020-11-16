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
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/internal"
)

const (
	errorTypeKey    = label.Key("error.type")
	errorMessageKey = label.Key("error.message")
	errorEventName  = "error"
)

// span implements apitrace.Span interface.
type span struct {
	// data contains information recorded about the span.
	//
	// It will be non-nil if we are exporting the span or recording events for it.
	// Otherwise, data is nil, and the span is simply a carrier for the
	// SpanContext, so that the trace ID is propagated.
	data        *export.SpanData
	mu          sync.Mutex // protects the contents of *data (but not the pointer value.)
	spanContext apitrace.SpanContext

	// attributes are capped at configured limit. When the capacity is reached an oldest entry
	// is removed to create room for a new entry.
	attributes *attributesMap

	// messageEvents are stored in FIFO queue capped by configured limit.
	messageEvents *evictedQueue

	// links are stored in FIFO queue capped by configured limit.
	links *evictedQueue

	// spanStore is the spanStore this span belongs to, if any, otherwise it is nil.
	//*spanStore
	endOnce sync.Once

	executionTracerTaskEnd func()  // ends the execution tracer span
	tracer                 *tracer // tracer used to create span.
}

var _ apitrace.Span = &span{}

func (s *span) SpanContext() apitrace.SpanContext {
	if s == nil {
		return apitrace.EmptySpanContext()
	}
	return s.spanContext
}

func (s *span) IsRecording() bool {
	if s == nil {
		return false
	}
	return s.data != nil
}

func (s *span) SetStatus(code codes.Code, msg string) {
	if s == nil {
		return
	}
	if !s.IsRecording() {
		return
	}
	s.mu.Lock()
	s.data.StatusCode = code
	s.data.StatusMessage = msg
	s.mu.Unlock()
}

func (s *span) SetAttributes(attributes ...label.KeyValue) {
	if !s.IsRecording() {
		return
	}
	s.copyToCappedAttributes(attributes...)
}

// End ends the span.
//
// The only SpanOption currently supported is WithTimestamp which will set the
// end time for a Span's life-cycle.
//
// If this method is called while panicking an error event is added to the
// Span before ending it and the panic is continued.
func (s *span) End(options ...apitrace.SpanOption) {
	if s == nil {
		return
	}

	if recovered := recover(); recovered != nil {
		// Record but don't stop the panic.
		defer panic(recovered)
		s.addEventWithTimestamp(
			time.Now(),
			errorEventName,
			errorTypeKey.String(typeStr(recovered)),
			errorMessageKey.String(fmt.Sprint(recovered)),
		)
	}

	if s.executionTracerTaskEnd != nil {
		s.executionTracerTaskEnd()
	}
	if !s.IsRecording() {
		return
	}
	config := apitrace.NewSpanConfig(options...)
	s.endOnce.Do(func() {
		sps, ok := s.tracer.provider.spanProcessors.Load().(spanProcessorStates)
		mustExportOrProcess := ok && len(sps) > 0
		if mustExportOrProcess {
			sd := s.makeSpanData()
			if config.Timestamp.IsZero() {
				sd.EndTime = internal.MonotonicEndTime(sd.StartTime)
			} else {
				sd.EndTime = config.Timestamp
			}
			for _, sp := range sps {
				sp.sp.OnEnd(sd)
			}
		}
	})
}

func (s *span) RecordError(ctx context.Context, err error, opts ...apitrace.ErrorOption) {
	if s == nil || err == nil {
		return
	}

	if !s.IsRecording() {
		return
	}

	cfg := apitrace.ErrorConfig{}

	for _, o := range opts {
		o(&cfg)
	}

	if cfg.Timestamp.IsZero() {
		cfg.Timestamp = time.Now()
	}

	if cfg.StatusCode != codes.Unset {
		s.SetStatus(cfg.StatusCode, "")
	}

	s.AddEventWithTimestamp(ctx, cfg.Timestamp, errorEventName,
		errorTypeKey.String(typeStr(err)),
		errorMessageKey.String(err.Error()),
	)
}

func typeStr(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.PkgPath() == "" && t.Name() == "" {
		// Likely a builtin type.
		return t.String()
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func (s *span) Tracer() apitrace.Tracer {
	return s.tracer
}

func (s *span) AddEvent(ctx context.Context, name string, attrs ...label.KeyValue) {
	if !s.IsRecording() {
		return
	}
	s.addEventWithTimestamp(time.Now(), name, attrs...)
}

func (s *span) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue) {
	if !s.IsRecording() {
		return
	}
	s.addEventWithTimestamp(timestamp, name, attrs...)
}

func (s *span) addEventWithTimestamp(timestamp time.Time, name string, attrs ...label.KeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageEvents.add(export.Event{
		Name:       name,
		Attributes: attrs,
		Time:       timestamp,
	})
}

var errUninitializedSpan = errors.New("failed to set name on uninitialized span")

func (s *span) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		global.Handle(errUninitializedSpan)
		return
	}
	s.data.Name = name
	// SAMPLING
	noParent := !s.data.ParentSpanID.IsValid()
	var ctx apitrace.SpanContext
	if noParent {
		ctx = apitrace.EmptySpanContext()
	} else {
		// FIXME: Where do we get the parent context from?
		// From SpanStore?
		ctx = s.data.SpanContext
	}
	data := samplingData{
		noParent:     noParent,
		remoteParent: s.data.HasRemoteParent,
		parent:       ctx,
		name:         name,
		cfg:          s.tracer.provider.config.Load().(*Config),
		span:         s,
		attributes:   s.data.Attributes,
		links:        s.data.Links,
		kind:         s.data.SpanKind,
	}
	sampled := makeSamplingDecision(data)

	// Adding attributes directly rather than using s.SetAttributes()
	// as s.mu is already locked and attempting to do so would deadlock.
	for _, a := range sampled.Attributes {
		s.attributes.add(a)
	}
}

func (s *span) addLink(link apitrace.Link) {
	if !s.IsRecording() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.links.add(link)
}

// makeSpanData produces a SpanData representing the current state of the span.
// It requires that s.data is non-nil.
func (s *span) makeSpanData() *export.SpanData {
	var sd export.SpanData
	s.mu.Lock()
	defer s.mu.Unlock()
	sd = *s.data

	s.attributes.toSpanData(&sd)

	if len(s.messageEvents.queue) > 0 {
		sd.MessageEvents = s.interfaceArrayToMessageEventArray()
		sd.DroppedMessageEventCount = s.messageEvents.droppedCount
	}
	if len(s.links.queue) > 0 {
		sd.Links = s.interfaceArrayToLinksArray()
		sd.DroppedLinkCount = s.links.droppedCount
	}
	return &sd
}

func (s *span) interfaceArrayToLinksArray() []apitrace.Link {
	linkArr := make([]apitrace.Link, 0)
	for _, value := range s.links.queue {
		linkArr = append(linkArr, value.(apitrace.Link))
	}
	return linkArr
}

func (s *span) interfaceArrayToMessageEventArray() []export.Event {
	messageEventArr := make([]export.Event, 0)
	for _, value := range s.messageEvents.queue {
		messageEventArr = append(messageEventArr, value.(export.Event))
	}
	return messageEventArr
}

func (s *span) copyToCappedAttributes(attributes ...label.KeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range attributes {
		if a.Value.Type() != label.INVALID {
			s.attributes.add(a)
		}
	}
}

func (s *span) addChild() {
	if !s.IsRecording() {
		return
	}
	s.mu.Lock()
	s.data.ChildSpanCount++
	s.mu.Unlock()
}

func startSpanInternal(tr *tracer, name string, parent apitrace.SpanContext, remoteParent bool, o *apitrace.SpanConfig) *span {
	var noParent bool
	span := &span{}
	span.spanContext = parent

	cfg := tr.provider.config.Load().(*Config)

	if parent == apitrace.EmptySpanContext() {
		span.spanContext.TraceID = cfg.IDGenerator.NewTraceID()
		noParent = true
	}
	span.spanContext.SpanID = cfg.IDGenerator.NewSpanID()
	data := samplingData{
		noParent:     noParent,
		remoteParent: remoteParent,
		parent:       parent,
		name:         name,
		cfg:          cfg,
		span:         span,
		attributes:   o.Attributes,
		links:        o.Links,
		kind:         o.SpanKind,
	}
	sampled := makeSamplingDecision(data)

	// TODO: [rghetia] restore when spanstore is added.
	// if !internal.LocalSpanStoreEnabled && !span.spanContext.IsSampled() && !o.Record {
	if !span.spanContext.IsSampled() && !o.Record {
		return span
	}

	startTime := o.Timestamp
	if startTime.IsZero() {
		startTime = time.Now()
	}
	span.data = &export.SpanData{
		SpanContext:            span.spanContext,
		StartTime:              startTime,
		SpanKind:               apitrace.ValidateSpanKind(o.SpanKind),
		Name:                   name,
		HasRemoteParent:        remoteParent,
		Resource:               cfg.Resource,
		InstrumentationLibrary: tr.instrumentationLibrary,
	}
	span.attributes = newAttributesMap(cfg.MaxAttributesPerSpan)
	span.messageEvents = newEvictedQueue(cfg.MaxEventsPerSpan)
	span.links = newEvictedQueue(cfg.MaxLinksPerSpan)

	span.SetAttributes(sampled.Attributes...)

	if !noParent {
		span.data.ParentSpanID = parent.SpanID
	}
	// TODO: [rghetia] restore when spanstore is added.
	//if internal.LocalSpanStoreEnabled {
	//	ss := spanStoreForNameCreateIfNew(name)
	//	if ss != nil {
	//		span.spanStore = ss
	//		ss.add(span)
	//	}
	//}

	return span
}

type samplingData struct {
	noParent     bool
	remoteParent bool
	parent       apitrace.SpanContext
	name         string
	cfg          *Config
	span         *span
	attributes   []label.KeyValue
	links        []apitrace.Link
	kind         apitrace.SpanKind
}

func makeSamplingDecision(data samplingData) SamplingResult {
	if data.noParent || data.remoteParent {
		// If this span is the child of a local span and no
		// Sampler is set in the options, keep the parent's
		// TraceFlags.
		//
		// Otherwise, consult the Sampler in the options if it
		// is non-nil, otherwise the default sampler.
		sampler := data.cfg.DefaultSampler
		//if o.Sampler != nil {
		//	sampler = o.Sampler
		//}
		spanContext := &data.span.spanContext
		sampled := sampler.ShouldSample(SamplingParameters{
			ParentContext:   data.parent,
			TraceID:         spanContext.TraceID,
			Name:            data.name,
			HasRemoteParent: data.remoteParent,
			Kind:            data.kind,
			Attributes:      data.attributes,
			Links:           data.links,
		})
		if sampled.Decision == RecordAndSample {
			spanContext.TraceFlags |= apitrace.FlagsSampled
		} else {
			spanContext.TraceFlags &^= apitrace.FlagsSampled
		}
		return sampled
	}
	if data.parent.TraceFlags&apitrace.FlagsSampled != 0 {
		return SamplingResult{Decision: RecordAndSample}
	}
	return SamplingResult{Decision: Drop}
}
