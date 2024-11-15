package packngo

import (
	"bytes"
	"context"
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
	baseURL         = "https://api.equinix.com/metal/v1/"
	mediaType       = "application/json"
	debugEnvVar     = "PACKNGO_DEBUG"

	headerRateLimit              = "X-RateLimit-Limit"
	headerRateRemaining          = "X-RateLimit-Remaining"
	headerRateReset              = "X-RateLimit-Reset"
	expectedAPIContentTypePrefix = "application/json"
)

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
	apiKeySet     bool
	header        http.Header

	RateLimit Rate

	// Equinix Metal Api Objects
	APIKeys                APIKeyService
	BGPConfig              BGPConfigService
	BGPSessions            BGPSessionService
	Batches                BatchService
	CapacityService        CapacityService
	Connections            ConnectionService
	DeviceIPs              DeviceIPService
	Devices                DeviceService
	Emails                 EmailService
	Events                 EventService
	Facilities             FacilityService
	HardwareReservations   HardwareReservationService
	Invitations            InvitationService
	Members                MemberService
	Metros                 MetroService
	Notifications          NotificationService
	OperatingSystems       OSService
	Organizations          OrganizationService
	Plans                  PlanService
	Ports                  PortService
	ProjectIPs             ProjectIPService
	ProjectVirtualNetworks ProjectVirtualNetworkService
	Projects               ProjectService
	SSHKeys                SSHKeyService
	SpotMarket             SpotMarketService
	SpotMarketRequests     SpotMarketRequestService
	MetalGateways          MetalGatewayService
	TwoFactorAuth          TwoFactorAuthService
	Users                  UserService
	VirtualCircuits        VirtualCircuitService
	VLANAssignments        VLANAssignmentService
	VolumeAttachments      VolumeAttachmentService
	Volumes                VolumeService
	VRFs                   VRFService

	// DevicePorts
	//
	// Deprecated: Use Client.Ports or Device methods
	DevicePorts DevicePortService
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

	req.Header = c.header.Clone()
	req.Header.Set("X-Auth-Token", c.APIKey)
	req.Header.Set("X-Consumer-Token", c.ConsumerToken)
	req.Header.Set("User-Agent", c.UserAgent)

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
	dumpDeprecation(response.Response)
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

// dumpDeprecation logs headers defined by
// https://tools.ietf.org/html/rfc8594
func dumpDeprecation(resp *http.Response) {
	uri := ""
	if resp.Request != nil {
		uri = resp.Request.Method + " " + resp.Request.URL.Path
	}

	deprecation := resp.Header.Get("Deprecation")
	if deprecation != "" {
		if deprecation == "true" {
			deprecation = ""
		} else {
			deprecation = " on " + deprecation
		}
		log.Printf("WARNING: %q reported deprecation%s", uri, deprecation)
	}

	sunset := resp.Header.Get("Sunset")
	if sunset != "" {
		log.Printf("WARNING: %q reported sunsetting on %s", uri, sunset)
	}

	links := resp.Header.Values("Link")

	for _, s := range links {
		for _, ss := range strings.Split(s, ",") {
			if strings.Contains(ss, "rel=\"sunset\"") {
				link := strings.Split(ss, ";")[0]
				log.Printf("WARNING: See %s for sunset details", link)
			} else if strings.Contains(ss, "rel=\"deprecation\"") {
				link := strings.Split(ss, ";")[0]
				log.Printf("WARNING: See %s for deprecation details", link)
			}
		}
	}
}

// from terraform-plugin-sdk/v2/helper/logging/transport.go
func prettyPrintJsonLines(b []byte) string {
	parts := strings.Split(string(b), "\n")
	for i, p := range parts {
		if b := []byte(p); json.Valid(b) {
			var out bytes.Buffer
			_ = json.Indent(&out, b, "", " ")
			parts[i] = out.String()
		}
	}
	return strings.Join(parts, "\n")
}

func dumpResponse(resp *http.Response) {
	o, _ := httputil.DumpResponse(resp, true)
	strResp := prettyPrintJsonLines(o)
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
	reqBodyStr := prettyPrintJsonLines(bbs)
	strReq := prettyPrintJsonLines(o)
	log.Printf("\n=======[REQUEST]=============\n%s%s\n", string(strReq), reqBodyStr)
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

// NewClientWithAuth initializes and returns a Client, use this to get an API Client to operate on
// N.B.: Equinix Metal's API certificate requires Go 1.5+ to successfully parse. If you are using
// an older version of Go, pass in a custom http.Client with a custom TLS configuration
// that sets "InsecureSkipVerify" to "true"
func NewClientWithAuth(consumerToken string, apiKey string, httpClient *http.Client) *Client {
	client, _ := NewClientWithBaseURL(consumerToken, apiKey, httpClient, baseURL)
	return client
}

// NewClientWithBaseURL returns a Client pointing to nonstandard API URL, e.g.
// for mocking the remote API
func NewClientWithBaseURL(consumerToken string, apiKey string, httpClient *http.Client, apiBaseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return NewClient(WithAuth(consumerToken, apiKey), WithHTTPClient(httpClient), WithBaseURL(apiBaseURL))
}

// NewClient initializes and returns a Client. The opts are functions such as WithAuth,
// WithHTTPClient, etc.
//
// An example:
//
//	c, err := NewClient()
//
// An alternative example, which avoids reading PACKET_AUTH_TOKEN environment variable:
//
//	c, err := NewClient(WithAuth("packngo lib", packetAuthToken))
func NewClient(opts ...ClientOpt) (*Client, error) {
	// set defaults, then let caller override them
	c := &Client{
		client:        http.DefaultClient,
		UserAgent:     UserAgent,
		ConsumerToken: "packngo lib",
		header:        http.Header{},
	}

	c.header.Set("Content-Type", mediaType)
	c.header.Set("Accept", mediaType)

	var err error
	c.BaseURL, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

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
	c.Invitations = &InvitationServiceOp{client: c}
	c.Members = &MemberServiceOp{client: c}
	c.Metros = &MetroServiceOp{client: c}
	c.Notifications = &NotificationServiceOp{client: c}
	c.OperatingSystems = &OSServiceOp{client: c}
	c.Organizations = &OrganizationServiceOp{client: c}
	c.Plans = &PlanServiceOp{client: c}
	c.Ports = &PortServiceOp{client: c}
	c.ProjectIPs = &ProjectIPServiceOp{client: c}
	c.ProjectVirtualNetworks = &ProjectVirtualNetworkServiceOp{client: c}
	c.Projects = &ProjectServiceOp{client: c}
	c.SSHKeys = &SSHKeyServiceOp{client: c}
	c.SpotMarket = &SpotMarketServiceOp{client: c}
	c.SpotMarketRequests = &SpotMarketRequestServiceOp{client: c}
	c.MetalGateways = &MetalGatewayServiceOp{client: c}
	c.TwoFactorAuth = &TwoFactorAuthServiceOp{client: c}
	c.Users = &UserServiceOp{client: c}
	c.VirtualCircuits = &VirtualCircuitServiceOp{client: c}
	c.VolumeAttachments = &VolumeAttachmentServiceOp{client: c}
	c.Volumes = &VolumeServiceOp{client: c}
	c.VRFs = &VRFServiceOp{client: c}
	c.VLANAssignments = &VLANAssignmentServiceOp{client: c}
	c.debug = os.Getenv(debugEnvVar) != ""

	for _, fn := range opts {
		err := fn(c)
		if err != nil {
			return nil, err
		}
	}

	if !c.apiKeySet {
		c.APIKey = os.Getenv(authTokenEnvVar)

		if c.APIKey == "" {
			return nil, fmt.Errorf("you must export %s", authTokenEnvVar)
		}

		c.apiKeySet = true
	}

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
