// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"

	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/version"
)

// Config contains the client's configuration options
type Config struct {
	// HCPConfig contains values that allow to interact with HCP APIs.
	config.HCPConfig

	// HostPath is the host name of the HCP API, without a scheme prefix.
	//
	// Deprecated: HCPConfig should be used instead
	HostPath string `json:"host_path"`

	// AuthURL is the URL of the authentication provider, inclusive of the 'https' scheme prefix.
	//
	// Deprecated: HCPConfig should be used instead
	AuthURL string `json:"auth_url"`

	// ClientID is the client ID of a Service Principal Key.
	//
	// Deprecated: HCPConfig should be used instead
	ClientID string `json:"client_id"`

	// ClientSecret is the client secret of a Service Principal Key, only provided on creation.
	//
	// Deprecated: HCPConfig should be used instead
	ClientSecret string `json:"client_secret"`

	// SourceChannel denotes the client (channel) that originiated the request.
	// This is synonymous to a user-agent.
	SourceChannel string `json:"source_channel"`

	// Client allows passing a custom http.Client to be used instead of the
	// cleanhttp default one for both Auth and API requests. This should be used only for testing
	// in testing to provide the httptest.Server's custom client that will trust
	// its TLS cert.
	//
	// Deprecated: HCPConfig should be used instead
	Client *http.Client
}

// New creates a client with the right base path to connect to any HCP API
func New(cfg Config) (runtime *httptransport.Runtime, err error) {
	// Populate default values where possible.
	if err = cfg.Canonicalize(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create a transport using the API TLS config.
	tlsTransport := cleanhttp.DefaultPooledTransport()
	tlsTransport.TLSClientConfig = cfg.APITLSConfig()

	// Wrap the transport in an oauth2.Transport.
	var transport http.RoundTripper = &oauth2.Transport{
		Base:   tlsTransport,
		Source: cfg,
	}

	var opts []MiddlewareOption

	if cfg.SourceChannel != "" {
		// Use custom transport in order to set the source channel header when it is present.
		sc := fmt.Sprintf("%s hcp-go-sdk/%s", cfg.SourceChannel, version.Version)

		opts = append(opts, withSourceChannel(sc))
	}

	if cfg.Profile().OrganizationID != "" && cfg.Profile().ProjectID != "" {

		opts = append(opts, withOrgAndProjectIDs(cfg.Profile().OrganizationID, cfg.Profile().ProjectID))
	}

	transport = &roundTripperWithMiddleware{
		OriginalRoundTripper: transport,
		MiddlewareOptions:    opts,
	}

	// Set the scheme based on the TLS configuration.
	scheme := "https"
	if cfg.APITLSConfig() == nil {
		scheme = "http"
	}

	// Create a runtime that can be used by the generated clients.
	runtime = httptransport.New(cfg.APIAddress(), "", []string{scheme})
	runtime.Transport = transport

	return runtime, nil
}

// Canonicalize populates default values for config fields that are otherwise unset.
func (c *Config) Canonicalize() error {
	// When a HCPConfig is provided it is expected to be canonical and
	// deprecated configuration values will be ignored.
	if c.HCPConfig != nil {
		return nil
	}

	tlsConfig := &tls.Config{}

	// Use the test client's TLS configuration, if one was provided
	if c.Client != nil {
		httpTransport, ok := c.Client.Transport.(*http.Transport)
		if ok {
			tlsConfig = httpTransport.TLSClientConfig
		}
	}

	// Otherwise, construct a HCPConfig
	options := []config.HCPConfigOption{config.FromEnv()}

	if c.HostPath != "" {
		// Allow https:// prefix on HostPath even though it's the only scheme we allow
		// as it's more natural to support the URL. Any other scheme we don't strip
		// which will fail validation.
		if strings.HasPrefix(strings.ToLower(c.HostPath), "https://") {
			c.HostPath = c.HostPath[8:]
		}

		options = append(options, config.WithAPI(c.HostPath, tlsConfig))
	}

	if c.AuthURL != "" {
		options = append(options, config.WithAuth(c.AuthURL, tlsConfig))
	}

	if c.ClientID != "" && c.ClientSecret != "" {
		options = append(options, config.WithClientCredentials(c.ClientID, c.ClientSecret))
	}

	var err error
	c.HCPConfig, err = config.NewHCPConfig(options...)
	if err != nil {
		return fmt.Errorf("failed to construct HCP config: %w", err)
	}

	return nil
}
