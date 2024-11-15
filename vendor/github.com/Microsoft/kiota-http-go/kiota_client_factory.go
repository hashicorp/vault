// Package nethttplibrary implements the Kiota abstractions with net/http to execute the requests.
// It also provides a middleware infrastructure with some default middleware handlers like the retry handler and the redirect handler.
package nethttplibrary

import (
	"errors"
	abs "github.com/microsoft/kiota-abstractions-go"
	nethttp "net/http"
	"net/url"
	"time"
)

// GetClientWithProxySettings creates a new default net/http client with a proxy url and default middleware
// Not providing any middleware would result in having default middleware provided
func GetClientWithProxySettings(proxyUrlStr string, middleware ...Middleware) (*nethttp.Client, error) {
	client := getDefaultClientWithoutMiddleware()

	transport, err := getTransportWithProxy(proxyUrlStr, nil, middleware...)
	if err != nil {
		return nil, err
	}
	client.Transport = transport
	return client, nil
}

// GetClientWithAuthenticatedProxySettings creates a new default net/http client with a proxy url and default middleware
// Not providing any middleware would result in having default middleware provided
func GetClientWithAuthenticatedProxySettings(proxyUrlStr string, username string, password string, middleware ...Middleware) (*nethttp.Client, error) {
	client := getDefaultClientWithoutMiddleware()

	user := url.UserPassword(username, password)
	transport, err := getTransportWithProxy(proxyUrlStr, user, middleware...)
	if err != nil {
		return nil, err
	}
	client.Transport = transport
	return client, nil
}

func getTransportWithProxy(proxyUrlStr string, user *url.Userinfo, middlewares ...Middleware) (nethttp.RoundTripper, error) {
	proxyURL, err := url.Parse(proxyUrlStr)
	if err != nil {
		return nil, err
	}

	if user != nil {
		proxyURL.User = user
	}

	transport := &nethttp.Transport{
		Proxy: nethttp.ProxyURL(proxyURL),
	}

	if len(middlewares) == 0 {
		middlewares = GetDefaultMiddlewares()
	}

	return NewCustomTransportWithParentTransport(transport, middlewares...), nil
}

// GetDefaultClient creates a new default net/http client with the options configured for the Kiota request adapter
func GetDefaultClient(middleware ...Middleware) *nethttp.Client {
	client := getDefaultClientWithoutMiddleware()
	client.Transport = NewCustomTransport(middleware...)
	return client
}

// used for internal unit testing
func getDefaultClientWithoutMiddleware() *nethttp.Client {
	// the default client doesn't come with any other settings than making a new one does, and using the default client impacts behavior for non-kiota requests
	return &nethttp.Client{
		CheckRedirect: func(req *nethttp.Request, via []*nethttp.Request) error {
			return nethttp.ErrUseLastResponse
		},
		Timeout: time.Second * 100,
	}
}

// GetDefaultMiddlewares creates a new default set of middlewares for the Kiota request adapter
func GetDefaultMiddlewares() []Middleware {
	return getDefaultMiddleWare(make(map[abs.RequestOptionKey]Middleware))
}

// GetDefaultMiddlewaresWithOptions creates a new default set of middlewares for the Kiota request adapter with options
func GetDefaultMiddlewaresWithOptions(requestOptions ...abs.RequestOption) ([]Middleware, error) {
	if len(requestOptions) == 0 {
		return GetDefaultMiddlewares(), nil
	}

	// map of middleware options
	middlewareMap := make(map[abs.RequestOptionKey]Middleware)

	for _, element := range requestOptions {
		switch v := element.(type) {
		case *RetryHandlerOptions:
			middlewareMap[retryKeyValue] = NewRetryHandlerWithOptions(*v)
		case *RedirectHandlerOptions:
			middlewareMap[redirectKeyValue] = NewRedirectHandlerWithOptions(*v)
		case *CompressionOptions:
			middlewareMap[compressKey] = NewCompressionHandlerWithOptions(*v)
		case *ParametersNameDecodingOptions:
			middlewareMap[parametersNameDecodingKeyValue] = NewParametersNameDecodingHandlerWithOptions(*v)
		case *UserAgentHandlerOptions:
			middlewareMap[userAgentKeyValue] = NewUserAgentHandlerWithOptions(v)
		case *HeadersInspectionOptions:
			middlewareMap[headersInspectionKeyValue] = NewHeadersInspectionHandlerWithOptions(*v)
		default:
			// none of the above types
			return nil, errors.New("unsupported option type")
		}
	}

	middleware := getDefaultMiddleWare(middlewareMap)
	return middleware, nil
}

// getDefaultMiddleWare creates a new default set of middlewares for the Kiota request adapter
func getDefaultMiddleWare(middlewareMap map[abs.RequestOptionKey]Middleware) []Middleware {
	middlewareSource := map[abs.RequestOptionKey]func() Middleware{
		retryKeyValue: func() Middleware {
			return NewRetryHandler()
		},
		redirectKeyValue: func() Middleware {
			return NewRedirectHandler()
		},
		compressKey: func() Middleware {
			return NewCompressionHandler()
		},
		parametersNameDecodingKeyValue: func() Middleware {
			return NewParametersNameDecodingHandler()
		},
		userAgentKeyValue: func() Middleware {
			return NewUserAgentHandler()
		},
		headersInspectionKeyValue: func() Middleware {
			return NewHeadersInspectionHandler()
		},
	}

	// loop over middlewareSource and add any middleware that wasn't provided in the requestOptions
	for key, value := range middlewareSource {
		if _, ok := middlewareMap[key]; !ok {
			middlewareMap[key] = value()
		}
	}

	var middleware []Middleware
	for _, value := range middlewareMap {
		middleware = append(middleware, value)
	}

	return middleware
}
