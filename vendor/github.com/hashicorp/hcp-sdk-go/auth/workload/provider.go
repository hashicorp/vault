// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workload

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"
)

// IdentityProviderConfig configures how to source a workload credential and
// exchange it for an HCP Service Principal access token using Workload Identity
// Federation.
type IdentityProviderConfig struct {
	// ProviderResourceName is the resource name of the workload identity
	// provider to exchange the access_token with.
	ProviderResourceName string `json:"provider_resource_name,omitempty"`

	// File sources the subject credential from a file.
	File *FileCredentialSource `json:"file,omitempty"`

	// Token sources the credentials from a directly supplied token.
	Token *CredentialTokenSource `json:"token,omitempty"`

	// EnvironmentVariable sources the subject credential from an environment
	// variable.
	EnvironmentVariable *EnvironmentVariableCredentialSource `json:"env,omitempty"`

	// URL sources the subject credential by making a HTTP request to the
	// provided URL.
	URL *URLCredentialSource `json:"url,omitempty"`

	// AWS uses the IMDS endpoint to retrieve the AWS Caller Identity.
	AWS *AWSCredentialSource `json:"aws,omitempty"`
}

// Validate validates the config.
func (c *IdentityProviderConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("workload identity provider config must not be nil")
	}

	if c.ProviderResourceName == "" {
		return fmt.Errorf("workload identity provider resource name must be set")
	}

	set := 0
	if c.File != nil {
		set++
		if err := c.File.Validate(); err != nil {
			return err
		}
	}

	if c.Token != nil {
		set++
		if err := c.Token.Validate(); err != nil {
			return err
		}
	}

	if c.EnvironmentVariable != nil {
		set++
		if err := c.EnvironmentVariable.Validate(); err != nil {
			return err
		}
	}

	if c.URL != nil {
		set++
		if err := c.URL.Validate(); err != nil {
			return err
		}
	}

	if c.AWS != nil {
		set++
	}

	if set == 0 {
		return fmt.Errorf("a credential source must be configured")
	} else if set > 1 {
		return fmt.Errorf("only one credential source may be configured")
	}

	return nil
}

// Provider sources a workload token and exchanges it for a HCP service
// principal access token. It implements the oauth2.TokenSource interface.
type Provider struct {
	// wipResourceName is the resource name of the workload identity provider to
	// exchange the workload subject token with.
	wipResourceName string

	// jwtProvider is set if the credential source retrieves an opaque JWT
	// token.
	jwtProvider jwtAccessTokenProvider

	// awsProvider is set if the credential source is AWS.
	awsProvider awsCallerIDProvider

	// apiInfo retrieves information on how to access the HCP API
	apiInfo hcpAPIInfo

	// httpClient is used to make requests to the exchange-token endpoint. It
	// should be retrieved using the getHTTPClient method.
	httpClient     *http.Client
	httpClientOnce *sync.Once
}

// New takes an IdentityProviderConfig and returns a Provider or an error if the
// configuration is invalid. The provider can then be used as an auth source
// when creating the HCP Configuration.
func New(c *IdentityProviderConfig) (*Provider, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	// Construct the provider
	p := &Provider{
		wipResourceName: c.ProviderResourceName,
		awsProvider:     c.AWS,
		httpClientOnce:  new(sync.Once),
	}
	if c.URL != nil {
		p.jwtProvider = c.URL
	} else if c.File != nil {
		p.jwtProvider = c.File
	} else if c.EnvironmentVariable != nil {
		p.jwtProvider = c.EnvironmentVariable
	} else if c.Token != nil {
		p.jwtProvider = c.Token
	}

	return p, nil
}

// ResourceName returns the resource name of the provider.
func (p *Provider) ResourceName() string {
	return p.wipResourceName
}

// SetAPI configures the HCP API to use. This will be called by the
// WithWorkloadIdentity helper.
func (p *Provider) SetAPI(info hcpAPIInfo) {
	p.apiInfo = info
}

// Token implements the oauth2.TokenSource interface. It retrieves the workload
// subject token using the configured credential source and then exchanges it
// for the HCP SP access_token.
func (p *Provider) Token() (*oauth2.Token, error) {
	if p.apiInfo == nil {
		return nil, fmt.Errorf("API info must be set before Token() can be called")
	}

	// Get the token
	exchangeReq := &exchangeTokenRequest{}
	if p.jwtProvider != nil {
		token, err := p.jwtProvider.token()
		if err != nil {
			return nil, err
		}
		exchangeReq.JwtToken = token
	} else {
		callerIdentityReq, err := p.awsProvider.getCallerIdentityReq(p.wipResourceName)
		if err != nil {
			return nil, err
		}
		exchangeReq.AwsGetCallerIDToken = callerIdentityReq
	}

	exchangeBytes, err := json.Marshal(exchangeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall exchange token request: %v", err)
	}

	// Make the exchange request
	resp, err := p.getHTTPClient().Post(p.exchangeTokenRequestURL(), "application/json", bytes.NewReader(exchangeBytes))
	if err != nil {
		return nil, fmt.Errorf("invalid response from exchange token endpoint: %v", err)
	}
	defer resp.Body.Close()

	// Read the body
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	// Check that we got a valid response
	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, fmt.Errorf("exchange token status code %d: %s", c, body)
	}

	var exchangeResp exchangeTokenResponse
	err = json.Unmarshal(body, &exchangeResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body from exchange token endpoint: %v", err)
	}

	token := &oauth2.Token{
		AccessToken: exchangeResp.AccessToken,
		Expiry:      time.Now().Add(time.Duration(exchangeResp.ExpiresIn) * time.Second),
	}

	return token, nil
}

// exchangeTokenRequestURL returns the URL for the exchange-token endpoint.
func (p *Provider) exchangeTokenRequestURL() string {
	u := &url.URL{
		Scheme: "https",
		Host:   p.apiInfo.APIAddress(),
		Path:   fmt.Sprintf("/2019-12-10/%s/exchange-token", p.wipResourceName),
	}

	return u.String()
}

// getHTTPClient returns the HTTP Client to use when dialing the HCP API.
func (p *Provider) getHTTPClient() *http.Client {
	p.httpClientOnce.Do(func() {
		transport := cleanhttp.DefaultPooledTransport()
		transport.TLSClientConfig = p.apiInfo.APITLSConfig().Clone()
		p.httpClient = &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		}
	})
	return p.httpClient
}

// exchangeRequest is used to exchange an external subject token for a service
// principal token
type exchangeTokenRequest struct {
	AwsGetCallerIDToken *callerIdentityRequest `json:"aws_get_caller_id_token,omitempty"`
	JwtToken            string                 `json:"jwt_token,omitempty"`
}

// exchangeTokenResponse is the response from the exchange token endpoint.
type exchangeTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// jwtAccessTokenProvider is the interface for providers that return a jwt
// access token.
type jwtAccessTokenProvider interface {
	// token returns the access_token or an error
	token() (string, error)
}

// awsCallerIDProvider is the interface for an AWS provider that return a Caller
// Identity Request.
type awsCallerIDProvider interface {
	// getCallerIdentityReq returns the signed AWS GetCallerIdentity request.
	getCallerIdentityReq(wipResourceName string) (*callerIdentityRequest, error)
}

// hcpAPIInfo returns the API information for accessing the HCP API.
type hcpAPIInfo interface {
	// APIAddress will return the HCP API address (<hostname>[:port]).
	APIAddress() string

	// APITLSConfig will return the API TLS configuration.
	APITLSConfig() *tls.Config
}
