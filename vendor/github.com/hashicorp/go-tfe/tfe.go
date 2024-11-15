// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/jsonapi"
	"golang.org/x/time/rate"

	slug "github.com/hashicorp/go-slug"
)

const (
	_userAgent         = "go-tfe"
	_headerRateLimit   = "X-RateLimit-Limit"
	_headerRateReset   = "X-RateLimit-Reset"
	_headerAppName     = "TFP-AppName"
	_headerAPIVersion  = "TFP-API-Version"
	_headerTFEVersion  = "X-TFE-Version"
	_includeQueryParam = "include"

	DefaultAddress      = "https://app.terraform.io"
	DefaultBasePath     = "/api/v2/"
	DefaultRegistryPath = "/api/registry/"
	// PingEndpoint is a no-op API endpoint used to configure the rate limiter
	PingEndpoint       = "ping"
	ContentTypeJSONAPI = "application/vnd.api+json"
)

// RetryLogHook allows a function to run before each retry.

type RetryLogHook func(attemptNum int, resp *http.Response)

// Config provides configuration details to the API client.

type Config struct {
	// The address of the Terraform Enterprise API.
	Address string

	// The base path on which the API is served.
	BasePath string

	// The base path for the Registry API
	RegistryBasePath string

	// API token used to access the Terraform Enterprise API.
	Token string

	// Headers that will be added to every request.
	Headers http.Header

	// A custom HTTP client to use.
	HTTPClient *http.Client

	// RetryLogHook is invoked each time a request is retried.
	RetryLogHook RetryLogHook

	// RetryServerErrors enables the retry logic in the client.
	RetryServerErrors bool
}

// DefaultConfig returns a default config structure.

func DefaultConfig() *Config {
	config := &Config{
		Address:           os.Getenv("TFE_ADDRESS"),
		BasePath:          DefaultBasePath,
		RegistryBasePath:  DefaultRegistryPath,
		Token:             os.Getenv("TFE_TOKEN"),
		Headers:           make(http.Header),
		HTTPClient:        cleanhttp.DefaultPooledClient(),
		RetryServerErrors: false,
	}

	// Set the default address if none is given.
	if config.Address == "" {
		if host := os.Getenv("TFE_HOSTNAME"); host != "" {
			config.Address = fmt.Sprintf("https://%s", host)
		} else {
			config.Address = DefaultAddress
		}
	}

	// Set the default user agent.
	config.Headers.Set("User-Agent", _userAgent)

	return config
}

// Client is the Terraform Enterprise API client. It provides the basic
// connectivity and configuration for accessing the TFE API
type Client struct {
	baseURL           *url.URL
	registryBaseURL   *url.URL
	token             string
	headers           http.Header
	http              *retryablehttp.Client
	limiter           *rate.Limiter
	retryLogHook      RetryLogHook
	retryServerErrors bool
	remoteAPIVersion  string
	remoteTFEVersion  string
	appName           string

	Admin                      Admin
	Agents                     Agents
	AgentPools                 AgentPools
	AgentTokens                AgentTokens
	Applies                    Applies
	AuditTrails                AuditTrails
	Comments                   Comments
	ConfigurationVersions      ConfigurationVersions
	CostEstimates              CostEstimates
	GHAInstallations           GHAInstallations
	GPGKeys                    GPGKeys
	NotificationConfigurations NotificationConfigurations
	OAuthClients               OAuthClients
	OAuthTokens                OAuthTokens
	Organizations              Organizations
	OrganizationMemberships    OrganizationMemberships
	OrganizationTags           OrganizationTags
	OrganizationTokens         OrganizationTokens
	Plans                      Plans
	PlanExports                PlanExports
	Policies                   Policies
	PolicyChecks               PolicyChecks
	PolicyEvaluations          PolicyEvaluations
	PolicySetOutcomes          PolicySetOutcomes
	PolicySetParameters        PolicySetParameters
	PolicySetVersions          PolicySetVersions
	PolicySets                 PolicySets
	RegistryModules            RegistryModules
	RegistryNoCodeModules      RegistryNoCodeModules
	RegistryProviders          RegistryProviders
	RegistryProviderPlatforms  RegistryProviderPlatforms
	RegistryProviderVersions   RegistryProviderVersions
	Runs                       Runs
	RunEvents                  RunEvents
	RunTasks                   RunTasks
	RunTasksIntegration        RunTasksIntegration
	RunTriggers                RunTriggers
	SSHKeys                    SSHKeys
	Stacks                     Stacks
	StackConfigurations        StackConfigurations
	StackDeployments           StackDeployments
	StackPlans                 StackPlans
	StackPlanOperations        StackPlanOperations
	StackSources               StackSources
	StateVersionOutputs        StateVersionOutputs
	StateVersions              StateVersions
	TaskResults                TaskResults
	TaskStages                 TaskStages
	Teams                      Teams
	TeamAccess                 TeamAccesses
	TeamMembers                TeamMembers
	TeamProjectAccess          TeamProjectAccesses
	TeamTokens                 TeamTokens
	TestRuns                   TestRuns
	TestVariables              TestVariables
	Users                      Users
	UserTokens                 UserTokens
	Variables                  Variables
	VariableSets               VariableSets
	VariableSetVariables       VariableSetVariables
	Workspaces                 Workspaces
	WorkspaceResources         WorkspaceResources
	WorkspaceRunTasks          WorkspaceRunTasks
	Projects                   Projects

	Meta Meta
}

// Admin is the the Terraform Enterprise Admin API. It provides access to site
// wide admin settings. These are only available for Terraform Enterprise and
// do not function against HCP Terraform
type Admin struct {
	Organizations     AdminOrganizations
	Workspaces        AdminWorkspaces
	Runs              AdminRuns
	TerraformVersions AdminTerraformVersions
	OPAVersions       AdminOPAVersions
	SentinelVersions  AdminSentinelVersions
	Users             AdminUsers
	Settings          *AdminSettings
}

// Meta contains any HCP Terraform APIs which provide data about the API itself.
type Meta struct {
	IPRanges IPRanges
}

// doForeignPUTRequest performs a PUT request using the specific data body. The Content-Type
// header is set to application/octet-stream but no Authentication header is sent. No response
// body is decoded.
func (c *Client) doForeignPUTRequest(ctx context.Context, foreignURL string, data io.Reader) error {
	u, err := url.Parse(foreignURL)
	if err != nil {
		return fmt.Errorf("specified URL was not valid: %w", err)
	}

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json, */*")
	reqHeaders.Set("Content-Type", "application/octet-stream")

	req, err := retryablehttp.NewRequest("PUT", u.String(), data)
	if err != nil {
		return err
	}

	// Set the default headers.
	for k, v := range c.headers {
		req.Header[k] = v
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	request := &ClientRequest{
		retryableRequest: req,
		http:             c.http,
		Header:           req.Header,
	}

	return request.DoJSON(ctx, nil)
}

// NewRequest performs some basic API request preparation based on the method
// specified. For GET requests, the reqBody is encoded as query parameters.
// For DELETE, PATCH, and POST requests, the request body is serialized as JSONAPI.
// For PUT requests, the request body is sent as a stream of bytes.
func (c *Client) NewRequest(method, path string, reqBody any) (*ClientRequest, error) {
	return c.NewRequestWithAdditionalQueryParams(method, path, reqBody, nil)
}

// NewRequestWithAdditionalQueryParams performs some basic API request
// preparation based on the method specified. For GET requests, the reqBody is
// encoded as query parameters. For DELETE, PATCH, and POST requests, the
// request body is serialized as JSONAPI. For PUT requests, the request body is
// sent as a stream of bytes. Additional query parameters can be added to the
// request as a string map. Note that if a key exists in both the reqBody and
// additionalQueryParams, the value in additionalQueryParams will be used.
func (c *Client) NewRequestWithAdditionalQueryParams(method, path string, reqBody any, additionalQueryParams map[string][]string) (*ClientRequest, error) {
	var u *url.URL
	var err error
	if strings.Contains(path, "/api/registry/") {
		u, err = c.registryBaseURL.Parse(path)
		if err != nil {
			return nil, err
		}
	} else {
		u, err = c.baseURL.Parse(path)
		if err != nil {
			return nil, err
		}
	}

	// Will contain combined query values from path parsing and
	// additionalQueryParams parameter
	q := make(url.Values)

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Authorization", "Bearer "+c.token)

	var body any
	switch method {
	case "GET":
		reqHeaders.Set("Accept", ContentTypeJSONAPI)

		// Encode the reqBody as query parameters
		if reqBody != nil {
			q, err = query.Values(reqBody)
			if err != nil {
				return nil, err
			}
		}
	case "DELETE", "PATCH", "POST":
		reqHeaders.Set("Accept", ContentTypeJSONAPI)
		reqHeaders.Set("Content-Type", ContentTypeJSONAPI)

		if reqBody != nil {
			if body, err = serializeRequestBody(reqBody); err != nil {
				return nil, err
			}
		}
	case "PUT":
		reqHeaders.Set("Accept", "application/json")
		reqHeaders.Set("Content-Type", "application/octet-stream")
		body = reqBody
	}

	for k, v := range u.Query() {
		q[k] = v
	}
	for k, v := range additionalQueryParams {
		q[k] = v
	}

	u.RawQuery = encodeQueryParams(q)

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// Set the default headers.
	for k, v := range c.headers {
		req.Header[k] = v
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return &ClientRequest{
		retryableRequest: req,
		http:             c.http,
		limiter:          c.limiter,
		Header:           req.Header,
	}, nil
}

// NewClient creates a new Terraform Enterprise API client.
func NewClient(cfg *Config) (*Client, error) {
	config := DefaultConfig()

	// Layer in the provided config for any non-blank values.
	if cfg != nil { // nolint
		if cfg.Address != "" {
			config.Address = cfg.Address
		}
		if cfg.BasePath != "" {
			config.BasePath = cfg.BasePath
		}
		if cfg.RegistryBasePath != "" {
			config.RegistryBasePath = cfg.RegistryBasePath
		}
		if cfg.Token != "" {
			config.Token = cfg.Token
		}
		for k, v := range cfg.Headers {
			config.Headers[k] = v
		}
		if cfg.HTTPClient != nil {
			config.HTTPClient = cfg.HTTPClient
		}
		if cfg.RetryLogHook != nil {
			config.RetryLogHook = cfg.RetryLogHook
		}
		config.RetryServerErrors = cfg.RetryServerErrors
	}

	// Parse the address to make sure its a valid URL.
	baseURL, err := url.Parse(config.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	baseURL.Path = config.BasePath
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	registryURL, err := url.Parse(config.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	registryURL.Path = config.RegistryBasePath
	if !strings.HasSuffix(registryURL.Path, "/") {
		registryURL.Path += "/"
	}

	// This value must be provided by the user.
	if config.Token == "" {
		return nil, fmt.Errorf("missing API token")
	}

	// Create the client.
	client := &Client{
		baseURL:           baseURL,
		registryBaseURL:   registryURL,
		token:             config.Token,
		headers:           config.Headers,
		retryLogHook:      config.RetryLogHook,
		retryServerErrors: config.RetryServerErrors,
	}

	client.http = &retryablehttp.Client{
		Backoff:      client.retryHTTPBackoff,
		CheckRetry:   client.retryHTTPCheck,
		ErrorHandler: retryablehttp.PassthroughErrorHandler,
		HTTPClient:   config.HTTPClient,
		RetryWaitMin: 100 * time.Millisecond,
		RetryWaitMax: 400 * time.Millisecond,
		RetryMax:     30,
	}

	meta, err := client.getRawAPIMetadata()
	if err != nil {
		return nil, err
	}

	// Configure the rate limiter.
	client.configureLimiter(meta.RateLimit)

	// Save the API version so we can return it from the RemoteAPIVersion
	// method later.
	client.remoteAPIVersion = meta.APIVersion

	// Save the TFE version
	client.remoteTFEVersion = meta.TFEVersion

	// Save the app name
	client.appName = meta.AppName

	// Create Admin
	client.Admin = Admin{
		Organizations:     &adminOrganizations{client: client},
		Workspaces:        &adminWorkspaces{client: client},
		Runs:              &adminRuns{client: client},
		Settings:          newAdminSettings(client),
		TerraformVersions: &adminTerraformVersions{client: client},
		OPAVersions:       &adminOPAVersions{client: client},
		SentinelVersions:  &adminSentinelVersions{client: client},
		Users:             &adminUsers{client: client},
	}

	// Create the services.
	client.AgentPools = &agentPools{client: client}
	client.Agents = &agents{client: client}
	client.AgentTokens = &agentTokens{client: client}
	client.Applies = &applies{client: client}
	client.AuditTrails = &auditTrails{client: client}
	client.Comments = &comments{client: client}
	client.ConfigurationVersions = &configurationVersions{client: client}
	client.GHAInstallations = &gHAInstallations{client: client}
	client.CostEstimates = &costEstimates{client: client}
	client.GPGKeys = &gpgKeys{client: client}
	client.RegistryNoCodeModules = &registryNoCodeModules{client: client}
	client.NotificationConfigurations = &notificationConfigurations{client: client}
	client.OAuthClients = &oAuthClients{client: client}
	client.OAuthTokens = &oAuthTokens{client: client}
	client.OrganizationMemberships = &organizationMemberships{client: client}
	client.Organizations = &organizations{client: client}
	client.OrganizationTags = &organizationTags{client: client}
	client.OrganizationTokens = &organizationTokens{client: client}
	client.PlanExports = &planExports{client: client}
	client.Plans = &plans{client: client}
	client.Policies = &policies{client: client}
	client.PolicyChecks = &policyChecks{client: client}
	client.PolicyEvaluations = &policyEvaluation{client: client}
	client.PolicySetOutcomes = &policySetOutcome{client: client}
	client.PolicySetParameters = &policySetParameters{client: client}
	client.PolicySets = &policySets{client: client}
	client.PolicySetVersions = &policySetVersions{client: client}
	client.Projects = &projects{client: client}
	client.RegistryModules = &registryModules{client: client}
	client.RegistryProviderPlatforms = &registryProviderPlatforms{client: client}
	client.RegistryProviders = &registryProviders{client: client}
	client.RegistryProviderVersions = &registryProviderVersions{client: client}
	client.Runs = &runs{client: client}
	client.RunEvents = &runEvents{client: client}
	client.RunTasks = &runTasks{client: client}
	client.RunTasksIntegration = &runTaskIntegration{client: client}
	client.RunTriggers = &runTriggers{client: client}
	client.SSHKeys = &sshKeys{client: client}
	client.Stacks = &stacks{client: client}
	client.StackConfigurations = &stackConfigurations{client: client}
	client.StackDeployments = &stackDeployments{client: client}
	client.StackPlans = &stackPlans{client: client}
	client.StackPlanOperations = &stackPlanOperations{client: client}
	client.StackSources = &stackSources{client: client}
	client.StateVersionOutputs = &stateVersionOutputs{client: client}
	client.StateVersions = &stateVersions{client: client}
	client.TaskResults = &taskResults{client: client}
	client.TaskStages = &taskStages{client: client}
	client.TeamAccess = &teamAccesses{client: client}
	client.TeamMembers = &teamMembers{client: client}
	client.TeamProjectAccess = &teamProjectAccesses{client: client}
	client.Teams = &teams{client: client}
	client.TeamTokens = &teamTokens{client: client}
	client.TestRuns = &testRuns{client: client}
	client.TestVariables = &testVariables{client: client}
	client.Users = &users{client: client}
	client.UserTokens = &userTokens{client: client}
	client.Variables = &variables{client: client}
	client.VariableSets = &variableSets{client: client}
	client.VariableSetVariables = &variableSetVariables{client: client}
	client.WorkspaceRunTasks = &workspaceRunTasks{client: client}
	client.Workspaces = &workspaces{client: client}
	client.WorkspaceResources = &workspaceResources{client: client}

	client.Meta = Meta{
		IPRanges: &ipRanges{client: client},
	}

	return client, nil
}

// AppName returns the name of the instance.
func (c Client) AppName() string {
	return c.appName
}

// IsCloud returns true if the client is configured against a HCP Terraform
// instance.
//
// Whether an instance is HCP Terraform or Terraform Enterprise is derived from the TFP-AppName header.
func (c Client) IsCloud() bool {
	return c.appName == "HCP Terraform"
}

// IsEnterprise returns true if the client is configured against a Terraform
// Enterprise instance.
//
// Whether an instance is HCP Terraform or TFE is derived from the TFP-AppName header. Note:
// not all TFE releases include this header in API responses.
func (c Client) IsEnterprise() bool {
	return !c.IsCloud()
}

// RemoteAPIVersion returns the server's declared API version string.
//
// A HCP Terraform or Enterprise API server returns its API version in an
// HTTP header field in all responses. The NewClient function saves the
// version number returned in its initial setup request and RemoteAPIVersion
// returns that cached value.
//
// The API protocol calls for this string to be a dotted-decimal version number
// like 2.3.0, where the first number indicates the API major version while the
// second indicates a minor version which may have introduced some
// backward-compatible additional features compared to its predecessor.
//
// Explicit API versioning was added to the HCP Terraform and Enterprise
// APIs as a later addition, so older servers will not return version
// information. In that case, this function returns an empty string as the
// version.
func (c Client) RemoteAPIVersion() string {
	return c.remoteAPIVersion
}

// BaseURL returns the base URL as configured in the client
func (c Client) BaseURL() url.URL {
	return *c.baseURL
}

// BaseRegistryURL returns the registry base URL as configured in the client
func (c Client) BaseRegistryURL() url.URL {
	return *c.registryBaseURL
}

// SetFakeRemoteAPIVersion allows setting a given string as the client's remoteAPIVersion,
// overriding the value pulled from the API header during client initialization.
//
// This is intended for use in tests, when you may want to configure your TFE client to
// return something different than the actual API version in order to test error handling.
func (c *Client) SetFakeRemoteAPIVersion(fakeAPIVersion string) {
	c.remoteAPIVersion = fakeAPIVersion
}

// RemoteTFEVersion returns the server's declared TFE version string.
//
// A Terraform Enterprise API server includes its current version in an
// HTTP header field in all responses. This value is saved by the client
// during the initial setup request and RemoteTFEVersion returns that cached
// value. This function returns an empty string for any Terraform Enterprise version
// earlier than v202208-3 and for HCP Terraform.
func (c Client) RemoteTFEVersion() string {
	return c.remoteTFEVersion
}

// RetryServerErrors configures the retry HTTP check to also retry
// unexpected errors or requests that failed with a server error.
func (c *Client) RetryServerErrors(retry bool) {
	c.retryServerErrors = retry
}

// retryHTTPCheck provides a callback for Client.CheckRetry which
// will retry both rate limit (429) and server (>= 500) errors.
func (c *Client) retryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return c.retryServerErrors, err
	}
	if resp.StatusCode == 429 || (c.retryServerErrors && resp.StatusCode >= 500) {
		return true, nil
	}
	return false, nil
}

// retryHTTPBackoff provides a generic callback for Client.Backoff which
// will pass through all calls based on the status code of the response.
func (c *Client) retryHTTPBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if c.retryLogHook != nil {
		c.retryLogHook(attemptNum, resp)
	}

	// Use the rate limit backoff function when we are rate limited.
	if resp != nil && resp.StatusCode == 429 {
		return rateLimitBackoff(min, max, resp)
	}

	// Set custom duration's when we experience a service interruption.
	min = 700 * time.Millisecond
	max = 900 * time.Millisecond

	return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
}

// rateLimitBackoff provides a callback for Client.Backoff which will use the
// X-RateLimit_Reset header to determine the time to wait. We add some jitter
// to prevent a thundering herd.
//
// min and max are mainly used for bounding the jitter that will be added to
// the reset time retrieved from the headers. But if the final wait time is
// less than min, min will be used instead.
func rateLimitBackoff(min, max time.Duration, resp *http.Response) time.Duration {
	// rnd is used to generate pseudo-random numbers.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// First create some jitter bounded by the min and max durations.
	jitter := time.Duration(rnd.Float64() * float64(max-min))

	if resp != nil && resp.Header.Get(_headerRateReset) != "" {
		v := resp.Header.Get(_headerRateReset)
		reset, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Fatal(err)
		}
		// Only update min if the given time to wait is longer
		if reset > 0 && time.Duration(reset*1e9) > min {
			min = time.Duration(reset * 1e9)
		}
	}

	return min + jitter
}

type rawAPIMetadata struct {
	// APIVersion is the raw API version string reported by the server in the
	// TFP-API-Version response header, or an empty string if that header
	// field was not included in the response.
	APIVersion string

	// TFEVersion is the raw TFE version string reported by the server in the
	// X-TFE-Version response header, or an empty string if that header
	// field was not included in the response.
	TFEVersion string

	// RateLimit is the raw API version string reported by the server in the
	// X-RateLimit-Limit response header, or an empty string if that header
	// field was not included in the response.
	RateLimit string

	// AppName is either 'HCP Terraform' or 'Terraform Enterprise'
	AppName string
}

func (c *Client) getRawAPIMetadata() (rawAPIMetadata, error) {
	var meta rawAPIMetadata

	// Create a new request.
	u, err := c.baseURL.Parse(PingEndpoint)
	if err != nil {
		return meta, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return meta, err
	}

	// Attach the default headers.
	for k, v := range c.headers {
		req.Header[k] = v
	}
	req.Header.Set("Accept", ContentTypeJSONAPI)
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Make a single request to retrieve the rate limit headers.
	resp, err := c.http.HTTPClient.Do(req)
	if err != nil {
		return meta, err
	}
	resp.Body.Close()

	meta.APIVersion = resp.Header.Get(_headerAPIVersion)
	meta.RateLimit = resp.Header.Get(_headerRateLimit)
	meta.TFEVersion = resp.Header.Get(_headerTFEVersion)
	meta.AppName = resp.Header.Get(_headerAppName)

	return meta, nil
}

// configureLimiter configures the rate limiter.
func (c *Client) configureLimiter(rawLimit string) {
	// Set default values for when rate limiting is disabled.
	limit := rate.Inf
	burst := 0

	if v := rawLimit; v != "" {
		if rateLimit, err := strconv.ParseFloat(v, 64); rateLimit > 0 {
			if err != nil {
				log.Fatal(err)
			}
			// Configure the limit and burst using a split of 2/3 for the limit and
			// 1/3 for the burst. This enables clients to burst 1/3 of the allowed
			// calls before the limiter kicks in. The remaining calls will then be
			// spread out evenly using intervals of time.Second / limit which should
			// prevent hitting the rate limit.
			limit = rate.Limit(rateLimit * 0.66)
			burst = int(rateLimit * 0.33)
		}
	}

	// Create a new limiter using the calculated values.
	c.limiter = rate.NewLimiter(limit, burst)
}

// encodeQueryParams encodes the values into "URL encoded" form
// ("bar=baz&foo=quux") sorted by key. This version behaves as url.Values
// Encode, except that it encodes certain keys as comma-separated values instead
// of using multiple keys.
func encodeQueryParams(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		if len(vs) > 1 && validSliceKey(k) {
			val := strings.Join(vs, ",")
			vs = vs[:0]
			vs = append(vs, val)
		}
		keyEscaped := url.QueryEscape(k)

		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// decodeQueryParams types an object and converts the struct fields into
// Query Parameters, which can be used with NewRequestWithAdditionalQueryParams
// Note that a field without a `url` annotation will be converted into a query
// parameter. Use url:"-" to ignore struct fields.
func decodeQueryParams(v any) (url.Values, error) {
	if v == nil {
		return make(url.Values, 0), nil
	}
	return query.Values(v)
}

// serializeRequestBody serializes the given ptr or ptr slice into a JSON
// request. It automatically uses jsonapi or json serialization, depending
// on the body type's tags.
func serializeRequestBody(v interface{}) (interface{}, error) {
	// The body can be a slice of pointers or a pointer. In either
	// case we want to choose the serialization type based on the
	// individual record type. To determine that type, we need
	// to either follow the pointer or examine the slice element type.
	// There are other theoretical possibilities (e. g. maps,
	// non-pointers) but they wouldn't work anyway because the
	// json-api library doesn't support serializing other things.
	var modelType reflect.Type
	bodyType := reflect.TypeOf(v)
	switch bodyType.Kind() {
	case reflect.Slice:
		sliceElem := bodyType.Elem()
		if sliceElem.Kind() != reflect.Ptr {
			return nil, ErrInvalidRequestBody
		}
		modelType = sliceElem.Elem()
	case reflect.Ptr:
		modelType = reflect.ValueOf(v).Elem().Type()
	default:
		return nil, ErrInvalidRequestBody
	}

	// Infer whether the request uses jsonapi or regular json
	// serialization based on how the fields are tagged.
	jsonAPIFields := 0
	jsonFields := 0
	for i := 0; i < modelType.NumField(); i++ {
		structField := modelType.Field(i)
		if structField.Tag.Get("jsonapi") != "" {
			jsonAPIFields++
		}
		if structField.Tag.Get("json") != "" {
			jsonFields++
		}
	}
	if jsonAPIFields > 0 && jsonFields > 0 {
		// Defining a struct with both json and jsonapi tags doesn't
		// make sense, because a struct can only be serialized
		// as one or another. If this does happen, it's a bug
		// in the library that should be fixed at development time
		return nil, ErrInvalidStructFormat
	}

	if jsonFields > 0 {
		return json.Marshal(v)
	}
	buf := bytes.NewBuffer(nil)
	if err := jsonapi.MarshalPayloadWithoutIncluded(buf, v); err != nil {
		return nil, err
	}
	return buf, nil
}

func unmarshalResponse(responseBody io.Reader, model interface{}) error {
	// Get the value of model so we can test if it's a struct.
	dst := reflect.Indirect(reflect.ValueOf(model))

	// Return an error if model is not a struct or an io.Writer.
	if dst.Kind() != reflect.Struct {
		return fmt.Errorf("%v must be a struct or an io.Writer", dst)
	}

	// Try to get the Items and Pagination struct fields.
	items := dst.FieldByName("Items")
	pagination := dst.FieldByName("Pagination")

	// Unmarshal a single value if model does not contain the
	// Items and Pagination struct fields.
	if !items.IsValid() || !pagination.IsValid() {
		return jsonapi.UnmarshalPayload(responseBody, model)
	}

	// Return an error if model.Items is not a slice.
	if items.Type().Kind() != reflect.Slice {
		return ErrItemsMustBeSlice
	}

	// Create a temporary buffer and copy all the read data into it.
	body := bytes.NewBuffer(nil)
	reader := io.TeeReader(responseBody, body)

	// Unmarshal as a list of values as model.Items is a slice.
	raw, err := jsonapi.UnmarshalManyPayload(reader, items.Type().Elem())
	if err != nil {
		return err
	}

	// Make a new slice to hold the results.
	sliceType := reflect.SliceOf(items.Type().Elem())
	result := reflect.MakeSlice(sliceType, 0, len(raw))

	// Add all of the results to the new slice.
	for _, v := range raw {
		result = reflect.Append(result, reflect.ValueOf(v))
	}

	// Pointer-swap the result.
	items.Set(result)

	// As we are getting a list of values, we need to decode
	// the pagination details out of the response body.
	p, err := parsePagination(body)
	if err != nil {
		return err
	}

	// Pointer-swap the decoded pagination details.
	pagination.Set(reflect.ValueOf(p))

	return nil
}

// ListOptions is used to specify pagination options when making API requests.
// Pagination allows breaking up large result sets into chunks, or "pages".
type ListOptions struct {
	// The page number to request. The results vary based on the PageSize.
	PageNumber int `url:"page[number],omitempty"`

	// The number of elements returned in a single page.
	PageSize int `url:"page[size],omitempty"`
}

// Pagination is used to return the pagination details of an API request.
type Pagination struct {
	CurrentPage  int `json:"current-page"`
	PreviousPage int `json:"prev-page"`
	NextPage     int `json:"next-page"`
	TotalPages   int `json:"total-pages"`
	TotalCount   int `json:"total-count"`
}

func parsePagination(body io.Reader) (*Pagination, error) {
	var raw struct {
		Meta struct {
			Pagination Pagination `jsonapi:"pagination"`
		} `jsonapi:"meta"`
	}

	// JSON decode the raw response.
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return &Pagination{}, err
	}

	return &raw.Meta.Pagination, nil
}

// checkResponseCode refines typical API errors into more specific errors
// if possible. It returns nil if the response code < 400
func checkResponseCode(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 399 {
		return nil
	}

	var errs []string
	var err error

	switch r.StatusCode {
	case 400:
		errs, err = decodeErrorPayload(r)
		if err != nil {
			return err
		}

		if errorPayloadContains(errs, "Invalid include parameter") {
			return ErrInvalidIncludeValue
		}
		return errors.New(strings.Join(errs, "\n"))
	case 401:
		return ErrUnauthorized
	case 404:
		return ErrResourceNotFound
	case 409:
		switch {
		case strings.HasSuffix(r.Request.URL.Path, "actions/lock"):
			return ErrWorkspaceLocked
		case strings.HasSuffix(r.Request.URL.Path, "actions/unlock"):
			errs, err = decodeErrorPayload(r)
			if err != nil {
				return err
			}

			if errorPayloadContains(errs, "is locked by Run") {
				return ErrWorkspaceLockedByRun
			}

			if errorPayloadContains(errs, "is locked by Team") {
				return ErrWorkspaceLockedByTeam
			}

			if errorPayloadContains(errs, "is locked by User") {
				return ErrWorkspaceLockedByUser
			}

			return ErrWorkspaceNotLocked
		case strings.HasSuffix(r.Request.URL.Path, "actions/force-unlock"):
			return ErrWorkspaceNotLocked
		case strings.HasSuffix(r.Request.URL.Path, "actions/safe-delete"):
			errs, err = decodeErrorPayload(r)
			if err != nil {
				return err
			}
			if errorPayloadContains(errs, "locked") {
				return ErrWorkspaceLockedCannotDelete
			}
			if errorPayloadContains(errs, "being processed") {
				return ErrWorkspaceStillProcessing
			}

			return ErrWorkspaceNotSafeToDelete
		}
	}

	errs, err = decodeErrorPayload(r)
	if err != nil {
		return err
	}

	return errors.New(strings.Join(errs, "\n"))
}

func decodeErrorPayload(r *http.Response) ([]string, error) {
	// Decode the error payload.
	var errs []string
	errPayload := &jsonapi.ErrorsPayload{}
	err := json.NewDecoder(r.Body).Decode(errPayload)
	if err != nil || len(errPayload.Errors) == 0 {
		return errs, errors.New(r.Status)
	}

	// Parse and format the errors.
	for _, e := range errPayload.Errors {
		if e.Detail == "" {
			errs = append(errs, e.Title)
		} else {
			errs = append(errs, fmt.Sprintf("%s\n\n%s", e.Title, e.Detail))
		}
	}

	return errs, nil
}

func errorPayloadContains(payloadErrors []string, match string) bool {
	for _, e := range payloadErrors {
		if strings.Contains(e, match) {
			return true
		}
	}
	return false
}

func packContents(path string) (*bytes.Buffer, error) {
	body := bytes.NewBuffer(nil)

	file, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return body, fmt.Errorf(`failed to find files under the path "%v": %w`, path, err)
		}
		return body, fmt.Errorf(`unable to upload files from the path "%v": %w`, path, err)
	}

	if !file.Mode().IsDir() {
		return body, ErrMissingDirectory
	}

	_, errSlug := slug.Pack(path, body, true)
	if errSlug != nil {
		return body, errSlug
	}

	return body, nil
}

func validSliceKey(key string) bool {
	return key == _includeQueryParam || strings.Contains(key, "filter[")
}
