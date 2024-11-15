package authentication

import (
	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	absauth "github.com/microsoft/kiota-abstractions-go/authentication"
	kiotaidentity "github.com/microsoft/kiota-authentication-azure-go"
)

// AzureIdentityAuthenticationProvider is a wrapper around the AzureIdentityAuthenticationProvider that sets default values for Microsoft Graph.
type AzureIdentityAuthenticationProvider struct {
	kiotaidentity.AzureIdentityAuthenticationProvider
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
	return NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptions(credential, scopes, validHosts, kiotaidentity.ObservabilityOptions{})
}

// NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptions creates a new instance of the AzureIdentityAuthenticationProvider.
func NewAzureIdentityAuthenticationProviderWithScopesAndValidHostsAndObservabilityOptions(credential azcore.TokenCredential, scopes []string, validHosts []string, observabilityOptions kiotaidentity.ObservabilityOptions) (*AzureIdentityAuthenticationProvider, error) {
	accessTokenProvider, err := NewAzureIdentityAccessTokenProviderWithScopesAndValidHostsAndObservabilityOptions(credential, scopes, validHosts, observabilityOptions)
	if err != nil {
		return nil, err
	}
	baseBearer := absauth.NewBaseBearerTokenAuthenticationProvider(accessTokenProvider)
	result := &AzureIdentityAuthenticationProvider{
		AzureIdentityAuthenticationProvider: kiotaidentity.AzureIdentityAuthenticationProvider{
			BaseBearerTokenAuthenticationProvider: *baseBearer,
		},
	}
	return result, nil
}
