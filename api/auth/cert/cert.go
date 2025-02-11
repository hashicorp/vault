package cert

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-rootcerts"
	"github.com/hashicorp/vault/api"
)

type CertAuth struct {
	role               string
	caCert             string
	caCertBytes        []byte
	clientCert         string
	clientKey          string
	insecureSkipVerify bool
}

var _ api.AuthMethod = (*CertAuth)(nil)

type LoginOption func(a *CertAuth) error

// NewCertAuth initializes a new cert auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithClientCertAndKey, WithInsecure
//
// https://developer.hashicorp.com/vault/api-docs/auth/cert#login-with-tls-certificate-method
func NewCertAuth(roleName string, opts ...LoginOption) (*CertAuth, error) {
	a := &CertAuth{
		role: roleName,
	}

	for _, opt := range opts {
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	return a, nil
}

// Login sets up the required request body for the cert auth method's /login
// endpoint, and performs a write to it.
// It adds the client cert and key to the request.
func (a *CertAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	c, err := a.httpClient()
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"name": a.role,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal login data: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/cert/login", client.Address())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with cert auth: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unable to log in with cert auth, response code: %d. response body: %s", resp.StatusCode, string(body))
	}

	var secret api.Secret
	if err := json.Unmarshal(body, &secret); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	return &secret, nil
}

func (a *CertAuth) httpClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(a.clientCert, a.clientKey)
	if err != nil {
		return nil, fmt.Errorf("unable to load cert: %w", err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: a.insecureSkipVerify,
		Certificates:       []tls.Certificate{cert},
	}

	if a.caCert != "" || len(a.caCertBytes) > 0 {
		err = rootcerts.ConfigureTLS(tlsConfig, &rootcerts.Config{
			CAPath:        a.caCert,
			CACertificate: a.caCertBytes,
		})

		if err != nil {
			return nil, fmt.Errorf("unable to configure TLS: %w", err)
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}, nil
}

// WithCACert sets the CA cert to be used for the login request.
// caCert is the path to the CA cert file.
func WithCACert(caCert string) LoginOption {
	return func(a *CertAuth) error {
		a.caCert = caCert
		return nil
	}
}

// WithCACertBytes sets the CA cert to be used for the login request.
// caCertBytes is the bytes of the CA cert.
// caCertBytes takes precedence over caCert.
func WithCACertBytes(caCertBytes []byte) LoginOption {
	return func(a *CertAuth) error {
		a.caCertBytes = caCertBytes
		return nil
	}
}

// WithClientCertAndKey sets the client cert and key to be used for the login request.
func WithClientCertAndKey(clientCert, clientKey string) LoginOption {
	return func(a *CertAuth) error {
		a.clientCert = clientCert
		a.clientKey = clientKey
		return nil
	}
}

// WithInsecure skips the verification of the server's certificate chain and host name.
func WithInsecure() LoginOption {
	return func(a *CertAuth) error {
		a.insecureSkipVerify = true
		return nil
	}
}
