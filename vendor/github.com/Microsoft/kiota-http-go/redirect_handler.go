package nethttplibrary

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"net/url"
	"strings"

	abs "github.com/microsoft/kiota-abstractions-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RedirectHandler handles redirect responses and follows them according to the options specified.
type RedirectHandler struct {
	// options to use when evaluating whether to redirect or not
	options RedirectHandlerOptions
}

// NewRedirectHandler creates a new redirect handler with the default options.
func NewRedirectHandler() *RedirectHandler {
	return NewRedirectHandlerWithOptions(RedirectHandlerOptions{
		MaxRedirects: defaultMaxRedirects,
		ShouldRedirect: func(req *nethttp.Request, res *nethttp.Response) bool {
			return true
		},
	})
}

// NewRedirectHandlerWithOptions creates a new redirect handler with the specified options.
func NewRedirectHandlerWithOptions(options RedirectHandlerOptions) *RedirectHandler {
	return &RedirectHandler{options: options}
}

// RedirectHandlerOptions to use when evaluating whether to redirect or not.
type RedirectHandlerOptions struct {
	// A callback that determines whether to redirect or not.
	ShouldRedirect func(req *nethttp.Request, res *nethttp.Response) bool
	// The maximum number of redirects to follow.
	MaxRedirects int
}

var redirectKeyValue = abs.RequestOptionKey{
	Key: "RedirectHandler",
}

type redirectHandlerOptionsInt interface {
	abs.RequestOption
	GetShouldRedirect() func(req *nethttp.Request, res *nethttp.Response) bool
	GetMaxRedirect() int
}

// GetKey returns the key value to be used when the option is added to the request context
func (options *RedirectHandlerOptions) GetKey() abs.RequestOptionKey {
	return redirectKeyValue
}

// GetShouldRedirect returns the redirection evaluation function.
func (options *RedirectHandlerOptions) GetShouldRedirect() func(req *nethttp.Request, res *nethttp.Response) bool {
	return options.ShouldRedirect
}

// GetMaxRedirect returns the maximum number of redirects to follow.
func (options *RedirectHandlerOptions) GetMaxRedirect() int {
	if options == nil || options.MaxRedirects < 1 {
		return defaultMaxRedirects
	} else if options.MaxRedirects > absoluteMaxRedirects {
		return absoluteMaxRedirects
	} else {
		return options.MaxRedirects
	}
}

const defaultMaxRedirects = 5
const absoluteMaxRedirects = 20
const movedPermanently = 301
const found = 302
const seeOther = 303
const temporaryRedirect = 307
const permanentRedirect = 308
const locationHeader = "Location"

// Intercept implements the interface and evaluates whether to follow a redirect response.
func (middleware RedirectHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	var observabilityName string
	if obsOptions != nil {
		observabilityName = obsOptions.GetTracerInstrumentationName()
		ctx, span = otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "RedirectHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.redirect.enable", true))
		defer span.End()
		req = req.WithContext(ctx)
	}
	response, err := pipeline.Next(req, middlewareIndex)
	if err != nil {
		return response, err
	}
	reqOption, ok := req.Context().Value(redirectKeyValue).(redirectHandlerOptionsInt)
	if !ok {
		reqOption = &middleware.options
	}
	return middleware.redirectRequest(ctx, pipeline, middlewareIndex, reqOption, req, response, 0, observabilityName)
}

func (middleware RedirectHandler) redirectRequest(ctx context.Context, pipeline Pipeline, middlewareIndex int, reqOption redirectHandlerOptionsInt, req *nethttp.Request, response *nethttp.Response, redirectCount int, observabilityName string) (*nethttp.Response, error) {
	shouldRedirect := reqOption.GetShouldRedirect() != nil && reqOption.GetShouldRedirect()(req, response) || reqOption.GetShouldRedirect() == nil
	if middleware.isRedirectResponse(response) &&
		redirectCount < reqOption.GetMaxRedirect() &&
		shouldRedirect {
		redirectCount++
		redirectRequest, err := middleware.getRedirectRequest(req, response)
		if err != nil {
			return response, err
		}
		if observabilityName != "" {
			ctx, span := otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "RedirectHandler_Intercept - redirect "+fmt.Sprint(redirectCount))
			span.SetAttributes(attribute.Int("com.microsoft.kiota.handler.redirect.count", redirectCount),
				attribute.Int("http.status_code", response.StatusCode),
			)
			defer span.End()
			redirectRequest = redirectRequest.WithContext(ctx)
		}

		result, err := pipeline.Next(redirectRequest, middlewareIndex)
		if err != nil {
			return result, err
		}
		return middleware.redirectRequest(ctx, pipeline, middlewareIndex, reqOption, redirectRequest, result, redirectCount, observabilityName)
	}
	return response, nil
}

func (middleware RedirectHandler) isRedirectResponse(response *nethttp.Response) bool {
	if response == nil {
		return false
	}
	locationHeader := response.Header.Get(locationHeader)
	if locationHeader == "" {
		return false
	}
	statusCode := response.StatusCode
	return statusCode == movedPermanently || statusCode == found || statusCode == seeOther || statusCode == temporaryRedirect || statusCode == permanentRedirect
}

func (middleware RedirectHandler) getRedirectRequest(request *nethttp.Request, response *nethttp.Response) (*nethttp.Request, error) {
	if request == nil || response == nil {
		return nil, errors.New("request or response is nil")
	}
	locationHeaderValue := response.Header.Get(locationHeader)
	if locationHeaderValue[0] == '/' {
		locationHeaderValue = request.URL.Scheme + "://" + request.URL.Host + locationHeaderValue
	}
	result := request.Clone(request.Context())
	targetUrl, err := url.Parse(locationHeaderValue)
	if err != nil {
		return nil, err
	}
	result.URL = targetUrl
	if result.Host != targetUrl.Host {
		result.Host = targetUrl.Host
	}
	sameHost := strings.EqualFold(targetUrl.Host, request.URL.Host)
	sameScheme := strings.EqualFold(targetUrl.Scheme, request.URL.Scheme)
	if !sameHost || !sameScheme {
		result.Header.Del("Authorization")
	}
	if response.StatusCode == seeOther {
		result.Method = nethttp.MethodGet
		result.Header.Del("Content-Type")
		result.Header.Del("Content-Length")
		result.Body = nil
	}
	return result, nil
}
