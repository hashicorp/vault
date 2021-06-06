package gocbcore

import (
	"fmt"
)

// RequestTracer describes the tracing abstraction in the SDK.
type RequestTracer interface {
	StartSpan(operationName string, parentContext RequestSpanContext) RequestSpan
}

// RequestSpan is the interface for spans that are created by a RequestTracer.
type RequestSpan interface {
	Finish()
	Context() RequestSpanContext
	SetTag(key string, value interface{}) RequestSpan
}

// RequestSpanContext is the interface for for external span contexts that can be passed in into the SDK option blocks.
type RequestSpanContext interface {
}

type noopSpan struct{}
type noopSpanContext struct{}

var (
	defaultNoopSpanContext = noopSpanContext{}
	defaultNoopSpan        = noopSpan{}
)

type noopTracer struct {
}

func (tracer noopTracer) StartSpan(operationName string, parentContext RequestSpanContext) RequestSpan {
	return defaultNoopSpan
}

func (span noopSpan) Finish() {
}

func (span noopSpan) Context() RequestSpanContext {
	return defaultNoopSpanContext
}

func (span noopSpan) SetTag(key string, value interface{}) RequestSpan {
	return defaultNoopSpan
}

type opTracer struct {
	parentContext RequestSpanContext
	opSpan        RequestSpan
}

func (tracer *opTracer) Finish() {
	if tracer.opSpan != nil {
		tracer.opSpan.Finish()
	}
}

func (tracer *opTracer) RootContext() RequestSpanContext {
	if tracer.opSpan != nil {
		return tracer.opSpan.Context()
	}

	return tracer.parentContext
}

type tracerManager interface {
	CreateOpTrace(operationName string, parentContext RequestSpanContext) *opTracer
	StartHTTPSpan(req *httpRequest, name string) RequestSpan
	StartCmdTrace(req *memdQRequest)
	StartNetTrace(req *memdQRequest)
}

type tracerComponent struct {
	tracer           RequestTracer
	bucket           string
	noRootTraceSpans bool
}

func newTracerComponent(tracer RequestTracer, bucket string, noRootTraceSpans bool) *tracerComponent {
	return &tracerComponent{
		tracer:           tracer,
		bucket:           bucket,
		noRootTraceSpans: noRootTraceSpans,
	}
}

func (tc *tracerComponent) CreateOpTrace(operationName string, parentContext RequestSpanContext) *opTracer {
	if tc.noRootTraceSpans {
		return &opTracer{
			parentContext: parentContext,
			opSpan:        nil,
		}
	}

	opSpan := tc.tracer.StartSpan(operationName, parentContext).
		SetTag("component", "couchbase-go-sdk").
		SetTag("db.instance", tc.bucket).
		SetTag("span.kind", "client")

	return &opTracer{
		parentContext: parentContext,
		opSpan:        opSpan,
	}
}

func (tc *tracerComponent) StartHTTPSpan(req *httpRequest, name string) RequestSpan {
	return tc.tracer.StartSpan(name, req.RootTraceContext).
		SetTag("retry", req.RetryAttempts())
}

func (tc *tracerComponent) StartCmdTrace(req *memdQRequest) {
	if req.cmdTraceSpan != nil {
		logWarnf("Attempted to start tracing on traced request")
		return
	}

	if req.RootTraceContext == nil {
		return
	}

	req.processingLock.Lock()
	req.cmdTraceSpan = tc.tracer.StartSpan(req.Packet.Command.Name(), req.RootTraceContext).
		SetTag("retry", req.RetryAttempts())

	req.processingLock.Unlock()
}

func (tc *tracerComponent) StartNetTrace(req *memdQRequest) {
	if req.cmdTraceSpan == nil {
		return
	}

	if req.netTraceSpan != nil {
		logWarnf("Attempted to start net tracing on traced request")
		return
	}

	req.processingLock.Lock()
	req.netTraceSpan = tc.tracer.StartSpan("rpc", req.cmdTraceSpan.Context()).
		SetTag("span.kind", "client")
	req.processingLock.Unlock()
}

func stopCmdTrace(req *memdQRequest) {
	if req.RootTraceContext == nil {
		return
	}

	if req.cmdTraceSpan == nil {
		logWarnf("Attempted to stop tracing on untraced request")
		return
	}

	req.cmdTraceSpan.Finish()
	req.cmdTraceSpan = nil
}

func cancelReqTrace(req *memdQRequest) {
	if req.cmdTraceSpan != nil {
		if req.netTraceSpan != nil {
			req.netTraceSpan.Finish()
		}

		req.cmdTraceSpan.Finish()
	}
}

func stopNetTrace(req *memdQRequest, resp *memdQResponse, localAddress, remoteAddress string) {
	if req.cmdTraceSpan == nil {
		return
	}

	if req.netTraceSpan == nil {
		logWarnf("Attempted to stop net tracing on an untraced request")
		return
	}

	req.netTraceSpan.SetTag("couchbase.operation_id", fmt.Sprintf("0x%x", resp.Opaque))
	req.netTraceSpan.SetTag("couchbase.local_id", resp.sourceConnID)
	if isLogRedactionLevelNone() {
		req.netTraceSpan.SetTag("couchbase.document_key", string(req.Key))
	}
	req.netTraceSpan.SetTag("local.address", localAddress)
	req.netTraceSpan.SetTag("peer.address", remoteAddress)
	if resp.Packet.ServerDurationFrame != nil {
		req.netTraceSpan.SetTag("server_duration", resp.Packet.ServerDurationFrame.ServerDuration)
	}

	req.netTraceSpan.Finish()
	req.netTraceSpan = nil
}
