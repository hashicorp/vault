package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-querystring/query"

	"reflect"
)

const (
	libraryVersion            = "1"
	userAgent                 = "oktasdk-go/" + libraryVersion
	productionDomain          = "okta.com"
	previewDomain             = "oktapreview.com"
	urlFormat                 = "https://%s.%s/api/v1/"
	headerRateLimit           = "X-Rate-Limit-Limit"
	headerRateRemaining       = "X-Rate-Limit-Remaining"
	headerRateReset           = "X-Rate-Limit-Reset"
	headerOKTARequestID       = "X-Okta-Request-Id"
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "SSWS %v"
	mediaTypeJSON             = "application/json"
	defaultLimit              = 50
	// FilterEqualOperator Filter Operatorid for "equal"
	FilterEqualOperator = "eq"
	// FilterStartsWithOperator - filter operator for "starts with"
	FilterStartsWithOperator = "sw"
	// FilterGreaterThanOperator - filter operator for "greater than"
	FilterGreaterThanOperator = "gt"
	// FilterLessThanOperator - filter operator for "less than"
	FilterLessThanOperator = "lt"

	// If the API returns a "X-Rate-Limit-Remaining" header less than this the SDK will either pause
	//  Or throw  RateLimitError depending on the client.PauseOnRateLimit value
	defaultRateRemainingFloor = 100
)

// A Client manages communication with the API.
type Client struct {
	clientMu sync.Mutex   // clientMu protects the client during calls that modify the CheckRedirect func.
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests.
	//  This will be built automatically based on inputs to NewClient
	//  If needed you can override this if needed (your URL is not *.okta.com or *.oktapreview.com)
	BaseURL *url.URL

	// User agent used when communicating with the GitHub API.
	UserAgent string

	apiKey                   string
	authorizationHeaderValue string
	PauseOnRateLimit         bool

	// RateRemainingFloor - If the API returns a "X-Rate-Limit-Remaining" header less than this the SDK will either pause
	//  Or throw  RateLimitError depending on the client.PauseOnRateLimit value. It defaults to 30
	// One client doing too much work can lock out all API Access for every other client
	// We are trying to be a "good API User Citizen"
	RateRemainingFloor int

	rateMu         sync.Mutex
	mostRecentRate Rate

	Limit int
	// mostRecent rateLimitCategory

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the  API.
	// Service for Working with Users
	Users *UsersService

	// Service for Working with Groups
	Groups *GroupsService

	// Service for Working with Apps
	Apps *AppsService
}

type service struct {
	client *Client
}

// NewClient returns a new OKTA API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client, orgName string, apiToken string, isProduction bool) *Client {
	var baseDomain string
	if isProduction {
		baseDomain = productionDomain
	} else {
		baseDomain = previewDomain
	}
	client, _ := NewClientWithDomain(httpClient, orgName, baseDomain, apiToken)
	return client
}

// NewClientWithDomain creates a client based on the organziation name and
// base domain for requests (okta.com, okta-emea.com, oktapreview.com, etc).
func NewClientWithDomain(httpClient *http.Client, orgName string, domain string, apiToken string) (*Client, error) {
	baseURL, err := url.Parse(fmt.Sprintf(urlFormat, orgName, domain))
	if err != nil {
		return nil, err
	}
	return NewClientWithBaseURL(httpClient, baseURL, apiToken), nil
}

// NewClientWithBaseURL creates a client based on the full base URL and api
// token
func NewClientWithBaseURL(httpClient *http.Client, baseURL *url.URL, apiToken string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.PauseOnRateLimit = true // If rate limit found it will block until that time. If false then Error will be returned
	c.authorizationHeaderValue = fmt.Sprintf(headerAuthorizationFormat, apiToken)
	c.apiKey = apiToken
	c.Limit = defaultLimit
	c.RateRemainingFloor = defaultRateRemainingFloor
	c.common.client = c

	c.Users = (*UsersService)(&c.common)
	c.Groups = (*GroupsService)(&c.common)
	c.Apps = (*AppsService)(&c.common)
	return c
}

// Rate represents the rate limit for the current client.
type Rate struct {
	// The number of requests per minute the client is currently limited to.
	RatePerMinuteLimit int

	// The number of remaining requests the client can make this minute
	Remaining int

	// The time at which the current rate limit will reset.
	ResetTime time.Time
}

// Response is a OKTA API response.  This wraps the standard http.Response
// returned from OKTA and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results.

	NextURL *url.URL
	// PrevURL       *url.URL
	SelfURL       *url.URL
	OKTARequestID string
	Rate
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}

	response.OKTARequestID = r.Header.Get(headerOKTARequestID)

	response.populatePaginationURLS()
	response.Rate = parseRate(r)
	return response
}

// populatePageValues parses the HTTP Link response headers and populates the
// various pagination link values in the Response.

// OKTA LINK Header takes this form:
// 		Link: <https://yoursubdomain.okta.com/api/v1/users?after=00ubfjQEMYBLRUWIEDKK>; rel="next",
// 			<https://yoursubdomain.okta.com/api/v1/users?after=00ub4tTFYKXCCZJSGFKM>; rel="self"

func (r *Response) populatePaginationURLS() {

	for k, v := range r.Header {

		if k == "Link" {
			nextRegex := regexp.MustCompile(`<(.*?)>; rel="next"`)
			// prevRegex := regexp.MustCompile(`<(.*?)>; rel="prev"`)
			selfRegex := regexp.MustCompile(`<(.*?)>; rel="self"`)

			for _, linkValue := range v {
				nextLinkMatch := nextRegex.FindStringSubmatch(linkValue)
				if len(nextLinkMatch) != 0 {
					r.NextURL, _ = url.Parse(nextLinkMatch[1])
				}
				selfLinkMatch := selfRegex.FindStringSubmatch(linkValue)
				if len(selfLinkMatch) != 0 {
					r.SelfURL, _ = url.Parse(selfLinkMatch[1])
				}
				// prevLinkMatch := prevRegex.FindStringSubmatch(linkValue)
				// if len(prevLinkMatch) != 0 {
				// 	r.PrevURL, _ = url.Parse(prevLinkMatch[1])
				// }
			}
		}
	}

}

// parseRate parses the rate related headers.
func parseRate(r *http.Response) Rate {
	var rate Rate

	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.RatePerMinuteLimit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.ResetTime = time.Unix(v, 0)
		}
	}
	return rate
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.  If rate limit is exceeded and reset time is in the future,
// Do returns rate immediately without making a network API call.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {

	// If we've hit rate limit, don't make further requests before Reset time.
	if err := c.checkRateLimitBeforeDo(req); err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := newResponse(resp)

	c.rateMu.Lock()
	c.mostRecentRate.RatePerMinuteLimit = response.Rate.RatePerMinuteLimit
	c.mostRecentRate.Remaining = response.Rate.Remaining
	c.mostRecentRate.ResetTime = response.Rate.ResetTime
	c.rateMu.Unlock()

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		// fmt.Printf("Error after sdk.Do return\n")

		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// checkRateLimitBeforeDo does not make any network calls, but uses existing knowledge from
// current client state in order to quickly check if *RateLimitError can be immediately returned
// from Client.Do, and if so, returns it so that Client.Do can skip making a network API call unnecessarily.
// Otherwise it returns nil, and Client.Do should proceed normally.
// http://developer.okta.com/docs/api/getting_started/design_principles.html#rate-limiting
func (c *Client) checkRateLimitBeforeDo(req *http.Request) error {

	c.rateMu.Lock()
	mostRecentRate := c.mostRecentRate
	c.rateMu.Unlock()
	// fmt.Printf("checkRateLimitBeforeDo: \t Remaining = %d, \t ResetTime = %s\n", mostRecentRate.Remaining, mostRecentRate.ResetTime.String())
	if !mostRecentRate.ResetTime.IsZero() && mostRecentRate.Remaining < c.RateRemainingFloor && time.Now().Before(mostRecentRate.ResetTime) {

		if c.PauseOnRateLimit {
			// If rate limit is hitting threshold then pause until the rate limit resets
			//   This behavior is controlled by the client PauseOnRateLimit value
			// fmt.Printf("checkRateLimitBeforeDo: \t ***pause**** \t Time Now = %s \tPause After = %s\n", time.Now().String(), mostRecentRate.ResetTime.Sub(time.Now().Add(2*time.Second)).String())
			<-time.After(mostRecentRate.ResetTime.Sub(time.Now().Add(2 * time.Second)))
		} else {
			// fmt.Printf("checkRateLimitBeforeDo: \t ***error****\n")

			return &RateLimitError{
				Rate: mostRecentRate,
			}
		}

	}

	return nil
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// and *TwoFactorAuthError for two-factor authentication errors.
// TODO - check un-authorized
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResp := &errorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, &errorResp.ErrorDetail)
	}
	switch {
	case r.StatusCode == http.StatusTooManyRequests:

		return &RateLimitError{
			Rate:        parseRate(r),
			Response:    r,
			ErrorDetail: errorResp.ErrorDetail}

	default:
		return errorResp
	}

}

type apiError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorLink    string `json:"errorLink"`
	ErrorID      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
}

type errorResponse struct {
	Response    *http.Response //
	ErrorDetail apiError
}

func (r *errorResponse) Error() string {
	return fmt.Sprintf("HTTP Method: %v - URL: %v: - HTTP Status Code: %d, OKTA Error Code: %v, OKTA Error Summary: %v, OKTA Error Causes: %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.ErrorDetail.ErrorCode, r.ErrorDetail.ErrorSummary, r.ErrorDetail.ErrorCauses)
}

// RateLimitError occurs when OKTA returns 429 "Too Many Requests" response with a rate limit
// remaining value of 0, and error message starts with "API rate limit exceeded for ".
type RateLimitError struct {
	Rate        Rate // Rate specifies last known rate limit for the client
	ErrorDetail apiError
	Response    *http.Response //
}

func (r *RateLimitError) Error() string {

	return fmt.Sprintf("rate reset in %v", r.Rate.ResetTime.Sub(time.Now()))

}

// Code stolen from Github api libary
// Stringify attempts to create a reasonable string representation of types in
// the library.  It does things like resolve pointers to their values
// and omits struct fields with nil values.
func stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

// stringifyValue was heavily inspired by the goprotobuf library.

func stringifyValue(w io.Writer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i))
		}

		w.Write([]byte{']'})
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}
		w.Write([]byte{'{'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(", "))
			} else {
				sep = true
			}

			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv)
		}

		w.Write([]byte{'}'})
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerAuthorization, fmt.Sprintf(headerAuthorizationFormat, c.apiKey))

	if body != nil {
		req.Header.Set("Content-Type", mediaTypeJSON)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

type dateFilter struct {
	Value    time.Time
	Operator string
}
