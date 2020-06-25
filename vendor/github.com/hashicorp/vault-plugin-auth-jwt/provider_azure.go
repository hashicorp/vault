package jwtauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	log "github.com/hashicorp/go-hclog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// The old MS graph API requires setting an api-version query parameter
	windowsGraphHost  = "graph.windows.net"
	windowsAPIVersion = "1.6"

	// Distributed claim fields
	claimNamesField   = "_claim_names"
	claimSourcesField = "_claim_sources"
)

// AzureProvider is used for Azure-specific configuration
type AzureProvider struct {
	// Context for azure calls
	ctx context.Context

	// OIDC provider
	provider *oidc.Provider
}

// Initialize anything in the AzureProvider struct - satisfying the CustomProvider interface
func (a *AzureProvider) Initialize(jc *jwtConfig) error {
	return nil
}

// SensitiveKeys - satisfying the CustomProvider interface
func (a *AzureProvider) SensitiveKeys() []string {
	return []string{}
}

// FetchGroups - custom groups fetching for azure - satisfying GroupsFetcher interface
func (a *AzureProvider) FetchGroups(b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole) (interface{}, error) {
	groupsClaimRaw := getClaim(b.Logger(), allClaims, role.GroupsClaim)

	if groupsClaimRaw == nil {
		// If the "groups" claim is missing, it might be because the user is a
		// member of more than 200 groups, which means the token contains
		// distributed claim information. Attempt to look that up here.
		azureClaimSourcesURL, err := a.getClaimSource(b.Logger(), allClaims, role)
		if err != nil {
			return nil, fmt.Errorf("unable to get claim sources: %s", err)
		}

		// Get provider because we'll need to get a new token for microsoft's
		// graph API, specifically the old graph API
		provider, err := b.getProvider(b.cachedConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to get provider: %s", err)
		}
		a.provider = provider

		a.ctx, err = b.createCAContext(b.providerCtx, b.cachedConfig.OIDCDiscoveryCAPEM)
		if err != nil {
			return nil, fmt.Errorf("unable to create CA Context: %s", err)
		}

		azureGroups, err := a.getAzureGroups(azureClaimSourcesURL, b.cachedConfig)
		if err != nil {
			return nil, fmt.Errorf("%q claim not found in token: %v", role.GroupsClaim, err)
		}
		groupsClaimRaw = azureGroups
	}
	b.Logger().Debug(fmt.Sprintf("groups claim raw is %v", groupsClaimRaw))
	return groupsClaimRaw, nil
}

// In Azure, if you are indirectly member of more than 200 groups, they will
// send _claim_names and _claim_sources instead of the groups, per OIDC Core
// 1.0, section 5.6.2:
// https://openid.net/specs/openid-connect-core-1_0.html#AggregatedDistributedClaims
// In the future this could be used with other providers as well. Example:
//
// {
// 	 "_claim_names": {
// 	   "groups": "src1"
// 	 },
// 	 "_claim_sources": {
// 	   "src1": {
// 	     "endpoint": "https://graph.windows.net...."
// 	   }
//   }
// }
//
// For this to work, "profile" should be set in "oidc_scopes" in the vault oidc role.
//
func (a *AzureProvider) getClaimSource(logger log.Logger, allClaims map[string]interface{}, role *jwtRole) (string, error) {
	// Get the source key for the groups claim
	name := fmt.Sprintf("/%s/%s", claimNamesField, role.GroupsClaim)
	groupsClaimSource := getClaim(logger, allClaims, name)
	if groupsClaimSource == nil {
		return "", fmt.Errorf("unable to locate groups claim %q in %s", role.GroupsClaim, claimNamesField)
	}
	// Get the endpoint source for the groups claim
	endpoint := fmt.Sprintf("/%s/%s/endpoint", claimSourcesField, groupsClaimSource.(string))
	val := getClaim(logger, allClaims, endpoint)
	if val == nil {
		return "", fmt.Errorf("unable to locate %s in claims", endpoint)
	}
	logger.Debug(fmt.Sprintf("found Azure Graph API endpoint for group membership: %v", val))
	return fmt.Sprintf("%v", val), nil
}

// Fetch user groups from the Azure AD Graph API
func (a *AzureProvider) getAzureGroups(groupsURL string, c *jwtConfig) (interface{}, error) {
	urlParsed, err := url.Parse(groupsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse distributed groups source url %s: %s", groupsURL, err)
	}
	token, err := a.getAzureToken(c, urlParsed.Host)
	if err != nil {
		return nil, fmt.Errorf("unable to get token: %s", err)
	}
	payload := strings.NewReader("{\"securityEnabledOnly\": false}")
	req, err := http.NewRequest("POST", groupsURL, payload)
	if err != nil {
		return nil, fmt.Errorf("error constructing groups endpoint request: %s", err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	// If endpoint is the old windows graph api, add api-version
	if urlParsed.Host == windowsGraphHost {
		query := req.URL.Query()
		query.Add("api-version", windowsAPIVersion)
		req.URL.RawQuery = query.Encode()
	}
	client := http.DefaultClient
	if c, ok := a.ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = c
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to call Azure AD Graph API: %s", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Azure AD Graph API response: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get groups: %s", string(body))
	}

	var target azureGroups
	if err := json.Unmarshal(body, &target); err != nil {
		return nil, fmt.Errorf("unabled to decode response: %s", err)
	}
	return target.Value, nil
}

// Login to Azure, using client id and secret.
func (a *AzureProvider) getAzureToken(c *jwtConfig, host string) (string, error) {
	config := &clientcredentials.Config{
		ClientID:     c.OIDCClientID,
		ClientSecret: c.OIDCClientSecret,
		TokenURL:     a.provider.Endpoint().TokenURL,
		Scopes: []string{
			"openid",
			"profile",
			"https://" + host + "/.default",
		},
	}
	token, err := config.Token(a.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Azure token: %s", err)
	}
	return token.AccessToken, nil
}

type azureGroups struct {
	Value []interface{} `json:"value"`
}
