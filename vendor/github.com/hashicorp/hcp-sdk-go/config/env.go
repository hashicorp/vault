// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcp-sdk-go/profile"
)

// The following constants contain the names of environment variables that can
// be set to provide configuration values.
const (
	envVarAuthURL        = "HCP_AUTH_URL"
	envVarOAuth2ClientID = "HCP_OAUTH_CLIENT_ID"

	envVarClientID     = "HCP_CLIENT_ID"
	envVarClientSecret = "HCP_CLIENT_SECRET"

	envVarPortalURL = "HCP_PORTAL_URL"

	envVarAPIAddress = "HCP_API_ADDRESS"

	// envVarAPIHostnameLegacy is necessary as `HCP_API_HOST` has been used in
	// the past and needs to be supported for backward compatibility.
	envVarAPIHostnameLegacy = "HCP_API_HOST"

	envVarAPITLS = "HCP_API_TLS"

	envVarAuthTLS = "HCP_AUTH_TLS"

	envVarSCADAAddress = "HCP_SCADA_ADDRESS"

	envVarSCADATLS = "HCP_SCADA_TLS"

	envVarHCPOrganizationID = "HCP_ORGANIZATION_ID"

	envVarHCPProjectID = "HCP_PROJECT_ID"
)

const (
	// tlsSettingInsecure is the value TLS environment variables can be set to
	// if the verification of the server certificate should be skipped.
	//
	// This should only be needed for development purposes.
	tlsSettingInsecure = "insecure"

	// tlsSettingDisabled is the value TLS environment variables can be set to
	// if the communication should happen in plain-text and TLS should be
	// disabled.
	//
	// This should only be needed for development purposes.
	tlsSettingDisabled = "disabled"
)

// FromEnv will return a HCPConfigOption that will populate the configuration
// with values from the environment.
//
// It will not fail if no or only part of the variables are present.
func FromEnv() HCPConfigOption {
	return func(config *hcpConfig) error {

		// Read client credentials from the environment, the values will only be
		// used if both are provided.
		clientID, clientIDOK := os.LookupEnv(envVarClientID)
		clientSecret, clientSecretOK := os.LookupEnv(envVarClientSecret)

		if clientIDOK && clientSecretOK {
			if err := apply(config, WithClientCredentials(clientID, clientSecret)); err != nil {
				return fmt.Errorf("failed to set client credentials from environment variables (%s, %s): %w", envVarClientID, envVarClientSecret, err)
			}
		}

		// Read auth URL from environment
		if authURL, ok := os.LookupEnv(envVarAuthURL); ok {
			if err := apply(config, WithAuth(authURL, nil)); err != nil {
				return fmt.Errorf("failed to parse environment variable %s: %w", envVarAuthURL, err)
			}
		}

		// Read oauth2ClientID from environment
		if oauth2ClientID, ok := os.LookupEnv(envVarOAuth2ClientID); ok {
			if err := apply(config, WithOAuth2ClientID(oauth2ClientID)); err != nil {
				return fmt.Errorf("failed to parse environment variable %s: %w", envVarOAuth2ClientID, err)
			}
		}

		// Read portal URL from environment
		if portalURL, ok := os.LookupEnv(envVarPortalURL); ok {
			if err := apply(config, WithPortalURL(portalURL)); err != nil {
				return fmt.Errorf("failed to parse environment variable %s: %w", envVarPortalURL, err)
			}
		}

		// Read API address from environment
		if apiAddress, ok := os.LookupEnv(envVarAPIAddress); ok {
			config.apiAddress = apiAddress
		}

		// Read legacy API hostname from environment
		if legacyAPIHostname, ok := os.LookupEnv(envVarAPIHostnameLegacy); ok {
			// Allow https:// prefix even though it's the only scheme we allow
			// as it's more natural to support the URL. Any other scheme we
			// don't strip which will fail validation.
			if strings.HasPrefix(strings.ToLower(legacyAPIHostname), "https://") {
				legacyAPIHostname = legacyAPIHostname[8:]
			}

			config.apiAddress = legacyAPIHostname
		}

		// Read API TLS setting from environment
		if apiTLSSetting, ok := os.LookupEnv(envVarAPITLS); ok {
			apiTLSConfig, err := tlsConfigForSetting(apiTLSSetting)
			if err != nil {
				return fmt.Errorf("failed to configure TLS based on environment variable %s: %w", envVarAPITLS, err)
			}
			config.apiTLSConfig = apiTLSConfig
		}

		// Read Auth TLS setting from environment
		if authTLSSetting, ok := os.LookupEnv(envVarAuthTLS); ok {
			authTLSConfig, err := tlsConfigForSetting(authTLSSetting)
			if err != nil {
				return fmt.Errorf("failed to configure TLS based on environment variable %s: %w", envVarAuthTLS, err)
			}
			config.authTLSConfig = authTLSConfig
		}

		// Read SCADA address from environment
		if scadaAddress, ok := os.LookupEnv(envVarSCADAAddress); ok {
			config.scadaAddress = scadaAddress
		}

		// Read SCADA TLS setting from environment
		if scadaTLSSetting, ok := os.LookupEnv(envVarSCADATLS); ok {
			scadaTLSConfig, err := tlsConfigForSetting(scadaTLSSetting)
			if err != nil {
				return fmt.Errorf("failed to configure TLS based on environment variable %s: %w", envVarSCADATLS, err)
			}
			config.scadaTLSConfig = scadaTLSConfig
		}

		// Read user profile information from the environment, the values will only be
		// used if both fields are provided.
		hcpOrganizationID, hcpOrganizationIDOK := os.LookupEnv(envVarHCPOrganizationID)
		hcpProjectID, hcpProjectIDOK := os.LookupEnv(envVarHCPProjectID)

		if hcpOrganizationIDOK && hcpProjectIDOK {
			userProfile := profile.UserProfile{OrganizationID: hcpOrganizationID, ProjectID: hcpProjectID}
			if err := apply(config, WithProfile(&userProfile)); err != nil {
				return fmt.Errorf("failed to configure profile fields based on environment variables (%s, %s): %w", envVarHCPOrganizationID, envVarHCPProjectID, err)
			}
		}

		return nil
	}
}

func tlsConfigForSetting(setting string) (*tls.Config, error) {
	switch setting {
	case tlsSettingDisabled:
		return nil, nil
	case tlsSettingInsecure:
		return &tls.Config{InsecureSkipVerify: true}, nil
	default:
		return nil, fmt.Errorf("invalid TLS setting value: %q", setting)
	}
}
