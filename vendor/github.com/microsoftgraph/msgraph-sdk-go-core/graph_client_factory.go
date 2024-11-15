package msgraphgocore

import (
	nethttp "net/http"

	khttp "github.com/microsoft/kiota-http-go"
)

var ReplacementPairs = map[string]string{"/users/me-token-to-replace": "/me"}

// GetDefaultMiddlewaresWithOptions creates a default slice of middleware for the Graph Client.
func GetDefaultMiddlewaresWithOptions(options *GraphClientOptions) []khttp.Middleware {
	kiotaMiddlewares := khttp.GetDefaultMiddlewares()
	graphMiddlewares := []khttp.Middleware{
		NewGraphTelemetryHandler(options),
		khttp.NewUrlReplaceHandler(true, ReplacementPairs),
	}
	graphMiddlewaresLen := len(graphMiddlewares)
	resultMiddlewares := make([]khttp.Middleware, len(kiotaMiddlewares)+graphMiddlewaresLen)
	copy(resultMiddlewares, graphMiddlewares)
	copy(resultMiddlewares[graphMiddlewaresLen:], kiotaMiddlewares)
	return resultMiddlewares
}

// GetDefaultClient creates a new http client with a preconfigured middleware pipeline
func GetDefaultClient(options *GraphClientOptions, middleware ...khttp.Middleware) *nethttp.Client {
	if len(middleware) == 0 {
		middleware = GetDefaultMiddlewaresWithOptions(options)
	}
	return khttp.GetDefaultClient(middleware...)
}
