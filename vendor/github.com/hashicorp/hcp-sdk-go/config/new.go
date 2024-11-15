// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/hashicorp/hcp-sdk-go/profile"
	"golang.org/x/oauth2"
)

const (
	// defaultAuthURL is the URL of the production auth endpoint.
	defaultAuthURL = "https://auth.idp.hashicorp.com"

	// defaultOAuth2ClientID is the client ID of the production auth application.
	defaultOAuth2ClientID = "4edd6521-6eb9-4d78-9039-7ce8569d667c"

	// defaultPortalURL is the URL of the production portal.
	defaultPortalURL = "https://portal.cloud.hashicorp.com"

	// defaultAPIAddress is the address of the production API.
	defaultAPIAddress = "api.cloud.hashicorp.com:443"

	// defaultSCADAAddress is the address of the production SCADA endpoint.
	defaultSCADAAddress = "scada.hashicorp.cloud:7224"

	// The audience is the API identifier configured in the auth provider and
	// must be provided when requesting an access token for the API. The value
	// is the same regardless of environment.
	aud = "https://api.hashicorp.cloud"

	// tokenPath is the path used to retrieve the access token.
	tokenPath string = "/oauth2/token"
)

// NewHCPConfig will return a HCPConfig. The configuration will be constructed
// from default values and the passed options.
//
// Based on the provided options the configuration values can be provided
// directly or will be read from other places (e.g. the environment).
//
// The configuration will default to values for the HCP production environment,
// but can be overwritten for development purposes.
//
// In addition to the default values the configuration requires client
// credentials to be provided through one of the passed options (e.g. by using
// WithCredentials or by using FromEnv and providing the client credentials via
// environment variables).
func NewHCPConfig(opts ...HCPConfigOption) (HCPConfig, error) {
	// Parse default URLsG
	authURL, _ := url.Parse(defaultAuthURL)
	portalURL, _ := url.Parse(defaultPortalURL)

	// Prepare basic config with default values.
	config := &hcpConfig{
		authURL:       authURL,
		authTLSConfig: &tls.Config{},
		oauth2Config: oauth2.Config{
			ClientID: defaultOAuth2ClientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:  defaultAuthURL + "/oauth2/auth",
				TokenURL: defaultAuthURL + "/oauth2/token",
			},
			RedirectURL: "http://localhost:8443/oidc/callback",
			Scopes:      []string{"openid", "offline_access"},
		},
		profile: &profile.UserProfile{},

		portalURL: portalURL,

		apiAddress:   defaultAPIAddress,
		apiTLSConfig: &tls.Config{},

		scadaAddress:   defaultSCADAAddress,
		scadaTLSConfig: &tls.Config{},
	}

	if config.suppressLogging {
		log.SetOutput(io.Discard)
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, fmt.Errorf("failed to apply configuration option: %w", err)
		}
	}

	// Configure the token source
	if err := config.setTokenSource(); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("the configuration is not valid: %w", err)
	}

	return config, nil
}
