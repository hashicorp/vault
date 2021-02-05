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
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	rootcerts "github.com/hashicorp/go-rootcerts"
)

var (
	// ClientConnTimeout is the timeout applied when attempting to contact a
	// client directly before switching to a connection through the Nomad
	// server.
	ClientConnTimeout = 1 * time.Second
)

const (
	// AllNamespacesNamespace is a sentinel Namespace value to indicate that api should search for
	// jobs and allocations in all the namespaces the requester can access.
	AllNamespacesNamespace = "*"
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

	// AuthToken is the secret ID of an ACL token
	AuthToken string

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

	// ctx is an optional context pass through to the underlying HTTP
	// request layer. Use Context() and WithContext() to manage this.
	ctx context.Context
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

func defaultHttpClient() *http.Client {
	httpClient := cleanhttp.DefaultClient()
	transport := httpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	return httpClient
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
		return nil, fmt.Errorf("nil HTTP client")
	} else if httpClient.Transport == nil {
		return nil, fmt.Errorf("nil HTTP client transport")
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

// ConfigureTLS applies a set of TLS configurations to the the HTTP client.
func ConfigureTLS(httpClient *http.Client, tlsConfig *TLSConfig) error {
	if tlsConfig == nil {
		return nil
	}
	if httpClient == nil {
		return fmt.Errorf("config HTTP Client must be set")
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
			return fmt.Errorf("Both client cert and client key must be provided")
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
			return fmt.Errorf("Both client cert and client key must be provided")
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
	// bootstrap the config
	defConfig := DefaultConfig()

	if config.Address == "" {
		config.Address = defConfig.Address
	} else if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
	}

	httpClient := config.HttpClient
	if httpClient == nil {
		httpClient = defaultHttpClient()
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
	for k, v := range q.Params {
		r.params.Set(k, v)
	}
	r.ctx = q.Context()
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
	r.ctx = q.Context()
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
	base, _ := url.Parse(c.config.Address)
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	r := &request{
		config: &c.config,
		method: method,
		url: &url.URL{
			Scheme:  base.Scheme,
			User:    base.User,
			Host:    base.Host,
			Path:    u.Path,
			RawPath: u.RawPath,
		},
		header: make(http.Header),
		params: make(map[string][]string),
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

	if c.config.Headers != nil {
		r.header = c.config.Headers
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
	diff := time.Now().Sub(start)

	// If the response is compressed, we swap the body's reader.
	if resp != nil && resp.Header != nil {
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			greader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return 0, nil, err
			}

			// The gzip reader doesn't close the wrapped reader so we use
			// multiCloser.
			reader = &multiCloser{
				reader:       greader,
				inorderClose: []io.Closer{greader, resp.Body},
			}
		default:
			reader = resp.Body
		}
		resp.Body = reader
	}

	return diff, resp, err
}

// rawQuery makes a GET request to the specified endpoint but returns just the
// response body.
func (c *Client) rawQuery(endpoint string, q *QueryOptions) (io.ReadCloser, error) {
	r, err := c.newRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	_, resp, err := requireOK(c.doRequest(r))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// websocket makes a websocket request to the specific endpoint
func (c *Client) websocket(endpoint string, q *QueryOptions) (*websocket.Conn, *http.Response, error) {

	transport, ok := c.httpClient.Transport.(*http.Transport)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported transport")
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
	if resp != nil && resp.StatusCode != 101 {
		var buf bytes.Buffer

		if resp.Header.Get("Content-Encoding") == "gzip" {
			greader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return nil, nil, fmt.Errorf("Unexpected response code: %d", resp.StatusCode)
			}
			io.Copy(&buf, greader)
		} else {
			io.Copy(&buf, resp.Body)
		}
		resp.Body.Close()

		return nil, nil, fmt.Errorf("Unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}

	return conn, resp, err
}

// query is used to do a GET request against an endpoint
// and deserialize the response into an interface using
// standard Nomad conventions.
func (c *Client) query(endpoint string, out interface{}, q *QueryOptions) (*QueryMeta, error) {
	r, err := c.newRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(c.doRequest(r))
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

// putQuery is used to do a PUT request when doing a read against an endpoint
// and deserialize the response into an interface using standard Nomad
// conventions.
func (c *Client) putQuery(endpoint string, in, out interface{}, q *QueryOptions) (*QueryMeta, error) {
	r, err := c.newRequest("PUT", endpoint)
	if err != nil {
		return nil, err
	}
	r.setQueryOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r))
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

// write is used to do a PUT request against an endpoint
// and serialize/deserialized using the standard Nomad conventions.
func (c *Client) write(endpoint string, in, out interface{}, q *WriteOptions) (*WriteMeta, error) {
	r, err := c.newRequest("PUT", endpoint)
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)
	r.obj = in
	rtt, resp, err := requireOK(c.doRequest(r))
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

// delete is used to do a DELETE request against an endpoint
// and serialize/deserialized using the standard Nomad conventions.
func (c *Client) delete(endpoint string, out interface{}, q *WriteOptions) (*WriteMeta, error) {
	r, err := c.newRequest("DELETE", endpoint)
	if err != nil {
		return nil, err
	}
	r.setWriteOptions(q)
	rtt, resp, err := requireOK(c.doRequest(r))
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
	q.LastContact = time.Duration(last) * time.Millisecond

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

// requireOK is used to wrap doRequest and check for a 200
func requireOK(d time.Duration, resp *http.Response, e error) (time.Duration, *http.Response, error) {
	if e != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return d, nil, e
	}
	if resp.StatusCode != 200 {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		return d, nil, fmt.Errorf("Unexpected response code: %d (%s)", resp.StatusCode, buf.Bytes())
	}
	return d, resp, nil
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
