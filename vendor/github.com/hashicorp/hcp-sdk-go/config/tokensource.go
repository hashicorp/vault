// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"context"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/hcp-sdk-go/auth"
	"github.com/hashicorp/hcp-sdk-go/auth/tokencache"
	"github.com/hashicorp/hcp-sdk-go/auth/workload"
	"github.com/hashicorp/hcp-sdk-go/config/files"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type sourceType = string

var (
	sourceTypeLogin            = sourceType("login")
	sourceTypeServicePrincipal = sourceType("service-principal")
	sourceTypeWorkload         = sourceType("workload")
)

func (c *hcpConfig) setTokenSource() error {
	// Check if a custom token source has been provided
	if c.tokenSource != nil {
		return nil
	}

	// Get the credential cache path
	// TODO: make this configurable
	cacheFile, err := files.TokenCacheFile()
	if err != nil {
		return err
	}

	tokenSource, sourceType, sourceIdentifier, err := c.getTokenSource()
	if err != nil {
		return err
	}

	switch sourceType {
	case sourceTypeLogin:
		c.tokenSource = tokencache.NewLoginTokenSource(cacheFile, tokenSource, &c.oauth2Config)
	case sourceTypeServicePrincipal:
		c.tokenSource = tokencache.NewServicePrincipalTokenSource(
			cacheFile,
			sourceIdentifier,
			tokenSource,
		)
	case sourceTypeWorkload:
		c.tokenSource = tokencache.NewWorkloadTokenSource(
			cacheFile,
			sourceIdentifier,
			tokenSource,
		)
	}

	return nil
}

// getTokenSource gets the token source. The order of precedence is:
//
// 1. Configured client credentials (either explicit or through environment
// variables).
// 2. Via credential file (sourced first via environment variable and then
// default file location).
// 3. Interactive session.
func (c *hcpConfig) getTokenSource() (oauth2.TokenSource, sourceType, string, error) {
	// Set up a token context with the custom auth TLS config
	tokenTransport := cleanhttp.DefaultPooledTransport()
	tokenTransport.TLSClientConfig = c.authTLSConfig
	ctx := context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		&http.Client{Transport: tokenTransport},
	)

	// Set client credentials token URL based on auth URL.
	tokenURL := c.authURL
	tokenURL.Path = tokenPath

	clientCredentials := clientcredentials.Config{
		EndpointParams: url.Values{"audience": {aud}},
		TokenURL:       tokenURL.String(),
	}

	// Set access token via configured client credentials.
	if c.clientID != "" && c.clientSecret != "" {
		// Create token source for client secrets
		clientCredentials.ClientID = c.clientID
		clientCredentials.ClientSecret = c.clientSecret

		return clientCredentials.TokenSource(ctx), sourceTypeServicePrincipal, clientCredentials.ClientID, nil
	}

	// Use workload provider config if it was provided
	if c.workloadProviderConfig != nil {
		provider, err := workload.New(c.workloadProviderConfig)
		if err != nil {
			return nil, "", "", err
		}
		provider.SetAPI(c)
		return oauth2.ReuseTokenSource(nil, provider), sourceTypeWorkload, c.workloadProviderConfig.ProviderResourceName, nil
	}

	// If we haven't been given an explicit credential file to use, try to load
	// the credential file from the environment or default location.
	if c.credentialFile == nil {
		credFile, err := auth.GetDefaultCredentialFile()
		if err != nil {
			return nil, "", "", err
		}
		c.credentialFile = credFile
	}

	// If we found a credential file use it as a credential source
	if c.credentialFile != nil {
		if c.credentialFile.Scheme == auth.CredentialFileSchemeServicePrincipal {
			// Set credentials on client credentials configuration
			clientCredentials.ClientID = c.credentialFile.Oauth.ClientID
			clientCredentials.ClientSecret = c.credentialFile.Oauth.ClientSecret

			// Create token source from the client credentials configuration.
			return clientCredentials.TokenSource(ctx), sourceTypeServicePrincipal, clientCredentials.ClientID, nil
		} else if c.credentialFile.Scheme == auth.CredentialFileSchemeWorkload {
			w, err := workload.New(c.credentialFile.Workload)
			if err != nil {
				return nil, "", "", err
			}

			// Set the API info
			w.SetAPI(c)

			// Use the workload provider as the token source
			return w, sourceTypeWorkload, c.credentialFile.Workload.ProviderResourceName, nil
		}
	}

	var loginTokenSource oauth2.TokenSource
	if !c.noBrowserLogin {
		loginTokenSource = auth.NewBrowserLogin(&c.oauth2Config, c.noDefaultBrowser)
	}

	return loginTokenSource, sourceTypeLogin, "", nil
}
