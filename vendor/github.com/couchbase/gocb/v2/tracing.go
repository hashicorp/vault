package gocb

import (
	"github.com/couchbase/gocbcore/v10"
	"time"
)

func tracerAddRef(tracer RequestTracer) {
	if tracer == nil {
		return
	}
	if refTracer, ok := tracer.(interface {
		AddRef() int32
	}); ok {
		refTracer.AddRef()
	}
}

func tracerDecRef(tracer RequestTracer) {
	if tracer == nil {
		return
	}
	if refTracer, ok := tracer.(interface {
		DecRef() int32
	}); ok {
		refTracer.DecRef()
	}
}

// RequestTracer describes the tracing abstraction in the SDK.
type RequestTracer interface {
	RequestSpan(parentContext RequestSpanContext, operationName string) RequestSpan
}

// RequestSpan is the interface for spans that are created by a RequestTracer.
type RequestSpan interface {
	End()
	Context() RequestSpanContext
	AddEvent(name string, timestamp time.Time)
	SetAttribute(key string, value interface{})
}

// RequestSpanContext is the interface for for external span contexts that can be passed in into the SDK option blocks.
type RequestSpanContext interface {
}

type coreRequestTracerWrapper struct {
	tracer RequestTracer
}

func (tracer *coreRequestTracerWrapper) RequestSpan(parentContext gocbcore.RequestSpanContext, operationName string) gocbcore.RequestSpan {
	return &coreRequestSpanWrapper{
		span: tracer.tracer.RequestSpan(parentContext, operationName),
	}
}

type coreRequestSpanWrapper struct {
	span RequestSpan
}

func (span *coreRequestSpanWrapper) End() {
	span.span.End()
}

func (span *coreRequestSpanWrapper) Context() gocbcore.RequestSpanContext {
	return span.span.Context()
}

func (span *coreRequestSpanWrapper) SetAttribute(key string, value interface{}) {
	span.span.SetAttribute(key, value)
}

func (span *coreRequestSpanWrapper) AddEvent(key string, timestamp time.Time) {
	span.span.SetAttribute(key, timestamp)
}

type noopSpan struct{}
type noopSpanContext struct{}

var (
	defaultNoopSpanContext = noopSpanContext{}
	defaultNoopSpan        = noopSpan{}
)

// NoopTracer is a RequestTracer implementation that does not perform any tracing.
type NoopTracer struct { // nolint: unused
}

// RequestSpan creates a new RequestSpan.
func (tracer *NoopTracer) RequestSpan(parentContext RequestSpanContext, operationName string) RequestSpan {
	return defaultNoopSpan
}

// End completes the span.
func (span noopSpan) End() {
}

// Context returns the RequestSpanContext for this span.
func (span noopSpan) Context() RequestSpanContext {
	return defaultNoopSpanContext
}

// SetAttribute adds an attribute to this span.
func (span noopSpan) SetAttribute(key string, value interface{}) {
}

// AddEvent adds an event to this span.
func (span noopSpan) AddEvent(key string, timestamp time.Time) {
}

func createSpan(tracer RequestTracer, parent RequestSpan, operationType, service string) RequestSpan {
	var tracectx RequestSpanContext
	if parent != nil {
		tracectx = parent.Context()
	}

	span := tracer.RequestSpan(tracectx, operationType)
	span.SetAttribute(spanAttribDBSystemKey, spanAttribDBSystemValue)
	if service != "" {
		span.SetAttribute(spanAttribServiceKey, service)
	}

	return span
}
