package nethttplibrary

import (
	nethttp "net/http"

	abs "github.com/microsoft/kiota-abstractions-go"
)

// ObservabilityOptions holds the tracing, metrics and logging configuration for the request adapter
type ObservabilityOptions struct {
	// Whether to include attributes which could contains EUII information like URLs
	IncludeEUIIAttributes bool
}

// GetTracerInstrumentationName returns the observability name to use for the tracer
func (o *ObservabilityOptions) GetTracerInstrumentationName() string {
	return "github.com/microsoft/kiota-http-go"
}

// GetIncludeEUIIAttributes returns whether to include attributes which could contains EUII information
func (o *ObservabilityOptions) GetIncludeEUIIAttributes() bool {
	return o.IncludeEUIIAttributes
}

// SetIncludeEUIIAttributes set whether to include attributes which could contains EUII information
func (o *ObservabilityOptions) SetIncludeEUIIAttributes(value bool) {
	o.IncludeEUIIAttributes = value
}

// ObservabilityOptionsInt defines the options contract for handlers
type ObservabilityOptionsInt interface {
	abs.RequestOption
	GetTracerInstrumentationName() string
	GetIncludeEUIIAttributes() bool
	SetIncludeEUIIAttributes(value bool)
}

func (*ObservabilityOptions) GetKey() abs.RequestOptionKey {
	return observabilityOptionsKeyValue
}

var observabilityOptionsKeyValue = abs.RequestOptionKey{
	Key: "ObservabilityOptions",
}

// GetObservabilityOptionsFromRequest returns the observability options from the request context
func GetObservabilityOptionsFromRequest(req *nethttp.Request) ObservabilityOptionsInt {
	if options, ok := req.Context().Value(observabilityOptionsKeyValue).(ObservabilityOptionsInt); ok {
		return options
	}
	return nil
}
