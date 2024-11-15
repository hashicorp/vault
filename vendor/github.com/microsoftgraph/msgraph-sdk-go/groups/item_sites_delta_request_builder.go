package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSitesDeltaRequestBuilder provides operations to call the delta method.
type ItemSitesDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSitesDeltaRequestBuilderGetQueryParameters get newly created, updated, or deleted sites without having to perform a full read of the entire sites collection. A delta function call for sites is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls,you can query for incremental changes in the sites. It allows you to maintain and synchronize a local store of a user's sites without having to fetch all the sites from the server every time.The application calls the API without specifying any parameters.The service begins enumerating sites and returns pages of changes to these sites, accompanied by either an @odata.nextLink or an @odata.deltaLink.Your application should continue making calls using the @odata.nextLink until there's an @odata.deltaLink  in the response. After you receive all the changes, you can apply them to your local state.To monitor future changes, call the delta API by using the @odata.deltaLink in the previous response. Any resources marked as deleted should be removed from your local state.
type ItemSitesDeltaRequestBuilderGetQueryParameters struct {
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
// ItemSitesDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSitesDeltaRequestBuilderGetQueryParameters
}
// NewItemSitesDeltaRequestBuilderInternal instantiates a new ItemSitesDeltaRequestBuilder and sets the default values.
func NewItemSitesDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesDeltaRequestBuilder) {
    m := &ItemSitesDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/sites/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSitesDeltaRequestBuilder instantiates a new ItemSitesDeltaRequestBuilder and sets the default values.
func NewItemSitesDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSitesDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get newly created, updated, or deleted sites without having to perform a full read of the entire sites collection. A delta function call for sites is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls,you can query for incremental changes in the sites. It allows you to maintain and synchronize a local store of a user's sites without having to fetch all the sites from the server every time.The application calls the API without specifying any parameters.The service begins enumerating sites and returns pages of changes to these sites, accompanied by either an @odata.nextLink or an @odata.deltaLink.Your application should continue making calls using the @odata.nextLink until there's an @odata.deltaLink  in the response. After you receive all the changes, you can apply them to your local state.To monitor future changes, call the delta API by using the @odata.deltaLink in the previous response. Any resources marked as deleted should be removed from your local state.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ItemSitesDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/site-delta?view=graph-rest-1.0
func (m *ItemSitesDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSitesDeltaRequestBuilderGetRequestConfiguration)(ItemSitesDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSitesDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSitesDeltaResponseable), nil
}
// GetAsDeltaGetResponse get newly created, updated, or deleted sites without having to perform a full read of the entire sites collection. A delta function call for sites is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls,you can query for incremental changes in the sites. It allows you to maintain and synchronize a local store of a user's sites without having to fetch all the sites from the server every time.The application calls the API without specifying any parameters.The service begins enumerating sites and returns pages of changes to these sites, accompanied by either an @odata.nextLink or an @odata.deltaLink.Your application should continue making calls using the @odata.nextLink until there's an @odata.deltaLink  in the response. After you receive all the changes, you can apply them to your local state.To monitor future changes, call the delta API by using the @odata.deltaLink in the previous response. Any resources marked as deleted should be removed from your local state.
// returns a ItemSitesDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/site-delta?view=graph-rest-1.0
func (m *ItemSitesDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ItemSitesDeltaRequestBuilderGetRequestConfiguration)(ItemSitesDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSitesDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSitesDeltaGetResponseable), nil
}
// ToGetRequestInformation get newly created, updated, or deleted sites without having to perform a full read of the entire sites collection. A delta function call for sites is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls,you can query for incremental changes in the sites. It allows you to maintain and synchronize a local store of a user's sites without having to fetch all the sites from the server every time.The application calls the API without specifying any parameters.The service begins enumerating sites and returns pages of changes to these sites, accompanied by either an @odata.nextLink or an @odata.deltaLink.Your application should continue making calls using the @odata.nextLink until there's an @odata.deltaLink  in the response. After you receive all the changes, you can apply them to your local state.To monitor future changes, call the delta API by using the @odata.deltaLink in the previous response. Any resources marked as deleted should be removed from your local state.
// returns a *RequestInformation when successful
func (m *ItemSitesDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSitesDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSitesDeltaRequestBuilder when successful
func (m *ItemSitesDeltaRequestBuilder) WithUrl(rawUrl string)(*ItemSitesDeltaRequestBuilder) {
    return NewItemSitesDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
