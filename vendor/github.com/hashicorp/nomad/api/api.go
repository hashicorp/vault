// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-rootcerts"
)

var (
	// ClientConnTimeout is the timeout applied when attempting to contact a
	// client directly before switching to a connection through the Nomad
	// server. For cluster topologies where API consumers don't have network
	// access to Nomad clients, set this to a small value (ex 1ms) to avoid
	// pausing on client APIs such as AllocFS.
	ClientConnTimeout = 1 * time.Second
)

const (
	// AllNamespacesNamespace is a sentinel Namespace value to indicate that api should search for
	// jobs and allocations in all the namespaces the requester can access.
	AllNamespacesNamespace = "*"

	// PermissionDeniedErrorContent is the string content of an error returned
	// by the API which indicates the caller does not have permission to
	// perform the action.
	PermissionDeniedErrorContent = "Permission denied"
)

// QueryOptions are used to parametrize a query
type QueryOptions struct {
	// Providing a datacenter overwrites the region provided
	// by the Config
	Region string

	// Namespace is the target namespace for the query.
	Namespace string

	// AllowStale allows any Nomad server (non-leader) to service
	// a read. This allows for lower latency and higher throughput
	AllowStale bool

	// WaitIndex is used to enable a blocking query. Waits
	// until the timeout or the next index is reached
	WaitIndex uint64

	// WaitTime is used to bound the duration of a wait.
	// Defaults to that of the Config, but can be overridden.
	WaitTime time.Duration

	// If set, used as prefix for resource list searches
	Prefix string

	// Set HTTP parameters on the query.
	Params map[string]string

	// Set HTTP headers on the query.
	Headers map[string]string

	// AuthToken is the secret ID of an ACL token
	AuthToken string

	// Filter specifies the go-bexpr filter expression to be used for
	// filtering the data prior to returning a response
	Filter string

	// PerPage is the number of entries to be returned in queries that support
	// paginated lists.
	PerPage int32

	// NextToken is the token used to indicate where to start paging
	// for queries that support paginated lists. This token should be
	// the ID of the next object after the last one seen in the
	// previous response.
	NextToken string

	// Reverse is used to reverse the default order of list results.
	//
	// Currently only supported by specific endpoints.
	Reverse bool

	// ctx is an optional context pass through to the underlying HTTP
	// request layer. Use Context() and WithContext() to manage this.
	ctx context.Context
}

// WriteOptions are used to parametrize a write
type WriteOptions struct {
	// Providing a datacenter overwrites the region provided
	// by the Config
	Region string

	// Namespace is the target namespace for the write.
	Namespace string

	// AuthToken is the secret ID of an ACL token
	AuthToken string

	// Set HTTP headers on the query.
	Headers map[string]string

	// ctx is an optional context pass through to the underlying HTTP
	// request layer. Use Context() and WithContext() to manage this.
	ctx context.Context

	// IdempotencyToken can be used to ensure the write is idempotent.
	IdempotencyToken string
}

// QueryMeta is used to return meta data about a query
type QueryMeta struct {
	// LastIndex. This can be used as a WaitIndex to perform
	// a blocking query
	LastIndex uint64

	// Time of last contact from the leader for the
	// server servicing the request
	LastContact time.Duration

	// Is there a known leader
	KnownLeader bool

	// How long did the request take
	RequestTime time.Duration

	// NextToken is the token used to indicate where to start paging
	// for queries that support paginated lists. To resume paging from
	// this point, pass this token in the next request's QueryOptions
	NextToken string
}

// WriteMeta is used to return meta data about a write
type WriteMeta struct {
	// LastIndex. This can be used as a WaitIndex to perform
	// a blocking query
	LastIndex uint64

	// How long did the request take
	RequestTime time.Duration
}

// HttpBasicAuth is used to authenticate http client with HTTP Basic Authentication
type HttpBasicAuth struct {
	// Username to use for HTTP Basic Authentication
	Username string

	// Password to use for HTTP Basic Authentication
	Password string
}

// Config is used to configure the creation of a client
type Config struct {
	// Address is the address of the Nomad agent
	Address string

	// Region to use. If not provided, the default agent region is used.
	Region string

	// SecretID to use. This can be overwritten per request.
	SecretID string

	// Namespace to use. If not provided the default namespace is used.
	Namespace string

	// HttpClient is the client to use. Default will be used if not provided.
	//
	// If set, it expected to be configured for tls already, and TLSConfig is ignored.
	// You may use ConfigureTLS() function to aid with initialization.
	HttpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	HttpAuth *HttpBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	// TLSConfig provides the various TLS related configurations for the http
	// client.
	//
	// TLSConfig is ignored if HttpClient is set.
	TLSConfig *TLSConfig

	Headers http.Header

	// retryOptions holds the configuration necessary to perform retries
	// on put calls.
	retryOptions *retryOptions

	// url is populated with the initial parsed address and is not modified in the
	// case of a unix:// URL, as opposed to Address.
	url *url.URL
}

// URL returns a copy of the initial parsed address and is not modified in the
// case of a `unix://` URL, as opposed to Address.
func (c *Config) URL() *url.URL {
	return c.url
}

// ClientConfig copies the configuration with a new client address, region, and
// whether the client has TLS enabled.
func (c *Config) ClientConfig(region, address string, tlsEnabled bool) *Config {
	scheme := "http"
	if tlsEnabled {
		scheme = "https"
	}

	config := &Config{
		Address:    fmt.Sprintf("%s://%s", scheme, address),
		Region:     region,
		Namespace:  c.Namespace,
		HttpClient: c.HttpClient,
		SecretID:   c.SecretID,
		HttpAuth:   c.HttpAuth,
		WaitTime:   c.WaitTime,
		TLSConfig:  c.TLSConfig.Copy(),
		url:        copyURL(c.url),
	}

	// Update the tls server name for connecting to a client
	if tlsEnabled && config.TLSConfig != nil {
		config.TLSConfig.TLSServerName = fmt.Sprintf("client.%s.nomad", region)
	}

	return config
}

// TLSConfig contains the parameters needed to configure TLS on the HTTP client
// used to communicate with Nomad.
type TLSConfig struct {
	// CACert is the path to a PEM-encoded CA cert file to use to verify the
	// Nomad server SSL certificate.
	CACert string

	// CAPath is the path to a directory of PEM-encoded CA cert files to verify
	// the Nomad server SSL certificate.
	CAPath string

	// CACertPem is the PEM-encoded CA cert to use to verify the Nomad server
	// SSL certificate.
	CACertPEM []byte

	// ClientCert is the path to the certificate for Nomad communication
	ClientCert string

	// ClientCertPEM is the PEM-encoded certificate for Nomad communication
	ClientCertPEM []byte

	// ClientKey is the path to the private key for Nomad communication
	ClientKey string

	// ClientKeyPEM is the PEM-encoded private key for Nomad communication
	ClientKeyPEM []byte

	// TLSServerName, if set, is used to set the SNI host when connecting via
	// TLS.
	TLSServerName string

	// Insecure enables or disables SSL verification
	Insecure bool
}

func (t *TLSConfig) Copy() *TLSConfig {
	if t == nil {
		return nil
	}

	nt := new(TLSConfig)
	*nt = *t
	return nt
}

// defaultUDSClient creates a unix domain socket client. Errors return a nil
// http.Client, which is tested for in ConfigureTLS. This function expects that
// the Address has already been parsed into the config.url value.
func defaultUDSClient(config *Config) *http.Client {

	config.Address = "http://127.0.0.1"

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", config.url.EscapedPath())
			},
		},
	}
	return defaultClient(httpClient)
}

func defaultHttpClient() *http.Client {
	httpClient := cleanhttp.DefaultPooledClient()
	return defaultClient(httpClient)
}

func defaultClient(c *http.Client) *http.Client {
	transport := c.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Default to http/1: alloc exec/websocket aren't supported in http/2
	// well yet: https://github.com/gorilla/websocket/issues/417
	transport.ForceAttemptHTTP2 = false

	return c
}

// DefaultConfig returns a default configuration for the client
func DefaultConfig() *Config {
	config := &Config{
		Address:   "http://127.0.0.1:4646",
		TLSConfig: &TLSConfig{},
	}
	if addr := os.Getenv("NOMAD_ADDR"); addr != "" {
		config.Address = addr
	}
	if v := os.Getenv("NOMAD_REGION"); v != "" {
		config.Region = v
	}
	if v := os.Getenv("NOMAD_NAMESPACE"); v != "" {
		config.Namespace = v
	}
	if auth := os.Getenv("NOMAD_HTTP_AUTH"); auth != "" {
		var username, password string
		if strings.Contains(auth, ":") {
			split := strings.SplitN(auth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = auth
		}

		config.HttpAuth = &HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}

	// Read TLS specific env vars
	if v := os.Getenv("NOMAD_CACERT"); v != "" {
		config.TLSConfig.CACert = v
	}
	if v := os.Getenv("NOMAD_CAPATH"); v != "" {
		config.TLSConfig.CAPath = v
	}
	if v := os.Getenv("NOMAD_CLIENT_CERT"); v != "" {
		config.TLSConfig.ClientCert = v
	}
	if v := os.Getenv("NOMAD_CLIENT_KEY"); v != "" {
		config.TLSConfig.ClientKey = v
	}
	if v := os.Getenv("NOMAD_TLS_SERVER_NAME"); v != "" {
		config.TLSConfig.TLSServerName = v
	}
	if v := os.Getenv("NOMAD_SKIP_VERIFY"); v != "" {
		if insecure, err := strconv.ParseBool(v); err == nil {
			config.TLSConfig.Insecure = insecure
		}
	}
	if v := os.Getenv("NOMAD_TOKEN"); v != "" {
		config.SecretID = v
	}
	return config
}

// cloneWithTimeout returns a cloned httpClient with set timeout if positive;
// otherwise, returns the same client
func cloneWithTimeout(httpClient *http.Client, t time.Duration) (*http.Client, error) {
	if httpClient == nil {
		return nil, errors.New("nil HTTP client")
	} else if httpClient.Transport == nil {
		return nil, errors.New("nil HTTP client transport")
	}

	if t.Nanoseconds() < 0 {
		return httpClient, nil
	}

	tr, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("unexpected HTTP transport: %T", httpClient.Transport)
	}

	// copy all public fields, to avoid copying transient state and locks
	ntr := &http.Transport{
		Proxy:                  tr.Proxy,
		DialContext:            tr.DialContext,
		Dial:                   tr.Dial,
		DialTLS:                tr.DialTLS,
		TLSClientConfig:        tr.TLSClientConfig,
		TLSHandshakeTimeout:    tr.TLSHandshakeTimeout,
		DisableKeepAlives:      tr.DisableKeepAlives,
		DisableCompression:     tr.DisableCompression,
		MaxIdleConns:           tr.MaxIdleConns,
		MaxIdleConnsPerHost:    tr.MaxIdleConnsPerHost,
		MaxConnsPerHost:        tr.MaxConnsPerHost,
		IdleConnTimeout:        tr.IdleConnTimeout,
		ResponseHeaderTimeout:  tr.ResponseHeaderTimeout,
		ExpectContinueTimeout:  tr.ExpectContinueTimeout,
		TLSNextProto:           tr.TLSNextProto,
		ProxyConnectHeader:     tr.ProxyConnectHeader,
		MaxResponseHeaderBytes: tr.MaxResponseHeaderBytes,
	}

	// apply timeout
	ntr.DialContext = (&net.Dialer{
		Timeout:   t,
		KeepAlive: 30 * time.Second,
	}).DialContext

	// clone http client with new transport
	nc := *httpClient
	nc.Transport = ntr
	return &nc, nil
}

// ConfigureTLS applies a set of TLS configurations to the HTTP client.
func ConfigureTLS(httpClient *http.Client, tlsConfig *TLSConfig) error {
	if tlsConfig == nil {
		return nil
	}
	if httpClient == nil {
		return errors.New("config HTTP Client must be set")
	}

	var clientCert tls.Certificate
	foundClientCert := false
	if tlsConfig.ClientCert != "" || tlsConfig.ClientKey != "" {
		if tlsConfig.ClientCert != "" && tlsConfig.ClientKey != "" {
			var err error
			clientCert, err = tls.LoadX509KeyPair(tlsConfig.ClientCert, tlsConfig.ClientKey)
			if err != nil {
				return err
			}
			foundClientCert = true
		} else {
			return errors.New("Both client cert and client key must be provided")
		}
	} else if len(tlsConfig.ClientCertPEM) != 0 || len(tlsConfig.ClientKeyPEM) != 0 {
		if len(tlsConfig.ClientCertPEM) != 0 && len(tlsConfig.ClientKeyPEM) != 0 {
			var err error
			clientCert, err = tls.X509KeyPair(tlsConfig.ClientCertPEM, tlsConfig.ClientKeyPEM)
			if err != nil {
				return err
			}
			foundClientCert = true
		} else {
			return errors.New("Both client cert and client key must be provided")
		}
	}

	clientTLSConfig := httpClient.Transport.(*http.Transport).TLSClientConfig
	rootConfig := &rootcerts.Config{
		CAFile:        tlsConfig.CACert,
		CAPath:        tlsConfig.CAPath,
		CACertificate: tlsConfig.CACertPEM,
	}
	if err := rootcerts.ConfigureTLS(clientTLSConfig, rootConfig); err != nil {
		return err
	}

	clientTLSConfig.InsecureSkipVerify = tlsConfig.Insecure

	if foundClientCert {
		clientTLSConfig.Certificates = []tls.Certificate{clientCert}
	}
	if tlsConfig.TLSServerName != "" {
		clientTLSConfig.ServerName = tlsConfig.TLSServerName
	}

	return nil
}

// Client provides a client to the Nomad API
type Client struct {
	httpClient *http.Client
	config     Config
}

// NewClient returns a new client
func NewClient(config *Config) (*Client, error) {
	var err error
	// bootstrap the config
	defConfig := DefaultConfig()

	if config.Address == "" {
		config.Address = defConfig.Address
	}

	// we have to test the address that comes from DefaultConfig, because it
	// could be the value of NOMAD_ADDR which is applied without testing. But
	// only on the first use of this Config, otherwise we'll have mutated the
	// address
	if config.url == nil {
		if config.url, err = url.Parse(config.Address); err != nil {
			return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
		}
	}

	httpClient := config.HttpClient
	if httpClient == nil {
		switch {
		case config.url.Scheme == "unix":
			httpClient = defaultUDSClient(config) // mutates config
		default:
			httpClient = defaultHttpClient()
		}

		if err := ConfigureTLS(httpClient, config.TLSConfig); err != nil {
			return nil, err
		}
	}

	client := &Client{
		config:     *config,
		httpClient: httpClient,
	}
	return client, nil
}

// Close closes the client's idle keep-alived connections. The default
// client configuration uses keep-alive to maintain connections and
// you should instantiate a single Client and reuse it for all
// requests from the same host. Connections will be closed
// automatically once the client is garbage collected. If you are
// creating multiple clients on the same host (for example, for
// testing), it may be useful to call Close() to avoid hitting
// connection limits.
func (c *Client) Close() {
	c.httpClient.CloseIdleConnections()
}

// Address return the address of the Nomad agent
func (c *Client) Address() string {
	return c.config.Address
}

// SetRegion sets the region to forward API requests to.
func (c *Client) SetRegion(region string) {
	c.config.Region = region
}

// SetNamespace sets the namespace to forward API requests to.
func (c *Client) SetNamespace(namespace string) {
	c.config.Namespace = namespace
}

// GetNodeClient returns a new Client that will dial the specified node. If the
// QueryOptions is set, its region will be used.
func (c *Client) GetNodeClient(nodeID string, q *QueryOptions) (*Client, error) {
	return c.getNodeClientImpl(nodeID, -1, q, c.Nodes().Info)
}

// GetNodeClientWithTimeout returns a new Client that will dial the specified
// node using the specified timeout. If the QueryOptions is set, its region will
// be used.
func (c *Client) GetNodeClientWithTimeout(
	nodeID string, timeout time.Duration, q *QueryOptions) (*Client, error) {
	return c.getNodeClientImpl(nodeID, timeout, q, c.Nodes().Info)
}

// nodeLookup is the definition of a function used to lookup a node. This is
// largely used to mock the lookup in tests.
type nodeLookup func(nodeID string, q *QueryOptions) (*Node, *QueryMeta, error)

// getNodeClientImpl is the implementation of creating a API client for
// contacting a node. It takes a function to lookup the node such that it can be
// mocked during tests.
func (c *Client) getNodeClientImpl(nodeID string, timeout time.Duration, q *QueryOptions, lookup nodeLookup) (*Client, error) {
	node, _, err := lookup(nodeID, q)
	if err != nil {
		return nil, err
	}
	if node.Status == "down" {
		return nil, NodeDownErr
	}
	if node.HTTPAddr == "" {
		return nil, fmt.Errorf("http addr of node %q (%s) is not advertised", node.Name, nodeID)
	}

	var region string
	switch {
	case q != nil && q.Region != "":
		// Prefer the region set in the query parameter
		region = q.Region
	case c.config.Region != "":
		// If the client is configured for a particular region use that
		region = c.config.Region
	default:
		// No region information is given so use GlobalRegion as the default.
		region = GlobalRegion
	}

	// Get an API client for the node
	conf := c.config.ClientConfig(region, node.HTTPAddr, node.TLSEnabled)

	// set timeout - preserve old behavior where errors are ignored and use untimed one
	httpClient, err := cloneWithTimeout(c.httpClient, timeout)
	// on error, fallback to using current http client
	if err != nil {
		httpClient = c.httpClient
	}
	conf.HttpClient = httpClient

	return NewClient(conf)
}

// SetSecretID sets the ACL token secret for API requests.
func (c *Client) SetSecretID(secretID string) {
	c.config.SecretID = secretID
}

func (c *Client) configureRetries(ro *retryOptions) {

	c.config.retryOptions = &retryOptions{
		maxRetries:      defaultNumberOfRetries,
		maxBackoffDelay: defaultMaxBackoffDelay,
		delayBase:       defaultDelayTimeBase,
	}

	if ro.delayBase != 0 {
		c.config.retryOptions.delayBase = ro.delayBase
	}

	if ro.maxRetries != defaultNumberOfRetries {
		c.config.retryOptions.maxRetries = ro.maxRetries
	}

	if ro.maxBackoffDelay != 0 {
		c.config.retryOptions.maxBackoffDelay = ro.maxBackoffDelay
	}

	if ro.maxToLastCall != 0 {
		c.config.retryOptions.maxToLastCall = ro.maxToLastCall
	}

	if ro.fixedDelay != 0 {
		c.config.retryOptions.fixedDelay = ro.fixedDelay
	}

	// Ensure that a big attempt number or a big delayBase number will not cause
	// a negative delay by overflowing the delay increase.
	c.config.retryOptions.maxValidAttempt = int64(math.Log2(float64(math.MaxInt64 /
		c.config.retryOptions.delayBase.Nanoseconds())))
}

// request is used to help build up a request
type request struct {
	config *Config
	method string
	url    *url.URL
	params url.Values
	token  string
	body   io.Reader
	obj    interface{}
	ctx    context.Context
	header http.Header
}

// setQueryOptions is used to annotate the request with
// additional query options
func (r *request) setQueryOptions(q *QueryOptions) {
	if q == nil {
		return
	}
	if q.Region != "" {
		r.params.Set("region", q.Region)
	}
	if q.Namespace != "" {
		r.params.Set("namespace", q.Namespace)
	}
	if q.AuthToken != "" {
		r.token = q.AuthToken
	}
	if q.AllowStale {
		r.params.Set("stale", "")
	}
	if q.WaitIndex != 0 {
		r.params.Set("index", strconv.FormatUint(q.WaitIndex, 10))
	}
	if q.WaitTime != 0 {
		r.params.Set("wait", durToMsec(q.WaitTime))
	}
	if q.Prefix != "" {
		r.params.Set("prefix", q.Prefix)
	}
	if q.Filter != "" {
		r.params.Set("filter", q.Filter)
	}
	if q.PerPage != 0 {
		r.params.Set("per_page", fmt.Sprint(q.PerPage))
	}
	if q.NextToken != "" {
		r.params.Set("next_token", q.NextToken)
	}
	if q.Reverse {
		r.params.Set("reverse", "true")
	}
	for k, v := range q.Params {
		r.params.Set(k, v)
	}
	r.ctx = q.Context()

	for k, v := range q.Headers {
		r.header.Set(k, v)
	}
}

// durToMsec converts a duration to a millisecond specified string
func durToMsec(dur time.Duration) string {
	return fmt.Sprintf("%dms", dur/time.Millisecond)
}

// setWriteOptions is used to annotate the request with
// additional write options
func (r *request) setWriteOptions(q *WriteOptions) {
	if q == nil {
		return
	}
	if q.Region != "" {
		r.params.Set("region", q.Region)
	}
	if q.Namespace != "" {
		r.params.Set("namespace", q.Namespace)
	}
	if q.AuthToken != "" {
		r.token = q.AuthToken
	}
	if q.IdempotencyToken != "" {
		r.params.Set("idempotency_token", q.IdempotencyToken)
	}
	r.ctx = q.Context()

	for k, v := range q.Headers {
		r.header.Set(k, v)
	}
}

// toHTTP converts the request to an HTTP request
func (r *request) toHTTP() (*http.Request, error) {
	// Encode the query parameters
	r.url.RawQuery = r.params.Encode()

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		if b, err := encodeBody(r.obj); err != nil {
			return nil, err
		} else {
			r.body = b
		}
	}

	ctx := func() context.Context {
		if r.ctx != nil {
			return r.ctx
		}
		return context.Background()
	}()

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

	req.Header = r.header

	// Optionally configure HTTP basic authentication
	if r.url.User != nil {
		username := r.url.User.Username()
		password, _ := r.url.User.Password()
		req.SetBasicAuth(username, password)
	} else if r.config.HttpAuth != nil {
		req.SetBasicAuth(r.config.HttpAuth.Username, r.config.HttpAuth.Password)
	}

	req.Header.Add("Accept-Encoding", "gzip")
	if r.token != "" {
		req.Header.Set("X-Nomad-Token", r.token)
	}

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host
	return req, nil
}

// newRequest is used to create a new request
func (c *Client) newRequest(method, path string) (*request, error) {

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	r := &request{
		config: &c.config,
		method: method,
		url: &url.URL{
			Scheme:  c.config.url.Scheme,
			User:    c.config.url.User,
			Host:    c.config.url.Host,
			Path:    u.Path,
			RawPath: u.RawPath,
		},
		header: make(http.Header),
		params: make(map[string][]string),
	}

	// fixup socket paths
	if r.url.Scheme == "unix" {
		r.url.Scheme = "http"
		r.url.Host = "127.0.0.1"
	}

	if c.config.Region != "" {
		r.params.Set("region", c.config.Region)
	}
	if c.config.Namespace != "" {
		r.params.Set("namespace", c.config.Namespace)
	}
	if c.config.WaitTime != 0 {
		r.params.Set("wait", durToMsec(r.config.WaitTime))
	}
	if c.config.SecretID != "" {
		r.token = r.config.SecretID
	}

	// Add in the query parameters, if any
	for key, values := range u.Query() {
		for _, value := range values {
			r.params.Add(key, value)
		}
	}

	for key, values := range c.config.Headers {
		r.header[key] = values
	}

	return r, nil
}

// multiCloser is to wrap a ReadCloser such that when close is called, multiple
// Closes occur.
type multiCloser struct {
	reader       io.Reader
	inorderClose []io.Closer
}

func (m *multiCloser) Close() error {
	for _, c := range m.inorderClose {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiCloser) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

// doRequest runs a request with our client
func (c *Client) doRequest(r *request) (time.Duration, *http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return 0, nil, err
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	diff := time.Since(start)

	// If the response is compressed, we swap the body's reader.
	if zipErr := c.autoUnzip(resp); zipErr != nil {
		return 0, nil, zipErr
	}

	return diff, resp, err
}

// autoUnzip modifies resp in-place, wrapping the response body with a gzip
// reader if the Content-Encoding of the response is "gzip".
func (*Client) autoUnzip(resp *http.Response) error {
	if resp == nil || resp.Header == nil {
		return nil
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		zReader, err := gzip.NewReader(resp.Body)
		if err == io.EOF {
			// zero length response, do not wrap
			return nil
		} else if err != nil {
			// some other error (e.g. corrupt)
			return err
		}

		// The gzip reader does not close an underlying reader, so use a
		// multiCloser to make sure response body does get closed.
		resp.Body = &multiCloser{
			reader:       zReader,
			inorderClose: []io.Closer{zReader, resp.Body},
		}
	}

	return nil
}

// rawQuery makes a GET request to the specified endpoint but returns just the
// response body.
func (c *Client) rawQuery(endpoint string, q *QueryOptions) (io.ReadCloser, error) {
	r, err := c.newRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	_, resp, err := requireOK(c.doRequest(r)) //nolint:bodyclose // Closing the body is the caller's responsibility.
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// websocket makes a websocket request to the specific endpoint
func (c *Client) websocket(endpoint string, q *QueryOptions) (*websocket.Conn, *http.Response, error) {

	transport, ok := c.httpClient.Transport.(*http.Transport)
	if !ok {
		return nil, nil, errors.New("unsupported transport")
	}
	dialer := websocket.Dialer{
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
		HandshakeTimeout: c.httpClient.Timeout,

		// values to inherit from http client configuration
		NetDial:         transport.Dial,
		NetDialContext:  transport.DialContext,
		Proxy:           transport.Proxy,
		TLSClientConfig: transport.TLSClientConfig,
	}

	// build request object for header and parameters
	r, err := c.newRequest("GET", endpoint)
	if err != nil {
		return nil, nil, err
	}
	r.setQueryOptions(q)

	rhttp, err := r.toHTTP()
	if err != nil {
		return nil, nil, err
	}

	// convert scheme
	wsScheme := ""
	switch rhttp.URL.Scheme {
	case "http":
		wsScheme = "ws"
	case "https":
		wsScheme = "wss"
	default:
		return nil, nil, fmt.Errorf("unsupported scheme: %v", rhttp.URL.Scheme)
	}
	rhttp.URL.Scheme = wsScheme

	conn, resp, err := dialer.Dial(rhttp.URL.String(), rhttp.Header)

	// check resp status code, as it's more informative than handshake error we get from ws library
	if resp != nil {
		switch resp.StatusCode {
		case http.StatusSwitchingProtocols:
			// Connection upgrade was successful.

		case http.StatusPermanentRedirect, http.StatusTemporaryRedirect, http.StatusMovedPermanently:
			loc := resp.Header.Get("Location")
			u, err := url.Parse(loc)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid redirect location %q: %w", loc, err)
			}
			return c.websocket(u.Path, q)

		default:
			var buf bytes.Buffer

			if resp.Header.Get("Content-Encoding") == "gzip" {
				greader, err := gzip.NewReader(resp.Body)
				if err != nil {
					return nil, nil, newUnexpectedResponseError(
						fromStatusCode(resp.StatusCode),
						withExpectedStatuses([]int{http.StatusSwitchingProtocols}),
						withError(err))
				}
				_, _ = io.Copy(&buf, greader)
			} else {
				_, _ = io.Copy(&buf, resp.Body)
			}
			_ = resp.Body.Close()

			return nil, nil, newUnexpectedResponseError(
				fromStatusCode(resp.StatusCode),
				withExpectedStatuses([]int{http.StatusSwitchingProtocols}),
				withBody(buf.String()),
			)
		}
	}

	return conn, resp, err
}

// query is used to do a GET request against an endpoint
// and deserialize the response into an interface using
// standard Nomad conventions.
func (c *Client) query(endpoint string, out any, q *QueryOptions) (*QueryMeta, error) {
	r, err := c.newRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(c.doRequest(r)) //nolint:bodyclose // Closing the body is the caller's responsibility.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	qm := &QueryMeta{}
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	if err := decodeBody(resp, out); err != nil {
		return nil, err
	}
	return qm, nil
}

// putQuery is used to do a PUT request when doing a "write" to a Client RPC.
// Client RPCs must use QueryOptions to allow setting AllowStale=true.
func (c *Client) putQuery(endpoint string, in, out any, q *QueryOptions) (*QueryMeta, error) {
	r, err := c.newRequest("PUT", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r)) //nolint:bodyclose // Closing the body is the caller's responsibility.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	qm := &QueryMeta{}
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	if err := decodeBody(resp, out); err != nil {
		return nil, err
	}
	return qm, nil
}

// put is used to do a PUT request against an endpoint and
// serialize/deserialized using the standard Nomad conventions.
func (c *Client) put(endpoint string, in, out any, q *WriteOptions) (*WriteMeta, error) {
	return c.write(http.MethodPut, endpoint, in, out, q)
}

// postQuery is used to do a POST request when doing a "write" to a Client RPC.
// Client RPCs must use QueryOptions to allow setting AllowStale=true.
func (c *Client) postQuery(endpoint string, in, out any, q *QueryOptions) (*QueryMeta, error) {
	r, err := c.newRequest("POST", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r)) //nolint:bodyclose // Closing the body is the caller's responsibility.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	qm := &QueryMeta{}
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	if err := decodeBody(resp, out); err != nil {
		return nil, err
	}
	return qm, nil
}

// post is used to do a POST request against an endpoint and
// serialize/deserialized using the standard Nomad conventions.
func (c *Client) post(endpoint string, in, out any, q *WriteOptions) (*WriteMeta, error) {
	return c.write(http.MethodPost, endpoint, in, out, q)
}

// write is used to do a write request against an endpoint and
// serialize/deserialized using the standard Nomad conventions.
//
// You probably want the delete, post, or put methods.
func (c *Client) write(verb, endpoint string, in, out any, q *WriteOptions) (*WriteMeta, error) {
	r, err := c.newRequest(verb, endpoint)
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r)) //nolint:bodyclose // Closing the body is the caller's responsibility.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	wm := &WriteMeta{RequestTime: rtt}
	parseWriteMeta(resp, wm)

	if out != nil {
		if err := decodeBody(resp, &out); err != nil {
			return nil, err
		}
	}
	return wm, nil
}

// delete is used to do a DELETE request against an endpoint and
// serialize/deserialized using the standard Nomad conventions.
func (c *Client) delete(endpoint string, in, out any, q *WriteOptions) (*WriteMeta, error) {
	r, err := c.newRequest("DELETE", endpoint)
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r)) //nolint:bodyclose // Closing the body is the caller's responsibility.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	wm := &WriteMeta{RequestTime: rtt}
	parseWriteMeta(resp, wm)

	if out != nil {
		if err := decodeBody(resp, &out); err != nil {
			return nil, err
		}
	}
	return wm, nil
}

// parseQueryMeta is used to help parse query meta-data
func parseQueryMeta(resp *http.Response, q *QueryMeta) error {
	header := resp.Header

	// Parse the X-Nomad-Index
	index, err := strconv.ParseUint(header.Get("X-Nomad-Index"), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse X-Nomad-Index: %v", err)
	}
	q.LastIndex = index

	// Parse the X-Nomad-LastContact
	last, err := strconv.ParseUint(header.Get("X-Nomad-LastContact"), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse X-Nomad-LastContact: %v", err)
	}
	if last > math.MaxInt64 {
		return fmt.Errorf("Last contact duration is out of range: %d", last)
	}
	q.LastContact = time.Duration(last) * time.Millisecond
	q.NextToken = header.Get("X-Nomad-NextToken")

	// Parse the X-Nomad-KnownLeader
	switch header.Get("X-Nomad-KnownLeader") {
	case "true":
		q.KnownLeader = true
	default:
		q.KnownLeader = false
	}
	return nil
}

// parseWriteMeta is used to help parse write meta-data
func parseWriteMeta(resp *http.Response, q *WriteMeta) error {
	header := resp.Header

	// Parse the X-Nomad-Index
	index, err := strconv.ParseUint(header.Get("X-Nomad-Index"), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse X-Nomad-Index: %v", err)
	}
	q.LastIndex = index
	return nil
}

// decodeBody is used to JSON decode a body
func decodeBody(resp *http.Response, out interface{}) error {
	switch resp.ContentLength {
	case 0:
		if out == nil {
			return nil
		}
		return errors.New("Got 0 byte response with non-nil decode object")
	default:
		dec := json.NewDecoder(resp.Body)
		return dec.Decode(out)
	}
}

// encodeBody prepares the reader to serve as the request body.
//
// Returns the `obj` input if it is a raw io.Reader object; otherwise
// returns a reader of the json format of the passed argument.
func encodeBody(obj interface{}) (io.Reader, error) {
	if reader, ok := obj.(io.Reader); ok {
		return reader, nil
	}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

// Context returns the context used for canceling HTTP requests related to this query
func (o *QueryOptions) Context() context.Context {
	if o != nil && o.ctx != nil {
		return o.ctx
	}
	return context.Background()
}

// WithContext creates a copy of the query options using the provided context to cancel related HTTP requests
func (o *QueryOptions) WithContext(ctx context.Context) *QueryOptions {
	o2 := new(QueryOptions)
	if o != nil {
		*o2 = *o
	}
	o2.ctx = ctx
	return o2
}

// Context returns the context used for canceling HTTP requests related to this write
func (o *WriteOptions) Context() context.Context {
	if o != nil && o.ctx != nil {
		return o.ctx
	}
	return context.Background()
}

// WithContext creates a copy of the write options using the provided context to cancel related HTTP requests
func (o *WriteOptions) WithContext(ctx context.Context) *WriteOptions {
	o2 := new(WriteOptions)
	if o != nil {
		*o2 = *o
	}
	o2.ctx = ctx
	return o2
}

// copyURL makes a deep copy of a net/url.URL
func copyURL(u1 *url.URL) *url.URL {
	if u1 == nil {
		return nil
	}
	o := *u1
	if o.User != nil {
		ou := *u1.User
		o.User = &ou
	}
	return &o
}
