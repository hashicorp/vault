package packngo

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	authTokenEnvVar = "PACKET_AUTH_TOKEN"
	libraryVersion  = "0.6.0"
	baseURL         = "https://api.equinix.com/metal/v1/"
	userAgent       = "packngo/" + libraryVersion
	mediaType       = "application/json"
	debugEnvVar     = "PACKNGO_DEBUG"

	headerRateLimit              = "X-RateLimit-Limit"
	headerRateRemaining          = "X-RateLimit-Remaining"
	headerRateReset              = "X-RateLimit-Reset"
	expectedAPIContentTypePrefix = "application/json"
)

var redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

// meta contains pagination information
type meta struct {
	Self           *Href `json:"self"`
	First          *Href `json:"first"`
	Last           *Href `json:"last"`
	Previous       *Href `json:"previous,omitempty"`
	Next           *Href `json:"next,omitempty"`
	Total          int   `json:"total"`
	CurrentPageNum int   `json:"current_page"`
	LastPageNum    int   `json:"last_page"`
}

// Response is the http response from api calls
type Response struct {
	*http.Response
	Rate
}

// Href is an API link
type Href struct {
	Href string `json:"href"`
}

func (r *Response) populateRate() {
	// parse the rate limit headers and populate Response.Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		r.Rate.RequestLimit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		r.Rate.RequestsRemaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			r.Rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
}

// ErrorResponse is the http response used on errors
type ErrorResponse struct {
	Response    *http.Response
	Errors      []string `json:"errors"`
	SingleError string   `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, strings.Join(r.Errors, ", "), r.SingleError)
}

// Client is the base API Client
type Client struct {
	client *http.Client
	debug  bool

	BaseURL *url.URL

	UserAgent     string
	ConsumerToken string
	APIKey        string

	RateLimit Rate

	// Equinix Metal Api Objects
	APIKeys                APIKeyService
	BGPConfig              BGPConfigService
	BGPSessions            BGPSessionService
	Batches                BatchService
	CapacityService        CapacityService
	Connections            ConnectionService
	DeviceIPs              DeviceIPService
	DevicePorts            DevicePortService
	Devices                DeviceService
	Emails                 EmailService
	Events                 EventService
	Facilities             FacilityService
	HardwareReservations   HardwareReservationService
	Notifications          NotificationService
	OperatingSystems       OSService
	Organizations          OrganizationService
	Plans                  PlanService
	ProjectIPs             ProjectIPService
	ProjectVirtualNetworks ProjectVirtualNetworkService
	Projects               ProjectService
	SSHKeys                SSHKeyService
	SpotMarket             SpotMarketService
	SpotMarketRequests     SpotMarketRequestService
	TwoFactorAuth          TwoFactorAuthService
	Users                  UserService
	VPN                    VPNService
	VirtualCircuits        VirtualCircuitService
	VolumeAttachments      VolumeAttachmentService
	Volumes                VolumeService
}

// requestDoer provides methods for making HTTP requests and receiving the
// response, errors, and a structured result
//
// This interface is used in *ServiceOp as a mockable alternative to a full
// Client object.
type requestDoer interface {
	NewRequest(method, path string, body interface{}) (*http.Request, error)
	Do(req *http.Request, v interface{}) (*Response, error)
	DoRequest(method, path string, body, v interface{}) (*Response, error)
	DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error)
}

// NewRequest inits a new http request with the proper headers
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	// relative path to append to the endpoint url, no leading slash please
	if path[0] == '/' {
		path = path[1:]
	}
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	// json encode the request body, if any
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Close = true

	req.Header.Add("X-Auth-Token", c.APIKey)
	req.Header.Add("X-Consumer-Token", c.ConsumerToken)

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do executes the http request
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := Response{Response: resp}
	response.populateRate()
	if c.debug {
		dumpResponse(response.Response)
	}
	c.RateLimit = response.Rate

	err = checkResponse(resp)
	// if the response is an error, return the ErrorResponse
	if err != nil {
		return &response, err
	}

	if v != nil {
		// if v implements the io.Writer interface, return the raw response
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return &response, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return &response, err
			}
		}
	}

	return &response, err
}

func dumpResponse(resp *http.Response) {
	o, _ := httputil.DumpResponse(resp, true)
	strResp := string(o)
	reg, _ := regexp.Compile(`"token":(.+?),`)
	reMatches := reg.FindStringSubmatch(strResp)
	if len(reMatches) == 2 {
		strResp = strings.Replace(strResp, reMatches[1], strings.Repeat("-", len(reMatches[1])), 1)
	}
	log.Printf("\n=======[RESPONSE]============\n%s\n\n", strResp)
}

func dumpRequest(req *http.Request) {
	r := req.Clone(context.TODO())
	r.Body, _ = req.GetBody()
	h := r.Header
	if len(h.Get("X-Auth-Token")) != 0 {
		h.Set("X-Auth-Token", "**REDACTED**")
	}
	defer r.Body.Close()

	o, _ := httputil.DumpRequestOut(r, false)
	bbs, _ := ioutil.ReadAll(r.Body)

	strReq := string(o)
	log.Printf("\n=======[REQUEST]=============\n%s%s\n", string(strReq), string(bbs))
}

// DoRequest is a convenience method, it calls NewRequest followed by Do
// v is the interface to unmarshal the response JSON into
func (c *Client) DoRequest(method, path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, path, body)
	if c.debug {
		dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// DoRequestWithHeader same as DoRequest
func (c *Client) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, path, body)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if c.debug {
		dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// NewClient initializes and returns a Client
func NewClient() (*Client, error) {
	apiToken := os.Getenv(authTokenEnvVar)
	if apiToken == "" {
		return nil, fmt.Errorf("you must export %s", authTokenEnvVar)
	}
	c := NewClientWithAuth("packngo lib", apiToken, nil)
	return c, nil

}

// NewClientWithAuth initializes and returns a Client, use this to get an API Client to operate on
// N.B.: Equinix Metal's API certificate requires Go 1.5+ to successfully parse. If you are using
// an older version of Go, pass in a custom http.Client with a custom TLS configuration
// that sets "InsecureSkipVerify" to "true"
func NewClientWithAuth(consumerToken string, apiKey string, httpClient *http.Client) *Client {
	client, _ := NewClientWithBaseURL(consumerToken, apiKey, httpClient, baseURL)
	return client
}

// RetryPolicy determines if the supplied http Response and error can be safely
// retried (for use with github.com/hashicorp/go-retryablehttp clients)
//
//    retryClient := retryablehttp.NewClient()
//    retryClient.CheckRetry = packngo.RetryPolicy
func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, nil
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, nil
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	//if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
	//	return true, nil
	//}

	return false, nil
}

// NewClientWithBaseURL returns a Client pointing to nonstandard API URL, e.g.
// for mocking the remote API
func NewClientWithBaseURL(consumerToken string, apiKey string, httpClient *http.Client, apiBaseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, BaseURL: u, UserAgent: userAgent, ConsumerToken: consumerToken, APIKey: apiKey}
	c.APIKeys = &APIKeyServiceOp{client: c}
	c.BGPConfig = &BGPConfigServiceOp{client: c}
	c.BGPSessions = &BGPSessionServiceOp{client: c}
	c.Batches = &BatchServiceOp{client: c}
	c.CapacityService = &CapacityServiceOp{client: c}
	c.Connections = &ConnectionServiceOp{client: c}
	c.DeviceIPs = &DeviceIPServiceOp{client: c}
	c.DevicePorts = &DevicePortServiceOp{client: c}
	c.Devices = &DeviceServiceOp{client: c}
	c.Emails = &EmailServiceOp{client: c}
	c.Events = &EventServiceOp{client: c}
	c.Facilities = &FacilityServiceOp{client: c}
	c.HardwareReservations = &HardwareReservationServiceOp{client: c}
	c.Notifications = &NotificationServiceOp{client: c}
	c.OperatingSystems = &OSServiceOp{client: c}
	c.Organizations = &OrganizationServiceOp{client: c}
	c.Plans = &PlanServiceOp{client: c}
	c.ProjectIPs = &ProjectIPServiceOp{client: c}
	c.ProjectVirtualNetworks = &ProjectVirtualNetworkServiceOp{client: c}
	c.Projects = &ProjectServiceOp{client: c}
	c.SSHKeys = &SSHKeyServiceOp{client: c}
	c.SpotMarket = &SpotMarketServiceOp{client: c}
	c.SpotMarketRequests = &SpotMarketRequestServiceOp{client: c}
	c.TwoFactorAuth = &TwoFactorAuthServiceOp{client: c}
	c.Users = &UserServiceOp{client: c}
	c.VirtualCircuits = &VirtualCircuitServiceOp{client: c}
	c.VPN = &VPNServiceOp{client: c}
	c.VolumeAttachments = &VolumeAttachmentServiceOp{client: c}
	c.Volumes = &VolumeServiceOp{client: c}
	c.debug = os.Getenv(debugEnvVar) != ""

	return c, nil
}

func checkResponse(r *http.Response) error {

	if s := r.StatusCode; s >= 200 && s <= 299 {
		// response is good, return
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	// if the response has a body, populate the message in errorResponse
	if err != nil {
		return err
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, expectedAPIContentTypePrefix) {
		errorResponse.SingleError = fmt.Sprintf("Unexpected Content-Type %s with status %s", ct, r.Status)
		return errorResponse
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}
