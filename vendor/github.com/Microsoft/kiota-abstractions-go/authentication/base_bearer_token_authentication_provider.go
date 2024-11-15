package authentication

import (
	"context"
	"errors"

	abs "github.com/microsoft/kiota-abstractions-go"
)

const authorizationHeader = "Authorization"
const claimsKey = "claims"

// BaseBearerTokenAuthenticationProvider provides a base class implementing AuthenticationProvider for Bearer token scheme.
type BaseBearerTokenAuthenticationProvider struct {
	// accessTokenProvider is called by the BaseBearerTokenAuthenticationProvider class to authenticate the request via the returned access token.
	accessTokenProvider AccessTokenProvider
}

// NewBaseBearerTokenAuthenticationProvider creates a new instance of the BaseBearerTokenAuthenticationProvider class.
func NewBaseBearerTokenAuthenticationProvider(accessTokenProvider AccessTokenProvider) *BaseBearerTokenAuthenticationProvider {
	return &BaseBearerTokenAuthenticationProvider{accessTokenProvider}
}

// AuthenticateRequest authenticates the provided RequestInformation instance using the provided authorization token callback.
func (provider *BaseBearerTokenAuthenticationProvider) AuthenticateRequest(ctx context.Context, request *abs.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	if request == nil {
		return errors.New("request is nil")
	}
	if request.Headers == nil {
		request.Headers = abs.NewRequestHeaders()
	}
	if provider.accessTokenProvider == nil {
		return errors.New("this class needs to be initialized with an access token provider")
	}
	if len(additionalAuthenticationContext) > 0 &&
		additionalAuthenticationContext[claimsKey] != nil &&
		request.Headers.ContainsKey(authorizationHeader) {
		request.Headers.Remove(authorizationHeader)
	}
	if !request.Headers.ContainsKey(authorizationHeader) {
		uri, err := request.GetUri()
		if err != nil {
			return err
		}
		token, err := provider.accessTokenProvider.GetAuthorizationToken(ctx, uri, additionalAuthenticationContext)
		if err != nil {
			return err
		}
		if token != "" {
			request.Headers.Add(authorizationHeader, "Bearer "+token)
		}
	}

	return nil
}

// GetAuthorizationTokenProvider returns the access token provider the BaseBearerTokenAuthenticationProvider class uses to authenticate the request.
func (provider *BaseBearerTokenAuthenticationProvider) GetAuthorizationTokenProvider() AccessTokenProvider {
	return provider.accessTokenProvider
}
