package gocb

import (
	"github.com/couchbase/gocbcore/v9"
)

func tracerAddRef(tracer requestTracer) {
	if tracer == nil {
		return
	}
	if refTracer, ok := tracer.(interface {
		AddRef() int32
	}); ok {
		refTracer.AddRef()
	}
}

func tracerDecRef(tracer requestTracer) {
	if tracer == nil {
		return
	}
	if refTracer, ok := tracer.(interface {
		DecRef() int32
	}); ok {
		refTracer.DecRef()
	}
}

// requestTracer describes the tracing abstraction in the SDK.
type requestTracer interface {
	StartSpan(operationName string, parentContext requestSpanContext) requestSpan
}

// requestSpan is the interface for spans that are created by a requestTracer.
type requestSpan interface {
	Finish()
	Context() requestSpanContext
	SetTag(key string, value interface{}) requestSpan
}

// requestSpanContext is the interface for for external span contexts that can be passed in into the SDK option blocks.
type requestSpanContext interface {
}

type requestTracerWrapper struct {
	tracer requestTracer
}

func (tracer *requestTracerWrapper) StartSpan(operationName string, parentContext gocbcore.RequestSpanContext) gocbcore.RequestSpan {
	return requestSpanWrapper{
		span: tracer.tracer.StartSpan(operationName, parentContext),
	}
}

type requestSpanWrapper struct {
	span requestSpan
}

func (span requestSpanWrapper) Finish() {
	span.span.Finish()
}

func (span requestSpanWrapper) Context() gocbcore.RequestSpanContext {
	return span.span.Context()
}

func (span requestSpanWrapper) SetTag(key string, value interface{}) gocbcore.RequestSpan {
	span.span = span.span.SetTag(key, value)
	return span
}

type noopSpan struct{}
type noopSpanContext struct{}

var (
	defaultNoopSpanContext = noopSpanContext{}
	defaultNoopSpan        = noopSpan{}
)

// noopTracer will have a future use so we tell the linter not to flag it.
type noopTracer struct { // nolint: unused
}

func (tracer *noopTracer) StartSpan(operationName string, parentContext requestSpanContext) requestSpan {
	return defaultNoopSpan
}

func (span noopSpan) Finish() {
}

func (span noopSpan) Context() requestSpanContext {
	return defaultNoopSpanContext
}

func (span noopSpan) SetTag(key string, value interface{}) requestSpan {
	return defaultNoopSpan
}
