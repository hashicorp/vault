package nethttplibrary

import (
	nethttp "net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Pipeline contract for middleware infrastructure
type Pipeline interface {
	// Next moves the request object through middlewares in the pipeline
	Next(req *nethttp.Request, middlewareIndex int) (*nethttp.Response, error)
}

// custom transport for net/http with a middleware pipeline
type customTransport struct {
	// middleware pipeline in use for the client
	middlewarePipeline *middlewarePipeline
}

// middleware pipeline implementation using a roundtripper from net/http
type middlewarePipeline struct {
	// the round tripper to use to execute the request
	transport nethttp.RoundTripper
	// the middlewares to execute
	middlewares []Middleware
}

func newMiddlewarePipeline(middlewares []Middleware, transport nethttp.RoundTripper) *middlewarePipeline {
	return &middlewarePipeline{
		transport:   transport,
		middlewares: middlewares,
	}
}

// Next moves the request object through middlewares in the pipeline
func (pipeline *middlewarePipeline) Next(req *nethttp.Request, middlewareIndex int) (*nethttp.Response, error) {
	if middlewareIndex < len(pipeline.middlewares) {
		middleware := pipeline.middlewares[middlewareIndex]
		return middleware.Intercept(pipeline, middlewareIndex+1, req)
	}
	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	var observabilityName string
	if obsOptions != nil {
		observabilityName = obsOptions.GetTracerInstrumentationName()
		ctx, span = otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "request_transport")
		defer span.End()
		req = req.WithContext(ctx)
	}
	return pipeline.transport.RoundTrip(req)
}

// RoundTrip executes the the next middleware and returns a response
func (transport *customTransport) RoundTrip(req *nethttp.Request) (*nethttp.Response, error) {
	return transport.middlewarePipeline.Next(req, 0)
}

// GetDefaultTransport returns the default http transport used by the library
func GetDefaultTransport() nethttp.RoundTripper {
	defaultTransport, ok := nethttp.DefaultTransport.(*nethttp.Transport)
	if !ok {
		return nethttp.DefaultTransport
	}
	defaultTransport = defaultTransport.Clone()
	defaultTransport.ForceAttemptHTTP2 = true
	defaultTransport.DisableCompression = false
	return defaultTransport
}

// NewCustomTransport creates a new custom transport for http client with the provided set of middleware
func NewCustomTransport(middlewares ...Middleware) *customTransport {
	return NewCustomTransportWithParentTransport(nil, middlewares...)
}

// NewCustomTransportWithParentTransport creates a new custom transport which relies on the provided transport for http client with the provided set of middleware
func NewCustomTransportWithParentTransport(parentTransport nethttp.RoundTripper, middlewares ...Middleware) *customTransport {
	if len(middlewares) == 0 {
		middlewares = GetDefaultMiddlewares()
	}
	if parentTransport == nil {
		parentTransport = GetDefaultTransport()
	}
	return &customTransport{
		middlewarePipeline: newMiddlewarePipeline(middlewares, parentTransport),
	}
}
