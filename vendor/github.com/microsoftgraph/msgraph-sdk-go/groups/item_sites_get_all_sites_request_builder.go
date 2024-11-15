package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSitesGetAllSitesRequestBuilder provides operations to call the getAllSites method.
type ItemSitesGetAllSitesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSitesGetAllSitesRequestBuilderGetQueryParameters list sites across geographies in an organization. This API can also be used to enumerate all sites in a non-multi-geo tenant. For more information, see Best practices for discovering files and detecting changes at scale.
type ItemSitesGetAllSitesRequestBuilderGetQueryParameters struct {
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
// ItemSitesGetAllSitesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesGetAllSitesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSitesGetAllSitesRequestBuilderGetQueryParameters
}
// NewItemSitesGetAllSitesRequestBuilderInternal instantiates a new ItemSitesGetAllSitesRequestBuilder and sets the default values.
func NewItemSitesGetAllSitesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesGetAllSitesRequestBuilder) {
    m := &ItemSitesGetAllSitesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/sites/getAllSites(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSitesGetAllSitesRequestBuilder instantiates a new ItemSitesGetAllSitesRequestBuilder and sets the default values.
func NewItemSitesGetAllSitesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesGetAllSitesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSitesGetAllSitesRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list sites across geographies in an organization. This API can also be used to enumerate all sites in a non-multi-geo tenant. For more information, see Best practices for discovering files and detecting changes at scale.
// Deprecated: This method is obsolete. Use GetAsGetAllSitesGetResponse instead.
// returns a ItemSitesGetAllSitesResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/site-getallsites?view=graph-rest-1.0
func (m *ItemSitesGetAllSitesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSitesGetAllSitesRequestBuilderGetRequestConfiguration)(ItemSitesGetAllSitesResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSitesGetAllSitesResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSitesGetAllSitesResponseable), nil
}
// GetAsGetAllSitesGetResponse list sites across geographies in an organization. This API can also be used to enumerate all sites in a non-multi-geo tenant. For more information, see Best practices for discovering files and detecting changes at scale.
// returns a ItemSitesGetAllSitesGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/site-getallsites?view=graph-rest-1.0
func (m *ItemSitesGetAllSitesRequestBuilder) GetAsGetAllSitesGetResponse(ctx context.Context, requestConfiguration *ItemSitesGetAllSitesRequestBuilderGetRequestConfiguration)(ItemSitesGetAllSitesGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSitesGetAllSitesGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSitesGetAllSitesGetResponseable), nil
}
// ToGetRequestInformation list sites across geographies in an organization. This API can also be used to enumerate all sites in a non-multi-geo tenant. For more information, see Best practices for discovering files and detecting changes at scale.
// returns a *RequestInformation when successful
func (m *ItemSitesGetAllSitesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSitesGetAllSitesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSitesGetAllSitesRequestBuilder when successful
func (m *ItemSitesGetAllSitesRequestBuilder) WithUrl(rawUrl string)(*ItemSitesGetAllSitesRequestBuilder) {
    return NewItemSitesGetAllSitesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
