package mongodbatlas // import "go.mongodb.org/atlas/mongodbatlas"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://cloud.mongodb.com/api/atlas/v1.0/"
	jsonMediaType  = "application/json"
	gzipMediaType  = "application/gzip"
	libraryName    = "go-mongodbatlas"
	// Version the version of the current API client
	Version = "0.5.0" // Should be set to the next version planned to be released
)

var (
	userAgent = fmt.Sprintf("%s/%s (%s;%s)", libraryName, Version, runtime.GOOS, runtime.GOARCH)
)

// Doer basic interface of a client to be able to do a request
type Doer interface {
	Do(context.Context, *http.Request, interface{}) (*Response, error)
}

// Completer interface for clients with callback
type Completer interface {
	OnRequestCompleted(RequestCompletionCallback)
}

// RequestDoer minimum interface for any service of the client
type RequestDoer interface {
	Doer
	Completer
	NewRequest(context.Context, string, string, interface{}) (*http.Request, error)
}

// GZipRequestDoer minimum interface for any service of the client that should handle gzip downloads
type GZipRequestDoer interface {
	Doer
	Completer
	NewGZipRequest(context.Context, string, string) (*http.Request, error)
}

// Client manages communication with MongoDBAtlas v1.0 API
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	// Services used for communicating with the API
	CustomDBRoles                       CustomDBRolesService
	DatabaseUsers                       DatabaseUsersService
	ProjectIPWhitelist                  ProjectIPWhitelistService
	ProjectIPAccessList                 ProjectIPAccessListService
	Organizations                       OrganizationsService
	Projects                            ProjectsService
	Clusters                            ClustersService
	CloudProviderSnapshots              CloudProviderSnapshotsService
	APIKeys                             APIKeysService
	ProjectAPIKeys                      ProjectAPIKeysService
	CloudProviderSnapshotRestoreJobs    CloudProviderSnapshotRestoreJobsService
	Peers                               PeersService
	Containers                          ContainersService
	EncryptionsAtRest                   EncryptionsAtRestService
	WhitelistAPIKeys                    WhitelistAPIKeysService
	PrivateIPMode                       PrivateIPModeService
	MaintenanceWindows                  MaintenanceWindowsService
	Teams                               TeamsService
	AtlasUsers                          AtlasUsersService
	GlobalClusters                      GlobalClustersService
	Auditing                            AuditingsService
	AlertConfigurations                 AlertConfigurationsService
	PrivateEndpoints                    PrivateEndpointsService
	X509AuthDBUsers                     X509AuthDBUsersService
	ContinuousSnapshots                 ContinuousSnapshotsService
	ContinuousRestoreJobs               ContinuousRestoreJobsService
	Checkpoints                         CheckpointsService
	Alerts                              AlertsService
	CloudProviderSnapshotBackupPolicies CloudProviderSnapshotBackupPoliciesService
	Events                              EventsService
	Processes                           ProcessesService
	ProcessMeasurements                 ProcessMeasurementsService
	ProcessDisks                        ProcessDisksService
	ProcessDiskMeasurements             ProcessDiskMeasurementsService
	ProcessDatabases                    ProcessDatabasesService
	ProcessDatabaseMeasurements         ProcessDatabaseMeasurementsService
	Indexes                             IndexesService
	Logs                                LogsService
	DataLakes                           DataLakeService
	OnlineArchives                      OnlineArchiveService
	Search                              SearchService
	CustomAWSDNS                        AWSCustomDNSService
	Integrations                        IntegrationsService
	LDAPConfigurations                  LDAPConfigurationsService
	PerformanceAdvisor                  PerformanceAdvisorService

	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

type service struct {
	Client RequestDoer
}

// Response is a MongoDBAtlas response. This wraps the standard http.Response returned from MongoDBAtlas API.
type Response struct {
	*http.Response

	// Links that were returned with the response.
	Links []*Link `json:"links"`
}

// ListOptions specifies the optional parameters to List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	PageNum int `url:"pageNum,omitempty"`

	// For paginated result sets, the number of results to include per page.
	ItemsPerPage int `url:"itemsPerPage,omitempty"`
}

// ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
	// The error code, which is simply the HTTP status code.
	ErrorCode int `json:"Error"`

	// A short description of the error, which is simply the HTTP status phrase.
	Reason string `json:"reason"`

	// A more detailed description of the error.
	Detail string `json:"detail,omitempty"`
}

func (resp *Response) getCurrentPageLink() (*Link, error) {
	if link := resp.getLinkByRef("self"); link != nil {
		return link, nil
	}
	return nil, errors.New("no self link found")
}

func (resp *Response) getLinkByRef(ref string) *Link {
	for i := range resp.Links {
		if resp.Links[i].Rel == ref {
			return resp.Links[i]
		}
	}
	return nil
}

// IsLastPage returns true if the current page is the last page
func (resp *Response) IsLastPage() bool {
	return resp.getLinkByRef("next") == nil
}

// CurrentPage gets the current page for list pagination request.
func (resp *Response) CurrentPage() (int, error) {
	link, err := resp.getCurrentPageLink()
	if err != nil {
		return 0, err
	}

	pageNumStr, err := link.getHrefQueryParam("pageNum")
	if err != nil {
		return 0, err
	}

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		return 0, fmt.Errorf("error getting current page: %s", err)
	}

	return pageNum, nil
}

// NewClient returns a new MongoDBAtlas API Client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.APIKeys = &APIKeysServiceOp{Client: c}
	c.CloudProviderSnapshots = &CloudProviderSnapshotsServiceOp{Client: c}
	c.ContinuousSnapshots = &ContinuousSnapshotsServiceOp{Client: c}
	c.CloudProviderSnapshotRestoreJobs = &CloudProviderSnapshotRestoreJobsServiceOp{Client: c}
	c.Clusters = &ClustersServiceOp{Client: c}
	c.Containers = &ContainersServiceOp{Client: c}
	c.CustomDBRoles = &CustomDBRolesServiceOp{Client: c}
	c.DatabaseUsers = &DatabaseUsersServiceOp{Client: c}
	c.EncryptionsAtRest = &EncryptionsAtRestServiceOp{Client: c}
	c.Organizations = &OrganizationsServiceOp{Client: c}
	c.Projects = &ProjectsServiceOp{Client: c}
	c.ProjectAPIKeys = &ProjectAPIKeysOp{Client: c}
	c.Peers = &PeersServiceOp{Client: c}
	c.ProjectIPWhitelist = &ProjectIPWhitelistServiceOp{Client: c}
	c.ProjectIPAccessList = &ProjectIPAccessListServiceOp{Client: c}
	c.WhitelistAPIKeys = &WhitelistAPIKeysServiceOp{Client: c}
	c.PrivateIPMode = &PrivateIPModeServiceOp{Client: c}
	c.MaintenanceWindows = &MaintenanceWindowsServiceOp{Client: c}
	c.Teams = &TeamsServiceOp{Client: c}
	c.AtlasUsers = &AtlasUsersServiceOp{Client: c}
	c.GlobalClusters = &GlobalClustersServiceOp{Client: c}
	c.Auditing = &AuditingsServiceOp{Client: c}
	c.AlertConfigurations = &AlertConfigurationsServiceOp{Client: c}
	c.PrivateEndpoints = &PrivateEndpointsServiceOp{Client: c}
	c.X509AuthDBUsers = &X509AuthDBUsersServiceOp{Client: c}
	c.ContinuousRestoreJobs = &ContinuousRestoreJobsServiceOp{Client: c}
	c.Checkpoints = &CheckpointsServiceOp{Client: c}
	c.Alerts = &AlertsServiceOp{Client: c}
	c.CloudProviderSnapshotBackupPolicies = &CloudProviderSnapshotBackupPoliciesServiceOp{Client: c}
	c.Events = &EventsServiceOp{Client: c}
	c.Processes = &ProcessesServiceOp{Client: c}
	c.ProcessMeasurements = &ProcessMeasurementsServiceOp{Client: c}
	c.ProcessDisks = &ProcessDisksServiceOp{Client: c}
	c.ProcessDiskMeasurements = &ProcessDiskMeasurementsServiceOp{Client: c}
	c.ProcessDatabases = &ProcessDatabasesServiceOp{Client: c}
	c.ProcessDatabaseMeasurements = &ProcessDatabaseMeasurementsServiceOp{Client: c}
	c.Indexes = &IndexesServiceOp{Client: c}
	c.Logs = &LogsServiceOp{Client: c}
	c.DataLakes = &DataLakeServiceOp{Client: c}
	c.OnlineArchives = &OnlineArchiveServiceOp{Client: c}
	c.Search = &SearchServiceOp{Client: c}
	c.CustomAWSDNS = &AWSCustomDNSServiceOp{Client: c}
	c.Integrations = &IntegrationsServiceOp{Client: c}
	c.LDAPConfigurations = &LDAPConfigurationsServiceOp{Client: c}
	c.PerformanceAdvisor = &PerformanceAdvisorServiceOp{Client: c}

	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new MongoDBAtlas API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
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
		c.UserAgent = fmt.Sprintf("%s %s", ua, userAgent)
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("base URL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var buf io.Reader
	if body != nil {
		if buf, err = c.newEncodedBody(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", jsonMediaType)
	}
	req.Header.Add("Accept", jsonMediaType)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// newEncodedBody returns an ReadWriter object containing the body of the http request
func (c *Client) newEncodedBody(body interface{}) (io.Reader, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	return buf, err
}

// NewGZipRequest creates an API request that accepts gzip. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash.
func (c *Client) NewGZipRequest(ctx context.Context, method, urlStr string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("base URL must have a trailing slash, but %q does not", c.BaseURL)
	}
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", gzipMediaType)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := &Response{Response: resp}

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d (request %q) %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Reason, r.Detail)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			log.Printf("[DEBUG] unmarshal error response: %s", err)
			errorResponse.Reason = string(data)
		}
	}

	return errorResponse
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func setListOptions(s string, opt interface{}) (string, error) {
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
