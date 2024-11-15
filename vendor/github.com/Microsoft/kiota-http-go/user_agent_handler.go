package nethttplibrary

import (
	"fmt"
	nethttp "net/http"
	"strings"

	abs "github.com/microsoft/kiota-abstractions-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UserAgentHandler adds the product to the user agent header.
type UserAgentHandler struct {
	options UserAgentHandlerOptions
}

// NewUserAgentHandler creates a new user agent handler with the default options.
func NewUserAgentHandler() *UserAgentHandler {
	return NewUserAgentHandlerWithOptions(nil)
}

// NewUserAgentHandlerWithOptions creates a new user agent handler with the specified options.
func NewUserAgentHandlerWithOptions(options *UserAgentHandlerOptions) *UserAgentHandler {
	if options == nil {
		options = NewUserAgentHandlerOptions()
	}
	return &UserAgentHandler{
		options: *options,
	}
}

// UserAgentHandlerOptions to use when adding the product to the user agent header.
type UserAgentHandlerOptions struct {
	Enabled        bool
	ProductName    string
	ProductVersion string
}

// NewUserAgentHandlerOptions creates a new user agent handler options with the default values.
func NewUserAgentHandlerOptions() *UserAgentHandlerOptions {
	return &UserAgentHandlerOptions{
		Enabled:        true,
		ProductName:    "kiota-go",
		ProductVersion: "1.4.5",
	}
}

var userAgentKeyValue = abs.RequestOptionKey{
	Key: "UserAgentHandler",
}

type userAgentHandlerOptionsInt interface {
	abs.RequestOption
	GetEnabled() bool
	GetProductName() string
	GetProductVersion() string
}

// GetKey returns the key value to be used when the option is added to the request context
func (options *UserAgentHandlerOptions) GetKey() abs.RequestOptionKey {
	return userAgentKeyValue
}

// GetEnabled returns the value of the enabled property
func (options *UserAgentHandlerOptions) GetEnabled() bool {
	return options.Enabled
}

// GetProductName returns the value of the product name property
func (options *UserAgentHandlerOptions) GetProductName() string {
	return options.ProductName
}

// GetProductVersion returns the value of the product version property
func (options *UserAgentHandlerOptions) GetProductVersion() string {
	return options.ProductVersion
}

const userAgentHeaderKey = "User-Agent"

func (middleware UserAgentHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	obsOptions := GetObservabilityOptionsFromRequest(req)
	if obsOptions != nil {
		observabilityName := obsOptions.GetTracerInstrumentationName()
		ctx := req.Context()
		ctx, span := otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "UserAgentHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.useragent.enable", true))
		defer span.End()
		req = req.WithContext(ctx)
	}
	options, ok := req.Context().Value(userAgentKeyValue).(userAgentHandlerOptionsInt)
	if !ok {
		options = &middleware.options
	}
	if options.GetEnabled() {
		additionalValue := fmt.Sprintf("%s/%s", options.GetProductName(), options.GetProductVersion())
		currentValue := req.Header.Get(userAgentHeaderKey)
		if currentValue == "" {
			req.Header.Set(userAgentHeaderKey, additionalValue)
		} else if !strings.Contains(currentValue, additionalValue) {
			req.Header.Set(userAgentHeaderKey, fmt.Sprintf("%s %s", currentValue, additionalValue))
		}
	}
	return pipeline.Next(req, middlewareIndex)
}
