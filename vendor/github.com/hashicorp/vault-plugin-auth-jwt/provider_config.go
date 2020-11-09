package jwtauth

import (
	"fmt"
)

// Provider-specific configuration interfaces
// All providers must implement the CustomProvider interface, and may implement
// others as needed.

// ProviderMap returns a map of provider names to custom types
func ProviderMap() map[string]CustomProvider {
	return map[string]CustomProvider{
		"azure":  &AzureProvider{},
		"gsuite": &GSuiteProvider{},
	}
}

// CustomProvider - Any custom provider must implement this interface
type CustomProvider interface {
	// Initialize should validate jwtConfig.ProviderConfig, set internal values
	// and run any initialization necessary for subsequent calls to interface
	// functions the provider implements
	Initialize(*jwtConfig) error

	// SensitiveKeys returns any fields in a provider's jwtConfig.ProviderConfig
	// that should be masked or omitted when output
	SensitiveKeys() []string
}

// NewProviderConfig - returns appropriate provider struct if provider_config is
// specified in jwtConfig. The provider map is provider name -to- instance of a
// CustomProvider.
func NewProviderConfig(jc *jwtConfig, providerMap map[string]CustomProvider) (CustomProvider, error) {
	if len(jc.ProviderConfig) == 0 {
		return nil, nil
	}
	provider, ok := jc.ProviderConfig["provider"].(string)
	if !ok {
		return nil, fmt.Errorf("'provider' field not found in provider_config")
	}
	newCustomProvider, ok := providerMap[provider]
	if !ok {
		return nil, fmt.Errorf("provider %q not found in custom providers", provider)
	}
	if err := newCustomProvider.Initialize(jc); err != nil {
		return nil, fmt.Errorf("error initializing %q provider_config: %s", provider, err)
	}
	return newCustomProvider, nil
}

// UserInfoFetcher - Optional support for custom user info handling
type UserInfoFetcher interface {
	FetchUserInfo(*jwtAuthBackend, map[string]interface{}, *jwtRole) error
}

// GroupsFetcher - Optional support for custom groups handling
type GroupsFetcher interface {
	// FetchGroups queries for groups claims during login
	FetchGroups(*jwtAuthBackend, map[string]interface{}, *jwtRole) (interface{}, error)
}
