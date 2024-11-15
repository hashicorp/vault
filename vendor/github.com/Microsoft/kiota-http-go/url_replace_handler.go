package nethttplibrary

import (
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

var urlReplaceOptionKey = abstractions.RequestOptionKey{Key: "UrlReplaceOptionKey"}

// UrlReplaceHandler is a middleware handler that replaces url segments in the uri path.
type UrlReplaceHandler struct {
	options UrlReplaceOptions
}

// NewUrlReplaceHandler creates a configuration object for the CompressionHandler
func NewUrlReplaceHandler(enabled bool, replacementPairs map[string]string) *UrlReplaceHandler {
	return &UrlReplaceHandler{UrlReplaceOptions{Enabled: enabled, ReplacementPairs: replacementPairs}}
}

// UrlReplaceOptions is a configuration object for the UrlReplaceHandler middleware
type UrlReplaceOptions struct {
	Enabled          bool
	ReplacementPairs map[string]string
}

// GetKey returns UrlReplaceOptions unique name in context object
func (u *UrlReplaceOptions) GetKey() abstractions.RequestOptionKey {
	return urlReplaceOptionKey
}

// GetReplacementPairs reads ReplacementPairs settings from UrlReplaceOptions
func (u *UrlReplaceOptions) GetReplacementPairs() map[string]string {
	return u.ReplacementPairs
}

// IsEnabled reads Enabled setting from UrlReplaceOptions
func (u *UrlReplaceOptions) IsEnabled() bool {
	return u.Enabled
}

type urlReplaceOptionsInt interface {
	abstractions.RequestOption
	IsEnabled() bool
	GetReplacementPairs() map[string]string
}

// Intercept is invoked by the middleware pipeline to either move the request/response
// to the next middleware in the pipeline
func (c *UrlReplaceHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *http.Request) (*http.Response, error) {
	reqOption, ok := req.Context().Value(urlReplaceOptionKey).(urlReplaceOptionsInt)
	if !ok {
		reqOption = &c.options
	}

	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	if obsOptions != nil {
		ctx, span = otel.GetTracerProvider().Tracer(obsOptions.GetTracerInstrumentationName()).Start(ctx, "UrlReplaceHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.url_replacer.enable", true))
		defer span.End()
		req = req.WithContext(ctx)
	}

	if !reqOption.IsEnabled() || len(reqOption.GetReplacementPairs()) == 0 {
		return pipeline.Next(req, middlewareIndex)
	}

	req.URL.Path = ReplacePathTokens(req.URL.Path, reqOption.GetReplacementPairs())

	if span != nil {
		span.SetAttributes(attribute.String("http.request_url", req.RequestURI))
	}

	return pipeline.Next(req, middlewareIndex)
}

// ReplacePathTokens invokes token replacement logic on the given url path
func ReplacePathTokens(path string, replacementPairs map[string]string) string {
	for key, value := range replacementPairs {
		path = strings.Replace(path, key, value, 1)
	}
	return path
}
