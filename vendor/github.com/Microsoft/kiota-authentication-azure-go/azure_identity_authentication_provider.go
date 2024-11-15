// Package microsoft_kiota_authentication_azure implements Kiota abstractions for authentication using the Azure Core library.
// In order to use this package, you must also add the github.com/Azure/azure-sdk-for-go/sdk/azidentity.
package microsoft_kiota_authentication_azure

import (
	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	auth "github.com/microsoft/kiota-abstractions-go/authentication"
)

// AzureIdentityAuthenticationProvider implementation of AuthenticationProvider that supports implementations of TokenCredential from Azure.Identity.
type AzureIdentityAuthenticationProvider struct {
	auth.BaseBearerTokenAuthenticationProvider
}

// NewAzureIdentityAuthenticationProvider creates a new instance of the AzureIdentityAuthenticationProvider using "https://graph.microsoft.com/.default" as the default scope.
func NewAzureIdentityAuthenticationProvider(credential azcore.TokenCredential) (*AzureIdentityAuthenticationProvider, error) {
	return NewAzureIdentityAuthenticationProviderWithScopes(credential, nil)
}

// NewAzureIdentityAuthenticationProviderWithScopes creates a new instance of the AzureIdentityAuthenticationProvider.
func NewAzureIdentityAuthenticationProviderWithScopes(credential azcore.TokenCredential, scopes []string) (*AzureIdentityAuthenticationProvider, error) {
	return NewAzureIdentityAuthenticationProviderWithScopesAndValidHosts(credential, scopes, nil)
}

// NewAzureIdentityAuthenticationProviderWithScopesAndValidHosts creates a new instance of the AzureIdentityAuthenticationProvider.
func NewAzureIdentityAuthenticationProviderWithScopesAndValidHosts(credential azcore.TokenCredential, scopes []string, validHosts []string) (*AzureIdentityAuthenticationProvider, error) {
	return NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptions(credential, scopes, validHosts, ObservabilityOptions{})
}

// NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptions creates a new instance of the AzureIdentityAuthenticationProvider.
func NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptions(credential azcore.TokenCredential, scopes []string, validHosts []string, observabilityOptions ObservabilityOptions) (*AzureIdentityAuthenticationProvider, error) {
	return NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptionsAndIsCaeEnabled(credential, scopes, validHosts, observabilityOptions, true)
}

// NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptionsAndIsCaeEnabled creates a new instance of the AzureIdentityAuthenticationProvider.
func NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptionsAndIsCaeEnabled(credential azcore.TokenCredential, scopes []string, validHosts []string, observabilityOptions ObservabilityOptions, isCaeEnabled bool) (*AzureIdentityAuthenticationProvider, error) {
	accessTokenProvider, err := NewAzureIdentityAccessTokenProviderWithScopesAndValidHostsAndObservabilityOptionsAndIsCaeEnabled(credential, scopes, validHosts, observabilityOptions, isCaeEnabled)
	if err != nil {
		return nil, err
	}
	baseBearer := auth.NewBaseBearerTokenAuthenticationProvider(accessTokenProvider)
	result := &AzureIdentityAuthenticationProvider{
		BaseBearerTokenAuthenticationProvider: *baseBearer,
	}
	return result, nil
}
