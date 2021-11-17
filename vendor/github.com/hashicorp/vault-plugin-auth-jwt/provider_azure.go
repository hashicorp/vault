package jwtauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"golang.org/x/oauth2"
)

const (
	// Deprecated: The host of the Azure Active Directory (AAD) graph API
	azureADGraphHost = "graph.windows.net"

	// The host and version of the Microsoft Graph API
	microsoftGraphHost       = "graph.microsoft.com"
	microsoftGraphAPIVersion = "/v1.0"

	// Distributed claim fields
	claimNamesField   = "_claim_names"
	claimSourcesField = "_claim_sources"
)

// AzureProvider is used for Azure-specific configuration
type AzureProvider struct {
	// Context for azure calls
	ctx context.Context
}

// Initialize anything in the AzureProvider struct - satisfying the CustomProvider interface
func (a *AzureProvider) Initialize(_ context.Context, _ *jwtConfig) error {
	return nil
}

// SensitiveKeys - satisfying the CustomProvider interface
func (a *AzureProvider) SensitiveKeys() []string {
	return []string{}
}

// FetchGroups - custom groups fetching for azure - satisfying GroupsFetcher interface
func (a *AzureProvider) FetchGroups(_ context.Context, b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole, tokenSource oauth2.TokenSource) (interface{}, error) {
	groupsClaimRaw := getClaim(b.Logger(), allClaims, role.GroupsClaim)

	if groupsClaimRaw == nil {
		// If the "groups" claim is missing, it might be because the user is a
		// member of more than 200 groups, which means the token contains
		// distributed claim information. Attempt to look that up here.
		azureClaimSourcesURL, err := a.getClaimSource(b.Logger(), allClaims, role)
		if err != nil {
			return nil, fmt.Errorf("unable to get claim sources: %s", err)
		}

		a.ctx, err = b.createCAContext(b.providerCtx, b.cachedConfig.OIDCDiscoveryCAPEM)
		if err != nil {
			return nil, fmt.Errorf("unable to create CA Context: %s", err)
		}

		azureGroups, err := a.getAzureGroups(azureClaimSourcesURL, tokenSource)
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

	urlParsed, err := url.Parse(fmt.Sprintf("%v", val))
	if err != nil {
		return "", fmt.Errorf("unable to parse claim source URL: %w", err)
	}

	// If the endpoint source for the groups claim has a host of the deprecated AAD graph API,
	// then replace it to instead use the Microsoft graph API. The AAD graph API is deprecated
	// and will eventually stop servicing requests. See details at:
	// - https://developer.microsoft.com/en-us/office/blogs/microsoft-graph-or-azure-ad-graph/
	// - https://docs.microsoft.com/en-us/graph/api/overview?view=graph-rest-1.0
	if urlParsed.Host == azureADGraphHost {
		urlParsed.Host = microsoftGraphHost
		urlParsed.Path = microsoftGraphAPIVersion + urlParsed.Path
	}

	logger.Debug(fmt.Sprintf("found Azure Graph API endpoint for group membership: %v", urlParsed.String()))
	return urlParsed.String(), nil
}

// Fetch user groups from the Microsoft Graph API
func (a *AzureProvider) getAzureGroups(groupsURL string, tokenSource oauth2.TokenSource) (interface{}, error) {
	urlParsed, err := url.Parse(groupsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse distributed groups source url %s: %s", groupsURL, err)
	}

	// Use the Access Token that was pre-negotiated between the Claims Provider and RP
	// via https://openid.net/specs/openid-connect-core-1_0.html#AggregatedDistributedClaims.
	if tokenSource == nil {
		return nil, errors.New("token unavailable to call Microsoft Graph API")
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to get token: %s", err)
	}
	payload := strings.NewReader("{\"securityEnabledOnly\": false}")
	req, err := http.NewRequest("POST", urlParsed.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("error constructing groups endpoint request: %s", err)
	}
	req.Header.Add("content-type", "application/json")
	token.SetAuthHeader(req)

	client := http.DefaultClient
	if c, ok := a.ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = c
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to call Microsoft Graph API: %s", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Microsoft Graph API response: %s", err)
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

type azureGroups struct {
	Value []interface{} `json:"value"`
}
