package mongodbatlas

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
	"strconv"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://cloud.mongodb.com/api/atlas/v1.0/"
	userAgent      = "go-mongodbatlas" + libraryVersion
	mediaType      = "application/json"
)

// Client manages communication with MongoDBAtlas v1.0 API
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	//Services used for communicating with the API
	CustomDBRoles                    CustomDBRolesService
	DatabaseUsers                    DatabaseUsersService
	ProjectIPWhitelist               ProjectIPWhitelistService
	Projects                         ProjectsService
	Clusters                         ClustersService
	CloudProviderSnapshots           CloudProviderSnapshotsService
	APIKeys                          APIKeysService
	ProjectAPIKeys                   ProjectAPIKeysService
	CloudProviderSnapshotRestoreJobs CloudProviderSnapshotRestoreJobsService
	Peers                            PeersService
	Containers                       ContainersService
	EncryptionsAtRest                EncryptionsAtRestService
	WhitelistAPIKeys                 WhitelistAPIKeysService
	PrivateIPMode                    PrivateIpModeService
	MaintenanceWindows               MaintenanceWindowsService
	Teams                            TeamsService
	AtlasUsers                       AtlasUsersService
	GlobalClusters                   GlobalClustersService
	Auditing                         AuditingsService

	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

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

//ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
	//The error code, which is simply the HTTP status code.
	ErrorCode int `json:"Error"`

	//A short description of the error, which is simply the HTTP status phrase.
	Reason string `json:"reason"`

	//A more detailed description of the error.
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

//IsLastPage returns true if the current page is the last page
func (resp *Response) IsLastPage() bool {
	return resp.getLinkByRef("next") == nil
}

//CurrentPage gets the current page for list pagination request.
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

	c.APIKeys = &APIKeysServiceOp{client: c}
	c.CloudProviderSnapshots = &CloudProviderSnapshotsServiceOp{client: c}
	c.CloudProviderSnapshotRestoreJobs = &CloudProviderSnapshotRestoreJobsServiceOp{client: c}
	c.Clusters = &ClustersServiceOp{client: c}
	c.Containers = &ContainersServiceOp{client: c}
	c.CustomDBRoles = &CustomDBRolesServiceOp{client: c}
	c.DatabaseUsers = &DatabaseUsersServiceOp{client: c}
	c.EncryptionsAtRest = &EncryptionsAtRestServiceOp{client: c}
	c.Projects = &ProjectsServiceOp{client: c}
	c.ProjectAPIKeys = &ProjectAPIKeysOp{client: c}
	c.Peers = &PeersServiceOp{client: c}
	c.ProjectIPWhitelist = &ProjectIPWhitelistServiceOp{client: c}
	c.WhitelistAPIKeys = &WhitelistAPIKeysServiceOp{client: c}
	c.PrivateIPMode = &PrivateIpModeServiceOp{client: c}
	c.MaintenanceWindows = &MaintenanceWindowsServiceOp{client: c}
	c.Teams = &TeamsServiceOp{client: c}
	c.AtlasUsers = &AtlasUsersServiceOp{client: c}
	c.GlobalClusters = &GlobalClustersServiceOp{client: c}
	c.Auditing = &AuditingsServiceOp{client: c}

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
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
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

	response := newResponse(resp)

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
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
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

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
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
