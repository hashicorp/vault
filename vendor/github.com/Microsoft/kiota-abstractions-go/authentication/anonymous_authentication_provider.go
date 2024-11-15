package authentication

import (
	"context"
	abs "github.com/microsoft/kiota-abstractions-go"
)

// AnonymousAuthenticationProvider implements the AuthenticationProvider interface does not perform any authentication.
type AnonymousAuthenticationProvider struct {
}

// AuthenticateRequest is a placeholder method that "authenticates" the RequestInformation instance: no-op.
func (provider *AnonymousAuthenticationProvider) AuthenticateRequest(context context.Context, request *abs.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	return nil
}
