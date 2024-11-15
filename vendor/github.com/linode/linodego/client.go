package linodego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// APIConfigEnvVar environment var to get path to Linode config
	APIConfigEnvVar = "LINODE_CONFIG"
	// APIConfigProfileEnvVar specifies the profile to use when loading from a Linode config
	APIConfigProfileEnvVar = "LINODE_PROFILE"
	// APIHost Linode API hostname
	APIHost = "api.linode.com"
	// APIHostVar environment var to check for alternate API URL
	APIHostVar = "LINODE_URL"
	// APIHostCert environment var containing path to CA cert to validate against
	APIHostCert = "LINODE_CA"
	// APIVersion Linode API version
	APIVersion = "v4"
	// APIVersionVar environment var to check for alternate API Version
	APIVersionVar = "LINODE_API_VERSION"
	// APIProto connect to API with http(s)
	APIProto = "https"
	// APIEnvVar environment var to check for API token
	APIEnvVar = "LINODE_TOKEN"
	// APISecondsPerPoll how frequently to poll for new Events or Status in WaitFor functions
	APISecondsPerPoll = 3
	// Maximum wait time for retries
	APIRetryMaxWaitTime       = time.Duration(30) * time.Second
	APIDefaultCacheExpiration = time.Minute * 15
)

//nolint:unused
var (
	reqLogTemplate = template.Must(template.New("request").Parse(`Sending request:
Method: {{.Method}}
URL: {{.URL}}
Headers: {{.Headers}}
Body: {{.Body}}`))

	respLogTemplate = template.Must(template.New("response").Parse(`Received response:
Status: {{.Status}}
Headers: {{.Headers}}
Body: {{.Body}}`))
)

var envDebug = false

// Client is a wrapper around the Resty client
type Client struct {
	resty             *resty.Client
	userAgent         string
	debug             bool
	retryConditionals []RetryConditional

	pollInterval time.Duration

	baseURL         string
	apiVersion      string
	apiProto        string
	selectedProfile string
	loadedProfile   string

	configProfiles map[string]ConfigProfile

	// Fields for caching endpoint responses
	shouldCache     bool
	cacheExpiration time.Duration
	cachedEntries   map[string]clientCacheEntry
	cachedEntryLock *sync.RWMutex
}

type EnvDefaults struct {
	Token   string
	Profile string
}

type clientCacheEntry struct {
	Created time.Time
	Data    any
	// If != nil, use this instead of the
	// global expiry
	ExpiryOverride *time.Duration
}

type (
	Request  = resty.Request
	Response = resty.Response
	Logger   = resty.Logger
)

func init() {
	// Whether we will enable Resty debugging output
	if apiDebug, ok := os.LookupEnv("LINODE_DEBUG"); ok {
		if parsed, err := strconv.ParseBool(apiDebug); err == nil {
			envDebug = parsed
			log.Println("[INFO] LINODE_DEBUG being set to", envDebug)
		} else {
			log.Println("[WARN] LINODE_DEBUG should be an integer, 0 or 1")
		}
	}
}

// SetUserAgent sets a custom user-agent for HTTP requests
func (c *Client) SetUserAgent(ua string) *Client {
	c.userAgent = ua
	c.resty.SetHeader("User-Agent", c.userAgent)

	return c
}

type RequestParams struct {
	Body     any
	Response any
}

// Generic helper to execute HTTP requests using the net/http package
//
// nolint:unused, funlen, gocognit
func (c *httpClient) doRequest(ctx context.Context, method, url string, params RequestParams) error {
	var (
		req        *http.Request
		bodyBuffer *bytes.Buffer
		resp       *http.Response
		err        error
	)

	for range httpDefaultRetryCount {
		req, bodyBuffer, err = c.createRequest(ctx, method, url, params)
		if err != nil {
			return err
		}

		if err = c.applyBeforeRequest(req); err != nil {
			return err
		}

		if c.debug && c.logger != nil {
			c.logRequest(req, method, url, bodyBuffer)
		}

		processResponse := func() error {
			defer func() {
				closeErr := resp.Body.Close()
				if closeErr != nil && err == nil {
					err = closeErr
				}
			}()
			if err = c.checkHTTPError(resp); err != nil {
				return err
			}
			if c.debug && c.logger != nil {
				var logErr error
				resp, logErr = c.logResponse(resp)
				if logErr != nil {
					return logErr
				}
			}
			if params.Response != nil {
				if err = c.decodeResponseBody(resp, params.Response); err != nil {
					return err
				}
			}

			// Apply after-response mutations
			if err = c.applyAfterResponse(resp); err != nil {
				return err
			}

			return nil
		}

		resp, err = c.sendRequest(req)
		if err == nil {
			if err = processResponse(); err == nil {
				return nil
			}
		}

		if !c.shouldRetry(resp, err) {
			break
		}

		retryAfter, retryErr := c.retryAfter(resp)
		if retryErr != nil {
			return retryErr
		}

		// Sleep for the specified duration before retrying.
		// If retryAfter is 0 (i.e., Retry-After header is not found),
		// no delay is applied.
		time.Sleep(retryAfter)
	}

	return err
}

// nolint:unused
func (c *httpClient) shouldRetry(resp *http.Response, err error) bool {
	for _, retryConditional := range c.retryConditionals {
		if retryConditional(resp, err) {
			return true
		}
	}
	return false
}

// nolint:unused
func (c *httpClient) createRequest(ctx context.Context, method, url string, params RequestParams) (*http.Request, *bytes.Buffer, error) {
	var bodyReader io.Reader
	var bodyBuffer *bytes.Buffer

	if params.Body != nil {
		bodyBuffer = new(bytes.Buffer)
		if err := json.NewEncoder(bodyBuffer).Encode(params.Body); err != nil {
			if c.debug && c.logger != nil {
				c.logger.Errorf("failed to encode body: %v", err)
			}
			return nil, nil, fmt.Errorf("failed to encode body: %w", err)
		}
		bodyReader = bodyBuffer
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		if c.debug && c.logger != nil {
			c.logger.Errorf("failed to create request: %v", err)
		}
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, bodyBuffer, nil
}

// nolint:unused
func (c *httpClient) applyBeforeRequest(req *http.Request) error {
	for _, mutate := range c.onBeforeRequest {
		if err := mutate(req); err != nil {
			if c.debug && c.logger != nil {
				c.logger.Errorf("failed to mutate before request: %v", err)
			}
			return fmt.Errorf("failed to mutate before request: %w", err)
		}
	}
	return nil
}

// nolint:unused
func (c *httpClient) applyAfterResponse(resp *http.Response) error {
	for _, mutate := range c.onAfterResponse {
		if err := mutate(resp); err != nil {
			if c.debug && c.logger != nil {
				c.logger.Errorf("failed to mutate after response: %v", err)
			}
			return fmt.Errorf("failed to mutate after response: %w", err)
		}
	}
	return nil
}

// nolint:unused
func (c *httpClient) logRequest(req *http.Request, method, url string, bodyBuffer *bytes.Buffer) {
	var reqBody string
	if bodyBuffer != nil {
		reqBody = bodyBuffer.String()
	} else {
		reqBody = "nil"
	}

	var logBuf bytes.Buffer
	err := reqLogTemplate.Execute(&logBuf, map[string]interface{}{
		"Method":  method,
		"URL":     url,
		"Headers": req.Header,
		"Body":    reqBody,
	})
	if err == nil {
		c.logger.Debugf(logBuf.String())
	}
}

// nolint:unused
func (c *httpClient) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.debug && c.logger != nil {
			c.logger.Errorf("failed to send request: %v", err)
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return resp, nil
}

// nolint:unused
func (c *httpClient) checkHTTPError(resp *http.Response) error {
	_, err := coupleAPIErrorsHTTP(resp, nil)
	if err != nil {
		if c.debug && c.logger != nil {
			c.logger.Errorf("received HTTP error: %v", err)
		}
		return err
	}
	return nil
}

// nolint:unused
func (c *httpClient) logResponse(resp *http.Response) (*http.Response, error) {
	var respBody bytes.Buffer
	if _, err := io.Copy(&respBody, resp.Body); err != nil {
		c.logger.Errorf("failed to read response body: %v", err)
	}

	var logBuf bytes.Buffer
	err := respLogTemplate.Execute(&logBuf, map[string]interface{}{
		"Status":  resp.Status,
		"Headers": resp.Header,
		"Body":    respBody.String(),
	})
	if err == nil {
		c.logger.Debugf(logBuf.String())
	}

	resp.Body = io.NopCloser(bytes.NewReader(respBody.Bytes()))
	return resp, nil
}

// nolint:unused
func (c *httpClient) decodeResponseBody(resp *http.Response, response interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		if c.debug && c.logger != nil {
			c.logger.Errorf("failed to decode response: %v", err)
		}
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

// R wraps resty's R method
func (c *Client) R(ctx context.Context) *resty.Request {
	return c.resty.R().
		ExpectContentType("application/json").
		SetHeader("Content-Type", "application/json").
		SetContext(ctx).
		SetError(APIError{})
}

// SetDebug sets the debug on resty's client
func (c *Client) SetDebug(debug bool) *Client {
	c.debug = debug
	c.resty.SetDebug(debug)

	return c
}

// SetLogger allows the user to override the output
// logger for debug logs.
func (c *Client) SetLogger(logger Logger) *Client {
	c.resty.SetLogger(logger)

	return c
}

//nolint:unused
func (c *httpClient) httpSetDebug(debug bool) *httpClient {
	c.debug = debug

	return c
}

//nolint:unused
func (c *httpClient) httpSetLogger(logger httpLogger) *httpClient {
	c.logger = logger

	return c
}

// OnBeforeRequest adds a handler to the request body to run before the request is sent
func (c *Client) OnBeforeRequest(m func(request *Request) error) {
	c.resty.OnBeforeRequest(func(_ *resty.Client, req *resty.Request) error {
		return m(req)
	})
}

// OnAfterResponse adds a handler to the request body to run before the request is sent
func (c *Client) OnAfterResponse(m func(response *Response) error) {
	c.resty.OnAfterResponse(func(_ *resty.Client, req *resty.Response) error {
		return m(req)
	})
}

// nolint:unused
func (c *httpClient) httpOnBeforeRequest(m func(*http.Request) error) *httpClient {
	c.onBeforeRequest = append(c.onBeforeRequest, m)

	return c
}

// nolint:unused
func (c *httpClient) httpOnAfterResponse(m func(*http.Response) error) *httpClient {
	c.onAfterResponse = append(c.onAfterResponse, m)

	return c
}

// UseURL parses the individual components of the given API URL and configures the client
// accordingly. For example, a valid URL.
// For example:
//
//	client.UseURL("https://api.test.linode.com/v4beta")
func (c *Client) UseURL(apiURL string) (*Client, error) {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Create a new URL excluding the path to use as the base URL
	baseURL := &url.URL{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
	}

	c.SetBaseURL(baseURL.String())

	versionMatches := regexp.MustCompile(`/v[a-zA-Z0-9]+`).FindAllString(parsedURL.Path, -1)

	// Only set the version if a version is found in the URL, else use the default
	if len(versionMatches) > 0 {
		c.SetAPIVersion(
			strings.Trim(versionMatches[len(versionMatches)-1], "/"),
		)
	}

	return c, nil
}

// SetBaseURL sets the base URL of the Linode v4 API (https://api.linode.com/v4)
func (c *Client) SetBaseURL(baseURL string) *Client {
	baseURLPath, _ := url.Parse(baseURL)

	c.baseURL = path.Join(baseURLPath.Host, baseURLPath.Path)
	c.apiProto = baseURLPath.Scheme

	c.updateHostURL()

	return c
}

// SetAPIVersion sets the version of the API to interface with
func (c *Client) SetAPIVersion(apiVersion string) *Client {
	c.apiVersion = apiVersion

	c.updateHostURL()

	return c
}

func (c *Client) updateHostURL() {
	apiProto := APIProto
	baseURL := APIHost
	apiVersion := APIVersion

	if c.baseURL != "" {
		baseURL = c.baseURL
	}

	if c.apiVersion != "" {
		apiVersion = c.apiVersion
	}

	if c.apiProto != "" {
		apiProto = c.apiProto
	}

	c.resty.SetBaseURL(
		fmt.Sprintf(
			"%s://%s/%s",
			apiProto,
			baseURL,
			url.PathEscape(apiVersion),
		),
	)
}

// SetRootCertificate adds a root certificate to the underlying TLS client config
func (c *Client) SetRootCertificate(path string) *Client {
	c.resty.SetRootCertificate(path)
	return c
}

// SetToken sets the API token for all requests from this client
// Only necessary if you haven't already provided the http client to NewClient() configured with the token.
func (c *Client) SetToken(token string) *Client {
	c.resty.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return c
}

// SetRetries adds retry conditions for "Linode Busy." errors and 429s.
func (c *Client) SetRetries() *Client {
	c.
		addRetryConditional(linodeBusyRetryCondition).
		addRetryConditional(tooManyRequestsRetryCondition).
		addRetryConditional(serviceUnavailableRetryCondition).
		addRetryConditional(requestTimeoutRetryCondition).
		addRetryConditional(requestGOAWAYRetryCondition).
		addRetryConditional(requestNGINXRetryCondition).
		SetRetryMaxWaitTime(APIRetryMaxWaitTime)
	configureRetries(c)
	return c
}

// AddRetryCondition adds a RetryConditional function to the Client
func (c *Client) AddRetryCondition(retryCondition RetryConditional) *Client {
	c.resty.AddRetryCondition(resty.RetryConditionFunc(retryCondition))
	return c
}

func (c *Client) addRetryConditional(retryConditional RetryConditional) *Client {
	c.retryConditionals = append(c.retryConditionals, retryConditional)
	return c
}

func (c *Client) addCachedResponse(endpoint string, response any, expiry *time.Duration) {
	if !c.shouldCache {
		return
	}

	responseValue := reflect.ValueOf(response)

	entry := clientCacheEntry{
		Created:        time.Now(),
		ExpiryOverride: expiry,
	}

	switch responseValue.Kind() {
	case reflect.Ptr:
		// We want to automatically deref pointers to
		// avoid caching mutable data.
		entry.Data = responseValue.Elem().Interface()
	default:
		entry.Data = response
	}

	c.cachedEntryLock.Lock()
	defer c.cachedEntryLock.Unlock()

	c.cachedEntries[endpoint] = entry
}

func (c *Client) getCachedResponse(endpoint string) any {
	if !c.shouldCache {
		return nil
	}

	c.cachedEntryLock.RLock()

	// Hacky logic to dynamically RUnlock
	// only if it is still locked by the
	// end of the function.
	// This is necessary as we take write
	// access if the entry has expired.
	rLocked := true
	defer func() {
		if rLocked {
			c.cachedEntryLock.RUnlock()
		}
	}()

	entry, ok := c.cachedEntries[endpoint]
	if !ok {
		return nil
	}

	// Handle expired entries
	elapsedTime := time.Since(entry.Created)

	hasExpired := elapsedTime > c.cacheExpiration
	if entry.ExpiryOverride != nil {
		hasExpired = elapsedTime > *entry.ExpiryOverride
	}

	if hasExpired {
		// We need to give up our read access and request read-write access
		c.cachedEntryLock.RUnlock()
		rLocked = false

		c.cachedEntryLock.Lock()
		defer c.cachedEntryLock.Unlock()

		delete(c.cachedEntries, endpoint)
		return nil
	}

	return c.cachedEntries[endpoint].Data
}

// InvalidateCache clears all cached responses for all endpoints.
func (c *Client) InvalidateCache() {
	c.cachedEntryLock.Lock()
	defer c.cachedEntryLock.Unlock()

	// GC will handle the old map
	c.cachedEntries = make(map[string]clientCacheEntry)
}

// InvalidateCacheEndpoint invalidates a single cached endpoint.
func (c *Client) InvalidateCacheEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse URL for caching: %w", err)
	}

	c.cachedEntryLock.Lock()
	defer c.cachedEntryLock.Unlock()

	delete(c.cachedEntries, u.Path)

	return nil
}

// SetGlobalCacheExpiration sets the desired time for any cached response
// to be valid for.
func (c *Client) SetGlobalCacheExpiration(expiryTime time.Duration) {
	c.cacheExpiration = expiryTime
}

// UseCache sets whether response caching should be used
func (c *Client) UseCache(value bool) {
	c.shouldCache = value
}

// SetRetryMaxWaitTime sets the maximum delay before retrying a request.
func (c *Client) SetRetryMaxWaitTime(maxWaitTime time.Duration) *Client {
	c.resty.SetRetryMaxWaitTime(maxWaitTime)
	return c
}

// SetRetryWaitTime sets the default (minimum) delay before retrying a request.
func (c *Client) SetRetryWaitTime(minWaitTime time.Duration) *Client {
	c.resty.SetRetryWaitTime(minWaitTime)
	return c
}

// SetRetryAfter sets the callback function to be invoked with a failed request
// to determine wben it should be retried.
func (c *Client) SetRetryAfter(callback RetryAfter) *Client {
	c.resty.SetRetryAfter(resty.RetryAfterFunc(callback))
	return c
}

// SetRetryCount sets the maximum retry attempts before aborting.
func (c *Client) SetRetryCount(count int) *Client {
	c.resty.SetRetryCount(count)
	return c
}

// SetPollDelay sets the number of milliseconds to wait between events or status polls.
// Affects all WaitFor* functions and retries.
func (c *Client) SetPollDelay(delay time.Duration) *Client {
	c.pollInterval = delay
	return c
}

// GetPollDelay gets the number of milliseconds to wait between events or status polls.
// Affects all WaitFor* functions and retries.
func (c *Client) GetPollDelay() time.Duration {
	return c.pollInterval
}

// SetHeader sets a custom header to be used in all API requests made with the current
// client.
// NOTE: Some headers may be overridden by the individual request functions.
func (c *Client) SetHeader(name, value string) {
	c.resty.SetHeader(name, value)
}

func (c *Client) enableLogSanitization() *Client {
	c.resty.OnRequestLog(func(r *resty.RequestLog) error {
		// masking authorization header
		r.Header.Set("Authorization", "Bearer *******************************")
		return nil
	})

	return c
}

// NewClient factory to create new Client struct
func NewClient(hc *http.Client) (client Client) {
	if hc != nil {
		client.resty = resty.NewWithClient(hc)
	} else {
		client.resty = resty.New()
	}

	client.shouldCache = true
	client.cacheExpiration = APIDefaultCacheExpiration
	client.cachedEntries = make(map[string]clientCacheEntry)
	client.cachedEntryLock = &sync.RWMutex{}

	client.SetUserAgent(DefaultUserAgent)

	baseURL, baseURLExists := os.LookupEnv(APIHostVar)

	if baseURLExists {
		client.SetBaseURL(baseURL)
	}
	apiVersion, apiVersionExists := os.LookupEnv(APIVersionVar)
	if apiVersionExists {
		client.SetAPIVersion(apiVersion)
	} else {
		client.SetAPIVersion(APIVersion)
	}

	certPath, certPathExists := os.LookupEnv(APIHostCert)

	if certPathExists {
		cert, err := os.ReadFile(filepath.Clean(certPath))
		if err != nil {
			log.Fatalf("[ERROR] Error when reading cert at %s: %s\n", certPath, err.Error())
		}

		client.SetRootCertificate(certPath)

		if envDebug {
			log.Printf("[DEBUG] Set API root certificate to %s with contents %s\n", certPath, cert)
		}
	}

	client.
		SetRetryWaitTime(APISecondsPerPoll * time.Second).
		SetPollDelay(APISecondsPerPoll * time.Second).
		SetRetries().
		SetDebug(envDebug).
		enableLogSanitization()

	return
}

// NewClientFromEnv creates a Client and initializes it with values
// from the LINODE_CONFIG file and the LINODE_TOKEN environment variable.
func NewClientFromEnv(hc *http.Client) (*Client, error) {
	client := NewClient(hc)

	// Users are expected to chain NewClient(...) and LoadConfig(...) to customize these options
	configPath, err := resolveValidConfigPath()
	if err != nil {
		return nil, err
	}

	// Populate the token from the environment.
	// Tokens should be first priority to maintain backwards compatibility
	if token, ok := os.LookupEnv(APIEnvVar); ok && token != "" {
		client.SetToken(token)
		return &client, nil
	}

	if p, ok := os.LookupEnv(APIConfigEnvVar); ok {
		configPath = p
	} else if !ok && configPath == "" {
		return nil, fmt.Errorf("no linode config file or token found")
	}

	configProfile := DefaultConfigProfile

	if p, ok := os.LookupEnv(APIConfigProfileEnvVar); ok {
		configProfile = p
	}

	client.selectedProfile = configProfile

	// We should only load the config if the config file exists
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("error loading config file %s: %w", configPath, err)
	}

	err = client.preLoadConfig(configPath)
	return &client, err
}

func (c *Client) preLoadConfig(configPath string) error {
	if envDebug {
		log.Printf("[INFO] Loading profile from %s\n", configPath)
	}

	if err := c.LoadConfig(&LoadConfigOptions{
		Path:            configPath,
		SkipLoadProfile: true,
	}); err != nil {
		return err
	}

	// We don't want to load the profile until the user is actually making requests
	c.OnBeforeRequest(func(_ *Request) error {
		if c.loadedProfile != c.selectedProfile {
			if err := c.UseProfile(c.selectedProfile); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func copyBool(bPtr *bool) *bool {
	if bPtr == nil {
		return nil
	}

	t := *bPtr

	return &t
}

func copyInt(iPtr *int) *int {
	if iPtr == nil {
		return nil
	}

	t := *iPtr

	return &t
}

func copyString(sPtr *string) *string {
	if sPtr == nil {
		return nil
	}

	t := *sPtr

	return &t
}

func copyTime(tPtr *time.Time) *time.Time {
	if tPtr == nil {
		return nil
	}

	t := *tPtr

	return &t
}

func generateListCacheURL(endpoint string, opts *ListOptions) (string, error) {
	if opts == nil {
		return endpoint, nil
	}

	hashedOpts, err := opts.Hash()
	if err != nil {
		return endpoint, err
	}

	return fmt.Sprintf("%s:%s", endpoint, hashedOpts), nil
}
