package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	errRedirect            = errors.New("redirect")
	defaultHTTPClientSetup sync.Once
	defaultHTTPClient      = &http.Client{}
)

// Config is used to configure the creation of the client.
type Config struct {
	// Address is the address of the Vault server. This should be a complete
	// URL such as "http://vault.example.com". If you need a custom SSL
	// cert or want to enable insecure mode, you need to specify a custom
	// HttpClient.
	Address string

	// HttpClient is the HTTP client to use, which will currently always have the
	// same values as http.DefaultClient. This is used to control redirect behavior.
	HttpClient *http.Client
}

// DefaultConfig returns a default configuration for the client. It is
// safe to modify the return value of this function.
//
// The default Address is https://127.0.0.1:8200, but this can be overridden by
// setting the `VAULT_ADDR` environment variable.
func DefaultConfig() *Config {
	config := &Config{
		Address:    "https://127.0.0.1:8200",
		HttpClient: defaultHTTPClient,
	}

	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		config.Address = addr
	}

	return config
}

// Client is the client to the Vault API. Create a client with
// NewClient.
type Client struct {
	addr   *url.URL
	config *Config
	token  string
}

// NewClient returns a new client for the given configuration.
//
// If the environment variable `VAULT_TOKEN` is present, the token will be
// automatically added to the client. Otherwise, you must manually call
// `SetToken()`.
func NewClient(c *Config) (*Client, error) {
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	if c.HttpClient == defaultHTTPClient {
		defaultHTTPClientSetup.Do(func() {
			// Ensure redirects are not automatically followed
			c.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return errRedirect
			}
		})
	}

	client := &Client{
		addr:   u,
		config: c,
	}

	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	}

	return client, nil
}

// Token returns the access token being used by this client. It will
// return the empty string if there is no token set.
func (c *Client) Token() string {
	return c.token
}

// SetToken sets the token directly. This won't perform any auth
// verification, it simply sets the token properly for future requests.
func (c *Client) SetToken(v string) {
	c.token = v
}

// ClearToken deletes the token if it is set or does nothing otherwise.
func (c *Client) ClearToken() {
	c.token = ""
}

// NewRequest creates a new raw request object to query the Vault server
// configured for this client. This is an advanced method and generally
// doesn't need to be called externally.
func (c *Client) NewRequest(method, path string) *Request {
	req := &Request{
		Method: method,
		URL: &url.URL{
			Scheme: c.addr.Scheme,
			Host:   c.addr.Host,
			Path:   path,
		},
		ClientToken: c.token,
		Params:      make(map[string][]string),
	}

	return req
}

// RawRequest performs the raw request given. This request may be against
// a Vault server not configured with this client. This is an advanced operation
// that generally won't need to be called externally.
func (c *Client) RawRequest(r *Request) (*Response, error) {
	redirectCount := 0
START:
	req, err := r.ToHTTP()
	if err != nil {
		return nil, err
	}

	var result *Response
	resp, err := c.config.HttpClient.Do(req)
	if resp != nil {
		result = &Response{Response: resp}
	}
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && urlErr.Err == errRedirect {
			err = nil
		} else if strings.Contains(err.Error(), "tls: oversized") {
			err = fmt.Errorf(
				"%s\n\n"+
					"This error usually means that the server is running with TLS disabled\n"+
					"but the client is configured to use TLS. Please either enable TLS\n"+
					"on the server or run the client with -address set to an address\n"+
					"that uses the http protocol:\n\n"+
					"    vault <command> -address http://<address>\n\n"+
					"You can also set the VAULT_ADDR environment variable:\n\n\n"+
					"    VAULT_ADDR=http://<address> vault <command>\n\n"+
					"where <address> is replaced by the actual address to the server.",
				err)
		}
	}
	if err != nil {
		return result, err
	}

	// Check for a redirect, only allowing for a single redirect
	if (resp.StatusCode == 302 || resp.StatusCode == 307) && redirectCount == 0 {
		// Parse the updated location
		respLoc, err := resp.Location()
		if err != nil {
			return result, err
		}

		// Ensure a protocol downgrade doesn't happen
		if req.URL.Scheme == "https" && respLoc.Scheme != "https" {
			return result, fmt.Errorf("redirect would cause protocol downgrade")
		}

		// Update the request
		r.URL = respLoc

		// Reset the request body if any
		if err := r.ResetJSONBody(); err != nil {
			return result, err
		}

		// Retry the request
		redirectCount++
		goto START
	}

	if err := result.Error(); err != nil {
		return result, err
	}

	return result, nil
}
