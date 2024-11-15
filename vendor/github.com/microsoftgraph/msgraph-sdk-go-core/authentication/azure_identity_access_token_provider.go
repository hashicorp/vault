package authentication

import (
	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	kiotaidentity "github.com/microsoft/kiota-authentication-azure-go"
)

// AzureIdentityAccessTokenProvider is a wrapper around the AzureIdentityAccessTokenProvider from the Kiota library with Microsoft Graph default valid hosts.
type AzureIdentityAccessTokenProvider struct {
	kiotaidentity.AzureIdentityAccessTokenProvider
}

// NewAzureIdentityAccessTokenProvider creates a new instance of the AzureIdentityAccessTokenProvider using "<scheme>://<host>/.default" as the default scope.
func NewAzureIdentityAccessTokenProvider(credential azcore.TokenCredential) (*AzureIdentityAccessTokenProvider, error) {
	return NewAzureIdentityAccessTokenProviderWithScopes(credential, nil)
}

// NewAzureIdentityAccessTokenProviderWithScopes creates a new instance of the AzureIdentityAccessTokenProvider.
func NewAzureIdentityAccessTokenProviderWithScopes(credential azcore.TokenCredential, scopes []string) (*AzureIdentityAccessTokenProvider, error) {
	return NewAzureIdentityAccessTokenProviderWithScopesAndValidHosts(credential, scopes, nil)
}

// NewAzureIdentityAccessTokenProviderWithScopesAndValidHosts creates a new instance of the AzureIdentityAccessTokenProvider.
func NewAzureIdentityAccessTokenProviderWithScopesAndValidHosts(credential azcore.TokenCredential, scopes []string, validHosts []string) (*AzureIdentityAccessTokenProvider, error) {
	return NewAzureIdentityAccessTokenProviderWithScopesAndValidHostsAndObservabilityOptions(credential, scopes, validHosts, kiotaidentity.ObservabilityOptions{})
}

// NewAzureIdentityAccessTokenProviderWithScopesAndValidHosts creates a new instance of the AzureIdentityAccessTokenProvider.
func NewAzureIdentityAccessTokenProviderWithScopesAndValidHostsAndObservabilityOptions(credential azcore.TokenCredential, scopes []string, validHosts []string, observabilityOptions kiotaidentity.ObservabilityOptions) (*AzureIdentityAccessTokenProvider, error) {
	base, err := kiotaidentity.NewAzureIdentityAccessTokenProviderWithScopesAndValidHostsAndObservabilityOptions(credential, scopes, validHosts, observabilityOptions)
	if err != nil {
		return nil, err
	}
	if len(validHosts) == 0 {
		base.GetAllowedHostsValidator().SetAllowedHosts([]string{"graph.microsoft.com", "graph.microsoft.us", "dod-graph.microsoft.us", "graph.microsoft.de", "microsoftgraph.chinacloudapi.cn", "canary.graph.microsoft.com"})
	}
	result := &AzureIdentityAccessTokenProvider{
		AzureIdentityAccessTokenProvider: *base,
	}

	return result, nil
}
