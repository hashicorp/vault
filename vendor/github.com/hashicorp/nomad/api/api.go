package api

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	rootcerts "github.com/hashicorp/go-rootcerts"
)

// QueryOptions are used to parameterize a query
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

	// SecretID is the secret ID of an ACL token
	SecretID string
}

// WriteOptions are used to parameterize a write
type WriteOptions struct {
	// Providing a datacenter overwrites the region provided
	// by the Config
	Region string

	// Namespace is the target namespace for the write.
	Namespace string

	// SecretID is the secret ID of an ACL token
	SecretID string
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

	// httpClient is the client to use. Default will be used if not provided.
	httpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	HttpAuth *HttpBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime time.Duration

	// TLSConfig provides the various TLS related configurations for the http
	// client
	TLSConfig *TLSConfig
}

// ClientConfig copies the configuration with a new client address, region, and
// whether the client has TLS enabled.
func (c *Config) ClientConfig(region, address string, tlsEnabled bool) *Config {
	scheme := "http"
	if tlsEnabled {
		scheme = "https"
	}
	defaultConfig := DefaultConfig()
	config := &Config{
		Address:    fmt.Sprintf("%s://%s", scheme, address),
		Region:     region,
		Namespace:  c.Namespace,
		httpClient: defaultConfig.httpClient,
		SecretID:   c.SecretID,
		HttpAuth:   c.HttpAuth,
		WaitTime:   c.WaitTime,
		TLSConfig:  c.TLSConfig.Copy(),
	}
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

	// ClientCert is the path to the certificate for Nomad communication
	ClientCert string

	// ClientKey is the path to the private key for Nomad communication
	ClientKey string

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

// DefaultConfig returns a default configuration for the client
func DefaultConfig() *Config {
	config := &Config{
		Address:    "http://127.0.0.1:4646",
		httpClient: cleanhttp.DefaultClient(),
		TLSConfig:  &TLSConfig{},
	}
	transport := config.httpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
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

// ConfigureTLS applies a set of TLS configurations to the the HTTP client.
func (c *Config) ConfigureTLS() error {
	if c.TLSConfig == nil {
		return nil
	}
	if c.httpClient == nil {
		return fmt.Errorf("config HTTP Client must be set")
	}

	var clientCert tls.Certificate
	foundClientCert := false
	if c.TLSConfig.ClientCert != "" || c.TLSConfig.ClientKey != "" {
		if c.TLSConfig.ClientCert != "" && c.TLSConfig.ClientKey != "" {
			var err error
			clientCert, err = tls.LoadX509KeyPair(c.TLSConfig.ClientCert, c.TLSConfig.ClientKey)
			if err != nil {
				return err
			}
			foundClientCert = true
		} else {
			return fmt.Errorf("Both client cert and client key must be provided")
		}
	}

	clientTLSConfig := c.httpClient.Transport.(*http.Transport).TLSClientConfig
	rootConfig := &rootcerts.Config{
		CAFile: c.TLSConfig.CACert,
		CAPath: c.TLSConfig.CAPath,
	}
	if err := rootcerts.ConfigureTLS(clientTLSConfig, rootConfig); err != nil {
		return err
	}

	clientTLSConfig.InsecureSkipVerify = c.TLSConfig.Insecure

	if foundClientCert {
		clientTLSConfig.Certificates = []tls.Certificate{clientCert}
	}
	if c.TLSConfig.TLSServerName != "" {
		clientTLSConfig.ServerName = c.TLSConfig.TLSServerName
	}

	return nil
}

// Client provides a client to the Nomad API
type Client struct {
	config Config
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

	if config.httpClient == nil {
		config.httpClient = defConfig.httpClient
	}

	// Configure the TLS cofigurations
	if err := config.ConfigureTLS(); err != nil {
		return nil, err
	}

	client := &Client{
		config: *config,
	}
	return client, nil
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
	return c.getNodeClientImpl(nodeID, q, c.Nodes().Info)
}

// nodeLookup is the definition of a function used to lookup a node. This is
// largely used to mock the lookup in tests.
type nodeLookup func(nodeID string, q *QueryOptions) (*Node, *QueryMeta, error)

// getNodeClientImpl is the implementation of creating a API client for
// contacting a node. It takes a function to lookup the node such that it can be
// mocked during tests.
func (c *Client) getNodeClientImpl(nodeID string, q *QueryOptions, lookup nodeLookup) (*Client, error) {
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
		// No region information is given so use the default.
		region = "global"
	}

	// Get an API client for the node
	conf := c.config.ClientConfig(region, node.HTTPAddr, node.TLSEnabled)
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
	if q.SecretID != "" {
		r.token = q.SecretID
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
	if q.SecretID != "" {
		r.token = q.SecretID
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

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

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
			Scheme: base.Scheme,
			User:   base.User,
			Host:   base.Host,
			Path:   u.Path,
		},
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
	resp, err := c.config.httpClient.Do(req)
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
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
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
