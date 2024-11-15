package authentication

import (
	"context"
	"errors"
	"strings"

	abs "github.com/microsoft/kiota-abstractions-go"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// ApiKeyAuthenticationProvider implements the AuthenticationProvider interface and adds an API key to the request.
type ApiKeyAuthenticationProvider struct {
	apiKey        string
	parameterName string
	keyLocation   KeyLocation
	validator     *AllowedHostsValidator
}

type KeyLocation int

const (
	// QUERYPARAMETER_KEYLOCATION is the value for the key location to be used as a query parameter.
	QUERYPARAMETER_KEYLOCATION KeyLocation = iota
	// HEADER_KEYLOCATION is the value for the key location to be used as a header.
	HEADER_KEYLOCATION
)

// NewApiKeyAuthenticationProvider creates a new ApiKeyAuthenticationProvider instance
func NewApiKeyAuthenticationProvider(apiKey string, parameterName string, keyLocation KeyLocation) (*ApiKeyAuthenticationProvider, error) {
	return NewApiKeyAuthenticationProviderWithValidHosts(apiKey, parameterName, keyLocation, nil)
}

// NewApiKeyAuthenticationProviderWithValidHosts creates a new ApiKeyAuthenticationProvider instance while specifying a list of valid hosts
func NewApiKeyAuthenticationProviderWithValidHosts(apiKey string, parameterName string, keyLocation KeyLocation, validHosts []string) (*ApiKeyAuthenticationProvider, error) {
	if len(apiKey) == 0 {
		return nil, errors.New("apiKey cannot be empty")
	}
	if len(parameterName) == 0 {
		return nil, errors.New("parameterName cannot be empty")
	}

	validator, err := NewAllowedHostsValidatorErrorCheck(validHosts)
	if err != nil {
		return nil, err
	}
	return &ApiKeyAuthenticationProvider{
		apiKey:        apiKey,
		parameterName: parameterName,
		keyLocation:   keyLocation,
		validator:     validator,
	}, nil
}

// AuthenticateRequest adds the API key to the request.
func (p *ApiKeyAuthenticationProvider) AuthenticateRequest(ctx context.Context, request *abs.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	ctx, span := otel.GetTracerProvider().Tracer("github.com/microsoft/kiota-abstractions-go").Start(ctx, "GetAuthorizationToken")
	defer span.End()
	if request == nil {
		return errors.New("request cannot be nil")
	}

	url, err := request.GetUri()

	if err != nil {
		return err
	}

	if !(*(p.validator)).IsUrlHostValid(url) {
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.authentication.is_url_valid", false))
		return nil
	}
	if !strings.EqualFold(url.Scheme, "https") {
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.authentication.is_url_valid", false))
		err := errors.New("url scheme must be https")
		span.RecordError(err)
		return err
	}
	span.SetAttributes(attribute.Bool("com.microsoft.kiota.authentication.is_url_valid", true))

	switch p.keyLocation {
	case QUERYPARAMETER_KEYLOCATION:
		query := url.Query()
		query.Set(p.parameterName, p.apiKey)
		url.RawQuery = query.Encode()
		request.SetUri(*url)
	case HEADER_KEYLOCATION:
		request.Headers.Add(p.parameterName, p.apiKey)
	}

	return nil
}
