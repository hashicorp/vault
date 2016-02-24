package api

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const EnvVaultAddress = "VAULT_ADDR"
const EnvVaultCACert = "VAULT_CACERT"
const EnvVaultCAPath = "VAULT_CAPATH"
const EnvVaultClientCert = "VAULT_CLIENT_CERT"
const EnvVaultClientKey = "VAULT_CLIENT_KEY"
const EnvVaultInsecure = "VAULT_SKIP_VERIFY"
const EnvVaultTLSServerName = "VAULT_TLS_SERVER_NAME"

var (
	errRedirect = errors.New("redirect")
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

	redirectSetup sync.Once
}

// DefaultConfig returns a default configuration for the client. It is
// safe to modify the return value of this function.
//
// The default Address is https://127.0.0.1:8200, but this can be overridden by
// setting the `VAULT_ADDR` environment variable.
func DefaultConfig() *Config {
	config := &Config{
		Address: "https://127.0.0.1:8200",

		HttpClient: cleanhttp.DefaultClient(),
	}
	config.HttpClient.Timeout = time.Second * 60
	transport := config.HttpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if v := os.Getenv(EnvVaultAddress); v != "" {
		config.Address = v
	}

	return config
}

// ReadEnvironment reads configuration information from the
// environment. If there is an error, no configuration value
// is updated.
func (c *Config) ReadEnvironment() error {
	var envAddress string
	var envCACert string
	var envCAPath string
	var envClientCert string
	var envClientKey string
	var envInsecure bool
	var foundInsecure bool
	var envTLSServerName string

	var newCertPool *x509.CertPool
	var clientCert tls.Certificate
	var foundClientCert bool

	if v := os.Getenv(EnvVaultAddress); v != "" {
		envAddress = v
	}
	if v := os.Getenv(EnvVaultCACert); v != "" {
		envCACert = v
	}
	if v := os.Getenv(EnvVaultCAPath); v != "" {
		envCAPath = v
	}
	if v := os.Getenv(EnvVaultClientCert); v != "" {
		envClientCert = v
	}
	if v := os.Getenv(EnvVaultClientKey); v != "" {
		envClientKey = v
	}
	if v := os.Getenv(EnvVaultInsecure); v != "" {
		var err error
		envInsecure, err = strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("Could not parse VAULT_SKIP_VERIFY")
		}
		foundInsecure = true
	}
	if v := os.Getenv(EnvVaultTLSServerName); v != "" {
		envTLSServerName = v
	}
	// If we need custom TLS configuration, then set it
	if envCACert != "" || envCAPath != "" || envClientCert != "" || envClientKey != "" || envInsecure {
		var err error
		if envCACert != "" {
			newCertPool, err = LoadCACert(envCACert)
		} else if envCAPath != "" {
			newCertPool, err = LoadCAPath(envCAPath)
		}
		if err != nil {
			return fmt.Errorf("Error setting up CA path: %s", err)
		}

		if envClientCert != "" && envClientKey != "" {
			clientCert, err = tls.LoadX509KeyPair(envClientCert, envClientKey)
			if err != nil {
				return err
			}
			foundClientCert = true
		} else if envClientCert != "" || envClientKey != "" {
			return fmt.Errorf("Both client cert and client key must be provided")
		}
	}

	if envAddress != "" {
		c.Address = envAddress
	}

	clientTLSConfig := c.HttpClient.Transport.(*http.Transport).TLSClientConfig
	if foundInsecure {
		clientTLSConfig.InsecureSkipVerify = envInsecure
	}
	if newCertPool != nil {
		clientTLSConfig.RootCAs = newCertPool
	}
	if foundClientCert {
		clientTLSConfig.Certificates = []tls.Certificate{clientCert}
	}
	if envTLSServerName != "" {
		clientTLSConfig.ServerName = envTLSServerName
	}

	return nil
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

	if c.HttpClient == nil {
		c.HttpClient = DefaultConfig().HttpClient
	}

	redirFunc := func() {
		// Ensure redirects are not automatically followed
		// Note that this is sane for the API client as it has its own
		// redirect handling logic (and thus also for command/meta),
		// but in e.g. http_test actual redirect handling is necessary
		c.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return errRedirect
		}
	}

	c.redirectSetup.Do(redirFunc)

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
	if (resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 307) && redirectCount == 0 {
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

// Loads the certificate from given path and creates a certificate pool from it.
func LoadCACert(path string) (*x509.CertPool, error) {
	certs, err := loadCertFromPEM(path)
	if err != nil {
		return nil, err
	}

	result := x509.NewCertPool()
	for _, cert := range certs {
		result.AddCert(cert)
	}

	return result, nil
}

// Loads the certificates present in the given directory and creates a
// certificate pool from it.
func LoadCAPath(path string) (*x509.CertPool, error) {
	result := x509.NewCertPool()
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		certs, err := loadCertFromPEM(path)
		if err != nil {
			return err
		}

		for _, cert := range certs {
			result.AddCert(cert)
		}
		return nil
	}

	return result, filepath.Walk(path, fn)
}

// Creates a certificate from the given path
func loadCertFromPEM(path string) ([]*x509.Certificate, error) {
	pemCerts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certs := make([]*x509.Certificate, 0, 5)
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		certs = append(certs, cert)
	}

	return certs, nil
}
