package nethttplibrary

import (
	nethttp "net/http"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HeadersInspectionHandlerOptions is the options to use when inspecting headers
type HeadersInspectionOptions struct {
	InspectRequestHeaders  bool
	InspectResponseHeaders bool
	RequestHeaders         *abstractions.RequestHeaders
	ResponseHeaders        *abstractions.ResponseHeaders
}

// NewHeadersInspectionOptions creates a new HeadersInspectionOptions with default options
func NewHeadersInspectionOptions() *HeadersInspectionOptions {
	return &HeadersInspectionOptions{
		RequestHeaders:  abstractions.NewRequestHeaders(),
		ResponseHeaders: abstractions.NewResponseHeaders(),
	}
}

type headersInspectionOptionsInt interface {
	abstractions.RequestOption
	GetInspectRequestHeaders() bool
	GetInspectResponseHeaders() bool
	GetRequestHeaders() *abstractions.RequestHeaders
	GetResponseHeaders() *abstractions.ResponseHeaders
}

var headersInspectionKeyValue = abstractions.RequestOptionKey{
	Key: "nethttplibrary.HeadersInspectionOptions",
}

// GetInspectRequestHeaders returns true if the request headers should be inspected
func (o *HeadersInspectionOptions) GetInspectRequestHeaders() bool {
	return o.InspectRequestHeaders
}

// GetInspectResponseHeaders returns true if the response headers should be inspected
func (o *HeadersInspectionOptions) GetInspectResponseHeaders() bool {
	return o.InspectResponseHeaders
}

// GetRequestHeaders returns the request headers
func (o *HeadersInspectionOptions) GetRequestHeaders() *abstractions.RequestHeaders {
	return o.RequestHeaders
}

// GetResponseHeaders returns the response headers
func (o *HeadersInspectionOptions) GetResponseHeaders() *abstractions.ResponseHeaders {
	return o.ResponseHeaders
}

// GetKey returns the key for the HeadersInspectionOptions
func (o *HeadersInspectionOptions) GetKey() abstractions.RequestOptionKey {
	return headersInspectionKeyValue
}

// HeadersInspectionHandler allows inspecting of the headers of the request and response via a request option
type HeadersInspectionHandler struct {
	options HeadersInspectionOptions
}

// NewHeadersInspectionHandler creates a new HeadersInspectionHandler with default options
func NewHeadersInspectionHandler() *HeadersInspectionHandler {
	return NewHeadersInspectionHandlerWithOptions(*NewHeadersInspectionOptions())
}

// NewHeadersInspectionHandlerWithOptions creates a new HeadersInspectionHandler with the given options
func NewHeadersInspectionHandlerWithOptions(options HeadersInspectionOptions) *HeadersInspectionHandler {
	return &HeadersInspectionHandler{options: options}
}

// Intercept implements the interface and evaluates whether to retry a failed request.
func (middleware HeadersInspectionHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	var observabilityName string
	if obsOptions != nil {
		observabilityName = obsOptions.GetTracerInstrumentationName()
		ctx, span = otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "HeadersInspectionHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.headersInspection.enable", true))
		defer span.End()
		req = req.WithContext(ctx)
	}
	reqOption, ok := req.Context().Value(headersInspectionKeyValue).(headersInspectionOptionsInt)
	if !ok {
		reqOption = &middleware.options
	}
	if reqOption.GetInspectRequestHeaders() {
		for k, v := range req.Header {
			if len(v) == 1 {
				reqOption.GetRequestHeaders().Add(k, v[0])
			} else {
				reqOption.GetRequestHeaders().Add(k, v[0], v[1:]...)
			}
		}
	}
	response, err := pipeline.Next(req, middlewareIndex)
	if reqOption.GetInspectResponseHeaders() {
		for k, v := range response.Header {
			if len(v) == 1 {
				reqOption.GetResponseHeaders().Add(k, v[0])
			} else {
				reqOption.GetResponseHeaders().Add(k, v[0], v[1:]...)
			}
		}
	}
	if err != nil {
		return response, err
	}
	return response, err
}
