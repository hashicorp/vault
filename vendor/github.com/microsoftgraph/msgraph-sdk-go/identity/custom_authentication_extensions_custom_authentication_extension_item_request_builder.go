package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder provides operations to manage the customAuthenticationExtensions property of the microsoft.graph.identityContainer entity.
type CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetQueryParameters read the properties and relationships of a customAuthenticationExtension object. The following derived types are currently supported.
type CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetQueryParameters
}
// CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderInternal instantiates a new CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder and sets the default values.
func NewCustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) {
    m := &CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/customAuthenticationExtensions/{customAuthenticationExtension%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewCustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder instantiates a new CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder and sets the default values.
func NewCustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete a customAuthenticationExtension object. The following derived types are currently supported.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/customauthenticationextension-delete?view=graph-rest-1.0
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get read the properties and relationships of a customAuthenticationExtension object. The following derived types are currently supported.
// returns a CustomAuthenticationExtensionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/customauthenticationextension-get?view=graph-rest-1.0
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) Get(ctx context.Context, requestConfiguration *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomAuthenticationExtensionable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomAuthenticationExtensionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomAuthenticationExtensionable), nil
}
// Patch update the properties of a customAuthenticationExtension object. The following derived types are currently supported.
// returns a CustomAuthenticationExtensionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/customauthenticationextension-update?view=graph-rest-1.0
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomAuthenticationExtensionable, requestConfiguration *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomAuthenticationExtensionable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomAuthenticationExtensionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomAuthenticationExtensionable), nil
}
// ToDeleteRequestInformation delete a customAuthenticationExtension object. The following derived types are currently supported.
// returns a *RequestInformation when successful
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of a customAuthenticationExtension object. The following derived types are currently supported.
// returns a *RequestInformation when successful
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPatchRequestInformation update the properties of a customAuthenticationExtension object. The following derived types are currently supported.
// returns a *RequestInformation when successful
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomAuthenticationExtensionable, requestConfiguration *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// ValidateAuthenticationConfiguration provides operations to call the validateAuthenticationConfiguration method.
// returns a *CustomAuthenticationExtensionsItemValidateAuthenticationConfigurationRequestBuilder when successful
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) ValidateAuthenticationConfiguration()(*CustomAuthenticationExtensionsItemValidateAuthenticationConfigurationRequestBuilder) {
    return NewCustomAuthenticationExtensionsItemValidateAuthenticationConfigurationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder when successful
func (m *CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) WithUrl(rawUrl string)(*CustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder) {
    return NewCustomAuthenticationExtensionsCustomAuthenticationExtensionItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
