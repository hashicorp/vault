package authentication

import (
	"context"
	abs "github.com/microsoft/kiota-abstractions-go"
)

// AuthenticationProvider authenticates the RequestInformation request.
type AuthenticationProvider interface {
	// AuthenticateRequest authenticates the provided RequestInformation.
	AuthenticateRequest(context context.Context, request *abs.RequestInformation, additionalAuthenticationContext map[string]interface{}) error
}
