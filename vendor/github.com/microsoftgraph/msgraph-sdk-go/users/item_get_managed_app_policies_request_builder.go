package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemGetManagedAppPoliciesRequestBuilder provides operations to call the getManagedAppPolicies method.
type ItemGetManagedAppPoliciesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemGetManagedAppPoliciesRequestBuilderGetQueryParameters gets app restrictions for a given user.
type ItemGetManagedAppPoliciesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// ItemGetManagedAppPoliciesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemGetManagedAppPoliciesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemGetManagedAppPoliciesRequestBuilderGetQueryParameters
}
// NewItemGetManagedAppPoliciesRequestBuilderInternal instantiates a new ItemGetManagedAppPoliciesRequestBuilder and sets the default values.
func NewItemGetManagedAppPoliciesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGetManagedAppPoliciesRequestBuilder) {
    m := &ItemGetManagedAppPoliciesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/getManagedAppPolicies(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemGetManagedAppPoliciesRequestBuilder instantiates a new ItemGetManagedAppPoliciesRequestBuilder and sets the default values.
func NewItemGetManagedAppPoliciesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemGetManagedAppPoliciesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemGetManagedAppPoliciesRequestBuilderInternal(urlParams, requestAdapter)
}
// Get gets app restrictions for a given user.
// Deprecated: This method is obsolete. Use GetAsGetManagedAppPoliciesGetResponse instead.
// returns a ItemGetManagedAppPoliciesResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-mam-user-getmanagedapppolicies?view=graph-rest-1.0
func (m *ItemGetManagedAppPoliciesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemGetManagedAppPoliciesRequestBuilderGetRequestConfiguration)(ItemGetManagedAppPoliciesResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemGetManagedAppPoliciesResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemGetManagedAppPoliciesResponseable), nil
}
// GetAsGetManagedAppPoliciesGetResponse gets app restrictions for a given user.
// returns a ItemGetManagedAppPoliciesGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-mam-user-getmanagedapppolicies?view=graph-rest-1.0
func (m *ItemGetManagedAppPoliciesRequestBuilder) GetAsGetManagedAppPoliciesGetResponse(ctx context.Context, requestConfiguration *ItemGetManagedAppPoliciesRequestBuilderGetRequestConfiguration)(ItemGetManagedAppPoliciesGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemGetManagedAppPoliciesGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemGetManagedAppPoliciesGetResponseable), nil
}
// ToGetRequestInformation gets app restrictions for a given user.
// returns a *RequestInformation when successful
func (m *ItemGetManagedAppPoliciesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemGetManagedAppPoliciesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemGetManagedAppPoliciesRequestBuilder when successful
func (m *ItemGetManagedAppPoliciesRequestBuilder) WithUrl(rawUrl string)(*ItemGetManagedAppPoliciesRequestBuilder) {
    return NewItemGetManagedAppPoliciesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
