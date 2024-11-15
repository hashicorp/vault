package godo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const (
	libraryVersion = "1.130.0"
	defaultBaseURL = "https://api.digitalocean.com/"
	userAgent      = "godo/" + libraryVersion
	mediaType      = "application/json"

	headerRateLimit             = "RateLimit-Limit"
	headerRateRemaining         = "RateLimit-Remaining"
	headerRateReset             = "RateLimit-Reset"
	headerRequestID             = "x-request-id"
	internalHeaderRetryAttempts = "X-Godo-Retry-Attempts"

	defaultRetryMax     = 4
	defaultRetryWaitMax = 30
	defaultRetryWaitMin = 1
)

// Client manages communication with DigitalOcean V2 API.
type Client struct {
	// HTTP client used to communicate with the DO API.
	HTTPClient *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Rate contains the current rate limit for the client as determined by the most recent
	// API call. It is not thread-safe. Please consider using GetRate() instead.
	Rate    Rate
	ratemtx sync.Mutex

	// Services used for communicating with the API
	Account           AccountService
	Actions           ActionsService
	Apps              AppsService
	Balance           BalanceService
	BillingHistory    BillingHistoryService
	CDNs              CDNService
	Certificates      CertificatesService
	Databases         DatabasesService
	Domains           DomainsService
	Droplets          DropletsService
	DropletActions    DropletActionsService
	DropletAutoscale  DropletAutoscaleService
	Firewalls         FirewallsService
	FloatingIPs       FloatingIPsService
	FloatingIPActions FloatingIPActionsService
	Functions         FunctionsService
	Images            ImagesService
	ImageActions      ImageActionsService
	Invoices          InvoicesService
	Keys              KeysService
	Kubernetes        KubernetesService
	LoadBalancers     LoadBalancersService
	Monitoring        MonitoringService
	OneClick          OneClickService
	Projects          ProjectsService
	Regions           RegionsService
	Registry          RegistryService
	Registries        RegistriesService
	ReservedIPs       ReservedIPsService
	ReservedIPActions ReservedIPActionsService
	Sizes             SizesService
	Snapshots         SnapshotsService
	Storage           StorageService
	StorageActions    StorageActionsService
	Tags              TagsService
	UptimeChecks      UptimeChecksService
	VPCs              VPCsService

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string

	// Optional rate limiter to ensure QoS.
	rateLimiter *rate.Limiter

	// Optional retry values. Setting the RetryConfig.RetryMax value enables automatically retrying requests
	// that fail with 429 or 500-level response codes using the go-retryablehttp client
	RetryConfig RetryConfig
}

// RetryConfig sets the values used for enabling retries and backoffs for
// requests that fail with 429 or 500-level response codes using the go-retryablehttp client.
// RetryConfig.RetryMax must be configured to enable this behavior. RetryConfig.RetryWaitMin and
// RetryConfig.RetryWaitMax are optional, with the default values being 1.0 and 30.0, respectively.
//
// You can use
//
//	godo.PtrTo(1.0)
//
// to explicitly set the RetryWaitMin and RetryWaitMax values.
//
// Note: Opting to use the go-retryablehttp client will overwrite any custom HTTP client passed into New().
// Only the oauth2.TokenSource and Timeout will be maintained.
type RetryConfig struct {
	RetryMax     int
	RetryWaitMin *float64    // Minimum time to wait
	RetryWaitMax *float64    // Maximum time to wait
	Logger       interface{} // Customer logger instance. Must implement either go-retryablehttp.Logger or go-retryablehttp.LeveledLogger
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`

	// Whether App responses should include project_id fields. The field will be empty if false or if omitted. (ListApps)
	WithProjects bool `url:"with_projects,omitempty"`
}

// TokenListOptions specifies the optional parameters to various List methods that support token pagination.
type TokenListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`

	// For paginated result sets which support tokens, the token provided by the last set
	// of results in order to retrieve the next set of results. This is expected to be faster
	// than incrementing or decrementing the page number.
	Token string `url:"page_token,omitempty"`
}

// Response is a DigitalOcean response. This wraps the standard http.Response returned from DigitalOcean.
type Response struct {
	*http.Response

	// Links that were returned with the response. These are parsed from
	// request body and not the header.
	Links *Links

	// Meta describes generic information about the response.
	Meta *Meta

	// Monitoring URI
	// Deprecated: This field is not populated. To poll for the status of a
	// newly created Droplet, use Links.Actions[0].HREF
	Monitor string

	Rate
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	// RequestID returned from the API, useful to contact support.
	RequestID string `json:"request_id"`

	// Attempts is the number of times the request was attempted when retries are enabled.
	Attempts int
}

// Rate contains the rate limit for the current client.
type Rate struct {
	// The number of request per hour the client is currently limited to.
	Limit int `json:"limit"`

	// The number of remaining requests the client can make this hour.
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset Timestamp `json:"reset"`
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// NewFromToken returns a new DigitalOcean API client with the given API
// token.
func NewFromToken(token string) *Client {
	cleanToken := strings.Trim(strings.TrimSpace(token), "'")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cleanToken})

	oauthClient := oauth2.NewClient(ctx, ts)
	client, err := New(oauthClient, WithRetryAndBackoffs(
		RetryConfig{
			RetryMax:     defaultRetryMax,
			RetryWaitMin: PtrTo(float64(defaultRetryWaitMin)),
			RetryWaitMax: PtrTo(float64(defaultRetryWaitMax)),
		},
	))
	if err != nil {
		panic(err)
	}

	return client
}

// NewClient returns a new DigitalOcean API client, using the given
// http.Client to perform all requests.
//
// Users who wish to pass their own http.Client should use this method. If
// you're in need of further customization, the godo.New method allows more
// options, such as setting a custom URL or a custom user agent string.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{HTTPClient: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.Account = &AccountServiceOp{client: c}
	c.Actions = &ActionsServiceOp{client: c}
	c.Apps = &AppsServiceOp{client: c}
	c.Balance = &BalanceServiceOp{client: c}
	c.BillingHistory = &BillingHistoryServiceOp{client: c}
	c.CDNs = &CDNServiceOp{client: c}
	c.Certificates = &CertificatesServiceOp{client: c}
	c.Databases = &DatabasesServiceOp{client: c}
	c.Domains = &DomainsServiceOp{client: c}
	c.Droplets = &DropletsServiceOp{client: c}
	c.DropletActions = &DropletActionsServiceOp{client: c}
	c.DropletAutoscale = &DropletAutoscaleServiceOp{client: c}
	c.Firewalls = &FirewallsServiceOp{client: c}
	c.FloatingIPs = &FloatingIPsServiceOp{client: c}
	c.FloatingIPActions = &FloatingIPActionsServiceOp{client: c}
	c.Functions = &FunctionsServiceOp{client: c}
	c.Images = &ImagesServiceOp{client: c}
	c.ImageActions = &ImageActionsServiceOp{client: c}
	c.Invoices = &InvoicesServiceOp{client: c}
	c.Keys = &KeysServiceOp{client: c}
	c.Kubernetes = &KubernetesServiceOp{client: c}
	c.LoadBalancers = &LoadBalancersServiceOp{client: c}
	c.Monitoring = &MonitoringServiceOp{client: c}
	c.OneClick = &OneClickServiceOp{client: c}
	c.Projects = &ProjectsServiceOp{client: c}
	c.Regions = &RegionsServiceOp{client: c}
	c.Registry = &RegistryServiceOp{client: c}
	c.Registries = &RegistriesServiceOp{client: c}
	c.ReservedIPs = &ReservedIPsServiceOp{client: c}
	c.ReservedIPActions = &ReservedIPActionsServiceOp{client: c}
	c.Sizes = &SizesServiceOp{client: c}
	c.Snapshots = &SnapshotsServiceOp{client: c}
	c.Storage = &StorageServiceOp{client: c}
	c.StorageActions = &StorageActionsServiceOp{client: c}
	c.Tags = &TagsServiceOp{client: c}
	c.UptimeChecks = &UptimeChecksServiceOp{client: c}
	c.VPCs = &VPCsServiceOp{client: c}

	c.headers = make(map[string]string)

	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new DigitalOcean API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// if retryMax is set it will use the retryablehttp client.
	if c.RetryConfig.RetryMax > 0 {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryMax = c.RetryConfig.RetryMax

		if c.RetryConfig.RetryWaitMin != nil {
			retryableClient.RetryWaitMin = time.Duration(*c.RetryConfig.RetryWaitMin * float64(time.Second))
		}
		if c.RetryConfig.RetryWaitMax != nil {
			retryableClient.RetryWaitMax = time.Duration(*c.RetryConfig.RetryWaitMax * float64(time.Second))
		}

		// By default this is nil and does not log.
		retryableClient.Logger = c.RetryConfig.Logger

		// if timeout is set, it is maintained before overwriting client with StandardClient()
		retryableClient.HTTPClient.Timeout = c.HTTPClient.Timeout

		// This custom ErrorHandler is required to provide errors that are consistent
		// with a *godo.ErrorResponse and a non-nil *godo.Response while providing
		// insight into retries using an internal header.
		retryableClient.ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
			if resp != nil {
				resp.Header.Add(internalHeaderRetryAttempts, strconv.Itoa(numTries))

				return resp, err
			}

			return resp, err
		}

		retryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			// In addition to the default retry policy, we also retry HTTP/2 INTERNAL_ERROR errors.
			// See: https://github.com/golang/go/issues/51323
			if err != nil && strings.Contains(err.Error(), "INTERNAL_ERROR") && strings.Contains(reflect.TypeOf(err).String(), "http2") {
				return true, nil
			}

			return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
		}

		var source *oauth2.Transport
		if _, ok := c.HTTPClient.Transport.(*oauth2.Transport); ok {
			source = c.HTTPClient.Transport.(*oauth2.Transport)
		}
		c.HTTPClient = retryableClient.StandardClient()
		c.HTTPClient.Transport = &oauth2.Transport{
			Base:   c.HTTPClient.Transport,
			Source: source.Source,
		}

	}

	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// SetRequestHeaders sets optional HTTP headers on the client that are
// sent on each HTTP request.
func SetRequestHeaders(headers map[string]string) ClientOpt {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

// SetStaticRateLimit sets an optional client-side rate limiter that restricts
// the number of queries per second that the client can send to enforce QoS.
func SetStaticRateLimit(rps float64) ClientOpt {
	return func(c *Client) error {
		c.rateLimiter = rate.NewLimiter(rate.Limit(rps), 1)
		return nil
	}
}

// WithRetryAndBackoffs sets retry values. Setting the RetryConfig.RetryMax value enables automatically retrying requests
// that fail with 429 or 500-level response codes using the go-retryablehttp client
func WithRetryAndBackoffs(retryConfig RetryConfig) ClientOpt {
	return func(c *Client) error {
		c.RetryConfig.RetryMax = retryConfig.RetryMax
		c.RetryConfig.RetryWaitMax = retryConfig.RetryWaitMax
		c.RetryConfig.RetryWaitMin = retryConfig.RetryWaitMin
		c.RetryConfig.Logger = retryConfig.Logger
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// GetRate returns the current rate limit for the client as determined by the most recent
// API call. It is thread-safe.
func (c *Client) GetRate() Rate {
	c.ratemtx.Lock()
	defer c.ratemtx.Unlock()
	return c.Rate
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	response.populateRate()

	return &response
}

// populateRate parses the rate related headers and populates the response Rate.
func (r *Response) populateRate() {
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		r.Rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		r.Rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			r.Rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if c.rateLimiter != nil {
		err := c.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	resp, err := DoRequestWithClient(ctx, c.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
		}

		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)
	c.ratemtx.Lock()
	c.Rate = response.Rate
	c.ratemtx.Unlock()

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != http.StatusNoContent && v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	var attempted string
	if r.Attempts > 0 {
		attempted = fmt.Sprintf("; giving up after %d attempt(s)", r.Attempts)
	}

	if r.RequestID != "" {
		return fmt.Sprintf("%v %v: %d (request %q) %v%s",
			r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.RequestID, r.Message, attempted)
	}
	return fmt.Sprintf("%v %v: %d %v%s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, attempted)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
// If the API error response does not include the request ID in its body, the one from its header will be used.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	if errorResponse.RequestID == "" {
		errorResponse.RequestID = r.Header.Get(headerRequestID)
	}

	attempts, strconvErr := strconv.Atoi(r.Header.Get(internalHeaderRetryAttempts))
	if strconvErr == nil {
		errorResponse.Attempts = attempts
	}

	return errorResponse
}

func (r Rate) String() string {
	return Stringify(r)
}

// PtrTo returns a pointer to the provided input.
func PtrTo[T any](v T) *T {
	return &v
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
//
// Deprecated: Use PtrTo instead.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
//
// Deprecated: Use PtrTo instead.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
//
// Deprecated: Use PtrTo instead.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}
