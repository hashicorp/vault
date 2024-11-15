package identity

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder provides operations to manage the identityProviders property of the microsoft.graph.b2xIdentityUserFlow entity.
type B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetQueryParameters the identity providers included in the user flow.
type B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetQueryParameters
}
// NewB2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderInternal instantiates a new B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder and sets the default values.
func NewB2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder) {
    m := &B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identity/b2xUserFlows/{b2xIdentityUserFlow%2Did}/identityProviders/{identityProvider%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewB2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder instantiates a new B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder and sets the default values.
func NewB2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewB2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the identity providers included in the user flow.
// Deprecated: The identityProvider API is deprecated and will stop returning data on March 2023. Please use the new identityProviderBase API. as of 2021-05/identityProvider
// returns a IdentityProviderable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder) Get(ctx context.Context, requestConfiguration *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProviderable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentityProviderFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentityProviderable), nil
}
// ToGetRequestInformation the identity providers included in the user flow.
// Deprecated: The identityProvider API is deprecated and will stop returning data on March 2023. Please use the new identityProviderBase API. as of 2021-05/identityProvider
// returns a *RequestInformation when successful
func (m *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Deprecated: The identityProvider API is deprecated and will stop returning data on March 2023. Please use the new identityProviderBase API. as of 2021-05/identityProvider
// returns a *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder when successful
func (m *B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder) WithUrl(rawUrl string)(*B2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder) {
    return NewB2xUserFlowsItemIdentityProvidersIdentityProviderItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
