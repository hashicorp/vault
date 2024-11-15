package nethttplibrary

import nethttp "net/http"

// Middleware interface for cross cutting concerns with HTTP requests and responses.
type Middleware interface {
	// Intercept intercepts the request and returns the response. The implementer MUST call pipeline.Next()
	Intercept(Pipeline, int, *nethttp.Request) (*nethttp.Response, error)
}
