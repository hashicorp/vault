package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder provides operations to call the validateAuthenticationConfiguration method.
type CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderInternal instantiates a new CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder and sets the default values.
func NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder) {
    m := &CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/customAuthenticationExtensions/validateAuthenticationConfiguration", pathParameters),
    }
    return m
}
// NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder instantiates a new CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder and sets the default values.
func NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderInternal(urlParams, requestAdapter)
}
// Post invoke action validateAuthenticationConfiguration
// returns a AuthenticationConfigurationValidationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder) Post(ctx context.Context, body CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBodyable, requestConfiguration *CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationConfigurationValidationable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAuthenticationConfigurationValidationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AuthenticationConfigurationValidationable), nil
}
// ToPostRequestInformation invoke action validateAuthenticationConfiguration
// returns a *RequestInformation when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder) ToPostRequestInformation(ctx context.Context, body CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBodyable, requestConfiguration *CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder) WithUrl(rawUrl string)(*CustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder) {
    return NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
