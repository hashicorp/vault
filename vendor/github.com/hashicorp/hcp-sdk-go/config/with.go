// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/hashicorp/hcp-sdk-go/auth"
	"github.com/hashicorp/hcp-sdk-go/auth/workload"
	"github.com/hashicorp/hcp-sdk-go/profile"
	"golang.org/x/oauth2"
)

// WithClientCredentials credentials is an option that can be used to set
// HCP client credentials on the configuration.
func WithClientCredentials(clientID, clientSecret string) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.clientID = clientID
		config.clientSecret = clientSecret

		return nil
	}
}

// WithWorkloadIdentity exchanges a workload identity provider credentials for
// an HCP Service Principal token. The Workload Identity Provider can be AWS or
// any OIDC based identity provider.
func WithWorkloadIdentity(providerConfig *workload.IdentityProviderConfig) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.workloadProviderConfig = providerConfig

		return nil
	}
}

// WithAPI credentials is an option that can be used to provide a custom
// configuration for the API endpoint.
//
// If nil is provided for the tlsConfig value, TLS will be disabled.
//
// This should only be necessary for development purposes.
func WithAPI(address string, tlsConfig *tls.Config) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.apiAddress = address
		config.apiTLSConfig = tlsConfig.Clone()

		return nil
	}
}

// WithSCADA credentials is an option that can be used to provide a custom
// configuration for the SCADA endpoint.
//
// If nil is provided for the tlsConfig value, TLS will be disabled.
//
// This should only be necessary for development purposes.
func WithSCADA(address string, tlsConfig *tls.Config) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.scadaAddress = address
		config.scadaTLSConfig = tlsConfig.Clone()

		return nil
	}
}

// WithPortalURL credentials is an option that can be used to provide a custom
// URL for the portal.
//
// This should only be necessary for development purposes.
func WithPortalURL(portalURL string) HCPConfigOption {
	return func(config *hcpConfig) error {
		parsedPortalURL, err := url.Parse(portalURL)
		if err != nil {
			return fmt.Errorf("failed to parse portal URL: %w", err)
		}

		config.portalURL = parsedPortalURL

		return nil
	}
}

// WithAuth credentials is an option that can be used to provide a custom URL
// for the auth endpoint.
//
// An alternative TLS configuration can be provided, if none is provided the
// default TLS configuration will be used. It is not possible to disable TLS for
// the auth endpoint.
//
// This should only be necessary for development purposes.
func WithAuth(authURL string, tlsConfig *tls.Config) HCPConfigOption {
	return func(config *hcpConfig) error {
		parsedAuthURL, err := url.Parse(authURL)
		if err != nil {
			return fmt.Errorf("failed to parse auth URL: %w", err)
		}

		// Ensure a TLS configuration is set, as the auth endpoint should always
		// use TLS.
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}

		config.authURL = parsedAuthURL
		config.authTLSConfig = tlsConfig.Clone()

		// Ensure the OAuth2 endpoints are updated with the new auth URL
		config.oauth2Config.Endpoint.AuthURL = authURL + "/oauth2/auth"
		config.oauth2Config.Endpoint.TokenURL = authURL + "/oauth2/token"

		return nil
	}
}

// WithOAuth2ClientID credentials is an option that can be used to provide a
// custom OAuth2 Client ID.
//
// An alternative OAuth2 ClientID can be provided, if none is provided the
// default OAuth2 Client ID will be used.
//
// This should only be necessary for development purposes.
func WithOAuth2ClientID(oauth2ClientID string) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.oauth2Config.ClientID = oauth2ClientID

		return nil
	}
}

// WithProfile is an option that can be used to provide a custom UserProfile struct.
func WithProfile(p *profile.UserProfile) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.profile = p
		return nil
	}
}

// WithTokenSource can be used to set a token source. This should only be necessary for testing.
// Tokens from a custom token source will not be cached.
func WithTokenSource(tokenSource oauth2.TokenSource) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.tokenSource = tokenSource
		return nil
	}
}

// WithoutBrowserLogin disables the automatic opening of the browser login.
func WithoutBrowserLogin() HCPConfigOption {
	return func(config *hcpConfig) error {
		config.noBrowserLogin = true
		return nil
	}
}

// WithoutOpenDefaultBrowser disables opening the default browser when
// browser login is enabled.
func WithoutOpenDefaultBrowser() HCPConfigOption {
	return func(config *hcpConfig) error {
		config.noDefaultBrowser = true
		return nil
	}
}

// WithoutLogging disables this SDK from printing of any kind, this is necessary
// since there is not a consistent logger that is used throughout the project so
// a log level option is not sufficient.
func WithoutLogging() HCPConfigOption {
	return func(config *hcpConfig) error {
		config.suppressLogging = true
		return nil
	}
}

// WithCredentialFile sets the given credential file to be used as an
// authentication source.
func WithCredentialFile(cf *auth.CredentialFile) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.credentialFile = cf
		return config.credentialFile.Validate()
	}
}

// WithCredentialFilePath will search for a credential file at the given path to
// be used as an authentication source.
func WithCredentialFilePath(p string) HCPConfigOption {
	return func(config *hcpConfig) error {
		cf, err := auth.ReadCredentialFile(p)
		config.credentialFile = cf
		return err
	}
}
