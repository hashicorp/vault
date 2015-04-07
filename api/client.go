package api

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	vaultHttp "github.com/hashicorp/vault/http"
)

// Config is used to configure the creation of the client.
type Config struct {
	// Address is the address of the Vault server. This should be a complete
	// URL such as "http://vault.example.com". If you need a custom SSL
	// cert or want to enable insecure mode, you need to specify a custom
	// HttpClient.
	Address string

	// HttpClient is the HTTP client to use. http.DefaultClient will be
	// used if not specified. The HTTP client must have the cookie jar set
	// to be able to store cookies, otherwise authentication (login) will
	// not work properly. If the jar is nil, a default empty cookie jar
	// will be set.
	HttpClient *http.Client
}

// DefaultConfig returns a default configuration for the client. It is
// safe to modify the return value of this function.
func DefaultConfig() Config {
	config := Config{
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

	// Make a copy of the HTTP client so we can configure it without
	// affecting the original
	//
	// If no cookie jar is set on the client, we set a default empty
	// cookie jar.
	if c.HttpClient.Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{})
		if err != nil {
			return nil, err
		}

		c.HttpClient.Jar = jar
	}

	return &Client{
		addr:   u,
		config: c,
	}, nil
}

// Token returns the access token being used by this client. It will
// return the empty string if there is no token set.
func (c *Client) Token() string {
	r := c.NewRequest("GET", "/")
	for _, cookie := range c.config.HttpClient.Jar.Cookies(r.URL) {
		if cookie.Name == vaultHttp.AuthCookieName {
			return cookie.Value
		}
	}

	return ""
}

// SetToken sets the token directly. This won't perform any auth
// verification, it simply sets the cookie properly for future requests.
func (c *Client) SetToken(v string) {
	r := c.NewRequest("GET", "/")
	c.config.HttpClient.Jar.SetCookies(r.URL, []*http.Cookie{
		&http.Cookie{
			Name:    vaultHttp.AuthCookieName,
			Value:   v,
			Path:    "/",
			Expires: time.Now().Add(365 * 24 * time.Hour),
		},
	})
}

// ClearToken deletes the token cookie if it is set or does nothing otherwise.
func (c *Client) ClearToken() {
	r := c.NewRequest("GET", "/")
	c.config.HttpClient.Jar.SetCookies(r.URL, []*http.Cookie{
		&http.Cookie{
			Name:    vaultHttp.AuthCookieName,
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
		},
	})
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

	var result *Response
	resp, err := c.config.HttpClient.Do(req)
	if resp != nil {
		result = &Response{Response: resp}
	}
	if err != nil {
		return result, err
	}

	if err := result.Error(); err != nil {
		return result, err
	}

	return result, nil
}
