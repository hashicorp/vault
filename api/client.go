package api

import (
	"net/http"
	"net/url"
)

// Config is used to configure the creation of the client.
type Config struct {
	// Address is the address of the Vault server. This should be a complete
	// URL such as "http://vault.example.com". If you need a custom SSL
	// cert or want to enable insecure mode, you need to specify a custom
	// HttpClient.
	Address string

	// HttpClient is the HTTP client to use. http.DefaultClient will be
	// used if not specified.
	HttpClient *http.Client
}

// DefaultConfig returns a default configuration for the client. It is
// safe to modify the return value of this function.
func DefaultConfig() *Config {
	config := &Config{
		Address:    "https://127.0.0.1:8200",
		HttpClient: http.DefaultClient,
	}

	return config
}

// Client is the client to the Vault API. Create a client with
// NewClient.
type Client struct {
	addr   *url.URL
	config Config
}

// NewClient returns a new client for the given configuration.
func NewClient(c Config) (*Client, error) {
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	return &Client{
		addr:   u,
		config: c,
	}, nil
}

// NewRequest creates a new raw request object to query the Vault server
// configured for this client. This is an advanced method and generally
// doesn't need to be called externally.
func (c *Client) NewRequest(method, path string) *Request {
	return &Request{
		Method: method,
		URL: &url.URL{
			Scheme: c.addr.Scheme,
			Host:   c.addr.Host,
			Path:   path,
		},
		Params: make(map[string][]string),
	}
}

// RawRequest performs the raw request given. This request may be against
// a Vault server not configured with this client. This is an advanced operation
// that generally won't need to be called externally.
func (c *Client) RawRequest(r *Request) (*Response, error) {
	req, err := r.ToHTTP()
	if err != nil {
		return nil, err
	}

	resp, err := c.config.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{Response: resp}, nil
}
