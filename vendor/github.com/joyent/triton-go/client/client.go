//
// Copyright 2020 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/errors"
	tritonutils "github.com/joyent/triton-go/utils"
	pkgerrors "github.com/pkg/errors"
)

var (
	ErrDefaultAuth = pkgerrors.New("default SSH agent authentication requires SDC_KEY_ID / TRITON_KEY_ID and SSH_AUTH_SOCK")
	ErrAccountName = pkgerrors.New("missing account name")
	ErrMissingURL  = pkgerrors.New("missing API URL")

	InvalidTritonURL   = "invalid format of Triton URL"
	InvalidMantaURL    = "invalid format of Manta URL"
	InvalidServicesURL = "invalid format of Triton Service Groups URL"
	InvalidDCInURL     = "invalid data center in URL"

	knownDCFormats = []string{
		`https?://(.*).api.joyent.com`,
		`https?://(.*).api.joyentcloud.com`,
		`https?://(.*).api.samsungcloud.io`,
	}

	jpcFormatURL = "https://tsg.%s.svc.joyent.zone"
	spcFormatURL = "https://tsg.%s.svc.samsungcloud.zone"

	tritonTransportHttpTraceChecked = false
	tritonTransportHttpTraceEnabled = false
)

// Client represents a connection to the Triton Compute or Object Storage APIs.
type Client struct {
	HTTPClient    *http.Client
	RequestHeader *http.Header
	Authorizers   []authentication.Signer
	TritonURL     url.URL
	MantaURL      url.URL
	ServicesURL   url.URL
	AccountName   string
	Username      string
}

func isPrivateInstall(url string) bool {
	for _, pattern := range knownDCFormats {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return false
		}
	}

	return true
}

func wrapTritonTransport(in http.RoundTripper) http.RoundTripper {
	if tritonTransportHttpTraceChecked == false {
		tritonTransportHttpTraceChecked = true
		if os.Getenv("TRITON_TRACE_HTTP") != "" {
			tritonTransportHttpTraceEnabled = true
		}
	}

	if tritonTransportHttpTraceEnabled {
		return tritonutils.TraceRoundTripper(in)
	}

	return in
}

// parseDC parses out the data center commonly found in Triton URLs. Returns an
// error if the Triton URL does not include a known data center name, in which
// case a URL override (TRITON_TSG_URL) must be provided.
func parseDC(url string) (string, bool, error) {
	isSamsung := false
	if strings.Contains(url, "samsung") {
		isSamsung = true
	}

	for _, pattern := range knownDCFormats {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1], isSamsung, nil
		}
	}

	return "", isSamsung, fmt.Errorf("failed to parse datacenter from '%s'", url)
}

// New is used to construct a Client in order to make API
// requests to the Triton API.
//
// At least one signer must be provided - example signers include
// authentication.PrivateKeySigner and authentication.SSHAgentSigner.
func New(tritonURL string, mantaURL string, accountName string, signers ...authentication.Signer) (*Client, error) {
	if accountName == "" {
		return nil, ErrAccountName
	}

	if tritonURL == "" && mantaURL == "" {
		return nil, ErrMissingURL
	}

	cloudURL, err := url.Parse(tritonURL)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, InvalidTritonURL)
	}

	storageURL, err := url.Parse(mantaURL)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, InvalidMantaURL)
	}

	// Generate the Services URL (TSG) based on the current datacenter used in
	// the Triton URL (if TritonURL is available). If TRITON_TSG_URL environment
	// variable is available than override using that value instead.
	tsgURL := triton.GetEnv("TSG_URL")
	if tsgURL == "" && tritonURL != "" && !isPrivateInstall(tritonURL) {
		currentDC, isSamsung, err := parseDC(tritonURL)
		if err != nil {
			return nil, pkgerrors.Wrapf(err, InvalidDCInURL)
		}

		tsgURL = fmt.Sprintf(jpcFormatURL, currentDC)
		if isSamsung {
			tsgURL = fmt.Sprintf(spcFormatURL, currentDC)
		}
	}

	servicesURL, err := url.Parse(tsgURL)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, InvalidServicesURL)
	}

	authorizers := make([]authentication.Signer, 0)
	for _, key := range signers {
		if key != nil {
			authorizers = append(authorizers, key)
		}
	}

	skipTLSVerify := triton.GetEnv("SKIP_TLS_VERIFY") != ""

	newClient := &Client{
		HTTPClient: &http.Client{
			Transport:     httpTransport(skipTLSVerify),
			CheckRedirect: doNotFollowRedirects,
		},
		Authorizers: authorizers,
		TritonURL:   *cloudURL,
		MantaURL:    *storageURL,
		ServicesURL: *servicesURL,
		AccountName: accountName,
	}

	// Allow wrapping of the HTTP Transport.
	newClient.HTTPClient.Transport = wrapTritonTransport(newClient.HTTPClient.Transport)

	// Default to constructing an SSHAgentSigner if there are no other signers
	// passed into NewClient and there's an TRITON_KEY_ID and SSH_AUTH_SOCK
	// available in the user's environ(7).
	if len(newClient.Authorizers) == 0 {
		if err := newClient.DefaultAuth(); err != nil {
			return nil, err
		}
	}

	return newClient, nil
}

// initDefaultAuth provides a default key signer for a client. This should only
// be used internally if the client has no other key signer for authenticating
// with Triton. We first look for both `SDC_KEY_ID` and `SSH_AUTH_SOCK` in the
// user's environ(7). If so we default to the SSH agent key signer.
func (c *Client) DefaultAuth() error {
	tritonKeyId := triton.GetEnv("KEY_ID")
	if tritonKeyId != "" {
		input := authentication.SSHAgentSignerInput{
			KeyID:       tritonKeyId,
			AccountName: c.AccountName,
			Username:    c.Username,
		}
		defaultSigner, err := authentication.NewSSHAgentSigner(input)
		if err != nil {
			return pkgerrors.Wrapf(err, "unable to initialize NewSSHAgentSigner")
		}
		c.Authorizers = append(c.Authorizers, defaultSigner)
	}

	return ErrDefaultAuth
}

// InsecureSkipTLSVerify turns off TLS verification for the client connection. This
// allows connection to an endpoint with a certificate which was signed by a non-
// trusted CA, such as self-signed certificates. This can be useful when connecting
// to temporary Triton installations such as Triton Cloud-On-A-Laptop.
func (c *Client) InsecureSkipTLSVerify() {
	if c.HTTPClient == nil {
		return
	}

	c.HTTPClient.Transport = httpTransport(true)
	// Allow wrapping of the HTTP Transport.
	c.HTTPClient.Transport = wrapTritonTransport(c.HTTPClient.Transport)
}

// httpTransport is responsible for setting up our HTTP client's transport
// settings
func httpTransport(insecureSkipTLSVerify bool) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        10,
		IdleConnTimeout:     15 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipTLSVerify,
		},
	}
}

func doNotFollowRedirects(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

// DecodeError decodes a backend Triton error into a more usable Go error type
func (c *Client) DecodeError(resp *http.Response, requestMethod string, consumeBody bool) error {
	err := &errors.APIError{
		StatusCode: resp.StatusCode,
	}

	if requestMethod != http.MethodHead && resp.Body != nil && consumeBody {
		errorDecoder := json.NewDecoder(resp.Body)
		if err := errorDecoder.Decode(err); err != nil {
			return pkgerrors.Wrapf(err, "unable to decode error response")
		}
	}

	if err.Message == "" {
		err.Message = fmt.Sprintf("HTTP response returned status code %d", err.StatusCode)
	}

	return err
}

// overrideHeader overrides the header of the passed in HTTP request
func (c *Client) overrideHeader(req *http.Request) {
	if c.RequestHeader != nil {
		for k := range *c.RequestHeader {
			req.Header.Set(k, c.RequestHeader.Get(k))
		}
	}
}

// resetHeader will reset the struct field that stores custom header
// information
func (c *Client) resetHeader() {
	c.RequestHeader = nil
}

// -----------------------------------------------------------------------------

type RequestInput struct {
	Method  string
	Path    string
	Query   *url.Values
	Headers *http.Header
	Body    interface{}

	// If the response has the HTTP status code 410 (i.e., "Gone"), should we preserve the contents of the body for the caller?
	PreserveGone bool
}

func (c *Client) ExecuteRequestURIParams(ctx context.Context, inputs RequestInput) (io.ReadCloser, error) {
	defer c.resetHeader()

	method := inputs.Method
	path := inputs.Path
	body := inputs.Body
	query := inputs.Query

	var requestBody io.Reader
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	endpoint := c.TritonURL
	endpoint.Path = path
	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to construct HTTP request")
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	// NewClient ensures there's always an authorizer (unless this is called
	// outside that constructor).
	authHeader, err := c.Authorizers[0].Sign(dateHeader, false)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to sign HTTP request")
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Version", triton.CloudAPIMajorVersion)
	req.Header.Set("User-Agent", triton.UserAgent())

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.overrideHeader(req)

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to execute HTTP request")
	}

	// We will only return a response from the API it is in the HTTP StatusCode
	// 2xx range
	// StatusMultipleChoices is StatusCode 300
	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, nil
	}

	return nil, c.DecodeError(resp, req.Method, true)
}

func (c *Client) ExecuteRequest(ctx context.Context, inputs RequestInput) (io.ReadCloser, error) {
	return c.ExecuteRequestURIParams(ctx, inputs)
}

func (c *Client) ExecuteRequestRaw(ctx context.Context, inputs RequestInput) (*http.Response, error) {
	defer c.resetHeader()

	method := inputs.Method
	path := inputs.Path
	body := inputs.Body
	query := inputs.Query

	var requestBody io.Reader
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	endpoint := c.TritonURL
	endpoint.Path = path
	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to construct HTTP request")
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	// NewClient ensures there's always an authorizer (unless this is called
	// outside that constructor).
	authHeader, err := c.Authorizers[0].Sign(dateHeader, false)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to sign HTTP request")
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Version", triton.CloudAPIMajorVersion)
	req.Header.Set("User-Agent", triton.UserAgent())

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.overrideHeader(req)

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to execute HTTP request")
	}

	// We will only return a response from the API it is in the HTTP StatusCode
	// 2xx range
	// StatusMultipleChoices is StatusCode 300
	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		return resp, nil
	}

	// GetMachine returns a HTTP 410 response for deleted instances, but the body of the response is still a valid machine object with a State value of "deleted". Return the object to the caller as well as an error.
	if inputs.PreserveGone && resp.StatusCode == http.StatusGone {
		// Do not consume the response body.
		return resp, c.DecodeError(resp, req.Method, false)
	}

	return nil, c.DecodeError(resp, req.Method, true)
}

func (c *Client) ExecuteRequestStorage(ctx context.Context, inputs RequestInput) (io.ReadCloser, http.Header, error) {
	defer c.resetHeader()

	method := inputs.Method
	path := inputs.Path
	query := inputs.Query
	headers := inputs.Headers
	body := inputs.Body

	endpoint := c.MantaURL
	endpoint.Path = path

	var requestBody io.Reader
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "unable to construct HTTP request")
	}

	if body != nil && (headers == nil || headers.Get("Content-Type") == "") {
		req.Header.Set("Content-Type", "application/json")
	}
	if headers != nil {
		for key, values := range *headers {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := c.Authorizers[0].Sign(dateHeader, true)
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "unable to sign HTTP request")
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", triton.UserAgent())

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	c.overrideHeader(req)

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "unable to execute HTTP request")
	}

	// We will only return a response from the API it is in the HTTP StatusCode
	// 2xx range
	// StatusMultipleChoices is StatusCode 300
	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	return nil, nil, c.DecodeError(resp, req.Method, true)
}

type RequestNoEncodeInput struct {
	Method  string
	Path    string
	Query   *url.Values
	Headers *http.Header
	Body    io.Reader
}

func (c *Client) ExecuteRequestNoEncode(ctx context.Context, inputs RequestNoEncodeInput) (io.ReadCloser, http.Header, error) {
	defer c.resetHeader()

	method := inputs.Method
	path := inputs.Path
	query := inputs.Query
	headers := inputs.Headers
	body := inputs.Body

	endpoint := c.MantaURL
	endpoint.Path = path

	req, err := http.NewRequest(method, endpoint.String(), body)
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "unable to construct HTTP request")
	}

	if headers != nil {
		for key, values := range *headers {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	authHeader, err := c.Authorizers[0].Sign(dateHeader, true)
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "unable to sign HTTP request")
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Version", triton.CloudAPIMajorVersion)
	req.Header.Set("User-Agent", triton.UserAgent())

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	c.overrideHeader(req)

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, pkgerrors.Wrapf(err, "unable to execute HTTP request")
	}

	// We will only return a response from the API it is in the HTTP StatusCode
	// 2xx range
	// StatusMultipleChoices is StatusCode 300
	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, resp.Header, nil
	}

	return nil, nil, c.DecodeError(resp, req.Method, true)
}

func (c *Client) ExecuteRequestTSG(ctx context.Context, inputs RequestInput) (io.ReadCloser, error) {
	defer c.resetHeader()

	method := inputs.Method
	path := inputs.Path
	body := inputs.Body
	query := inputs.Query

	var requestBody io.Reader
	if body != nil {
		marshaled, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	endpoint := c.ServicesURL
	endpoint.Path = path
	if query != nil {
		endpoint.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, endpoint.String(), requestBody)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to construct HTTP request")
	}

	dateHeader := time.Now().UTC().Format(time.RFC1123)
	req.Header.Set("date", dateHeader)

	// NewClient ensures there's always an authorizer (unless this is called
	// outside that constructor).
	authHeader, err := c.Authorizers[0].Sign(dateHeader, false)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to sign HTTP request")
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Version", triton.CloudAPIMajorVersion)
	req.Header.Set("User-Agent", triton.UserAgent())

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.overrideHeader(req)

	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "unable to execute HTTP request")
	}

	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		return resp.Body, nil
	}

	return nil, fmt.Errorf("could not process backend TSG request")
}
