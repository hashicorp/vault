// Copyright (c) 2015-2024 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

var (
	hdrUserAgentKey       = http.CanonicalHeaderKey("User-Agent")
	hdrAcceptKey          = http.CanonicalHeaderKey("Accept")
	hdrContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	hdrContentLengthKey   = http.CanonicalHeaderKey("Content-Length")
	hdrContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
	hdrLocationKey        = http.CanonicalHeaderKey("Location")
	hdrAuthorizationKey   = http.CanonicalHeaderKey("Authorization")
	hdrWwwAuthenticateKey = http.CanonicalHeaderKey("WWW-Authenticate")

	plainTextType   = "text/plain; charset=utf-8"
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"

	jsonCheck = regexp.MustCompile(`(?i:(application|text)/(.*json.*)(;|$))`)
	xmlCheck  = regexp.MustCompile(`(?i:(application|text)/(.*xml.*)(;|$))`)

	hdrUserAgentValue = "go-resty/" + Version + " (https://github.com/go-resty/resty)"
	bufPool           = &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
)

type (
	// RequestMiddleware type is for request middleware, called before a request is sent
	RequestMiddleware func(*Client, *Request) error

	// ResponseMiddleware type is for response middleware, called after a response has been received
	ResponseMiddleware func(*Client, *Response) error

	// PreRequestHook type is for the request hook, called right before the request is sent
	PreRequestHook func(*Client, *http.Request) error

	// RequestLogCallback type is for request logs, called before the request is logged
	RequestLogCallback func(*RequestLog) error

	// ResponseLogCallback type is for response logs, called before the response is logged
	ResponseLogCallback func(*ResponseLog) error

	// ErrorHook type is for reacting to request errors, called after all retries were attempted
	ErrorHook func(*Request, error)

	// SuccessHook type is for reacting to request success
	SuccessHook func(*Client, *Response)
)

// Client struct is used to create a Resty client with client-level settings,
// these settings apply to all the requests raised from the client.
//
// Resty also provides an option to override most of the client settings
// at [Request] level.
type Client struct {
	BaseURL               string
	HostURL               string // Deprecated: use BaseURL instead. To be removed in v3.0.0 release.
	QueryParam            url.Values
	FormData              url.Values
	PathParams            map[string]string
	RawPathParams         map[string]string
	Header                http.Header
	UserInfo              *User
	Token                 string
	AuthScheme            string
	Cookies               []*http.Cookie
	Error                 reflect.Type
	Debug                 bool
	DisableWarn           bool
	AllowGetMethodPayload bool
	RetryCount            int
	RetryWaitTime         time.Duration
	RetryMaxWaitTime      time.Duration
	RetryConditions       []RetryConditionFunc
	RetryHooks            []OnRetryFunc
	RetryAfter            RetryAfterFunc
	RetryResetReaders     bool
	JSONMarshal           func(v interface{}) ([]byte, error)
	JSONUnmarshal         func(data []byte, v interface{}) error
	XMLMarshal            func(v interface{}) ([]byte, error)
	XMLUnmarshal          func(data []byte, v interface{}) error

	// HeaderAuthorizationKey is used to set/access Request Authorization header
	// value when `SetAuthToken` option is used.
	HeaderAuthorizationKey string
	ResponseBodyLimit      int

	jsonEscapeHTML      bool
	setContentLength    bool
	closeConnection     bool
	notParseResponse    bool
	trace               bool
	debugBodySizeLimit  int64
	outputDirectory     string
	scheme              string
	log                 Logger
	httpClient          *http.Client
	proxyURL            *url.URL
	beforeRequest       []RequestMiddleware
	udBeforeRequest     []RequestMiddleware
	udBeforeRequestLock *sync.RWMutex
	preReqHook          PreRequestHook
	successHooks        []SuccessHook
	afterResponse       []ResponseMiddleware
	afterResponseLock   *sync.RWMutex
	requestLog          RequestLogCallback
	responseLog         ResponseLogCallback
	errorHooks          []ErrorHook
	invalidHooks        []ErrorHook
	panicHooks          []ErrorHook
	rateLimiter         RateLimiter
	generateCurlOnDebug bool
	unescapeQueryParams bool
}

// User type is to hold an username and password information
type User struct {
	Username, Password string
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Client methods
//___________________________________

// SetHostURL method sets the Host URL in the client instance. It will be used with a request
// raised from this client with a relative URL
//
//	// Setting HTTP address
//	client.SetHostURL("http://myjeeva.com")
//
//	// Setting HTTPS address
//	client.SetHostURL("https://myjeeva.com")
//
// Deprecated: use [Client.SetBaseURL] instead. To be removed in the v3.0.0 release.
func (c *Client) SetHostURL(url string) *Client {
	c.SetBaseURL(url)
	return c
}

// SetBaseURL method sets the Base URL in the client instance. It will be used with a request
// raised from this client with a relative URL
//
//	// Setting HTTP address
//	client.SetBaseURL("http://myjeeva.com")
//
//	// Setting HTTPS address
//	client.SetBaseURL("https://myjeeva.com")
func (c *Client) SetBaseURL(url string) *Client {
	c.BaseURL = strings.TrimRight(url, "/")
	c.HostURL = c.BaseURL
	return c
}

// SetHeader method sets a single header field and its value in the client instance.
// These headers will be applied to all requests from this client instance.
// Also, it can be overridden by request-level header options.
//
// See [Request.SetHeader] or [Request.SetHeaders].
//
// For Example: To set `Content-Type` and `Accept` as `application/json`
//
//	client.
//		SetHeader("Content-Type", "application/json").
//		SetHeader("Accept", "application/json")
func (c *Client) SetHeader(header, value string) *Client {
	c.Header.Set(header, value)
	return c
}

// SetHeaders method sets multiple header fields and their values at one go in the client instance.
// These headers will be applied to all requests from this client instance. Also, it can be
// overridden at request level headers options.
//
// See [Request.SetHeaders] or [Request.SetHeader].
//
// For Example: To set `Content-Type` and `Accept` as `application/json`
//
//	client.SetHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept": "application/json",
//		})
func (c *Client) SetHeaders(headers map[string]string) *Client {
	for h, v := range headers {
		c.Header.Set(h, v)
	}
	return c
}

// SetHeaderVerbatim method sets a single header field and its value verbatim in the current request.
//
// For Example: To set `all_lowercase` and `UPPERCASE` as `available`.
//
//	client.
//		SetHeaderVerbatim("all_lowercase", "available").
//		SetHeaderVerbatim("UPPERCASE", "available")
func (c *Client) SetHeaderVerbatim(header, value string) *Client {
	c.Header[header] = []string{value}
	return c
}

// SetCookieJar method sets custom [http.CookieJar] in the resty client. It's a way to override the default.
//
// For Example, sometimes we don't want to save cookies in API mode so that we can remove the default
// CookieJar in resty client.
//
//	client.SetCookieJar(nil)
func (c *Client) SetCookieJar(jar http.CookieJar) *Client {
	c.httpClient.Jar = jar
	return c
}

// SetCookie method appends a single cookie to the client instance.
// These cookies will be added to all the requests from this client instance.
//
//	client.SetCookie(&http.Cookie{
//				Name:"go-resty",
//				Value:"This is cookie value",
//			})
func (c *Client) SetCookie(hc *http.Cookie) *Client {
	c.Cookies = append(c.Cookies, hc)
	return c
}

// SetCookies method sets an array of cookies in the client instance.
// These cookies will be added to all the requests from this client instance.
//
//	cookies := []*http.Cookie{
//		&http.Cookie{
//			Name:"go-resty-1",
//			Value:"This is cookie 1 value",
//		},
//		&http.Cookie{
//			Name:"go-resty-2",
//			Value:"This is cookie 2 value",
//		},
//	}
//
//	// Setting a cookies into resty
//	client.SetCookies(cookies)
func (c *Client) SetCookies(cs []*http.Cookie) *Client {
	c.Cookies = append(c.Cookies, cs...)
	return c
}

// SetQueryParam method sets a single parameter and its value in the client instance.
// It will be formed as a query string for the request.
//
//	For Example: `search=kitchen%20papers&size=large`
//
// In the URL after the `?` mark. These query params will be added to all the requests raised from
// this client instance. Also, it can be overridden at the request level.
//
// See [Request.SetQueryParam] or [Request.SetQueryParams].
//
//	client.
//		SetQueryParam("search", "kitchen papers").
//		SetQueryParam("size", "large")
func (c *Client) SetQueryParam(param, value string) *Client {
	c.QueryParam.Set(param, value)
	return c
}

// SetQueryParams method sets multiple parameters and their values at one go in the client instance.
// It will be formed as a query string for the request.
//
//	For Example: `search=kitchen%20papers&size=large`
//
// In the URL after the `?` mark. These query params will be added to all the requests raised from this
// client instance. Also, it can be overridden at the request level.
//
// See [Request.SetQueryParams] or [Request.SetQueryParam].
//
//	client.SetQueryParams(map[string]string{
//			"search": "kitchen papers",
//			"size": "large",
//		})
func (c *Client) SetQueryParams(params map[string]string) *Client {
	for p, v := range params {
		c.SetQueryParam(p, v)
	}
	return c
}

// SetUnescapeQueryParams method sets the unescape query parameters choice for request URL.
// To prevent broken URL, resty replaces space (" ") with "+" in the query parameters.
//
// See [Request.SetUnescapeQueryParams]
//
// NOTE: Request failure is possible due to non-standard usage of Unescaped Query Parameters.
func (c *Client) SetUnescapeQueryParams(unescape bool) *Client {
	c.unescapeQueryParams = unescape
	return c
}

// SetFormData method sets Form parameters and their values in the client instance.
// It applies only to HTTP methods `POST` and `PUT`, and the request content type would be set as
// `application/x-www-form-urlencoded`. These form data will be added to all the requests raised from
// this client instance. Also, it can be overridden at the request level.
//
// See [Request.SetFormData].
//
//	client.SetFormData(map[string]string{
//			"access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
//			"user_id": "3455454545",
//		})
func (c *Client) SetFormData(data map[string]string) *Client {
	for k, v := range data {
		c.FormData.Set(k, v)
	}
	return c
}

// SetBasicAuth method sets the basic authentication header in the HTTP request. For Example:
//
//	Authorization: Basic <base64-encoded-value>
//
// For Example: To set the header for username "go-resty" and password "welcome"
//
//	client.SetBasicAuth("go-resty", "welcome")
//
// This basic auth information is added to all requests from this client instance.
// It can also be overridden at the request level.
//
// See [Request.SetBasicAuth].
func (c *Client) SetBasicAuth(username, password string) *Client {
	c.UserInfo = &User{Username: username, Password: password}
	return c
}

// SetAuthToken method sets the auth token of the `Authorization` header for all HTTP requests.
// The default auth scheme is `Bearer`; it can be customized with the method [Client.SetAuthScheme]. For Example:
//
//	Authorization: <auth-scheme> <auth-token-value>
//
// For Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F
//
//	client.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
//
// This auth token gets added to all the requests raised from this client instance.
// Also, it can be overridden at the request level.
//
// See [Request.SetAuthToken].
func (c *Client) SetAuthToken(token string) *Client {
	c.Token = token
	return c
}

// SetAuthScheme method sets the auth scheme type in the HTTP request. For Example:
//
//	Authorization: <auth-scheme-value> <auth-token-value>
//
// For Example: To set the scheme to use OAuth
//
//	client.SetAuthScheme("OAuth")
//
// This auth scheme gets added to all the requests raised from this client instance.
// Also, it can be overridden at the request level.
//
// Information about auth schemes can be found in [RFC 7235], IANA [HTTP Auth schemes].
//
// See [Request.SetAuthToken].
//
// [RFC 7235]: https://tools.ietf.org/html/rfc7235
// [HTTP Auth schemes]: https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml#authschemes
func (c *Client) SetAuthScheme(scheme string) *Client {
	c.AuthScheme = scheme
	return c
}

// SetDigestAuth method sets the Digest Access auth scheme for the client. If a server responds with 401 and sends
// a Digest challenge in the WWW-Authenticate Header, requests will be resent with the appropriate Authorization Header.
//
// For Example: To set the Digest scheme with user "Mufasa" and password "Circle Of Life"
//
//	client.SetDigestAuth("Mufasa", "Circle Of Life")
//
// Information about Digest Access Authentication can be found in [RFC 7616].
//
// See [Request.SetDigestAuth].
//
// [RFC 7616]: https://datatracker.ietf.org/doc/html/rfc7616
func (c *Client) SetDigestAuth(username, password string) *Client {
	oldTransport := c.httpClient.Transport
	c.OnBeforeRequest(func(c *Client, _ *Request) error {
		c.httpClient.Transport = &digestTransport{
			digestCredentials: digestCredentials{username, password},
			transport:         oldTransport,
		}
		return nil
	})
	c.OnAfterResponse(func(c *Client, _ *Response) error {
		c.httpClient.Transport = oldTransport
		return nil
	})
	return c
}

// R method creates a new request instance; it's used for Get, Post, Put, Delete, Patch, Head, Options, etc.
func (c *Client) R() *Request {
	r := &Request{
		QueryParam:    url.Values{},
		FormData:      url.Values{},
		Header:        http.Header{},
		Cookies:       make([]*http.Cookie, 0),
		PathParams:    map[string]string{},
		RawPathParams: map[string]string{},
		Debug:         c.Debug,

		client:              c,
		multipartFiles:      []*File{},
		multipartFields:     []*MultipartField{},
		jsonEscapeHTML:      c.jsonEscapeHTML,
		log:                 c.log,
		responseBodyLimit:   c.ResponseBodyLimit,
		generateCurlOnDebug: c.generateCurlOnDebug,
		unescapeQueryParams: c.unescapeQueryParams,
	}
	return r
}

// NewRequest method is an alias for method `R()`.
func (c *Client) NewRequest() *Request {
	return c.R()
}

// OnBeforeRequest method appends a request middleware to the before request chain.
// The user-defined middlewares are applied before the default Resty request middlewares.
// After all middlewares have been applied, the request is sent from Resty to the host server.
//
//	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
//			// Now you have access to the Client and Request instance
//			// manipulate it as per your need
//
//			return nil 	// if its successful otherwise return error
//		})
func (c *Client) OnBeforeRequest(m RequestMiddleware) *Client {
	c.udBeforeRequestLock.Lock()
	defer c.udBeforeRequestLock.Unlock()

	c.udBeforeRequest = append(c.udBeforeRequest, m)

	return c
}

// OnAfterResponse method appends response middleware to the after-response chain.
// Once we receive a response from the host server, the default Resty response middleware
// gets applied, and then the user-assigned response middleware is applied.
//
//	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
//			// Now you have access to the Client and Response instance
//			// manipulate it as per your need
//
//			return nil 	// if its successful otherwise return error
//		})
func (c *Client) OnAfterResponse(m ResponseMiddleware) *Client {
	c.afterResponseLock.Lock()
	defer c.afterResponseLock.Unlock()

	c.afterResponse = append(c.afterResponse, m)

	return c
}

// OnError method adds a callback that will be run whenever a request execution fails.
// This is called after all retries have been attempted (if any).
// If there was a response from the server, the error will be wrapped in [ResponseError]
// which has the last response received from the server.
//
//	client.OnError(func(req *resty.Request, err error) {
//		if v, ok := err.(*resty.ResponseError); ok {
//			// Do something with v.Response
//		}
//		// Log the error, increment a metric, etc...
//	})
//
// Out of the [Client.OnSuccess], [Client.OnError], [Client.OnInvalid], [Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [Request.Execute] that completes.
func (c *Client) OnError(h ErrorHook) *Client {
	c.errorHooks = append(c.errorHooks, h)
	return c
}

// OnSuccess method adds a callback that will be run whenever a request execution
// succeeds.  This is called after all retries have been attempted (if any).
//
// Out of the [Client.OnSuccess], [Client.OnError], [Client.OnInvalid], [Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [Request.Execute] that completes.
func (c *Client) OnSuccess(h SuccessHook) *Client {
	c.successHooks = append(c.successHooks, h)
	return c
}

// OnInvalid method adds a callback that will be run whenever a request execution
// fails before it starts because the request is invalid.
//
// Out of the [Client.OnSuccess], [Client.OnError], [Client.OnInvalid], [Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [Request.Execute] that completes.
func (c *Client) OnInvalid(h ErrorHook) *Client {
	c.invalidHooks = append(c.invalidHooks, h)
	return c
}

// OnPanic method adds a callback that will be run whenever a request execution
// panics.
//
// Out of the [Client.OnSuccess], [Client.OnError], [Client.OnInvalid], [Client.OnPanic]
// callbacks, exactly one set will be invoked for each call to [Request.Execute] that completes.
//
// If an [Client.OnSuccess], [Client.OnError], or [Client.OnInvalid] callback panics,
// then exactly one rule can be violated.
func (c *Client) OnPanic(h ErrorHook) *Client {
	c.panicHooks = append(c.panicHooks, h)
	return c
}

// SetPreRequestHook method sets the given pre-request function into a resty client.
// It is called right before the request is fired.
//
// NOTE: Only one pre-request hook can be registered. Use [Client.OnBeforeRequest] for multiple.
func (c *Client) SetPreRequestHook(h PreRequestHook) *Client {
	if c.preReqHook != nil {
		c.log.Warnf("Overwriting an existing pre-request hook: %s", functionName(h))
	}
	c.preReqHook = h
	return c
}

// SetDebug method enables the debug mode on the Resty client. The client logs details
// of every request and response.
//
//	client.SetDebug(true)
//
// Also, it can be enabled at the request level for a particular request; see [Request.SetDebug].
//   - For [Request], it logs information such as HTTP verb, Relative URL path,
//     Host, Headers, and Body if it has one.
//   - For [Response], it logs information such as Status, Response Time, Headers,
//     and Body if it has one.
func (c *Client) SetDebug(d bool) *Client {
	c.Debug = d
	return c
}

// SetDebugBodyLimit sets the maximum size in bytes for which the response and
// request body will be logged in debug mode.
//
//	client.SetDebugBodyLimit(1000000)
func (c *Client) SetDebugBodyLimit(sl int64) *Client {
	c.debugBodySizeLimit = sl
	return c
}

// OnRequestLog method sets the request log callback to Resty. Registered callback gets
// called before the resty logs the information.
func (c *Client) OnRequestLog(rl RequestLogCallback) *Client {
	if c.requestLog != nil {
		c.log.Warnf("Overwriting an existing on-request-log callback from=%s to=%s",
			functionName(c.requestLog), functionName(rl))
	}
	c.requestLog = rl
	return c
}

// OnResponseLog method sets the response log callback to Resty. Registered callback gets
// called before the resty logs the information.
func (c *Client) OnResponseLog(rl ResponseLogCallback) *Client {
	if c.responseLog != nil {
		c.log.Warnf("Overwriting an existing on-response-log callback from=%s to=%s",
			functionName(c.responseLog), functionName(rl))
	}
	c.responseLog = rl
	return c
}

// SetDisableWarn method disables the warning log message on the Resty client.
//
// For example, Resty warns users when BasicAuth is used in non-TLS mode.
//
//	client.SetDisableWarn(true)
func (c *Client) SetDisableWarn(d bool) *Client {
	c.DisableWarn = d
	return c
}

// SetAllowGetMethodPayload method allows the GET method with payload on the Resty client.
//
// For example, Resty allows the user to send a request with a payload using the HTTP GET method.
//
//	client.SetAllowGetMethodPayload(true)
func (c *Client) SetAllowGetMethodPayload(a bool) *Client {
	c.AllowGetMethodPayload = a
	return c
}

// SetLogger method sets given writer for logging Resty request and response details.
//
// Compliant to interface [resty.Logger]
func (c *Client) SetLogger(l Logger) *Client {
	c.log = l
	return c
}

// SetContentLength method enables the HTTP header `Content-Length` value for every request.
// By default, Resty won't set `Content-Length`.
//
//	client.SetContentLength(true)
//
// Also, you have the option to enable a particular request. See [Request.SetContentLength]
func (c *Client) SetContentLength(l bool) *Client {
	c.setContentLength = l
	return c
}

// SetTimeout method sets the timeout for a request raised by the client.
//
//	client.SetTimeout(time.Duration(1 * time.Minute))
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// SetError method registers the global or client common `Error` object into Resty.
// It is used for automatic unmarshalling if the response status code is greater than 399 and
// content type is JSON or XML. It can be a pointer or a non-pointer.
//
//	client.SetError(&Error{})
//	// OR
//	client.SetError(Error{})
func (c *Client) SetError(err interface{}) *Client {
	c.Error = typeOf(err)
	return c
}

// SetRedirectPolicy method sets the redirect policy for the client. Resty provides ready-to-use
// redirect policies. Wanna create one for yourself, refer to `redirect.go`.
//
//	client.SetRedirectPolicy(FlexibleRedirectPolicy(20))
//
//	// Need multiple redirect policies together
//	client.SetRedirectPolicy(FlexibleRedirectPolicy(20), DomainCheckRedirectPolicy("host1.com", "host2.net"))
func (c *Client) SetRedirectPolicy(policies ...interface{}) *Client {
	for _, p := range policies {
		if _, ok := p.(RedirectPolicy); !ok {
			c.log.Errorf("%v does not implement resty.RedirectPolicy (missing Apply method)",
				functionName(p))
		}
	}

	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for _, p := range policies {
			if err := p.(RedirectPolicy).Apply(req, via); err != nil {
				return err
			}
		}
		return nil // looks good, go ahead
	}

	return c
}

// SetRetryCount method enables retry on Resty client and allows you
// to set no. of retry count. Resty uses a Backoff mechanism.
func (c *Client) SetRetryCount(count int) *Client {
	c.RetryCount = count
	return c
}

// SetRetryWaitTime method sets the default wait time for sleep before retrying
// request.
//
// Default is 100 milliseconds.
func (c *Client) SetRetryWaitTime(waitTime time.Duration) *Client {
	c.RetryWaitTime = waitTime
	return c
}

// SetRetryMaxWaitTime method sets the max wait time for sleep before retrying
// request.
//
// Default is 2 seconds.
func (c *Client) SetRetryMaxWaitTime(maxWaitTime time.Duration) *Client {
	c.RetryMaxWaitTime = maxWaitTime
	return c
}

// SetRetryAfter sets a callback to calculate the wait time between retries.
// Default (nil) implies exponential backoff with jitter
func (c *Client) SetRetryAfter(callback RetryAfterFunc) *Client {
	c.RetryAfter = callback
	return c
}

// SetJSONMarshaler method sets the JSON marshaler function to marshal the request body.
// By default, Resty uses [encoding/json] package to marshal the request body.
func (c *Client) SetJSONMarshaler(marshaler func(v interface{}) ([]byte, error)) *Client {
	c.JSONMarshal = marshaler
	return c
}

// SetJSONUnmarshaler method sets the JSON unmarshaler function to unmarshal the response body.
// By default, Resty uses [encoding/json] package to unmarshal the response body.
func (c *Client) SetJSONUnmarshaler(unmarshaler func(data []byte, v interface{}) error) *Client {
	c.JSONUnmarshal = unmarshaler
	return c
}

// SetXMLMarshaler method sets the XML marshaler function to marshal the request body.
// By default, Resty uses [encoding/xml] package to marshal the request body.
func (c *Client) SetXMLMarshaler(marshaler func(v interface{}) ([]byte, error)) *Client {
	c.XMLMarshal = marshaler
	return c
}

// SetXMLUnmarshaler method sets the XML unmarshaler function to unmarshal the response body.
// By default, Resty uses [encoding/xml] package to unmarshal the response body.
func (c *Client) SetXMLUnmarshaler(unmarshaler func(data []byte, v interface{}) error) *Client {
	c.XMLUnmarshal = unmarshaler
	return c
}

// AddRetryCondition method adds a retry condition function to an array of functions
// that are checked to determine if the request is retried. The request will
// retry if any functions return true and the error is nil.
//
// NOTE: These retry conditions are applied on all requests made using this Client.
// For [Request] specific retry conditions, check [Request.AddRetryCondition]
func (c *Client) AddRetryCondition(condition RetryConditionFunc) *Client {
	c.RetryConditions = append(c.RetryConditions, condition)
	return c
}

// AddRetryAfterErrorCondition adds the basic condition of retrying after encountering
// an error from the HTTP response
func (c *Client) AddRetryAfterErrorCondition() *Client {
	c.AddRetryCondition(func(response *Response, err error) bool {
		return response.IsError()
	})
	return c
}

// AddRetryHook adds a side-effecting retry hook to an array of hooks
// that will be executed on each retry.
func (c *Client) AddRetryHook(hook OnRetryFunc) *Client {
	c.RetryHooks = append(c.RetryHooks, hook)
	return c
}

// SetRetryResetReaders method enables the Resty client to seek the start of all
// file readers are given as multipart files if the object implements [io.ReadSeeker].
func (c *Client) SetRetryResetReaders(b bool) *Client {
	c.RetryResetReaders = b
	return c
}

// SetTLSClientConfig method sets TLSClientConfig for underlying client Transport.
//
// For Example:
//
//	// One can set a custom root certificate. Refer: http://golang.org/pkg/crypto/tls/#example_Dial
//	client.SetTLSClientConfig(&tls.Config{ RootCAs: roots })
//
//	// or One can disable security check (https)
//	client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: true })
//
// NOTE: This method overwrites existing [http.Transport.TLSClientConfig]
func (c *Client) SetTLSClientConfig(config *tls.Config) *Client {
	transport, err := c.Transport()
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}
	transport.TLSClientConfig = config
	return c
}

// SetProxy method sets the Proxy URL and Port for the Resty client.
//
//	client.SetProxy("http://proxyserver:8888")
//
// OR you could also set Proxy via environment variable, refer to [http.ProxyFromEnvironment]
func (c *Client) SetProxy(proxyURL string) *Client {
	transport, err := c.Transport()
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}

	pURL, err := url.Parse(proxyURL)
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}

	c.proxyURL = pURL
	transport.Proxy = http.ProxyURL(c.proxyURL)
	return c
}

// RemoveProxy method removes the proxy configuration from the Resty client
//
//	client.RemoveProxy()
func (c *Client) RemoveProxy() *Client {
	transport, err := c.Transport()
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}
	c.proxyURL = nil
	transport.Proxy = nil
	return c
}

// SetCertificates method helps to conveniently set client certificates into Resty.
func (c *Client) SetCertificates(certs ...tls.Certificate) *Client {
	config, err := c.tlsConfig()
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}
	config.Certificates = append(config.Certificates, certs...)
	return c
}

// SetRootCertificate method helps to add one or more root certificates into the Resty client
//
//	client.SetRootCertificate("/path/to/root/pemFile.pem")
func (c *Client) SetRootCertificate(pemFilePath string) *Client {
	rootPemData, err := os.ReadFile(pemFilePath)
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}
	c.handleCAs("root", rootPemData)
	return c
}

// SetRootCertificateFromString method helps to add one or more root certificates
// into the Resty client
//
//	client.SetRootCertificateFromString("pem certs content")
func (c *Client) SetRootCertificateFromString(pemCerts string) *Client {
	c.handleCAs("root", []byte(pemCerts))
	return c
}

// SetClientRootCertificate method helps to add one or more client's root
// certificates into the Resty client
//
//	client.SetClientRootCertificate("/path/to/root/pemFile.pem")
func (c *Client) SetClientRootCertificate(pemFilePath string) *Client {
	rootPemData, err := os.ReadFile(pemFilePath)
	if err != nil {
		c.log.Errorf("%v", err)
		return c
	}
	c.handleCAs("client", rootPemData)
	return c
}

// SetClientRootCertificateFromString method helps to add one or more clients
// root certificates into the Resty client
//
//	client.SetClientRootCertificateFromString("pem certs content")
func (c *Client) SetClientRootCertificateFromString(pemCerts string) *Client {
	c.handleCAs("client", []byte(pemCerts))
	return c
}

func (c *Client) handleCAs(scope string, permCerts []byte) {
	config, err := c.tlsConfig()
	if err != nil {
		c.log.Errorf("%v", err)
		return
	}

	switch scope {
	case "root":
		if config.RootCAs == nil {
			config.RootCAs = x509.NewCertPool()
		}
		config.RootCAs.AppendCertsFromPEM(permCerts)
	case "client":
		if config.ClientCAs == nil {
			config.ClientCAs = x509.NewCertPool()
		}
		config.ClientCAs.AppendCertsFromPEM(permCerts)
	}
}

// SetOutputDirectory method sets the output directory for saving HTTP responses in a file.
// Resty creates one if the output directory does not exist. This setting is optional,
// if you plan to use the absolute path in [Request.SetOutput] and can used together.
//
//	client.SetOutputDirectory("/save/http/response/here")
func (c *Client) SetOutputDirectory(dirPath string) *Client {
	c.outputDirectory = dirPath
	return c
}

// SetRateLimiter sets an optional [RateLimiter]. If set, the rate limiter will control
// all requests were made by this client.
func (c *Client) SetRateLimiter(rl RateLimiter) *Client {
	c.rateLimiter = rl
	return c
}

// SetTransport method sets custom [http.Transport] or any [http.RoundTripper]
// compatible interface implementation in the Resty client.
//
//	transport := &http.Transport{
//		// something like Proxying to httptest.Server, etc...
//		Proxy: func(req *http.Request) (*url.URL, error) {
//			return url.Parse(server.URL)
//		},
//	}
//	client.SetTransport(transport)
//
// NOTE:
//   - If transport is not the type of `*http.Transport`, then you may not be able to
//     take advantage of some of the Resty client settings.
//   - It overwrites the Resty client transport instance and its configurations.
func (c *Client) SetTransport(transport http.RoundTripper) *Client {
	if transport != nil {
		c.httpClient.Transport = transport
	}
	return c
}

// SetScheme method sets a custom scheme for the Resty client. It's a way to override the default.
//
//	client.SetScheme("http")
func (c *Client) SetScheme(scheme string) *Client {
	if !IsStringEmpty(scheme) {
		c.scheme = strings.TrimSpace(scheme)
	}
	return c
}

// SetCloseConnection method sets variable `Close` in HTTP request struct with the given
// value. More info: https://golang.org/src/net/http/request.go
func (c *Client) SetCloseConnection(close bool) *Client {
	c.closeConnection = close
	return c
}

// SetDoNotParseResponse method instructs Resty not to parse the response body automatically.
// Resty exposes the raw response body as [io.ReadCloser]. If you use it, do not
// forget to close the body, otherwise, you might get into connection leaks, and connection
// reuse may not happen.
//
// NOTE: [Response] middlewares are not executed using this option. You have
// taken over the control of response parsing from Resty.
func (c *Client) SetDoNotParseResponse(notParse bool) *Client {
	c.notParseResponse = notParse
	return c
}

// SetPathParam method sets a single URL path key-value pair in the
// Resty client instance.
//
//	client.SetPathParam("userId", "sample@sample.com")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/sample@sample.com/details
//
// It replaces the value of the key while composing the request URL.
// The value will be escaped using [url.PathEscape] function.
//
// It can be overridden at the request level,
// see [Request.SetPathParam] or [Request.SetPathParams]
func (c *Client) SetPathParam(param, value string) *Client {
	c.PathParams[param] = value
	return c
}

// SetPathParams method sets multiple URL path key-value pairs at one go in the
// Resty client instance.
//
//	client.SetPathParams(map[string]string{
//		"userId":       "sample@sample.com",
//		"subAccountId": "100002",
//		"path":         "groups/developers",
//	})
//
//	Result:
//	   URL - /v1/users/{userId}/{subAccountId}/{path}/details
//	   Composed URL - /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
//
// It replaces the value of the key while composing the request URL.
// The values will be escaped using [url.PathEscape] function.
//
// It can be overridden at the request level,
// see [Request.SetPathParam] or [Request.SetPathParams]
func (c *Client) SetPathParams(params map[string]string) *Client {
	for p, v := range params {
		c.SetPathParam(p, v)
	}
	return c
}

// SetRawPathParam method sets a single URL path key-value pair in the
// Resty client instance.
//
//	client.SetPathParam("userId", "sample@sample.com")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/sample@sample.com/details
//
//	client.SetPathParam("path", "groups/developers")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/groups%2Fdevelopers/details
//
// It replaces the value of the key while composing the request URL.
// The value will be used as it is and will not be escaped.
//
// It can be overridden at the request level,
// see [Request.SetRawPathParam] or [Request.SetRawPathParams]
func (c *Client) SetRawPathParam(param, value string) *Client {
	c.RawPathParams[param] = value
	return c
}

// SetRawPathParams method sets multiple URL path key-value pairs at one go in the
// Resty client instance.
//
//	client.SetPathParams(map[string]string{
//		"userId":       "sample@sample.com",
//		"subAccountId": "100002",
//		"path":         "groups/developers",
//	})
//
//	Result:
//	   URL - /v1/users/{userId}/{subAccountId}/{path}/details
//	   Composed URL - /v1/users/sample@sample.com/100002/groups/developers/details
//
// It replaces the value of the key while composing the request URL.
// The values will be used as they are and will not be escaped.
//
// It can be overridden at the request level,
// see [Request.SetRawPathParam] or [Request.SetRawPathParams]
func (c *Client) SetRawPathParams(params map[string]string) *Client {
	for p, v := range params {
		c.SetRawPathParam(p, v)
	}
	return c
}

// SetJSONEscapeHTML method enables or disables the HTML escape on JSON marshal.
// By default, escape HTML is false.
//
// NOTE: This option only applies to the standard JSON Marshaller used by Resty.
//
// It can be overridden at the request level, see [Client.SetJSONEscapeHTML]
func (c *Client) SetJSONEscapeHTML(b bool) *Client {
	c.jsonEscapeHTML = b
	return c
}

// SetResponseBodyLimit method sets a maximum body size limit in bytes on response,
// avoid reading too much data to memory.
//
// Client will return [resty.ErrResponseBodyTooLarge] if the body size of the body
// in the uncompressed response is larger than the limit.
// Body size limit will not be enforced in the following cases:
//   - ResponseBodyLimit <= 0, which is the default behavior.
//   - [Request.SetOutput] is called to save response data to the file.
//   - "DoNotParseResponse" is set for client or request.
//
// It can be overridden at the request level; see [Request.SetResponseBodyLimit]
func (c *Client) SetResponseBodyLimit(v int) *Client {
	c.ResponseBodyLimit = v
	return c
}

// EnableTrace method enables the Resty client trace for the requests fired from
// the client using [httptrace.ClientTrace] and provides insights.
//
//	client := resty.New().EnableTrace()
//
//	resp, err := client.R().Get("https://httpbin.org/get")
//	fmt.Println("Error:", err)
//	fmt.Println("Trace Info:", resp.Request.TraceInfo())
//
// The method [Request.EnableTrace] is also available to get trace info for a single request.
func (c *Client) EnableTrace() *Client {
	c.trace = true
	return c
}

// DisableTrace method disables the Resty client trace. Refer to [Client.EnableTrace].
func (c *Client) DisableTrace() *Client {
	c.trace = false
	return c
}

// EnableGenerateCurlOnDebug method enables the generation of CURL commands in the debug log.
// It works in conjunction with debug mode.
//
// NOTE: Use with care.
//   - Potential to leak sensitive data from [Request] and [Response] in the debug log.
//   - Beware of memory usage since the request body is reread.
func (c *Client) EnableGenerateCurlOnDebug() *Client {
	c.generateCurlOnDebug = true
	return c
}

// DisableGenerateCurlOnDebug method disables the option set by [Client.EnableGenerateCurlOnDebug].
func (c *Client) DisableGenerateCurlOnDebug() *Client {
	c.generateCurlOnDebug = false
	return c
}

// IsProxySet method returns the true is proxy is set from the Resty client; otherwise
// false. By default, the proxy is set from the environment variable; refer to [http.ProxyFromEnvironment].
func (c *Client) IsProxySet() bool {
	return c.proxyURL != nil
}

// GetClient method returns the underlying [http.Client] used by the Resty.
func (c *Client) GetClient() *http.Client {
	return c.httpClient
}

// Clone returns a clone of the original client.
//
// NOTE: Use with care:
//   - Interface values are not deeply cloned. Thus, both the original and the
//     clone will use the same value.
//   - This function is not safe for concurrent use. You should only use this method
//     when you are sure that any other goroutine is not using the client.
func (c *Client) Clone() *Client {
	// dereference the pointer and copy the value
	cc := *c

	// lock values should not be copied - thus new values are used.
	cc.afterResponseLock = &sync.RWMutex{}
	cc.udBeforeRequestLock = &sync.RWMutex{}
	return &cc
}

func (c *Client) executeBefore(req *Request) error {
	// Lock the user-defined pre-request hooks.
	c.udBeforeRequestLock.RLock()
	defer c.udBeforeRequestLock.RUnlock()

	// Lock the post-request hooks.
	c.afterResponseLock.RLock()
	defer c.afterResponseLock.RUnlock()

	// Apply Request middleware
	var err error

	// user defined on before request methods
	// to modify the *resty.Request object
	for _, f := range c.udBeforeRequest {
		if err = f(c, req); err != nil {
			return wrapNoRetryErr(err)
		}
	}

	// If there is a rate limiter set for this client, the Execute call
	// will return an error if the rate limit is exceeded.
	if req.client.rateLimiter != nil {
		if !req.client.rateLimiter.Allow() {
			return wrapNoRetryErr(ErrRateLimitExceeded)
		}
	}

	// resty middlewares
	for _, f := range c.beforeRequest {
		if err = f(c, req); err != nil {
			return wrapNoRetryErr(err)
		}
	}

	if hostHeader := req.Header.Get("Host"); hostHeader != "" {
		req.RawRequest.Host = hostHeader
	}

	// call pre-request if defined
	if c.preReqHook != nil {
		if err = c.preReqHook(c, req.RawRequest); err != nil {
			return wrapNoRetryErr(err)
		}
	}

	if err = requestLogger(c, req); err != nil {
		return wrapNoRetryErr(err)
	}

	req.RawRequest.Body = newRequestBodyReleaser(req.RawRequest.Body, req.bodyBuf)
	return nil
}

// Executes method executes the given `Request` object and returns
// response or error.
func (c *Client) execute(req *Request) (*Response, error) {
	if err := c.executeBefore(req); err != nil {
		return nil, err
	}

	req.Time = time.Now()
	resp, err := c.httpClient.Do(req.RawRequest)

	response := &Response{
		Request:     req,
		RawResponse: resp,
	}

	if err != nil || req.notParseResponse || c.notParseResponse {
		response.setReceivedAt()
		if logErr := responseLogger(c, response); logErr != nil {
			return response, wrapErrors(logErr, err)
		}
		if err != nil {
			return response, err
		}
		return response, nil
	}

	if !req.isSaveResponse {
		defer closeq(resp.Body)
		body := resp.Body

		// GitHub #142 & #187
		if strings.EqualFold(resp.Header.Get(hdrContentEncodingKey), "gzip") && resp.ContentLength != 0 {
			if _, ok := body.(*gzip.Reader); !ok {
				body, err = gzip.NewReader(body)
				if err != nil {
					err = wrapErrors(responseLogger(c, response), err)
					response.setReceivedAt()
					return response, err
				}
				defer closeq(body)
			}
		}

		if response.body, err = readAllWithLimit(body, req.responseBodyLimit); err != nil {
			err = wrapErrors(responseLogger(c, response), err)
			response.setReceivedAt()
			return response, err
		}

		response.size = int64(len(response.body))
	}

	response.setReceivedAt() // after we read the body

	// Apply Response middleware
	err = responseLogger(c, response)
	if err != nil {
		return response, wrapNoRetryErr(err)
	}

	for _, f := range c.afterResponse {
		if err = f(c, response); err != nil {
			break
		}
	}

	return response, wrapNoRetryErr(err)
}

var ErrResponseBodyTooLarge = errors.New("resty: response body too large")

// https://github.com/golang/go/issues/51115
// [io.LimitedReader] can only return [io.EOF]
func readAllWithLimit(r io.Reader, maxSize int) ([]byte, error) {
	if maxSize <= 0 {
		return io.ReadAll(r)
	}

	var buf [512]byte // make buf stack allocated
	result := make([]byte, 0, 512)
	total := 0
	for {
		n, err := r.Read(buf[:])
		total += n
		if total > maxSize {
			return nil, ErrResponseBodyTooLarge
		}

		if err != nil {
			if err == io.EOF {
				result = append(result, buf[:n]...)
				break
			}
			return nil, err
		}

		result = append(result, buf[:n]...)
	}

	return result, nil
}

// getting TLS client config if not exists then create one
func (c *Client) tlsConfig() (*tls.Config, error) {
	transport, err := c.Transport()
	if err != nil {
		return nil, err
	}
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	return transport.TLSClientConfig, nil
}

// Transport method returns [http.Transport] currently in use or error
// in case the currently used `transport` is not a [http.Transport].
//
// Since v2.8.0 has become exported method.
func (c *Client) Transport() (*http.Transport, error) {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		return transport, nil
	}
	return nil, errors.New("current transport is not an *http.Transport instance")
}

// just an internal helper method
func (c *Client) outputLogTo(w io.Writer) *Client {
	c.log.(*logger).l.SetOutput(w)
	return c
}

// ResponseError is a wrapper that includes the server response with an error.
// Neither the err nor the response should be nil.
type ResponseError struct {
	Response *Response
	Err      error
}

func (e *ResponseError) Error() string {
	return e.Err.Error()
}

func (e *ResponseError) Unwrap() error {
	return e.Err
}

// Helper to run errorHooks hooks.
// It wraps the error in a [ResponseError] if the resp is not nil
// so hooks can access it.
func (c *Client) onErrorHooks(req *Request, resp *Response, err error) {
	if err != nil {
		if resp != nil { // wrap with ResponseError
			err = &ResponseError{Response: resp, Err: err}
		}
		for _, h := range c.errorHooks {
			h(req, err)
		}
	} else {
		for _, h := range c.successHooks {
			h(c, resp)
		}
	}
}

// Helper to run panicHooks hooks.
func (c *Client) onPanicHooks(req *Request, err error) {
	for _, h := range c.panicHooks {
		h(req, err)
	}
}

// Helper to run invalidHooks hooks.
func (c *Client) onInvalidHooks(req *Request, err error) {
	for _, h := range c.invalidHooks {
		h(req, err)
	}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// File struct and its methods
//_______________________________________________________________________

// File struct represents file information for multipart request
type File struct {
	Name      string
	ParamName string
	io.Reader
}

// String method returns the string value of current file details
func (f *File) String() string {
	return fmt.Sprintf("ParamName: %v; FileName: %v", f.ParamName, f.Name)
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// MultipartField struct
//_______________________________________________________________________

// MultipartField struct represents the custom data part for a multipart request
type MultipartField struct {
	Param       string
	FileName    string
	ContentType string
	io.Reader
}

func createClient(hc *http.Client) *Client {
	if hc.Transport == nil {
		hc.Transport = createTransport(nil)
	}

	c := &Client{ // not setting lang default values
		QueryParam:             url.Values{},
		FormData:               url.Values{},
		Header:                 http.Header{},
		Cookies:                make([]*http.Cookie, 0),
		RetryWaitTime:          defaultWaitTime,
		RetryMaxWaitTime:       defaultMaxWaitTime,
		PathParams:             make(map[string]string),
		RawPathParams:          make(map[string]string),
		JSONMarshal:            json.Marshal,
		JSONUnmarshal:          json.Unmarshal,
		XMLMarshal:             xml.Marshal,
		XMLUnmarshal:           xml.Unmarshal,
		HeaderAuthorizationKey: http.CanonicalHeaderKey("Authorization"),

		jsonEscapeHTML:      true,
		httpClient:          hc,
		debugBodySizeLimit:  math.MaxInt32,
		udBeforeRequestLock: &sync.RWMutex{},
		afterResponseLock:   &sync.RWMutex{},
	}

	// Logger
	c.SetLogger(createLogger())

	// default before request middlewares
	c.beforeRequest = []RequestMiddleware{
		parseRequestURL,
		parseRequestHeader,
		parseRequestBody,
		createHTTPRequest,
		addCredentials,
		createCurlCmd,
	}

	// user defined request middlewares
	c.udBeforeRequest = []RequestMiddleware{}

	// default after response middlewares
	c.afterResponse = []ResponseMiddleware{
		parseResponseBody,
		saveResponseIntoFile,
	}

	return c
}
