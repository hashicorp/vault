package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// FederationConfigurationsAvailableProviderTypesRequestBuilder provides operations to call the availableProviderTypes method.
type FederationConfigurationsAvailableProviderTypesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// FederationConfigurationsAvailableProviderTypesRequestBuilderGetQueryParameters get all identity providers supported in a directory.
type FederationConfigurationsAvailableProviderTypesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// FederationConfigurationsAvailableProviderTypesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type FederationConfigurationsAvailableProviderTypesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *FederationConfigurationsAvailableProviderTypesRequestBuilderGetQueryParameters
}
// NewFederationConfigurationsAvailableProviderTypesRequestBuilderInternal instantiates a new FederationConfigurationsAvailableProviderTypesRequestBuilder and sets the default values.
func NewFederationConfigurationsAvailableProviderTypesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*FederationConfigurationsAvailableProviderTypesRequestBuilder) {
    m := &FederationConfigurationsAvailableProviderTypesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/federationConfigurations/availableProviderTypes(){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewFederationConfigurationsAvailableProviderTypesRequestBuilder instantiates a new FederationConfigurationsAvailableProviderTypesRequestBuilder and sets the default values.
func NewFederationConfigurationsAvailableProviderTypesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*FederationConfigurationsAvailableProviderTypesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewFederationConfigurationsAvailableProviderTypesRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get all identity providers supported in a directory.
// Deprecated: This method is obsolete. Use GetAsAvailableProviderTypesGetResponse instead.
// returns a FederationConfigurationsAvailableProviderTypesResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/identityproviderbase-availableprovidertypes?view=graph-rest-1.0
func (m *FederationConfigurationsAvailableProviderTypesRequestBuilder) Get(ctx context.Context, requestConfiguration *FederationConfigurationsAvailableProviderTypesRequestBuilderGetRequestConfiguration)(FederationConfigurationsAvailableProviderTypesResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateFederationConfigurationsAvailableProviderTypesResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(FederationConfigurationsAvailableProviderTypesResponseable), nil
}
// GetAsAvailableProviderTypesGetResponse get all identity providers supported in a directory.
// returns a FederationConfigurationsAvailableProviderTypesGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/identityproviderbase-availableprovidertypes?view=graph-rest-1.0
func (m *FederationConfigurationsAvailableProviderTypesRequestBuilder) GetAsAvailableProviderTypesGetResponse(ctx context.Context, requestConfiguration *FederationConfigurationsAvailableProviderTypesRequestBuilderGetRequestConfiguration)(FederationConfigurationsAvailableProviderTypesGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateFederationConfigurationsAvailableProviderTypesGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(FederationConfigurationsAvailableProviderTypesGetResponseable), nil
}
// ToGetRequestInformation get all identity providers supported in a directory.
// returns a *RequestInformation when successful
func (m *FederationConfigurationsAvailableProviderTypesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *FederationConfigurationsAvailableProviderTypesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *FederationConfigurationsAvailableProviderTypesRequestBuilder when successful
func (m *FederationConfigurationsAvailableProviderTypesRequestBuilder) WithUrl(rawUrl string)(*FederationConfigurationsAvailableProviderTypesRequestBuilder) {
    return NewFederationConfigurationsAvailableProviderTypesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
